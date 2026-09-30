// Package bidcheck provides optional bid-to-request compatibility checks between request views and response models.
//
// These are neither JSON Schema nor basic validation checks; findings are returned in validation.CheckResult.
// OK() means no ERROR findings; WARN findings (low prices, blocked ads, deadlines, etc.) may still yield OK. Use Warnings()/Has for policy decisions.
package bidcheck

import (
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"math"
	"strconv"
	"strings"

	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/view"
)

// Bid checks one Bid against a request snapshot. Without response currency or seat context, it skips price comparisons and deal seat allowlists; use Response for full checks.
func Bid(req *view.RequestView, bid *openrtb.Bid) validation.CheckResult {
	return bidAt(req, bid, "", nil, "bid", newMatchContext(req))
}

// Response checks an entire BidResponse. All responses are checked for request ID consistency; an empty seatbid may still produce a PAST_DEADLINE warning.
func Response(req *view.RequestView, res *openrtb.BidResponse) validation.CheckResult {
	var issues []validation.CheckIssue
	if req == nil || res == nil {
		return validation.CheckResult{Issues: []validation.CheckIssue{{
			Code: validation.CodeMalformed, Severity: validation.SeverityError, Path: "", Message: "nil req or res",
		}}}
	}
	if res.Id != req.AuctionID() {
		issues = append(issues, validation.CheckIssue{Code: validation.CodeRequestIDMismatch, Severity: validation.SeverityError, Path: "BidResponse.id", Message: "response id does not match request id"})
	}
	if req.PastDeadline() {
		issues = append(issues, validation.CheckIssue{
			Code: validation.CodePastDeadline, Severity: validation.SeverityWarn, Path: "tmax",
			Message: "past request deadline (85% of tmax)",
		})
	}
	if len(res.Seatbid) == 0 {
		return validation.CheckResult{Issues: issues}
	}
	// Same blank rule as BasicValidation / checkFloor: trim; blank skips CUR + floor.
	resCur := strings.TrimSpace(res.Cur)
	if resCur != "" {
		allowed := false
		for _, c := range req.Currencies() {
			if asciiLower(resCur) == asciiLower(strings.TrimSpace(c)) {
				allowed = true
				break
			}
		}
		if !allowed {
			issues = append(issues, validation.CheckIssue{
				Code: validation.CodeCurNotAllowed, Severity: validation.SeverityWarn, Path: "BidResponse.cur",
				Message: "response currency not in BidRequest.cur",
			})
		}
	}
	ctx := newMatchContext(req)
	for i, sb := range res.Seatbid {
		if sb == nil {
			issues = append(issues, validation.CheckIssue{
				Code: validation.CodeMalformed, Severity: validation.SeverityError, Path: "seatbid[" + itoa(i) + "]",
				Message: "seatbid entry is nil",
			})
			continue
		}
		if len(sb.Bid) == 0 {
			issues = append(issues, validation.CheckIssue{
				Code: validation.CodeMalformed, Severity: validation.SeverityError, Path: "seatbid[" + itoa(i) + "].bid",
				Message: "seatbid.bid must be a non-empty array",
			})
			continue
		}
		for j, bid := range sb.Bid {
			if bid == nil {
				issues = append(issues, validation.CheckIssue{
					Code: validation.CodeMalformed, Severity: validation.SeverityError,
					Path:    "seatbid[" + itoa(i) + "].bid[" + itoa(j) + "]",
					Message: "bid is nil",
				})
				continue
			}
			path := "seatbid[" + itoa(i) + "].bid[" + itoa(j) + "]"
			one := bidAt(req, bid, resCur, &sb.Seat, path, ctx)
			issues = append(issues, one.Issues...)
		}
	}
	return validation.CheckResult{Issues: issues}
}

