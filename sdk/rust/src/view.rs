//! Lightweight BidRequest / BidResponse query views for matching, bidding, and response decisions.
//!
//! Uses strongly typed models generated from the same proto definitions; does **not** perform full JSON Schema validation.
//! Use [`crate::jsonschema`] for contract validation.
//!
//! Core concepts:
//! - [`MarkupMask`]: bitmask of presentation types on an Imp (corresponding to `Bid.mtype` 1–4),
//!   distinct from `Banner.format[]` (size list).
//! - [`Inventory`]: BidRequest-level inventory type (site/app/dooh, mutually exclusive),
//!   distinct from protobuf `Content.Channel` or the OpenRTB `Content` object.

use std::sync::Arc;
use std::time::{Duration, Instant};

use crate::proto::*;
use crate::validation;

/// Bitmask of presentation types on an Imp (banner / video / audio / native).
///
/// Corresponds to OpenRTB `Bid.mtype`: `BANNER=1`, `VIDEO=2`, `AUDIO=3`, `NATIVE=4`.
///
/// **Note**: identifies presentation object types attached to an Imp, distinct from `Banner.format[]`
/// (the allowed size list), which describes Banner creative dimensions rather than presentation types.
#[derive(Clone, Copy, Debug, Default, Eq, PartialEq)]
pub struct MarkupMask(u8);

impl MarkupMask {
    pub fn from_imp(imp: &Imp) -> MarkupMask {
        let mut mask = MarkupMask::NONE;
        if imp.banner.is_some() {
            mask |= MarkupMask::BANNER;
        }
        if imp.video.is_some() {
            mask |= MarkupMask::VIDEO;
        }
        if imp.audio.is_some() {
            mask |= MarkupMask::AUDIO;
        }
        if imp.native.is_some() {
            mask |= MarkupMask::NATIVE;
        }
        mask
    }
    pub fn from_mtype(mtype: i32) -> MarkupMask {
        match mtype {
            1 => MarkupMask::BANNER,
            2 => MarkupMask::VIDEO,
            3 => MarkupMask::AUDIO,
            4 => MarkupMask::NATIVE,
            _ => MarkupMask::NONE,
        }
    }

    /// Empty mask (no presentation types).
    pub const NONE: MarkupMask = MarkupMask(0);
    /// Banner (`Bid.mtype = 1`).
    pub const BANNER: MarkupMask = MarkupMask(1 << 0);
    /// Video (`Bid.mtype = 2`).
    pub const VIDEO: MarkupMask = MarkupMask(1 << 1);
    /// Audio (`Bid.mtype = 3`).
    pub const AUDIO: MarkupMask = MarkupMask(1 << 2);
    /// Native (`Bid.mtype = 4`).
    pub const NATIVE: MarkupMask = MarkupMask(1 << 3);

    /// Returns the underlying bits.
    pub fn bits(self) -> u8 {
        self.0
    }

    /// Reports whether the type represented by `flag` is present.
    pub fn has(self, flag: MarkupMask) -> bool {
        self.0 & flag.0 != 0
    }

    /// Reports whether the Banner bit is set.
    pub fn has_banner(self) -> bool {
        self.has(Self::BANNER)
    }
    /// Reports whether the Video bit is set.
    pub fn has_video(self) -> bool {
        self.has(Self::VIDEO)
    }
    /// Reports whether the Audio bit is set.
    pub fn has_audio(self) -> bool {
        self.has(Self::AUDIO)
    }
    /// Reports whether the Native bit is set.
    pub fn has_native(self) -> bool {
        self.has(Self::NATIVE)
    }

    /// Number of set type bits.
    pub fn count(self) -> u32 {
        self.0.count_ones()
    }

    /// Returns the mask only when exactly one type is set, otherwise [`Self::NONE`].
    pub fn primary(self) -> MarkupMask {
        if self.count() == 1 {
            self
        } else {
            Self::NONE
        }
    }

    /// Infers OpenRTB `Bid.mtype`; returns 0 for multiple or no formats.
    pub fn mtype(self) -> i32 {
        match self.primary() {
            Self::BANNER => 1,
            Self::VIDEO => 2,
            Self::AUDIO => 3,
            Self::NATIVE => 4,
            _ => 0,
        }
    }
}

impl std::ops::BitOr for MarkupMask {
    type Output = MarkupMask;
    fn bitor(self, rhs: MarkupMask) -> MarkupMask {
        MarkupMask(self.0 | rhs.0)
    }
}

