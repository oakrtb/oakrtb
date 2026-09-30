# OakRTB 接入指南

本指南面向将 OakRTB 接入现有 SSP、Exchange 或 DSP 的开发者。先按 [项目首页](../README.md#快速开始) 运行一个完整示例，再将其中的构建、解析与检查步骤接入自己的服务。

## 1. 引入 SDK

以下方式使用本地源码，适用于联调和尚未选择发布版本的项目。示例版本与当前仓库 `VERSION` 一致，为 `0.2.0`；升级时同步检查版本与 API 变更。

### Go

在你的 Go 项目目录中执行，将路径替换为实际的 OakRTB 绝对路径：

```sh
go mod edit -require=github.com/oakrtb/oakrtb/sdk/go@v0.0.0
go mod edit -replace=github.com/oakrtb/oakrtb/sdk/go=/absolute/path/to/oakrtb/sdk/go
```

添加 SDK 的 import 后执行 `go mod tidy`。这里的 `v0.0.0` 是本地替换所需的占位版本，不会下载对应远程版本。应用需要 Go 1.25 或更高版本。

```go
import (
    "github.com/oakrtb/oakrtb/sdk/go/builder"
    "github.com/oakrtb/oakrtb/sdk/go/codec"
    "github.com/oakrtb/oakrtb/sdk/go/view"
)
```

完整实现见 [Go 示例](../examples/go/main.go)。Go 模型已提交到仓库，正常使用 SDK 无需本地生成 protobuf。

### Java

使用 JDK 21，在 OakRTB 根目录安装 SDK 到本地 Maven 仓库：

```sh
mvn -f sdk/java/pom.xml install -DskipTests
```

随后在应用的 `pom.xml` 中添加依赖：

```xml
<dependency>
  <groupId>com.oakrtb</groupId>
  <artifactId>oakrtb-sdk</artifactId>
  <version>0.2.0</version>
</dependency>
```

Maven 应用使用普通 jar 及其传递依赖。需要直接以 classpath 运行时，可以用 `make jar` 生成包含依赖的 `gen/java/dist/oakrtb-sdk-0.2.0-all.jar`，用法见 [Java 示例](../examples/java/AuctionExample.java)。

### Rust

在应用的 `Cargo.toml` 中添加本地依赖，路径相对于该文件：

```toml
[dependencies]
oakrtb-sdk = { path = "../oakrtb/sdk/rust" }
```

使用 Rust 1.88+（建议当前 stable），并安装 `protoc`。SDK 在构建时从 protobuf 生成模型。完整实现见 [Rust 示例](../examples/rust/src/main.rs)。

如果应用只需要模型、编解码、Builder、View 和竞价检查，可以关闭默认的完整 Schema 校验功能：

```toml
[dependencies]
oakrtb-sdk = { path = "../oakrtb/sdk/rust", default-features = false }
```

关闭后不能使用 `jsonschema`、`ValidatedPayload` 或 `build_validated`；完整接入示例使用这些 API，需保留默认 feature。

公共包仓库的发布流程见 [publishing.md](publishing.md)。采用已发布包时，应确认版本实际存在且包含所需 API；本指南中的本地依赖不要求公共包发布。

## 2. 选择入站路径

| 输入 | 处理方式 |
|---|---|
| OpenRTB JSON | 可选完整 Schema 校验 → SDK codec 解码 → RequestView |
| protobuf 字节 | 原生 protobuf 解码 → RequestView |
| 应用内新建请求 | Builder 构建 → RequestView，或编码后发送 |

三种语言的常用 API 对照：

| 操作 | Go | Java | Rust |
|---|---|---|---|
| JSON 解码请求 | `codec.UnmarshalBidRequest` | `Json.parseBidRequest` | `codec::parse_bid_request` |
| 构造请求视图 | `view.NewRequest` | `RequestView.of` | `RequestView::new` |
| 查找展示位 | `FindImp` | `findImp` | `find_imp` |
| 构造响应 | `builder.NewBidResponse` | `BidResponseBuilder.create` | `BidResponseBuilder::new` |
| 检查整个响应 | `bidcheck.Response` | `BidCheck.response` | `bidcheck::response` |
| JSON 编码 | `codec.MarshalJSON` | `Json.toJsonBytes` | `codec::to_json` |
| 完整请求 Schema 校验 | `jsonschema.Request` | `Schema.request` | `jsonschema::request` |

Java 的 `Json` 位于 `com.oakrtb.sdk.codec`，`Schema` 位于 `com.oakrtb.sdk.jsonschema`。生成模型分别位于 Go 的 `oakrtb/v2`、Java 的 `com.oakrtb.openrtb.v2` 和 Rust 的 `oakrtb_sdk::proto`。

使用 SDK codec 处理 OpenRTB JSON：模型中的 `ext` 是 JSON 对象字符串，线上 JSON 中的 `ext` 是对象，codec 负责转换。原生 ProtoJSON 的输出不能直接替代这个过程。

## 3. 执行业务决策

`RequestView` 创建时执行基础校验，成功后可按 ID 查找展示位，或按 Banner、Video、Audio、Native 筛选。需要格式就绪检查时调用 `validation` 中的对应方法，基础结构合法不代表素材选择条件已经满足。

在构造 Bid 前，业务通常需要确定展示位、支持的币种、创意尺寸和类型、广告主域名、出价及 Deal。Bid 的 `impid` 必须对应请求中的展示位，响应 `id` 回显请求 `id`；多格式展示位应显式指定 `mtype`。

Builder 的 `Build/build` 执行构建器级检查并返回强类型模型。`BuildValidated/buildValidated/build_validated` 会额外生成 JSON 并运行完整 Schema 校验；调用后仍需检查返回的 Report。

## 4. 处理校验与竞价检查结果

| 检查 | 结果含义 | 应用处理 |
|---|---|---|
| codec 解码 | JSON 语法、字段类型等是否可解析 | 返回请求错误或记录入站失败 |
| 基础校验 / View 构造 | ID、必填结构、库存互斥、数值等是否可用 | 停止当前无效请求或响应的处理 |
| jsonschema | 报文是否满足完整 JSON 合同 | 检查 `Report.ok` 与 `errors` |
| bidcheck ERROR | 请求与响应的关联或竞价约束错误 | 拒绝该响应或修正后重新检查 |
| bidcheck WARN | 低价、屏蔽、超时或币种差异等风险 | 根据错误码执行明确的业务策略 |

`CheckResult.OK/ok` 只表示无 ERROR。例如低于底价可返回 WARN，同时 `ok` 仍为 true。[完整示例](../examples/README.md) 展示了如何把 WARN 转换为业务拒投。

优先检查整个响应。单 Bid 检查缺少响应币种和席位上下文，会跳过价格比较与 Deal 席位白名单。底价比较要求双方币种非空且一致，SDK 不执行汇率转换。PMP 的私有拍卖、交易白名单及固定价规则见 [SDK 文档](sdk.md#私有交易检查)。

## 5. 返回 HTTP 响应

SDK 的返回值由你的 HTTP 框架映射为状态码与正文：

| 结果 | HTTP 响应 |
|---|---|
| 正常出价 | `200`，正文为编码后的 BidResponse |
| 结构化 no-bid | `200`，正文包含 `id`、`cur` 和 `nbr` |
| 无正文 no-bid | `204`，不写响应正文 |
| 无法解析或请求校验失败 | `400`，可返回统一 Report |

协议路径为 `POST /openrtb/v2/auction`，部署时可由双方约定调整。JSON 使用 `Content-Type: application/json`；protobuf 使用 `application/x-protobuf`。HTTP 请求版本头为 `x-openrtb-version: 2.6`，与 SDK 包版本 `0.2.0` 是两个概念。

请求与响应采用相同编码。压缩、超时、回调和业务 HTTP 服务由接入方实现，完整约定见 [transport.md](transport.md)。

## 6. 数据与生命周期

- **缺失与零值**：optional 数值字段区分缺失和显式零值。Go 使用指针，Java 使用 `hasX`，Rust 使用 `Option`。
- **扩展字段**：自定义数据放入 `ext`；未知顶层字段不应作为持久保存业务扩展的方式。
- **Go View**：借用原始模型，使用期间保持只读。需要隔离后续输入修改时使用 `NewRequestCopy/NewResponseCopy`。
- **Java View**：读取不可变 protobuf 消息；修改原 Builder 不影响已构建消息。
- **Rust View**：借用不可变模型，由生命周期约束其使用。
- **时间预算**：View 在创建时固定自己的截止时间。HTTP 服务仍需计入接收、排队和解码耗时，并遵守上游 `tmax`。

## 7. 验证接入

先运行三种场景的 [完整示例](../examples/README.md)，再用真实接入方的报文验证字段、币种、素材与 Deal 约束。示例中的域名和创意地址是占位值，流量被标记为测试流量。

仓库根目录可执行：

```sh
make check-copies check-architecture
make validate
make sdk-test
```

`make validate` 需要先安装 `scripts/requirements.txt`。`make sdk-test` 需要三种语言的工具链，也可选择对应的 `sdk-test-go`、`sdk-test-java` 或 `sdk-test-rust`。

JSON 接入样例在 `examples/bid-request/` 与 `examples/bid-response/`；全字段、非法报文和跨语言一致性数据在 `testdata/`。修改 SDK 行为时，应同步检查受影响的样例与测试。

## 查询协议字段

字段用途、类型、必填性、尺寸单位和 Native 内嵌属性见 [完整字段手册](fields.md)；常用对象摘要见 [对象字典](objects.md)。
