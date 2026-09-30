# View 使用示例

三语言均使用生成模型。视图构造执行基础校验；BidCheck 负责请求与响应的关系检查；完整 JSON Schema 独立可选。分层与迁移见 [sdk.md](sdk.md)。

## Go

```go
import (
    "fmt"

    "github.com/oakrtb/oakrtb/sdk/go/builder"
    "github.com/oakrtb/oakrtb/sdk/go/codec"
    "github.com/oakrtb/oakrtb/sdk/go/bidcheck"
    "github.com/oakrtb/oakrtb/sdk/go/view"
    "github.com/oakrtb/oakrtb/sdk/go/validation"
)

req, err := codec.UnmarshalBidRequest(requestBytes)
if err != nil { return err }
requestView, err := view.NewRequest(req)
if err != nil { return err }

imp := requestView.FindImp("1")
if imp == nil { return fmt.Errorf("missing imp 1") }
if ready := validation.ImpReadyMtype(imp.Imp, 1); !ready.OK() {
    return fmt.Errorf("imp is not ready: %v", ready.Issues)
}
res, err := builder.NewBidResponse(requestView.AuctionID()).
    Currency("USD").
    AddSeatBid("seat", builder.NewBid("b", "1", 2).Banner().Build()).Build()
if err != nil { return err }
result := bidcheck.Response(requestView, res)
if !result.OK() { return fmt.Errorf("response mismatch: %v", result.Errors()) }
// Apply business policy to result.Warnings(); OK does not mean policy approval.
responseBytes, err := codec.MarshalJSON(res)
```

片段位于接收 `requestBytes` 的函数中，按实际返回类型处理结果。构建 View 后不要修改 req 或其子对象。若需要修改原始请求，改用 `view.NewRequestCopy(req)`。

## Java

```java
import com.oakrtb.sdk.codec.Json;
import com.oakrtb.sdk.builder.BidResponseBuilder;
import com.oakrtb.sdk.view.RequestView;
import com.oakrtb.sdk.view.ResponseView;
import com.oakrtb.sdk.bidcheck.BidCheck;
import com.oakrtb.openrtb.v2.Bid;

var req = Json.parseBidRequest(requestJson);
var requestView = RequestView.of(req);
var bid = Bid.newBuilder().setId("b").setImpid("1").setPrice(2).setMtypeValue(1).build();
var res = BidResponseBuilder.create(requestView.auctionId())
    .currency("USD").addSeatBid("seat", bid).build();
var result = BidCheck.response(requestView, res);
if (!result.ok()) throw new IllegalArgumentException(result.errors().toString());
// Apply business policy to result.warnings().
var responseView = ResponseView.of(res);
byte[] responseBytes = Json.toJsonBytes(responseView.response());
```

生成模型不可变，视图不依赖可变 builder。`RequestView.facts()` 返回 `RequestView.ImpFact`，响应侧对应 `ResponseView.BidFact`。

## Rust

```rust
use oakrtb_sdk::{builder::{BidBuilder, BidResponseBuilder}, codec, bidcheck, view::RequestView};

fn respond(request_bytes: &[u8]) -> Result<Vec<u8>, String> {
    let req = codec::parse_bid_request(request_bytes).map_err(|e| e.to_string())?;
    let request_view = RequestView::new(&req)?;
    let res = BidResponseBuilder::new(request_view.auction_id())
        .currency("USD")
        .add_seat_bid("seat", vec![BidBuilder::new("b", "1", 2.0).banner().build()])
        .build()?;
    let result = bidcheck::response(&request_view, &res);
    if !result.ok() { return Err(format!("response mismatch: {:?}", result.issues)); }
    // Apply business policy to result.warnings().
    codec::to_json(&res).map_err(|e| e.to_string())
}
```

同一个 req 也可由 `prost::Message::decode` 得到，View/BidCheck 无需变更。Builder 的 `build()` 返回模型；需要 JSON 字节时使用 `build_json()`。

## 查询能力

| 能力 | Go | Java | Rust |
|---|---|---|---|
| 查找展示位 | `FindImp` | `findImp` | `find_imp` |
| 按格式筛选 | `ImpsWith` | `impsWith` | `imps_with` |
| 查找出价 | `FindBid` | `findBid` | `find_bid` |
| 查找展示位的出价 | `BidsForImp` | `bidsForImp` | `bids_for_imp` |
| 紧凑字段摘要 | `Facts` | `facts` | `facts` |

`MarkupMask` 表示 banner/video/audio/native，与 Banner.format 尺寸数组不同。`Inventory` 表示 site/app/dooh，与 Content.Channel 不同。请求可包含多种格式；多格式出价须指明 mtype。

## 校验与 HTTP 边界

- 解码失败：报文语法或模型字段类型不合法。
- 基础校验失败：缺 id/at/cur/imp、库存互斥冲突、无效价格等。
- BidCheck ERROR：请求 ID 不匹配、impid 不存在、mtype 不匹配等。
- BidCheck WARN：低于底价、屏蔽项、货币不匹配、超过建议截止时间等，是否拒绝由业务决定。

响应币种与 imp.bidfloorcur 均非空且相等才比较底价，不隐式补 USD；单条 BidCheck.bid 没有响应币种，跳过底价。no-bid 仍核对响应请求 ID。

SDK 不实现 HTTP 服务、拍卖引擎或政策执行。HTTP 204 无正文表示 no-bid；结构化 no-bid 可返回带 id/cur/nbr 的 JSON。传输头、压缩、状态码详见 [transport.md](transport.md)。