func bidAt(req *view.RequestView, bid *openrtb.Bid, responseCur string, seat *string, path string, ctx matchContext) validation.CheckResult {
	var issues []validation.CheckIssue
	if req == nil || bid == nil {
		return validation.CheckResult{Issues: []validation.CheckIssue{{
			Code: validation.CodeMalformed, Severity: validation.SeverityError, Path: path, Message: "nil req or bid",
		}}}
	}
	if price := bid.GetPrice(); math.IsNaN(price) || math.IsInf(price, 0) || price <= 0 {
		return validation.CheckResult{Issues: []validation.CheckIssue{{Code: validation.CodeMalformed, Severity: validation.SeverityError, Path: path + ".price", Message: "price must be finite and > 0"}}}
	}
	imp := ctx.findImp(bid.Impid)
	if imp == nil {
		return validation.CheckResult{Issues: []validation.CheckIssue{{
			Code: validation.CodeImpNotFound, Severity: validation.SeverityError, Path: path + ".impid",
			Message: "impid not found in request",
		}}}
	}
	mtype := int32(bid.GetMtype())
	markup := imp.Markup

	if mtype >= 500 {
		issues = append(issues, validation.CheckIssue{
			Code: validation.CodeMtypeVendor, Severity: validation.SeverityWarn, Path: path + ".mtype",
			Message: "vendor mtype >=500",
		})
	} else if mtype != 0 && (mtype < 1 || mtype > 4) {
		issues = append(issues, validation.CheckIssue{
			Code: validation.CodeMtypeUnknown, Severity: validation.SeverityError, Path: path + ".mtype",
			Message: "mtype must be 0–4 or >=500",
		})
	} else if mtype == 0 && markup.Count() > 1 {
		issues = append(issues, validation.CheckIssue{
			Code: validation.CodeMtypeRequired, Severity: validation.SeverityError, Path: path + ".mtype",
			Message: "multi-format Imp requires Bid.mtype",
		})
	} else if mtype >= 1 && mtype <= 4 {
		chosen := view.MarkupFromMtype(openrtb.MarkupType(mtype))
		if !markup.Has(chosen) {
			issues = append(issues, validation.CheckIssue{
				Code: validation.CodeMtypeMismatch, Severity: validation.SeverityError, Path: path + ".mtype",
				Message: "mtype does not match Imp formats",
			})
		}
	}

	eff := view.MarkupNone
	if mtype >= 1 && mtype <= 4 {
		eff = view.MarkupFromMtype(openrtb.MarkupType(mtype))
	} else if markup.Count() == 1 {
		eff = markup
	}

	issues = append(issues, checkFloor(imp, bid, responseCur, seat, path)...)
	issues = append(issues, checkAttr(imp, bid, eff, path)...)
	issues = append(issues, checkBlocks(ctx, bid, path)...)
	return validation.CheckResult{Issues: issues}
}

func checkFloor(imp *view.ImpView, bid *openrtb.Bid, responseCur string, seat *string, path string) []validation.CheckIssue {
	floor := imp.BidFloor
	currency := imp.BidFloorCur
	fixed := false
	fail := func(code, suffix, message string) []validation.CheckIssue {
		return []validation.CheckIssue{{Code: code, Severity: validation.SeverityError, Path: path + suffix, Message: message}}
	}
	if bid.Dealid == "" {
		if imp.Pmp != nil && imp.Pmp.GetPrivateAuction() == 1 {
			return fail(validation.CodeDealRequired, ".dealid", "private auction requires a deal")
		}
	} else {
		var deal *openrtb.Deal
		if imp.Pmp != nil {
			for _, candidate := range imp.Pmp.Deals {
				if candidate != nil && candidate.Id == bid.Dealid {
					deal = candidate
					break
				}
			}
		}
		if deal == nil {
			return fail(validation.CodeDealNotFound, ".dealid", "dealid not found on impression")
		}
		if seat != nil && len(deal.Wseat) > 0 {
			allowed := false
			for _, value := range deal.Wseat {
				if value == *seat {
					allowed = true
					break
				}
			}
			if !allowed {
				return fail(validation.CodeDealSeatNotAllowed, ".dealid", "seat not allowed by deal")
			}
		}
		if len(deal.Wadomain) > 0 {
			allowed := len(bid.Adomain) > 0
			for _, domain := range bid.Adomain {
				found := false
				for _, value := range deal.Wadomain {
					if asciiLower(domain) == asciiLower(value) {
						found = true
						break
					}
				}
				allowed = allowed && found
			}
			if !allowed {
				return fail(validation.CodeDealAdomainNotAllowed, ".adomain", "advertiser domain not allowed by deal")
			}
		}
		if deal.Bidfloor != nil {
			floor = deal.GetBidfloor()
			currency = deal.Bidfloorcur
		}
		fixed = deal.GetAt() == 3 && deal.Bidfloor != nil
	}
	if floor <= 0 && !fixed {
		return nil
	}
	// Require both response cur and bidfloorcur before numeric comparison (bidcheck.Bid has no response cur).
	respCur := strings.TrimSpace(responseCur)
	floorCur := strings.TrimSpace(currency)
	if respCur == "" || floorCur == "" {
		return nil
	}
	if asciiLower(respCur) != asciiLower(floorCur) {
		return []validation.CheckIssue{{
			Code: validation.CodeFloorCurDiff, Severity: validation.SeverityWarn, Path: path + ".price",
			Message: "bid currency differs from applicable bidfloorcur; skip floor compare",
		}}
	}
	if fixed && bid.GetPrice() != floor {
		return fail(validation.CodeDealPriceMismatch, ".price", "price differs from fixed deal price")
	}
	if bid.GetPrice() < floor {
		return []validation.CheckIssue{{
			Code: validation.CodePriceBelowFloor, Severity: validation.SeverityWarn, Path: path + ".price",
			Message: "price below applicable bidfloor",
		}}
	}
	return nil
}

