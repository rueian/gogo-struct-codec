package structbson

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

var DefaultStructCodec = StructCodec{}

type StructCodec struct{}

func (c StructCodec) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoStructType {
		return bsoncodec.ValueEncoderError{Name: "StructCodec.EncodeValue", Types: []reflect.Type{ProtoStructType}, Received: val}
	}

	dw, err := vw.WriteDocument()
	if err != nil {
		return err
	}

	fieldsField := val.Field(0) // the 'Fields' field
	keys := fieldsField.MapKeys()

	var evw bsonrw.ValueWriter
	for _, key := range keys {
		evw, err = dw.WriteDocumentElement(key.String())
		if err != nil {
			return err
		}

		vv := fieldsField.MapIndex(key)
		if vv.IsNil() {
			err = evw.WriteNull()
		} else {
			err = DefaultValueCodec.EncodeValue(ec, evw, vv.Elem())
		}
		if err != nil {
			return err
		}
	}
	return dw.WriteDocumentEnd()
}

func (c StructCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || val.Type() != ProtoStructType {
		return bsoncodec.ValueDecoderError{Name: "StructCodec.DecodeValue", Types: []reflect.Type{ProtoStructType}, Received: val}
	}

	fieldsField := val.Field(0) // the 'Fields' field

	if fieldsField.IsNil() {
		fieldsField.Set(reflect.MakeMap(fieldsField.Type()))
	}

	mapDecoder, err := dc.LookupDecoder(fieldsField.Type())
	if err != nil {
		return err
	}

	return mapDecoder.DecodeValue(dc, vr, fieldsField)
}