impl std::ops::BitOrAssign for MarkupMask {
    fn bitor_assign(&mut self, rhs: MarkupMask) {
        self.0 |= rhs.0;
    }
}

/// BidRequest inventory type (site / app / dooh, mutually exclusive).
///
/// Inferred from top-level `site`, `app`, and `dooh` fields to identify the **inventory medium**.
///
/// **Note**: distinct from protobuf `Content.Channel`, the OpenRTB `Content` object,
/// and `ContentBuilder`, which describe content metadata rather than inventory type.
#[derive(Clone, Copy, Debug, Default, Eq, PartialEq)]
pub enum Inventory {
    /// No site/app/dooh is specified.
    #[default]
    None,
    /// Website inventory (`BidRequest.site`).
    Site,
    /// Mobile/app inventory (`BidRequest.app`).
    App,
    /// Digital out-of-home inventory (`BidRequest.dooh`).
    Dooh,
}

impl Inventory {
    /// Returns `"site"` / `"app"` / `"dooh"` / `"none"`.
    pub fn as_str(self) -> &'static str {
        match self {
            Inventory::Site => "site",
            Inventory::App => "app",
            Inventory::Dooh => "dooh",
            Inventory::None => "none",
        }
    }
}

/// View of one Imp, including its [`MarkupMask`].
#[derive(Clone, Debug)]
pub struct ImpView<'a> {
    /// Original Imp model.
    pub imp: &'a Imp,
    /// Imp id.
    pub id: &'a str,
    /// Presentation type mask for this Imp.
    pub markup: MarkupMask,
    /// Placement tag ID.
    pub tagid: Option<&'a str>,
    /// Price floor.
    pub bidfloor: f64,
    /// Floor currency.
    pub bidfloorcur: Option<&'a str>,
    /// Whether this is interstitial inventory.
    pub instl: i32,
    /// Whether HTTPS is required.
    pub secure: i32,
    /// Whether this is rewarded inventory.
    pub rwdd: i32,
    /// SSAI flag.
    pub ssai: i32,
    /// `banner` subobject.
    pub banner: Option<&'a Banner>,
    /// `video` subobject.
    pub video: Option<&'a Video>,
    /// `audio` subobject.
    pub audio: Option<&'a Audio>,
    /// `native` subobject.
    pub native: Option<&'a Native>,
    /// `pmp` subobject.
    pub pmp: Option<&'a Pmp>,
}

fn nonempty(s: &str) -> Option<&str> {
    if s.is_empty() {
        None
    } else {
        Some(s)
    }
}

/// Builds an [`ImpView`] list for all Imps.
fn imps(req: &BidRequest) -> Vec<ImpView<'_>> {
    req.imp
        .iter()
        .map(|imp| ImpView {
            imp,
            id: &imp.id,
            markup: MarkupMask::from_imp(imp),
            tagid: nonempty(&imp.tagid),
            bidfloor: imp.bidfloor.unwrap_or(0.0),
            bidfloorcur: nonempty(&imp.bidfloorcur),
            instl: imp.instl.unwrap_or(0),
            secure: imp.secure.unwrap_or(0),
            rwdd: imp.rwdd.unwrap_or(0),
            ssai: imp.ssai.unwrap_or(0),
            banner: imp.banner.as_ref(),
            video: imp.video.as_ref(),
            audio: imp.audio.as_ref(),
            native: imp.native.as_ref(),
            pmp: imp.pmp.as_ref(),
        })
        .collect()
}

/// Query view borrowing an immutable BidRequest.
#[derive(Clone, Debug)]
pub struct RequestView<'a> {
    /// Original request and fixed deadline.
    request: &'a BidRequest,
    deadline: Option<Instant>,
    /// Imp views.
    imps: Vec<ImpView<'a>>,
}

