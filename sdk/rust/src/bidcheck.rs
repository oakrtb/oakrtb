//! Optional bid-to-request compatibility checks based on [`RequestView`] and response models.
//!
//! These are neither JSON Schema nor basic validation checks; findings are returned in [`CheckResult`].
//! [`CheckResult::ok`] only indicates the absence of ERROR findings; WARN findings (low prices, blocked ads, deadlines, etc.) may still yield ok.

use crate::proto::{Bid, BidResponse};
use crate::view::{ImpView, MarkupMask, RequestView};
use std::collections::{HashMap, HashSet};

use crate::validation::*;

/// Checks one Bid model against a request snapshot.
/// Without response currency or seat context, skips price comparisons and deal seat allowlists; use [`response`] for full checks.
pub fn bid(req: &RequestView<'_>, bid: &Bid) -> CheckResult {
    bid_at(bid, None, None, "bid", &MatchContext::new(req))
}

/// Checks an entire BidResponse. An empty seatbid still undergoes request ID checks and may produce a PAST_DEADLINE warning.
pub fn response(req: &RequestView<'_>, res: &BidResponse) -> CheckResult {
    let mut issues = Vec::new();
    if req.past_deadline() {
        issues.push(CheckIssue::new(
            CODE_PAST_DEADLINE,
            Severity::Warn,
            "tmax",
            "past request deadline (85% of tmax)",
        ));
    }
    if res.id != req.auction_id() {
        issues.push(CheckIssue::new(
            CODE_REQUEST_ID_MISMATCH,
            Severity::Error,
            "BidResponse.id",
            "response id does not match request id",
        ));
    }
    if res.seatbid.is_empty() {
        return CheckResult { issues };
    }
    let cur = res.cur.trim();
    let res_cur = if cur.is_empty() { None } else { Some(cur) };
    if let Some(cur) = res_cur {
        if !req
            .currencies()
            .iter()
            .any(|c| c.trim().eq_ignore_ascii_case(cur))
        {
            issues.push(CheckIssue::new(
                CODE_CUR_NOT_ALLOWED,
                Severity::Warn,
                "BidResponse.cur",
                "response currency not in BidRequest.cur",
            ));
        }
    }
    let context = MatchContext::new(req);
    for (i, seat) in res.seatbid.iter().enumerate() {
        if seat.bid.is_empty() {
            issues.push(CheckIssue::new(
                CODE_MALFORMED,
                Severity::Error,
                format!("seatbid[{i}].bid"),
                "seatbid.bid must be a non-empty array",
            ));
        }
        for (j, bid) in seat.bid.iter().enumerate() {
            issues.extend(
                bid_at(
                    bid,
                    res_cur,
                    Some(&seat.seat),
                    &format!("seatbid[{i}].bid[{j}]"),
                    &context,
                )
                .issues,
            );
        }
    }
    CheckResult { issues }
}

