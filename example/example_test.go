package example

import (
	"encoding/json"
	"github.com/gogo/protobuf/types"
	"github.com/rueian/gogo-struct-codec/ptypes"
	"reflect"
	"testing"
)

func TestMyMessage(t *testing.T) {
	m := MyMessage{Payload: &ptypes.Struct{}, Value: &ptypes.Value{}}
	m.Payload.Fields = map[string]*types.Value{
		"k": {Kind: &types.Value_StringValue{StringValue: "v"}},
	}
	m.Value.Val.Kind = &types.Value_ListValue{ListValue: &types.ListValue{
		Values: []*types.Value{
			{Kind: &types.Value_NullValue{NullValue: 0}},
			{Kind: &types.Value_StringValue{StringValue: "zzz"}},
			{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: map[string]*types.Value{
				"n": {Kind: &types.Value_NumberValue{NumberValue: 5}},
			}}}},
		},
	}}

	bs, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	if res := string(bs); res != `{"payload":{"k":"v"},"value":[null,"zzz",{"n":5}]}` {
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
