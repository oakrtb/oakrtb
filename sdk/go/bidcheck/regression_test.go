package bidcheck

import (
	"github.com/oakrtb/oakrtb/sdk/go/builder"
	rt "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"github.com/oakrtb/oakrtb/sdk/go/view"
	"google.golang.org/protobuf/proto"
	"math"
	"testing"
)

func TestResponseRequestID(t *testing.T) {
	req := bannerSnap(t)
	for _, id := range []string{"a1", "other", ""} {
		for _, noBid := range []bool{true, false} {
			res := &rt.BidResponse{Id: id, Cur: "USD"}
			if !noBid {
				res.Seatbid = []*rt.SeatBid{{Bid: []*rt.Bid{{Id: "b", Impid: "1", Price: proto.Float64(2)}}}}
			}
			result := Response(req, res)
			if result.Has(validation.CodeRequestIDMismatch) != (id != "a1") {
				t.Fatalf("id=%q noBid=%v: %+v", id, noBid, result)
			}
			if id != "a1" && result.OK() {
				t.Fatal("mismatched response accepted")
			}
		}
	}
}
func TestRejectInvalidPrices(t *testing.T) {
	req := bannerSnap(t)
	for _, price := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 0, -1} {
		bid := &rt.Bid{Id: "b", Impid: "1", Price: proto.Float64(price)}
		res := &rt.BidResponse{Id: "a1", Cur: "USD", Seatbid: []*rt.SeatBid{{Bid: []*rt.Bid{bid}}}}
		if _, err := view.NewResponse(res); err == nil {
			t.Fatalf("BasicValidation accepted %v", price)
		}
		if r := Bid(req, bid); r.OK() || !r.Has(validation.CodeMalformed) {
			t.Fatalf("竞价检查 accepted %v: %+v", price, r)
		}
		if r := Response(req, res); r.OK() {
			t.Fatalf("竞价检查.Response accepted %v", price)
		}
		if _, err := builder.NewBidResponse("a1").AddSeatBid("s", bid).Build(); err == nil {
			t.Fatalf("Builder accepted %v", price)
		}
	}
}

func TestIndexesPreserveMatchingAndArePerCall(t *testing.T) {
	req := &rt.BidRequest{Id: "a", At: proto.Int32(1), Cur: []string{"USD"}, Badv: []string{"BLOCKED.COM"}, Bcat: []string{"IAB25"}, Bapp: []string{"COM.BLOCKED"}, Imp: []*rt.Imp{
		{Id: "1", Banner: &rt.Banner{}}, {Id: "2", Video: &rt.Video{}},
	}}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	res := &rt.BidResponse{Id: "a", Cur: "USD", Seatbid: []*rt.SeatBid{{Bid: []*rt.Bid{
		{Id: "b1", Impid: "1", Price: proto.Float64(1), Mtype: rt.MarkupType_MARKUP_TYPE_BANNER, Adomain: []string{"blocked.com"}, Cat: []string{"iab25"}, Bundle: "com.blocked"},
		{Id: "b2", Impid: "2", Price: proto.Float64(1), Mtype: rt.MarkupType_MARKUP_TYPE_VIDEO},
	}}}}
	result := Response(snap, res)
	if !result.OK() || !result.Has(validation.CodeAdomainBlocked) || !result.Has(validation.CodeCatBlocked) || !result.Has(validation.CodeBundleBlocked) {
		t.Fatalf("unexpected indexed result: %+v", result)
	}
	// A separate snapshot with the same auction/imp IDs must not reuse prior blocklists.
	other := proto.Clone(req).(*rt.BidRequest)
	other.Badv = nil
	other.Bcat = nil
	other.Bapp = nil
	clean, err := view.NewRequest(other)
	if err != nil {
		t.Fatal(err)
	}
	if r := Response(clean, res); len(r.Issues) != 0 {
		t.Fatalf("stale indexes: %+v", r)
	}
}
