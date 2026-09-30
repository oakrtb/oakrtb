package builder

import (
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"google.golang.org/protobuf/proto"
)

// Full fixtures: every message field set (except ext) with schema-friendly values.
// Generated for exhaustive set/get and JSON round-trip coverage.

func fullFormat() *openrtb.Format {
	return &openrtb.Format{
		W:      proto.Int32(300),
		H:      proto.Int32(250),
		Wratio: proto.Int32(16),
		Hratio: proto.Int32(9),
		Wmin:   proto.Int32(100),
	}
}

func fullMetric() *openrtb.Metric {
	return &openrtb.Metric{
		Type:   "Metric.type.v",
		Value:  proto.Float64(3.14),
		Vendor: "Metric.vendor.v",
	}
}

func fullDurFloors() *openrtb.DurFloors {
	return &openrtb.DurFloors{
		Mindur:      proto.Int32(1),
		Maxdur:      proto.Int32(15),
		Bidfloor:    proto.Float64(0.5),
		Bidfloorcur: "USD",
	}
}

func fullRefSettings() *openrtb.RefSettings {
	return &openrtb.RefSettings{
		Reftype: proto.Int32(3),
		Minint:  proto.Int32(30),
	}
}

func fullRefresh() *openrtb.Refresh {
	return &openrtb.Refresh{
		Refsettings: []*openrtb.RefSettings{fullRefSettings()},
		Count:       proto.Int32(2),
	}
}

func fullQty() *openrtb.Qty {
	return &openrtb.Qty{
		Multiplier: proto.Float64(2.5),
		Sourcetype: proto.Int32(1),
		Vendor:     "Qty.vendor.v",
	}
}

func fullDeal() *openrtb.Deal {
	return &openrtb.Deal{
		Id:           "Deal.id.v",
		Bidfloor:     proto.Float64(0.5),
		Bidfloorcur:  "USD",
		At:           proto.Int32(1),
		Wseat:        []string{"Deal.wseat.a", "Deal.wseat.b"},
		Wadomain:     []string{"Deal.wadomain.a", "Deal.wadomain.b"},
		Guar:         proto.Int32(1),
		Mincpmpersec: proto.Float64(3.14),
		Durfloors:    []*openrtb.DurFloors{fullDurFloors()},
	}
}

func fullPmp() *openrtb.Pmp {
	return &openrtb.Pmp{
		PrivateAuction: proto.Int32(1),
		Deals:          []*openrtb.Deal{fullDeal()},
	}
}

func fullBrandVersion() *openrtb.BrandVersion {
	return &openrtb.BrandVersion{
		Brand:   "BrandVersion.brand.v",
		Version: []string{"BrandVersion.version.a", "BrandVersion.version.b"},
	}
}

func fullUserAgent() *openrtb.UserAgent {
	return &openrtb.UserAgent{
		Browsers:     []*openrtb.BrandVersion{fullBrandVersion()},
		Platform:     fullBrandVersion(),
		Mobile:       proto.Int32(1),
		Architecture: "UserAgent.architecture.v",
		Bitness:      "UserAgent.bitness.v",
		Model:        "UserAgent.model.v",
		Source:       proto.Int32(2),
	}
}

func fullGeo() *openrtb.Geo {
	return &openrtb.Geo{
		Lat:           proto.Float64(37.77),
		Lon:           proto.Float64(-122.42),
		Type:          proto.Int32(1),
		Accuracy:      proto.Int32(50),
		Lastfix:       proto.Int32(10),
		Ipservice:     proto.Int32(3),
		Country:       "USA",
		Region:        "Geo.region.v",
		Metro:         "Geo.metro.v",
		City:          "Geo.city.v",
		Zip:           "Geo.zip.v",
		Utcoffset:     proto.Int32(480),
		Regionfips104: "Geo.regionfips104.v",
	}
}

func fullSegment() *openrtb.Segment {
	return &openrtb.Segment{
		Id:    "Segment.id.v",
		Name:  "Segment.name.v",
		Value: "Segment.value.v",
	}
}

func fullData() *openrtb.Data {
	return &openrtb.Data{
		Id:      "Data.id.v",
		Name:    "Data.name.v",
		Cids:    []string{"Data.cids.a", "Data.cids.b"},
		Segment: []*openrtb.Segment{fullSegment()},
	}
}

func fullUID() *openrtb.UID {
	return &openrtb.UID{
		Id:    "UID.id.v",
		Atype: proto.Int32(1),
	}
}

