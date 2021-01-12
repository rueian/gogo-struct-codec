package ptypes

import (
	"bytes"
	"encoding/json"
	"testing"
)

var tests = []string{"null", "number", "string", "array", "struct", "bool"}

func TestValue_MarshalUnMarshalJSON(t *testing.T) {
	for _, ty := range tests {
		f := Value{Val: *fixture().Fields[ty]}
		bs, err := json.Marshal(&f)
		if err != nil {
			t.Error(err)
		}
		v := Value{}
		if err := json.Unmarshal(bs, &v); err != nil {
			t.Error(err)
		}
		if !v.Equal(f) {
			t.Errorf("%s type result not match with fact, got: %s", ty, string(bs))
		}
	}
}

func TestValue_MarshalUnmarshalJSONPB(t *testing.T) {
	for _, ty := range tests {
		f := Value{Val: *fixture().Fields[ty]}
		bs, err := defaultJSONPBMarshaler.MarshalToString(&f)
		if err != nil {
			t.Error(err)
		}
		v := Value{}
		if err := defaultJSONPBUnmarshaler.Unmarshal(bytes.NewBufferString(bs), &v); err != nil {
			t.Error(err)
		}
		if !v.Equal(f) {
			t.Errorf("%s type result not match with fact, got: %s", ty, string(bs))
		}
	}
}

func TestValue_SQL(t *testing.T) {
	for _, ty := range tests {
		f := Value{Val: *fixture().Fields[ty]}
		dv, err := f.Value()
		if err != nil {
			t.Error(err)
		}
		v := Value{}
		if err := v.Scan(dv); err != nil {
			t.Error(err)
		}
		if !v.Equal(f) {
			t.Errorf("%s type result not match with fact", ty)
		}
	}
}
