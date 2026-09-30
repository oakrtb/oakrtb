package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.BidResponse;
import com.oakrtb.openrtb.v2.MarkupType;
import com.oakrtb.sdk.builder.BidResponseBuilder;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ResponseViewTest {

  @Test
  void markupFromMtypeMapsKnownTypes() {
    assertTrue(MarkupMask.fromMtype(MarkupType.MARKUP_TYPE_BANNER).hasBanner());
    assertTrue(MarkupMask.fromMtype(MarkupType.MARKUP_TYPE_VIDEO).hasVideo());
    assertTrue(MarkupMask.fromMtype(MarkupType.MARKUP_TYPE_AUDIO).hasAudio());
    assertTrue(MarkupMask.fromMtype(MarkupType.MARKUP_TYPE_NATIVE).hasNative());
    assertEquals(0, MarkupMask.fromMtype((MarkupType) null).bits());
  }

  @Test
  void markupFromMtypeValue() {
    assertTrue(MarkupMask.fromMtype(1).hasBanner());
    assertTrue(MarkupMask.fromMtype(4).hasNative());
    assertEquals(0, MarkupMask.fromMtype(0).bits());
    assertEquals(0, MarkupMask.fromMtype(99).bits());
  }

  @Test
  void sharedNoBid() {
    BidResponse res = BidResponseBuilder.create("req-1").noBid(3).build();
    ResponseView sv = ResponseView.of(res);
    assertEquals("req-1", sv.requestId());
    assertTrue(sv.noBid());
    assertEquals(3, sv.nbr());
  }

  @Test
  void pipelineBidsFlattensSeatBids() {
    BidResponse res =
        BidResponseBuilder.create("req-1")
            .addSeatBid(
                "seat-a",
                BidResponseBuilder.bid("b1", "1", 1.0).banner().build(),
                BidResponseBuilder.bid("b2", "2", 2.0).video().build())
            .build();
    ResponseView snap = ResponseView.of(res);
    var bids = snap.bids();
    assertEquals(2, bids.size());
    assertEquals("b1", bids.get(0).id());
    assertEquals("seat-a", bids.get(0).seat());
    assertTrue(bids.get(0).markup().hasBanner());
    assertTrue(bids.get(1).markup().hasVideo());
  }
  @Test
  void projectionsShareBidsAndExposeImmutableLists() {
    var response = BidResponseBuilder.create("a").addSeatBid("seat",
        BidResponseBuilder.bid("b", "1", 1).banner().build()).build();
    var view = ResponseView.of(response);
    assertSame(response, view.response());
    assertSame(view.bids().getFirst(), view.seatBids().getFirst().bids().getFirst());
    assertThrows(UnsupportedOperationException.class, () -> view.bids().clear());
    assertThrows(UnsupportedOperationException.class, () -> view.seatBids().clear());
    assertThrows(UnsupportedOperationException.class, () -> view.seatBids().getFirst().bids().clear());
  }
}
