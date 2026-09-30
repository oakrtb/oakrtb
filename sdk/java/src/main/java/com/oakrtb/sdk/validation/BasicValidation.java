package com.oakrtb.sdk.validation;
import com.oakrtb.openrtb.v2.*;
import java.util.Objects;
/** Basic model checks, independent of views, builders and JSON Schema. */
public final class BasicValidation {
 private BasicValidation() {}
  public static void validateRequest(BidRequest req) {
    Objects.requireNonNull(req, "BidRequest");
    if (req.getId().isBlank()) {
      throw new IllegalArgumentException("validation: BidRequest.id is required");
    }
    if (req.getAt() != 1 && req.getAt() != 2 && req.getAt() < 500) {
      throw new IllegalArgumentException("validation: BidRequest.at must be 1, 2, or >=500");
    }
    if (req.getCurCount() == 0) {
      throw new IllegalArgumentException("validation: BidRequest.cur is required");
    }
    for (int i = 0; i < req.getCurCount(); i++) {
      if (req.getCur(i).isBlank()) {
        throw new IllegalArgumentException("validation: BidRequest.cur[" + i + "] is blank");
      }
    }
    if (req.getImpCount() == 0) {
      throw new IllegalArgumentException("validation: BidRequest.imp requires at least one Imp");
    }
    int n = 0;
    if (req.hasSite()) {
      n++;
    }
    if (req.hasApp()) {
      n++;
    }
    if (req.hasDooh()) {
      n++;
    }
    if (n > 1) {
      throw new IllegalArgumentException("validation: site/app/dooh are mutually exclusive");
    }
    var ids = new java.util.HashSet<String>();
    for (int i = 0; i < req.getImpCount(); i++) {
      Imp imp = req.getImp(i);
      if (imp.getId().isBlank()) {
        throw new IllegalArgumentException("validation: imp[" + i + "].id is required");
      }
      if (!ids.add(imp.getId())) throw new IllegalArgumentException("validation: duplicate imp.id " + imp.getId());
      if (imp.hasPmp()) {
        var deals = new java.util.HashSet<String>();
        for (var d : imp.getPmp().getDealsList()) {
          if (d.getId().isBlank() || !deals.add(d.getId())) throw new IllegalArgumentException("validation: deal.id must be nonblank and unique");
          if (d.hasBidfloor() && (!Double.isFinite(d.getBidfloor()) || d.getBidfloor() < 0)) throw new IllegalArgumentException("validation: deal.bidfloor must be finite and nonnegative");
          if (d.getAt() == 3 && !d.hasBidfloor()) throw new IllegalArgumentException("validation: fixed-price deal requires bidfloor");
        }
      }
      if (!imp.hasBanner() && !imp.hasVideo() && !imp.hasAudio() && !imp.hasNative()) {
        throw new IllegalArgumentException(
            "validation: imp[" + i + "] needs banner, video, audio, or native");
      }
    }
  }
  public static void validateResponse(BidResponse res) {
    Objects.requireNonNull(res, "BidResponse");
    if (res.getId().isBlank()) {
      throw new IllegalArgumentException("validation: BidResponse.id is required");
    }
    if (res.getCur().isBlank()) {
      throw new IllegalArgumentException("validation: BidResponse.cur is required");
    }
    if (res.getSeatbidCount() == 0) {
      return;
    }
    for (int i = 0; i < res.getSeatbidCount(); i++) {
      SeatBid sb = res.getSeatbid(i);
      if (sb.getBidCount() == 0) {
        throw new IllegalArgumentException("validation: seatbid[" + i + "] needs at least one bid");
      }
      for (int j = 0; j < sb.getBidCount(); j++) {
        Bid bid = sb.getBid(j);
        if (bid.getId().isBlank() || bid.getImpid().isBlank()) {
          throw new IllegalArgumentException(
              "validation: seatbid[" + i + "].bid[" + j + "] requires id and impid");
        }
        if (!Double.isFinite(bid.getPrice()) || bid.getPrice() <= 0) {
          throw new IllegalArgumentException(
              "validation: seatbid[" + i + "].bid[" + j + "].price must be finite and > 0");
        }
      }
    }
  }
}
