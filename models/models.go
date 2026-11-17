package models

import "errors"

var (
	ErrExists   = errors.New("already exists")
	ErrNotFound = errors.New("not found")
	ErrState    = errors.New("invalid state")
)

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeliveryPartner struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Busy    bool   `json:"busy"`
	OrderID string `json:"orderId,omitempty"`
}

type Order struct {
	ID                string `json:"id"`
	CustomerID        string `json:"customerId"`
	Item              string `json:"item"`
	Status            string `json:"status"`
	DeliveryPartnerID string `json:"deliveryPartnerId,omitempty"`
}
