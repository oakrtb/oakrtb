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

// BidRequestBuilder assembles a BidRequest with inventory and impressions.
type BidRequestBuilder struct {
	req *openrtb.BidRequest
	err error
}

// NewBidRequest creates a request builder; id is required (auction identifier).
func NewBidRequest(id string) *BidRequestBuilder {
	b := &BidRequestBuilder{req: &openrtb.BidRequest{Id: id}}
	if strings.TrimSpace(id) == "" {
		b.err = fmt.Errorf("builder: BidRequest.id is required")
	}
	return b
}

// FirstPrice sets at=1 (first-price auction).
func (b *BidRequestBuilder) FirstPrice() *BidRequestBuilder {
	return b.AuctionType(int32(openrtb.AuctionType_AUCTION_TYPE_FIRST_PRICE))
}

// SecondPricePlus sets at=2 (the OpenRTB default auction type).
func (b *BidRequestBuilder) SecondPricePlus() *BidRequestBuilder {
	return b.AuctionType(int32(openrtb.AuctionType_AUCTION_TYPE_SECOND_PRICE_PLUS))
}

// AuctionType sets BidRequest.at.
func (b *BidRequestBuilder) AuctionType(at int32) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.At = proto.Int32(at)
	return b
}

// Tmax sets the timeout in milliseconds, including network latency.
func (b *BidRequestBuilder) Tmax(ms int32) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Tmax = proto.Int32(ms)
	return b
}

// Currency sets acceptable bid currencies (ISO-4217), replacing the existing list.
func (b *BidRequestBuilder) Currency(codes ...string) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Cur = append([]string{}, codes...)
	return b
}

// Test marks test traffic (test=1).
func (b *BidRequestBuilder) Test() *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Test = proto.Int32(1)
	return b
}

// Bcat sets blocked advertiser categories.
func (b *BidRequestBuilder) Bcat(cats ...string) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Bcat = append([]string{}, cats...)
	return b
}

// Badv sets blocked advertiser domains.
func (b *BidRequestBuilder) Badv(domains ...string) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Badv = append([]string{}, domains...)
	return b
}

// Site sets website inventory (clears app/dooh).
func (b *BidRequestBuilder) Site(site *openrtb.Site) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Site = site
	b.req.App = nil
	b.req.Dooh = nil
	return b
}

// App sets app inventory (clears site/dooh).
func (b *BidRequestBuilder) App(app *openrtb.App) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.App = app
	b.req.Site = nil
	b.req.Dooh = nil
	return b
}

// Dooh sets digital out-of-home inventory (clears site/app).
func (b *BidRequestBuilder) Dooh(dooh *openrtb.Dooh) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Dooh = dooh
	b.req.Site = nil
	b.req.App = nil
	return b
}

// Device sets device context.
func (b *BidRequestBuilder) Device(device *openrtb.Device) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Device = device
	return b
}

// User sets user/audience context.
func (b *BidRequestBuilder) User(user *openrtb.User) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.User = user
	return b
}

// Regs sets privacy/regulatory signals.
func (b *BidRequestBuilder) Regs(regs *openrtb.Regs) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Regs = regs
	return b
}

// Source sets the upstream source/supply chain.
func (b *BidRequestBuilder) Source(source *openrtb.Source) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	b.req.Source = source
	return b
}

// AddImp appends an impression opportunity.
func (b *BidRequestBuilder) AddImp(imp *openrtb.Imp) *BidRequestBuilder {
	if b.err != nil {
		return b
	}
	if imp == nil {
		b.err = fmt.Errorf("builder: nil Imp")
		return b
	}
	b.req.Imp = append(b.req.Imp, imp)
	return b
}

// Build returns a BidRequest after structural validation (requires id and at least one imp with a format).
func (b *BidRequestBuilder) Build() (*openrtb.BidRequest, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := validation.Request(b.req); err != nil {
		return nil, err
	}
	for i, imp := range b.req.Imp {
		if err := checkImp(imp, i); err != nil {
			return nil, err
		}
	}
	return proto.Clone(b.req).(*openrtb.BidRequest), nil
}

// MustBuild behaves like Build but panics on error.
func (b *BidRequestBuilder) MustBuild() *openrtb.BidRequest {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

// BuildJSON builds and serializes OpenRTB JSON.
func (b *BidRequestBuilder) BuildJSON() ([]byte, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSON(req)
}

// BuildValidated builds, serializes, and validates against JSON Schema.
func (b *BidRequestBuilder) BuildValidated() ([]byte, jsonschema.Report, error) {
	raw, err := b.BuildJSON()
	if err != nil {
		return nil, jsonschema.Report{}, err
	}
	return raw, jsonschema.Request(raw), nil
}

func checkImp(imp *openrtb.Imp, i int) error {
	result := validation.ImpReady(imp)
	if errors := result.Errors(); len(errors) > 0 {
		return fmt.Errorf("builder: imp[%d]: %s", i, errors[0].Message)
	}
	return nil
}
