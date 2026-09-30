package builder

import (
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"google.golang.org/protobuf/proto"
)

// --- Imp factories -----------------------------------------------------------

// NewBannerImp creates a display Imp (sets Imp.banner).
func NewBannerImp(id string) *ImpBuilder {
	return &ImpBuilder{imp: &openrtb.Imp{Id: id, Banner: &openrtb.Banner{}}}
}

// NewVideoImp creates a video Imp (sets Imp.video).
func NewVideoImp(id string) *ImpBuilder {
	return &ImpBuilder{imp: &openrtb.Imp{Id: id, Video: &openrtb.Video{}}}
}

// NewAudioImp creates an audio Imp (sets Imp.audio).
func NewAudioImp(id string) *ImpBuilder {
	return &ImpBuilder{imp: &openrtb.Imp{Id: id, Audio: &openrtb.Audio{}}}
}

// NewNativeImp creates a native Imp (sets Imp.native).
func NewNativeImp(id string) *ImpBuilder {
	return &ImpBuilder{imp: &openrtb.Imp{Id: id, Native: &openrtb.Native{Ver: "1.2"}}}
}

// ImpBuilder configures an Imp and its format object.
type ImpBuilder struct {
	imp *openrtb.Imp
}

// Floor sets the CPM floor and currency.
func (b *ImpBuilder) Floor(bidfloor float64, cur string) *ImpBuilder {
	b.imp.Bidfloor = proto.Float64(bidfloor)
	b.imp.Bidfloorcur = cur
	return b
}

// Secure requires HTTPS creatives (secure=1).
func (b *ImpBuilder) Secure() *ImpBuilder {
	b.imp.Secure = proto.Int32(1)
	return b
}

// TagID sets the publisher's placement ID.
func (b *ImpBuilder) TagID(tagid string) *ImpBuilder {
	b.imp.Tagid = tagid
	return b
}

// Interstitial marks interstitial inventory (instl=1).
func (b *ImpBuilder) Interstitial() *ImpBuilder {
	b.imp.Instl = proto.Int32(1)
	return b
}

// Rewarded marks rewarded inventory (rwdd=1).
func (b *ImpBuilder) Rewarded() *ImpBuilder {
	b.imp.Rwdd = proto.Int32(1)
	return b
}

// Size sets the Banner or Video player width and height.
func (b *ImpBuilder) Size(w, h int32) *ImpBuilder {
	if b.imp.Banner != nil {
		b.imp.Banner.W = proto.Int32(w)
		b.imp.Banner.H = proto.Int32(h)
	}
	if b.imp.Video != nil {
		b.imp.Video.W = proto.Int32(w)
		b.imp.Video.H = proto.Int32(h)
	}
	return b
}

// BannerPos sets Banner.pos (AdPosition).
func (b *ImpBuilder) BannerPos(pos int32) *ImpBuilder {
	if b.imp.Banner != nil {
		b.imp.Banner.Pos = proto.Int32(pos)
	}
	return b
}

// BannerMimes sets Banner.mimes.
func (b *ImpBuilder) BannerMimes(mimes ...string) *ImpBuilder {
	if b.imp.Banner != nil {
		b.imp.Banner.Mimes = append([]string{}, mimes...)
	}
	return b
}

// VideoMimes sets the required Video.mimes.
func (b *ImpBuilder) VideoMimes(mimes ...string) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Mimes = append([]string{}, mimes...)
	}
	return b
}

// VideoDuration sets the minimum and maximum duration in seconds.
func (b *ImpBuilder) VideoDuration(min, max int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Minduration = proto.Int32(min)
		b.imp.Video.Maxduration = proto.Int32(max)
	}
	return b
}

// VideoProtocols sets Video.protocols (integer VAST/DAAST Protocol enum values).
func (b *ImpBuilder) VideoProtocols(protocols ...int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Protocols = append([]int32{}, protocols...)
	}
	return b
}

// StartDelay sets Video.startdelay (0 pre-roll / -1 mid-roll / -2 post-roll / >0 mid-roll delay).
func (b *ImpBuilder) StartDelay(v int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Startdelay = proto.Int32(v)
	}
	return b
}

