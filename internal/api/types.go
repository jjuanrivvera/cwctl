package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// ID unmarshals from a JSON string OR number and always marshals back as a string.
// Chatwoot ids are integers, but identifiers elsewhere (inbox identifier, source_id,
// conversation uuid) are strings; one flexible type lets every record render consistently
// and avoids float64 precision loss above 2^53 that a naive decode would suffer.
type ID string

func (id *ID) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*id = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*id = ID(s)
		return nil
	}
	// Number: json.Number keeps the exact textual form of large integers.
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	*id = ID(n.String())
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }

func (id ID) String() string { return string(id) }

// FlexInt accepts a JSON number or a numeric string and stores an int64. Chatwoot returns
// meta counts as numbers on some endpoints and strings on others (contacts
// meta.current_page is "1"). Int64 is parsed before Float64 so integers above 2^53 keep
// exact values; NaN/Inf and non-numeric strings are rejected.
type FlexInt int64

func (n *FlexInt) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	s := string(b)
	if b[0] == '"' {
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*n = 0
			return nil
		}
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		*n = FlexInt(v)
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf("FlexInt: %q is not a number", s)
	}
	*n = FlexInt(int64(f))
	return nil
}

// FlexBool accepts a real JSON bool, "true"/"false"/"1"/"0"/"yes"/"no" strings, or 0/1
// numbers — the usual drift across a Rails API's boolean-ish fields.
type FlexBool bool

func (fb *FlexBool) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch string(b) {
	case "", "null":
		*fb = false
		return nil
	case "true", "1":
		*fb = true
		return nil
	case "false", "0":
		*fb = false
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true", "1", "yes":
			*fb = true
			return nil
		case "false", "0", "no", "":
			*fb = false
			return nil
		}
		return fmt.Errorf("FlexBool: %q is not a boolean", s)
	}
	return fmt.Errorf("FlexBool: %s is not a boolean", string(b))
}

// FlexTime accepts Chatwoot's two timestamp dialects — unix seconds as a number
// (conversations.timestamp, messages.created_at, sometimes fractional) and ISO-8601 /
// RFC3339 strings (contacts.created_at) — and normalizes to a time.Time that marshals
// back as RFC3339.
type FlexTime struct {
	time.Time
}

func (t *FlexTime) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" || string(b) == `""` {
		t.Time = time.Time{}
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05 MST", "2006-01-02"} {
			if parsed, err := time.Parse(layout, s); err == nil {
				t.Time = parsed
				return nil
			}
		}
		// Unix seconds arriving as a string.
		if secs, err := strconv.ParseFloat(s, 64); err == nil {
			t.Time = timeFromUnixFloat(secs)
			return nil
		}
		return fmt.Errorf("FlexTime: unrecognized time %q", s)
	}
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf("FlexTime: invalid numeric time")
	}
	t.Time = timeFromUnixFloat(f)
	return nil
}

func timeFromUnixFloat(secs float64) time.Time {
	sec, frac := math.Modf(secs)
	return time.Unix(int64(sec), int64(frac*1e9)).UTC()
}

func (t FlexTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.UTC().Format(time.RFC3339))
}

func (t FlexTime) String() string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// StringOrSlice accepts a JSON string ("x") or an array of strings (["x","y"]) and
// normalizes to a slice — a common shape drift in real-world APIs.
type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*s = nil
		return nil
	}
	if b[0] == '[' {
		var arr []string
		if err := json.Unmarshal(b, &arr); err != nil {
			return err
		}
		*s = arr
		return nil
	}
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = []string{str}
	return nil
}
