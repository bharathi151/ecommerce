package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"quickcommerce/models"
	"quickcommerce/service"
)

type Controller struct {
	svc *service.Service
}

func NewController(svc *service.Service) *Controller { return &Controller{svc: svc} }

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	if errors.Is(err, models.ErrNotFound) {
		code = http.StatusNotFound
	}
	write(w, code, map[string]string{"error": err.Error()})
}

func (c *Controller) SignupCustomer(w http.ResponseWriter, r *http.Request) {
	var b struct{ ID, Name string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.SignupCustomer(b.ID, b.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 201, res)
}

func (c *Controller) SignupDeliveryPartner(w http.ResponseWriter, r *http.Request) {
	var b struct{ ID, Name string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.SignupDeliveryPartner(b.ID, b.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 201, res)
}

func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var b struct{ CustomerID, Item string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.CreateOrder(b.CustomerID, b.Item)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 201, res)
}

func (c *Controller) PickUp(w http.ResponseWriter, r *http.Request) {
	var b struct{ DeliveryPartnerID, OrderID string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.PickUp(b.DeliveryPartnerID, b.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 200, res)
}

func (c *Controller) Deliver(w http.ResponseWriter, r *http.Request) {
	var b struct{ DeliveryPartnerID, OrderID string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.Deliver(b.DeliveryPartnerID, b.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 200, res)
}

func (c *Controller) Cancel(w http.ResponseWriter, r *http.Request) {
	var b struct{ OrderID string }
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		write(w, 400, map[string]string{"error": "bad json"})
		return
	}
	res, err := c.svc.Cancel(b.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 200, res)
}

func (c *Controller) GetOrder(w http.ResponseWriter, r *http.Request) {
	res, err := c.svc.GetOrder(r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 200, res)
}

func (c *Controller) GetDeliveryPartner(w http.ResponseWriter, r *http.Request) {
	res, err := c.svc.GetDeliveryPartner(r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, 200, res)
}
