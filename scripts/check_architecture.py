#!/usr/bin/env python3
"""Check SDK production dependency boundaries (tests and compatibility entry points stay local)."""
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
ALLOWED = {
    "codec": set(),
    "validation": set(),
    "jsonschema": set(),
    "builder": {"codec", "validation", "jsonschema"},
    "view": {"validation"},
    "bidcheck": {"view", "validation"},
}


def dependencies(language: str, source: str) -> set[str]:
    # Tests may construct models with builders without introducing a production edge.
    source = re.split(r"#\[cfg\((?:all\()?test", source)[0]
    source = re.sub(r"/\*.*?\*/", "", source, flags=re.S)
    source = re.sub(r"^\s*//.*$", "", source, flags=re.M)
    if language == "go":
        return set(re.findall(r'"github.com/oakrtb/oakrtb/sdk/go/(\w+)', source)) & ALLOWED.keys()
    if language == "java":
        return set(re.findall(r"\bcom\.oakrtb\.sdk\.(\w+)", source)) & ALLOWED.keys()
    found = set(re.findall(r"\bcrate::(\w+)", source))
    for group in re.findall(r"\buse crate::\{([^;]+)\};", source):
        found.update(re.findall(r"\b(codec|validation|jsonschema|builder|view|bidcheck)\b", group))
    return found & ALLOWED.keys()


def violations(layer: str, dependencies: set[str]) -> set[str]:
    return dependencies - ALLOWED[layer] - {layer}


def check(root: Path) -> list[str]:
    errors = []
    trees = [("go", root / "sdk/go", "*.go"),
             ("java", root / "sdk/java/src/main/java/com/oakrtb/sdk", "*.java"),
             ("rust", root / "sdk/rust/src", "*.rs")]
    for language, tree, glob in trees:
        for path in tree.rglob(glob):
            if path.name.endswith(("_test.go", "_test.rs")):
                continue
            relative = path.relative_to(tree)
            layer = relative.parts[0].split(".")[0]
            if layer not in ALLOWED:
                continue
            for dependency in sorted(violations(layer, dependencies(language, path.read_text()))):
                errors.append(f"{path.relative_to(root)}: forbidden {layer} -> {dependency}")
    legacy_types = r"(?:Request|Response)(?:Pipeline|Snapshot|SharedView)"
    for language, tree, glob in trees:
        if language not in {"go", "rust"}:
            continue
        for path in tree.rglob(glob):
            if path.name.endswith(("_test.go", "_test.rs")):
                continue
            source = path.read_text()
            if re.search(r"\b(?:type|struct)\s+" + legacy_types + r"\b", source):
                errors.append(f"{path.relative_to(root)}: legacy view type is forbidden")
            if re.search(r"\b(?:func|fn)\s+(?:RunRequest(?:Copy)?|RunResponse(?:Copy)?|OfRequest|OfResponse|LightGateRequest|LightGateResponse|run_request|run_response|light_gate_request|light_gate_response)\s*\(", source):
                errors.append(f"{path.relative_to(root)}: legacy view factory is forbidden")
    for name in ("sdk/go/builder/json.go", "sdk/java/src/main/java/com/oakrtb/sdk/builder/Json.java"):
        if (root / name).exists():
            errors.append(f"{name}: legacy codec facade is forbidden")
    java_view = root / "sdk/java/src/main/java/com/oakrtb/sdk/view"
    for name in ("RequestViews", "ResponseViews", "RequestPipeline", "ResponsePipeline"):
        if (java_view / f"{name}.java").exists():
            errors.append(f"Java view must not reintroduce legacy {name}")
    return errors


if __name__ == "__main__":
    errors = check(ROOT)
    if errors:
        print("\n".join(errors), file=sys.stderr)
        sys.exit(1)
    print("SDK dependency boundaries verified")
