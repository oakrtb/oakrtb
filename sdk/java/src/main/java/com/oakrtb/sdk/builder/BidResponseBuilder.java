package com.oakrtb.sdk.builder;
import com.oakrtb.sdk.codec.Json;

import com.oakrtb.openrtb.v2.Bid;
import com.oakrtb.openrtb.v2.BidResponse;
import com.oakrtb.openrtb.v2.MarkupType;
import com.oakrtb.openrtb.v2.SeatBid;
import com.oakrtb.sdk.jsonschema.Report;
import com.oakrtb.sdk.jsonschema.Schema;

import java.util.Objects;

/**
 * Fluent builder for an OpenRTB {@link BidResponse}.
 *
 * <p>For DSP bid responses or structured no-bids after an auction; {@link #build()} validates required fields before returning.
 */
public final class BidResponseBuilder {
  private final BidResponse.Builder res = BidResponse.newBuilder().setCur("USD");
  private boolean noBidSet;
  private String error;

  private BidResponseBuilder(String requestId) {
    if (requestId == null || requestId.isBlank()) {
      error = "builder: BidResponse.id is required (echo BidRequest.id)";
    } else {
      res.setId(requestId);
    }
  }

  /**
   * Creates a response builder; {@code id} must echo the corresponding {@link BidRequest} ID.
   *
   * @param requestId auction request ID (required and nonempty)
   * @return a new builder instance
   */
  public static BidResponseBuilder create(String requestId) {
    return new BidResponseBuilder(requestId);
  }

  /**
   * Sets the optional response-level bid ID for logging/reconciliation.
   *
   * @param bidid bid ID
   * @return this builder for chaining
   */
  public BidResponseBuilder bidId(String bidid) {
    res.setBidid(bidid);
    return this;
  }

  /**
   * Sets the response currency (ISO-4217, e.g. {@code USD}).
   *
   * @param cur currency code
   * @return this builder for chaining
   */
  public BidResponseBuilder currency(String cur) {
    res.setCur(cur);
    return this;
  }

  /**
   * Declares a structured no-bid and clears existing {@code seatbid} entries.
   *
   * @param nbr no-bid reason code (OpenRTB enum value)
   * @return this builder for chaining
   */
  public BidResponseBuilder noBid(int nbr) {
    res.setNbr(nbr);
    res.clearSeatbid();
    noBidSet = true;
    return this;
  }

  /**
   * Appends a {@link SeatBid} containing one or more {@link Bid} entries.
   *
   * <p>Clears the no-bid state; {@code bids} must contain at least one entry.
   *
   * @param seat seat identifier
   * @param bids bids for this seat
   * @return this builder for chaining
   * @throws IllegalArgumentException if {@code bids} is empty
   */
  public BidResponseBuilder addSeatBid(String seat, Bid... bids) {
    if (bids == null || bids.length == 0) {
      throw new IllegalArgumentException("builder: SeatBid requires at least one Bid");
    }
    // Switching to a bid clears structured no-bid.
    noBidSet = false;
    res.clearNbr();
    SeatBid.Builder sb = SeatBid.newBuilder().setSeat(seat);
    for (Bid bid : bids) {
      sb.addBid(Objects.requireNonNull(bid));
    }
    res.addSeatbid(sb);
    return this;
  }

  /**
   * Builds a protobuf {@link BidResponse} and validates id, cur, seatbid/no-bid, and each bid's id/impid/price.
   *
   * @return the constructed response
   * @throws IllegalStateException if required fields are missing or a bid is invalid
   */
  public BidResponse build() {
    if (error != null) throw new IllegalStateException(error);
    BidResponse result = res.build();
    try { com.oakrtb.sdk.validation.BasicValidation.validateResponse(result); }
    catch (IllegalArgumentException e) { throw new IllegalStateException(e.getMessage(), e); }
    if (result.getSeatbidCount() == 0 && !noBidSet) throw new IllegalStateException("builder: BidResponse needs seatbid[] or noBid(nbr)");
    return result;
  }

  /**
   * Builds the response and serializes it as OpenRTB JSON bytes.
   *
   * @return UTF-8 JSON bytes
   */
  public byte[] buildJson() {
    return com.oakrtb.sdk.codec.Json.toJsonBytes(build());
  }

