use crate::proto::*;
use crate::{codec, validation};

#[cfg(feature = "jsonschema")]
use crate::jsonschema::request;

#[cfg(feature = "jsonschema")]
use super::ValidatedPayload;

/// Fluent BidRequest builder that produces a BidRequest model.
#[derive(Clone, Debug, Default)]
pub struct BidRequestBuilder {
    req: BidRequest,
    error: Option<String>,
}

impl BidRequestBuilder {
    /// Creates a builder; `id` must not be empty.
    pub fn new(id: impl Into<String>) -> Self {
        let id = id.into();
        let error = if id.trim().is_empty() {
            Some("builder: BidRequest.id is required".into())
        } else {
            None
        };
        Self {
            req: BidRequest {
                id,
                ..Default::default()
            },
            error,
        }
    }

    /// First-price auction (`at = 1`).
    pub fn first_price(mut self) -> Self {
        self.req.at = Some(1);
        self
    }

    /// Second-price-plus auction (`at = 2`).
    pub fn second_price_plus(mut self) -> Self {
        self.req.at = Some(2);
        self
    }

    /// Sets the auction type `at`.
    pub fn auction_type(mut self, at: i32) -> Self {
        self.req.at = Some(at);
        self
    }

    /// Sets the maximum processing time `tmax` in milliseconds.
    pub fn tmax(mut self, ms: i32) -> Self {
        self.req.tmax = Some(ms);
        self
    }

    /// Sets allowed currency codes (at least one).
    pub fn currency(mut self, codes: &[&str]) -> Self {
        self.req.cur = codes.iter().map(|s| (*s).to_string()).collect();
        self
    }

    /// Marks a test request (`test = 1`).
    pub fn test(mut self) -> Self {
        self.req.test = Some(1);
        self
    }

    /// Sets blocked IAB categories `bcat`.
    pub fn bcat(mut self, cats: &[&str]) -> Self {
        self.req.bcat = cats.iter().map(|s| (*s).to_string()).collect();
        self
    }

    /// Sets blocked advertiser domains `badv`.
    pub fn badv(mut self, domains: &[&str]) -> Self {
        self.req.badv = domains.iter().map(|s| (*s).to_string()).collect();
        self
    }

    /// Sets website inventory `site` (mutually exclusive with app/dooh).
    pub fn site(mut self, site: Site) -> Self {
        self.req.site = Some(site);
        self.req.app = None;
        self.req.dooh = None;
        self
    }

    /// Sets app inventory `app` (mutually exclusive with site/dooh).
    pub fn app(mut self, app: App) -> Self {
        self.req.app = Some(app);
        self.req.site = None;
        self.req.dooh = None;
        self
    }

    /// Sets DOOH inventory `dooh` (mutually exclusive with site/app).
    pub fn dooh(mut self, dooh: Dooh) -> Self {
        self.req.dooh = Some(dooh);
        self.req.site = None;
        self.req.app = None;
        self
    }

    /// Sets the `device` object.
    pub fn device(mut self, device: Device) -> Self {
        self.req.device = Some(device);
        self
    }

    /// Sets the `user` object.
    pub fn user(mut self, user: User) -> Self {
        self.req.user = Some(user);
        self
    }

    /// Sets the `regs` object.
    pub fn regs(mut self, regs: Regs) -> Self {
        self.req.regs = Some(regs);
        self
    }

    /// Sets the `source` object.
    pub fn source(mut self, source: Source) -> Self {
        self.req.source = Some(source);
        self
    }

    /// Appends an Imp, which must contain banner/video/audio/native.
    pub fn add_imp(mut self, imp: Imp) -> Self {
        self.req.imp.push(imp);
        self
    }

    /// Builds a strongly typed model and performs builder-level validation.
    pub fn build(self) -> Result<BidRequest, String> {
        if let Some(error) = self.error {
            return Err(error);
        }
        validation::request(&self.req)?;
        for (i, imp) in self.req.imp.iter().enumerate() {
            check_imp(imp, i)?;
        }
        Ok(self.req)
    }