func fullEID() *openrtb.EID {
	return &openrtb.EID{
		Source:   "EID.source.v",
		Uids:     []*openrtb.UID{fullUID()},
		Inserter: "EID.inserter.v",
		Matcher:  "EID.matcher.v",
		Mm:       proto.Int32(2),
	}
}

func fullPublisher() *openrtb.Publisher {
	return &openrtb.Publisher{
		Id:     "Publisher.id.v",
		Name:   "Publisher.name.v",
		Cattax: proto.Int32(7),
		Cat:    []string{"Publisher.cat.a", "Publisher.cat.b"},
		Domain: "Publisher.domain.v",
	}
}

func fullProducer() *openrtb.Producer {
	return &openrtb.Producer{
		Id:     "Producer.id.v",
		Name:   "Producer.name.v",
		Cattax: proto.Int32(7),
		Cat:    []string{"Producer.cat.a", "Producer.cat.b"},
		Domain: "Producer.domain.v",
	}
}

func fullNetwork() *openrtb.Network {
	return &openrtb.Network{
		Id:     "Network.id.v",
		Name:   "Network.name.v",
		Domain: "Network.domain.v",
	}
}

func fullChannel() *openrtb.Channel {
	return &openrtb.Channel{
		Id:     "Channel.id.v",
		Name:   "Channel.name.v",
		Domain: "Channel.domain.v",
	}
}

func fullContent() *openrtb.Content {
	return &openrtb.Content{
		Id:                 "Content.id.v",
		Episode:            proto.Int32(3),
		Title:              "Content.title.v",
		Series:             "Content.series.v",
		Season:             "Content.season.v",
		Artist:             "Content.artist.v",
		Genre:              "Content.genre.v",
		Album:              "Content.album.v",
		Isrc:               "Content.isrc.v",
		Producer:           fullProducer(),
		Url:                "https://example.com/Content/url",
		Cattax:             proto.Int32(6),
		Cat:                []string{"Content.cat.a", "Content.cat.b"},
		Prodq:              proto.Int32(1),
		Context:            proto.Int32(1),
		Contentrating:      "Content.contentrating.v",
		Userrating:         "Content.userrating.v",
		Qagmediarating:     proto.Int32(1),
		Keywords:           "Content.keywords.v",
		Kwarray:            []string{"Content.kwarray.a", "Content.kwarray.b"},
		Livestream:         proto.Int32(1),
		Sourcerelationship: proto.Int32(1),
		Len:                proto.Int32(1800),
		Language:           "en",
		Langb:              "en-US",
		Embeddable:         proto.Int32(1),
		Data:               []*openrtb.Data{fullData()},
		Network:            fullNetwork(),
		Channel:            fullChannel(),
		Gtax:               proto.Int32(1),
		Genres:             []string{"Content.genres.a", "Content.genres.b"},
		Realtime:           proto.Int32(1),
		Firstbroadcast:     proto.Int32(1),
	}
}

func fullBanner() *openrtb.Banner {
	return &openrtb.Banner{
		Format:   []*openrtb.Format{fullFormat()},
		W:        proto.Int32(300),
		H:        proto.Int32(250),
		Btype:    []int32{3, 4},
		Battr:    []int32{1, 2},
		Pos:      proto.Int32(1),
		Mimes:    []string{"image/jpeg", "image/png"},
		Topframe: proto.Int32(1),
		Expdir:   []int32{1, 2},
		Api:      []int32{3, 5, 7},
		Id:       "Banner.id.v",
		Vcm:      proto.Int32(1),
	}
}

func fullVideo() *openrtb.Video {
	return &openrtb.Video{
		Mimes:          []string{"video/mp4", "video/webm"},
		Minduration:    proto.Int32(5),
		Maxduration:    proto.Int32(30),
		Startdelay:     proto.Int32(-1),
		Maxseq:         proto.Int32(3),
		Poddur:         proto.Int32(90),
		Protocols:      []int32{2, 3, 7},
		W:              proto.Int32(1920),
		H:              proto.Int32(1080),
		Podid:          "Video.podid.v",
		Podseq:         proto.Int32(1),
		Rqddurs:        []int32{15, 30},
		Plcmt:          proto.Int32(1),
		Linearity:      proto.Int32(1),
		Skip:           proto.Int32(1),
		Skipmin:        proto.Int32(5),
		Skipafter:      proto.Int32(5),
		Slotinpod:      proto.Int32(1),
		Mincpmpersec:   proto.Float64(3.14),
		Battr:          []int32{1, 2},
		Maxextended:    proto.Int32(15),
		Minbitrate:     proto.Int32(200),
		Maxbitrate:     proto.Int32(5000),
		Boxingallowed:  proto.Int32(1),
		Playbackmethod: []int32{1, 2},
		Playbackend:    proto.Int32(1),
		Delivery:       []int32{1, 2},
		Pos:            proto.Int32(1),
		Companionad:    []*openrtb.Banner{fullBanner()},
		Api:            []int32{3, 5, 7},
		Companiontype:  []int32{1, 2},
		Placement:      proto.Int32(1),
		Poddedupe:      []int32{1, 3},
		Durfloors:      []*openrtb.DurFloors{fullDurFloors()},
	}
}

