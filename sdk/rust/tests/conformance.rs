use oakrtb_sdk::{
    bidcheck, codec,
    proto::{BidRequest, BidResponse},
    validation,
    view::RequestView,
};
use prost::Message;
use serde_json::Value;
#[test]
fn shared_contract() {
    let path = std::path::Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../testdata/conformance/cases.json");
    let groups: Value = serde_json::from_slice(&std::fs::read(path).unwrap()).unwrap();
    for c in groups["roundtrip"].as_array().unwrap() {
        let raw = serde_json::to_vec(&c["input"]).unwrap();
        let encoded = if c["kind"] == "request" {
            let model = codec::parse_bid_request(&raw).unwrap();
            let back = BidRequest::decode(model.encode_to_vec().as_slice()).unwrap();
            codec::to_json(&back).unwrap()
        } else {
            let model = codec::parse_bid_response(&raw).unwrap();
            let back = BidResponse::decode(model.encode_to_vec().as_slice()).unwrap();
            codec::to_json(&back).unwrap()
        };
        assert_eq!(
            serde_json::from_slice::<Value>(&encoded).unwrap(),
            c["expected"],
            "{}",
            c["name"]
        );
    }
    for c in groups["reject"].as_array().unwrap() {
        let raw = serde_json::to_vec(&c["input"]).unwrap();
        let rejected = if c["kind"] == "request" {
            codec::parse_bid_request(&raw).is_err()
        } else {
            codec::parse_bid_response(&raw).is_err()
        };
        assert!(rejected, "{}", c["name"]);
    }
    for c in groups["basic"].as_array().unwrap() {
        let raw = serde_json::to_vec(&c["input"]).unwrap();
        let valid = if c["kind"] == "request" {
            validation::request(&codec::parse_bid_request(&raw).unwrap()).is_ok()
        } else {
            validation::response(&codec::parse_bid_response(&raw).unwrap()).is_ok()
        };
        assert_eq!(valid, c["ok"].as_bool().unwrap(), "{}", c["name"]);
    }
    for c in groups["bidcheck"].as_array().unwrap() {
        let req = codec::parse_bid_request(&serde_json::to_vec(&c["request"]).unwrap()).unwrap();
        let res = codec::parse_bid_response(&serde_json::to_vec(&c["response"]).unwrap()).unwrap();
        let view = RequestView::new(&req).unwrap();
        let result = bidcheck::response(&view, &res);
        assert_eq!(result.ok(), c["ok"].as_bool().unwrap());
        let mut codes: Vec<_> = result.issues.iter().map(|i| i.code.as_str()).collect();
        codes.sort();
        let mut expected: Vec<_> = c["codes"]
            .as_array()
            .unwrap()
            .iter()
            .map(|v| v.as_str().unwrap())
            .collect();
        expected.sort();
        assert_eq!(codes, expected, "{}", c["name"]);
    }

    for c in groups["facts"].as_array().unwrap() {
        let req = codec::parse_bid_request(&serde_json::to_vec(&c["request"]).unwrap()).unwrap();
        let v = RequestView::new(&req).unwrap();
        let f = v.facts();
        assert_eq!(f[0].banner_w, Some(c["banner_w"].as_i64().unwrap() as i32));
        assert_eq!(f[0].banner_h, None);
        assert_eq!(f[0].native_request, c["native_request"].as_str());
        let res = codec::parse_bid_response(&serde_json::to_vec(&c["response"]).unwrap()).unwrap();
        let r = oakrtb_sdk::view::ResponseView::new(&res).unwrap();
        assert_eq!(r.facts()[0].has_adm, c["has_adm"].as_bool().unwrap());
    }
    for c in groups["ready"].as_array().unwrap() {
        let req = codec::parse_bid_request(&serde_json::to_vec(&c["request"]).unwrap()).unwrap();
        let result = validation::imp_ready_mtype(&req.imp[0], c["mtype"].as_i64().unwrap() as i32);
        assert_eq!(result.ok(), c["ok"].as_bool().unwrap(), "{}", c["name"]);
        let mut codes: Vec<_> = result.issues.iter().map(|i| i.code.as_str()).collect();
        let mut expected: Vec<_> = c["codes"]
            .as_array()
            .unwrap()
            .iter()
            .map(|v| v.as_str().unwrap())
            .collect();
        codes.sort();
        expected.sort();
        assert_eq!(codes, expected);
        let builder_ok = c["builder_ok"].as_bool().unwrap();
        assert_eq!(validation::imp_ready(&req.imp[0]).ok(), builder_ok);
        assert_eq!(
            oakrtb_sdk::builder::BidRequestBuilder::new(&req.id)
                .auction_type(req.at.unwrap())
                .currency(&["USD"])
                .add_imp(req.imp[0].clone())
                .build()
                .is_ok(),
            builder_ok
        );
    }
}
