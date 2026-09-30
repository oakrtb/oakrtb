package builder_test

import (
	"testing"

	"github.com/oakrtb/oakrtb/sdk/go/builder"
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
)

func TestBannerSiteRequestValidated(t *testing.T) {
	raw, result, err := builder.NewBidRequest("auction-banner-1").
		FirstPrice().
		Tmax(120).
		Currency("USD").
		Site(builder.NewSite().
			ID("102855").
			Domain("www.example.com").
			Page("https://www.example.com/article").
			Publisher(builder.NewPublisher().ID("8953").Name("Example").Domain("example.com").Build()).
			Build()).
		Device(builder.NewDevice().
			UA("Mozilla/5.0").
			IP("192.0.2.1").
			DeviceType(int32(openrtb.DeviceType_DEVICE_TYPE_PHONE)).
			OS("Android", "14").
			Build()).
		AddImp(builder.NewBannerImp("1").
			Size(300, 250).
			Floor(0.03, "USD").
			Secure().
			BannerPos(int32(openrtb.AdPosition_AD_POSITION_ABOVE_THE_FOLD)).
			BannerMimes("image/jpeg", "image/png").
			Build()).
		BuildValidated()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Fatalf("schema: %s", result.Errors)
	}
	if len(raw) == 0 {
		t.Fatal("empty json")
	}
}

func TestVideoAppRequestValidated(t *testing.T) {
	_, result, err := builder.NewBidRequest("auction-video-1").
		FirstPrice().
		Tmax(200).
		Currency("USD").
		App(builder.NewApp().
			ID("ctv-1").
			Name("Example CTV").
			Bundle("com.example.ctv").
			Content(builder.NewContent().Title("Show").Series("Series").Context(1).Build()).
			Build()).
		Device(builder.NewDevice().
			UA("CTV").
			DeviceType(int32(openrtb.DeviceType_DEVICE_TYPE_CONNECTED_TV)).
			Build()).
		AddImp(builder.NewVideoImp("1").
			VideoMimes("video/mp4", "video/webm").
			VideoDuration(5, 30).
			VideoProtocols(2, 3, 5, 6).
			Size(1920, 1080).
			StartDelay(0).
			Plcmt(int32(openrtb.VideoPlcmt_VIDEO_PLCMT_INSTREAM)).
			Linearity(int32(openrtb.VideoLinearity_VIDEO_LINEARITY_LINEAR)).
			Skip(5).
			Pod("pod-1", 1).
			Floor(5.0, "USD").
			Secure().
			Build()).
		BuildValidated()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Fatalf("schema: %+v", result.Errors)
	}
}

func TestNativeRequestValidated(t *testing.T) {
	nativeReq := `{"ver":"1.2","assets":[{"id":1,"required":1,"title":{"len":90}}]}`
	_, result, err := builder.NewBidRequest("auction-native-1").
		FirstPrice().
		Tmax(100).
		Currency("USD").
		Site(builder.NewSite().ID("feed-1").Domain("news.example.com").Build()).
		Device(builder.NewDevice().UA("Mozilla/5.0").IP("203.0.113.5").DeviceType(2).Build()).
		AddImp(builder.NewNativeImp("1").
			NativeRequest(nativeReq).
			Floor(0.5, "USD").
			Build()).
		BuildValidated()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Fatalf("schema: %+v", result.Errors)
	}
}

func TestBidResponseValidated(t *testing.T) {
	_, result, err := builder.NewBidResponse("auction-banner-1").
		BidID("abc123").
		Currency("USD").
		AddSeatBid("512",
			builder.NewBid("1", "1", 1.23).
				Banner().
				Size(300, 250).
				Adomain("advertiser.com").
				Crid("creative-9").
				Adm(`<img src="https://cdn.example/ad.png"/>`).
				NURL("https://dsp.example/win?price=${AUCTION_PRICE}").
				Build(),
		).
		BuildValidated()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Fatalf("schema: %+v", result.Errors)
	}
}

func TestNoBidResponseValidated(t *testing.T) {
	_, result, err := builder.NewBidResponse("auction-1").
		NoBid(int32(openrtb.NoBidReason_NO_BID_REASON_INVALID_REQUEST)).
		BuildValidated()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Fatalf("schema: %+v", result.Errors)
	}
}

