package healthx

import (
	"github.com/nats-io/nats.go"
)

type NatsConnectionCheck struct {
	Conn *nats.Conn
	ID   string
}

func (c *NatsConnectionCheck) Run() HealthCheckStatusContract {
	result := HealthCheckStatus{
		StatusText: "connected",
		Failed:     false,
		ID:         c.ID,
	}
	switch c.Conn.Status() {
	case nats.CONNECTED:
		result.StatusText = "connected"
	case nats.CLOSED:
		result.StatusText = "closed"
		result.Failed = true
	default:
		result.StatusText = "other"
		result.Failed = true
	}

	return result
}
