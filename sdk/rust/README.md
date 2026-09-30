# oakrtb-sdk (Rust)

OpenRTB-compatible models, builders, `codec`, basic `validation`, query `view`, optional `bidcheck`, and off-path `jsonschema` validation for [OakRTB](https://github.com/oakrtb/oakrtb).

```toml
[dependencies]
oakrtb-sdk = "0.2.0"
```

Requires **Rust 1.88+** and a `protoc` on `PATH` for the first build (prost).

Schema JSON and `openrtb.proto` are **vendored** under `schemas/` and `proto/` (refreshed by `make sync-schemas` in the monorepo).

See the [repository docs](https://github.com/oakrtb/oakrtb/blob/main/docs/sdk.md) and [publishing guide](https://github.com/oakrtb/oakrtb/blob/main/docs/publishing.md).

See [SDK architecture and migration](../../docs/sdk.md) for the shared typed model, codec entry points, View API and migration requirements.

Start with the [integration guide](../../docs/getting-started.md) and [runnable examples](../../examples/README.md).
