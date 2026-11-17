package models

import "errors"

var (
	ErrExists       = errors.New("already exists")
	ErrNotFound     = errors.New("not found")
	ErrState        = errors.New("invalid state")
	ErrCannotCancel = errors.New("cannot cancel order")
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusAssigned  OrderStatus = "ASSIGNED"
	OrderStatusPickedUp  OrderStatus = "PICKED_UP"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type PartnerStatus string

const (
	PartnerStatusAvailable PartnerStatus = "AVAILABLE"
	PartnerStatusBusy      PartnerStatus = "BUSY"
)

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeliveryPartner struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Status  PartnerStatus `json:"status"`
	OrderID string        `json:"orderId,omitempty"`
}

type Order struct {
	ID                string      `json:"id"`
	CustomerID        string      `json:"customerId"`
	Item              string      `json:"item"`
	Status            OrderStatus `json:"status"`
	DeliveryPartnerID string      `json:"deliveryPartnerId,omitempty"`
}
