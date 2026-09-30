use super::*;
use crate::proto::{Banner, Imp};

/// Checks that Imp has the format selected by OpenRTB `Bid.mtype` (1–4) and meets builder-level requirements.
pub fn imp_ready_mtype(imp: &Imp, mtype: i32) -> CheckResult {
    let path = format!("imp[{}]", imp.id);
    if !(1..=4).contains(&mtype) {
        return CheckResult {
            issues: vec![CheckIssue::new(
                CODE_MTYPE_UNKNOWN,
                Severity::Error,
                format!("{path}.mtype"),
                "mtype must be 1–4 for ImpReady",
            )],
        };
    }
    imp_ready_chosen(imp, 1 << (mtype - 1), &path)
}

/// Like [`imp_ready_mtype`], but selects a format with a single-bit format mask.
///
/// **Note**: the format bitmask identifies presentation types on an Imp (banner/video/audio/native),
/// distinct from `Banner.format[]` (size list).
pub fn imp_ready_markup(imp: &Imp, chosen: u8) -> CheckResult {
    let path = format!("imp[{}]", imp.id);
    if chosen.count_ones() != 1 {
        return CheckResult {
            issues: vec![CheckIssue::new(
                CODE_MULTI_NEEDS_CHOICE,
                Severity::Error,
                path,
                "chosen MarkupMask must be exactly one bit",
            )],
        };
    }
    imp_ready_chosen(imp, chosen, &path)
}

fn imp_ready_chosen(imp: &Imp, chosen: u8, path: &str) -> CheckResult {
    let mut issues = Vec::new();
    if !format_present(imp, chosen) {
        return CheckResult {
            issues: vec![CheckIssue::new(
                CODE_FORMAT_NOT_ON_IMP,
                Severity::Error,
                path,
                "chosen format is not present on Imp",
            )],
        };
    }
    if chosen == 1 {
        let size_ok = imp.banner.as_ref().map(banner_size_ok).unwrap_or(false);
        if !size_ok {
            issues.push(CheckIssue::new(
                CODE_BANNER_SIZE_MISSING,
                Severity::Error,
                format!("{path}.banner"),
                "banner needs w/h or format[]",
            ));
        }
    }
    if chosen == 2 {
        let v = imp.video.as_ref();
        let mimes = v.map(|v| &v.mimes);
        if mimes.map(|a| a.is_empty()).unwrap_or(true) {
            issues.push(CheckIssue::new(
                CODE_VIDEO_MIMES_MISSING,
                Severity::Error,
                format!("{path}.video.mimes"),
                "video.mimes is required",
            ));
        } else if let Some(v) = v {
            let mind = v.minduration.unwrap_or(0);
            let maxd = v.maxduration.unwrap_or(0);
            if mind == 0 && maxd == 0 {
                issues.push(CheckIssue::new(
                    CODE_VIDEO_DURATION_UNSET,
                    Severity::Warn,
                    format!("{path}.video"),
                    "video minduration/maxduration unset",
                ));
            }
            let protocols = Some(&v.protocols);
            if protocols.map(|a| a.is_empty()).unwrap_or(true) {
                issues.push(CheckIssue::new(
                    CODE_VIDEO_PROTOCOLS_UNSET,
                    Severity::Warn,
                    format!("{path}.video.protocols"),
                    "video.protocols unset",
                ));
            }
        }
    }
    if chosen == 4 {
        let mimes = imp.audio.as_ref().map(|o| &o.mimes);
        if mimes.map(|a| a.is_empty()).unwrap_or(true) {
            issues.push(CheckIssue::new(
                CODE_AUDIO_MIMES_MISSING,
                Severity::Error,
                format!("{path}.audio.mimes"),
                "audio.mimes is required",
            ));
        }
    }
    if chosen == 8 {
        let req = imp
            .native
            .as_ref()
            .map(|o| o.request.as_str())
            .unwrap_or("");
        if req.trim().is_empty() {
            issues.push(CheckIssue::new(
                CODE_NATIVE_REQUEST_MISSING,
                Severity::Error,
                format!("{path}.native.request"),
                "native.request is required",
            ));
        }
    }
    CheckResult { issues }
}

fn banner_size_ok(b: &Banner) -> bool {
    (b.w.unwrap_or(0) > 0 && b.h.unwrap_or(0) > 0) || !b.format.is_empty()
}

/// Check all present formats. Warnings do not prevent building.
pub fn imp_ready(imp: &Imp) -> CheckResult {
    let mut result = CheckResult::default();
    for mtype in 1..=4 {
        if format_present(imp, 1 << (mtype - 1)) {
            result.issues.extend(imp_ready_mtype(imp, mtype).issues);
        }
    }
    result
}
fn format_present(imp: &Imp, chosen: u8) -> bool {
    match chosen {
        1 => imp.banner.is_some(),
        2 => imp.video.is_some(),
        4 => imp.audio.is_some(),
        8 => imp.native.is_some(),
        _ => false,
    }
}
