// Package codec encodes and decodes the generated OpenRTB models.
package codec

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
)

var marshalOpts = protojson.MarshalOptions{
	UseProtoNames:   true, // OpenRTB JSON uses proto field names (bidfloor, not bidFloor)
	UseEnumNumbers:  true, // OpenRTB enums are integers
	EmitUnpopulated: false,
}

var unmarshalOpts = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

// MarshalJSON encodes a protobuf message as OpenRTB JSON.
func MarshalJSON(m proto.Message) ([]byte, error) {
	if m == nil || !m.ProtoReflect().IsValid() {
		return nil, fmt.Errorf("codec: nil message")
	}
	raw, err := marshalOpts.Marshal(m)
	if err != nil {
		return nil, err
	}
	return wireJSON(raw, m.ProtoReflect().Descriptor(), true)
}

// UnmarshalBidRequest parses OpenRTB BidRequest JSON into a protobuf message.
func UnmarshalBidRequest(data []byte) (*openrtb.BidRequest, error) {
	out := &openrtb.BidRequest{}
	raw, err := wireJSON(data, out.ProtoReflect().Descriptor(), false)
	if err != nil {
		return nil, err
	}
	if err := unmarshalOpts.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnmarshalBidResponse parses OpenRTB BidResponse JSON into a protobuf message.
func UnmarshalBidResponse(data []byte) (*openrtb.BidResponse, error) {
	out := &openrtb.BidResponse{}
	raw, err := wireJSON(data, out.ProtoReflect().Descriptor(), false)
	if err != nil {
		return nil, err
	}
	if err := unmarshalOpts.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}
