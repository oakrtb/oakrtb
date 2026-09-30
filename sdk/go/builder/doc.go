// Package builder provides fluent builders for OakRTB BidRequest / BidResponse.
//
// Builders assemble protobuf models and populate the fields required by each ad format. codec.MarshalJSON emits
// OpenRTB-compatible JSON with numeric enums. Use BuildValidated for full contract validation.
//
// Typical workflow:
//
//	req := builder.NewBidRequest("auction-1").
//		FirstPrice().
//		Tmax(120).
//		Currency("USD").
//		Site(builder.NewSite().ID("s1").Domain("example.com").Page("https://example.com/a").Build()).
//		Device(builder.NewDevice().UA("Mozilla/5.0").IP("192.0.2.1").DeviceType(4).Build()).
//		AddImp(builder.NewBannerImp("1").Size(300, 250).Floor(0.03, "USD").Secure().Build()).
//		MustBuild()
//	raw, err := codec.MarshalJSON(req)
package builder
