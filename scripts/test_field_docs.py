"""Regression tests for documentation drift detection and safe schema annotation."""
import copy
import json
import unittest

from field_docs import (PROTO, SCHEMA, NATIVE, TEXT, parse_proto, render,
                        schema_text, validate_descriptions)


class FieldDocsTest(unittest.TestCase):
    def setUp(self):
        self.proto = PROTO.read_text()
        self.fields = parse_proto(self.proto)
        self.translations = json.loads(TEXT.read_text())

    def test_missing_or_stale_translation_is_rejected(self):
        missing = copy.deepcopy(self.translations['fields'])
        del missing['Bid.mtype']
        with self.assertRaisesRegex(ValueError, 'coverage mismatch'):
            validate_descriptions(self.fields, missing)
        stale = copy.deepcopy(self.translations['fields'])
        stale['Bid.mtype']['source'] = 'Old requirement'
        with self.assertRaisesRegex(ValueError, 'Stale translation: Bid.mtype'):
            validate_descriptions(self.fields, stale)

    def test_missing_schema_field_is_rejected(self):
        schema = json.loads(SCHEMA.read_text())
        del schema['$defs']['Bid']['properties']['mtype']
        with self.assertRaisesRegex(ValueError, 'Proto/schema field mismatch'):
            render(self.proto, schema, json.loads(NATIVE.read_text()), self.translations)

    def test_new_native_field_requires_description(self):
        native = json.loads(NATIVE.read_text())
        native['$defs']['Asset']['properties']['new_field'] = {'type': 'string'}
        with self.assertRaisesRegex(ValueError, 'Native translation coverage mismatch'):
            render(self.proto, json.loads(SCHEMA.read_text()), native, self.translations)

    def test_nested_descriptions_and_number_literals_are_preserved(self):
        original = '{"$defs":{"Bid":{"properties":{"price":{"minimum":0.123456789012345678901,"allOf":[{"description":"nested"}]}}}}}'
        expected = json.loads(original)
        expected['$defs']['Bid']['properties']['price']['description'] = 'CPM price'
        updated = schema_text(original, expected)
        self.assertIn('0.123456789012345678901', updated)
        self.assertIn('"description":"nested"', updated)
        self.assertEqual(expected, json.loads(updated))
        self.assertEqual(updated, schema_text(updated, expected))

    def test_comment_macro_does_not_end_message(self):
        fields = parse_proto('message Bid {\n  // ${AUCTION_PRICE}\n  string nurl = 1;\n  int32 w = 2; // Width.\n}\n')
        self.assertEqual({'nurl', 'w'}, set(fields['Bid']))
        self.assertEqual('${AUCTION_PRICE}', fields['Bid']['nurl']['description'])


if __name__ == '__main__':
    unittest.main()
