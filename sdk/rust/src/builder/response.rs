#[cfg(feature = "jsonschema")]
use super::ValidatedPayload;
#[cfg(feature = "jsonschema")]
use crate::jsonschema::response;
use crate::proto::*;
use crate::{codec, validation};

/// Fluent BidResponse builder that produces a BidResponse model.
#[derive(Clone, Debug)]
pub struct BidResponseBuilder {
    res: BidResponse,
    no_bid_set: bool,
    error: Option<String>,
}

impl BidResponseBuilder {
    /// Creates a builder; `request_id` must echo BidRequest.id.
    pub fn new(request_id: impl Into<String>) -> Self {
        let id = request_id.into();
        let error = if id.trim().is_empty() {
            Some("builder: BidResponse.id is required".into())
        } else {
            None
        };
        Self {
            res: BidResponse {
                id,
                cur: "USD".into(),
                ..Default::default()
            },
            no_bid_set: false,
            error,
        }
    }

    /// Sets the optional DSP-side bid ID.
    pub fn bid_id(mut self, bidid: impl Into<String>) -> Self {
        self.res.bidid = bidid.into();
        self
    }

    /// Sets response currency (ISO-4217, default `"USD"`).
    pub fn currency(mut self, cur: impl Into<String>) -> Self {
        self.res.cur = cur.into();
        self
    }

    /// Structured no-bid: sets `nbr` and clears seatbid.
    pub fn no_bid(mut self, nbr: i32) -> Self {
        self.res.nbr = Some(nbr);
        self.res.seatbid.clear();
        self.no_bid_set = true;
        self
    }

    /// Appends a SeatBid; `bids` must contain at least one Bid.
    pub fn add_seat_bid(mut self, seat: &str, bids: Vec<Bid>) -> Self {
        if bids.is_empty() {
            self.error = Some("builder: SeatBid requires at least one Bid".into());
            return self;
        }
        // Switching to a bid clears structured no-bid.
        self.no_bid_set = false;
        self.res.nbr = None;
        self.res.seatbid.push(SeatBid {
            seat: seat.to_owned(),
            bid: bids,
            ..Default::default()
        });
        self
    }

    /// Builds a strongly typed model and performs builder-level validation.
    pub fn build(self) -> Result<BidResponse, String> {
        if let Some(error) = self.error {
            return Err(error);
        }
        validation::response(&self.res)?;
        if self.res.seatbid.is_empty() && !self.no_bid_set {
            return Err("builder: BidResponse needs seatbid[] or noBid(nbr)".into());
        }
        Ok(self.res)
    }

    /// Builds and serializes JSON bytes.
    pub fn build_json(self) -> Result<Vec<u8>, String> {
        let v = self.build()?;
        codec::to_json(&v).map_err(|e| e.to_string())
    }

    /// Builds JSON and validates it against the BidResponse Schema.
    #[cfg(feature = "jsonschema")]
    pub fn build_validated(self) -> Result<ValidatedPayload, String> {
        let json = self.build_json()?;
        let result = response(&json);
        Ok(ValidatedPayload { json, result })
    }
}

/// Fluent builder for one Bid.
#[derive(Clone, Debug)]
pub struct BidBuilder {
    o: Bid,
}

impl BidBuilder {
    /// Creates a Bid; `price` must be > 0 (validated by BidResponseBuilder).
    pub fn new(id: &str, impid: &str, price: f64) -> Self {
        Self {
            o: Bid {
                id: id.to_owned(),
                impid: impid.to_owned(),
                price: Some(price),
                ..Default::default()
            },
        }
    }

    /// Sets adm (ad markup).
    pub fn adm(mut self, adm: &str) -> Self {
        self.o.adm = adm.to_owned();
        self
    }
    /// Sets nurl (win notice URL).
    pub fn nurl(mut self, u: &str) -> Self {
        self.o.nurl = u.to_owned();
        self
    }
    /// Sets burl (billing notice URL).
    pub fn burl(mut self, u: &str) -> Self {
        self.o.burl = u.to_owned();
        self
    }
    /// Sets crid (creative ID).
    pub fn crid(mut self, crid: &str) -> Self {
        self.o.crid = crid.to_owned();
        self
    }
    /// Sets cid (campaign/ad group ID).
    pub fn cid(mut self, cid: &str) -> Self {
        self.o.cid = cid.to_owned();
        self
    }
    /// Sets advertiser domains.
    pub fn adomain(mut self, domains: &[&str]) -> Self {
        self.o.adomain = domains.iter().map(|s| (*s).to_owned()).collect();
        self
    }
    /// Sets creative width and height.
    pub fn size(mut self, w: i32, h: i32) -> Self {
        self.o.w = Some(w);
        self.o.h = Some(h);
        self
    }
    /// Sets the PMP deal ID.
    pub fn deal_id(mut self, id: &str) -> Self {
        self.o.dealid = id.to_owned();
        self
    }
    /// Sets OpenRTB `mtype` (1=banner, 2=video, 3=audio, 4=native).
    pub fn markup_type(mut self, m: i32) -> Self {
        self.o.mtype = m;
        self
    }
    /// Equivalent to `markup_type(1)`.
    pub fn banner(self) -> Self {
        self.markup_type(1)
    }
    /// Equivalent to `markup_type(2)`.
    pub fn video(self) -> Self {
        self.markup_type(2)
    }
    /// Equivalent to `markup_type(3)`.
    pub fn audio(self) -> Self {
        self.markup_type(3)
    }
    /// Equivalent to `markup_type(4)`.
    pub fn native(self) -> Self {
        self.markup_type(4)
    }
    /// Sets creative duration in seconds, typically for video/audio.
    pub fn dur(mut self, seconds: i32) -> Self {
        self.o.dur = Some(seconds);
        self
    }
    /// Produces a Bid model.
    pub fn build(self) -> Bid {
        self.o
    }
}