    /// Builds and serializes JSON bytes.
    pub fn build_json(self) -> Result<Vec<u8>, String> {
        let v = self.build()?;
        codec::to_json(&v).map_err(|e| e.to_string())
    }

    /// Builds JSON and validates it against the BidRequest Schema, including embedded native content.
    #[cfg(feature = "jsonschema")]
    pub fn build_validated(self) -> Result<ValidatedPayload, String> {
        let json = self.build_json()?;
        let result = request(&json);
        Ok(ValidatedPayload { json, result })
    }
}

fn check_imp(imp: &Imp, i: usize) -> Result<(), String> {
    let result = crate::validation::imp_ready(imp);
    if let Some(error) = result.errors().next() {
        return Err(format!("builder: imp[{i}]: {}", error.message));
    }
    Ok(())
}

// --- Imp / inventory helpers -------------------------------------------------

fn base_imp(id: &str) -> Imp {
    Imp {
        id: id.to_owned(),
        ..Default::default()
    }
}

macro_rules! imp_common {
    ($s:ident) => {
        /// Sets the floor and currency.
        pub fn floor(mut self, bidfloor: f64, cur: &str) -> Self {
            self.imp.bidfloor = Some(bidfloor);
            self.imp.bidfloorcur = cur.to_owned();
            self
        }
        /// Requires HTTPS creatives (`secure = 1`).
        pub fn secure(mut self) -> Self {
            self.imp.secure = Some(1);
            self
        }
        /// Sets the placement tag ID.
        pub fn tag_id(mut self, tagid: &str) -> Self {
            self.imp.tagid = tagid.to_owned();
            self
        }
        /// Interstitial inventory (`instl = 1`).
        pub fn interstitial(mut self) -> Self {
            self.imp.instl = Some(1);
            self
        }
        /// Rewarded inventory (`rwdd = 1`).
        pub fn rewarded(mut self) -> Self {
            self.imp.rwdd = Some(1);
            self
        }
    };
}

/// Banner Imp builder.
#[derive(Clone, Debug)]
pub struct BannerImpBuilder {
    imp: Imp,
    banner: Banner,
}

