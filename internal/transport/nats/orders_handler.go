package nats

import (
	"encoding/json"
	"errors"
	"github.com/nats-io/stan.go"
	"log"
	"project-L0/internal/models"
	"project-L0/internal/service"
)

type OrdersHandler struct {
	service service.Orders
}

func New(service service.Orders) *OrdersHandler {
	return &OrdersHandler{service: service}
}

func (h *OrdersHandler) Create(m *stan.Msg) {
	log.Println("Message received")

	var order models.Order

	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		log.Println(errors.New("invalid input structure"))
		return
	}

	err = h.service.Create(&order)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Order created")
}
