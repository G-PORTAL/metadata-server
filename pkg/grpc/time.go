package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}

	return &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   0,
	}
}
