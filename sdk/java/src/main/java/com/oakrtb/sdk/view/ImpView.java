package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.Audio;
import com.oakrtb.openrtb.v2.Banner;
import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.Native;
import com.oakrtb.openrtb.v2.Pmp;
import com.oakrtb.openrtb.v2.Video;

/** Immutable projection of one impression. */
public record ImpView(
      Imp imp,
      String id,
      MarkupMask markup,
      String tagId,
      double bidFloor,
      String bidFloorCur,
      int instl,
      int secure,
      int rwdd,
      int ssai,
      Banner banner,
      Video video,
      Audio audio,
      Native nativeAd,
      Pmp pmp) {
  static ImpView of(Imp imp) {
    return new ImpView(
              imp,
              imp.getId(),
              MarkupMask.fromImp(imp),
              imp.getTagid(),
              imp.getBidfloor(),
              imp.getBidfloorcur(),
              imp.getInstl(),
              imp.getSecure(),
              imp.getRwdd(),
              imp.getSsai(),
              imp.hasBanner() ? imp.getBanner() : null,
              imp.hasVideo() ? imp.getVideo() : null,
              imp.hasAudio() ? imp.getAudio() : null,
              imp.hasNative() ? imp.getNative() : null,
              imp.hasPmp() ? imp.getPmp() : null);
  }
}
