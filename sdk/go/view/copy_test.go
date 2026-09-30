package view

import (
	rt "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestRequestCopyIsolated(t *testing.T) {
	req := &rt.BidRequest{Id: "a", At: proto.Int32(1), Cur: []string{"USD"}, Imp: []*rt.Imp{{Id: "1", Banner: &rt.Banner{W: proto.Int32(300)}}}}
	snap, err := NewRequestCopy(req)
	if err != nil {
		t.Fatal(err)
	}
	req.Id = "changed"
	req.Cur[0] = "EUR"
	req.Imp[0].Id = "changed"
	*req.Imp[0].Banner.W = 900
	if snap.AuctionID() != "a" || snap.Request().Id != "a" || snap.Currencies()[0] != "USD" || snap.FindImp("1").Banner.GetW() != 300 {
		t.Fatal("snapshot changed with source")
	}
}
func TestResponseCopyIsolated(t *testing.T) {
	res := &rt.BidResponse{Id: "a", Cur: "USD", Seatbid: []*rt.SeatBid{{Bid: []*rt.Bid{{Id: "b", Impid: "1", Price: proto.Float64(1), Adomain: []string{"example.com"}}}}}}
	snap, err := NewResponseCopy(res)
	if err != nil {
		t.Fatal(err)
	}
	res.Id = "changed"
	*res.Seatbid[0].Bid[0].Price = 9
	res.Seatbid[0].Bid[0].Adomain[0] = "changed"
	if snap.RequestID() != "a" || snap.Response().Id != "a" || snap.FindBid("b").Bid.GetPrice() != 1 || snap.FindBid("b").Adomain[0] != "example.com" {
		t.Fatal("snapshot changed with source")
	}
}
func TestCopyNil(t *testing.T) {
	if _, err := NewRequestCopy(nil); err == nil {
		t.Fatal("nil request accepted")
	}
	if _, err := NewResponseCopy(nil); err == nil {
		t.Fatal("nil response accepted")
	}
}

func TestSeatViewsShareFlatStorage(t *testing.T) {
	res := &rt.BidResponse{Id: "a", Cur: "USD", Seatbid: []*rt.SeatBid{
		{Seat: "first", Bid: []*rt.Bid{{Id: "b1", Impid: "1", Price: proto.Float64(1)}}},
		{Seat: "second", Bid: []*rt.Bid{{Id: "b2", Impid: "1", Price: proto.Float64(2)}}},
	}}
	snap, err := NewResponse(res)
	if err != nil {
		t.Fatal(err)
	}
	if &snap.SeatBids()[0].Bids()[0] != &snap.Bids()[0] || &snap.SeatBids()[1].Bids()[0] != &snap.Bids()[1] {
		t.Fatal("bid views were copied")
	}
	extended := append(snap.SeatBids()[0].Bids(), BidView{ID: "new"})
	if extended[1].ID != "new" || snap.SeatBids()[1].Bids()[0].ID != "b2" {
		t.Fatal("append overwrote neighbouring seat")
	}
}
