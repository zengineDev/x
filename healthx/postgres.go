package healthx

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresConnectionCheck struct {
	Conn *pgxpool.Conn
	ID   string
}

func (c *PostgresConnectionCheck) Run() HealthCheckStatusContract {
	result := HealthCheckStatus{
		StatusText: "connected",
		Failed:     false,
		ID:         c.ID,
	}

	err := c.Conn.Conn().Ping(context.Background())
	if err != nil {
		result.Failed = true
		result.StatusText = "disconnected"
	}

	return result
}
