package ptypes

import (
	"bytes"
	"encoding/json"
	"github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func fixture() Struct {
	return Struct{Struct: types.Struct{Fields: map[string]*types.Value{
		"null":   {Kind: &types.Value_NullValue{NullValue: 0}},
		"number": {Kind: &types.Value_NumberValue{NumberValue: 123}},
		"string": {Kind: &types.Value_StringValue{StringValue: "456"}},
		"array": {Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: []*types.Value{
			{Kind: &types.Value_NullValue{NullValue: 0}},
			{Kind: &types.Value_NumberValue{NumberValue: 123}},
			{Kind: &types.Value_StringValue{StringValue: "456"}},
			{Kind: &types.Value_BoolValue{BoolValue: true}},
			{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: map[string]*types.Value{
				"nested1": {Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: map[string]*types.Value{
					"nested2": {Kind: &types.Value_StringValue{StringValue: "789"}},
				}}}},
			}}}},
			{Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: []*types.Value{
				{Kind: &types.Value_NumberValue{NumberValue: 10000}},
			}}}},
		}}}},
		"struct": {Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: map[string]*types.Value{
			"struct": {Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: []*types.Value{
				{Kind: &types.Value_NumberValue{NumberValue: 20000}},
			}}}},
		}}}},
	}}}
}

const fact = `{"array":[null,123,"456",true,{"nested1":{"nested2":"789"}},[10000]],"null":null,"number":123,"string":"456","struct":{"struct":[20000]}}`

func TestStruct_MarshalJSON(t *testing.T) {
	bs, err := json.Marshal(fixture())
	if err != nil {
		t.Error(err)
	}

	if string(bs) != fact {
		t.Errorf("result not match with fact, got: %s", string(bs))
	}
}

func TestStruct_UnmarshalJSON(t *testing.T) {
	c := Struct{}
	if err := json.Unmarshal([]byte(fact), &c); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(c, fixture()) {
		t.Fatalf("result not match with fact, got %v", c)
	}
}

func TestStruct_MarshalJSONPB(t *testing.T) {
	f := fixture()
	res, err := defaultJSONPBMarshaler.MarshalToString(&f)
	if err != nil {
		t.Error(err)
	}

	if res != fact {
		t.Errorf("result not match with fact, got: %s", res)
	}
}

func TestStruct_UnmarshalJSONPB(t *testing.T) {
	c := Struct{}
	if err := defaultJSONPBUnmarshaler.Unmarshal(bytes.NewBufferString(fact), &c); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(c, fixture()) {
		t.Fatalf("result not match with fact, got %v", c)
	}
}

func TestStruct_BSON(t *testing.T) {
	f := fixture()
	bs, err := bson.Marshal(f)
	if err != nil {
		t.Error(err)
	}

	c := Struct{}
	err = bson.Unmarshal(bs, &c)
	if err != nil {
		t.Error(err)
	}

	bs, err = json.Marshal(c)
	if err != nil {
		t.Error(err)
	}

	if string(bs) != fact {
		t.Errorf("result not match with fact, got: %s", string(bs))
	}
}

func TestStruct_SQL(t *testing.T) {
	v, err := fixture().Value()
	if err != nil {
		t.Error(err)
	}

	c := Struct{}
	if err = c.Scan(v); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(c, fixture()) {
		t.Fatalf("result not match with fact, got %v", c)
	}
}
