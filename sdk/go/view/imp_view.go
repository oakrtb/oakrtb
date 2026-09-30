package view

import openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"

// ImpView is a zero-copy, read-only view of an Imp and its markup bitmask.
type ImpView struct {
	Imp *openrtb.Imp

	ID          string
	Markup      MarkupMask
	TagID       string
	BidFloor    float64
	BidFloorCur string
	Instl       int32
	Secure      int32
	Rwdd        int32
	Ssai        int32

	Banner *openrtb.Banner
	Video  *openrtb.Video
	Audio  *openrtb.Audio
	Native *openrtb.Native
	Pmp    *openrtb.Pmp
}

func viewImp(imp *openrtb.Imp) ImpView {
	return ImpView{
		Imp:         imp,
		ID:          imp.Id,
		Markup:      MarkupFromImp(imp),
		TagID:       imp.Tagid,
		BidFloor:    imp.GetBidfloor(),
		BidFloorCur: imp.Bidfloorcur,
		Instl:       imp.GetInstl(),
		Secure:      imp.GetSecure(),
		Rwdd:        imp.GetRwdd(),
		Ssai:        imp.GetSsai(),
		Banner:      imp.Banner,
		Video:       imp.Video,
		Audio:       imp.Audio,
		Native:      imp.Native,
		Pmp:         imp.Pmp,
	}
}
