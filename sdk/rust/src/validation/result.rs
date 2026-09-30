/// Finding severity.
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub enum Severity {
    /// Error: [`CheckResult::ok`] is false.
    Error,
    /// Warning: does not affect `ok()`, but requires attention.
    Warn,
}

impl Severity {
    /// Returns the string `"ERROR"` or `"WARN"`.
    pub fn as_str(self) -> &'static str {
        match self {
            Severity::Error => "ERROR",
            Severity::Warn => "WARN",
        }
    }
}

/// The selected presentation format is absent from Imp.
pub const CODE_FORMAT_NOT_ON_IMP: &str = "FORMAT_NOT_ON_IMP";
/// The format bitmask must select exactly one format.
pub const CODE_MULTI_NEEDS_CHOICE: &str = "MULTI_NEEDS_CHOICE";
/// Banner is missing w/h or format[].
pub const CODE_BANNER_SIZE_MISSING: &str = "BANNER_SIZE_MISSING";
/// Video is missing mimes.
pub const CODE_VIDEO_MIMES_MISSING: &str = "VIDEO_MIMES_MISSING";
/// Audio is missing mimes.
pub const CODE_AUDIO_MIMES_MISSING: &str = "AUDIO_MIMES_MISSING";
/// Native is missing the request string.
pub const CODE_NATIVE_REQUEST_MISSING: &str = "NATIVE_REQUEST_MISSING";
/// Video has neither minduration nor maxduration (warning).
pub const CODE_VIDEO_DURATION_UNSET: &str = "VIDEO_DURATION_UNSET";
/// Video protocols are unset (warning).
pub const CODE_VIDEO_PROTOCOLS_UNSET: &str = "VIDEO_PROTOCOLS_UNSET";
/// The response/payload structure cannot be parsed (e.g. not a JSON object).
pub const CODE_MALFORMED: &str = "MALFORMED";
/// The response ID does not match the request ID.
pub const CODE_REQUEST_ID_MISMATCH: &str = "REQUEST_ID_MISMATCH";
/// Bid.impid does not exist in the request.
pub const CODE_IMP_NOT_FOUND: &str = "IMP_NOT_FOUND";
/// A multi-format Imp requires Bid.mtype.
pub const CODE_MTYPE_REQUIRED: &str = "MTYPE_REQUIRED";
/// Bid.mtype does not match the Imp format.
pub const CODE_MTYPE_MISMATCH: &str = "MTYPE_MISMATCH";
/// mtype is outside both 1–4 and the vendor range.
pub const CODE_MTYPE_UNKNOWN: &str = "MTYPE_UNKNOWN";
/// Vendor-specific mtype ≥ 500 (warning).
pub const CODE_MTYPE_VENDOR: &str = "MTYPE_VENDOR";
/// The bid is below imp.bidfloor (warning).
pub const CODE_PRICE_BELOW_FLOOR: &str = "PRICE_BELOW_FLOOR";
/// Response currency differs from imp.bidfloorcur; floor comparison is skipped (warning).
pub const CODE_FLOOR_CUR_DIFF: &str = "FLOOR_CUR_DIFF";
/// bid.attr overlaps with format battr (warning).
pub const CODE_ATTR_BLOCKED: &str = "ATTR_BLOCKED";
/// adomain matches BidRequest.badv (warning).
pub const CODE_ADOMAIN_BLOCKED: &str = "ADOMAIN_BLOCKED";
/// bundle matches BidRequest.bapp (warning).
pub const CODE_BUNDLE_BLOCKED: &str = "BUNDLE_BLOCKED";
/// cat matches BidRequest.bcat (warning).
pub const CODE_CAT_BLOCKED: &str = "CAT_BLOCKED";
/// Response currency is not in BidRequest.cur (warning).
pub const CODE_CUR_NOT_ALLOWED: &str = "CUR_NOT_ALLOWED";
/// The deadline at 85% of request tmax has passed (warning).
pub const CODE_PAST_DEADLINE: &str = "PAST_DEADLINE";

/// A single bid check finding.
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct CheckIssue {
    /// Stable error code (see the `CODE_*` constants).
    pub code: String,
    /// Severity level.
    pub severity: Severity,
    /// JSON-style path (e.g. `seatbid[0].bid[1].price`).
    pub path: String,
    /// Human-readable explanation.
    pub message: String,
}

impl CheckIssue {
    pub(crate) fn new(
        code: &str,
        severity: Severity,
        path: impl Into<String>,
        message: impl Into<String>,
    ) -> Self {
        Self {
            code: code.to_string(),
            severity,
            path: path.into(),
            message: message.into(),
        }
    }
}

/// Aggregated bid check results.
#[derive(Clone, Debug, Default)]
pub struct CheckResult {
    /// All findings.
    pub issues: Vec<CheckIssue>,
}

impl CheckResult {
    /// True when there are no Error findings.
    pub fn ok(&self) -> bool {
        !self.issues.iter().any(|i| i.severity == Severity::Error)
    }

    /// Reports whether a finding with the specified code exists.
    pub fn has(&self, code: &str) -> bool {
        self.issues.iter().any(|i| i.code == code)
    }

    /// ERROR findings.
    pub fn errors(&self) -> impl Iterator<Item = &CheckIssue> {
        self.issues.iter().filter(|i| i.severity == Severity::Error)
    }

    /// WARN findings.
    pub fn warnings(&self) -> impl Iterator<Item = &CheckIssue> {
        self.issues.iter().filter(|i| i.severity == Severity::Warn)
    }
}

pub const CODE_DEAL_REQUIRED: &str = "DEAL_REQUIRED";

pub const CODE_DEAL_NOT_FOUND: &str = "DEAL_NOT_FOUND";

pub const CODE_DEAL_SEAT_NOT_ALLOWED: &str = "DEAL_SEAT_NOT_ALLOWED";

pub const CODE_DEAL_ADOMAIN_NOT_ALLOWED: &str = "DEAL_ADOMAIN_NOT_ALLOWED";

pub const CODE_DEAL_PRICE_MISMATCH: &str = "DEAL_PRICE_MISMATCH";
