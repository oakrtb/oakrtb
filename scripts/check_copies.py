#!/usr/bin/env python3
"""Read-only check of publishable SDK copies against the canonical schemas/proto."""
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]


def copy_pairs(root: Path):
    for source in sorted((root / "schema/jsonschema").glob("*.json")):
        for directory in ("sdk/go/jsonschema/schemas", "sdk/java/src/main/resources/schema/jsonschema", "sdk/rust/schemas"):
            yield source, root / directory / source.name
    source = root / "proto/oakrtb/v2/openrtb.proto"
    for directory in ("sdk/java/src/main/proto", "sdk/rust/proto"):
        yield source, root / directory / "oakrtb/v2/openrtb.proto"


def check_copies(root: Path) -> list[str]:
    errors = []
    if not list((root / "schema/jsonschema").glob("*.json")):
        errors.append("canonical schemas are missing")
    expected = set()
    for source, target in copy_pairs(root):
        expected.add(target)
        if not source.is_file() or not target.is_file():
            errors.append(f"missing source or copy: {target.relative_to(root)}")
        elif source.read_bytes() != target.read_bytes():
            errors.append(f"outdated copy: {target.relative_to(root)}")
    for directory in ("sdk/go/jsonschema/schemas", "sdk/java/src/main/resources/schema/jsonschema", "sdk/rust/schemas"):
        for target in (root / directory).glob("*.json"):
            if target not in expected:
                errors.append(f"unexpected copy: {target.relative_to(root)}")
    return errors


if __name__ == "__main__":
    errors = check_copies(ROOT)
    if errors:
        print("\n".join(errors), file=sys.stderr)
        print("Run make sync-schemas and commit the updated copies.", file=sys.stderr)
        sys.exit(1)
    print("SDK schema/proto copies match canonical sources")