func fullAudio() *openrtb.Audio {
	return &openrtb.Audio{
		Mimes:         []string{"audio/mpeg"},
		Minduration:   proto.Int32(5),
		Maxduration:   proto.Int32(30),
		Poddur:        proto.Int32(7),
		Protocols:     []int32{2, 3, 7},
		Startdelay:    proto.Int32(7),
		Rqddurs:       []int32{15, 30},
		Podid:         "Audio.podid.v",
		Podseq:        proto.Int32(1),
		Slotinpod:     proto.Int32(1),
		Mincpmpersec:  proto.Float64(3.14),
		Battr:         []int32{1, 2},
		Maxextended:   proto.Int32(7),
		Minbitrate:    proto.Int32(7),
		Maxbitrate:    proto.Int32(7),
		Delivery:      []int32{1, 2},
		Companionad:   []*openrtb.Banner{fullBanner()},
		Api:           []int32{3, 5, 7},
		Companiontype: []int32{1, 2},
		Maxseq:        proto.Int32(7),
		Feed:          proto.Int32(3),
		Stitched:      proto.Int32(1),
		Nvol:          proto.Int32(1),
		Durfloors:     []*openrtb.DurFloors{fullDurFloors()},
	}
}

func fullNative() *openrtb.Native {
	return &openrtb.Native{
		Request: "{\"ver\":\"1.2\",\"assets\":[{\"id\":1,\"required\":1,\"title\":{\"len\":90}}]}",
		Ver:     "1.2",
		Api:     []int32{3, 5, 7},
		Battr:   []int32{1, 2},
	}
}

func fullImp() *openrtb.Imp {
	return &openrtb.Imp{
		Id:                "Imp.id.v",
		Metric:            []*openrtb.Metric{fullMetric()},
		Banner:            fullBanner(),
		Video:             fullVideo(),
		Audio:             fullAudio(),
		Native:            fullNative(),
		Pmp:               fullPmp(),
		Displaymanager:    "Imp.displaymanager.v",
		Displaymanagerver: "Imp.displaymanagerver.v",
		Instl:             proto.Int32(1),
		Tagid:             "Imp.tagid.v",
		Bidfloor:          proto.Float64(0.5),
		Bidfloorcur:       "USD",
		Clickbrowser:      proto.Int32(1),
		Secure:            proto.Int32(1),
		Iframebuster:      []string{"Imp.iframebuster.a", "Imp.iframebuster.b"},
		Rwdd:              proto.Int32(1),
		Ssai:              proto.Int32(2),
		Exp:               proto.Int32(7),
		Qty:               fullQty(),
		Dt:                proto.Float64(3.14),
		Refresh:           fullRefresh(),
	}
}

func fullSite() *openrtb.Site {
	return &openrtb.Site{
		Id:                     "Site.id.v",
		Name:                   "Site.name.v",
		Domain:                 "Site.domain.v",
		Cattax:                 proto.Int32(7),
		Cat:                    []string{"Site.cat.a", "Site.cat.b"},
		Sectioncat:             []string{"Site.sectioncat.a", "Site.sectioncat.b"},
		Pagecat:                []string{"Site.pagecat.a", "Site.pagecat.b"},
		Page:                   "https://example.com/Site/page",
		Ref:                    "https://example.com/Site/ref",
		Search:                 "Site.search.v",
		Mobile:                 proto.Int32(1),
		Privacypolicy:          proto.Int32(1),
		Publisher:              fullPublisher(),
		Content:                fullContent(),
		Keywords:               "Site.keywords.v",
		Kwarray:                []string{"Site.kwarray.a", "Site.kwarray.b"},
		Inventorypartnerdomain: "Site.inventorypartnerdomain.v",
	}
}