impl BannerImpBuilder {
    /// Creates a Banner Imp; `id` identifies the Imp.
    pub fn new(id: impl Into<String>) -> Self {
        Self {
            imp: base_imp(&id.into()),
            banner: Banner::default(),
        }
    }
    imp_common!(self);
    /// Sets Banner width and height (as an alternative or addition to `format[]`).
    pub fn size(mut self, w: i32, h: i32) -> Self {
        self.banner.w = Some(w);
        self.banner.h = Some(h);
        self
    }
    /// Sets ad position `pos`.
    pub fn pos(mut self, pos: i32) -> Self {
        self.banner.pos = Some(pos);
        self
    }
    /// Sets allowed MIME types.
    pub fn mimes(mut self, mimes: &[&str]) -> Self {
        self.banner.mimes = mimes.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Produces an Imp model.
    pub fn build(mut self) -> Imp {
        self.imp.banner = Some(self.banner);
        self.imp
    }
}

/// Video Imp builder.
#[derive(Clone, Debug)]
pub struct VideoImpBuilder {
    imp: Imp,
    video: Video,
}

impl VideoImpBuilder {
    /// Creates a Video Imp.
    pub fn new(id: impl Into<String>) -> Self {
        Self {
            imp: base_imp(&id.into()),
            video: Video::default(),
        }
    }
    imp_common!(self);
    /// Sets video MIME types (required).
    pub fn mimes(mut self, mimes: &[&str]) -> Self {
        self.video.mimes = mimes.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Sets the minimum and maximum duration in seconds.
    pub fn duration(mut self, min: i32, max: i32) -> Self {
        self.video.minduration = Some(min);
        self.video.maxduration = Some(max);
        self
    }
    /// Sets supported video protocols.
    pub fn protocols(mut self, protocols: &[i32]) -> Self {
        self.video.protocols = protocols.to_vec();
        self
    }
    /// Sets player width and height.
    pub fn size(mut self, w: i32, h: i32) -> Self {
        self.video.w = Some(w);
        self.video.h = Some(h);
        self
    }
    /// Sets `startdelay`.
    pub fn start_delay(mut self, v: i32) -> Self {
        self.video.startdelay = Some(v);
        self
    }
    /// Sets video placement `plcmt`.
    pub fn plcmt(mut self, plcmt: i32) -> Self {
        self.video.plcmt = Some(plcmt);
        self
    }
    /// Sets `linearity` (linear/nonlinear).
    pub fn linearity(mut self, v: i32) -> Self {
        self.video.linearity = Some(v);
        self
    }
    /// Enables skipping and sets `skipafter` in seconds.
    pub fn skip(mut self, skipafter: i32) -> Self {
        self.video.skip = Some(1);
        self.video.skipafter = Some(skipafter);
        self
    }
    /// Sets ad pod `podid` and `slotinpod`.
    pub fn pod(mut self, podid: &str, slotinpod: i32) -> Self {
        self.video.podid = podid.to_owned();
        self.video.slotinpod = Some(slotinpod);
        self
    }
    /// Sets the `playbackmethod` list.
    pub fn playback_method(mut self, methods: &[i32]) -> Self {
        self.video.playbackmethod = methods.to_vec();
        self
    }
    /// Produces an Imp model.
    pub fn build(mut self) -> Imp {
        self.imp.video = Some(self.video);
        self.imp
    }
}

/// Audio Imp builder.
#[derive(Clone, Debug)]
pub struct AudioImpBuilder {
    imp: Imp,
    audio: Audio,
}

impl AudioImpBuilder {
    /// Creates an Audio Imp.
    pub fn new(id: impl Into<String>) -> Self {
        Self {
            imp: base_imp(&id.into()),
            audio: Audio::default(),
        }
    }
    imp_common!(self);
    /// Sets audio MIME types (required).
    pub fn mimes(mut self, mimes: &[&str]) -> Self {
        self.audio.mimes = mimes.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Sets the minimum and maximum duration in seconds.
    pub fn duration(mut self, min: i32, max: i32) -> Self {
        self.audio.minduration = Some(min);
        self.audio.maxduration = Some(max);
        self
    }
    /// Sets supported audio protocols.
    pub fn protocols(mut self, protocols: &[i32]) -> Self {
        self.audio.protocols = protocols.to_vec();
        self
    }
    /// Sets the feed type.
    pub fn feed(mut self, feed: i32) -> Self {
        self.audio.feed = Some(feed);
        self
    }
    /// Produces an Imp model.
    pub fn build(mut self) -> Imp {
        self.imp.audio = Some(self.audio);
        self.imp
    }
}

/// Native Imp builder.
#[derive(Clone, Debug)]
pub struct NativeImpBuilder {
    imp: Imp,
    native: Native,
}

impl NativeImpBuilder {
    /// Creates a Native Imp (default `native.ver = "1.2"`).
    pub fn new(id: impl Into<String>) -> Self {
        let native = Native {
            ver: "1.2".into(),
            ..Default::default()
        };
        Self {
            imp: base_imp(&id.into()),
            native,
        }
    }
    imp_common!(self);
    /// Sets the Native Request JSON string (required).
    pub fn request(mut self, native_request_json: &str) -> Self {
        self.native.request = native_request_json.to_owned();
        self
    }
    /// Sets the Native protocol version.
    pub fn ver(mut self, ver: &str) -> Self {
        self.native.ver = ver.to_owned();
        self
    }
    /// Produces an Imp model.
    pub fn build(mut self) -> Imp {
        self.imp.native = Some(self.native);
        self.imp
    }
}

/// Site object builder (website inventory).
#[derive(Clone, Debug, Default)]
pub struct SiteBuilder {
    o: Site,
}
impl SiteBuilder {
    /// Creates an empty Site builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the site ID.
    pub fn id(mut self, v: &str) -> Self {
        self.o.id = v.to_owned();
        self
    }
    /// Sets the site name.
    pub fn name(mut self, v: &str) -> Self {
        self.o.name = v.to_owned();
        self
    }
    /// Sets the site domain.
    pub fn domain(mut self, v: &str) -> Self {
        self.o.domain = v.to_owned();
        self
    }
    /// Sets the current page URL.
    pub fn page(mut self, v: &str) -> Self {
        self.o.page = v.to_owned();
        self
    }
    /// Sets IAB content categories.
    pub fn cat(mut self, cats: &[&str]) -> Self {
        self.o.cat = cats.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Sets the Publisher subobject.
    pub fn publisher(mut self, p: Publisher) -> Self {
        self.o.publisher = Some(p);
        self
    }
    /// Produces a Site model.
    pub fn build(self) -> Site {
        self.o
    }
}

/// App object builder (mobile/app inventory).
#[derive(Clone, Debug, Default)]
pub struct AppBuilder {
    o: App,
}
impl AppBuilder {
    /// Creates an empty App builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the app ID.
    pub fn id(mut self, v: &str) -> Self {
        self.o.id = v.to_owned();
        self
    }
    /// Sets the app name.
    pub fn name(mut self, v: &str) -> Self {
        self.o.name = v.to_owned();
        self
    }
    /// Sets the package name/bundle.
    pub fn bundle(mut self, v: &str) -> Self {
        self.o.bundle = v.to_owned();
        self
    }
    /// Sets the app domain.
    pub fn domain(mut self, v: &str) -> Self {
        self.o.domain = v.to_owned();
        self
    }
    /// Sets the Publisher subobject.
    pub fn publisher(mut self, p: Publisher) -> Self {
        self.o.publisher = Some(p);
        self
    }
    /// Sets the Content subobject (distinct from [`crate::view::Inventory`]).
    pub fn content(mut self, c: Content) -> Self {
        self.o.content = Some(c);
        self
    }
    /// Produces an App model.
    pub fn build(self) -> App {
        self.o
    }
}

/// DOOH object builder (digital out-of-home inventory).
#[derive(Clone, Debug, Default)]
pub struct DoohBuilder {
    o: Dooh,
}
impl DoohBuilder {
    /// Creates an empty DOOH builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the dooh ID.
    pub fn id(mut self, v: &str) -> Self {
        self.o.id = v.to_owned();
        self
    }
    /// Sets the dooh name.
    pub fn name(mut self, v: &str) -> Self {
        self.o.name = v.to_owned();
        self
    }
    /// Sets venue type IDs.
    pub fn venue_type(mut self, ids: &[&str]) -> Self {
        self.o.venuetype = ids.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Sets the venue type taxonomy.
    pub fn venue_type_tax(mut self, tax: i32) -> Self {
        self.o.venuetypetax = Some(tax);
        self
    }
    /// Sets the Publisher subobject.
    pub fn publisher(mut self, p: Publisher) -> Self {
        self.o.publisher = Some(p);
        self
    }
    /// Produces a DOOH model.
    pub fn build(self) -> Dooh {
        self.o
    }
}

/// Device object builder.
#[derive(Clone, Debug, Default)]
pub struct DeviceBuilder {
    o: Device,
}
impl DeviceBuilder {
    /// Creates an empty Device builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the User-Agent.
    pub fn ua(mut self, v: &str) -> Self {
        self.o.ua = v.to_owned();
        self
    }
    /// Sets the IPv4 address.
    pub fn ip(mut self, v: &str) -> Self {
        self.o.ip = v.to_owned();
        self
    }
    /// Sets device type `devicetype`.
    pub fn device_type(mut self, t: i32) -> Self {
        self.o.devicetype = Some(t);
        self
    }
    /// Sets the device manufacturer.
    pub fn make(mut self, v: &str) -> Self {
        self.o.make = v.to_owned();
        self
    }
    /// Sets the device model.
    pub fn model(mut self, v: &str) -> Self {
        self.o.model = v.to_owned();
        self
    }
    /// Sets the operating system and version.
    pub fn os(mut self, os: &str, osv: &str) -> Self {
        self.o.os = os.to_owned();
        self.o.osv = osv.to_owned();
        self
    }
    /// Sets the IFA (advertising identifier).
    pub fn ifa(mut self, v: &str) -> Self {
        self.o.ifa = v.to_owned();
        self
    }
    /// Sets the Geo subobject.
    pub fn geo(mut self, g: Geo) -> Self {
        self.o.geo = Some(g);
        self
    }
    /// Produces a Device model.
    pub fn build(self) -> Device {
        self.o
    }
}

/// Publisher object builder.
#[derive(Clone, Debug, Default)]
pub struct PublisherBuilder {
    o: Publisher,
}
impl PublisherBuilder {
    /// Creates an empty Publisher builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the publisher ID.
    pub fn id(mut self, v: &str) -> Self {
        self.o.id = v.to_owned();
        self
    }
    /// Sets the publisher name.
    pub fn name(mut self, v: &str) -> Self {
        self.o.name = v.to_owned();
        self
    }
    /// Sets the publisher domain.
    pub fn domain(mut self, v: &str) -> Self {
        self.o.domain = v.to_owned();
        self
    }
    /// Produces a Publisher model.
    pub fn build(self) -> Publisher {
        self.o
    }
}

/// Content object builder (content metadata; distinct from [`crate::view::Inventory`]).
#[derive(Clone, Debug, Default)]
pub struct ContentBuilder {
    o: Content,
}
impl ContentBuilder {
    /// Creates an empty Content builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets the content title.
    pub fn title(mut self, v: &str) -> Self {
        self.o.title = v.to_owned();
        self
    }
    /// Sets the series name.
    pub fn series(mut self, v: &str) -> Self {
        self.o.series = v.to_owned();
        self
    }
    /// Sets the season.
    pub fn season(mut self, v: &str) -> Self {
        self.o.season = v.to_owned();
        self
    }
    /// Sets the episode number.
    pub fn episode(mut self, n: i32) -> Self {
        self.o.episode = Some(n);
        self
    }
    /// Sets content context `context`.
    pub fn context(mut self, v: i32) -> Self {
        self.o.context = Some(v);
        self
    }
    /// Sets whether this is a live stream.
    pub fn livestream(mut self, v: i32) -> Self {
        self.o.livestream = Some(v);
        self
    }
    /// Sets whether this is real-time content.
    pub fn realtime(mut self, v: i32) -> Self {
        self.o.realtime = Some(v);
        self
    }
    /// Produces a Content model.
    pub fn build(self) -> Content {
        self.o
    }
}

/// Geo object builder.
#[derive(Clone, Debug, Default)]
pub struct GeoBuilder {
    o: Geo,
}
impl GeoBuilder {
    /// Creates an empty Geo builder.
    pub fn new() -> Self {
        Self::default()
    }
    /// Sets latitude and longitude.
    pub fn lat_lon(mut self, lat: f64, lon: f64) -> Self {
        self.o.lat = Some(lat);
        self.o.lon = Some(lon);
        self
    }
    /// Sets location type `type`.
    pub fn type_(mut self, t: i32) -> Self {
        self.o.r#type = Some(t);
        self
    }
    /// Sets the country code.
    pub fn country(mut self, c: &str) -> Self {
        self.o.country = c.to_owned();
        self
    }
    /// Sets the city.
    pub fn city(mut self, c: &str) -> Self {
        self.o.city = c.to_owned();
        self
    }
    /// Produces a Geo model.
    pub fn build(self) -> Geo {
        self.o
    }
}
