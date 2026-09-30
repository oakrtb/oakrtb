use oakrtb_sdk::{
    bidcheck,
    builder::{BannerImpBuilder, BidBuilder, BidRequestBuilder, BidResponseBuilder, SiteBuilder},
    codec, jsonschema,
    view::{RequestView, ResponseView},
};

// Simulate the DSP boundary; price == 0 means no candidate was found.
fn respond(request_json: &[u8], price: f64) -> Result<Vec<u8>, String> {
    // Full contract validation is useful at integration boundaries.
    let report = jsonschema::request(request_json);
    if !report.ok {
        return Err(format!("invalid request: {:?}", report.errors));
    }
    let req = codec::parse_bid_request(request_json).map_err(|e| e.to_string())?;
    let request_view = RequestView::new(&req)?;
    let imp = request_view.find_imp("imp-1").ok_or("missing imp-1")?;
    if !imp.markup.has_banner() {
        return Err("expected Banner imp-1".into());
    }

    let mut response = BidResponseBuilder::new(request_view.auction_id()).currency("USD");
    if price == 0.0 {
        // Reason 0 means unknown; this demo has no more specific standard reason.
        response = response.no_bid(0);
    } else {
        let bid = BidBuilder::new("bid-1", imp.id, price).banner().size(300, 250)
            .crid("creative-1").adomain(&["advertiser.example"])
            .adm(r#"<a href="https://advertiser.example/"><img src="https://cdn.example/banner.png" width="300" height="250"></a>"#)
            .build();
        response = response.add_seat_bid("buyer-1", vec![bid]);
    }
    let mut res = response.build()?;
    let findings = bidcheck::response(&request_view, &res);
    if !findings.ok() {
        return Err(format!("response mismatch: {:?}", findings.issues));
    }
    // Demo policy: reject any warning, even when ok() is true.
    let warnings: Vec<_> = findings.warnings().collect();
    if !warnings.is_empty() {
        eprintln!("policy rejected bid: {:?}", warnings);
        res = BidResponseBuilder::new(request_view.auction_id())
            .currency("USD")
            .no_bid(0)
            .build()?;
    }
    let payload = codec::to_json(&res).map_err(|e| e.to_string())?;
    let output_report = jsonschema::response(&payload);
    if !output_report.ok {
        return Err(format!("invalid response: {:?}", output_report.errors));
    }
    Ok(payload)
}

fn main() -> Result<(), String> {
    // SSP: construct an auction request and serialize it for the DSP.
    let req = BidRequestBuilder::new("auction-1")
        .first_price()
        .currency(&["USD"])
        .test()
        .site(
            SiteBuilder::new()
                .id("site-1")
                .domain("publisher.example")
                .build(),
        )
        .add_imp(
            BannerImpBuilder::new("imp-1")
                .size(300, 250)
                .floor(1.0, "USD")
                .secure()
                .build(),
        )
        .build()?;
    let request_json = codec::to_json(&req).map_err(|e| e.to_string())?;
    println!("request: {}", String::from_utf8_lossy(&request_json));
    for price in [2.0, 0.5, 0.0] {
        let response_json = respond(&request_json, price)?;
        // SSP: decode the returned response and inspect its query view.
        let res = codec::parse_bid_response(&response_json).map_err(|e| e.to_string())?;
        let response_view = ResponseView::new(&res)?;
        if response_view.no_bid() != (price < 1.0) {
            return Err("unexpected decision".into());
        }
        println!(
            "price={price} no_bid={} response={}",
            response_view.no_bid(),
            String::from_utf8_lossy(&response_json)
        );
    }
    Ok(())
}
