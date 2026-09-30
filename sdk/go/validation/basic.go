// Package validation checks basic model structure without running JSON Schema.
package validation

import (
	"fmt"
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"math"
	"strings"
)

func Request(req *openrtb.BidRequest) error {
	if req == nil {
		return fmt.Errorf("validation: nil BidRequest")
	}
	if strings.TrimSpace(req.Id) == "" {
		return fmt.Errorf("validation: BidRequest.id is required")
	}
	if at := req.GetAt(); at != 1 && at != 2 && at < 500 {
		return fmt.Errorf("validation: BidRequest.at must be 1, 2, or >=500")
	}
	if len(req.Cur) == 0 {
		return fmt.Errorf("validation: BidRequest.cur is required")
	}
	for i, c := range req.Cur {
		if strings.TrimSpace(c) == "" {
			return fmt.Errorf("validation: BidRequest.cur[%d] is blank", i)
		}
	}
	if len(req.Imp) == 0 {
		return fmt.Errorf("validation: BidRequest.imp requires at least one Imp")
	}
	n := 0
	if req.Site != nil {
		n++
	}
	if req.App != nil {
		n++
	}
	if req.Dooh != nil {
		n++
	}
	if n > 1 {
		return fmt.Errorf("validation: site/app/dooh are mutually exclusive")
	}
	seen := make(map[string]struct{}, len(req.Imp))
	for i, imp := range req.Imp {
		if imp == nil {
			return fmt.Errorf("validation: imp[%d] is nil", i)
		}
		if strings.TrimSpace(imp.Id) == "" {
			return fmt.Errorf("validation: imp[%d].id is required", i)
		}
		if _, exists := seen[imp.Id]; exists {
			return fmt.Errorf("validation: duplicate imp.id %q", imp.Id)
		}
		seen[imp.Id] = struct{}{}
		if imp.Pmp != nil {
			deals := map[string]struct{}{}
			for _, d := range imp.Pmp.Deals {
				if d == nil || strings.TrimSpace(d.Id) == "" {
					return fmt.Errorf("validation: deal.id is required")
				}
				if _, exists := deals[d.Id]; exists {
					return fmt.Errorf("validation: duplicate deal.id %q", d.Id)
				}
				deals[d.Id] = struct{}{}
				if d.Bidfloor != nil && (math.IsNaN(d.GetBidfloor()) || math.IsInf(d.GetBidfloor(), 0) || d.GetBidfloor() < 0) {
					return fmt.Errorf("validation: deal.bidfloor must be finite and nonnegative")
				}
				if d.GetAt() == 3 && d.Bidfloor == nil {
					return fmt.Errorf("validation: fixed-price deal requires bidfloor")
				}
			}
		}
		if imp.Banner == nil && imp.Video == nil && imp.Audio == nil && imp.Native == nil {
			return fmt.Errorf("validation: imp[%d] needs banner, video, audio, or native", i)
		}
	}
	return nil
}
func Response(res *openrtb.BidResponse) error {
	if res == nil {
		return fmt.Errorf("validation: nil BidResponse")
	}
	if strings.TrimSpace(res.Id) == "" {
		return fmt.Errorf("validation: BidResponse.id is required")
	}
	if strings.TrimSpace(res.Cur) == "" {
		return fmt.Errorf("validation: BidResponse.cur is required")
	}
	if len(res.Seatbid) == 0 {
		return nil // Equivalent no-bid / empty-body forms.
	}
	for i, sb := range res.Seatbid {
		if sb == nil {
			return fmt.Errorf("validation: seatbid[%d] is nil", i)
		}
		if len(sb.Bid) == 0 {
			return fmt.Errorf("validation: seatbid[%d] needs at least one bid", i)
		}
		for j, bid := range sb.Bid {
			if bid == nil {
				return fmt.Errorf("validation: seatbid[%d].bid[%d] is nil", i, j)
			}
			if strings.TrimSpace(bid.Id) == "" || strings.TrimSpace(bid.Impid) == "" {
				return fmt.Errorf("validation: seatbid[%d].bid[%d] requires id and impid", i, j)
			}
			if math.IsNaN(bid.GetPrice()) || math.IsInf(bid.GetPrice(), 0) || bid.GetPrice() <= 0 {
				return fmt.Errorf("validation: seatbid[%d].bid[%d].price must be finite and > 0", i, j)
			}
		}
	}
	return nil
}
