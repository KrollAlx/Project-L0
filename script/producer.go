package main

import (
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

const (
	clusterID = "test-cluster"
	clientID  = "client-producer"
	natsUrl   = "0.0.0.0:4223"
)

func main() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsUrl))
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	//jsonFile, err := os.Open("model.json")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println("Successfully Opened model.json")

	content, err := ioutil.ReadFile("model.json")
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(content)

	err = sc.Publish("orders", content)
	if err != nil {
		log.Println(err)
		return
	}
}
