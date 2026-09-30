package main

import (
	"fmt"
	"os"

	"github.com/oakrtb/oakrtb/sdk/go/bidcheck"
	"github.com/oakrtb/oakrtb/sdk/go/builder"
	"github.com/oakrtb/oakrtb/sdk/go/codec"
	"github.com/oakrtb/oakrtb/sdk/go/jsonschema"
	"github.com/oakrtb/oakrtb/sdk/go/view"
)

// respond simulates the DSP boundary; price == 0 means no candidate was found.
func respond(requestJSON []byte, price float64) ([]byte, error) {
	// Full contract validation is useful at integration boundaries.
	if report := jsonschema.Request(requestJSON); !report.Ok {
		return nil, fmt.Errorf("invalid request: %+v", report.Errors)
	}
	req, err := codec.UnmarshalBidRequest(requestJSON)
	if err != nil {
		return nil, err
	}
	// The borrowed request must remain read-only for the lifetime of this view.
	requestView, err := view.NewRequest(req)
	if err != nil {
		return nil, err
	}
	imp := requestView.FindImp("imp-1")
	if imp == nil || !imp.Markup.HasBanner() {
		return nil, fmt.Errorf("expected Banner imp-1")
	}

	response := builder.NewBidResponse(requestView.AuctionID()).Currency("USD")
	if price == 0 {
		// Reason 0 means unknown; this demo has no more specific standard reason.
		response.NoBid(0)
	} else {
		bid := builder.NewBid("bid-1", imp.ID, price).Banner().Size(300, 250).
			Crid("creative-1").Adomain("advertiser.example").
			Adm(`<a href="https://advertiser.example/"><img src="https://cdn.example/banner.png" width="300" height="250"></a>`).Build()
		response.AddSeatBid("buyer-1", bid)
	}
	res, err := response.Build()
	if err != nil {
		return nil, err
	}
	findings := bidcheck.Response(requestView, res)
	if !findings.OK() {
		return nil, fmt.Errorf("response mismatch: %+v", findings.Errors())
	}
	// Demo policy: reject any warning, even when OK() is true.
	if warnings := findings.Warnings(); len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "policy rejected bid: %+v\n", warnings)
		res, err = builder.NewBidResponse(requestView.AuctionID()).Currency("USD").NoBid(0).Build()
		if err != nil {
			return nil, err
		}
	}
	payload, err := codec.MarshalJSON(res)
	if err != nil {
		return nil, err
	}
	if report := jsonschema.Response(payload); !report.Ok {
		return nil, fmt.Errorf("invalid response: %+v", report.Errors)
	}
	return payload, nil
}

func run() error {
	// SSP: construct an auction request and serialize it for the DSP.
	req, err := builder.NewBidRequest("auction-1").FirstPrice().Currency("USD").Test().
		Site(builder.NewSite().ID("site-1").Domain("publisher.example").Build()).
		AddImp(builder.NewBannerImp("imp-1").Size(300, 250).Floor(1, "USD").Secure().Build()).Build()
	if err != nil {
		return err
	}
	requestJSON, err := codec.MarshalJSON(req)
	if err != nil {
		return err
	}
	fmt.Printf("request: %s\n", requestJSON)
	for _, price := range []float64{2, 0.5, 0} {
		responseJSON, err := respond(requestJSON, price)
		if err != nil {
			return err
		}
		// SSP: decode the returned response and inspect its query view.
		res, err := codec.UnmarshalBidResponse(responseJSON)
		if err != nil {
			return err
		}
		responseView, err := view.NewResponse(res)
		if err != nil {
			return err
		}
		if responseView.NoBid() != (price < 1) {
			return fmt.Errorf("unexpected decision for price %g", price)
		}
		fmt.Printf("price=%g no_bid=%t response=%s\n", price, responseView.NoBid(), responseJSON)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
