package builder

import (
	"fmt"
	"github.com/oakrtb/oakrtb/sdk/go/codec"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"google.golang.org/protobuf/proto"
	"strings"

	"github.com/oakrtb/oakrtb/sdk/go/jsonschema"
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
)

// BidResponseBuilder assembles a BidResponse (bids or a structured no-bid).
type BidResponseBuilder struct {
	res      *openrtb.BidResponse
	noBidSet bool
	err      error
}

// NewBidResponse creates a response builder for the given BidRequest.id.
func NewBidResponse(requestID string) *BidResponseBuilder {
	b := &BidResponseBuilder{res: &openrtb.BidResponse{Id: requestID, Cur: "USD"}}
	if strings.TrimSpace(requestID) == "" {
		b.err = fmt.Errorf("builder: BidResponse.id is required (echo BidRequest.id)")
	}
	return b
}

// BidID sets the buyer-side response ID.
func (b *BidResponseBuilder) BidID(bidid string) *BidResponseBuilder {
	if b.err != nil {
		return b
	}
	b.res.Bidid = bidid
	return b
}

// Currency sets BidResponse.cur (ISO-4217).
func (b *BidResponseBuilder) Currency(cur string) *BidResponseBuilder {
	if b.err != nil {
		return b
	}
	b.res.Cur = cur
	return b
}

// NoBid sets a structured no-bid reason (nbr) and clears seatbid.
func (b *BidResponseBuilder) NoBid(nbr int32) *BidResponseBuilder {
	if b.err != nil {
		return b
	}
	b.res.Nbr = proto.Int32(nbr)
	b.res.Seatbid = nil
	b.noBidSet = true
	return b
}

// AddSeatBid appends bids for a seat and clears any previous NoBid state.
func (b *BidResponseBuilder) AddSeatBid(seat string, bids ...*openrtb.Bid) *BidResponseBuilder {
	if b.err != nil {
		return b
	}
	if len(bids) == 0 {
		b.err = fmt.Errorf("builder: SeatBid requires at least one Bid")
		return b
	}
	b.noBidSet = false
	b.res.Nbr = nil
	b.res.Seatbid = append(b.res.Seatbid, &openrtb.SeatBid{
		Seat: seat,
		Bid:  bids,
	})
	return b
}

// Build returns a BidResponse after structural validation.
func (b *BidResponseBuilder) Build() (*openrtb.BidResponse, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := validation.Response(b.res); err != nil {
		return nil, err
	}
	if len(b.res.Seatbid) == 0 && !b.noBidSet {
		return nil, fmt.Errorf("builder: BidResponse needs seatbid[] or noBid(nbr)")
	}
	return proto.Clone(b.res).(*openrtb.BidResponse), nil
}

// MustBuild behaves like Build but panics on error.
func (b *BidResponseBuilder) MustBuild() *openrtb.BidResponse {
	res, err := b.Build()
	if err != nil {
		panic(err)
	}
	return res
}

// BuildJSON builds and serializes OpenRTB JSON.
func (b *BidResponseBuilder) BuildJSON() ([]byte, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSON(res)
}

// BuildValidated builds, serializes, and validates against JSON Schema.
func (b *BidResponseBuilder) BuildValidated() ([]byte, jsonschema.Report, error) {
	raw, err := b.BuildJSON()
	if err != nil {
		return nil, jsonschema.Report{}, err
	}
	return raw, jsonschema.Response(raw), nil
}

// --- Bid factory -------------------------------------------------------------

// NewBid creates a single Bid (id, impid, and price are required).
func NewBid(id, impid string, price float64) *BidBuilder {
	return &BidBuilder{b: &openrtb.Bid{Id: id, Impid: impid, Price: proto.Float64(price)}}
}

// BidBuilder configures one Bid.
type BidBuilder struct{ b *openrtb.Bid }

// Adm sets the ad markup.
func (b *BidBuilder) Adm(adm string) *BidBuilder { b.b.Adm = adm; return b }

// NURL sets the win notice URL.
func (b *BidBuilder) NURL(u string) *BidBuilder { b.b.Nurl = u; return b }

// BURL sets the billing notice URL.
func (b *BidBuilder) BURL(u string) *BidBuilder { b.b.Burl = u; return b }

// LURL sets the loss notice URL.
func (b *BidBuilder) LURL(u string) *BidBuilder { b.b.Lurl = u; return b }

// Crid sets the creative ID.
func (b *BidBuilder) Crid(crid string) *BidBuilder {
	b.b.Crid = crid
	return b
}

// Cid sets the campaign ID.
func (b *BidBuilder) Cid(cid string) *BidBuilder { b.b.Cid = cid; return b }

// Adomain sets the advertiser domain list.
func (b *BidBuilder) Adomain(domains ...string) *BidBuilder {
	b.b.Adomain = append([]string{}, domains...)
	return b
}

// Size sets the creative width and height.
func (b *BidBuilder) Size(w, h int32) *BidBuilder {
	b.b.W = proto.Int32(w)
	b.b.H = proto.Int32(h)
	return b
}

// DealID sets the PMP deal ID.
func (b *BidBuilder) DealID(id string) *BidBuilder {
	b.b.Dealid = id
	return b
}

// MarkupType sets Bid.mtype.
func (b *BidBuilder) MarkupType(m openrtb.MarkupType) *BidBuilder {
	b.b.Mtype = m
	return b
}

// Banner sets mtype to Banner (1).
func (b *BidBuilder) Banner() *BidBuilder {
	return b.MarkupType(openrtb.MarkupType_MARKUP_TYPE_BANNER)
}

// Video sets mtype to Video (2).
func (b *BidBuilder) Video() *BidBuilder {
	return b.MarkupType(openrtb.MarkupType_MARKUP_TYPE_VIDEO)
}

// Audio sets mtype to Audio (3).
func (b *BidBuilder) Audio() *BidBuilder {
	return b.MarkupType(openrtb.MarkupType_MARKUP_TYPE_AUDIO)
}

// Native sets mtype to Native (4).
func (b *BidBuilder) Native() *BidBuilder {
	return b.MarkupType(openrtb.MarkupType_MARKUP_TYPE_NATIVE)
}

// Dur sets the creative duration in seconds.
func (b *BidBuilder) Dur(seconds int32) *BidBuilder { b.b.Dur = proto.Int32(seconds); return b }

// Build returns the configured Bid.
func (b *BidBuilder) Build() *openrtb.Bid { return proto.Clone(b.b).(*openrtb.Bid) }
