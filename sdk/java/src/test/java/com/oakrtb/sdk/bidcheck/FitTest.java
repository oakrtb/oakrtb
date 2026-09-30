package com.oakrtb.sdk.bidcheck;
import com.oakrtb.sdk.validation.*;

import com.oakrtb.openrtb.v2.Banner;
import com.oakrtb.openrtb.v2.Bid;
import com.oakrtb.openrtb.v2.BidRequest;
import com.oakrtb.openrtb.v2.BidResponse;
import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.MarkupType;
import com.oakrtb.openrtb.v2.Video;
import com.oakrtb.sdk.builder.BidRequestBuilder;
import com.oakrtb.sdk.builder.BidResponseBuilder;
import com.oakrtb.sdk.builder.ImpBuilders;
import com.oakrtb.sdk.builder.Parts;
import com.oakrtb.sdk.view.MarkupMask;
import com.oakrtb.sdk.view.RequestView;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FitTest {
  private static RequestView bannerSnap() {
    BidRequest req =
        BidRequestBuilder.create("a1")
            .firstPrice()
            .currency("USD")
            .site(Parts.site().id("s1").build())
            .addImp(ImpBuilders.banner("1").size(300, 250).floor(1.0, "USD").build())
            .build();
    return RequestView.of(req);
  }

  @Test
  void impReadyVideoMissingMimes() {
    Imp imp =
        Imp.newBuilder()
            .setId("1")
            .setVideo(Video.newBuilder().build())
            .build();
    BidRequest req =
        BidRequest.newBuilder().setId("x").setAt(1).addCur("USD").addImp(imp).build();
    var snap = RequestView.of(req);
    CheckResult r = com.oakrtb.sdk.validation.Readiness.impReady(snap.imps().get(0).imp(), 2);
    assertFalse(r.ok());
    assertTrue(r.has(IssueCode.VIDEO_MIMES_MISSING));
  }

  @Test
  void multiFormatNeedsMtype() {
    Imp multi =
        Imp.newBuilder()
            .setId("1")
            .setBanner(Banner.newBuilder().setW(320).setH(50).build())
            .setVideo(Video.newBuilder().addMimes("video/mp4").build())
            .build();
    BidRequest req =
        BidRequest.newBuilder().setId("m").setAt(1).addCur("USD").addImp(multi).build();
    var snap = RequestView.of(req);
    Bid bid = Bid.newBuilder().setId("b1").setImpid("1").setPrice(2.0).build();
    CheckResult r = BidCheck.bid(snap, bid);
    assertFalse(r.ok());
    assertTrue(r.has(IssueCode.MTYPE_REQUIRED));
  }

  @Test
  void mtypeMismatch() {
    var snap = bannerSnap();
    Bid bid =
        Bid.newBuilder()
            .setId("b1")
            .setImpid("1")
            .setPrice(2.0)
            .setMtype(MarkupType.MARKUP_TYPE_VIDEO)
            .build();
    CheckResult r = BidCheck.bid(snap, bid);
    assertFalse(r.ok());
    assertTrue(r.has(IssueCode.MTYPE_MISMATCH));
  }

  @Test
  void priceBelowFloorIsWarn() {
    var snap = bannerSnap();
    Bid bid = BidResponseBuilder.bid("b1", "1", 0.5).banner().build();
    assertFalse(BidCheck.bid(snap, bid).has(IssueCode.PRICE_BELOW_FLOOR));
    BidResponse res =
        BidResponseBuilder.create("a1").currency("USD").addSeatBid("s1", bid).build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
    assertTrue(r.has(IssueCode.PRICE_BELOW_FLOOR));
  }

  @Test
  void attrBlockedIsWarn() {
    Imp imp =
        Imp.newBuilder()
            .setId("1")
            .setBanner(Banner.newBuilder().setW(300).setH(250).addBattr(1).build())
            .build();
    BidRequest req =
        BidRequest.newBuilder().setId("x").setAt(1).addCur("USD").addImp(imp).build();
    var snap = RequestView.of(req);
    Bid bid =
        Bid.newBuilder()
            .setId("b1")
            .setImpid("1")
            .setPrice(2.0)
            .setMtype(MarkupType.MARKUP_TYPE_BANNER)
            .addAttr(1)
            .build();
    CheckResult r = BidCheck.bid(snap, bid);
    assertTrue(r.ok());
    assertTrue(r.has(IssueCode.ATTR_BLOCKED));
  }

  @Test
  void noBidResponseOk() {
    var snap = bannerSnap();
    BidResponse res = BidResponseBuilder.create("a1").currency("USD").noBid(0).build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
  }

  @Test
  void impidNotFound() {
    var snap = bannerSnap();
    Bid bid = BidResponseBuilder.bid("b1", "missing", 2.0).banner().build();
    CheckResult r = BidCheck.bid(snap, bid);
    assertFalse(r.ok());
    assertTrue(r.has(IssueCode.IMP_NOT_FOUND));
  }

  @Test
  void impReadyFormatBit() {
    var snap = bannerSnap();
    CheckResult r = com.oakrtb.sdk.validation.Readiness.impReadyMask(snap.imps().get(0).imp(), MarkupMask.BANNER);
    assertTrue(r.ok());
  }

  @Test
  void floorCurDiffSkipsCompare() {
    Imp imp =
        Imp.newBuilder()
            .setId("1")
            .setBidfloor(1.0)
            .setBidfloorcur("EUR")
            .setBanner(Banner.newBuilder().setW(1).setH(1).build())
            .build();
    BidRequest req =
        BidRequest.newBuilder().setId("a1").setAt(1).addCur("USD").addImp(imp).build();
    var snap = RequestView.of(req);
    Bid bid = BidResponseBuilder.bid("b1", "1", 0.1).banner().build();
    BidResponse res =
        BidResponseBuilder.create("a1").currency("USD").addSeatBid("s", bid).build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
    assertTrue(r.has(IssueCode.FLOOR_CUR_DIFF));
    assertFalse(r.has(IssueCode.PRICE_BELOW_FLOOR));
  }

  @Test
  void responseCurWhitespaceTrimmed() {
    var snap = bannerSnap();
    Bid bid = BidResponseBuilder.bid("b1", "1", 0.5).banner().build();
    BidResponse res =
        BidResponse.newBuilder()
            .setId("a1")
            .setCur("  USD  ")
            .addSeatbid(
                com.oakrtb.openrtb.v2.SeatBid.newBuilder().setSeat("s").addBid(bid).build())
            .build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
    assertFalse(r.has(IssueCode.CUR_NOT_ALLOWED));
    assertTrue(r.has(IssueCode.PRICE_BELOW_FLOOR));
  }

  @Test
  void responseEmptyBidMalformed() {
    var snap = bannerSnap();
    BidResponse res =
        BidResponse.newBuilder()
            .setId("a1")
            .setCur("USD")
            .addSeatbid(com.oakrtb.openrtb.v2.SeatBid.newBuilder().setSeat("s").build())
            .build();
    CheckResult r = BidCheck.response(snap, res);
    assertFalse(r.ok());
    assertTrue(r.has(IssueCode.MALFORMED));
  }

  @Test
  void responseBlankCurSkipsCurAndFloor() {
    var snap = bannerSnap();
    Bid bid = BidResponseBuilder.bid("b1", "1", 0.5).banner().build();
    BidResponse res =
        BidResponse.newBuilder()
            .setId("a1")
            .setCur("   ")
            .addSeatbid(
                com.oakrtb.openrtb.v2.SeatBid.newBuilder().setSeat("s").addBid(bid).build())
            .build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
    assertFalse(r.has(IssueCode.CUR_NOT_ALLOWED));
    assertFalse(r.has(IssueCode.PRICE_BELOW_FLOOR));
  }

  @Test
  void curNotAllowed() {
    var snap = bannerSnap();
    Bid bid = BidResponseBuilder.bid("b1", "1", 2.0).banner().build();
    BidResponse res =
        BidResponseBuilder.create("a1").currency("EUR").addSeatBid("s", bid).build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
    assertTrue(r.has(IssueCode.CUR_NOT_ALLOWED));
  }

  @Test
  void responseHappyPath() {
    var snap = bannerSnap();
    BidResponse res =
        BidResponseBuilder.create("a1")
            .currency("USD")
            .addSeatBid("", BidResponseBuilder.bid("b1", "1", 2.0).banner().adm("<a/>").build())
            .build();
    CheckResult r = BidCheck.response(snap, res);
    assertTrue(r.ok());
  }

  @Test void requestIdMustMatchIncludingNoBid() {
    var req = bannerSnap();
    for (String id : new String[]{"a1", "other", ""}) {
      for (boolean noBid : new boolean[]{true, false}) {
        var res = BidResponse.newBuilder().setId(id).setCur("USD");
        if (!noBid) res.addSeatbid(com.oakrtb.openrtb.v2.SeatBid.newBuilder()
            .addBid(Bid.newBuilder().setId("b").setImpid("1").setPrice(2)));
        var result = BidCheck.response(req, res.build());
        assertEquals(!id.equals("a1"), result.has(IssueCode.REQUEST_ID_MISMATCH));
        if (!id.equals("a1")) assertFalse(result.ok());
      }
    }
  }
  @Test void rejectsInvalidPrices() {
    var req = bannerSnap();
    for (double price : new double[]{Double.NaN, Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY, 0, -1}) {
      var bid = Bid.newBuilder().setId("b").setImpid("1").setPrice(price).build();
      var res = BidResponse.newBuilder().setId("a1").setCur("USD")
          .addSeatbid(com.oakrtb.openrtb.v2.SeatBid.newBuilder().addBid(bid)).build();
      assertThrows(IllegalArgumentException.class, () -> com.oakrtb.sdk.view.ResponseView.of(res));
      assertFalse(BidCheck.bid(req, bid).ok());
      assertFalse(BidCheck.response(req, res).ok());
      assertThrows(IllegalStateException.class, () -> BidResponseBuilder.create("a1").addSeatBid("s", bid).build());
    }
  }
}
