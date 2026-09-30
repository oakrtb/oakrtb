//! Fluent builders that produce strongly typed models generated from the same proto definitions.
//!
//! `build_json()` delegates to codec; `build_validated()` is available with the jsonschema feature.

mod request;
mod response;

pub use request::{
    AppBuilder, AudioImpBuilder, BannerImpBuilder, BidRequestBuilder, ContentBuilder,
    DeviceBuilder, DoohBuilder, GeoBuilder, NativeImpBuilder, PublisherBuilder, SiteBuilder,
    VideoImpBuilder,
};
pub use response::{BidBuilder, BidResponseBuilder};

#[cfg(all(test, feature = "jsonschema"))]
mod build_test;

#[cfg(all(test, feature = "jsonschema"))]
mod full_fields_test;

/// JSON bytes and their independent schema validation report.
#[cfg(feature = "jsonschema")]
#[derive(Debug)]
pub struct ValidatedPayload {
    pub json: Vec<u8>,
    pub result: crate::jsonschema::Report,
}
#[cfg(feature = "jsonschema")]
impl ValidatedPayload {
    pub fn ok(&self) -> bool {
        self.result.ok
    }
}
