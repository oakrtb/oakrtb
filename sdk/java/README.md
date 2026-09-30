# oakrtb-sdk (Java)

OpenRTB-compatible models, builders, `codec`, basic `validation`, query `view`, optional `bidcheck`, and off-path `jsonschema` validation for [OakRTB](https://github.com/oakrtb/oakrtb).

```xml
<dependency>
  <groupId>com.oakrtb</groupId>
  <artifactId>oakrtb-sdk</artifactId>
  <version>0.2.0</version>
</dependency>
```

Requires **JDK 21+**. JSON Schema resources are vendored under `src/main/resources/schema/jsonschema/` (refreshed by `make sync-schemas`).

Build from the monorepo:

```bash
make jar   # → gen/java/dist/oakrtb-sdk-0.2.0*.jar
```

Publishing to Maven Central: see [docs/publishing.md](https://github.com/oakrtb/oakrtb/blob/main/docs/publishing.md).

See [SDK architecture and migration](../../docs/sdk.md) for the shared typed model, codec entry points, View API and migration requirements.

Start with the [integration guide](../../docs/getting-started.md) and [runnable examples](../../examples/README.md).
