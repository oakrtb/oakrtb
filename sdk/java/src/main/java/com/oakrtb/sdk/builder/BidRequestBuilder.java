package com.oakrtb.sdk.builder;
import com.oakrtb.sdk.codec.Json;

import com.oakrtb.openrtb.v2.App;
import com.oakrtb.openrtb.v2.BidRequest;
import com.oakrtb.openrtb.v2.Device;
import com.oakrtb.openrtb.v2.Dooh;
import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.Regs;
import com.oakrtb.openrtb.v2.Site;
import com.oakrtb.openrtb.v2.Source;
import com.oakrtb.openrtb.v2.User;
import com.oakrtb.sdk.jsonschema.Report;
import com.oakrtb.sdk.jsonschema.Schema;

import java.util.Objects;

/**
 * Fluent builder for an OpenRTB {@link BidRequest}.
 *
 * <p>For SSP/Exchange auction request construction; {@link #build()} validates at, cur, imp, and each Imp's format fields.
 */
public final class BidRequestBuilder {
  private final BidRequest.Builder req = BidRequest.newBuilder();
  private String error;

  private BidRequestBuilder(String id) {
    if (id == null || id.isBlank()) {
      error = "builder: BidRequest.id is required";
    } else {
      req.setId(id);
    }
  }

  /**
   * Creates a request builder.
   *
   * @param id unique auction request ID (required and nonempty)
   * @return a new builder instance
   */
  public static BidRequestBuilder create(String id) {
    return new BidRequestBuilder(id);
  }

  /**
   * Sets a first-price auction (at=1).
   *
   * @return this builder
   */
  public BidRequestBuilder firstPrice() {
    return auctionType(1);
  }

  /**
   * Sets a second-price-plus auction (at=2).
   *
   * @return this builder
   */
  public BidRequestBuilder secondPricePlus() {
    return auctionType(2);
  }

  /**
   * Sets the auction type (at).
   *
   * @param at OpenRTB auction type enum value
   * @return this builder
   */
  public BidRequestBuilder auctionType(int at) {
    req.setAt(at);
    return this;
  }

  /**
   * Sets the maximum response time in milliseconds.
   *
   * @param ms tmax in milliseconds
   * @return this builder
   */
  public BidRequestBuilder tmax(int ms) {
    req.setTmax(ms);
    return this;
  }

  /**
   * Sets allowed bid currencies (ISO-4217, at least one).
   *
   * @param codes one or more currency codes
   * @return this builder
   */
  public BidRequestBuilder currency(String... codes) {
    req.clearCur();
    for (String c : codes) {
      req.addCur(c);
    }
    return this;
  }

  /**
   * Marks a test request (test=1).
   *
   * @return this builder
   */
  public BidRequestBuilder test() {
    req.setTest(1);
    return this;
  }

  /**
   * Sets blocked advertiser categories (bcat).
   *
   * @param cats content categories, such as IAB categories
   * @return this builder
   */
  public BidRequestBuilder bcat(String... cats) {
    req.clearBcat();
    for (String c : cats) {
      req.addBcat(c);
    }
    return this;
  }

  /**
   * Sets blocked advertiser domains (badv).
   *
   * @param domains domain list
   * @return this builder
   */
  public BidRequestBuilder badv(String... domains) {
    req.clearBadv();
    for (String d : domains) {
      req.addBadv(d);
    }
    return this;
  }

  /**
   * Sets website inventory (site) and clears app/dooh.
   *
   * @param site Site object
   * @return this builder
   */
  public BidRequestBuilder site(Site site) {
    req.clearApp().clearDooh().setSite(Objects.requireNonNull(site));
    return this;
  }

  /**
   * Sets mobile app inventory (app) and clears site/dooh.
   *
   * @param app App object
   * @return this builder
   */
  public BidRequestBuilder app(App app) {
    req.clearSite().clearDooh().setApp(Objects.requireNonNull(app));
    return this;
  }

  /**
   * Sets DOOH inventory (dooh) and clears site/app.
   *
   * @param dooh Dooh object
   * @return this builder
   */
  public BidRequestBuilder dooh(Dooh dooh) {
    req.clearSite().clearApp().setDooh(Objects.requireNonNull(dooh));
    return this;
  }

  /**
   * Sets device information (device).
   *
   * @param device Device object
   * @return this builder
   */
  public BidRequestBuilder device(Device device) {
    req.setDevice(Objects.requireNonNull(device));
    return this;
  }

  /**
   * Sets user information (user).
   *
   * @param user User object
   * @return this builder
   */
  public BidRequestBuilder user(User user) {
    req.setUser(Objects.requireNonNull(user));
    return this;
  }

  /**
   * Sets regulatory/privacy information (regs).
   *
   * @param regs Regs object
   * @return this builder
   */
  public BidRequestBuilder regs(Regs regs) {
    req.setRegs(Objects.requireNonNull(regs));
    return this;
  }

  /**
   * Sets the request source (source).
   *
   * @param source Source object
   * @return this builder
   */
  public BidRequestBuilder source(Source source) {
    req.setSource(Objects.requireNonNull(source));
    return this;
  }

  /**
   * Appends an impression (Imp).
   *
   * @param imp Imp object (typically built with {@link ImpBuilders})
   * @return this builder
   */
  public BidRequestBuilder addImp(Imp imp) {
    req.addImp(Objects.requireNonNull(imp));
    return this;
  }

  /**
   * Builds a protobuf {@link BidRequest} and validates id, at, cur, imp, and format fields.
   *
   * @return the constructed request
   * @throws IllegalStateException if required-field or format validation fails
   */
  public BidRequest build() {
    if (error != null) throw new IllegalStateException(error);
    BidRequest result = req.build();
    try { com.oakrtb.sdk.validation.BasicValidation.validateRequest(result); }
    catch (IllegalArgumentException e) { throw new IllegalStateException(e.getMessage(), e); }
    for (int i = 0; i < result.getImpCount(); i++) checkImp(result.getImp(i), i);
    return result;
  }

  /**
   * Builds the request and serializes it as OpenRTB JSON bytes.
   *
   * @return UTF-8 JSON bytes
   */
  public byte[] buildJson() {
    return com.oakrtb.sdk.codec.Json.toJsonBytes(build());
  }

  /**
   * Builds JSON and validates it against JSON Schema.
   *
   * @return JSON and validation results
   */
  public ValidatedPayload buildValidated() {
    byte[] json = buildJson();
    return new ValidatedPayload(json, Schema.request(json));
  }

  private static void checkImp(Imp imp, int i) {
    var result = com.oakrtb.sdk.validation.Readiness.impReady(imp);
    if (!result.ok()) throw new IllegalStateException("builder: imp[" + i + "]: " + result.errors().getFirst().message());
  }
}
