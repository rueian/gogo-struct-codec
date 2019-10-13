package example

import (
	"encoding/json"
	"github.com/gogo/protobuf/types"
	"github.com/rueian/gogo-struct-codec/ptypes"
	"reflect"
	"testing"
)

func TestMyMessage(t *testing.T) {
	m := MyMessage{Payload: &ptypes.Struct{}}
	m.Payload.Fields = map[string]*types.Value{
		"k": {Kind: &types.Value_StringValue{StringValue: "v"}},
	}

	bs, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	if res := string(bs); res != `{"payload":{"k":"v"}}` {
		t.Errorf("json mismatched, got: %v", res)
	}

	c := MyMessage{}
	if err = json.Unmarshal(bs, &c); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(m, c) {
		t.Errorf("unmarshal back mismatched, got: %v", c)
	}
}