// Plcmt sets Video.plcmt (VideoPlcmt).
func (b *ImpBuilder) Plcmt(plcmt int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Plcmt = proto.Int32(plcmt)
	}
	return b
}

// Linearity sets Video.linearity.
func (b *ImpBuilder) Linearity(v int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Linearity = proto.Int32(v)
	}
	return b
}

// Skip sets Video.skip and the optional skipafter delay in seconds.
func (b *ImpBuilder) Skip(skipafter int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Skip = proto.Int32(1)
		b.imp.Video.Skipafter = proto.Int32(skipafter)
	}
	return b
}

// Pod places video/audio in an ad pod (podid + slotinpod).
func (b *ImpBuilder) Pod(podid string, slotinpod int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Podid = podid
		b.imp.Video.Slotinpod = proto.Int32(slotinpod)
	}
	if b.imp.Audio != nil {
		b.imp.Audio.Podid = podid
		b.imp.Audio.Slotinpod = proto.Int32(slotinpod)
	}
	return b
}

// PlaybackMethod sets Video.playbackmethod.
func (b *ImpBuilder) PlaybackMethod(methods ...int32) *ImpBuilder {
	if b.imp.Video != nil {
		b.imp.Video.Playbackmethod = append([]int32{}, methods...)
	}
	return b
}

// AudioMimes sets the required Audio.mimes.
func (b *ImpBuilder) AudioMimes(mimes ...string) *ImpBuilder {
	if b.imp.Audio != nil {
		b.imp.Audio.Mimes = append([]string{}, mimes...)
	}
	return b
}

// AudioDuration sets the minimum and maximum Audio duration.
func (b *ImpBuilder) AudioDuration(min, max int32) *ImpBuilder {
	if b.imp.Audio != nil {
		b.imp.Audio.Minduration = proto.Int32(min)
		b.imp.Audio.Maxduration = proto.Int32(max)
	}
	return b
}

// AudioProtocols sets Audio.protocols.
func (b *ImpBuilder) AudioProtocols(protocols ...int32) *ImpBuilder {
	if b.imp.Audio != nil {
		b.imp.Audio.Protocols = append([]int32{}, protocols...)
	}
	return b
}

// Feed sets Audio.feed (FeedType).
func (b *ImpBuilder) Feed(feed int32) *ImpBuilder {
	if b.imp.Audio != nil {
		b.imp.Audio.Feed = proto.Int32(feed)
	}
	return b
}

// NativeRequest sets the Native.request JSON string (Native 1.2 markup request).
func (b *ImpBuilder) NativeRequest(requestJSON string) *ImpBuilder {
	if b.imp.Native != nil {
		b.imp.Native.Request = requestJSON
	}
	return b
}

// NativeVer sets Native.ver (default "1.2").
func (b *ImpBuilder) NativeVer(ver string) *ImpBuilder {
	if b.imp.Native != nil {
		b.imp.Native.Ver = ver
	}
	return b
}

// Build returns the configured Imp.
func (b *ImpBuilder) Build() *openrtb.Imp {
	return proto.Clone(b.imp).(*openrtb.Imp)
}

// --- Inventory / device helpers ---------------------------------------------

// NewSite creates a Site builder.
func NewSite() *SiteBuilder { return &SiteBuilder{s: &openrtb.Site{}} }

// SiteBuilder configures a Site object.
type SiteBuilder struct{ s *openrtb.Site }

// ID sets Site.id.
func (b *SiteBuilder) ID(id string) *SiteBuilder { b.s.Id = id; return b }

// Name sets Site.name.
func (b *SiteBuilder) Name(name string) *SiteBuilder { b.s.Name = name; return b }

// Domain sets Site.domain.
func (b *SiteBuilder) Domain(domain string) *SiteBuilder { b.s.Domain = domain; return b }

// Page sets Site.page.
func (b *SiteBuilder) Page(page string) *SiteBuilder { b.s.Page = page; return b }

// Cat sets the Site.cat category list.
func (b *SiteBuilder) Cat(cats ...string) *SiteBuilder {
	b.s.Cat = append([]string{}, cats...)
	return b
}

// Publisher sets Site.publisher.
func (b *SiteBuilder) Publisher(p *openrtb.Publisher) *SiteBuilder {
	b.s.Publisher = p
	return b
}

