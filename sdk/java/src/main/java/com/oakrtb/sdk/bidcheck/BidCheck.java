package com.oakrtb.sdk.bidcheck;

import com.oakrtb.sdk.validation.CheckIssue;
import com.oakrtb.sdk.validation.CheckResult;
import com.oakrtb.sdk.validation.IssueCode;
import com.oakrtb.sdk.validation.Severity;

import com.oakrtb.openrtb.v2.Bid;
import com.oakrtb.openrtb.v2.BidResponse;
import com.oakrtb.openrtb.v2.SeatBid;
import com.oakrtb.sdk.view.MarkupMask;
import com.oakrtb.sdk.view.ImpView;
import com.oakrtb.sdk.view.RequestView;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.Objects;
import java.util.Set;

/**
 * Optional bid-to-request consistency checks between RequestView and BidResponse.
 *
 * <p>These are neither JSON Schema nor basic validation checks; findings are returned in {@link CheckResult}. {@code null} arguments throw NPE; business mismatches do not throw exceptions.
 * {@link CheckResult#ok()} only indicates the absence of ERROR findings; WARN findings (low prices, blocked ads, deadlines, etc.) may still yield ok.
 */
public final class BidCheck {
  private BidCheck() {}


  /**
   * Checks one Bid against a request snapshot. Without response currency or seat context, skips price comparisons and deal seat checks; use {@link #response} for full checks.
   *
   * @param req request view snapshot
   * @param bid bid to check
   * @return bid check results
   */
  public static CheckResult bid(RequestView req, Bid bid) {
    return bid(req, bid, null, null, "bid", MatchContext.from(req));
  }

  /**
   * Checks an entire BidResponse. An empty seatbid (wire no-bid) still undergoes request ID checks;
   * adds response-level PAST_DEADLINE / CUR_NOT_ALLOWED findings when applicable.
   *
   * @param req request view snapshot
   * @param res bid response
   * @return bid check results
   */
  public static CheckResult response(RequestView req, BidResponse res) {
    Objects.requireNonNull(req, "req");
    Objects.requireNonNull(res, "res");
    List<CheckIssue> issues = new ArrayList<>();
    if (!res.getId().equals(req.auctionId())) {
      issues.add(new CheckIssue(IssueCode.REQUEST_ID_MISMATCH, Severity.ERROR, "BidResponse.id", "response id does not match request id"));
    }
    if (req.pastDeadline()) {
      issues.add(
          new CheckIssue(
              IssueCode.PAST_DEADLINE,
              Severity.WARN,
              "tmax",
              "past request deadline (85% of tmax)"));
    }
    if (res.getSeatbidCount() == 0) {
      return CheckResult.of(issues);
    }
    // Same blank rule as basic validation / checkFloor: strip; blank skips CUR + floor.
    String resCur = res.getCur() == null ? "" : res.getCur().strip();
    if (!resCur.isEmpty()) {
      boolean allowed = false;
      for (String c : req.currencies()) {
        String cc = c == null ? "" : c.strip();
        if (!cc.isEmpty() && asciiLower(resCur).equals(asciiLower(cc))) {
          allowed = true;
          break;
        }
      }
      if (!allowed) {
        issues.add(
            new CheckIssue(
                IssueCode.CUR_NOT_ALLOWED,
                Severity.WARN,
                "BidResponse.cur",
                "response currency not in BidRequest.cur"));
      }
    }
    MatchContext context = MatchContext.from(req);
    for (int i = 0; i < res.getSeatbidCount(); i++) {
      SeatBid sb = res.getSeatbid(i);
      if (sb.getBidCount() == 0) {
        issues.add(
            new CheckIssue(
                IssueCode.MALFORMED,
                Severity.ERROR,
                "seatbid[" + i + "].bid",
                "seatbid.bid must be a non-empty array"));
        continue;
      }
      for (int j = 0; j < sb.getBidCount(); j++) {
        String path = "seatbid[" + i + "].bid[" + j + "]";
        CheckResult one = bid(req, sb.getBid(j), resCur, sb.getSeat(), path, context);
        issues.addAll(one.issues());
      }
    }
    return CheckResult.of(issues);
  }

