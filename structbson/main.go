package structbson

import (
	"github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"reflect"
)

var (
	ProtoStructType       = reflect.TypeOf(types.Struct{})
	ProtoValueType        = reflect.TypeOf(types.Value{})
	ProtoListValueType    = reflect.TypeOf(types.ListValue{})
	ProtoValueStructType  = reflect.TypeOf(types.Value_StructValue{})
	ProtoValueNumberType  = reflect.TypeOf(types.Value_NumberValue{})
	ProtoValueStringType  = reflect.TypeOf(types.Value_StringValue{})
	ProtoValueBoolType    = reflect.TypeOf(types.Value_BoolValue{})
	ProtoValueNullType    = reflect.TypeOf(types.Value_NullValue{})
	ProtoValueListType    = reflect.TypeOf(types.Value_ListValue{})
	ProtoValueNullPtrType = reflect.TypeOf(&types.Value_NullValue{})
)

var DefaultRegistry *bsoncodec.Registry

func init() {
	rb := bson.NewRegistryBuilder()
	rb.RegisterCodec(ProtoStructType, StructCodec{})
	rb.RegisterCodec(ProtoValueType, ValueCodec{})
	rb.RegisterCodec(ProtoListValueType, ListCodec{})
	DefaultRegistry = rb.Build()
}
