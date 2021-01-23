package streamx

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/zengineDev/x/configx"
)

var once sync.Once

type NatsCon struct {
	Con *nats.Conn
}

var (
	instance *NatsCon
)

func Connection() *NatsCon {
	once.Do(func() {
		cfg := configx.GetConfig()

		con, err := nats.Connect(cfg.Nats.Url)

		if err != nil {
			log.Fatal("Failed to connect to nats: %s", err)
		}

		instance = &NatsCon{Con: con}
	})

	return instance
}

func Broadcast(data interface{}, subject string) error {
	n := Connection()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	natsErr := n.Con.Publish(subject, bytes)
	if natsErr != nil {
		log.Error(natsErr)
	}

	return nil
}
