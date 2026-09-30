package builder

import (
	"encoding/json"
	"github.com/oakrtb/oakrtb/sdk/go/codec"
	openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"testing"
)

// Every numeric scalar must distinguish omission from an explicitly supplied zero.
func TestNumericPresenceRoundTrip(t *testing.T) {
	messages := openrtb.File_oakrtb_v2_openrtb_proto.Messages()
	for i := 0; i < messages.Len(); i++ {
		descriptor := messages.Get(i)
		for j := 0; j < descriptor.Fields().Len(); j++ {
			field := descriptor.Fields().Get(j)
			if field.IsList() || (field.Kind() != protoreflect.Int32Kind && field.Kind() != protoreflect.DoubleKind) {
				continue
			}
			t.Run(string(descriptor.Name())+"/"+string(field.Name()), func(t *testing.T) {
				m := dynamicpb.NewMessage(descriptor)
				absent, err := codec.MarshalJSON(m)
				if err != nil {
					t.Fatal(err)
				}
				var obj map[string]json.RawMessage
				if err = json.Unmarshal(absent, &obj); err != nil {
					t.Fatal(err)
				}
				if _, ok := obj[string(field.Name())]; ok {
					t.Fatal("absent field emitted")
				}
				m.Set(field, field.Default())
				raw, err := codec.MarshalJSON(m)
				if err != nil {
					t.Fatal(err)
				}
				back := dynamicpb.NewMessage(descriptor)
				if err = protojson.Unmarshal(raw, back); err != nil {
					t.Fatal(err)
				}
				if !back.Has(field) {
					t.Fatalf("explicit zero lost: %s", raw)
				}
			})
		}
	}
}
func TestNoBidZeroPresence(t *testing.T) {
	raw, err := NewBidResponse("a").NoBid(0).BuildJSON()
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]json.RawMessage
	if err = json.Unmarshal(raw, &obj); err != nil {
		t.Fatal(err)
	}
	if string(obj["nbr"]) != "0" {
		t.Fatalf("zero no-bid reason lost: %s", raw)
	}
	raw, err = NewBidResponse("a").NoBid(0).AddSeatBid("s", NewBid("b", "1", 1).Build()).BuildJSON()
	if err != nil {
		t.Fatal(err)
	}
	obj = nil
	if err = json.Unmarshal(raw, &obj); err != nil {
		t.Fatal(err)
	}
	if _, ok := obj["nbr"]; ok {
		t.Fatalf("no-bid reason retained on bid: %s", raw)
	}
}
