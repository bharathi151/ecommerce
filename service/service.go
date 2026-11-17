package service

import (
	"quickcommerce/core"
	"quickcommerce/models"
	"quickcommerce/repo"
)

type Service struct {
	repository *repo.Repo
	assigner   *core.Core
}

func NewService(repository *repo.Repo, assigner *core.Core) *Service {
	return &Service{repository: repository, assigner: assigner}
}

func (s *Service) SignupCustomer(id, name string) (models.Customer, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	if _, ok := s.repository.Customers[id]; ok {
		return models.Customer{}, models.ErrExists
	}
	customer := models.Customer{ID: id, Name: name}
	s.repository.Customers[id] = customer
	return customer, nil
}

func (s *Service) SignupDeliveryPartner(id, name string) (models.DeliveryPartner, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	if _, ok := s.repository.DeliveryPartners[id]; ok {
		return models.DeliveryPartner{}, models.ErrExists
	}
	deliveryPartner := &models.DeliveryPartner{ID: id, Name: name, Status: models.PartnerStatusAvailable}
	s.repository.DeliveryPartners[id] = deliveryPartner
	s.repository.PartnerQueue = append(s.repository.PartnerQueue, id)
	s.assigner.Assign()
	return *deliveryPartner, nil
}

func (s *Service) CreateOrder(customerID, item string) (models.Order, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	if _, ok := s.repository.Customers[customerID]; !ok {
		return models.Order{}, models.ErrNotFound
	}
	order := &models.Order{
		ID:         s.repository.NextID(),
		CustomerID: customerID,
		Item:       item,
		Status:     models.OrderStatusCreated,
	}
	s.repository.Orders[order.ID] = order
	s.repository.OrderQueue = append(s.repository.OrderQueue, order.ID)
	s.assigner.Assign()
	return *order, nil
}

func (s *Service) PickUp(deliveryPartnerID, orderID string) (models.Order, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	order, orderExists := s.repository.Orders[orderID]
	_, partnerExists := s.repository.DeliveryPartners[deliveryPartnerID]
	if !orderExists || !partnerExists {
		return models.Order{}, models.ErrNotFound
	}
	if order.Status != models.OrderStatusAssigned || order.DeliveryPartnerID != deliveryPartnerID {
		return models.Order{}, models.ErrState
	}
	order.Status = models.OrderStatusPickedUp
	return *order, nil
}

func (s *Service) Deliver(deliveryPartnerID, orderID string) (models.Order, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	order, orderExists := s.repository.Orders[orderID]
	partner, partnerExists := s.repository.DeliveryPartners[deliveryPartnerID]
	if !orderExists || !partnerExists {
		return models.Order{}, models.ErrNotFound
	}
	if order.Status != models.OrderStatusPickedUp || order.DeliveryPartnerID != deliveryPartnerID {
		return models.Order{}, models.ErrState
	}
	order.Status = models.OrderStatusDelivered
	partner.Status = models.PartnerStatusAvailable
	partner.OrderID = ""
	s.repository.PartnerQueue = append(s.repository.PartnerQueue, deliveryPartnerID)
	s.assigner.Assign()
	return *order, nil
}

func (s *Service) Cancel(orderID string) (models.Order, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	order, ok := s.repository.Orders[orderID]
	if !ok {
		return models.Order{}, models.ErrNotFound
	}
	if order.Status == models.OrderStatusPickedUp ||
		order.Status == models.OrderStatusDelivered ||
		order.Status == models.OrderStatusCancelled {
		return models.Order{}, models.ErrCannotCancel
	}
	partnerID := order.DeliveryPartnerID
	order.Status = models.OrderStatusCancelled
	order.DeliveryPartnerID = ""
	if partner, ok := s.repository.DeliveryPartners[partnerID]; ok && partner.OrderID == orderID {
		partner.Status = models.PartnerStatusAvailable
		partner.OrderID = ""
		s.repository.PartnerQueue = append(s.repository.PartnerQueue, partnerID)
		s.assigner.Assign()
	}
	return *order, nil
}

func (s *Service) GetOrder(id string) (models.Order, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	order, ok := s.repository.Orders[id]
	if !ok {
		return models.Order{}, models.ErrNotFound
	}
	return *order, nil
}

func (s *Service) GetDeliveryPartner(id string) (models.DeliveryPartner, error) {
	s.repository.Lock()
	defer s.repository.Unlock()
	partner, ok := s.repository.DeliveryPartners[id]
	if !ok {
		return models.DeliveryPartner{}, models.ErrNotFound
	}
	return *partner, nil
}
