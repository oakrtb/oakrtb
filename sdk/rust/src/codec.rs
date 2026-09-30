//! JSON encoding/decoding of the same generated models used by builders, views and protobuf.
use crate::proto::{BidRequest, BidResponse};
use serde::Serialize;

pub fn to_json<T: Serialize>(model: &T) -> Result<Vec<u8>, serde_json::Error> {
    serde_json::to_vec(model)
}
pub fn parse_bid_request(json: &[u8]) -> Result<BidRequest, serde_json::Error> {
    parse_model(json)
}
pub fn parse_bid_response(json: &[u8]) -> Result<BidResponse, serde_json::Error> {
    parse_model(json)
}
fn parse_model<T: serde::de::DeserializeOwned>(json: &[u8]) -> Result<T, serde_json::Error> {
    let value: serde_json::Value = serde_json::from_slice(json)?;
    if !value.is_object() {
        return Err(<serde_json::Error as serde::de::Error>::custom(
            "expected JSON object",
        ));
    }
    serde_json::from_value(value)
}
