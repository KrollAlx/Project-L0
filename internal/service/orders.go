package service

import (
	"project-L0/internal/models"
	"project-L0/internal/repository"
)

type Orders interface {
	RestoreCache() error
	Create(order *models.Order) error
	Get(id int) models.Order
}

type OrdersService struct {
	repo        repository.Orders
	ordersCache []models.Order
}

func New(repo repository.Orders) *OrdersService {
	return &OrdersService{repo: repo}
}

func (s *OrdersService) RestoreCache() error {
	orders, err := s.repo.GetAll()
	s.ordersCache = orders
	return err
}

func (s *OrdersService) Create(order *models.Order) error {
	s.ordersCache = append(s.ordersCache, *order)
	return s.repo.Create(order)
}

// TODO: добавить ошибку ненайденного заказа
func (s *OrdersService) Get(id int) models.Order {
	for _, ord := range s.ordersCache {
		if ord.Id == id {
			return ord
		}
	}
	return models.Order{}
}