impl<'a> RequestView<'a> {
    /// Borrow an immutable typed request after basic validation.
    pub fn new(req: &'a BidRequest) -> Result<Self, String> {
        validation::request(req)?;
        Ok(Self {
            request: req,
            deadline: (req.tmax.unwrap_or(0) > 0).then(|| {
                Instant::now() + Duration::from_millis(req.tmax.unwrap() as u64 * 85 / 100)
            }),
            imps: imps(req),
        })
    }
    pub fn imps(&self) -> &[ImpView<'a>] {
        &self.imps
    }
    pub fn deadline(&self) -> Option<Instant> {
        self.deadline
    }
    pub fn request(&self) -> &'a BidRequest {
        self.request
    }
    pub fn currencies(&self) -> &'a [String] {
        &self.request.cur
    }

    /// Auction ID.
    pub fn auction_id(&self) -> &str {
        &self.request.id
    }
    /// Auction type `at`.
    pub fn auction_type(&self) -> i32 {
        self.request.at.unwrap_or(0)
    }
    /// Inventory type ([`Inventory`], distinct from Content.Channel).
    pub fn inventory(&self) -> Inventory {
        if self.request.site.is_some() {
            Inventory::Site
        } else if self.request.app.is_some() {
            Inventory::App
        } else if self.request.dooh.is_some() {
            Inventory::Dooh
        } else {
            Inventory::None
        }
    }
    /// Reports whether the deadline at 85% of tmax has passed.
    pub fn past_deadline(&self) -> bool {
        self.deadline.map(|d| Instant::now() > d).unwrap_or(false)
    }

    /// Finds an Imp by ID.
    pub fn find_imp(&self, id: &str) -> Option<&ImpView<'a>> {
        self.imps.iter().find(|i| i.id == id)
    }

    /// Returns Imps containing the specified [`MarkupMask`] type.
    pub fn imps_with(&self, flag: MarkupMask) -> Vec<&ImpView<'a>> {
        self.imps.iter().filter(|i| i.markup.has(flag)).collect()
    }

    /// Flattens each Imp into an [`ImpFact`] for downstream consumption.
    pub fn facts(&self) -> Vec<ImpFact<'a>> {
        self.imps.iter().map(ImpFact::from_view).collect()
    }
}

/// Flattened facts for one Imp.
#[derive(Clone, Debug)]
pub struct ImpFact<'a> {
    /// Imp id.
    pub id: &'a str,
    /// Presentation type mask.
    pub markup: MarkupMask,
    /// Inferred `Bid.mtype` (0 for multiple formats).
    pub mtype: i32,
    /// tag id.
    pub tagid: Option<&'a str>,
    /// Price floor.
    pub bidfloor: f64,
    /// Floor currency.
    pub bidfloorcur: Option<&'a str>,
    /// secure flag.
    pub secure: i32,
    /// Interstitial flag.
    pub instl: i32,
    /// Rewarded flag.
    pub rwdd: i32,
    /// SSAI flag.
    pub ssai: i32,
    /// Banner width, if present.
    pub banner_w: Option<i32>,
    /// Banner height, if present.
    pub banner_h: Option<i32>,
    /// Native request string, if present.
    pub native_request: Option<&'a str>,
}

impl<'a> ImpFact<'a> {
    fn from_view(iv: &ImpView<'a>) -> Self {
        let banner_w = iv.banner.and_then(|b| b.w);
        let banner_h = iv.banner.and_then(|b| b.h);
        let native_request = iv
            .native
            .map(|n| n.request.as_str())
            .filter(|s| !s.is_empty());
        Self {
            id: iv.id,
            markup: iv.markup,
            mtype: iv.markup.mtype(),
            tagid: iv.tagid,
            bidfloor: iv.bidfloor,
            bidfloorcur: iv.bidfloorcur,
            secure: iv.secure,
            instl: iv.instl,
            rwdd: iv.rwdd,
            ssai: iv.ssai,
            banner_w,
            banner_h,
            native_request,
        }
    }

    /// Reports whether the Banner type is present.
    pub fn has_banner(&self) -> bool {
        self.markup.has_banner()
    }
    /// Reports whether the Video type is present.
    pub fn has_video(&self) -> bool {
        self.markup.has_video()
    }
    /// Reports whether the Audio type is present.
    pub fn has_audio(&self) -> bool {
        self.markup.has_audio()
    }
    /// Reports whether the Native type is present.
    pub fn has_native(&self) -> bool {
        self.markup.has_native()
    }
}

// --- BidResponse view ---------------------------------------------------

/// One SeatBid and its nested Bid list.
#[derive(Clone, Debug)]
pub struct SeatBidView<'a> {
    /// Original SeatBid model.
    pub seatbid: &'a SeatBid,
    /// Buyer seat ID.
    pub seat: Option<&'a str>,
    /// group flag.
    pub group: i32,
    /// Bid views for this seat.
    bids: Arc<[BidView<'a>]>,
    range: std::ops::Range<usize>,
}

