package persistentx

import (
	"context"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"

	"github.com/zengineDev/x/configx"
)

var once sync.Once

type DriverPg struct {
	Con *pgx.Conn
}

var (
	instance *DriverPg
)

func Connection() *DriverPg {
	once.Do(func() {
		cfg := configx.GetConfig()
		conn, err := pgx.Connect(context.Background(), cfg.DB.ConnectionString())
		if err != nil {
			log.Error(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		//defer conn.Close(context.Background())

		instance = &DriverPg{Con: conn}
	})

	return instance
}
