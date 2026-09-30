# Publishing OakRTB SDKs

Authority stays in the monorepo (`schema/jsonschema/`, `proto/`). Language packages **vendor copies** via `make sync-schemas` so they can be published independently.

```bash
make sync-schemas   # refresh Go / Java / Rust vendored schema + proto
git add sdk/go/jsonschema/schemas sdk/java/src/main/resources sdk/java/src/main/proto sdk/rust/schemas sdk/rust/proto
```

## Go (`pkg.go.dev`)

Module path: `github.com/oakrtb/oakrtb/sdk/go`

The module lives in the `sdk/go` subdirectory, so its release tag must include that prefix: **`sdk/go/v0.2.0`**. A repository release tag such as `v0.2.0` alone does not publish this module version. See the [Go module version rules](https://go.dev/ref/mod#vcs-version).

After creating and pushing the module tag for the release commit:

```bash
go get github.com/oakrtb/oakrtb/sdk/go@v0.2.0
# Check whether the module version is available:
curl -s 'https://proxy.golang.org/github.com/oakrtb/oakrtb/sdk/go/@v/v0.2.0.info'
```

No extra publish step beyond the Git tag.

## Rust (`crates.io`)

Prerequisites:

1. [crates.io](https://crates.io) account linked to GitHub
2. API token: `cargo login`
3. Rust 1.88+ and `protoc` on `PATH` (for local dry-run / consumers’ first build)

```bash
make sync-schemas
make publish-rust-dry          # cargo publish --dry-run
cd sdk/rust && cargo publish   # real upload
```

Crate name: **`oakrtb-sdk`**. Includes vendored `schemas/*.json` and `proto/**/*.proto`.

## Java (Maven Central)

Prerequisites:

1. Sonatype Central Portal account: https://central.sonatype.com/  
   - Register namespace **`com.oakrtb`** (GitHub `oakrtb` org ownership verification)
2. Generate a **user token** on the portal
3. GPG key for signing (`gpg --gen-key`); publish public key to a keyserver

`~/.m2/settings.xml` (do not commit):

```xml
<settings>
  <servers>
    <server>
      <id>central</id>
      <username>YOUR_CENTRAL_USERNAME</username>
      <password>YOUR_CENTRAL_TOKEN</password>
    </server>
  </servers>
  <profiles>
    <profile>
      <id>gpg</id>
      <properties>
        <gpg.keyname>YOUR_KEY_ID</gpg.keyname>
      </properties>
    </profile>
  </profiles>
  <activeProfiles>
    <activeProfile>gpg</activeProfile>
  </activeProfiles>
</settings>
```

Deploy from monorepo:

```bash
make sync-schemas
cd sdk/java
mvn -Prelease clean deploy
```

Coordinates:

```xml
<groupId>com.oakrtb</groupId>
<artifactId>oakrtb-sdk</artifactId>
<version>0.2.0</version>
```

The `release` profile attaches sources + javadoc, GPG-signs, and uses
`central-publishing-maven-plugin`. The shaded `*-all.jar` is an optional classifier for fat-jar users; Central’s primary artifact is the thin jar.

Dry-run package (no deploy / may skip gpg if unset):

```bash
make publish-java-dry
```

## GitHub Actions (optional)

The existing [publish workflow](../.github/workflows/publish.yml) handles `v*` tags and manual dispatch. Both jobs currently have `if: ${{ false }}`. Configure the credentials below, then remove those gates when ready to enable publishing. Adding secrets alone does not enable the jobs. Go module tags (`sdk/go/v*`) are separate from the repository release tags.

| Secret | Used for |
|---|---|
| `CRATES_IO_TOKEN` | `cargo publish` |
| `CENTRAL_USERNAME` / `CENTRAL_TOKEN` | Maven Central |
| `GPG_PRIVATE_KEY` / `GPG_PASSPHRASE` | `maven-gpg-plugin` |

Do **not** commit tokens. Prefer manual first publish, then automate.

## Checklist per release

1. Bump `VERSION`, SDK `pom.xml` / `Cargo.toml` / `go.mod` docs, CHANGELOG
2. Run `make sync-schemas`; after proto changes, also run `make proto-go` and commit the generated Go models. Run `make validate proto-check sdk-test`.
3. Create the repository tag `vX.Y.Z` and Go module tag `sdk/go/vX.Y.Z` on the release commit, then create the GitHub Release
4. `cargo publish` + `mvn -Prelease deploy`
5. Update README install snippets if coordinates changed