/// View of one Bid, including the [`MarkupMask`] inferred from mtype.
#[derive(Clone, Debug)]
pub struct BidView<'a> {
    /// Original Bid model.
    pub bid: &'a Bid,
    /// Bid id.
    pub id: &'a str,
    /// Target Imp ID.
    pub impid: &'a str,
    /// Owning seat from the parent SeatBid.
    pub seat: Option<&'a str>,
    /// Bid price.
    pub price: f64,
    /// Presentation type inferred from mtype.
    pub markup: MarkupMask,
    /// Raw OpenRTB mtype value.
    pub mtype: i32,
    /// Creative ID.
    pub crid: Option<&'a str>,
    /// Campaign ID.
    pub cid: Option<&'a str>,
    /// Deal id.
    pub dealid: Option<&'a str>,
    /// Creative width.
    pub w: i32,
    /// Creative height.
    pub h: i32,
    /// Duration in seconds.
    pub dur: i32,
    /// adm markup.
    pub adm: Option<&'a str>,
    /// Win notice URL.
    pub nurl: Option<&'a str>,
    /// Billing notice URL.
    pub burl: Option<&'a str>,
    /// Loss notice URL.
    pub lurl: Option<&'a str>,
    /// Advertiser domain list.
    pub adomain: &'a [String],
}

/// Builds SeatBid views and a flattened Bid list.
fn view_seatbids(res: &BidResponse) -> (Vec<SeatBidView<'_>>, Arc<[BidView<'_>]>) {
    let flat: Arc<[BidView<'_>]> = res
        .seatbid
        .iter()
        .flat_map(|sb| {
            sb.bid
                .iter()
                .map(move |bid| view_bid(bid, nonempty(&sb.seat)))
        })
        .collect();
    let mut offset = 0;
    let seats = res
        .seatbid
        .iter()
        .map(|sb| {
            let range = offset..offset + sb.bid.len();
            offset = range.end;
            SeatBidView {
                seatbid: sb,
                seat: nonempty(&sb.seat),
                group: sb.group.unwrap_or(0),
                bids: Arc::clone(&flat),
                range,
            }
        })
        .collect();
    (seats, flat)
}

fn view_bid<'a>(bid: &'a Bid, seat: Option<&'a str>) -> BidView<'a> {
    BidView {
        bid,
        id: &bid.id,
        impid: &bid.impid,
        seat,
        price: bid.price.unwrap_or(0.0),
        markup: MarkupMask::from_mtype(bid.mtype),
        mtype: bid.mtype,
        adomain: &bid.adomain,
        crid: nonempty(&bid.crid),
        cid: nonempty(&bid.cid),
        dealid: nonempty(&bid.dealid),
        adm: nonempty(&bid.adm),
        nurl: nonempty(&bid.nurl),
        burl: nonempty(&bid.burl),
        lurl: nonempty(&bid.lurl),
        w: bid.w.unwrap_or(0),
        h: bid.h.unwrap_or(0),
        dur: bid.dur.unwrap_or(0),
    }
}

/// Query view borrowing an immutable BidResponse.
#[derive(Clone, Debug)]
pub struct ResponseView<'a> {
    /// Original response.
    response: &'a BidResponse,
    /// SeatBid views.
    seatbids: Vec<SeatBidView<'a>>,
    /// All Bids, flattened.
    bids: Arc<[BidView<'a>]>,
}

impl<'a> ResponseView<'a> {
    /// Borrow an immutable typed response after basic validation.
    pub fn new(res: &'a BidResponse) -> Result<Self, String> {
        validation::response(res)?;
        let (seat_bids, bids) = view_seatbids(res);
        Ok(Self {
            response: res,
            seatbids: seat_bids,
            bids,
        })
    }
    pub fn seatbids(&self) -> &[SeatBidView<'a>] {
        &self.seatbids
    }
    pub fn bids(&self) -> &[BidView<'a>] {
        &self.bids
    }
    pub fn bidid(&self) -> Option<&str> {
        nonempty(&self.response.bidid)
    }
    pub fn response(&self) -> &'a BidResponse {
        self.response
    }

    /// Echoed BidRequest ID.
    pub fn request_id(&self) -> &str {
        &self.response.id
    }
    /// Response currency.
    pub fn currency(&self) -> &str {
        &self.response.cur
    }
    /// Reports whether this is a no-bid.
    pub fn no_bid(&self) -> bool {
        self.response.seatbid.is_empty()
    }
    /// No-bid reason code.
    pub fn nbr(&self) -> i32 {
        self.response.nbr.unwrap_or(0)
    }

    /// Finds a Bid by ID.
    pub fn find_bid(&self, id: &str) -> Option<&BidView<'a>> {
        self.bids.iter().find(|b| b.id == id)
    }

    /// Returns all Bids targeting the specified Imp.
    pub fn bids_for_imp(&self, impid: &str) -> Vec<&BidView<'a>> {
        self.bids.iter().filter(|b| b.impid == impid).collect()
    }

    /// Returns Bids containing the specified [`MarkupMask`] type.
    pub fn bids_with(&self, flag: MarkupMask) -> Vec<&BidView<'a>> {
        self.bids.iter().filter(|b| b.markup.has(flag)).collect()
    }

    /// Flattens each Bid into a [`BidFact`].
    pub fn facts(&self) -> Vec<BidFact<'a>> {
        self.bids.iter().map(BidFact::from_view).collect()
    }
}

