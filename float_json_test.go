package log

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
)

func TestJSONFieldsNonFiniteNumbers(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		var b bytes.Buffer
		writeObjectWithFields(&b, []byte(`{"message":"test"}`), []logField{{key: "value", val: v}})
		if !json.Valid(b.Bytes()) {
			t.Fatalf("invalid JSON: %s", b.String())
		}
		var got map[string]any
		if err := json.Unmarshal(b.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got["value"] != nil {
			t.Fatalf("got %v", got)
		}
	}
}
