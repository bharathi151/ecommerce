package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"quickcommerce/models"
	"quickcommerce/service"
)

type Controller struct {
	service *service.Service
}

type customerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type deliveryPartnerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type createOrderRequest struct {
	CustomerID string `json:"customerId"`
	Item       string `json:"item"`
}

type orderActionRequest struct {
	DeliveryPartnerID string `json:"deliveryPartnerId"`
	OrderID           string `json:"orderId"`
}

type cancelOrderRequest struct {
	OrderID string `json:"orderId"`
}

func NewController(service *service.Service) *Controller {
	return &Controller{service: service}
}

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func writeErr(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	if errors.Is(err, models.ErrNotFound) {
		code = http.StatusNotFound
	}
	write(w, code, map[string]string{"error": err.Error()})
}

func hasEmptyField(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func writeBadRequest(w http.ResponseWriter, message string) {
	write(w, http.StatusBadRequest, map[string]string{"error": message})
}

func (c *Controller) SignupCustomer(w http.ResponseWriter, r *http.Request) {
	var request customerRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.ID = strings.TrimSpace(request.ID)
	request.Name = strings.TrimSpace(request.Name)
	if hasEmptyField(request.ID, request.Name) {
		writeBadRequest(w, "id and name are required")
		return
	}
	res, err := c.service.SignupCustomer(request.ID, request.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusCreated, res)
}

func (c *Controller) SignupDeliveryPartner(w http.ResponseWriter, r *http.Request) {
	var request deliveryPartnerRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.ID = strings.TrimSpace(request.ID)
	request.Name = strings.TrimSpace(request.Name)
	if hasEmptyField(request.ID, request.Name) {
		writeBadRequest(w, "id and name are required")
		return
	}
	res, err := c.service.SignupDeliveryPartner(request.ID, request.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusCreated, res)
}

func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request createOrderRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.CustomerID = strings.TrimSpace(request.CustomerID)
	request.Item = strings.TrimSpace(request.Item)
	if hasEmptyField(request.CustomerID, request.Item) {
		writeBadRequest(w, "customerId and item are required")
		return
	}
	res, err := c.service.CreateOrder(request.CustomerID, request.Item)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusCreated, res)
}

func (c *Controller) PickUp(w http.ResponseWriter, r *http.Request) {
	var request orderActionRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.DeliveryPartnerID = strings.TrimSpace(request.DeliveryPartnerID)
	request.OrderID = strings.TrimSpace(request.OrderID)
	if hasEmptyField(request.DeliveryPartnerID, request.OrderID) {
		writeBadRequest(w, "deliveryPartnerId and orderId are required")
		return
	}
	res, err := c.service.PickUp(request.DeliveryPartnerID, request.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusOK, res)
}

func (c *Controller) Deliver(w http.ResponseWriter, r *http.Request) {
	var request orderActionRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.DeliveryPartnerID = strings.TrimSpace(request.DeliveryPartnerID)
	request.OrderID = strings.TrimSpace(request.OrderID)
	if hasEmptyField(request.DeliveryPartnerID, request.OrderID) {
		writeBadRequest(w, "deliveryPartnerId and orderId are required")
		return
	}
	res, err := c.service.Deliver(request.DeliveryPartnerID, request.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusOK, res)
}

func (c *Controller) Cancel(w http.ResponseWriter, r *http.Request) {
	var request cancelOrderRequest
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		writeBadRequest(w, "bad json")
		return
	}
	request.OrderID = strings.TrimSpace(request.OrderID)
	if hasEmptyField(request.OrderID) {
		writeBadRequest(w, "orderId is required")
		return
	}
	res, err := c.service.Cancel(request.OrderID)
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusOK, res)
}

func (c *Controller) GetOrder(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetOrder(r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusOK, res)
}

func (c *Controller) GetDeliveryPartner(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetDeliveryPartner(r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	write(w, http.StatusOK, res)
}
