package server

import (
	"net/http"
	"quickcommerce/controller"
)

func New(ctrl *controller.Controller) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /customers", ctrl.SignupCustomer)
	mux.HandleFunc("POST /delivery-partners", ctrl.SignupDeliveryPartner)
	mux.HandleFunc("POST /orders", ctrl.CreateOrder)
	mux.HandleFunc("POST /orders/pickup", ctrl.PickUp)
	mux.HandleFunc("POST /orders/deliver", ctrl.Deliver)
	mux.HandleFunc("POST /orders/cancel", ctrl.Cancel)
	mux.HandleFunc("GET /orders/{id}", ctrl.GetOrder)
	mux.HandleFunc("GET /delivery-partners/{id}", ctrl.GetDeliveryPartner)
	return mux
}
