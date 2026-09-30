package codec_test

import (
	"bytes"
	"encoding/json"
	"github.com/oakrtb/oakrtb/sdk/go/bidcheck"
	"github.com/oakrtb/oakrtb/sdk/go/builder"
	"github.com/oakrtb/oakrtb/sdk/go/codec"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"github.com/oakrtb/oakrtb/sdk/go/view"
	"google.golang.org/protobuf/proto"
	"os"
	"reflect"
	"sort"
	"testing"
)

type example struct {
	Mtype                              int32
	Code                               string
	BannerW                            *int32 `json:"banner_w"`
	BannerH                            *int32 `json:"banner_h"`
	NativeRequest                      string `json:"native_request"`
	HasAdm                             bool   `json:"has_adm"`
	Name, Kind                         string
	Input, Expected, Request, Response json.RawMessage
	OK                                 bool `json:"ok"`
	BuilderOK                          bool `json:"builder_ok"`
	Codes                              []string
}

func decode(kind string, raw []byte) (proto.Message, error) {
	if kind == "request" {
		return codec.UnmarshalBidRequest(raw)
	}
	return codec.UnmarshalBidResponse(raw)
}
func jsonValue(t *testing.T, raw []byte) any {
	t.Helper()
	var v any
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		t.Fatal(err)
	}
	return v
}
func TestSharedConformance(t *testing.T) {
	raw, err := os.ReadFile("../../../testdata/conformance/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var groups map[string][]example
	if err = json.Unmarshal(raw, &groups); err != nil {
		t.Fatal(err)
	}
	for group, cases := range groups {
		for _, c := range cases {
			t.Run(group+"/"+c.Name, func(t *testing.T) {
				switch group {
				case "roundtrip":
					model, err := decode(c.Kind, c.Input)
					if err != nil {
						t.Fatal(err)
					}
					wire, err := proto.Marshal(model)
					if err != nil {
						t.Fatal(err)
					}
					back := model.ProtoReflect().New().Interface()
					if err = proto.Unmarshal(wire, back); err != nil {
						t.Fatal(err)
					}
					encoded, err := codec.MarshalJSON(back)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(jsonValue(t, encoded), jsonValue(t, c.Expected)) {
						t.Fatalf("got %s, expected %s", encoded, c.Expected)
					}
				case "reject":
					if _, err := decode(c.Kind, c.Input); err == nil {
						t.Fatal("invalid JSON accepted")
					}
				case "basic":
					var err error
					if c.Kind == "request" {
						r, e := codec.UnmarshalBidRequest(c.Input)
						if e != nil {
							t.Fatal(e)
						}
						err = validation.Request(r)
					} else {
						r, e := codec.UnmarshalBidResponse(c.Input)
						if e != nil {
							t.Fatal(e)
						}
						err = validation.Response(r)
					}
					if (err == nil) != c.OK {
						t.Fatalf("unexpected basic validation: %v", err)
					}
				case "ready":
					req, e := codec.UnmarshalBidRequest(c.Request)
					if e != nil {
						t.Fatal(e)
					}
					result := validation.ImpReadyMtype(req.Imp[0], c.Mtype)
					codes := make([]string, 0, len(result.Issues))
					for _, issue := range result.Issues {
						codes = append(codes, issue.Code)
					}
					sort.Strings(codes)
					sort.Strings(c.Codes)
					if result.OK() != c.OK || !reflect.DeepEqual(codes, c.Codes) {
						t.Fatalf("unexpected readiness: %+v", result)
					}
					if validation.ImpReady(req.Imp[0]).OK() != c.BuilderOK {
						t.Fatal("unexpected all-format readiness")
					}
					_, err := builder.NewBidRequest(req.Id).AuctionType(req.GetAt()).Currency(req.Cur...).AddImp(req.Imp[0]).Build()
					if (err == nil) != c.BuilderOK {
						t.Fatalf("builder disagrees with readiness: %v", err)
					}

				case "facts":
					req, e := codec.UnmarshalBidRequest(c.Request)
					if e != nil {
						t.Fatal(e)
					}
					v, e := view.NewRequest(req)
					if e != nil {
						t.Fatal(e)
					}
					f := v.Facts()[0]
					if !reflect.DeepEqual(f.BannerW, c.BannerW) || !reflect.DeepEqual(f.BannerH, c.BannerH) || f.NativeRequest != c.NativeRequest {
						t.Fatalf("unexpected facts: %+v", f)
					}
					res, e := codec.UnmarshalBidResponse(c.Response)
					if e != nil {
						t.Fatal(e)
					}
					rv, e := view.NewResponse(res)
					if e != nil {
						t.Fatal(e)
					}
					if rv.Facts()[0].HasAdm != c.HasAdm {
						t.Fatal("hasAdm differs")
					}
				case "bidcheck":
					req, err := codec.UnmarshalBidRequest(c.Request)
					if err != nil {
						t.Fatal(err)
					}
					res, err := codec.UnmarshalBidResponse(c.Response)
					if err != nil {
						t.Fatal(err)
					}
					v, err := view.NewRequest(req)
					if err != nil {
						t.Fatal(err)
					}
					result := bidcheck.Response(v, res)
					if result.OK() != c.OK {
						t.Fatalf("wrong severity: %+v", result)
					}
					codes := make([]string, 0, len(result.Issues))
					for _, issue := range result.Issues {
						codes = append(codes, issue.Code)
					}
					sort.Strings(codes)
					sort.Strings(c.Codes)
					if !reflect.DeepEqual(codes, c.Codes) {
						t.Fatalf("got %v want %v", codes, c.Codes)
					}
				}
			})
		}
	}
}