  private static CheckResult bid(
      RequestView req, Bid bid, String responseCur, String seat, String path, MatchContext context) {
    Objects.requireNonNull(req, "req");
    Objects.requireNonNull(bid, "bid");
    List<CheckIssue> issues = new ArrayList<>();
    if (!Double.isFinite(bid.getPrice()) || bid.getPrice() <= 0) {
      return CheckResult.of(List.of(new CheckIssue(IssueCode.MALFORMED, Severity.ERROR, path + ".price", "price must be finite and > 0")));
    }
    ImpView imp = context.imps().get(bid.getImpid());
    if (imp == null) {
      issues.add(
          new CheckIssue(
              IssueCode.IMP_NOT_FOUND,
              Severity.ERROR,
              path + ".impid",
              "impid not found in request"));
      return CheckResult.of(issues);
    }
    int mtype = bid.getMtypeValue();
    MarkupMask markup = imp.markup();

    if (mtype >= 500) {
      issues.add(
          new CheckIssue(
              IssueCode.MTYPE_VENDOR,
              Severity.WARN,
              path + ".mtype",
              "vendor mtype >=500"));
    } else if (mtype != 0 && (mtype < 1 || mtype > 4)) {
      issues.add(
          new CheckIssue(
              IssueCode.MTYPE_UNKNOWN,
              Severity.ERROR,
              path + ".mtype",
              "mtype must be 0–4 or >=500"));
    } else if (mtype == 0 && markup.count() > 1) {
      issues.add(
          new CheckIssue(
              IssueCode.MTYPE_REQUIRED,
              Severity.ERROR,
              path + ".mtype",
              "multi-format Imp requires Bid.mtype"));
    } else if (mtype >= 1 && mtype <= 4) {
      MarkupMask chosen = MarkupMask.fromMtype(mtype);
      if (!markup.has(chosen.bits())) {
        issues.add(
            new CheckIssue(
                IssueCode.MTYPE_MISMATCH,
                Severity.ERROR,
                path + ".mtype",
                "mtype does not match Imp formats"));
      }
    }

    // Effective format for blocklist checks
    MarkupMask eff =
        (mtype >= 1 && mtype <= 4)
            ? MarkupMask.fromMtype(mtype)
            : (markup.count() == 1 ? markup : MarkupMask.of(MarkupMask.NONE));

    checkFloor(issues, imp, bid, responseCur, seat, path);
    checkAttr(issues, imp, bid, eff, path);
    checkBlocks(issues, context, bid, path);
    return CheckResult.of(issues);
  }

  private static void checkFloor(
      List<CheckIssue> issues,
      ImpView imp,
      Bid bid,
      String responseCur,
      String seat,
      String path) {
    double floor = imp.bidFloor();
    String currency = imp.bidFloorCur();
    boolean fixed = false;
    if (bid.getDealid().isEmpty()) {
      if (imp.pmp() != null && imp.pmp().getPrivateAuction() == 1) {
        issues.add(new CheckIssue(IssueCode.DEAL_REQUIRED, Severity.ERROR, path + ".dealid", "private auction requires a deal")); return;
      }
    } else {
      com.oakrtb.openrtb.v2.Deal deal = null;
      if (imp.pmp() != null) for (var candidate : imp.pmp().getDealsList()) if (candidate.getId().equals(bid.getDealid())) { deal = candidate; break; }
      if (deal == null) { issues.add(new CheckIssue(IssueCode.DEAL_NOT_FOUND, Severity.ERROR, path + ".dealid", "dealid not found on impression")); return; }
      if (seat != null && deal.getWseatCount() > 0 && !deal.getWseatList().contains(seat)) {
        issues.add(new CheckIssue(IssueCode.DEAL_SEAT_NOT_ALLOWED, Severity.ERROR, path + ".dealid", "seat not allowed by deal")); return;
      }
      if (deal.getWadomainCount() > 0) {
        var allowed = toLowerSet(deal.getWadomainList());
        if (bid.getAdomainCount() == 0 || bid.getAdomainList().stream().anyMatch(d -> !allowed.contains(asciiLower(d)))) {
          issues.add(new CheckIssue(IssueCode.DEAL_ADOMAIN_NOT_ALLOWED, Severity.ERROR, path + ".adomain", "advertiser domain not allowed by deal")); return;
        }
      }
      if (deal.hasBidfloor()) { floor = deal.getBidfloor(); currency = deal.getBidfloorcur(); }
      fixed = deal.getAt() == 3 && deal.hasBidfloor();
    }
    if (floor <= 0 && !fixed) {
      return;
    }
    // Require both response cur and bidfloorcur before numeric compare (BidCheck.bid has no response cur).
    String respCur = responseCur == null ? "" : responseCur.strip();
    String floorCur = currency == null ? "" : currency.strip();
    if (respCur.isEmpty() || floorCur.isEmpty()) {
      return;
    }
    if (!asciiLower(respCur).equals(asciiLower(floorCur))) {
      issues.add(
          new CheckIssue(
              IssueCode.FLOOR_CUR_DIFF,
              Severity.WARN,
              path + ".price",
              "bid currency differs from applicable bidfloorcur; skip floor compare"));
      return;
    }
    if (fixed && bid.getPrice() != floor) { issues.add(new CheckIssue(IssueCode.DEAL_PRICE_MISMATCH, Severity.ERROR, path + ".price", "price differs from fixed deal price")); return; }
    if (bid.getPrice() < floor) {
      issues.add(
          new CheckIssue(
              IssueCode.PRICE_BELOW_FLOOR,
              Severity.WARN,
              path + ".price",
              "price below applicable bidfloor"));
    }
  }

