package view

import (
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"google.golang.org/protobuf/proto"
	"time"
)

// NewRequest validates basic structure and borrows a model for read-only queries.
// The model and all returned pointers/slices must remain unmodified while the view is used.
func NewRequest(req *openrtb.BidRequest) (*RequestView, error) {
	if err := validation.Request(req); err != nil {
		return nil, err
	}
	views := make([]ImpView, 0, len(req.Imp))
	for _, imp := range req.Imp {
		views = append(views, viewImp(imp))
	}
	var deadline time.Time
	if req.GetTmax() > 0 {
		deadline = time.Now().Add(time.Duration(int64(req.GetTmax())*85/100) * time.Millisecond)
	}
	return &RequestView{request: req, deadline: deadline, imps: views}, nil
}

// NewRequestCopy isolates the view from subsequent mutations to the original model.
func NewRequestCopy(input *openrtb.BidRequest) (*RequestView, error) {
	if input == nil {
		return NewRequest(nil)
	}
	return NewRequest(proto.Clone(input).(*openrtb.BidRequest))
}

// RequestView is a request query view that has passed basic validation.
type RequestView struct {
	request  *openrtb.BidRequest
	deadline time.Time
	imps     []ImpView
}

// AuctionID returns the auction ID (BidRequest.id).
func (s *RequestView) AuctionID() string { return s.request.Id }

// AuctionType returns the auction type (BidRequest.at).
func (s *RequestView) AuctionType() int32 { return s.request.GetAt() }

// Inventory returns the inventory type (site/app/dooh), distinct from proto Content.Channel.
func (s *RequestView) Inventory() Inventory {
	switch {
	case s.request.Site != nil:
		return InventorySite
	case s.request.App != nil:
		return InventoryApp
	case s.request.Dooh != nil:
		return InventoryDooh
	}
	return InventoryNone
}

// Currencies returns acceptable currencies (BidRequest.cur).
func (s *RequestView) Currencies() []string { return s.request.Cur }

// PastDeadline reports whether the view's time budget has expired.
func (s *RequestView) PastDeadline() bool {
	return !s.deadline.IsZero() && time.Now().After(s.deadline)
}

// Request returns the underlying BidRequest pointer.
func (s *RequestView) Request() *openrtb.BidRequest { return s.request }

// FindImp finds an ImpView by ID, returning nil if absent.
func (s *RequestView) FindImp(id string) *ImpView {
	for i := range s.imps {
		if s.imps[i].ID == id {
			return &s.imps[i]
		}
	}
	return nil
}

// ImpsWith returns ImpViews containing the specified markup flag.
func (s *RequestView) ImpsWith(flag MarkupMask) []ImpView {
	var out []ImpView
	for _, iv := range s.imps {
		if iv.Markup.Has(flag) {
			out = append(out, iv)
		}
	}
	return out
}

// Facts returns compact per-imp summaries for matching/bidding.
func (s *RequestView) Facts() []ImpFact {
	out := make([]ImpFact, 0, len(s.imps))
	for _, iv := range s.imps {
		out = append(out, factFrom(iv))
	}
	return out
}

// ImpFact is a flattened ImpView summary for callers.
type ImpFact struct {
	ID            string
	Markup        MarkupMask
	Mtype         int32
	TagID         string
	BidFloor      float64
	BidFloorCur   string
	Secure        int32
	Instl         int32
	Rwdd          int32
	Ssai          int32
	BannerW       *int32
	BannerH       *int32
	NativeRequest string // Empty when native is absent.
}

func factFrom(iv ImpView) ImpFact {
	f := ImpFact{
		ID:          iv.ID,
		Markup:      iv.Markup,
		Mtype:       iv.Markup.Mtype(),
		TagID:       iv.TagID,
		BidFloor:    iv.BidFloor,
		BidFloorCur: iv.BidFloorCur,
		Secure:      iv.Secure,
		Instl:       iv.Instl,
		Rwdd:        iv.Rwdd,
		Ssai:        iv.Ssai,
	}
	if iv.Banner != nil {
		if iv.Banner.W != nil {
			w := iv.Banner.GetW()
			f.BannerW = &w
		}
		if iv.Banner.H != nil {
			h := iv.Banner.GetH()
			f.BannerH = &h
		}
	}
	if iv.Native != nil && iv.Native.Request != "" {
		f.NativeRequest = iv.Native.Request
	}
	return f
}

// HasBanner reports whether ImpFact contains Banner markup.
func (f ImpFact) HasBanner() bool { return f.Markup.HasBanner() }

// HasVideo reports whether ImpFact contains Video markup.
func (f ImpFact) HasVideo() bool { return f.Markup.HasVideo() }

// HasAudio reports whether ImpFact contains Audio markup.
func (f ImpFact) HasAudio() bool { return f.Markup.HasAudio() }

// HasNative reports whether ImpFact contains Native markup.
func (f ImpFact) HasNative() bool { return f.Markup.HasNative() }

// Imps returns borrowed read-only projections.
func (s *RequestView) Imps() []ImpView { return s.imps }

// Deadline is fixed when the view is created.
func (s *RequestView) Deadline() time.Time { return s.deadline }
