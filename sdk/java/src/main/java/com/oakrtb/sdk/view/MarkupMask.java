package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.MarkupType;

/**
 * Imp/Bid markup type bitmask (banner / video / audio / native).
 *
 * <p>Quickly identifies ad formats in impressions or bids during view and bidcheck processing.
 *
 * <p><b>Note:</b> unrelated to proto {@code Banner.Format} (size format list); constants such as {@code BANNER} represent
 * OpenRTB ad format bits corresponding to {@code Bid.mtype} 1–4.
 */
public final class MarkupMask {
  /** No format bits. */
  public static final int NONE = 0;
  /** Banner format bit, corresponding to Bid.mtype=1. */
  public static final int BANNER = 1 << 0;
  /** Video format bit, corresponding to Bid.mtype=2. */
  public static final int VIDEO = 1 << 1;
  /** Audio format bit, corresponding to Bid.mtype=3. */
  public static final int AUDIO = 1 << 2;
  /** Native format bit, corresponding to Bid.mtype=4. */
  public static final int NATIVE = 1 << 3;

  private final int bits;

  private MarkupMask(int bits) {
    this.bits = bits;
  }

  /**
   * Creates a mask from raw bits.
   *
   * @param bits bit combination (typically the bitwise OR of constants such as {@link #BANNER})
   * @return a MarkupMask instance
   */
  public static MarkupMask of(int bits) {
    return new MarkupMask(bits);
  }

  /**
   * Returns the raw bits.
   *
   * @return the integer bitmask
   */
  public int bits() {
    return bits;
  }

  /**
   * Reports whether the specified bit flag is present.
   *
   * @param flag single-bit constant such as {@link #BANNER} or {@link #VIDEO}
   * @return {@code true} if present
   */
  public boolean has(int flag) {
    return (bits & flag) != 0;
  }

  /**
   * Reports whether Banner format is present.
   *
   * @return {@code true} if Banner is present
   */
  public boolean hasBanner() {
    return has(BANNER);
  }

  /**
   * Reports whether Video format is present.
   *
   * @return {@code true} if Video is present
   */
  public boolean hasVideo() {
    return has(VIDEO);
  }

  /**
   * Reports whether Audio format is present.
   *
   * @return {@code true} if Audio is present
   */
  public boolean hasAudio() {
    return has(AUDIO);
  }

  /**
   * Reports whether Native format is present.
   *
   * @return {@code true} if Native is present
   */
  public boolean hasNative() {
    return has(NATIVE);
  }

  /**
   * Counts the set format bits.
   *
   * @return the format count (0–4)
   */
  public int count() {
    return Integer.bitCount(bits);
  }

  /**
   * Returns the format bit only when exactly one format is set; otherwise returns {@link #NONE}.
   *
   * @return a single format bit or NONE
   */
  public int primary() {
    return count() == 1 ? bits : NONE;
  }

  /**
   * Maps a single format to OpenRTB {@code Bid.mtype} (1–4); returns 0 for multiple or no formats.
   *
   * @return the mtype value
   */
  public int mtype() {
    return switch (primary()) {
      case BANNER -> 1;
      case VIDEO -> 2;
      case AUDIO -> 3;
      case NATIVE -> 4;
      default -> 0;
    };
  }

  @Override
  public boolean equals(Object o) {
    return o instanceof MarkupMask f && f.bits == bits;
  }

  @Override
  public int hashCode() {
    return bits;
  }

  /**
   * Returns the debug string {@code MarkupMask(bits)}.
   *
   * @return the string representation
   */
  @Override
  public String toString() {
    return "MarkupMask(" + bits + ")";
  }
public static MarkupMask fromImp(Imp imp) {
    if (imp == null) {
      return MarkupMask.of(MarkupMask.NONE);
    }
    int bits = MarkupMask.NONE;
    if (imp.hasBanner()) {
      bits |= MarkupMask.BANNER;
    }
    if (imp.hasVideo()) {
      bits |= MarkupMask.VIDEO;
    }
    if (imp.hasAudio()) {
      bits |= MarkupMask.AUDIO;
    }
    if (imp.hasNative()) {
      bits |= MarkupMask.NATIVE;
    }
    return MarkupMask.of(bits);
  }
public static MarkupMask fromMtype(MarkupType m) {
    if (m == null) {
      return MarkupMask.of(MarkupMask.NONE);
    }
    return switch (m) {
      case MARKUP_TYPE_BANNER -> MarkupMask.of(MarkupMask.BANNER);
      case MARKUP_TYPE_VIDEO -> MarkupMask.of(MarkupMask.VIDEO);
      case MARKUP_TYPE_AUDIO -> MarkupMask.of(MarkupMask.AUDIO);
      case MARKUP_TYPE_NATIVE -> MarkupMask.of(MarkupMask.NATIVE);
      default -> MarkupMask.of(MarkupMask.NONE);
    };
  }
public static MarkupMask fromMtype(int mtype) {
    return switch (mtype) {
      case 1 -> MarkupMask.of(MarkupMask.BANNER);
      case 2 -> MarkupMask.of(MarkupMask.VIDEO);
      case 3 -> MarkupMask.of(MarkupMask.AUDIO);
      case 4 -> MarkupMask.of(MarkupMask.NATIVE);
      default -> MarkupMask.of(MarkupMask.NONE);
    };
  }
}
