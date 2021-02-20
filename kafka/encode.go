package kafka

import "google.golang.org/protobuf/proto"

// Marshal is same as proto.Marshal(pb)
func Marshal(pb proto.Message) ([]byte, error) {
	return proto.MarshalOptions{
		AllowPartial: true, // syntax = "proto3";
	}.Marshal(pb)
}

// Unmarshal is same as proto.Unmarshal(buf,pb)
func Unmarshal(buf []byte, pb proto.Message) error {
	return proto.UnmarshalOptions{
		AllowPartial: true, // syntax = "proto3";
	}.Unmarshal(buf, pb)
}

