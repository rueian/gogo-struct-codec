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
)

var (
	_ jsonpb.JSONPBMarshaler   = (*Value)(nil)
	_ jsonpb.JSONPBUnmarshaler = (*Value)(nil)
	_ json.Marshaler           = (*Value)(nil)
	_ json.Unmarshaler         = (*Value)(nil)
	_ driver.Valuer            = (*Value)(nil)
	_ sql.Scanner              = (*Value)(nil)
	_ proto.Message            = (*Value)(nil)
)

type Value struct {
	Val types.Value
}

func (v *Value) XXX_WellKnownType() string    { return v.Val.XXX_WellKnownType() }
func (v *Value) XXX_Unmarshal(b []byte) error { return v.Val.XXX_Unmarshal(b) }
func (v *Value) XXX_Marshal(b []byte, d bool) ([]byte, error) {
	return v.Val.XXX_Marshal(b, d)
}
func (v *Value) XXX_Merge(src proto.Message)      { v.Val.XXX_Merge(src) }
func (v *Value) XXX_Size() int                    { return v.Val.XXX_Size() }
func (v *Value) XXX_DiscardUnknown()              { v.Val.XXX_DiscardUnknown() }
func (v *Value) XXX_OneofWrappers() []interface{} { return v.Val.XXX_OneofWrappers() }
func (v *Value) XXX_MessageName() string          { return v.Val.XXX_MessageName() }

func (v *Value) Equal(that interface{}) bool {
	if that == nil {
		return v == nil
	}
	that1, ok := that.(*Value)
	if !ok {
		that2, ok := that.(Value)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return v == nil
	} else if v == nil {
		return false
	}
	return v.Val.Equal(that1.Val)
}

func (v *Value) Compare(that interface{}) int {
	if that == nil {
		if v == nil {
			return 0
		}
		return 1
	}

	that1, ok := that.(*Value)
	if !ok {
		that2, ok := that.(Value)
		if ok {
			that1 = &that2
		} else {
			return 1
		}
	}
	if that1 == nil {
		if v == nil {
			return 0
		}
		return 1
	} else if v == nil {
		return -1
	}
	return v.Val.Compare(that1.Val)
}

func (v *Value) Reset() {
	*v = Value{}
}

func (v *Value) String() string {
	return v.Val.String()
}

func (v *Value) Marshal() (dAtA []byte, err error) {
	return v.Val.Marshal()
}

func (v *Value) Unmarshal(dAtA []byte) error {
	return v.Val.Unmarshal(dAtA)
}
func (v *Value) Descriptor() ([]byte, []int) {
	return v.Val.Descriptor()
}

func (v *Value) GetNullValue() types.NullValue {
	return v.Val.GetNullValue()
}

func (v *Value) GetNumberValue() float64 {
	return v.Val.GetNumberValue()
}

func (v *Value) GetStringValue() string {
	return v.Val.GetStringValue()
}

func (v *Value) GetBoolValue() bool {
	return v.Val.GetBoolValue()
}

func (v *Value) GetStructValue() *types.Struct {
	return v.Val.GetStructValue()
}

func (v *Value) GetListValue() *types.ListValue {
	return v.Val.GetListValue()
}
func (v *Value) GoString() string {
	return v.Val.GoString()
}

func (v *Value) MarshalTo(dAtA []byte) (int, error) {
	return v.Val.MarshalTo(dAtA)
}

func (v *Value) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return v.Val.MarshalToSizedBuffer(dAtA)
}

func (v *Value) ProtoMessage() {
	v.Val.ProtoMessage()
}

func (v *Value) Size() (n int) {
	return v.Val.Size()
}

func (v *Value) Scan(data interface{}) error {
	switch data := data.(type) {
	case []byte:
		return v.UnmarshalJSON(data)
	case string:
		return v.UnmarshalJSON([]byte(data))
	case nil:
		return errors.New("types.Value: Scan on nil pointer")
	}
	return errors.New(fmt.Sprintf("types.Value: cannot convert %T", data))
}

func (v Value) Value() (driver.Value, error) {
	if bs, err := v.MarshalJSON(); err != nil {
		return nil, err
	} else {
		return string(bs), nil
	}
}

func (v *Value) UnmarshalJSON(bs []byte) error {
	return v.UnmarshalJSONPB(defaultJSONPBUnmarshaler, bs)
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return v.MarshalJSONPB(defaultJSONPBMarshaler)
}

func (v *Value) UnmarshalJSONPB(m *jsonpb.Unmarshaler, bs []byte) error {
	return m.Unmarshal(bytes.NewReader(bs), &v.Val)
}

func (v *Value) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := m.Marshal(buf, &v.Val); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