// Build returns the configured Site.
func (b *SiteBuilder) Build() *openrtb.Site { return proto.Clone(b.s).(*openrtb.Site) }

// NewApp creates an App builder.
func NewApp() *AppBuilder { return &AppBuilder{a: &openrtb.App{}} }

// AppBuilder configures an App object.
type AppBuilder struct{ a *openrtb.App }

// ID sets App.id.
func (b *AppBuilder) ID(id string) *AppBuilder { b.a.Id = id; return b }

// Name sets App.name.
func (b *AppBuilder) Name(name string) *AppBuilder { b.a.Name = name; return b }

// Bundle sets App.bundle.
func (b *AppBuilder) Bundle(bundle string) *AppBuilder { b.a.Bundle = bundle; return b }

// Domain sets App.domain.
func (b *AppBuilder) Domain(domain string) *AppBuilder { b.a.Domain = domain; return b }

// StoreURL sets App.storeurl.
func (b *AppBuilder) StoreURL(url string) *AppBuilder { b.a.Storeurl = url; return b }

// Publisher sets App.publisher.
func (b *AppBuilder) Publisher(p *openrtb.Publisher) *AppBuilder {
	b.a.Publisher = p
	return b
}

// Content sets App.content.
func (b *AppBuilder) Content(c *openrtb.Content) *AppBuilder {
	b.a.Content = c
	return b
}

// Build returns the configured App.
func (b *AppBuilder) Build() *openrtb.App { return proto.Clone(b.a).(*openrtb.App) }

// NewDooh creates a Dooh builder.
func NewDooh() *DoohBuilder { return &DoohBuilder{d: &openrtb.Dooh{}} }

// DoohBuilder configures a Dooh object.
type DoohBuilder struct{ d *openrtb.Dooh }

// ID sets Dooh.id.
func (b *DoohBuilder) ID(id string) *DoohBuilder { b.d.Id = id; return b }

// Name sets Dooh.name.
func (b *DoohBuilder) Name(name string) *DoohBuilder { b.d.Name = name; return b }

// VenueType sets Dooh.venuetype.
func (b *DoohBuilder) VenueType(ids ...string) *DoohBuilder {
	b.d.Venuetype = append([]string{}, ids...)
	return b
}

// VenueTypeTax sets Dooh.venuetypetax.
func (b *DoohBuilder) VenueTypeTax(tax int32) *DoohBuilder {
	b.d.Venuetypetax = proto.Int32(tax)
	return b
}

// Publisher sets Dooh.publisher.
func (b *DoohBuilder) Publisher(p *openrtb.Publisher) *DoohBuilder {
	b.d.Publisher = p
	return b
}

// Build returns the configured Dooh.
func (b *DoohBuilder) Build() *openrtb.Dooh { return proto.Clone(b.d).(*openrtb.Dooh) }

// NewPublisher creates a Publisher builder.
func NewPublisher() *PublisherBuilder { return &PublisherBuilder{p: &openrtb.Publisher{}} }

// PublisherBuilder configures a Publisher object.
type PublisherBuilder struct{ p *openrtb.Publisher }

// ID sets Publisher.id.
func (b *PublisherBuilder) ID(id string) *PublisherBuilder { b.p.Id = id; return b }

// Name sets Publisher.name.
func (b *PublisherBuilder) Name(name string) *PublisherBuilder { b.p.Name = name; return b }

// Domain sets Publisher.domain.
func (b *PublisherBuilder) Domain(domain string) *PublisherBuilder { b.p.Domain = domain; return b }

// Build returns the configured Publisher.
func (b *PublisherBuilder) Build() *openrtb.Publisher { return proto.Clone(b.p).(*openrtb.Publisher) }

// NewDevice creates a Device builder.
func NewDevice() *DeviceBuilder { return &DeviceBuilder{d: &openrtb.Device{}} }

// DeviceBuilder configures a Device object.
type DeviceBuilder struct{ d *openrtb.Device }

// UA sets Device.ua.
func (b *DeviceBuilder) UA(ua string) *DeviceBuilder { b.d.Ua = ua; return b }

// IP sets Device.ip.
func (b *DeviceBuilder) IP(ip string) *DeviceBuilder { b.d.Ip = ip; return b }

