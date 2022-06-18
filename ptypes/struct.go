package ptypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/rueian/gogo-struct-codec/structbson"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	_ jsonpb.JSONPBMarshaler   = (*Struct)(nil)
	_ jsonpb.JSONPBUnmarshaler = (*Struct)(nil)
	_ bson.Marshaler           = (*Struct)(nil)
	_ bson.Unmarshaler         = (*Struct)(nil)
	_ json.Marshaler           = (*Struct)(nil)
	_ json.Unmarshaler         = (*Struct)(nil)
	_ driver.Valuer            = (*Struct)(nil)
	_ sql.Scanner              = (*Struct)(nil)
	_ proto.Message            = (*Struct)(nil)
	_ msgpack.CustomDecoder    = (*Struct)(nil)
	_ msgpack.CustomEncoder    = (*Struct)(nil)
)

var (
	defaultJSONPBMarshaler   = &jsonpb.Marshaler{}
	defaultJSONPBUnmarshaler = &jsonpb.Unmarshaler{}
)

type Struct struct {
	types.Struct
}

func (s *Struct) EncodeMsgpack(encoder *msgpack.Encoder) error {
	raw, err := proto.Marshal(s)
	if err != nil {
		return err
	}

	return encoder.EncodeBytes(raw)
}

func (s *Struct) DecodeMsgpack(decoder *msgpack.Decoder) error {
	raw, err := decoder.DecodeBytes()
	if err != nil {
		return err
	}

	return proto.Unmarshal(raw, s)
}

func (s *Struct) GetFields() map[string]*types.Value {
	if s != nil {
		return s.Struct.GetFields()
	}
	return nil
}

func (s *Struct) Equal(that interface{}) bool {
	if that == nil {
		return s == nil
	}
	that1, ok := that.(*Struct)
	if !ok {
		that2, ok := that.(Struct)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return s == nil
	} else if s == nil {
		return false
	}
	return s.Struct.Equal(that1.Struct)
}

func (s *Struct) Compare(that interface{}) int {
	if that == nil {
		if s == nil {
			return 0
		}
		return 1
	}

	that1, ok := that.(*Struct)
	if !ok {
		that2, ok := that.(Struct)
		if ok {
			that1 = &that2
		} else {
			return 1
		}
	}
	if that1 == nil {
		if s == nil {
			return 0
		}
		return 1
	} else if s == nil {
		return -1
	}
	return s.Struct.Compare(that1.Struct)
}

func (s *Struct) Reset() { *s = Struct{} }

func (s Struct) Value() (driver.Value, error) {
	if bs, err := s.MarshalJSON(); err != nil {
		return nil, err
	} else {
		return string(bs), nil
	}
}

func (s *Struct) Scan(data interface{}) error {
	switch data := data.(type) {
	case []byte:
		return s.UnmarshalJSON(data)
	case string:
		return s.UnmarshalJSON([]byte(data))
	case nil:
		return errors.New("types.Struct: Scan on nil pointer")
	}
	return errors.New(fmt.Sprintf("types.Struct: cannot convert %T", data))
}

func (s *Struct) UnmarshalJSONPB(m *jsonpb.Unmarshaler, bs []byte) error {
	return m.Unmarshal(bytes.NewReader(bs), &s.Struct)
}

func (s Struct) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := m.Marshal(buf, &s.Struct); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Struct) UnmarshalBSON(bs []byte) error {
	return bson.UnmarshalWithRegistry(structbson.DefaultRegistry, bs, &s.Struct)
}

func (s Struct) MarshalBSON() ([]byte, error) {
	return bson.MarshalWithRegistry(structbson.DefaultRegistry, s.Struct)
}

func (s *Struct) UnmarshalJSON(bs []byte) error {
	return s.UnmarshalJSONPB(defaultJSONPBUnmarshaler, bs)
}

func (s Struct) MarshalJSON() ([]byte, error) {
	return s.MarshalJSONPB(defaultJSONPBMarshaler)
}
