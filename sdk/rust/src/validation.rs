//! Basic structural checks on generated models; no schema compilation or view construction.
use crate::proto::{BidRequest, BidResponse};
pub fn request(req: &BidRequest) -> Result<(), String> {
    if req.id.trim().is_empty() {
        return Err("validation: BidRequest.id is required".into());
    }
    let at = req.at.unwrap_or(0);
    if at != 1 && at != 2 && at < 500 {
        return Err("validation: BidRequest.at must be 1, 2, or >=500".into());
    }
    if req.cur.is_empty() || req.cur.iter().any(|c| c.trim().is_empty()) {
        return Err("validation: BidRequest.cur is required and must not contain blanks".into());
    }
    if req.imp.is_empty() {
        return Err("validation: BidRequest.imp requires at least one Imp".into());
    }
    if req.site.is_some() as u8 + req.app.is_some() as u8 + req.dooh.is_some() as u8 > 1 {
        return Err("validation: site/app/dooh are mutually exclusive".into());
    }
    let mut ids = std::collections::HashSet::new();
    for (i, imp) in req.imp.iter().enumerate() {
        if imp.id.trim().is_empty() {
            return Err(format!("validation: imp[{i}].id is required"));
        }
        if !ids.insert(&imp.id) {
            return Err(format!("validation: duplicate imp.id {}", imp.id));
        }
        if let Some(pmp) = &imp.pmp {
            let mut deals = std::collections::HashSet::new();
            for d in &pmp.deals {
                if d.id.trim().is_empty() || !deals.insert(&d.id) {
                    return Err("validation: deal.id must be nonblank and unique".into());
                }
                if d.bidfloor.is_some_and(|f| !f.is_finite() || f < 0.0) {
                    return Err("validation: deal.bidfloor must be finite and nonnegative".into());
                }
                if d.at == Some(3) && d.bidfloor.is_none() {
                    return Err("validation: fixed-price deal requires bidfloor".into());
                }
            }
        }
        if imp.banner.is_none()
            && imp.video.is_none()
            && imp.audio.is_none()
            && imp.native.is_none()
        {
            return Err(format!(
                "validation: imp[{i}] needs banner, video, audio, or native"
            ));
        }
    }
    Ok(())
}
pub fn response(res: &BidResponse) -> Result<(), String> {
    if res.id.trim().is_empty() {
        return Err("validation: BidResponse.id is required".into());
    }
    if res.cur.trim().is_empty() {
        return Err("validation: BidResponse.cur is required".into());
    }
    for (i, seat) in res.seatbid.iter().enumerate() {
        if seat.bid.is_empty() {
            return Err(format!("validation: seatbid[{i}] needs at least one bid"));
        }
        for (j, bid) in seat.bid.iter().enumerate() {
            if bid.id.trim().is_empty() || bid.impid.trim().is_empty() {
                return Err(format!(
                    "validation: seatbid[{i}].bid[{j}] requires id and impid"
                ));
            }
            let price = bid.price.unwrap_or(0.0);
            if !price.is_finite() || price <= 0.0 {
                return Err(format!(
                    "validation: seatbid[{i}].bid[{j}].price must be finite and > 0"
                ));
            }
        }
    }
    Ok(())
}

mod readiness;
pub use readiness::*;
mod result;
pub use result::*;
