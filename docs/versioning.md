# 版本策略

OakRTB 版本写在 `VERSION`，遵循 SemVer。

## 兼容 OpenRTB

线格式 JSON 字段名与 IAB OpenRTB 2.6 对齐。当前钉住的规范快照为 **`2.6-202606`**
（https://github.com/InteractiveAdvertisingBureau/openrtb2.x/releases/tag/2.6-202606）。

`x-openrtb-version` 仍填 `2.6`（IAB 线格式主次版本）；OakRTB 自己的 semver 写在 `VERSION`。
文档与 CHANGELOG 应同时写明所对齐的 dated snapshot。

## 什么算兼容变更

- 新增可选字段或对象
- 文档化既有、已允许的扩展值；新增枚举值需单独评估兼容性
- 文档澄清

这些只增加 `VERSION` 的 MINOR 或 PATCH，不改 HTTP 路径 `/openrtb/v2/auction`。

模型能解码未知枚举数值，不代表完整 Schema 或 BidCheck 接受它。例如 JSON Schema 将 `Bid.mtype` 限制为 1–4；BidCheck 对厂商值 ≥500 给出 WARN，对其他未知值给出 ERROR。新增枚举值只有在旧接收方的校验与业务处理允许时才可视为兼容。

## 什么算破坏性变更

- 删除字段
- 把可选改为必填
- 改变字段类型或含义（例如价格单位）
- `site`/`app`/`dooh` 互斥规则变化

稳定版本（1.0 起）的破坏性变更升 MAJOR；当前 0.x 阶段的破坏性变更在 MINOR 版本中发布，并在 CHANGELOG 标记 Breaking、提供迁移说明。SDK API 变更与 HTTP 协议版本分别评估；仅 HTTP 合同需要时才调整路径。

## Schema / Proto / OpenAPI

| 产物 | 角色 |
|---|---|
| `schema/jsonschema/` | **JSON 校验权威**（required、互斥、数值范围） |
| `proto/oakrtb/v2/openrtb.proto` | 二进制编解码；注释对齐语义，**不**强制 required / oneOf |
| `openapi/openrtb.yaml` | JSON HTTP 合同（`$ref` schema）；不建模 protobuf |

改对象时同一 PR 必须同时更新：

1. `schema/jsonschema/`
2. `proto/oakrtb/v2/openrtb.proto`（若该字段走 protobuf）
3. `openapi/openrtb.yaml`（若改 HTTP 语义、示例或头）
4. `docs/objects.md` / `docs/spec.md` / `docs/transport.md`
5. `examples/bid-request|bid-response/` 与 `testdata/invalid/`
6. `CHANGELOG.md`
7. 若改了 schema 或 proto：`make sync-schemas`（刷新 Go/Java/Rust vendored 副本）；修改 proto 后还需 `make proto-go`，提交 Go 模型及副本，再跑 `make validate proto-check sdk-test`
8. 发包步骤见 [publishing.md](publishing.md)
