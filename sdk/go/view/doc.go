// Package view provides read-only query views of BidRequest / BidResponse.
// NewRequest / NewResponse perform basic validation and construct RequestView / ResponseView,
// providing shared fields, impressions, bids, and summary queries. Each entry point validates and constructs the view once.
//
// Full JSON Schema validation is not performed; use the jsonschema package.
//
// Naming: view.MarkupMask is a bitmask of markup types on an Imp, distinct from the proto
// Format (Banner size entry); view.Inventory is the request inventory type (site/app/dooh),
// distinct from proto Content.Channel (content distribution channel).
//
// NewRequest/NewResponse borrow their input: the original object and view must remain read-only while in use.
// Use NewRequestCopy/NewResponseCopy if the input will be modified later; views themselves must still be treated as read-only.
package view
