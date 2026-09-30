package validation

// Severity indicates the severity of a CheckIssue.
type Severity string

const (
	// SeverityError indicates an error that should cause bid rejection.
	SeverityError Severity = "ERROR"
	// SeverityWarning indicates a warning that permits continued processing but requires attention.
	SeverityWarn Severity = "WARN"
)

// Stable issue codes aligned with Java / Rust.
const (
	CodeFormatNotOnImp       = "FORMAT_NOT_ON_IMP"
	CodeMultiNeedsChoice     = "MULTI_NEEDS_CHOICE"
	CodeBannerSizeMissing    = "BANNER_SIZE_MISSING"
	CodeVideoMimesMissing    = "VIDEO_MIMES_MISSING"
	CodeAudioMimesMissing    = "AUDIO_MIMES_MISSING"
	CodeNativeRequestMissing = "NATIVE_REQUEST_MISSING"
	CodeVideoDurationUnset   = "VIDEO_DURATION_UNSET"
	CodeVideoProtocolsUnset  = "VIDEO_PROTOCOLS_UNSET"

	CodeRequestIDMismatch = "REQUEST_ID_MISMATCH"
	CodeMalformed         = "MALFORMED"
	CodeImpNotFound       = "IMP_NOT_FOUND"
	CodeMtypeRequired     = "MTYPE_REQUIRED"
	CodeMtypeMismatch     = "MTYPE_MISMATCH"
	CodeMtypeUnknown      = "MTYPE_UNKNOWN"
	CodeMtypeVendor       = "MTYPE_VENDOR"
	CodePriceBelowFloor   = "PRICE_BELOW_FLOOR"
	CodeFloorCurDiff      = "FLOOR_CUR_DIFF"
	CodeAttrBlocked       = "ATTR_BLOCKED"
	CodeAdomainBlocked    = "ADOMAIN_BLOCKED"
	CodeBundleBlocked     = "BUNDLE_BLOCKED"
	CodeCatBlocked        = "CAT_BLOCKED"
	CodeCurNotAllowed     = "CUR_NOT_ALLOWED"
	CodePastDeadline      = "PAST_DEADLINE"
)

// CheckIssue is a finding produced by bid checks.
type CheckIssue struct {
	Code     string
	Severity Severity
	Path     string
	Message  string
}

// CheckResult aggregates findings; OK is true when there are no ERROR findings.
type CheckResult struct {
	Issues []CheckIssue
}

// OK returns true when there are no ERROR findings.
func (r CheckResult) OK() bool {
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			return false
		}
	}
	return true
}

// Has reports whether a finding with the specified code exists.
func (r CheckResult) Has(code string) bool {
	for _, i := range r.Issues {
		if i.Code == code {
			return true
		}
	}
	return false
}

// Errors returns all ERROR findings.
func (r CheckResult) Errors() []CheckIssue {
	var out []CheckIssue
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			out = append(out, i)
		}
	}
	return out
}

// Warnings returns all WARN findings.
func (r CheckResult) Warnings() []CheckIssue {
	var out []CheckIssue
	for _, i := range r.Issues {
		if i.Severity == SeverityWarn {
			out = append(out, i)
		}
	}
	return out
}

// Private marketplace contract violations.
const (
	CodeDealRequired          = "DEAL_REQUIRED"
	CodeDealNotFound          = "DEAL_NOT_FOUND"
	CodeDealSeatNotAllowed    = "DEAL_SEAT_NOT_ALLOWED"
	CodeDealAdomainNotAllowed = "DEAL_ADOMAIN_NOT_ALLOWED"
	CodeDealPriceMismatch     = "DEAL_PRICE_MISMATCH"
)
