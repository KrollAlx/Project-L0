package pub

import (
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

const (
	clusterID   = "test-cluster"
	clientID    = "client-publisher"
	natsUrl     = "0.0.0.0:4223"
	channelName = "orders"
)

func Run() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsUrl))
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	content, err := ioutil.ReadFile("assets/model.json")
	if err != nil {
		log.Println(err)
		return
	}

	err = sc.Publish(channelName, content)
	if err != nil {
		log.Println(err)
		return
	}
}
