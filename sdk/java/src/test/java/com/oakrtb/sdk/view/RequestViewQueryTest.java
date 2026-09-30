package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.BidRequest;
import com.oakrtb.sdk.builder.BidRequestBuilder;
import com.oakrtb.sdk.builder.ImpBuilders;
import com.oakrtb.sdk.builder.Parts;
import java.time.Instant;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class RequestViewQueryTest {
  @Test
  void runFacts() {
    BidRequest req =
        BidRequestBuilder.create("auction-1")
            .firstPrice()
            .tmax(120)
            .currency("USD")
            .site(Parts.site().id("s1").domain("example.com").page("https://example.com/a").build())
            .addImp(ImpBuilders.banner("1").size(300, 250).floor(0.03, "USD").secure().build())
            .build();

    RequestView snap = RequestView.of(req);
    assertEquals("auction-1", snap.auctionId());
    assertEquals(Inventory.SITE, snap.inventory());
    assertTrue(snap.findImp("1").isPresent());
    assertEquals(1, snap.impsWith(MarkupMask.BANNER).size());

    var facts = snap.facts();
    assertEquals(1, facts.size());
    RequestView.ImpFact f = facts.get(0);
    assertTrue(f.hasBanner());
    assertEquals(1, f.mtype());
    assertEquals(300, f.bannerW());
    assertEquals(0.03, f.bidFloor(), 1e-9);
    assertEquals(1, f.secure());
  }

  @Test
  void factoryValidatesStructure() {
    BidRequest req =
        BidRequestBuilder.create("x")
            .firstPrice()
            .currency("USD")
            .addImp(ImpBuilders.banner("1").size(1, 1).build())
            .build();

    RequestView snap = RequestView.of(req);
    assertThrows(IllegalArgumentException.class, () -> RequestView.of(req.toBuilder().clearId().build()));
    assertEquals("x", snap.auctionId());
  }

  @Test
  void sharedPinOnce() throws Exception {
    BidRequest req =
        BidRequestBuilder.create("x")
            .firstPrice()
            .tmax(200)
            .currency("USD")
            .addImp(ImpBuilders.banner("1").size(1, 1).build())
            .build();
    RequestView p = RequestView.of(req);
    Instant d1 = p.deadline();
    Thread.sleep(3);
    p.facts();
    Instant d2 = p.deadline();
    assertEquals(d1, d2);
    assertSame(req, p.request());
    assertThrows(UnsupportedOperationException.class, () -> p.imps().clear());
  }
}
