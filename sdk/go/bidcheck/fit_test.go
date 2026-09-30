package bidcheck

import (
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"google.golang.org/protobuf/proto"
	"testing"

	"github.com/oakrtb/oakrtb/sdk/go/builder"
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/view"
)

func bannerSnap(t *testing.T) *view.RequestView {
	t.Helper()
	req, err := builder.NewBidRequest("a1").
		FirstPrice().
		Currency("USD").
		Site(builder.NewSite().ID("s1").Build()).
		AddImp(builder.NewBannerImp("1").Size(300, 250).Floor(1.0, "USD").Build()).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestImpReadyVideoMissingMimes(t *testing.T) {
	req := &openrtb.BidRequest{
		Id:  "x",
		At:  proto.Int32(1),
		Cur: []string{"USD"},
		Imp: []*openrtb.Imp{{Id: "1", Video: &openrtb.Video{}}},
	}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	r := validation.ImpReadyMtype(snap.Imps()[0].Imp, 2)
	if r.OK() || !r.Has(validation.CodeVideoMimesMissing) {
		t.Fatalf("got %+v", r)
	}
}

func TestMultiFormatNeedsMtype(t *testing.T) {
	req := &openrtb.BidRequest{
		Id:  "m",
		At:  proto.Int32(1),
		Cur: []string{"USD"},
		Imp: []*openrtb.Imp{{
			Id:     "1",
			Banner: &openrtb.Banner{W: proto.Int32(320), H: proto.Int32(50)},
			Video:  &openrtb.Video{Mimes: []string{"video/mp4"}},
		}},
	}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	bid := &openrtb.Bid{Id: "b1", Impid: "1", Price: proto.Float64(2)}
	r := Bid(snap, bid)
	if r.OK() || !r.Has(validation.CodeMtypeRequired) {
		t.Fatalf("got %+v", r)
	}
}

func TestMtypeMismatch(t *testing.T) {
	snap := bannerSnap(t)
	bid := &openrtb.Bid{
		Id: "b1", Impid: "1", Price: proto.Float64(2),
		Mtype: openrtb.MarkupType_MARKUP_TYPE_VIDEO,
	}
	r := Bid(snap, bid)
	if r.OK() || !r.Has(validation.CodeMtypeMismatch) {
		t.Fatalf("got %+v", r)
	}
}

func TestPriceBelowFloorWarn(t *testing.T) {
	snap := bannerSnap(t)
	bid := builder.NewBid("b1", "1", 0.5).Banner().Build()
	if Bid(snap, bid).Has(validation.CodePriceBelowFloor) {
		t.Fatal("竞价检查.Bid without response cur must skip floor compare")
	}
	res, err := builder.NewBidResponse("a1").Currency("USD").AddSeatBid("s1", bid).Build()
	if err != nil {
		t.Fatal(err)
	}
	r := Response(snap, res)
	if !r.OK() || !r.Has(validation.CodePriceBelowFloor) {
		t.Fatalf("got %+v", r)
	}
}

func TestAttrBlockedWarn(t *testing.T) {
	req := &openrtb.BidRequest{
		Id:  "x",
		At:  proto.Int32(1),
		Cur: []string{"USD"},
		Imp: []*openrtb.Imp{{
			Id:     "1",
			Banner: &openrtb.Banner{W: proto.Int32(300), H: proto.Int32(250), Battr: []int32{1}},
		}},
	}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	bid := &openrtb.Bid{
		Id: "b1", Impid: "1", Price: proto.Float64(2),
		Mtype: openrtb.MarkupType_MARKUP_TYPE_BANNER,
		Attr:  []int32{1},
	}
	r := Bid(snap, bid)
	if !r.OK() || !r.Has(validation.CodeAttrBlocked) {
		t.Fatalf("got %+v", r)
	}
}

func TestNoBidResponse(t *testing.T) {
	snap := bannerSnap(t)
	res, err := builder.NewBidResponse("a1").Currency("USD").NoBid(0).Build()
	if err != nil {
		t.Fatal(err)
	}
	r := Response(snap, res)
	if !r.OK() {
		t.Fatalf("got %+v", r)
	}
}

func TestImpidNotFound(t *testing.T) {
	snap := bannerSnap(t)
	bid := builder.NewBid("b1", "missing", 2).Banner().Build()
	r := Bid(snap, bid)
	if r.OK() || !r.Has(validation.CodeImpNotFound) {
		t.Fatalf("got %+v", r)
	}
}

func TestImpReadyMarkupBit(t *testing.T) {
	snap := bannerSnap(t)
	r := validation.ImpReadyMarkup(snap.Imps()[0].Imp, uint8(view.MarkupBanner))
	if !r.OK() {
		t.Fatalf("got %+v", r)
	}
}

func TestResponseHappy(t *testing.T) {
	snap := bannerSnap(t)
	res, err := builder.NewBidResponse("a1").
		Currency("USD").
		AddSeatBid("", builder.NewBid("b1", "1", 2).Banner().Adm("<a/>").Build()).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	r := Response(snap, res)
	if !r.OK() {
		t.Fatalf("got %+v", r)
	}
}

func TestFloorCurDiffSkipsCompare(t *testing.T) {
	req := &openrtb.BidRequest{
		Id: "a1", At: proto.Int32(1), Cur: []string{"USD"},
		Imp: []*openrtb.Imp{{
			Id: "1", Bidfloor: proto.Float64(1.0), Bidfloorcur: "EUR",
			Banner: &openrtb.Banner{W: proto.Int32(1), H: proto.Int32(1)},
		}},
	}
	snap, err := view.NewRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	bid := builder.NewBid("b1", "1", 0.1).Banner().Build()
	res, err := builder.NewBidResponse("a1").Currency("USD").AddSeatBid("s", bid).Build()
	if err != nil {
		t.Fatal(err)
	}
	r := Response(snap, res)
	if !r.OK() || !r.Has(validation.CodeFloorCurDiff) || r.Has(validation.CodePriceBelowFloor) {
		t.Fatalf("got %+v", r)
	}
}

func TestResponseCurWhitespaceTrimmed(t *testing.T) {
	snap := bannerSnap(t)
	bid := builder.NewBid("b1", "1", 0.5).Banner().Build()
	res, err := builder.NewBidResponse("a1").Currency("USD").AddSeatBid("s", bid).Build()
	if err != nil {
		t.Fatal(err)
	}
	res.Cur = "  USD  "
	r := Response(snap, res)
	if !r.OK() || r.Has(validation.CodeCurNotAllowed) || !r.Has(validation.CodePriceBelowFloor) {
		t.Fatalf("got %+v", r)
	}
}

func TestResponseNilMalformed(t *testing.T) {
	r := Response(nil, nil)
	if r.OK() || !r.Has(validation.CodeMalformed) {
		t.Fatalf("got %+v", r)
	}
}

func TestResponseEmptyBidMalformed(t *testing.T) {
	snap := bannerSnap(t)
	res := &openrtb.BidResponse{Id: "a1", Cur: "USD", Seatbid: []*openrtb.SeatBid{{Seat: "s"}}}
	r := Response(snap, res)
	if r.OK() || !r.Has(validation.CodeMalformed) {
		t.Fatalf("got %+v", r)
	}
}

func TestResponseBlankCurSkipsCurAndFloor(t *testing.T) {
	snap := bannerSnap(t)
	bid := builder.NewBid("b1", "1", 0.5).Banner().Build()
	res, err := builder.NewBidResponse("a1").Currency("USD").AddSeatBid("s", bid).Build()
	if err != nil {
		t.Fatal(err)
	}
	res.Cur = "   "
	r := Response(snap, res)
	if !r.OK() || r.Has(validation.CodeCurNotAllowed) || r.Has(validation.CodePriceBelowFloor) {
		t.Fatalf("got %+v", r)
	}
}

func TestCurNotAllowed(t *testing.T) {
	snap := bannerSnap(t)
	bid := builder.NewBid("b1", "1", 2).Banner().Build()
	res, err := builder.NewBidResponse("a1").Currency("EUR").AddSeatBid("s", bid).Build()
	if err != nil {
		t.Fatal(err)
	}
	r := Response(snap, res)
	if !r.OK() || !r.Has(validation.CodeCurNotAllowed) {
		t.Fatalf("got %+v", r)
	}
}
