package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Adapt only declared fields. An ext object is opaque: its own keys are never transformed.
func wireJSON(raw []byte, descriptor protoreflect.MessageDescriptor, encode bool) ([]byte, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, err
	}
	if object == nil {
		return nil, fmt.Errorf("codec: expected JSON object")
	}
	for key, value := range object {
		field := descriptor.Fields().ByName(protoreflect.Name(key))
		if field == nil {
			field = descriptor.Fields().ByJSONName(key)
		}
		if field == nil || bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			continue
		}
		if field.Name() == "ext" && field.Kind() == protoreflect.StringKind {
			if encode {
				var inner string
				if err := json.Unmarshal(value, &inner); err != nil {
					return nil, err
				}
				value = []byte(inner)
			}
			var ext map[string]json.RawMessage
			if err := json.Unmarshal(value, &ext); err != nil || ext == nil {
				return nil, fmt.Errorf("codec: %s.ext must contain a JSON object", descriptor.FullName())
			}
			if encode {
				object[key] = value
			} else {
				object[key], _ = json.Marshal(string(value))
			}
		} else if field.Kind() == protoreflect.MessageKind {
			if field.IsList() {
				var list []json.RawMessage
				if err := json.Unmarshal(value, &list); err != nil {
					return nil, err
				}
				for i := range list {
					v, err := wireJSON(list[i], field.Message(), encode)
					if err != nil {
						return nil, err
					}
					list[i] = v
				}
				object[key], _ = json.Marshal(list)
			} else {
				v, err := wireJSON(value, field.Message(), encode)
				if err != nil {
					return nil, err
				}
				object[key] = v
			}
		} else if field.Kind() == protoreflect.DoubleKind || field.Kind() == protoreflect.Int32Kind || field.Kind() == protoreflect.EnumKind {
			// ProtoJSON's string NaN/Infinity values are not OpenRTB numbers.
			var values []json.RawMessage
			if field.IsList() {
				if err := json.Unmarshal(value, &values); err != nil {
					return nil, err
				}
			} else {
				values = []json.RawMessage{value}
			}
			for _, value := range values {
				var number float64
				if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
					return nil, fmt.Errorf("codec: null array element")
				}
				if err := json.Unmarshal(value, &number); err != nil {
					return nil, fmt.Errorf("codec: %s must be a finite number", field.FullName())
				}
			}
		}
	}
	return json.Marshal(object)
}
