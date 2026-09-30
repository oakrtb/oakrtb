package builder

import (
	"github.com/oakrtb/oakrtb/sdk/go/codec"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDumpFullFixtures(t *testing.T) {
	if os.Getenv("DUMP_FIXTURES") != "1" {
		t.Skip("set DUMP_FIXTURES=1 to write JSON fixtures")
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// sdk/go/builder/dump_fixtures_test.go → repo root
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "testdata", "full"))
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Log("writing to", root)
	write := func(name string, raw []byte) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, raw, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := codec.MarshalJSON(fullBidRequestWithApp())
	if err != nil {
		t.Fatal(err)
	}
	write("bid-request/app.json", raw)
	raw, err = codec.MarshalJSON(fullBidRequestWithSite())
	if err != nil {
		t.Fatal(err)
	}
	write("bid-request/site.json", raw)
	raw, err = codec.MarshalJSON(fullBidRequestWithDooh())
	if err != nil {
		t.Fatal(err)
	}
	write("bid-request/dooh.json", raw)
	raw, err = codec.MarshalJSON(fullBidResponse())
	if err != nil {
		t.Fatal(err)
	}
	write("bid-response/full.json", raw)
}
