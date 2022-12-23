package grpc

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func TimeToTimestamp(t time.Time) *timestamp.Timestamp {
	if t.IsZero() {
		return nil
	}

	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   0,
	}
}
