import unittest
from check_architecture import dependencies, violations


class ArchitectureTest(unittest.TestCase):
    def test_detects_upward_edges(self):
        for language, source in [
            ("go", 'import "github.com/oakrtb/oakrtb/sdk/go/builder"'),
            ("java", "com.oakrtb.sdk.builder.Json.toJson(model);"),
            ("rust", "use crate::{builder, validation};"),
        ]:
            self.assertIn("builder", violations("codec", dependencies(language, source)))

    def test_validation_is_independent_of_views_and_bid_checks(self):
        self.assertEqual({"view", "bidcheck", "builder"}, violations("validation", {"view", "bidcheck", "builder"}))
        self.assertFalse(violations("bidcheck", {"view", "validation"}))

    def test_allows_declared_edges_and_ignores_tests(self):
        self.assertFalse(violations("builder", {"codec", "jsonschema", "validation"}))
        self.assertEqual(set(), dependencies("rust", "#[cfg(test)]\nmod tests { use crate::builder; }"))


if __name__ == "__main__":
    unittest.main()