/// Flattened facts for one Bid.
#[derive(Clone, Debug)]
pub struct BidFact<'a> {
    /// Bid id.
    pub id: &'a str,
    /// Target Imp ID.
    pub impid: &'a str,
    /// Owning seat.
    pub seat: Option<&'a str>,
    /// Bid price.
    pub price: f64,
    /// Presentation type mask.
    pub markup: MarkupMask,
    /// OpenRTB mtype.
    pub mtype: i32,
    /// Creative ID.
    pub crid: Option<&'a str>,
    /// Deal id.
    pub dealid: Option<&'a str>,
    /// Creative width.
    pub w: i32,
    /// Creative height.
    pub h: i32,
    /// Duration.
    pub dur: i32,
    /// Whether adm is nonempty.
    pub has_adm: bool,
    /// Advertiser domains.
    pub adomain: &'a [String],
}

impl<'a> BidFact<'a> {
    fn from_view(b: &BidView<'a>) -> Self {
        Self {
            id: b.id,
            impid: b.impid,
            seat: b.seat,
            price: b.price,
            markup: b.markup,
            mtype: b.mtype,
            crid: b.crid,
            dealid: b.dealid,
            w: b.w,
            h: b.h,
            dur: b.dur,
            has_adm: b.adm.map(|s| !s.is_empty()).unwrap_or(false),
            adomain: b.adomain,
        }
    }

