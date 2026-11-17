package service

import (
	"quickcommerce/core"
	"quickcommerce/models"
	"quickcommerce/repo"
)

type Service struct {
	r *repo.Repo
	c *core.Core
}

func NewService(r *repo.Repo, c *core.Core) *Service {
	return &Service{r: r, c: c}
}

func (s *Service) SignupCustomer(id, name string) (models.Customer, error) {
	s.r.Lock()
	defer s.r.Unlock()
	if _, ok := s.r.Customers[id]; ok {
		return models.Customer{}, models.ErrExists
	}
	cust := models.Customer{ID: id, Name: name}
	s.r.Customers[id] = cust
	return cust, nil
}

func (s *Service) SignupDeliveryPartner(id, name string) (models.DeliveryPartner, error) {
	s.r.Lock()
	defer s.r.Unlock()
	if _, ok := s.r.DeliveryPartners[id]; ok {
		return models.DeliveryPartner{}, models.ErrExists
	}
	dp := &models.DeliveryPartner{ID: id, Name: name}
	s.r.DeliveryPartners[id] = dp
	s.r.PartnerQueue = append(s.r.PartnerQueue, id)
	s.c.Assign()
	return *dp, nil
}

func (s *Service) CreateOrder(customerID, item string) (models.Order, error) {
	s.r.Lock()
	defer s.r.Unlock()
	if _, ok := s.r.Customers[customerID]; !ok {
		return models.Order{}, models.ErrNotFound
	}
	o := &models.Order{ID: s.r.NextID(), CustomerID: customerID, Item: item, Status: "CREATED"}
	s.r.Orders[o.ID] = o
	s.r.OrderQueue = append(s.r.OrderQueue, o.ID)
	s.c.Assign()
	return *o, nil
}

func (s *Service) PickUp(dpID, orderID string) (models.Order, error) {
	s.r.Lock()
	defer s.r.Unlock()
	o, ok := s.r.Orders[orderID]
	if !ok {
		return models.Order{}, models.ErrNotFound
	}
	if o.Status != "ASSIGNED" || o.DeliveryPartnerID != dpID {
		return models.Order{}, models.ErrState
	}
	o.Status = "PICKED_UP"
	return *o, nil
}

func (s *Service) Deliver(dpID, orderID string) (models.Order, error) {
	s.r.Lock()
	defer s.r.Unlock()
	o, ok := s.r.Orders[orderID]
	dp, dpOk := s.r.DeliveryPartners[dpID]
	if !ok || !dpOk {
		return models.Order{}, models.ErrNotFound
	}
	if o.Status != "PICKED_UP" || o.DeliveryPartnerID != dpID {
		return models.Order{}, models.ErrState
	}
	o.Status = "DELIVERED"
	dp.Busy = false
	dp.OrderID = ""
	s.r.PartnerQueue = append(s.r.PartnerQueue, dpID)
	s.c.Assign()
	return *o, nil
}

func (s *Service) Cancel(orderID string) (models.Order, error) {
	s.r.Lock()
	defer s.r.Unlock()
	o, ok := s.r.Orders[orderID]
	if !ok {
		return models.Order{}, models.ErrNotFound
	}
	if o.Status == "PICKED_UP" || o.Status == "DELIVERED" || o.Status == "CANCELLED" {
		return models.Order{}, models.ErrState
	}
	dpID := o.DeliveryPartnerID
	o.Status = "CANCELLED"
	o.DeliveryPartnerID = ""
	if dp, ok := s.r.DeliveryPartners[dpID]; ok && dp.OrderID == orderID {
		dp.Busy = false
		dp.OrderID = ""
		s.r.PartnerQueue = append(s.r.PartnerQueue, dpID)
		s.c.Assign()
	}
	return *o, nil
}

func (s *Service) GetOrder(id string) (models.Order, error) {
	s.r.Lock()
	defer s.r.Unlock()
	o, ok := s.r.Orders[id]
	if !ok {
		return models.Order{}, models.ErrNotFound
	}
	return *o, nil
}

func (s *Service) GetDeliveryPartner(id string) (models.DeliveryPartner, error) {
	s.r.Lock()
	defer s.r.Unlock()
	dp, ok := s.r.DeliveryPartners[id]
	if !ok {
		return models.DeliveryPartner{}, models.ErrNotFound
	}
	return *dp, nil
}
