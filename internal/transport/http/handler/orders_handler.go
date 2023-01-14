package handler

import (
	"html/template"
	"log"
	"net/http"
	"project-L0/internal/service"
	"strconv"
	"strings"
)

const (
	orderTemplatePath = "web/template/order.html"
	ordersRoute       = "/orders/"
)

type OrdersHandler struct {
	service service.Orders
}

func New(service service.Orders) *OrdersHandler {
	return &OrdersHandler{service: service}
}

func (h *OrdersHandler) GetOrder(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, ordersRoute)
	orderId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}

	order := h.service.Get(orderId)

	tmpl, _ := template.ParseFiles(orderTemplatePath)
	err = tmpl.Execute(w, order)
	if err != nil {
		log.Println(err)
	}
}
