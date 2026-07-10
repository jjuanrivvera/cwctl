package api

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_UnmarshalStringOrNumber(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want ID
	}{
		{"string", `"abc"`, "abc"},
		{"int", `42`, "42"},
		{"big int beyond 2^53", `9007199254740993`, "9007199254740993"},
		{"null", `null`, ""},
		{"empty string", `""`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var id ID
			require.NoError(t, json.Unmarshal([]byte(tc.in), &id))
			assert.Equal(t, tc.want, id)
		})
	}
}

func TestID_MarshalAlwaysString(t *testing.T) {
	b, err := json.Marshal(ID("123"))
	require.NoError(t, err)
	assert.Equal(t, `"123"`, string(b))
}

func TestID_UnmarshalInvalid(t *testing.T) {
	var id ID
	assert.Error(t, json.Unmarshal([]byte(`{bad}`), &id))
}

func TestFlexInt(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    FlexInt
		wantErr bool
	}{
		{"number", `25`, 25, false},
		{"string number", `"25"`, 25, false},
		{"contacts current_page string", `"1"`, 1, false},
		{"float", `25.0`, 25, false},
		{"big int64", `9007199254740993`, 9007199254740993, false},
		{"null", `null`, 0, false},
		{"empty string", `""`, 0, false},
		{"words", `"abc"`, 0, true},
		{"NaN string", `"NaN"`, 0, true},
		{"Inf string", `"Inf"`, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var n FlexInt
			err := json.Unmarshal([]byte(tc.in), &n)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, n)
		})
	}
}

func TestFlexBool(t *testing.T) {
	cases := []struct {
		in      string
		want    FlexBool
		wantErr bool
	}{
		{`true`, true, false},
		{`false`, false, false},
		{`1`, true, false},
		{`0`, false, false},
		{`"true"`, true, false},
		{`"YES"`, true, false},
		{`"no"`, false, false},
		{`""`, false, false},
		{`null`, false, false},
		{`"maybe"`, false, true},
		{`2`, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			var b FlexBool
			err := json.Unmarshal([]byte(tc.in), &b)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, b)
		})
	}
}

func TestFlexTime(t *testing.T) {
	t.Run("unix seconds number", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`1735603200`), &ft))
		assert.Equal(t, time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), ft.Time)
	})
	t.Run("fractional unix seconds", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`1735603200.5`), &ft))
		assert.Equal(t, int64(1735603200), ft.Unix())
	})
	t.Run("RFC3339 string", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`"2024-12-31T00:00:00.000Z"`), &ft))
		assert.Equal(t, int64(1735603200), ft.Unix())
	})
	t.Run("date only", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`"2024-12-31"`), &ft))
		assert.Equal(t, 2024, ft.Year())
	})
	t.Run("unix seconds as string", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`"1735603200"`), &ft))
		assert.Equal(t, int64(1735603200), ft.Unix())
	})
	t.Run("null and empty", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`null`), &ft))
		assert.True(t, ft.IsZero())
		assert.Equal(t, "", ft.String())
		b, err := json.Marshal(ft)
		require.NoError(t, err)
		assert.Equal(t, `null`, string(b))
	})
	t.Run("marshal normalizes to RFC3339", func(t *testing.T) {
		var ft FlexTime
		require.NoError(t, json.Unmarshal([]byte(`1735603200`), &ft))
		b, err := json.Marshal(ft)
		require.NoError(t, err)
		assert.Equal(t, `"2024-12-31T00:00:00Z"`, string(b))
	})
	t.Run("garbage rejected", func(t *testing.T) {
		var ft FlexTime
		assert.Error(t, json.Unmarshal([]byte(`"not-a-time"`), &ft))
	})
}

func TestStringOrSlice(t *testing.T) {
	cases := []struct {
		in   string
		want StringOrSlice
	}{
		{`"x"`, StringOrSlice{"x"}},
		{`["x","y"]`, StringOrSlice{"x", "y"}},
		{`null`, nil},
		{`[]`, StringOrSlice{}},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			var s StringOrSlice
			require.NoError(t, json.Unmarshal([]byte(tc.in), &s))
			assert.Equal(t, tc.want, s)
		})
	}
	var s StringOrSlice
	assert.Error(t, json.Unmarshal([]byte(`123`), &s))
}

// Fuzz the flexible decoders: they must never panic and must either error or produce a
// value that re-marshals cleanly.

func FuzzID(f *testing.F) {
	for _, seed := range []string{`"a"`, `1`, `null`, `9007199254740993`, `-5`, `1.5`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var id ID
		if err := json.Unmarshal(data, &id); err == nil {
			if _, err := json.Marshal(id); err != nil {
				t.Fatalf("marshal after unmarshal: %v", err)
			}
		}
	})
}

func FuzzFlexInt(f *testing.F) {
	for _, seed := range []string{`1`, `"1"`, `null`, `""`, `1.9`, `"NaN"`, `9007199254740993`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var n FlexInt
		_ = json.Unmarshal(data, &n)
	})
}

func FuzzFlexTime(f *testing.F) {
	for _, seed := range []string{`1735603200`, `"2024-12-31T00:00:00Z"`, `null`, `"1735603200.5"`, `-1`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var ft FlexTime
		if err := json.Unmarshal(data, &ft); err == nil {
			if _, err := json.Marshal(ft); err != nil {
				t.Fatalf("marshal after unmarshal: %v", err)
			}
		}
	})
}

func FuzzStringOrSlice(f *testing.F) {
	for _, seed := range []string{`"x"`, `["x"]`, `null`, `[]`, `[1]`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var s StringOrSlice
		_ = json.Unmarshal(data, &s)
	})
}

func FuzzDecodeList(f *testing.F) {
	for _, seed := range []string{
		`[]`, `[{"id":1}]`, `{"payload":[{"id":1}]}`, `{"data":{"meta":{"count":1},"payload":[{"id":1}]}}`,
		`{"audit_logs":[{"id":1}],"total_entries":9}`, `{}`, `null`, `{"a":[],"b":[]}`,
	} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		type rec struct {
			ID ID `json:"id"`
		}
		_, _, _ = decodeList[rec](data)
	})
}
