package view

import openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"

// SeatBidView contains a SeatBid and its nested BidView list.
type SeatBidView struct {
	SeatBid *openrtb.SeatBid

	Seat  string
	Group int32
	bids  []BidView
}

// BidView is a zero-copy, read-only view of a Bid and the markup inferred from mtype.
type BidView struct {
	Bid *openrtb.Bid

	ID      string
	ImpID   string
	Seat    string // Parent seat.
	Price   float64
	Markup  MarkupMask // Derived from Bid.mtype.
	Mtype   int32      // Raw MarkupType value, 0–4.
	Crid    string
	Cid     string
	DealID  string
	W       int32
	H       int32
	Dur     int32
	Adm     string
	NURL    string
	BURL    string
	LURL    string
	Adomain []string
}

func viewSeatBids(res *openrtb.BidResponse) ([]SeatBidView, []BidView) {
	seats := make([]SeatBidView, 0, len(res.Seatbid))
	total := 0
	for _, sb := range res.Seatbid {
		total += len(sb.Bid)
	}
	flat := make([]BidView, 0, total)
	for _, sb := range res.Seatbid {
		start := len(flat)
		for _, bid := range sb.Bid {
			flat = append(flat, viewBid(bid, sb.Seat))
		}
		end := len(flat)
		// Cap each seat slice so appending cannot overwrite a neighbouring seat.
		seats = append(seats, SeatBidView{SeatBid: sb, Seat: sb.Seat, Group: sb.GetGroup(), bids: flat[start:end:end]})
	}
	return seats, flat
}

func viewBid(bid *openrtb.Bid, seat string) BidView {
	return BidView{
		Bid:     bid,
		ID:      bid.Id,
		ImpID:   bid.Impid,
		Seat:    seat,
		Price:   bid.GetPrice(),
		Markup:  MarkupFromMtype(bid.Mtype),
		Mtype:   int32(bid.Mtype),
		Crid:    bid.Crid,
		Cid:     bid.Cid,
		DealID:  bid.Dealid,
		W:       bid.GetW(),
		H:       bid.GetH(),
		Dur:     bid.GetDur(),
		Adm:     bid.Adm,
		NURL:    bid.Nurl,
		BURL:    bid.Burl,
		LURL:    bid.Lurl,
		Adomain: bid.Adomain,
	}
}

// Bids returns borrowed read-only projections for this seat.
func (s SeatBidView) Bids() []BidView { return s.bids }