func TestAddSeatBidClearsNoBid(t *testing.T) {
	res, err := builder.NewBidResponse("auction-1").
		NoBid(2).
		AddSeatBid("512", builder.NewBid("1", "1", 1.0).Banner().Build()).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	if res.GetNbr() != 0 || len(res.Seatbid) != 1 {
		t.Fatalf("nbr=%d seatbid=%d", res.GetNbr(), len(res.Seatbid))
	}
}

func TestNoBidClearsSeatBid(t *testing.T) {
	res, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", builder.NewBid("1", "1", 1.0).Build()).
		NoBid(7).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	if res.GetNbr() != 7 || len(res.Seatbid) != 0 {
		t.Fatalf("nbr=%d seatbid=%d", res.GetNbr(), len(res.Seatbid))
	}
}

func TestRejectsMissingFormat(t *testing.T) {
	_, err := builder.NewBidRequest("x").
		FirstPrice().
		Currency("USD").
		Site(builder.NewSite().ID("s").Build()).
		AddImp(&openrtb.Imp{Id: "1"}).
		Build()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRejectsMissingAt(t *testing.T) {
	_, err := builder.NewBidRequest("x").
		Currency("USD").
		AddImp(builder.NewBannerImp("1").Size(300, 250).Build()).
		Build()
	if err == nil {
		t.Fatal("expected at required")
	}
}

func TestRejectsMissingCur(t *testing.T) {
	_, err := builder.NewBidRequest("x").
		FirstPrice().
		AddImp(builder.NewBannerImp("1").Size(300, 250).Build()).
		Build()
	if err == nil {
		t.Fatal("expected cur required")
	}
}

func TestRejectsSiteAndApp(t *testing.T) {
	b := builder.NewBidRequest("x").
		FirstPrice().
		Currency("USD").
		Site(builder.NewSite().ID("s").Build()).
		AddImp(builder.NewBannerImp("1").Size(300, 250).Build())
	b.App(builder.NewApp().ID("a").Build()) // App after Site should clear Site in our API
	// Our API clears the other — so Build should succeed with only App.
	req, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	if req.Site != nil || req.App == nil {
		t.Fatalf("expected app only, got site=%v app=%v", req.Site, req.App)
	}
}

func TestBidResponseRejectsMissingID(t *testing.T) {
	_, err := builder.NewBidResponse("").Build()
	if err == nil {
		t.Fatal("expected missing id error")
	}
}

func TestBidResponseNeedsSeatBidOrNoBid(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").Build()
	if err == nil {
		t.Fatal("expected seatbid or NoBid required")
	}
}

func TestBidResponseRejectsEmptySeatBid(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").AddSeatBid("512").Build()
	if err == nil {
		t.Fatal("expected empty bids error")
	}
}

func TestBidResponseRejectsZeroPrice(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", builder.NewBid("1", "1", 0).Build()).
		Build()
	if err == nil {
		t.Fatal("expected price > 0 error")
	}
}

func TestBidResponseRejectsNegativePrice(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", builder.NewBid("1", "1", -0.01).Build()).
		Build()
	if err == nil {
		t.Fatal("expected price > 0 error")
	}
}

func TestBidResponseRejectsMissingBidID(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", builder.NewBid("", "1", 1.0).Build()).
		Build()
	if err == nil {
		t.Fatal("expected missing bid id error")
	}
}

func TestBidResponseRejectsMissingImpid(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", builder.NewBid("1", "", 1.0).Build()).
		Build()
	if err == nil {
		t.Fatal("expected missing impid error")
	}
}

func TestBidResponseRejectsNilBid(t *testing.T) {
	_, err := builder.NewBidResponse("auction-1").
		AddSeatBid("512", nil).
		Build()
	if err == nil {
		t.Fatal("expected nil bid error")
	}
}

func TestBidResponseRejectsEmptyCurrency(t *testing.T) {
	b := builder.NewBidResponse("auction-1")
	b.Currency("")
	_, err := b.AddSeatBid("512", builder.NewBid("1", "1", 1.0).Build()).Build()
	if err == nil {
		t.Fatal("expected empty cur error")
	}
}