    /// Reports whether this is the Banner type.
    pub fn has_banner(&self) -> bool {
        self.markup.has_banner()
    }
    /// Reports whether this is the Video type.
    pub fn has_video(&self) -> bool {
        self.markup.has_video()
    }
    /// Reports whether this is the Audio type.
    pub fn has_audio(&self) -> bool {
        self.markup.has_audio()
    }
    /// Reports whether this is the Native type.
    pub fn has_native(&self) -> bool {
        self.markup.has_native()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::proto::{BidRequest, BidResponse, Imp};
    use serde_json::Value;
    fn model<T: serde::de::DeserializeOwned>(v: Value) -> T {
        serde_json::from_value(v).unwrap()
    }
    use crate::builder::{
        BannerImpBuilder, BidBuilder, BidRequestBuilder, BidResponseBuilder, DeviceBuilder,
        SiteBuilder,
    };

    #[test]
    fn seat_views_share_flat_storage() {
        let res: BidResponse = model(serde_json::json!({"id":"a", "cur":"USD", "seatbid":[
            {"seat":"first","bid":[{"id":"b1","impid":"1","price":1}]},
            {"seat":"second","bid":[{"id":"b2","impid":"1","price":2}]}
        ]}));
        let view = ResponseView::new(&res).unwrap();
        let cloned = view.clone();
        assert!(std::ptr::eq(view.response(), &res));
        for i in 0..2 {
            assert_eq!(view.seatbids()[i].bids().len(), 1);
            assert!(std::ptr::eq(&view.seatbids()[i].bids()[0], &view.bids()[i]));
            assert!(std::ptr::eq(&cloned.bids()[i], &view.bids()[i]));
        }
        assert_eq!(view.seatbids()[1].bids()[0].id, "b2");
    }

    #[test]
    fn view_banner_request() {
        let req = BidRequestBuilder::new("auction-1")
            .first_price()
            .tmax(120)
            .currency(&["USD"])
            .site(
                SiteBuilder::new()
                    .id("s1")
                    .domain("example.com")
                    .page("https://example.com/a")
                    .build(),
            )
            .device(
                DeviceBuilder::new()
                    .ua("Mozilla/5.0")
                    .ip("192.0.2.1")
                    .device_type(4)
                    .build(),
            )
            .add_imp(
                BannerImpBuilder::new("1")
                    .size(300, 250)
                    .floor(0.03, "USD")
                    .secure()
                    .build(),
            )
            .build()
            .expect("build");

        let res = RequestView::new(&req).expect("view");
        assert_eq!(res.auction_id(), "auction-1");
        assert_eq!(res.auction_type(), 1);
        assert_eq!(res.inventory(), Inventory::Site);
        assert!(res.deadline().is_some());
        assert!(res.request().site.is_some());
        assert!(res.request().device.is_some());
        assert_eq!(res.imps.len(), 1);
        let iv = &res.imps[0];
        assert!(iv.markup.has_banner());
        assert_eq!(iv.markup.count(), 1);
        assert_eq!(iv.markup.mtype(), 1);
        assert_eq!(iv.secure, 1);
        assert!((iv.bidfloor - 0.03).abs() < 1e-9);
        assert_eq!(iv.banner.unwrap().w.unwrap(), 300);
    }

    #[test]
    fn multi_format_and_mutex() {
        let mut req: BidRequest = model(serde_json::json!({
            "id": "m",
            "at": 2,
            "cur": ["USD"],
            "app": {"id": "a1", "bundle": "com.example.app"},
            "imp": [{
                "id": "1",
                "banner": {"w": 320, "h": 50},
                "video": {"mimes": ["video/mp4"]}
            }]
        }));
        let res = RequestView::new(&req).expect("view");
        assert_eq!(res.inventory(), Inventory::App);
        let f = res.imps[0].markup;
        assert!(f.has_banner() && f.has_video());
        assert_eq!(f.count(), 2);
        assert_eq!(f.mtype(), 0);

        req.site = Some(model(serde_json::json!({"id":"s"})));
        assert!(RequestView::new(&req).is_err());
    }

    #[test]
    fn request_view_facts() {
        let req = BidRequestBuilder::new("auction-1")
            .first_price()
            .tmax(120)
            .currency(&["USD"])
            .site(
                SiteBuilder::new()
                    .id("s1")
                    .domain("example.com")
                    .page("https://example.com/a")
                    .build(),
            )
            .add_imp(
                BannerImpBuilder::new("1")
                    .size(300, 250)
                    .floor(0.03, "USD")
                    .secure()
                    .build(),
            )
            .build()
            .expect("build");

        let snap = RequestView::new(&req).expect("run");
        assert_eq!(snap.auction_id(), "auction-1");
        assert_eq!(snap.inventory(), Inventory::Site);
        assert!(snap.find_imp("1").is_some());
        assert_eq!(snap.imps_with(MarkupMask::BANNER).len(), 1);
        let facts = snap.facts();
        assert_eq!(facts.len(), 1);
        assert!(facts[0].has_banner());
        assert_eq!(facts[0].mtype, 1);
        assert_eq!(facts[0].banner_w, Some(300));
        assert!((facts[0].bidfloor - 0.03).abs() < 1e-9);
    }

    #[test]
    fn response_request_view_facts() {
        let res = BidResponseBuilder::new("auction-1")
            .currency("USD")
            .add_seat_bid(
                "512",
                vec![BidBuilder::new("1", "1", 1.23)
                    .banner()
                    .size(300, 250)
                    .adm("<img/>")
                    .crid("c1")
                    .adomain(&["adv.com"])
                    .build()],
            )
            .build()
            .expect("build");

        let snap = ResponseView::new(&res).expect("response view");
        assert_eq!(snap.request_id(), "auction-1");
        assert!(!snap.no_bid());
        assert_eq!(snap.currency(), "USD");
        assert!(snap.find_bid("1").is_some());
        assert_eq!(snap.bids_for_imp("1").len(), 1);
        assert_eq!(snap.bids_with(MarkupMask::BANNER).len(), 1);
        let facts = snap.facts();
        assert_eq!(facts.len(), 1);
        assert!(facts[0].has_banner());
        assert_eq!(facts[0].mtype, 1);
        assert_eq!(facts[0].w, 300);
        assert!(facts[0].has_adm);
        assert!((facts[0].price - 1.23).abs() < 1e-9);
    }

    #[test]
    fn response_view_no_bid() {
        let res = BidResponseBuilder::new("auction-1")
            .no_bid(2)
            .build()
            .expect("build");
        let snap = ResponseView::new(&res).expect("run");
        assert!(snap.no_bid());
        assert_eq!(snap.nbr(), 2);
        assert!(snap.bids.is_empty());
    }

    #[test]
    fn markup_mask_bits_and_mtype() {
        let m = MarkupMask::BANNER | MarkupMask::VIDEO;
        assert_eq!(m.count(), 2);
        assert!(m.has(MarkupMask::BANNER));
        assert_eq!(m.mtype(), 0);
        assert_eq!(MarkupMask::VIDEO.mtype(), 2);
        assert_eq!(MarkupMask::from_mtype(3), MarkupMask::AUDIO);
        assert_eq!(MarkupMask::from_mtype(99), MarkupMask::NONE);
    }

    #[test]
    fn inventory_dooh_and_as_str() {
        let req: BidRequest = model(serde_json::json!({
            "id": "d", "at": 1, "cur": ["USD"],
            "dooh": {"id": "screen-1"},
            "imp": [{"id": "1", "banner": {"w": 1, "h": 1}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        assert_eq!(snap.inventory(), Inventory::Dooh);
        assert_eq!(Inventory::Dooh.as_str(), "dooh");
        assert_eq!(Inventory::None.as_str(), "none");
    }

    #[test]
    fn markup_mask_not_banner_format_array() {
        let imp: Imp = model(serde_json::json!({
            "id": "1",
            "banner": {"format": [{"w": 300, "h": 250}]}
        }));
        let mask = MarkupMask::from_imp(&imp);
        assert!(mask.has_banner());
        assert_eq!(mask.count(), 1);
        let formats = &imp.banner.as_ref().unwrap().format;
        assert_eq!(formats.len(), 1);
    }

    #[test]
    fn request_factory_validates() {
        let req: BidRequest = model(serde_json::json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {"w": 1, "h": 1}}]
        }));
        assert!(RequestView::new(&req).is_ok());
        let mut invalid = req.clone();
        invalid.id.clear();
        assert!(RequestView::new(&invalid).is_err());
    }

    #[test]
    fn response_factory_validates() {
        let res: BidResponse = model(serde_json::json!({"id": "r", "cur": "USD"}));
        assert!(ResponseView::new(&res).is_ok());
        let mut invalid = res.clone();
        invalid.cur.clear();
        assert!(ResponseView::new(&invalid).is_err());
    }

    #[test]
    fn factory_response_missing_cur() {
        let res: BidResponse = model(serde_json::json!({"id": "r"}));
        assert!(ResponseView::new(&res).is_err());
    }

    #[test]
    fn factory_rejects_blank() {
        let mut req: BidRequest = model(serde_json::json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {"w": 1, "h": 1}}]
        }));
        req.id = "   ".into();
        assert!(RequestView::new(&req).is_err());
        req.id = "x".into();
        req.cur = vec!["   ".into()];
        assert!(RequestView::new(&req).is_err());
        req.cur = vec!["USD".into()];
        req.imp[0].id = "\t".into();
        assert!(RequestView::new(&req).is_err());
    }

    #[test]
    fn factory_response_seatbid_non_array() {
        assert!(
            crate::codec::parse_bid_response(br#"{"id":"r","cur":"USD","seatbid":{}}"#).is_err()
        );
    }

    #[test]
    fn deadline_fixed_at_construction() {
        let req: BidRequest = model(serde_json::json!({
            "id": "x", "at": 1, "tmax": 200, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {"w": 1, "h": 1}}]
        }));
        let p = RequestView::new(&req).unwrap();
        let d1 = p.deadline();
        std::thread::sleep(std::time::Duration::from_millis(3));
        p.facts();
        let d2 = p.deadline();
        assert_eq!(d1, d2);
    }
}

impl<'a> SeatBidView<'a> {
    pub fn bids(&self) -> &[BidView<'a>] {
        &self.bids[self.range.clone()]
    }
}
