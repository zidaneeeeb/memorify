package http

import (
	"errors"
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/content"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	errUnknownScope  = errors.New("unknown scope name")
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains finance HTTP handlers.
type Handler struct {
	handlers      map[string]*handler
	content       content.Service
	auth          auth.Service
	scopeSettings map[Scope]ScopeSetting
}

// handler is the HTTP handler wrapper.
type handler struct {
	h        http.Handler
	identity HandlerIdentity
}

// HandlerIdentity denotes the identity of an HTTP hanlder.
type HandlerIdentity struct {
	Name string
	URL  string
}

// Followings are the known HTTP handler identities
var (
	HandlerContent = HandlerIdentity{
		Name: "content",
		URL:  "/v1/contents/{id}",
	}
	HandlerContents = HandlerIdentity{
		Name: "contents",
		URL:  "/v1/contents",
	}
)

// Scope is a shared settings identifier.
//
// Registering a new Scope is done by adding a new Scope
// value and new entry in ScopeName and ScopeValue.
type Scope int

// Followings are the known scopes in finance HTTP
// handlers.
const (
	_ Scope = iota
	ScopeCreateContent
	ScopeGetContents
	ScopeGetContentByID
	ScopeUpdateContent
	ScopeDeleteContent
)

var (
	// ScopeName defines all the known scopes and their string
	// representation.
	ScopeName = map[Scope]string{
		ScopeCreateContent:  "CreateContent",
		ScopeGetContents:    "GetContents",
		ScopeGetContentByID: "GetContentByID",
		ScopeUpdateContent:  "UpdateContent",
		ScopeDeleteContent:  "DeleteContent",
	}

	// ScopeValue is the reverse-mapping of ScopeName.
	ScopeValue = map[string]Scope{
		ScopeName[ScopeCreateContent]:  ScopeCreateContent,
		ScopeName[ScopeGetContents]:    ScopeGetContents,
		ScopeName[ScopeGetContentByID]: ScopeGetContentByID,
		ScopeName[ScopeUpdateContent]:  ScopeUpdateContent,
		ScopeName[ScopeDeleteContent]:  ScopeDeleteContent,
	}
)

// ScopeSetting is the available configurations of a Scope.
type ScopeSetting struct {
	Timeout time.Duration
}

// Followings are default values for ScopeSetting fields.
const (
	defaultTimeout = 5000 * time.Millisecond
)

// getDefaultScopeSettings returns default scope settings
// for all scopes.
func getDefaultScopeSettings() map[Scope]ScopeSetting {
	defaultSettings := make(map[Scope]ScopeSetting)
	for _, scope := range ScopeValue {
		defaultSettings[scope] = ScopeSetting{
			Timeout: defaultTimeout,
		}
	}
	return defaultSettings
}

// Option controls the behavior of Handler.
type Option func(*Handler) error

// WithHandler returns Option to add HTTP handler.
func WithHandler(identity HandlerIdentity) Option {
	return Option(func(h *Handler) error {
		if h.handlers == nil {
			h.handlers = map[string]*handler{}
		}

		h.handlers[identity.Name] = &handler{
			identity: identity,
		}

		handler, err := h.createHTTPHandler(identity.Name)
		if err != nil {
			return err
		}

		h.handlers[identity.Name].h = handler
		return nil
	})
}

// WithScopeSetting returns Option to set scope setting for
// a specific scope name.
func WithScopeSetting(scopeName string, scopeSetting ScopeSetting) Option {
	return Option(func(h *Handler) error {
		scope, ok := ScopeValue[scopeName]
		if !ok {
			return errUnknownScope
		}

		// validate setting
		if scopeSetting.Timeout <= 0 {
			scopeSetting.Timeout = defaultTimeout
		}

		h.scopeSettings[scope] = scopeSetting
		return nil
	})
}

// New creates a new Handler.
//
// For the given Option, WithScopeSetting() should come first
// before WithHandler()
func New(content content.Service, auth auth.Service, options ...Option) (*Handler, error) {
	h := &Handler{
		handlers:      make(map[string]*handler),
		content:       content,
		auth:          auth,
		scopeSettings: getDefaultScopeSettings(),
	}

	// apply options
	for _, opt := range options {
		err := opt(h)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

// createHTTPHandler creates a new HTTP handler that
// implements http.Handler.
func (h *Handler) createHTTPHandler(configName string) (http.Handler, error) {
	var httpHandler http.Handler
	switch configName {
	case HandlerContent.Name:
		httpHandler = &contentHandler{
			content:       h.content,
			auth:          h.auth,
			scopeSettings: h.scopeSettings,
		}
	case HandlerContents.Name:
		httpHandler = &contentsHandler{
			content:       h.content,
			auth:          h.auth,
			scopeSettings: h.scopeSettings,
		}
	default:
		return httpHandler, errUnknownConfig
	}

	return httpHandler, nil
}

// Start starts all HTTP handlers.
func (h *Handler) Start(multiplexer *mux.Router) error {
	for _, handler := range h.handlers {
		multiplexer.Handle(handler.identity.URL, handler.h)
	}
	return nil
}
