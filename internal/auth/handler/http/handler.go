package http

import (
	"hbdtoyou/internal/auth"
	httplib "hbdtoyou/pkg/http"
	"net/http"

	"github.com/gorilla/mux"
)

type authHandler struct {
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleLoginSocial(w, r)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

type userHandler struct {
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID := vars["id"]

	switch r.Method {
	case http.MethodGet:
		h.handleGetUserByID(w, r, userID)
	case http.MethodPatch:
		h.handleUpdateUser(w, r, userID)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}
