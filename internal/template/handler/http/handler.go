package http

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/template"
	"net/http"

	httplib "hbdtoyou/pkg/http"

	"github.com/gorilla/mux"
)

type templatesHandler struct {
	template      template.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *templatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateTemplate(w, r)
	case http.MethodGet:
		h.handleGetTemplates(w, r)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

type templateHandler struct {
	template      template.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	templateID := vars["id"]

	switch r.Method {
	case http.MethodGet:
		h.handleGetTemplateByID(w, r, templateID)
	case http.MethodPatch:
		h.handleUpdateTemplate(w, r, templateID)
	case http.MethodDelete:
		h.handleDeleteTemplateByID(w, r, templateID)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}
