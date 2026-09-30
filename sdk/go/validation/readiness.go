package validation

import (
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"strings"
)

// ImpReadyMtype checks that the selected markup type exists in Imp and meets builder-level requirements.
func ImpReadyMtype(imp *openrtb.Imp, mtype int32) CheckResult {
	if imp == nil {
		return CheckResult{Issues: []CheckIssue{{
			Code: CodeMalformed, Severity: SeverityError, Path: "imp", Message: "nil Imp",
		}}}
	}
	path := "imp[" + imp.Id + "]"
	if mtype < 1 || mtype > 4 {
		return CheckResult{Issues: []CheckIssue{{
			Code: CodeMtypeUnknown, Severity: SeverityError, Path: path + ".mtype",
			Message: "mtype must be 1–4 for ImpReady",
		}}}
	}
	return impReadyChosen(imp, uint8(1<<(mtype-1)), path)
}

// ImpReadyMarkup is like ImpReadyMtype but selects a type with a single-bit MarkupMask.
func ImpReadyMarkup(imp *openrtb.Imp, chosen uint8) CheckResult {
	if imp == nil {
		return CheckResult{Issues: []CheckIssue{{
			Code: CodeMalformed, Severity: SeverityError, Path: "imp", Message: "nil Imp",
		}}}
	}
	path := "imp[" + imp.Id + "]"
	if chosen == 0 || chosen&(chosen-1) != 0 {
		return CheckResult{Issues: []CheckIssue{{
			Code: CodeMultiNeedsChoice, Severity: SeverityError, Path: path,
			Message: "chosen MarkupMask must be exactly one bit",
		}}}
	}
	return impReadyChosen(imp, chosen, path)
}

func impReadyChosen(imp *openrtb.Imp, chosen uint8, path string) CheckResult {
	var issues []CheckIssue
	if !formatPresent(imp, chosen) {
		return CheckResult{Issues: []CheckIssue{{
			Code: CodeFormatNotOnImp, Severity: SeverityError, Path: path,
			Message: "chosen format is not present on Imp",
		}}}
	}
	if chosen == 1 {
		b := imp.Banner
		sizeOK := b != nil && ((b.GetW() > 0 && b.GetH() > 0) || len(b.Format) > 0)
		if !sizeOK {
			issues = append(issues, CheckIssue{
				Code: CodeBannerSizeMissing, Severity: SeverityError, Path: path + ".banner",
				Message: "banner needs w/h or format[]",
			})
		}
	}
	if chosen == 2 {
		v := imp.Video
		if v == nil || len(v.Mimes) == 0 {
			issues = append(issues, CheckIssue{
				Code: CodeVideoMimesMissing, Severity: SeverityError, Path: path + ".video.mimes",
				Message: "video.mimes is required",
			})
		} else {
			if v.GetMinduration() == 0 && v.GetMaxduration() == 0 {
				issues = append(issues, CheckIssue{
					Code: CodeVideoDurationUnset, Severity: SeverityWarn, Path: path + ".video",
					Message: "video minduration/maxduration unset",
				})
			}
			if len(v.Protocols) == 0 {
				issues = append(issues, CheckIssue{
					Code: CodeVideoProtocolsUnset, Severity: SeverityWarn, Path: path + ".video.protocols",
					Message: "video.protocols unset",
				})
			}
		}
	}
	if chosen == 4 {
		a := imp.Audio
		if a == nil || len(a.Mimes) == 0 {
			issues = append(issues, CheckIssue{
				Code: CodeAudioMimesMissing, Severity: SeverityError, Path: path + ".audio.mimes",
				Message: "audio.mimes is required",
			})
		}
	}
	if chosen == 8 {
		n := imp.Native
		if n == nil || strings.TrimSpace(n.Request) == "" {
			issues = append(issues, CheckIssue{
				Code: CodeNativeRequestMissing, Severity: SeverityError, Path: path + ".native.request",
				Message: "native.request is required",
			})
		}
	}
	return CheckResult{Issues: issues}
}

// ImpReady checks every format present on an Imp. Warnings do not prevent building.
func ImpReady(imp *openrtb.Imp) CheckResult {
	if imp == nil {
		return ImpReadyMtype(nil, 1)
	}
	var result CheckResult
	for mtype := int32(1); mtype <= 4; mtype++ {
		if formatPresent(imp, 1<<(mtype-1)) {
			result.Issues = append(result.Issues, ImpReadyMtype(imp, mtype).Issues...)
		}
	}
	return result
}
func formatPresent(imp *openrtb.Imp, chosen uint8) bool {
	switch chosen {
	case 1:
		return imp.Banner != nil
	case 2:
		return imp.Video != nil
	case 4:
		return imp.Audio != nil
	case 8:
		return imp.Native != nil
	}
	return false
}
