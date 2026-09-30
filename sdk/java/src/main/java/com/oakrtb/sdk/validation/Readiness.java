package com.oakrtb.sdk.validation;

import com.oakrtb.openrtb.v2.Audio;
import com.oakrtb.openrtb.v2.Banner;
import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.Native;
import com.oakrtb.openrtb.v2.Video;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

/** Model format readiness checks, independent of query views. */
public final class Readiness {
 private Readiness() {}
  /**
   * Checks that the selected markup type exists on Imp and meets OakRTB builder-level requirements.
   *
   * @param imp Imp model
   * @param mtype OpenRTB mtype (1–4)
   * @return validation results
   */
  public static CheckResult impReady(Imp imp, int mtype) {
    Objects.requireNonNull(imp, "imp");
    List<CheckIssue> issues = new ArrayList<>();
    String path = "imp[" + imp.getId() + "]";
    int chosen = 1 << (mtype - 1);
    if (mtype < 1 || mtype > 4) {
      issues.add(
          new CheckIssue(
              IssueCode.MTYPE_UNKNOWN,
              Severity.ERROR,
              path + ".mtype",
              "mtype must be 1–4 for ImpReady"));
      return CheckResult.of(issues);
    }
    return impReadyChosen(imp, chosen, path, issues);
  }

  /**
   * Like {@link #impReady(Imp, int)}, but selects a format with a single-bit mask.
   *
   * @param imp Imp model
   * @param chosen mask containing exactly one format
   * @return validation results
   */
  public static CheckResult impReadyMask(Imp imp, int chosen) {
    Objects.requireNonNull(imp, "imp");

    List<CheckIssue> issues = new ArrayList<>();
    String path = "imp[" + imp.getId() + "]";
    if (Integer.bitCount(chosen) != 1) {
      issues.add(
          new CheckIssue(
              IssueCode.MULTI_NEEDS_CHOICE,
              Severity.ERROR,
              path,
              "chosen MarkupMask must be exactly one bit"));
      return CheckResult.of(issues);
    }
    return impReadyChosen(imp, chosen, path, issues);
  }

  private static CheckResult impReadyChosen(
      Imp imp, int chosen, String path, List<CheckIssue> issues) {
    if (!formatPresent(imp, chosen)) {
      issues.add(
          new CheckIssue(
              IssueCode.FORMAT_NOT_ON_IMP,
              Severity.ERROR,
              path,
              "chosen format is not present on Imp"));
      return CheckResult.of(issues);
    }
    if (chosen == 1) {
      Banner b = imp.getBanner();
      boolean sizeOk =
          b != null && ((b.getW() > 0 && b.getH() > 0) || b.getFormatCount() > 0);
      if (!sizeOk) {
        issues.add(
            new CheckIssue(
                IssueCode.BANNER_SIZE_MISSING,
                Severity.ERROR,
                path + ".banner",
                "banner needs w/h or format[]"));
      }
    }
    if (chosen == 2) {
      Video v = imp.getVideo();
      if (v == null || v.getMimesCount() == 0) {
        issues.add(
            new CheckIssue(
                IssueCode.VIDEO_MIMES_MISSING,
                Severity.ERROR,
                path + ".video.mimes",
                "video.mimes is required"));
      } else {
        if (v.getMinduration() == 0 && v.getMaxduration() == 0) {
          issues.add(
              new CheckIssue(
                  IssueCode.VIDEO_DURATION_UNSET,
                  Severity.WARN,
                  path + ".video",
                  "video minduration/maxduration unset"));
        }
        if (v.getProtocolsCount() == 0) {
          issues.add(
              new CheckIssue(
                  IssueCode.VIDEO_PROTOCOLS_UNSET,
                  Severity.WARN,
                  path + ".video.protocols",
                  "video.protocols unset"));
        }
      }
    }
    if (chosen == 4) {
      Audio a = imp.getAudio();
      if (a == null || a.getMimesCount() == 0) {
        issues.add(
            new CheckIssue(
                IssueCode.AUDIO_MIMES_MISSING,
                Severity.ERROR,
                path + ".audio.mimes",
                "audio.mimes is required"));
      }
    }
    if (chosen == 8) {
      Native n = imp.getNative();
      if (n == null || n.getRequest().isBlank()) {
        issues.add(
            new CheckIssue(
                IssueCode.NATIVE_REQUEST_MISSING,
                Severity.ERROR,
                path + ".native.request",
                "native.request is required"));
      }
    }
    return CheckResult.of(issues);
  }

  /** Check all formats present on an Imp; callers may inspect warnings separately. */
  public static CheckResult impReady(Imp imp) {
    Objects.requireNonNull(imp, "imp");
    List<CheckIssue> issues = new ArrayList<>();
    for (int mtype = 1; mtype <= 4; mtype++) {
      if (formatPresent(imp, 1 << (mtype - 1))) issues.addAll(impReady(imp, mtype).issues());
    }
    return CheckResult.of(issues);
  }
  private static boolean formatPresent(Imp imp, int chosen) {
    return switch (chosen) {
      case 1 -> imp.hasBanner();
      case 2 -> imp.hasVideo();
      case 4 -> imp.hasAudio();
      case 8 -> imp.hasNative();
      default -> false;
    };
  }
}
