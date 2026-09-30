package builder

import (
	rt "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestBuildResultsAreIndependent(t *testing.T) {
	imp := NewBannerImp("1").Size(300, 250).Build()
	b := NewBidRequest("a").FirstPrice().Currency("USD").AddImp(imp)
	first, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	second, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	b.Currency("EUR")
	imp.Banner.W = proto.Int32(900)
	second.Imp[0].Id = "changed"
	if first.Cur[0] != "USD" || first.Imp[0].Banner.GetW() != 300 || first.Imp[0].Id != "1" {
		t.Fatal("request result aliases builder, input or another result")
	}
	bid := NewBid("b", "1", 1).Adomain("original.com").Build()
	rb := NewBidResponse("a").Currency("USD").AddSeatBid("s", bid)
	response, err := rb.Build()
	if err != nil {
		t.Fatal(err)
	}
	rb.Currency("EUR").NoBid(0)
	bid.Adomain[0] = "changed.com"
	if response.Cur != "USD" || len(response.Seatbid) != 1 || response.Seatbid[0].Bid[0].Adomain[0] != "original.com" {
		t.Fatal("response result aliases builder or input")
	}
}

func TestComponentBuildersReturnIndependentModels(t *testing.T) {
	b := NewBid("b", "1", 1)
	first := b.Build()
	b.Adm("changed")
	if first.Adm != "" {
		t.Fatal("bid aliases builder")
	}
	ib := NewBannerImp("1").Size(1, 1)
	i := ib.Build()
	ib.Size(2, 2)
	if i.Banner.GetW() != 1 {
		t.Fatal("imp aliases builder")
	}
	pub := &rt.Publisher{Id: "original"}
	sb := NewSite().Publisher(pub)
	site := sb.Build()
	pub.Id = "changed"
	sb.ID("later")
	if site.Publisher.Id != "original" || site.Id != "" {
		t.Fatal("site aliases nested input or builder")
	}
}
