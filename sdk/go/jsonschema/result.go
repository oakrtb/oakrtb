package jsonschema

import "encoding/json"

// Report contains JSON Schema validation results (when Ok is false, it can serve as an HTTP 400 response body).
type Report struct {
	Ok     bool    `json:"ok"`
	Errors []Issue `json:"errors"`
}

// Issue represents a single Schema validation failure.
type Issue struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

// OK returns a successful result.
func OK() Report {
	return Report{Ok: true, Errors: []Issue{}}
}

// Fail returns a failed result.
func Fail(errs ...Issue) Report {
	if len(errs) == 0 {
		errs = []Issue{}
	}
	return Report{Ok: false, Errors: errs}
}

// MarshalJSON ensures errors is always a JSON array, never null.
func (r Report) MarshalJSON() ([]byte, error) {
	type alias Report
	out := alias(r)
	if out.Errors == nil {
		out.Errors = []Issue{}
	}
	return json.Marshal(out)
}
