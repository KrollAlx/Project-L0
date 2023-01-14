package server

import (
	"github.com/nats-io/stan.go"
	"project-L0/internal/transport/nats"
)

const (
	clusterID   = "test-cluster"
	clientID    = "client-subscriber"
	natsUrl     = "0.0.0.0:4223"
	channelName = "orders"
)

type NatsServer struct {
	h   *nats.OrdersHandler
	sc  stan.Conn
	sub stan.Subscription
}

func NewNatsServer(h *nats.OrdersHandler) *NatsServer {
	return &NatsServer{h: h}
}

func (s *NatsServer) Run() error {
	var err error
	s.sc, err = stan.Connect(clusterID, clientID, stan.NatsURL(natsUrl))
	if err != nil {
		return err
	}

	s.sub, err = s.sc.Subscribe(channelName, s.h.Create)
	if err != nil {
		return err
	}

	return nil
}

func (s *NatsServer) Shutdown() error {
	err := s.sub.Unsubscribe()
	if err != nil {
		return err
	}
	err = s.sc.Close()
	if err != nil {
		return err
	}
	return nil
}
