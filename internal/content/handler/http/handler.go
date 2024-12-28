package http

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/content"
	"net/http"

	httplib "hbdtoyou/pkg/http"

	"github.com/gorilla/mux"
)

type contentsHandler struct {
	content       content.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *contentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateContent(w, r)
	case http.MethodGet:
		h.handleGetContents(w, r)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

type contentHandler struct {
	content       content.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

func (h *contentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contentID := vars["id"]

	switch r.Method {
	case http.MethodGet:
		h.handleGetContentByID(w, r, contentID)
	case http.MethodPatch:
		h.handleUpdateContent(w, r, contentID)
	case http.MethodDelete:
		h.handleDeleteContentByID(w, r, contentID)
	default:
		httplib.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}
