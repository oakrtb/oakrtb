# 接入示例

这里同时提供 JSON 报文和可运行的 Java、Go、Rust 接入代码。代码直接使用本仓库的 SDK，便于修改后验证，不要求先发布 SDK。

将 SDK 引入现有应用的步骤见 [接入指南](../docs/getting-started.md)。

## 选择示例

| 内容 | 入口 | 用途 |
|---|---|---|
| 简洁请求 | [bid-request/](bid-request/) | Banner、Video、Native 请求格式 |
| 简洁响应 | [bid-response/](bid-response/) | Banner 出价、结构化 no-bid |
| Go 完整流程 | [go/main.go](go/main.go) | 构建、编解码、校验、决策和响应查询 |
| Java 完整流程 | [java/AuctionExample.java](java/AuctionExample.java) | 同一流程的 Java API 用法 |
| Rust 完整流程 | [rust/src/main.rs](rust/src/main.rs) | 同一流程的 Rust API 用法 |

全字段报文和边界测试仍放在 [testdata/](../testdata/)。

## 完整流程

三个程序均在本地模拟一次 SSP → DSP → SSP 的报文往返，不启动 HTTP 服务、不请求外部地址。

1. SSP 构建网站 Banner 请求：300×250、USD、第一价格拍卖、底价 1 CPM、HTTPS 创意、测试流量。
2. 使用 `codec` 输出 OpenRTB JSON。DSP 入口执行完整 JSON Schema 校验，再解码为模型。
3. 构造 `RequestView`，执行基础校验并查找展示位。
4. 用 `builder` 构造包含创意、广告主域名、尺寸及 `mtype` 的响应。
5. 用 `bidcheck` 检查响应与请求的关系；ERROR 作为处理错误，WARN 按示例策略拒投。
6. 校验并序列化响应，SSP 解码后通过 `ResponseView` 识别出价或 no-bid。

每次运行依次验证以下三种情况，程序对结果进行检查，结果不符时以失败状态退出：

| 候选价格 | 预期结果 | 原因 |
|---|---|---|
| 2 CPM | `no_bid=false` | 高于底价，返回 Banner 出价 |
| 0.5 CPM | `no_bid=true` | 底价检查产生 WARN，示例策略拒投 |
| 0 | `no_bid=true` | 仅作为本例“无候选广告”的内部标记，不创建零价格 Bid |

标准输出包含请求和三次响应；标准错误包含低价拒投的诊断。

## 运行

以下命令从仓库根目录执行。首次构建需要下载依赖；Rust 构建需要 `protoc`。

### Go

需要 Go 1.25 或更高版本，复用 SDK 的 `go.mod`：

```sh
(cd sdk/go && go run ../../examples/go/main.go)
```

### Java

需要 JDK 21 和 Maven。先构建包含依赖的 jar，再使用 Java 源文件启动模式运行：

```sh
mvn -f sdk/java/pom.xml package -DskipTests
java --class-path "sdk/java/target/oakrtb-sdk-$(cat VERSION)-all.jar" examples/java/AuctionExample.java
```

### Rust

使用 Rust 1.88+ 和 `protoc`，并复用其构建缓存：

```sh
cargo run --locked --manifest-path examples/rust/Cargo.toml --target-dir sdk/rust/target
```

## 接入真实服务时

- 将 `respond` 的输入替换为 HTTP 请求正文，输出作为响应正文。解码或请求校验失败可映射为 HTTP 400；SDK 内部错误与业务无出价应分别处理。
- 本例使用 HTTP 200 对应的结构化 no-bid（`id`、`cur`、`nbr=0`）。也可返回 HTTP 204，此时正文必须为空；`nbr=0` 表示未知原因，应按实际业务选择更准确的标准原因码。
- 示例选择“所有 WARN 均拒投”以展示策略入口。生产中应按错误码配置策略；`OK/ok` 仅表示无 ERROR。
- 示例固定选择 `imp-1`、USD 和 300×250 创意。生产应选择支持的展示位、币种、尺寸和创意，并处理 PMP Deal、屏蔽项与隐私信号。
- 完整 Schema 校验适合接入联调或需要严格合同校验的入口；热路径可按需求选择。View 的基础校验和 BidCheck 的关联检查各有职责。
- 本例不设置 `tmax`，使本地三种场景输出稳定。实际服务需遵守上游超时预算，并将排队、解码及业务处理耗时计入整体预算。
- Go 的 View 借用原始模型，使用期间不要修改输入；需要隔离后续修改时使用 `NewRequestCopy`。Java 模型不可变，Rust 通过不可变借用约束生命周期。
- 创意中的 `.example` 地址为占位值。生产应替换为真实 HTTPS 素材和落地页，并补充实际设备、发布商、供应链及隐私信息。

协议细节见 [transport.md](../docs/transport.md)，模块分层见 [sdk.md](../docs/sdk.md)。