  /**
   * Builds JSON and validates it against JSON Schema, returning the payload and validation results.
   *
   * @return JSON and a {@link Report}
   */
  public ValidatedPayload buildValidated() {
    byte[] json = buildJson();
    return new ValidatedPayload(json, Schema.response(json));
  }

  /**
   * Creates a fluent builder for one {@link Bid}.
   *
   * @param id bid ID (required)
   * @param impid corresponding Imp.id (required)
   * @param price CPM price, must be &gt; 0
   * @return a {@link BidBuilder} instance
   */
  public static BidBuilder bid(String id, String impid, double price) {
    return new BidBuilder(id, impid, price);
  }

  /**
   * Fluent builder for one {@link Bid}, for use with {@link #addSeatBid(String, Bid...)}.
   */
  public static final class BidBuilder {
    private final Bid.Builder b = Bid.newBuilder();

    BidBuilder(String id, String impid, double price) {
      b.setId(id).setImpid(impid).setPrice(price);
    }

    /**
     * Sets ad markup (adm).
     *
     * @param adm ad creative content
     * @return this builder
     */
    public BidBuilder adm(String adm) {
      b.setAdm(adm);
      return this;
    }

    /**
     * Sets the win notice URL (nurl).
     *
     * @param u win notice URL
     * @return this builder
     */
    public BidBuilder nurl(String u) {
      b.setNurl(u);
      return this;
    }

    /**
     * Sets the billing notice URL (burl).
     *
     * @param u billing notice URL
     * @return this builder
     */
    public BidBuilder burl(String u) {
      b.setBurl(u);
      return this;
    }

    /**
     * Sets the creative ID (crid).
     *
     * @param crid creative identifier
     * @return this builder
     */
    public BidBuilder crid(String crid) {
      b.setCrid(crid);
      return this;
    }

    /**
     * Sets the campaign ID (cid).
     *
     * @param cid campaign identifier
     * @return this builder
     */
    public BidBuilder cid(String cid) {
      b.setCid(cid);
      return this;
    }

    /**
     * Sets advertiser domains (adomain).
     *
     * @param domains one or more domains
     * @return this builder
     */
    public BidBuilder adomain(String... domains) {
      b.clearAdomain();
      for (String d : domains) {
        b.addAdomain(d);
      }
      return this;
    }

    /**
     * Sets the creative width and height (w/h).
     *
     * @param w width in pixels
     * @param h height in pixels
     * @return this builder
     */
    public BidBuilder size(int w, int h) {
      b.setW(w).setH(h);
      return this;
    }

    /**
     * Sets the PMP deal ID (dealid).
     *
     * @param id deal ID
     * @return this builder
     */
    public BidBuilder dealId(String id) {
      b.setDealid(id);
      return this;
    }

    /**
     * Sets the markup type (proto {@link MarkupType} / Bid.mtype).
     *
     * @param m markup type enum
     * @return this builder
     */
    public BidBuilder markupType(MarkupType m) {
      b.setMtype(m);
      return this;
    }

    /**
     * Sets the markup type to Banner (mtype=1).
     *
     * @return this builder
     */
    public BidBuilder banner() {
      return markupType(MarkupType.MARKUP_TYPE_BANNER);
    }

    /**
     * Sets the markup type to Video (mtype=2).
     *
     * @return this builder
     */
    public BidBuilder video() {
      return markupType(MarkupType.MARKUP_TYPE_VIDEO);
    }

    /**
     * Sets the markup type to Audio (mtype=3).
     *
     * @return this builder
     */
    public BidBuilder audio() {
      return markupType(MarkupType.MARKUP_TYPE_AUDIO);
    }

    /**
     * Sets the markup type to Native (mtype=4).
     *
     * @return this builder
     */
    public BidBuilder nativeAd() {
      return markupType(MarkupType.MARKUP_TYPE_NATIVE);
    }

    /**
     * Sets the creative duration in seconds, typically for video/audio.
     *
     * @param seconds duration in seconds
     * @return this builder
     */
    public BidBuilder dur(int seconds) {
      b.setDur(seconds);
      return this;
    }

    /**
     * Builds a protobuf {@link Bid}.
     *
     * @return the bid object
     */
    public Bid build() {
      return b.build();
    }
  }
}