  private static void checkAttr(
      List<CheckIssue> issues,
      ImpView imp,
      Bid bid,
      MarkupMask eff,
      String path) {
    if (bid.getAttrCount() == 0 || eff.bits() == MarkupMask.NONE) {
      return;
    }
    List<Integer> battr = battrFor(imp, eff);
    if (battr.isEmpty()) {
      return;
    }
    Set<Integer> blocked = new HashSet<>(battr);
    for (int a : bid.getAttrList()) {
      if (blocked.contains(a)) {
        issues.add(
            new CheckIssue(
                IssueCode.ATTR_BLOCKED,
                Severity.WARN,
                path + ".attr",
                "bid.attr intersects format battr"));
        return;
      }
    }
  }

  private static List<Integer> battrFor(ImpView imp, MarkupMask eff) {
    if (eff.hasBanner() && imp.banner() != null) {
      return imp.banner().getBattrList();
    }
    if (eff.hasVideo() && imp.video() != null) {
      return imp.video().getBattrList();
    }
    if (eff.hasAudio() && imp.audio() != null) {
      return imp.audio().getBattrList();
    }
    if (eff.hasNative() && imp.nativeAd() != null) {
      return imp.nativeAd().getBattrList();
    }
    return List.of();
  }

  private static void checkBlocks(
      List<CheckIssue> issues, MatchContext context, Bid bid, String path) {
    Set<String> badv = context.badv();
    if (!badv.isEmpty()) {
      Set<String> block = badv;
      for (String d : bid.getAdomainList()) {
        if (d != null && block.contains(asciiLower(d))) {
          issues.add(
              new CheckIssue(
                  IssueCode.ADOMAIN_BLOCKED,
                  Severity.WARN,
                  path + ".adomain",
                  "adomain hit BidRequest.badv"));
          break;
        }
      }
    }
    Set<String> bapp = context.bapp();
    if (!bapp.isEmpty() && bid.getBundle() != null && !bid.getBundle().isBlank()) {
      Set<String> block = bapp;
      String bundle = asciiLower(bid.getBundle().strip());
      if (!bundle.isEmpty() && block.contains(bundle)) {
        issues.add(
            new CheckIssue(
                IssueCode.BUNDLE_BLOCKED,
                Severity.WARN,
                path + ".bundle",
                "bundle hit BidRequest.bapp"));
      }
    }
    Set<String> bcat = context.bcat();
    if (!bcat.isEmpty()) {
      Set<String> block = bcat;
      for (String c : bid.getCatList()) {
        if (c != null && block.contains(asciiLower(c))) {
          issues.add(
              new CheckIssue(
                  IssueCode.CAT_BLOCKED,
                  Severity.WARN,
                  path + ".cat",
                  "cat hit BidRequest.bcat"));
          break;
        }
      }
    }
  }

  private static Set<String> toLowerSet(List<String> in) {
    Set<String> out = new HashSet<>();
    for (String s : in) {
      if (s != null) {
        out.add(asciiLower(s));
      }
    }
    return out;
  }
  private record MatchContext(Map<String, ImpView> imps,
      Set<String> badv, Set<String> bapp, Set<String> bcat) {
    static MatchContext from(RequestView req) {
      Objects.requireNonNull(req, "req");
      Map<String, ImpView> imps = new HashMap<>();
      for (var imp : req.imps()) imps.putIfAbsent(imp.id(), imp);
      return new MatchContext(imps, toLowerSet(req.request().getBadvList()),
          toLowerSet(req.request().getBappList()), toLowerSet(req.request().getBcatList()));
    }
  }
  private static String asciiLower(String value) {
    StringBuilder out = new StringBuilder(value.length());
    for (int i = 0; i < value.length(); i++) {
      char c = value.charAt(i);
      out.append(c >= 'A' && c <= 'Z' ? (char) (c + ('a' - 'A')) : c);
    }
    return out.toString();
  }
}
