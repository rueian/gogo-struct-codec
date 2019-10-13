package structbson

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"reflect"
)

var DefaultValueCodec = ValueCodec{}

type ValueCodec struct{}

func (c ValueCodec) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoValueType {
		return bsoncodec.ValueEncoderError{Name: "ValueCodec.EncodeValue", Types: []reflect.Type{ProtoValueType}, Received: val}
	}

	kindField := val.Field(0) // the 'Kind' field
	if kindField.IsNil() {
		return vw.WriteNull()
	}

	kv := kindField.Elem()

	var xv reflect.Value
	switch kv.Type() {
	case ProtoValueNullPtrType:
		return vw.WriteNull()
	default:
		xv = kv.Elem().Field(0) // the 'XXXValue' field
	}

	encoder, err := ec.LookupEncoder(xv.Type())
	if err != nil {
		return err
	}
	return encoder.EncodeValue(ec, vw, xv)
}

func (c ValueCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoValueType {
		return bsoncodec.ValueDecoderError{Name: "ValueCodec.DecodeValue", Types: []reflect.Type{ProtoValueType}, Received: val}
	}

	kindField := val.Field(0) // the 'Kind' field

	switch vr.Type() {
	case bsontype.Null:
		kindField.Set(reflect.ValueOf(&types.Value_NullValue{}))
		return vr.ReadNull()
	case bsontype.Undefined:
		kindField.Set(reflect.ValueOf(&types.Value_NullValue{}))
		return vr.ReadUndefined()
	case bsontype.Type(0):
		kindField.Set(reflect.ValueOf(&types.Value_NullValue{}))
		return nil
	case bsontype.EmbeddedDocument:
		value := &types.Value_StructValue{StructValue: &types.Struct{}}
		if err := DefaultStructCodec.DecodeValue(dc, vr, reflect.ValueOf(value.StructValue).Elem()); err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(value))
	case bsontype.Array:
		list := &types.Value_ListValue{ListValue: &types.ListValue{}}
		if err := DefaultListCodec.DecodeValue(dc, vr, reflect.ValueOf(list.ListValue).Elem()); err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(list))
	case bsontype.Double:
		v, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(&types.Value_NumberValue{NumberValue: v}))
	case bsontype.Int32:
		v, err := vr.ReadInt32()
		if err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(&types.Value_NumberValue{NumberValue: float64(v)}))
	case bsontype.Int64:
		v, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(&types.Value_NumberValue{NumberValue: float64(v)}))
	case bsontype.String:
		v, err := vr.ReadString()
		if err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(&types.Value_StringValue{StringValue: v}))
	case bsontype.Boolean:
		v, err := vr.ReadBoolean()
		if err != nil {
			return err
		}
		kindField.Set(reflect.ValueOf(&types.Value_BoolValue{BoolValue: v}))
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	}
	return nil
}
