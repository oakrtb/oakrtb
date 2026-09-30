package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.BidResponse;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.Optional;

/** Query view of an immutable generated model; no pipeline state. */
public final class ResponseView {
    private final BidResponse response;
    private final List<SeatBidView> seatBids;
    private final List<BidView> bids;

    private ResponseView(BidResponse response) {
      this.response = response;
      this.seatBids = response.getSeatbidList().stream().map(SeatBidView::of).toList();
      this.bids = seatBids.stream().flatMap(seat -> seat.bids().stream()).toList();
    }

    /** Validate once and construct seat and flat lists sharing the same BidView objects. */
    public static ResponseView of(BidResponse response) {
      com.oakrtb.sdk.validation.BasicValidation.validateResponse(response);
      return new ResponseView(response);
    }

    public List<SeatBidView> seatBids() { return seatBids; }
    public List<BidView> bids() { return bids; }

    /** Echoed request ID. */
    public String requestId() {
      return response.getId();
    }

    /** Response-level bid ID. */
    public String bidid() {
      return response.getBidid();
    }

    /** Response currency. */
    public String currency() {
      return response.getCur();
    }

    /** Reports whether this is a wire no-bid (no seatbid). */
    public boolean noBid() {
      return response.getSeatbidCount() == 0;
    }

    /** No-bid reason code. */
    public int nbr() {
      return response.getNbr();
    }

    /** Original BidResponse. */
    public BidResponse response() {
      return response;
    }

    /**
     * Finds a bid by ID.
     *
     * @param id bid id
     * @return the BidView, or empty if absent
     */
    public Optional<BidView> findBid(String id) {
      if (id == null) {
        return Optional.empty();
      }
      return bids.stream().filter(b -> id.equals(b.id())).findFirst();
    }

    /**
     * Filters bids by impid.
     *
     * @param impid Imp id
     * @return matching Bids
     */
    public List<BidView> bidsForImp(String impid) {
      List<BidView> out = new ArrayList<>();
      for (BidView b : bids) {
        if (Objects.equals(impid, b.impid())) {
          out.add(b);
        }
      }
      return Collections.unmodifiableList(out);
    }

    /**
     * Filters bids by the specified format bit.
     *
     * @param formatFlag single-bit {@link MarkupMask} constant
     * @return matching Bids
     */
    public List<BidView> bidsWith(int formatFlag) {
      List<BidView> out = new ArrayList<>();
      for (BidView b : bids) {
        if (b.markup().has(formatFlag)) {
          out.add(b);
        }
      }
      return Collections.unmodifiableList(out);
    }

    /**
     * Extracts compact per-Bid facts.
     *
     * @return the BidFact list
     */
    public List<BidFact> facts() {
      List<BidFact> out = new ArrayList<>(bids.size());
      for (BidView b : bids) {
        out.add(BidFact.from(b));
      }
      return Collections.unmodifiableList(out);
    }
    public record BidFact(
      String id,
      String impid,
      String seat,
      double price,
      MarkupMask markup,
      int mtype,
      String crid,
      String dealid,
      int w,
      int h,
      int dur,
      boolean hasAdm,
      List<String> adomain) {

    public BidFact { adomain = List.copyOf(adomain); }

    static BidFact from(BidView b) {
      return new BidFact(
          b.id(),
          b.impid(),
          b.seat(),
          b.price(),
          b.markup(),
          b.mtype(),
          b.crid(),
          b.dealid(),
          b.w(),
          b.h(),
          b.dur(),
          b.adm() != null && !b.adm().isEmpty(),
          b.adomain());
    }

    /** Reports whether this is a Banner bid. */
    public boolean hasBanner() {
      return markup.hasBanner();
    }

    /** Reports whether this is a Video bid. */
    public boolean hasVideo() {
      return markup.hasVideo();
    }

    /** Reports whether this is an Audio bid. */
    public boolean hasAudio() {
      return markup.hasAudio();
    }

    /** Reports whether this is a Native bid. */
    public boolean hasNative() {
      return markup.hasNative();
    }
  }
}
