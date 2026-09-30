package view

import (
	"google.golang.org/protobuf/proto"
	"testing"

	"github.com/oakrtb/oakrtb/sdk/go/builder"
)

func TestResponseViewRunFacts(t *testing.T) {
	res, err := builder.NewBidResponse("auction-1").
		Currency("USD").
		AddSeatBid("512",
			builder.NewBid("1", "1", 1.23).Banner().Size(300, 250).Adm("<img/>").Crid("c1").Adomain("adv.com").Build(),
		).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	snap, err := NewResponse(res)
	if err != nil {
		t.Fatal(err)
	}
	if snap.RequestID() != "auction-1" || snap.NoBid() || snap.Currency() != "USD" {
		t.Fatalf("%+v", snap)
	}
	if snap.FindBid("1") == nil {
		t.Fatal("missing bid")
	}
	if len(snap.BidsForImp("1")) != 1 {
		t.Fatal(snap.BidsForImp("1"))
	}
	if len(snap.BidsWith(MarkupBanner)) != 1 {
		t.Fatal(snap.BidsWith(MarkupBanner))
	}
	facts := snap.Facts()
	if len(facts) != 1 || !facts[0].HasBanner() || facts[0].Mtype != 1 {
		t.Fatalf("%+v", facts)
	}
	if facts[0].Price != 1.23 || !facts[0].HasAdm || facts[0].W != 300 {
		t.Fatalf("%+v", facts[0])
	}
}

func TestResponseViewNoBid(t *testing.T) {
	res, err := builder.NewBidResponse("auction-1").NoBid(2).Build()
	if err != nil {
		t.Fatal(err)
	}
	snap, err := NewResponse(res)
	if err != nil {
		t.Fatal(err)
	}
	if !snap.NoBid() || snap.Nbr() != 2 || len(snap.Bids()) != 0 {
		t.Fatalf("%+v", snap)
	}
}

func TestResponseBasicValidationRejects(t *testing.T) {
	res, err := builder.NewBidResponse("x").
		AddSeatBid("s", builder.NewBid("1", "1", 1.0).Build()).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	res.Seatbid[0].Bid[0].Price = proto.Float64(0)
	if _, err := NewResponse(res); err == nil {
		t.Fatal("expected price error")
	}
}

func TestResponseBasicValidationRejectsBlankCur(t *testing.T) {
	res, err := builder.NewBidResponse("x").Currency("USD").NoBid(0).Build()
	if err != nil {
		t.Fatal(err)
	}
	res.Cur = "   "
	if _, err := NewResponse(res); err == nil {
		t.Fatal("expected blank cur")
	}
}
