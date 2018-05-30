package common

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pborman/uuid"
	pb "github.com/ZTP/pnp/common/proto"
)

// NewReqHdrGenerateTraceAndMessageID will generate a common request header with the TraceID, MessageID generated.
// It will also set the request time to be the moment this is called.
func NewReqHdrGenerateTraceAndMessageID() *pb.RequestHeader {
	return &pb.RequestHeader{Identifiers: &pb.Identifiers{TraceID: NewTraceID(), MessageID: NewMessageID()}, RequestTimestamp: ptypes.TimestampNow()}
}

// NewReqHdrGenerateMessageID will generate a common request header with the MessageID generated. The TraceID is provided by the caller as
// it may be part of a wider transaction where the TraceID is already known.
// It will also set the request time to be the moment this is called.
func NewReqHdrGenerateMessageID(traceID string) *pb.RequestHeader {
	return &pb.RequestHeader{Identifiers: &pb.Identifiers{TraceID: traceID, MessageID: NewMessageID()}, RequestTimestamp: ptypes.TimestampNow()}
}

// NewTraceID returns a new Trace ID.
func NewTraceID() string {
	return uuid.New()
}

// NewMessageID returns a new Message ID.
func NewMessageID() string {
	return uuid.New()
}
