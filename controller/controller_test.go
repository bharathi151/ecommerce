package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"quickcommerce/controller"
	"quickcommerce/core"
	"quickcommerce/models"
	"quickcommerce/repo"
	"quickcommerce/service"
)

func newController() *controller.Controller {
	store := repo.NewRepo()
	assigner := core.NewCore(store)
	orderService := service.NewService(store, assigner)
	return controller.NewController(orderService)
}

func request(t *testing.T, handler http.HandlerFunc, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(method, path, bytes.NewBufferString(body)))
	return recorder
}

func TestHandlersValidateRequiredFields(t *testing.T) {
	ctrl := newController()
	tests := []struct {
		name    string
		handler http.HandlerFunc
		body    string
	}{
		{"customer", ctrl.SignupCustomer, `{"id":" ","name":"Ada"}`},
		{"delivery partner", ctrl.SignupDeliveryPartner, `{"id":"partner-1"}`},
		{"order", ctrl.CreateOrder, `{"customerId":"customer-1"}`},
		{"pickup", ctrl.PickUp, `{"deliveryPartnerId":"partner-1"}`},
		{"delivery", ctrl.Deliver, `{"deliveryPartnerId":"partner-1"}`},
		{"cancellation", ctrl.Cancel, `{}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := request(t, test.handler, http.MethodPost, "/", test.body)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestOrderFlowUsesCamelCaseJSONAndTypedStatuses(t *testing.T) {
	ctrl := newController()

	customerResponse := request(t, ctrl.SignupCustomer, http.MethodPost, "/customers", `{"id":"customer-1","name":"Ada"}`)
	if customerResponse.Code != http.StatusCreated {
		t.Fatalf("customer signup status = %d, want %d", customerResponse.Code, http.StatusCreated)
	}

	partnerResponse := request(t, ctrl.SignupDeliveryPartner, http.MethodPost, "/delivery-partners", `{"id":"partner-1","name":"Grace"}`)
	if partnerResponse.Code != http.StatusCreated {
		t.Fatalf("partner signup status = %d, want %d", partnerResponse.Code, http.StatusCreated)
	}
	var partner models.DeliveryPartner
	if err := json.Unmarshal(partnerResponse.Body.Bytes(), &partner); err != nil {
		t.Fatalf("decode partner response: %v", err)
	}
	if partner.Status != models.PartnerStatusAvailable {
		t.Fatalf("partner status = %q, want %q", partner.Status, models.PartnerStatusAvailable)
	}

	orderResponse := request(t, ctrl.CreateOrder, http.MethodPost, "/orders", `{"customerId":"customer-1","item":"book"}`)
	if orderResponse.Code != http.StatusCreated {
		t.Fatalf("order creation status = %d, want %d", orderResponse.Code, http.StatusCreated)
	}
	var order models.Order
	if err := json.Unmarshal(orderResponse.Body.Bytes(), &order); err != nil {
		t.Fatalf("decode order response: %v", err)
	}
	if order.Status != models.OrderStatusAssigned || order.DeliveryPartnerID != "partner-1" {
		t.Fatalf("created order = %+v, want assigned to partner-1", order)
	}

	missingPartnerResponse := request(t, ctrl.PickUp, http.MethodPost, "/orders/pickup",
		`{"deliveryPartnerId":"unknown","orderId":"`+order.ID+`"}`)
	if missingPartnerResponse.Code != http.StatusNotFound {
		t.Fatalf("pickup with unknown partner status = %d, want %d", missingPartnerResponse.Code, http.StatusNotFound)
	}

	pickupResponse := request(t, ctrl.PickUp, http.MethodPost, "/orders/pickup",
		`{"deliveryPartnerId":"partner-1","orderId":"`+order.ID+`"}`)
	if pickupResponse.Code != http.StatusOK {
		t.Fatalf("pickup status = %d, want %d", pickupResponse.Code, http.StatusOK)
	}
	if err := json.Unmarshal(pickupResponse.Body.Bytes(), &order); err != nil {
		t.Fatalf("decode pickup response: %v", err)
	}
	if order.Status != models.OrderStatusPickedUp {
		t.Fatalf("picked-up order status = %q, want %q", order.Status, models.OrderStatusPickedUp)
	}

	deliverResponse := request(t, ctrl.Deliver, http.MethodPost, "/orders/deliver",
		`{"deliveryPartnerId":"partner-1","orderId":"`+order.ID+`"}`)
	if deliverResponse.Code != http.StatusOK {
		t.Fatalf("delivery status = %d, want %d", deliverResponse.Code, http.StatusOK)
	}
	if err := json.Unmarshal(deliverResponse.Body.Bytes(), &order); err != nil {
		t.Fatalf("decode delivery response: %v", err)
	}
	if order.Status != models.OrderStatusDelivered {
		t.Fatalf("delivered order status = %q, want %q", order.Status, models.OrderStatusDelivered)
	}
}
