use std::env;
use std::fs;
use std::path::{Path, PathBuf};

fn main() {
    println!("cargo:rerun-if-changed=build.rs");
    let manifest_dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());

    // Prefer vendored copies (crates.io); fall back to monorepo roots.
    let proto = first_existing(&[
        manifest_dir.join("proto/oakrtb/v2/openrtb.proto"),
        manifest_dir.join("../../proto/oakrtb/v2/openrtb.proto"),
    ])
    .expect("openrtb.proto (run make sync-schemas or build inside oakrtb monorepo)");
    let proto_root = proto
        .parent()
        .and_then(|p| p.parent())
        .and_then(|p| p.parent())
        .expect("proto root")
        .to_path_buf();

    println!("cargo:rerun-if-changed={}", proto.display());
    let mut config = prost_build::Config::new();
    let descriptors = config
        .load_fds(&[proto.as_path()], &[proto_root.as_path()])
        .expect("load proto descriptors");
    // JSON and protobuf share exactly the same generated model; no hand-maintained DTOs.
    for file in &descriptors.file {
        for message in &file.message_type {
            let name = format!(".{}.{}", file.package(), message.name());
            config.type_attribute(&name, "#[derive(serde::Serialize, serde::Deserialize)]");
            for field in &message.field {
                let path = format!("{}.{}", name, field.name());
                let skip = if field.label() as i32 == 3 {
                    "Vec::is_empty"
                } else if field.proto3_optional.unwrap_or(false) || field.r#type() as i32 == 11 {
                    "Option::is_none"
                } else if field.r#type() as i32 == 9 {
                    "String::is_empty"
                } else {
                    "crate::model_json::is_zero"
                };
                config.field_attribute(
                    &path,
                    format!("#[serde(default, skip_serializing_if = \"{skip}\")]"),
                );
                let decode = if field.r#type() as i32 == 11 {
                    if field.label() as i32 == 3 {
                        "crate::model_json::repeated_message"
                    } else {
                        "crate::model_json::optional_message"
                    }
                } else if field.r#type() as i32 == 5 || field.r#type() as i32 == 14 {
                    if field.label() as i32 == 3 {
                        "crate::model_json::repeated_i32"
                    } else if field.proto3_optional.unwrap_or(false) {
                        "crate::model_json::optional_i32"
                    } else {
                        "crate::model_json::scalar_i32"
                    }
                } else {
                    "crate::model_json::default_null"
                };
                if field.name() != "ext" {
                    config.field_attribute(
                        &path,
                        format!("#[serde(deserialize_with = \"{decode}\")]"),
                    );
                }
                if field.name() == "ext" {
                    config.field_attribute(&path, "#[serde(with = \"crate::model_json::ext\")]");
                } else if field.r#type() as i32 == 1 {
                    config.field_attribute(
                        &path,
                        "#[serde(serialize_with = \"crate::model_json::finite_optional\")]",
                    );
                }
            }
        }
    }
    config.compile_fds(descriptors).expect("prost compile");

    if env::var_os("CARGO_FEATURE_JSONSCHEMA").is_some() {
        sync_schemas(&manifest_dir);
    }
}

fn first_existing(candidates: &[PathBuf]) -> Option<PathBuf> {
    candidates.iter().find_map(|p| p.canonicalize().ok())
}

/// Copy JSON schemas into `OUT_DIR/schemas/` for `include_str!`.
fn sync_schemas(manifest_dir: &Path) {
    let src = first_existing(&[
        manifest_dir.join("schemas"),
        manifest_dir.join("../../schema/jsonschema"),
    ])
    .expect("schemas/ (run make sync-schemas or build inside oakrtb monorepo)");
    let out = PathBuf::from(env::var("OUT_DIR").unwrap()).join("schemas");
    fs::create_dir_all(&out).expect("create OUT_DIR/schemas");

    for name in [
        "openrtb.schema.json",
        "bid-request.schema.json",
        "bid-response.schema.json",
        "native.schema.json",
        "validation-result.schema.json",
    ] {
        let from = src.join(name);
        println!("cargo:rerun-if-changed={}", from.display());
        fs::copy(&from, out.join(name)).unwrap_or_else(|e| {
            panic!("copy {name} from {}: {e}", from.display());
        });
    }
}
