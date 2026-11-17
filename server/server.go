package server

import (
	"net/http"
	"quickcommerce/controller"
)

func New(controller *controller.Controller) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /customers", controller.SignupCustomer)
	mux.HandleFunc("POST /delivery-partners", controller.SignupDeliveryPartner)
	mux.HandleFunc("POST /orders", controller.CreateOrder)
	mux.HandleFunc("POST /orders/pickup", controller.PickUp)
	mux.HandleFunc("POST /orders/deliver", controller.Deliver)
	mux.HandleFunc("POST /orders/cancel", controller.Cancel)
	mux.HandleFunc("GET /orders/{id}", controller.GetOrder)
	mux.HandleFunc("GET /delivery-partners/{id}", controller.GetDeliveryPartner)
	return mux
}
