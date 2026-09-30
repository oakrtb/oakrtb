# Contributing to OakRTB

感谢贡献。本仓库定义 **OpenRTB 兼容的竞价协议**（JSON Schema / OpenAPI / protobuf）与 Go / Java / Rust SDK；**不**实现 SSP、DSP 或拍卖引擎。

Thanks for contributing. This repo defines an **OpenRTB-compatible bidding protocol** and SDKs. It does **not** implement an exchange or bidder.

## 开发环境 / Dev setup

- Python 3.10+（CI 使用 3.12）（`scripts/validate.py`）
- `protoc`（`make proto-check`）
- Go **1.25+**、JDK **21**、Rust **1.88+**（跑 SDK 测试）

```bash
python3 -m pip install -r scripts/requirements.txt
make validate
make proto-check
make sdk-test          # Or sdk-test-go / sdk-test-java / sdk-test-rust
```

改 `schema/jsonschema/` 或 `proto/` 后必须：

```bash
python3 scripts/field_docs.py --write  # Regenerate field reference and schema descriptions.
make sync-schemas      # Refresh vendored schema/proto copies.
make proto-go          # Regenerate committed Go models when proto changes.
make sdk-test
```

字段或注释改变时，先更新 `docs/field-descriptions.json` 中对应的中文 `zh` 和英文原文 `source`。生成器会拒绝遗漏或未同步的释义；新增 Native 属性也必须补充中文说明。`make validate` 和 CI 会检查字段覆盖、Schema 描述与生成文档是否同步；自动检查不能代替语义审阅。

修改 proto 前需安装 `protoc-gen-go`（版本应与 `sdk/go/go.mod` 中的 protobuf 依赖匹配），并确保它在 `PATH` 中。提交生成的 `sdk/go/oakrtb/v2/`，以及更新后的 `sdk/go/jsonschema/schemas/`、`sdk/java/src/main/{resources,proto}/`、`sdk/rust/{schemas,proto}/`。

发包（crates.io / Maven Central）见 [docs/publishing.md](docs/publishing.md)。

## 改什么、怎么改

权威顺序见 [docs/versioning.md](docs/versioning.md)：

1. `schema/jsonschema/`（JSON 校验权威）
2. `proto/oakrtb/v2/openrtb.proto`（若字段走 protobuf）
3. `openapi/openrtb.yaml`（若改 HTTP 语义 / 示例 / 头）
4. `docs/`、`examples/`、`testdata/`
5. `CHANGELOG.md` + 必要时 `VERSION`

入站流程：`codec → 模型 → view → 业务决策`；出站使用 `builder → bidcheck → codec`，由业务处理检查结果。`view` 构造时执行 `validation` 基础校验；`jsonschema` 提供独立的完整合同校验，可按边界需求选择。模块依赖见 [SDK 架构](docs/sdk.md)。

破坏性 SDK / 协议变更须在 CHANGELOG 标明 **Breaking**，并遵循 SemVer（见 [版本策略](docs/versioning.md)）。

## Pull request

- 从最新 `main` 开分支；一个 PR 聚焦一件事
- 核心检查命令：`make validate && make proto-check && make sdk-test`
- 描述里写清：动机、破坏性与否、如何验证
- 勿提交密钥、本地 IDE 杂项；`gen/` 与 `target/` 为忽略的本地构建产物；仅提交指定的生成模型与协议副本

## 行为准则

请保持专业、尊重。严重问题可私下联系仓库维护者（见 [SECURITY.md](SECURITY.md)）。

## 许可

贡献默认按仓库 [LICENSE](LICENSE)（Apache-2.0）授权。对象名 / JSON 字段名仍遵循 IAB OpenRTB（见 [NOTICE](NOTICE)）；本项目为**独立实现，非 IAB 官方**。
