package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
    "bytes"
)


// ProtoToJSON converts a proto2 message to JSON bytes
func ProtoToJSON(pb proto.Message) ([]byte, error) {
	marshaler := &jsonpb.Marshaler{
		EmitDefaults: true, // Include fields even if they're at default values
		OrigName:     true, // Use camelCase naming in JSON
		Indent:       "  ",
	}
	json, err := marshaler.MarshalToString(pb)
	return []byte(json), err
}

// JSONToProto converts JSON bytes to a proto2 message
func JSONToProto(data []byte, pb proto.Message) error {
	unmarshaler := &jsonpb.Unmarshaler{
		AllowUnknownFields: false, // More forgiving JSON parsing
	}
	return unmarshaler.Unmarshal(bytes.NewReader(data), pb)
}