func fullApp() *openrtb.App {
	return &openrtb.App{
		Id:                     "App.id.v",
		Name:                   "App.name.v",
		Bundle:                 "App.bundle.v",
		Domain:                 "App.domain.v",
		Storeurl:               "https://example.com/App/storeurl",
		Cattax:                 proto.Int32(7),
		Cat:                    []string{"App.cat.a", "App.cat.b"},
		Sectioncat:             []string{"App.sectioncat.a", "App.sectioncat.b"},
		Pagecat:                []string{"App.pagecat.a", "App.pagecat.b"},
		Ver:                    "1.0",
		Privacypolicy:          proto.Int32(1),
		Paid:                   proto.Int32(1),
		Publisher:              fullPublisher(),
		Content:                fullContent(),
		Keywords:               "App.keywords.v",
		Kwarray:                []string{"App.kwarray.a", "App.kwarray.b"},
		Inventorypartnerdomain: "App.inventorypartnerdomain.v",
	}
}

func fullDooh() *openrtb.Dooh {
	return &openrtb.Dooh{
		Id:           "Dooh.id.v",
		Name:         "Dooh.name.v",
		Venue:        proto.Int32(1),
		Fixed:        proto.Int32(1),
		Publisher:    fullPublisher(),
		Domain:       "Dooh.domain.v",
		Keywords:     "Dooh.keywords.v",
		Kwarray:      []string{"Dooh.kwarray.a", "Dooh.kwarray.b"},
		Content:      fullContent(),
		Venuetype:    []string{"Dooh.venuetype.a", "Dooh.venuetype.b"},
		Venuetypetax: proto.Int32(1),
	}
}

func fullDevice() *openrtb.Device {
	return &openrtb.Device{
		Geo:            fullGeo(),
		Dnt:            proto.Int32(1),
		Lmt:            proto.Int32(1),
		Ua:             "Device.ua.v",
		Sua:            fullUserAgent(),
		Ip:             "192.0.2.10",
		Ipv6:           "2001:db8::1",
		Devicetype:     proto.Int32(4),
		Make:           "Device.make.v",
		Model:          "Device.model.v",
		Os:             "Device.os.v",
		Osv:            "Device.osv.v",
		Hwv:            "Device.hwv.v",
		H:              proto.Int32(1920),
		W:              proto.Int32(1080),
		Ppi:            proto.Int32(400),
		Pxratio:        proto.Float64(2.0),
		Js:             proto.Int32(1),
		Geofetch:       proto.Int32(1),
		Language:       "en",
		Langb:          "en-US",
		Carrier:        "Device.carrier.v",
		Mccmnc:         "Device.mccmnc.v",
		Connectiontype: proto.Int32(2),
		Ifa:            "Device.ifa.v",
	}
}

func fullUser() *openrtb.User {
	return &openrtb.User{
		Id:         "User.id.v",
		Buyeruid:   "User.buyeruid.v",
		Keywords:   "User.keywords.v",
		Kwarray:    []string{"User.kwarray.a", "User.kwarray.b"},
		Customdata: "User.customdata.v",
		Geo:        fullGeo(),
		Data:       []*openrtb.Data{fullData()},
		Consent:    "User.consent.v",
		Eids:       []*openrtb.EID{fullEID()},
	}
}

func fullSupplyChainNode() *openrtb.SupplyChainNode {
	return &openrtb.SupplyChainNode{
		Asi:    "SupplyChainNode.asi.v",
		Sid:    "SupplyChainNode.sid.v",
		Rid:    "SupplyChainNode.rid.v",
		Name:   "SupplyChainNode.name.v",
		Domain: "SupplyChainNode.domain.v",
		Hp:     proto.Int32(1),
	}
}

func fullSupplyChain() *openrtb.SupplyChain {
	return &openrtb.SupplyChain{
		Complete: proto.Int32(1),
		Nodes:    []*openrtb.SupplyChainNode{fullSupplyChainNode()},
		Ver:      "1.0",
	}
}

func fullSource() *openrtb.Source {
	return &openrtb.Source{
		Fd:     proto.Int32(1),
		Tid:    "Source.tid.v",
		Pchain: "Source.pchain.v",
		Schain: fullSupplyChain(),
	}
}

func fullRegs() *openrtb.Regs {
	return &openrtb.Regs{
		Coppa:     proto.Int32(1),
		Gdpr:      proto.Int32(1),
		UsPrivacy: "Regs.us_privacy.v",
		Gpp:       "Regs.gpp.v",
		GppSid:    []int32{2, 6},
	}
}