func checkAttr(imp *view.ImpView, bid *openrtb.Bid, eff view.MarkupMask, path string) []validation.CheckIssue {
	if len(bid.Attr) == 0 || eff == view.MarkupNone {
		return nil
	}
	battr := battrFor(imp, eff)
	if len(battr) == 0 {
		return nil
	}
	blocked := map[int32]struct{}{}
	for _, v := range battr {
		blocked[v] = struct{}{}
	}
	for _, a := range bid.Attr {
		if _, ok := blocked[a]; ok {
			return []validation.CheckIssue{{
				Code: validation.CodeAttrBlocked, Severity: validation.SeverityWarn, Path: path + ".attr",
				Message: "bid.attr intersects format battr",
			}}
		}
	}
	return nil
}

func battrFor(imp *view.ImpView, eff view.MarkupMask) []int32 {
	if eff.HasBanner() && imp.Banner != nil {
		return imp.Banner.Battr
	}
	if eff.HasVideo() && imp.Video != nil {
		return imp.Video.Battr
	}
	if eff.HasAudio() && imp.Audio != nil {
		return imp.Audio.Battr
	}
	if eff.HasNative() && imp.Native != nil {
		return imp.Native.Battr
	}
	return nil
}

func checkBlocks(ctx matchContext, bid *openrtb.Bid, path string) []validation.CheckIssue {
	var issues []validation.CheckIssue
	if len(ctx.badv) > 0 {
		block := ctx.badv
		for _, d := range bid.Adomain {
			if _, ok := block[asciiLower(d)]; ok {
				issues = append(issues, validation.CheckIssue{
					Code: validation.CodeAdomainBlocked, Severity: validation.SeverityWarn, Path: path + ".adomain",
					Message: "adomain hit BidRequest.badv",
				})
				break
			}
		}
	}
	if len(ctx.bapp) > 0 && strings.TrimSpace(bid.Bundle) != "" {
		block := ctx.bapp
		if _, ok := block[asciiLower(strings.TrimSpace(bid.Bundle))]; ok {
			issues = append(issues, validation.CheckIssue{
				Code: validation.CodeBundleBlocked, Severity: validation.SeverityWarn, Path: path + ".bundle",
				Message: "bundle hit BidRequest.bapp",
			})
		}
	}
	if len(ctx.bcat) > 0 {
		block := ctx.bcat
		for _, c := range bid.Cat {
			if _, ok := block[asciiLower(c)]; ok {
				issues = append(issues, validation.CheckIssue{
					Code: validation.CodeCatBlocked, Severity: validation.SeverityWarn, Path: path + ".cat",
					Message: "cat hit BidRequest.bcat",
				})
				break
			}
		}
	}
	return issues
}

func lowerSet(in []string) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for _, s := range in {
		out[asciiLower(s)] = struct{}{}
	}
	return out
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

// Per-call indexes preserve the public snapshots' ownership contract and cannot go stale across calls.
type matchContext struct {
	single           *view.ImpView
	imps             map[string]*view.ImpView
	badv, bapp, bcat map[string]struct{}
}

func newMatchContext(req *view.RequestView) matchContext {
	if req == nil {
		return matchContext{}
	}
	ctx := matchContext{badv: lowerSet(req.Request().Badv), bapp: lowerSet(req.Request().Bapp), bcat: lowerSet(req.Request().Bcat)}
	if len(req.Imps()) == 1 {
		ctx.single = &req.Imps()[0]
		return ctx
	}
	ctx.imps = make(map[string]*view.ImpView, len(req.Imps()))
	for i := range req.Imps() {
		imp := &req.Imps()[i]
		if _, exists := ctx.imps[imp.ID]; !exists {
			ctx.imps[imp.ID] = imp
		}
	}
	return ctx
}

func (ctx matchContext) findImp(id string) *view.ImpView {
	if ctx.single != nil {
		if ctx.single.ID == id {
			return ctx.single
		}
		return nil
	}
	return ctx.imps[id]
}

func asciiLower(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return r
	}, s)
}
