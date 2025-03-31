package sources

import (
	"fmt"
	"net"
	"time"
)

type ReportMessageLevelType string

const (
	ReportMessageLevelTypeInfo    ReportMessageLevelType = "INFO"
	ReportMessageLevelTypeWarning ReportMessageLevelType = "WARNING"
	ReportMessageLevelTypeError   ReportMessageLevelType = "ERROR"
)

type ReportMessage struct {
	IP        net.IP                 `json:"ip_address"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Level     ReportMessageLevelType `json:"level"`
}

type CloudInitEventType string

const (
	CloudInitEventTypeStart CloudInitEventType = "start"
	CloudInitEventTypeEnd   CloudInitEventType = "finish"
)

type CloudInitResultType string

const (
	CloudInitResultTypeSuccess CloudInitResultType = "SUCCESS"
	CloudInitResultTypeWarn    CloudInitResultType = "WARN"
	CloudInitResultTypeFail    CloudInitResultType = "FAIL"
)

type CloudInitReport struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	EventType   CloudInitEventType  `json:"event_type"`
	Origin      string              `json:"origin"`
	Timestamp   float64             `json:"timestamp"`
	Result      CloudInitResultType `json:"result"`
}

func (c *CloudInitReport) ToReportMessage() ReportMessage {
	level := ReportMessageLevelTypeInfo
	switch c.Result {
	case CloudInitResultTypeFail:
		level = ReportMessageLevelTypeError
	case CloudInitResultTypeWarn:
		level = ReportMessageLevelTypeWarning
	case CloudInitResultTypeSuccess:
		break
	}

	return ReportMessage{
		IP: nil,
		Message: fmt.Sprintf("Cloud-Init name:%s description:%s origin:%s result:%s type:%s",
			c.Name, c.Description, c.Origin, c.Result, c.EventType),
		Timestamp: time.Unix(int64(c.Timestamp), 0),
		Level:     level,
	}
}
