package bidcheck

import (
	rt "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"
	"github.com/oakrtb/oakrtb/sdk/go/validation"
	"github.com/oakrtb/oakrtb/sdk/go/view"
	"google.golang.org/protobuf/proto"
	"strconv"
	"testing"
)

var benchmarkResult validation.CheckResult

func BenchmarkResponse(b *testing.B) {
	for _, size := range []int{1, 128} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			req := &rt.BidRequest{Id: "a", At: proto.Int32(1), Cur: []string{"USD"}}
			res := &rt.BidResponse{Id: "a", Cur: "USD", Seatbid: []*rt.SeatBid{{}}}
			for i := 0; i < size; i++ {
				id := strconv.Itoa(i)
				req.Imp = append(req.Imp, &rt.Imp{Id: id, Banner: &rt.Banner{}})
				res.Seatbid[0].Bid = append(res.Seatbid[0].Bid, &rt.Bid{Id: id, Impid: id, Price: proto.Float64(1), Adomain: []string{"safe"}, Cat: []string{"safe"}, Bundle: "safe"})
			}
			for i := 0; i < 256; i++ {
				s := strconv.Itoa(i)
				req.Badv = append(req.Badv, s)
				req.Bcat = append(req.Bcat, s)
				req.Bapp = append(req.Bapp, s)
			}
			snap, err := view.NewRequest(req)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkResult = Response(snap, res)
			}
		})
	}
}