fn bid_at(
    bid: &Bid,
    response_cur: Option<&str>,
    seat: Option<&str>,
    path: &str,
    context: &MatchContext<'_, '_>,
) -> CheckResult {
    let mut issues = Vec::new();
    let price = bid.price.unwrap_or(0.0);
    if !price.is_finite() || price <= 0.0 {
        return CheckResult {
            issues: vec![CheckIssue::new(
                CODE_MALFORMED,
                Severity::Error,
                format!("{path}.price"),
                "price must be finite and > 0",
            )],
        };
    }
    let impid = bid.impid.as_str();
    let Some(imp) = context.imps.get(impid) else {
        return CheckResult {
            issues: vec![CheckIssue::new(
                CODE_IMP_NOT_FOUND,
                Severity::Error,
                format!("{path}.impid"),
                "impid not found in request",
            )],
        };
    };
    let mtype = bid.mtype;
    let markup = imp.markup;

    if mtype >= 500 {
        issues.push(CheckIssue::new(
            CODE_MTYPE_VENDOR,
            Severity::Warn,
            format!("{path}.mtype"),
            "vendor mtype >=500",
        ));
    } else if mtype != 0 && !(1..=4).contains(&mtype) {
        issues.push(CheckIssue::new(
            CODE_MTYPE_UNKNOWN,
            Severity::Error,
            format!("{path}.mtype"),
            "mtype must be 0–4 or >=500",
        ));
    } else if mtype == 0 && markup.count() > 1 {
        issues.push(CheckIssue::new(
            CODE_MTYPE_REQUIRED,
            Severity::Error,
            format!("{path}.mtype"),
            "multi-format Imp requires Bid.mtype",
        ));
    } else if (1..=4).contains(&mtype) {
        let chosen = MarkupMask::from_mtype(mtype);
        if !markup.has(chosen) {
            issues.push(CheckIssue::new(
                CODE_MTYPE_MISMATCH,
                Severity::Error,
                format!("{path}.mtype"),
                "mtype does not match Imp formats",
            ));
        }
    }

    let eff = if (1..=4).contains(&mtype) {
        MarkupMask::from_mtype(mtype)
    } else if markup.count() == 1 {
        markup
    } else {
        MarkupMask::NONE
    };

    issues.extend(check_floor(imp, bid, response_cur, seat, path));
    issues.extend(check_attr(imp, bid, eff, path));
    issues.extend(check_blocks(context, bid, path));
    CheckResult { issues }
}

fn check_floor(
    imp: &ImpView<'_>,
    bid: &Bid,
    response_cur: Option<&str>,
    seat: Option<&str>,
    path: &str,
) -> Vec<CheckIssue> {
    let mut floor = imp.bidfloor;
    let mut currency = imp.bidfloorcur.unwrap_or("");
    let mut fixed = false;
    let fail = |code, suffix: &str, message| {
        vec![CheckIssue::new(
            code,
            Severity::Error,
            format!("{path}{suffix}"),
            message,
        )]
    };
    if bid.dealid.is_empty() {
        if imp.pmp.is_some_and(|p| p.private_auction == Some(1)) {
            return fail(
                CODE_DEAL_REQUIRED,
                ".dealid",
                "private auction requires a deal",
            );
        }
    } else {
        let Some(deal) = imp
            .pmp
            .and_then(|p| p.deals.iter().find(|d| d.id == bid.dealid))
        else {
            return fail(
                CODE_DEAL_NOT_FOUND,
                ".dealid",
                "dealid not found on impression",
            );
        };
        if seat.is_some_and(|seat| !deal.wseat.is_empty() && !deal.wseat.iter().any(|s| s == seat))
        {
            return fail(
                CODE_DEAL_SEAT_NOT_ALLOWED,
                ".dealid",
                "seat not allowed by deal",
            );
        }
        if !deal.wadomain.is_empty()
            && (bid.adomain.is_empty()
                || bid
                    .adomain
                    .iter()
                    .any(|d| !deal.wadomain.iter().any(|a| a.eq_ignore_ascii_case(d))))
        {
            return fail(
                CODE_DEAL_ADOMAIN_NOT_ALLOWED,
                ".adomain",
                "advertiser domain not allowed by deal",
            );
        }
        if let Some(value) = deal.bidfloor {
            floor = value;
            currency = &deal.bidfloorcur;
        }
        fixed = deal.at == Some(3) && deal.bidfloor.is_some();
    }
    if floor <= 0.0 && !fixed {
        return vec![];
    }
    // Require both response cur and bidfloorcur before numeric compare (`bid` has no response cur).
    let resp_cur = response_cur.map(str::trim).unwrap_or("");
    let floor_cur = currency.trim();
    if resp_cur.is_empty() || floor_cur.is_empty() {
        return vec![];
    }
    let price = bid.price.unwrap_or(0.0);
    if !resp_cur.eq_ignore_ascii_case(floor_cur) {
        return vec![CheckIssue::new(
            CODE_FLOOR_CUR_DIFF,
            Severity::Warn,
            format!("{path}.price"),
            "bid currency differs from applicable bidfloorcur; skip floor compare",
        )];
    }
    if fixed && price != floor {
        return fail(
            CODE_DEAL_PRICE_MISMATCH,
            ".price",
            "price differs from fixed deal price",
        );
    }
    if price < floor {
        return vec![CheckIssue::new(
            CODE_PRICE_BELOW_FLOOR,
            Severity::Warn,
            format!("{path}.price"),
            "price below applicable bidfloor",
        )];
    }
    vec![]
}

