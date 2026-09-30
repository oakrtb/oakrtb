package com.oakrtb.sdk.validation;

/**
 * Stable string error codes for the validation API, aligned with the Go/Rust SDKs.
 *
 * <p>Used by {@link CheckIssue#code()} and {@link CheckResult#has(String)}.
 */
public final class IssueCode {
  private IssueCode() {}

  /** The selected format is absent from Imp. */
  public static final String FORMAT_NOT_ON_IMP = "FORMAT_NOT_ON_IMP";
  /** MarkupMask must contain a single bit. */
  public static final String MULTI_NEEDS_CHOICE = "MULTI_NEEDS_CHOICE";
  /** Banner is missing w/h or format[]. */
  public static final String BANNER_SIZE_MISSING = "BANNER_SIZE_MISSING";
  /** Video is missing mimes. */
  public static final String VIDEO_MIMES_MISSING = "VIDEO_MIMES_MISSING";
  /** Audio is missing mimes. */
  public static final String AUDIO_MIMES_MISSING = "AUDIO_MIMES_MISSING";
  /** Native is missing request. */
  public static final String NATIVE_REQUEST_MISSING = "NATIVE_REQUEST_MISSING";
  /** Video has no duration range (warning). */
  public static final String VIDEO_DURATION_UNSET = "VIDEO_DURATION_UNSET";
  /** Video has no protocols (warning). */
  public static final String VIDEO_PROTOCOLS_UNSET = "VIDEO_PROTOCOLS_UNSET";

  /** The response/payload structure cannot be parsed (e.g. not a JSON object). */
  public static final String MALFORMED = "MALFORMED";
  /** The response ID does not match the request ID. */
  public static final String REQUEST_ID_MISMATCH = "REQUEST_ID_MISMATCH";
  /** bid.impid does not exist in the request. */
  public static final String IMP_NOT_FOUND = "IMP_NOT_FOUND";
  /** A multi-format Imp requires Bid.mtype. */
  public static final String MTYPE_REQUIRED = "MTYPE_REQUIRED";
  /** Bid.mtype does not match the Imp format. */
  public static final String MTYPE_MISMATCH = "MTYPE_MISMATCH";
  /** mtype is outside both 1–4 and the vendor range. */
  public static final String MTYPE_UNKNOWN = "MTYPE_UNKNOWN";
  /** Vendor-specific mtype ≥500 (warning). */
  public static final String MTYPE_VENDOR = "MTYPE_VENDOR";
  /** The bid is below the price floor (warning). */
  public static final String PRICE_BELOW_FLOOR = "PRICE_BELOW_FLOOR";
  /** Bid currency differs from imp.bidfloorcur; floor comparison is skipped (warning). */
  public static final String FLOOR_CUR_DIFF = "FLOOR_CUR_DIFF";
  /** bid.attr conflicts with format battr (warning). */
  public static final String ATTR_BLOCKED = "ATTR_BLOCKED";
  /** adomain matches BidRequest.badv (warning). */
  public static final String ADOMAIN_BLOCKED = "ADOMAIN_BLOCKED";
  /** bundle matches BidRequest.bapp (warning). */
  public static final String BUNDLE_BLOCKED = "BUNDLE_BLOCKED";
  /** cat matches BidRequest.bcat (warning). */
  public static final String CAT_BLOCKED = "CAT_BLOCKED";
  /** Response currency is not in BidRequest.cur (warning). */
  public static final String CUR_NOT_ALLOWED = "CUR_NOT_ALLOWED";
  /** The request deadline has passed (warning). */
  public static final String PAST_DEADLINE = "PAST_DEADLINE";
  public static final String DEAL_REQUIRED = "DEAL_REQUIRED";
  public static final String DEAL_NOT_FOUND = "DEAL_NOT_FOUND";
  public static final String DEAL_SEAT_NOT_ALLOWED = "DEAL_SEAT_NOT_ALLOWED";
  public static final String DEAL_ADOMAIN_NOT_ALLOWED = "DEAL_ADOMAIN_NOT_ALLOWED";
  public static final String DEAL_PRICE_MISMATCH = "DEAL_PRICE_MISMATCH";
}
