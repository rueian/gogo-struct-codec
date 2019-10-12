package structbson

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"reflect"
)

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
	case bsontype.Type(0):
		kindField.Set(reflect.New(ProtoValueNullType))
	case bsontype.EmbeddedDocument:
		kindField.Set(reflect.New(ProtoValueStructType))
		valueField := kindField.Elem().Elem().Field(0)
		valueField.Set(reflect.New(ProtoStructType))
		structValue := valueField.Elem()
		decoder, err := dc.LookupDecoder(structValue.Type())
		if err != nil {
			return err
		}
		return decoder.DecodeValue(dc, vr, structValue)
	case bsontype.Array:
		kindField.Set(reflect.New(ProtoValueListType))
		valueField := kindField.Elem().Elem().Field(0)
		valueField.Set(reflect.New(ProtoListValueType))
		listValue := valueField.Elem()
		decoder, err := dc.LookupDecoder(listValue.Type())
		if err != nil {
			return err
		}
		return decoder.DecodeValue(dc, vr, listValue)
	case bsontype.Double:
		kindField.Set(reflect.New(ProtoValueNumberType))
		v, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		kindField.Elem().Elem().Field(0).SetFloat(v)
	case bsontype.Int32:
		kindField.Set(reflect.New(ProtoValueNumberType))
		v, err := vr.ReadInt32()
		if err != nil {
			return err
		}
		kindField.Elem().Elem().Field(0).SetFloat(float64(v))
	case bsontype.Int64:
		kindField.Set(reflect.New(ProtoValueNumberType))
		v, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		kindField.Elem().Elem().Field(0).SetFloat(float64(v))
	case bsontype.String:
		kindField.Set(reflect.New(ProtoValueStringType))
		v, err := vr.ReadString()
		if err != nil {
			return err
		}
		kindField.Elem().Elem().Field(0).SetString(v)
	case bsontype.Boolean:
		kindField.Set(reflect.New(ProtoValueBoolType))
		v, err := vr.ReadBoolean()
		if err != nil {
			return err
		}
		kindField.Elem().Elem().Field(0).SetBool(v)
	case bsontype.Null:
		kindField.Set(reflect.New(ProtoValueNullType))
		return vr.ReadNull()
	case bsontype.Undefined:
		kindField.Set(reflect.New(ProtoValueNullType))
		return vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	}
	return nil
}