fn check_attr(imp: &ImpView<'_>, bid: &Bid, eff: MarkupMask, path: &str) -> Vec<CheckIssue> {
    let blocked = battr_for(imp, eff);
    if bid.attr.iter().any(|a| blocked.contains(a)) {
        vec![CheckIssue::new(
            CODE_ATTR_BLOCKED,
            Severity::Warn,
            format!("{path}.attr"),
            "bid.attr intersects format battr",
        )]
    } else {
        vec![]
    }
}

fn battr_for<'a>(imp: &ImpView<'a>, eff: MarkupMask) -> &'a [i32] {
    if eff.has_banner() {
        imp.banner.map(|b| b.battr.as_slice()).unwrap_or(&[])
    } else if eff.has_video() {
        imp.video.map(|b| b.battr.as_slice()).unwrap_or(&[])
    } else if eff.has_audio() {
        imp.audio.map(|b| b.battr.as_slice()).unwrap_or(&[])
    } else if eff.has_native() {
        imp.native.map(|b| b.battr.as_slice()).unwrap_or(&[])
    } else {
        &[]
    }
}

fn check_blocks(context: &MatchContext<'_, '_>, bid: &Bid, path: &str) -> Vec<CheckIssue> {
    let mut issues = Vec::new();
    if bid
        .adomain
        .iter()
        .any(|d| context.badv.contains(&d.to_ascii_lowercase()))
    {
        issues.push(CheckIssue::new(
            CODE_ADOMAIN_BLOCKED,
            Severity::Warn,
            format!("{path}.adomain"),
            "adomain hit BidRequest.badv",
        ));
    }
    if !bid.bundle.trim().is_empty()
        && context
            .bapp
            .contains(&bid.bundle.trim().to_ascii_lowercase())
    {
        issues.push(CheckIssue::new(
            CODE_BUNDLE_BLOCKED,
            Severity::Warn,
            format!("{path}.bundle"),
            "bundle hit BidRequest.bapp",
        ));
    }
    if bid
        .cat
        .iter()
        .any(|c| context.bcat.contains(&c.to_ascii_lowercase()))
    {
        issues.push(CheckIssue::new(
            CODE_CAT_BLOCKED,
            Severity::Warn,
            format!("{path}.cat"),
            "cat hit BidRequest.bcat",
        ));
    }
    issues
}

