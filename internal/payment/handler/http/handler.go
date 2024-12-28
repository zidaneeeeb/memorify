package http

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/payment"
	"net/http"

	httplib "hbdtoyou/pkg/http"

	"github.com/gorilla/mux"
)

type paymentsHandler struct {
	payment       payment.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *paymentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreatePayment(w, r)
	case http.MethodGet:
		h.handleGetPayments(w, r)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

type paymentHandler struct {
	payment       payment.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *paymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	paymentID := vars["id"]

	switch r.Method {
	case http.MethodGet:
		h.handleGetPaymentByID(w, r, paymentID)
	case http.MethodPatch:
		h.handleUpdatePayment(w, r, paymentID)
	// case http.MethodDelete:
	// 	h.handleDeletePaymentByID(w, r, paymentID)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}
