package com.oakrtb.sdk.view;

/**
 * Auction request inventory type (site / app / dooh, mutually exclusive).
 *
 * <p>Used by {@link RequestView} to identify the traffic source type.
 *
 * <p><b>Note:</b> unrelated to proto {@code Content.Channel}; identifies only the top-level BidRequest inventory object type.
 */
public enum Inventory {
  /** No site/app/dooh is set. */
  NONE,
  /** Website inventory. */
  SITE,
  /** Mobile app inventory. */
  APP,
  /** Digital out-of-home (DOOH) inventory. */
  DOOH;

  /**
   * Returns a lowercase wire-style string (site/app/dooh/none).
   *
   * @return the inventory type string
   */
  @Override
  public String toString() {
    return switch (this) {
      case SITE -> "site";
      case APP -> "app";
      case DOOH -> "dooh";
      default -> "none";
    };
  }
}
