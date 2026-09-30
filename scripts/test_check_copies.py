import tempfile
import unittest
from pathlib import Path
from check_copies import check_copies, copy_pairs


class CopyCheckTest(unittest.TestCase):
    def test_detects_drift_without_rewriting_it(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            for name in ("schema/jsonschema/test.json", "proto/oakrtb/v2/openrtb.proto"):
                source = root / name
                source.parent.mkdir(parents=True, exist_ok=True)
                source.write_text("canonical")
            for source, target in copy_pairs(root):
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(source.read_bytes())
            self.assertEqual([], check_copies(root))
            targets = [target for _, target in copy_pairs(root)]
            for target in targets:
                target.write_text("stale")
                self.assertTrue(check_copies(root))
                self.assertEqual("stale", target.read_text())
                target.write_text("canonical")
            targets[0].unlink()
            self.assertTrue(check_copies(root))


if __name__ == "__main__":
    unittest.main()