struct MatchContext<'a, 'r> {
    imps: HashMap<&'a str, &'a ImpView<'r>>,
    badv: HashSet<String>,
    bapp: HashSet<String>,
    bcat: HashSet<String>,
}
impl<'a, 'r> MatchContext<'a, 'r> {
    fn new(req: &'a RequestView<'r>) -> Self {
        let mut imps = HashMap::with_capacity(req.imps().len());
        for imp in req.imps() {
            imps.entry(imp.id).or_insert(imp);
        }
        let lower = |values: &[String]| {
            values
                .iter()
                .map(|s| s.to_ascii_lowercase())
                .collect::<HashSet<_>>()
        };
        Self {
            imps,
            badv: lower(&req.request().badv),
            bapp: lower(&req.request().bapp),
            bcat: lower(&req.request().bcat),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::proto::{Bid, BidRequest, BidResponse};
    use serde_json::Value;
    fn model<T: serde::de::DeserializeOwned>(v: Value) -> T {
        serde_json::from_value(v).unwrap()
    }
    use crate::view::RequestView;
    use serde_json::json;

    fn banner_req() -> BidRequest {
        model(json!({
            "id": "a1",
            "at": 1,
            "cur": ["USD"],
            "site": {"id": "s1"},
            "imp": [{
                "id": "1",
                "bidfloor": 1.0,
                "bidfloorcur": "USD",
                "banner": {"w": 300, "h": 250}
            }]
        }))
    }

    #[test]
    fn video_mimes_missing() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "video": {}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_mtype(snap.imps()[0].imp, 2);
        assert!(!r.ok());
        assert!(r.has(CODE_VIDEO_MIMES_MISSING));
    }

    #[test]
    fn multi_needs_mtype() {
        let req: BidRequest = model(json!({
            "id": "m", "at": 1, "cur": ["USD"],
            "imp": [{
                "id": "1",
                "banner": {"w": 320, "h": 50},
                "video": {"mimes": ["video/mp4"]}
            }]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 2.0}));
        let r = super::bid(&snap, &bid);
        assert!(!r.ok());
        assert!(r.has(CODE_MTYPE_REQUIRED));
    }

    #[test]
    fn mtype_mismatch() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 2.0, "mtype": 2}));
        let r = super::bid(&snap, &bid);
        assert!(!r.ok());
        assert!(r.has(CODE_MTYPE_MISMATCH));
    }

    #[test]
    fn price_below_floor_warn() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 0.5, "mtype": 1}));
        // bidcheck::bid has no response cur, so skip the floor check; use the response path.
        assert!(!super::bid(&snap, &bid).has(CODE_PRICE_BELOW_FLOOR));
        let res: BidResponse = model(json!({
            "id": "a1", "cur": "USD",
            "seatbid": [{"bid": [bid]}]
        }));
        let r = response(&snap, &res);
        assert!(r.ok());
        assert!(r.has(CODE_PRICE_BELOW_FLOOR));
    }

    #[test]
    fn response_non_object_malformed() {
        assert!(crate::codec::parse_bid_response(b"[]").is_err());
    }

    #[test]
    fn response_seatbid_non_array_malformed() {
        assert!(
            crate::codec::parse_bid_response(br#"{"id":"a1","cur":"USD","seatbid":{}}"#).is_err()
        );
    }

    #[test]
    fn response_seatbid_empty_bid_malformed() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        for seatbid in [serde_json::json!([{}]), serde_json::json!([{"bid":[]}])] {
            let res: BidResponse =
                model(serde_json::json!({"id":"a1","cur":"USD","seatbid":seatbid}));
            assert!(response(&snap, &res).has(CODE_MALFORMED));
        }
        for seatbid in [
            serde_json::json!([{"bid":1}]),
            serde_json::json!([{"bid":[1]}]),
        ] {
            assert!(serde_json::from_value::<BidResponse>(
                serde_json::json!({"id":"a1","cur":"USD","seatbid":seatbid})
            )
            .is_err());
        }
    }

    #[test]
    fn bid_non_object_malformed() {
        assert!(crate::codec::parse_bid_response(br#"{"seatbid":[{"bid":[[]]}]}"#).is_err());
    }

    #[test]
    fn response_blank_cur_skips_cur_and_floor() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 0.5, "mtype": 1}));
        let res: BidResponse = model(json!({
            "id": "a1", "cur": "   ",
            "seatbid": [{"bid": [bid]}]
        }));
        let r = response(&snap, &res);
        assert!(r.ok());
        assert!(!r.has(CODE_CUR_NOT_ALLOWED));
        assert!(!r.has(CODE_PRICE_BELOW_FLOOR));
    }

    #[test]
    fn response_non_string_cur_malformed() {
        assert!(crate::codec::parse_bid_response(br#"{"id":"a1","cur":1}"#).is_err());
    }

    #[test]
    fn fit_result_warnings_errors() {
        let mut r = CheckResult {
            issues: vec![
                CheckIssue::new(CODE_IMP_NOT_FOUND, Severity::Error, "x", "e"),
                CheckIssue::new(CODE_PRICE_BELOW_FLOOR, Severity::Warn, "y", "w"),
            ],
        };
        assert_eq!(r.errors().count(), 1);
        assert_eq!(r.warnings().count(), 1);
        assert!(!r.ok());
        let _ = &mut r;
    }

    #[test]
    fn response_cur_whitespace_trimmed_for_allow_and_floor() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 0.5, "mtype": 1}));
        let res: BidResponse = model(json!({
            "id": "a1", "cur": "  USD  ",
            "seatbid": [{"bid": [bid]}]
        }));
        let r = response(&snap, &res);
        assert!(r.ok());
        assert!(!r.has(CODE_CUR_NOT_ALLOWED));
        assert!(r.has(CODE_PRICE_BELOW_FLOOR));
    }

    #[test]
    fn attr_blocked_warn() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {"w": 300, "h": 250, "battr": [1]}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid =
            model(json!({"id": "b1", "impid": "1", "price": 2.0, "mtype": 1, "attr": [1]}));
        let r = super::bid(&snap, &bid);
        assert!(r.ok());
        assert!(r.has(CODE_ATTR_BLOCKED));
    }

    #[test]
    fn no_bid_response() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let res: BidResponse = model(json!({"id": "a1", "cur": "USD", "nbr": 0}));
        let r = response(&snap, &res);
        assert!(r.ok());
    }

    #[test]
    fn impid_not_found() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "missing", "price": 2.0, "mtype": 1}));
        let r = super::bid(&snap, &bid);
        assert!(!r.ok());
        assert!(r.has(CODE_IMP_NOT_FOUND));
    }

    #[test]
    fn response_happy() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let res: BidResponse = model(json!({
            "id": "a1",
            "cur": "USD",
            "seatbid": [{"bid": [{
                "id": "b1", "impid": "1", "price": 2.0, "mtype": 1, "adm": "<a/>"
            }]}]
        }));
        let r = response(&snap, &res);
        assert!(r.ok());
    }

    #[test]
    fn imp_ready_mtype_unknown() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_mtype(snap.imps()[0].imp, 9);
        assert!(!r.ok());
        assert!(r.has(CODE_MTYPE_UNKNOWN));
    }

    #[test]
    fn imp_ready_markup_multi_bit() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_markup(snap.imps()[0].imp, 1 | 2);
        assert!(!r.ok());
        assert!(r.has(CODE_MULTI_NEEDS_CHOICE));
    }

    #[test]
    fn imp_ready_banner_size_missing() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_markup(snap.imps()[0].imp, 1);
        assert!(!r.ok());
        assert!(r.has(CODE_BANNER_SIZE_MISSING));
    }

    #[test]
    fn imp_ready_banner_format_array_ok() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "banner": {"format": [{"w": 300, "h": 250}]}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_markup(snap.imps()[0].imp, 1);
        assert!(r.ok());
    }

    #[test]
    fn audio_mimes_missing() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "audio": {}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_mtype(snap.imps()[0].imp, 3);
        assert!(!r.ok());
        assert!(r.has(CODE_AUDIO_MIMES_MISSING));
    }

    #[test]
    fn native_request_missing() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "native": {"ver": "1.2"}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let r = imp_ready_mtype(snap.imps()[0].imp, 4);
        assert!(!r.ok());
        assert!(r.has(CODE_NATIVE_REQUEST_MISSING));
    }

    #[test]
    fn floor_cur_diff_skips_compare() {
        let req: BidRequest = model(json!({
            "id": "a1", "at": 1, "cur": ["USD"],
            "imp": [{"id": "1", "bidfloor": 1.0, "bidfloorcur": "EUR", "banner": {"w": 1, "h": 1}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({"id": "b1", "impid": "1", "price": 0.1, "mtype": 1}));
        let res: BidResponse = model(json!({
            "id": "a1", "cur": "USD",
            "seatbid": [{"bid": [bid.clone()]}]
        }));
        let r = response(&snap, &res);
        assert!(r.ok());
        assert!(r.has(CODE_FLOOR_CUR_DIFF));
        assert!(!r.has(CODE_PRICE_BELOW_FLOOR));
    }

    #[test]
    fn adomain_blocked() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"], "badv": ["blocked.com"],
            "imp": [{"id": "1", "banner": {"w": 300, "h": 250}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({
            "id": "b1", "impid": "1", "price": 2.0, "mtype": 1,
            "adomain": ["blocked.com"]
        }));
        let r = super::bid(&snap, &bid);
        assert!(r.has(CODE_ADOMAIN_BLOCKED));
    }

    #[test]
    fn bundle_blocked() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"], "bapp": ["com.blocked"],
            "imp": [{"id": "1", "banner": {"w": 300, "h": 250}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({
            "id": "b1", "impid": "1", "price": 2.0, "mtype": 1,
            "bundle": "com.blocked"
        }));
        let r = super::bid(&snap, &bid);
        assert!(r.has(CODE_BUNDLE_BLOCKED));
    }

    #[test]
    fn cat_blocked() {
        let req: BidRequest = model(json!({
            "id": "x", "at": 1, "cur": ["USD"], "bcat": ["IAB25"],
            "imp": [{"id": "1", "banner": {"w": 300, "h": 250}}]
        }));
        let snap = RequestView::new(&req).unwrap();
        let bid: Bid = model(json!({
            "id": "b1", "impid": "1", "price": 2.0, "mtype": 1,
            "cat": ["IAB25"]
        }));
        let r = super::bid(&snap, &bid);
        assert!(r.has(CODE_CAT_BLOCKED));
    }

    #[test]
    fn cur_not_allowed() {
        let req = banner_req();
        let snap = RequestView::new(&req).unwrap();
        let res: BidResponse = model(json!({
            "id": "a1", "cur": "EUR",
            "seatbid": [{"bid": [{
                "id": "b1", "impid": "1", "price": 2.0, "mtype": 1
            }]}]
        }));
        let r = response(&snap, &res);
        assert!(r.has(CODE_CUR_NOT_ALLOWED));
    }

    #[test]
    fn fit_result_has_and_severity() {
        let r = CheckResult {
            issues: vec![CheckIssue::new(
                CODE_IMP_NOT_FOUND,
                Severity::Error,
                "p",
                "m",
            )],
        };
        assert!(!r.ok());
        assert!(r.has(CODE_IMP_NOT_FOUND));
        assert_eq!(Severity::Error.as_str(), "ERROR");
        assert_eq!(Severity::Warn.as_str(), "WARN");
    }

    #[test]
    fn request_id_must_match_including_no_bid() {
        let raw = banner_req();
        let req = RequestView::new(&raw).unwrap();
        for id in ["a1", "other", ""] {
            for no_bid in [true, false] {
                let mut res: BidResponse = model(json!({"id": id, "cur": "USD"}));
                if !no_bid {
                    res.seatbid = vec![model(json!({"bid": [{"id":"b", "impid":"1", "price":2}]}))];
                }
                let result = response(&req, &res);
                assert_eq!(result.has(CODE_REQUEST_ID_MISMATCH), id != "a1");
                if id != "a1" {
                    assert!(!result.ok());
                }
            }
        }
    }
    #[test]
    fn invalid_prices_are_rejected() {
        let raw = banner_req();
        let req = RequestView::new(&raw).unwrap();
        for price in [f64::NAN, f64::INFINITY, f64::NEG_INFINITY, 0.0, -1.0] {
            let bid = Bid {
                id: "b".into(),
                impid: "1".into(),
                price: Some(price),
                ..Default::default()
            };
            assert!(!super::bid(&req, &bid).ok());
            let res = BidResponse {
                id: "a1".into(),
                cur: "USD".into(),
                seatbid: vec![crate::proto::SeatBid {
                    bid: vec![bid],
                    ..Default::default()
                }],
                ..Default::default()
            };
            assert!(crate::validation::response(&res).is_err());
            assert!(!response(&req, &res).ok());
        }
    }
}
