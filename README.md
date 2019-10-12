# gogo-struct-codec

codec of the Well Known `google.protobuf.Struct` type which generated by [gogo protobuf](https://github.com/gogo/protobuf)

## BSON

Convert the gogo's `types.Struct` from and to [bson](https://github.com/mongodb/mongo-go-driver)

### MongoDB example

```go
package main

import (
    "context"
    "github.com/gogo/protobuf/types"
    "github.com/rueian/gogo-struct-codec/structbson"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    
    // MUST first replace bson's DefaultRegistry with the one in structbson
    bson.DefaultRegistry = structbson.DefaultRegistry

    s := types.Struct{Fields:map[string]*types.Value{
        "_id": {Kind: &types.Value_StringValue{StringValue: "str"}},
    }}
    
    // connect mongoclient

    res, err := mongoclient.Database("db").Collection("collection").InsertOne(ctx, s)
    
    // handle err and response 
}
```

### BSON manipulation

[more details in test cases](./structbson/main_test.go)