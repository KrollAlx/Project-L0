package nats

import (
	"encoding/json"
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
	var order models.Order

	log.Printf("Received a message: %s\n", string(m.Data))

	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.service.Create(&order)
	if err != nil {
		log.Println(err)
		return
	}
}
