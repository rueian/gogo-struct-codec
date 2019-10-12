package structbson

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_Convert(t *testing.T) {
	bson.DefaultRegistry = DefaultRegistry

	fixture := &types.Struct{Fields: map[string]*types.Value{
		"nullValue": nil,
		"sliceValue": {
			Kind: &types.Value_ListValue{
				ListValue: &types.ListValue{
					Values: []*types.Value{
						{
							Kind: &types.Value_ListValue{},
						},
						{
							Kind: &types.Value_ListValue{ListValue: nil},
						},
						{
							Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: []*types.Value{}}},
						},
						{
							Kind: &types.Value_ListValue{
								ListValue: &types.ListValue{
									Values: []*types.Value{
										{
											Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: nil}},
										},
										{
											Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: map[string]*types.Value{}}},
										},
										{
											Kind: &types.Value_StructValue{
												StructValue: &types.Struct{
													Fields: map[string]*types.Value{
														"nullValue": nil,
														"zeroValue": {
															Kind: &types.Value_StructValue{},
														},
														"emptyValue": {
															Kind: &types.Value_StructValue{StructValue: &types.Struct{}},
														},
														"structValue": {Kind: &types.Value_StructValue{
															StructValue: &types.Struct{
																Fields: map[string]*types.Value{
																	"string":  {Kind: &types.Value_StringValue{StringValue: "str"}},
																	"number":  {Kind: &types.Value_NumberValue{NumberValue: 1234}},
																	"boolean": {Kind: &types.Value_BoolValue{BoolValue: true}},
																	"null":    {Kind: &types.Value_NullValue{NullValue: 0}},
																}},
														}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}}

	bs, err := bson.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}

	target := &types.Struct{}
	if err := bson.Unmarshal(bs, target); err != nil {
		t.Fatal(err)
	}

	jsonMarshaler := jsonpb.Marshaler{}

	jb, err := jsonMarshaler.MarshalToString(fixture)
	if err != nil {
		t.Fatal(err)
	}

	jt, err := jsonMarshaler.MarshalToString(target)
	if err != nil {
		t.Fatal(err)
	}

	if jt != jb {
		t.Fail()
	}
}
