package streamx

import (
	"encoding/json"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/zengineDev/x/configx"
	"sync"
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

type BroadcastableContract interface {
	ToCloudEvent() event.Event
	OnSubject() string
}

func Broadcast(e BroadcastableContract) error {
	n := Connection()

	bytes, err := json.Marshal(e.ToCloudEvent())
	if err != nil {
		return err
	}

	natsErr := n.Con.Publish(e.OnSubject(), bytes)
	if natsErr != nil {
		log.Error(natsErr)
	}

	return nil
}