func fullBid() *openrtb.Bid {
	return &openrtb.Bid{
		Id:             "Bid.id.v",
		Impid:          "Bid.impid.v",
		Price:          proto.Float64(1.23),
		Nurl:           "https://example.com/Bid/nurl",
		Burl:           "https://example.com/Bid/burl",
		Lurl:           "https://example.com/Bid/lurl",
		Adm:            "Bid.adm.v",
		Adid:           "Bid.adid.v",
		Adomain:        []string{"Bid.adomain.a", "Bid.adomain.b"},
		Bundle:         "Bid.bundle.v",
		Iurl:           "https://example.com/Bid/iurl",
		Cid:            "Bid.cid.v",
		Crid:           "Bid.crid.v",
		Tactic:         "Bid.tactic.v",
		Cattax:         proto.Int32(6),
		Cat:            []string{"Bid.cat.a", "Bid.cat.b"},
		Attr:           []int32{1, 2},
		Apis:           []int32{3, 5, 7},
		Protocol:       proto.Int32(3),
		Qagmediarating: proto.Int32(1),
		Language:       "en",
		Langb:          "en-US",
		Dealid:         "Bid.dealid.v",
		W:              proto.Int32(300),
		H:              proto.Int32(250),
		Wratio:         proto.Int32(16),
		Hratio:         proto.Int32(9),
		Exp:            proto.Int32(60),
		Dur:            proto.Int32(15),
		Mtype:          openrtb.MarkupType_MARKUP_TYPE_BANNER,
		Slotinpod:      proto.Int32(1),
	}
}

func fullSeatBid() *openrtb.SeatBid {
	return &openrtb.SeatBid{
		Bid:   []*openrtb.Bid{fullBid()},
		Seat:  "SeatBid.seat.v",
		Group: proto.Int32(1),
	}
}

func fullBidResponse() *openrtb.BidResponse {
	return &openrtb.BidResponse{
		Id:         "BidResponse.id.v",
		Seatbid:    []*openrtb.SeatBid{fullSeatBid()},
		Bidid:      "BidResponse.bidid.v",
		Cur:        "USD",
		Customdata: "BidResponse.customdata.v",
		Nbr:        proto.Int32(7),
	}
}

func fullBidRequestWithApp() *openrtb.BidRequest {
	return &openrtb.BidRequest{
		Id:      "BidRequest.id.v",
		Imp:     []*openrtb.Imp{fullImp()},
		App:     fullApp(),
		Device:  fullDevice(),
		User:    fullUser(),
		Test:    proto.Int32(1),
		At:      proto.Int32(1),
		Tmax:    proto.Int32(120),
		Wseat:   []string{"BidRequest.wseat.a", "BidRequest.wseat.b"},
		Bseat:   []string{"BidRequest.bseat.a", "BidRequest.bseat.b"},
		Allimps: proto.Int32(1),
		Cur:     []string{"USD", "EUR"},
		Wlang:   []string{"BidRequest.wlang.a", "BidRequest.wlang.b"},
		Wlangb:  []string{"BidRequest.wlangb.a", "BidRequest.wlangb.b"},
		Acat:    []string{"BidRequest.acat.a", "BidRequest.acat.b"},
		Bcat:    []string{"BidRequest.bcat.a", "BidRequest.bcat.b"},
		Cattax:  proto.Int32(6),
		Badv:    []string{"BidRequest.badv.a", "BidRequest.badv.b"},
		Bapp:    []string{"BidRequest.bapp.a", "BidRequest.bapp.b"},
		Source:  fullSource(),
		Regs:    fullRegs(),
	}
}

func fullBidRequestWithSite() *openrtb.BidRequest {
	req := fullBidRequestWithApp()
	req.App = nil
	req.Site = fullSite()
	return req
}

func fullBidRequestWithDooh() *openrtb.BidRequest {
	req := fullBidRequestWithApp()
	req.App = nil
	req.Dooh = fullDooh()
	return req
}

// FieldCounts maps message name → number of proto fields excluding ext.
var FieldCounts = map[string]int{
	"App":             17,
	"Audio":           24,
	"Banner":          12,
	"Bid":             31,
	"BidRequest":      23,
	"BidResponse":     6,
	"BrandVersion":    2,
	"Channel":         3,
	"Content":         33,
	"Data":            4,
	"Deal":            9,
	"Device":          25,
	"Dooh":            11,
	"DurFloors":       4,
	"EID":             5,
	"Format":          5,
	"Geo":             13,
	"Imp":             22,
	"Metric":          3,
	"Native":          4,
	"Network":         3,
	"Pmp":             2,
	"Producer":        5,
	"Publisher":       5,
	"Qty":             3,
	"RefSettings":     2,
	"Refresh":         2,
	"Regs":            5,
	"SeatBid":         3,
	"Segment":         3,
	"Site":            17,
	"Source":          4,
	"SupplyChain":     3,
	"SupplyChainNode": 6,
	"UID":             2,
	"User":            9,
	"UserAgent":       7,
	"Video":           34,
}
