package view

import (
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"google.golang.org/protobuf/proto"
)

// NewResponse validates basic structure and borrows a model for read-only queries.
// The model and all returned pointers/slices must remain unmodified while the view is used.
func NewResponse(res *openrtb.BidResponse) (*ResponseView, error) {
	if err := validation.Response(res); err != nil {
		return nil, err
	}
	seats, bids := viewSeatBids(res)
	return &ResponseView{response: res, seatBids: seats, bids: bids}, nil
}

// NewResponseCopy isolates the view from subsequent mutations to the original model.
func NewResponseCopy(input *openrtb.BidResponse) (*ResponseView, error) {
	if input == nil {
		return NewResponse(nil)
	}
	return NewResponse(proto.Clone(input).(*openrtb.BidResponse))
}

// ResponseView is a response query view that has passed basic validation.
type ResponseView struct {
	response *openrtb.BidResponse
	seatBids []SeatBidView
	bids     []BidView
}

// RequestID returns the corresponding request ID (BidResponse.id).
func (s *ResponseView) RequestID() string { return s.response.Id }

// BidID returns the buyer response ID (BidResponse.bidid).
func (s *ResponseView) BidID() string { return s.response.Bidid }

// Currency returns the response currency (BidResponse.cur).
func (s *ResponseView) Currency() string { return s.response.Cur }

// NoBid reports whether the response is a no-bid (no seatbid).
func (s *ResponseView) NoBid() bool { return len(s.response.Seatbid) == 0 }

// Nbr returns the structured no-bid reason code (BidResponse.nbr).
func (s *ResponseView) Nbr() int32 { return s.response.GetNbr() }

// Response returns the underlying BidResponse pointer.
func (s *ResponseView) Response() *openrtb.BidResponse { return s.response }

// FindBid finds a BidView by bid ID, returning nil if absent.
func (s *ResponseView) FindBid(id string) *BidView {
	for i := range s.bids {
		if s.bids[i].ID == id {
			return &s.bids[i]
		}
	}
	return nil
}

// BidsForImp returns all BidViews for the specified impid.
func (s *ResponseView) BidsForImp(impid string) []BidView {
	var out []BidView
	for _, b := range s.bids {
		if b.ImpID == impid {
			out = append(out, b)
		}
	}
	return out
}

// BidsWith returns BidViews containing the specified markup flag.
func (s *ResponseView) BidsWith(flag MarkupMask) []BidView {
	var out []BidView
	for _, b := range s.bids {
		if b.Markup.Has(flag) {
			out = append(out, b)
		}
	}
	return out
}

// Facts returns compact per-bid summaries for downstream processing.
func (s *ResponseView) Facts() []BidFact {
	out := make([]BidFact, 0, len(s.bids))
	for _, b := range s.bids {
		out = append(out, bidFactFrom(b))
	}
	return out
}

// BidFact is a flattened BidView summary for callers.
type BidFact struct {
	ID      string
	ImpID   string
	Seat    string
	Price   float64
	Markup  MarkupMask
	Mtype   int32
	Crid    string
	DealID  string
	W       int32
	H       int32
	Dur     int32
	HasAdm  bool
	Adomain []string
}

func bidFactFrom(b BidView) BidFact {
	return BidFact{
		ID:      b.ID,
		ImpID:   b.ImpID,
		Seat:    b.Seat,
		Price:   b.Price,
		Markup:  b.Markup,
		Mtype:   b.Mtype,
		Crid:    b.Crid,
		DealID:  b.DealID,
		W:       b.W,
		H:       b.H,
		Dur:     b.Dur,
		HasAdm:  b.Adm != "",
		Adomain: b.Adomain,
	}
}

// HasBanner reports whether BidFact uses Banner markup.
func (f BidFact) HasBanner() bool { return f.Markup.HasBanner() }

// HasVideo reports whether BidFact uses Video markup.
func (f BidFact) HasVideo() bool { return f.Markup.HasVideo() }

// HasAudio reports whether BidFact uses Audio markup.
func (f BidFact) HasAudio() bool { return f.Markup.HasAudio() }

// HasNative reports whether BidFact uses Native markup.
func (f BidFact) HasNative() bool { return f.Markup.HasNative() }

// SeatBids returns borrowed read-only seat projections.
func (s *ResponseView) SeatBids() []SeatBidView { return s.seatBids }

// Bids returns borrowed read-only bid projections.
func (s *ResponseView) Bids() []BidView { return s.bids }
