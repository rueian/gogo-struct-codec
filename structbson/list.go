package structbson

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

var DefaultListCodec = ListCodec{}

type ListCodec struct{}

func (c ListCodec) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoListValueType {
		return bsoncodec.ValueEncoderError{Name: "ListCodec.EncodeValue", Types: []reflect.Type{ProtoListValueType}, Received: val}
	}

	valuesField := val.Field(0) // the 'Values' field

	encoder, err := ec.LookupEncoder(valuesField.Type())
	if err != nil {
		return err
	}
	return encoder.EncodeValue(ec, vw, valuesField)
}

func (c ListCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoListValueType {
		return bsoncodec.ValueDecoderError{Name: "ListCodec.DecodeValue", Types: []reflect.Type{ProtoListValueType}, Received: val}
	}

	ar, err := vr.ReadArray()
	if err != nil {
		return err
	}

	var elems []reflect.Value
	for {
		vr, err := ar.ReadValue()
		if err == bsonrw.ErrEOA {
			break
		}
		if err != nil {
			return err
		}

		elem := reflect.New(ProtoValueType)

		err = DefaultValueCodec.DecodeValue(dc, vr, elem.Elem())
		if err != nil {
			return err
		}
		elems = append(elems, elem)
	}

	valuesField := val.Field(0) // the 'Values' field
	if valuesField.IsNil() {
		valuesField.Set(reflect.MakeSlice(valuesField.Type(), 0, len(elems)))
	}

	valuesField.SetLen(0)
	valuesField.Set(reflect.Append(valuesField, elems...))
	return nil
}
