package healthx

import (
	"context"

	persistentx "github.com/zengineDev/x/storage/persistent"
)

type PostgresConnectionCheck struct {
	Conn *persistentx.DriverPg
	ID   string
}

func (c *PostgresConnectionCheck) Run() HealthCheckStatusContract {
	result := HealthCheckStatus{
		StatusText: "connected",
		Failed:     false,
		ID:         c.ID,
	}

	err := c.Conn.Con.Ping(context.Background())
	if err != nil {
		result.Failed = true
		result.StatusText = "disconnected"
	}

	return result
}
