package persistentx

import (
	"context"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
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
			log.Error(err)
		}
		//defer conn.Close(context.Background())

		instance = &DriverPg{Con: conn}
	})

	return instance
}
