package server

import (
	"net/http"
	"project-L0/internal/transport/http/handler"
)

type Server struct {
	h *handler.OrdersHandler
}

func NewHttpServer(h *handler.OrdersHandler) *Server {
	return &Server{h: h}
}

func (s *Server) Run() error {
	http.HandleFunc("/orders/", s.h.GetOrder)
	return http.ListenAndServe(":3000", nil)
}
