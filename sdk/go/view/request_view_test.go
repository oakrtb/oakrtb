package view

import (
	"testing"

	"github.com/oakrtb/oakrtb/sdk/go/builder"
)

func TestRequestViewFacts(t *testing.T) {
	req, err := builder.NewBidRequest("auction-1").
		FirstPrice().
		Tmax(120).
		Currency("USD").
		Site(builder.NewSite().ID("s1").Domain("example.com").Page("https://example.com/a").Build()).
		AddImp(builder.NewBannerImp("1").Size(300, 250).Floor(0.03, "USD").Secure().Build()).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	snap, err := NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if snap.AuctionID() != "auction-1" || snap.Inventory() != InventorySite {
		t.Fatalf("%+v", snap)
	}
	if snap.FindImp("1") == nil {
		t.Fatal("missing imp")
	}
	if len(snap.ImpsWith(MarkupBanner)) != 1 {
		t.Fatal(snap.ImpsWith(MarkupBanner))
	}
	facts := snap.Facts()
	if len(facts) != 1 || !facts[0].HasBanner() || facts[0].Mtype != 1 {
		t.Fatalf("%+v", facts)
	}
	if facts[0].BannerW == nil || *facts[0].BannerW != 300 {
		t.Fatalf("w=%v", facts[0].BannerW)
	}
	if facts[0].BidFloor != 0.03 || facts[0].Secure != 1 {
		t.Fatalf("%+v", facts[0])
	}
}

func TestRequestViewConstruction(t *testing.T) {
	req, err := builder.NewBidRequest("x").FirstPrice().Currency("USD").
		AddImp(builder.NewBannerImp("1").Size(1, 1).Build()).Build()
	if err != nil {
		t.Fatal(err)
	}
	snap, err := NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if snap.AuctionID() != "x" {
		t.Fatal(snap.AuctionID())
	}
}