// IPv6 sets Device.ipv6.
func (b *DeviceBuilder) IPv6(ip string) *DeviceBuilder { b.d.Ipv6 = ip; return b }

// DeviceType sets Device.devicetype.
func (b *DeviceBuilder) DeviceType(t int32) *DeviceBuilder {
	b.d.Devicetype = proto.Int32(t)
	return b
}

// Make sets Device.make.
func (b *DeviceBuilder) Make(make string) *DeviceBuilder { b.d.Make = make; return b }

// Model sets Device.model.
func (b *DeviceBuilder) Model(model string) *DeviceBuilder { b.d.Model = model; return b }

// OS sets Device.os and Device.osv.
func (b *DeviceBuilder) OS(os, osv string) *DeviceBuilder { b.d.Os = os; b.d.Osv = osv; return b }

// IFA sets Device.ifa.
func (b *DeviceBuilder) IFA(ifa string) *DeviceBuilder { b.d.Ifa = ifa; return b }

// ConnectionType sets Device.connectiontype.
func (b *DeviceBuilder) ConnectionType(t int32) *DeviceBuilder {
	b.d.Connectiontype = proto.Int32(t)
	return b
}

// Geo sets Device.geo.
func (b *DeviceBuilder) Geo(g *openrtb.Geo) *DeviceBuilder { b.d.Geo = g; return b }

// Build returns the configured Device.
func (b *DeviceBuilder) Build() *openrtb.Device { return proto.Clone(b.d).(*openrtb.Device) }

// NewGeo creates a Geo builder.
func NewGeo() *GeoBuilder { return &GeoBuilder{g: &openrtb.Geo{}} }

// GeoBuilder configures a Geo object.
type GeoBuilder struct{ g *openrtb.Geo }

// LatLon sets Geo.lat and Geo.lon.
func (b *GeoBuilder) LatLon(lat, lon float64) *GeoBuilder {
	b.g.Lat = proto.Float64(lat)
	b.g.Lon = proto.Float64(lon)
	return b
}

// Type sets Geo.type.
func (b *GeoBuilder) Type(t int32) *GeoBuilder { b.g.Type = proto.Int32(t); return b }

// Country sets Geo.country.
func (b *GeoBuilder) Country(c string) *GeoBuilder { b.g.Country = c; return b }

// Region sets Geo.region.
func (b *GeoBuilder) Region(r string) *GeoBuilder { b.g.Region = r; return b }

// City sets Geo.city.
func (b *GeoBuilder) City(c string) *GeoBuilder { b.g.City = c; return b }

// Build returns the configured Geo.
func (b *GeoBuilder) Build() *openrtb.Geo { return proto.Clone(b.g).(*openrtb.Geo) }

// NewContent creates a Content builder.
func NewContent() *ContentBuilder { return &ContentBuilder{c: &openrtb.Content{}} }

// ContentBuilder configures a Content object.
type ContentBuilder struct{ c *openrtb.Content }

// Title sets Content.title.
func (b *ContentBuilder) Title(t string) *ContentBuilder { b.c.Title = t; return b }

// Series sets Content.series.
func (b *ContentBuilder) Series(s string) *ContentBuilder { b.c.Series = s; return b }

// Season sets Content.season.
func (b *ContentBuilder) Season(s string) *ContentBuilder { b.c.Season = s; return b }

// Episode sets Content.episode.
func (b *ContentBuilder) Episode(n int32) *ContentBuilder {
	b.c.Episode = proto.Int32(n)
	return b
}

// Context sets Content.context.
func (b *ContentBuilder) Context(v int32) *ContentBuilder {
	b.c.Context = proto.Int32(v)
	return b
}

// LiveStream sets Content.livestream.
func (b *ContentBuilder) LiveStream(v int32) *ContentBuilder {
	b.c.Livestream = proto.Int32(v)
	return b
}

// Realtime sets Content.realtime.
func (b *ContentBuilder) Realtime(v int32) *ContentBuilder {
	b.c.Realtime = proto.Int32(v)
	return b
}

// Build returns the configured Content.
func (b *ContentBuilder) Build() *openrtb.Content { return proto.Clone(b.c).(*openrtb.Content) }
