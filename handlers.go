package httpserver

import (
	"log"
	"net/http"
	"sync"

	"github.com/elos/data"
	"github.com/elos/transfer"
)

// NullHandler (Testing) {{{

type NullHandler struct {
	Handled map[*http.Request]bool
	m       sync.Mutex
}

func (h *NullHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()
	h.Handled[r] = true
	h.m.Unlock()
}

func (h *NullHandler) Reset() *NullHandler {
	h.m.Lock()
	h.Handled = make(map[*http.Request]bool)
	h.m.Unlock()
	return h
}

func NewNullHandler() *NullHandler {
	return (&NullHandler{}).Reset()
}

// NullHandler (Testing) }}}

//  ErrorHandler {{{

// Allows route to handle an error
type ErrorHandlerConstructor func(error) http.Handler

// underlying information needed to form an error response
type ErrorHandler struct {
	Err error
}

// implement http.Handler
func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transfer.ServerError(w, h.Err)
}

// Returns a handler capable of serving the error information
func NewErrorHandler(err error) http.Handler {
	return &ErrorHandler{
		Err: err,
	}
}

// }}}

// ResourceHandler {{{

// Allows route to handle returning a json resource
type ResourceHandlerConstructor func(int, interface{}) http.Handler

// underlying information needed to return a json resource
type ResourceHandler struct {
	Code     int
	Resource interface{}
}

// implements http.Handler
func (h *ResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transfer.WriteResourceResponse(w, h.Code, h.Resource)
}

// Returns a handler capable of serving the resource
func NewResourceHandler(code int, resource interface{}) http.Handler {
	return &ResourceHandler{
		Code:     code,
		Resource: resource,
	}
}

// }}}

// BadMethodHandler {{{

/*
	Allows route to handle a suspected invalid method
	- Should only be used by HTTPMethodHandler
*/
type BadMethodHandlerConstructor func(*http.Request) http.Handler

// underlying information need to notify user of invalid method
type BadMethodHandler struct {
	RequestedMethod string
}

// implemens http.Handler
func (h *BadMethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transfer.InvalidMethod(w)
}

// Returns a handler capable of notifying the user of the invalid method
func NewBadMethodHandler(r *http.Request) http.Handler {
	return &BadMethodHandler{
		RequestedMethod: r.Method,
	}
}

// InvalidMethodHandler }}}

// UnauthorizedHandler {{{

// Allows a route to indicate the agent is unauthorized
type UnauthorizedHandlerConstructor func(string) http.Handler

/*
	underlying information necessary to inform client of
	lack of credentials
*/
type UnauthorizedHandler struct {
	Reason string
}

// implements http.Handler
func (h *UnauthorizedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transfer.Unauthorized(w)
}

// Returns a handler capable of serving the unauthorized error
func NewUnauthorizedHandler(reason string) http.Handler {
	return &UnauthorizedHandler{
		Reason: reason,
	}
}

// UnauthorizedHandler }}}

// Authenticators {{{

var DefaultAuthenticator transfer.RequestAuthenticator = transfer.AuthenticateRequest

// Authenticators }}}

// AuthenticationHandler {{{

type AuthenticationHandler struct {
	data.Store
	Authenticator          transfer.RequestAuthenticator
	NewErrorHandler        ErrorHandlerConstructor
	NewUnauthorizedHandler UnauthorizedHandlerConstructor
	AuthenticatedHandler   AuthenticatedHandler
}

func NewAuthenticationHandler(s data.Store, a transfer.RequestAuthenticator, eh ErrorHandlerConstructor,
	uh UnauthorizedHandlerConstructor, t AuthenticatedHandler) http.Handler {
	foo := &AuthenticationHandler{
		Authenticator:          a,
		NewErrorHandler:        eh,
		NewUnauthorizedHandler: uh,
		AuthenticatedHandler:   t,
	}

	foo.Store = s

	return foo
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agent, authenticated, err := h.Authenticator(h.Store, r)

	if err != nil {
		log.Printf("An error occurred during authentication, err: %s", err)
		h.NewErrorHandler(err).ServeHTTP(w, r)
		return
	}

	if authenticated {
		h.AuthenticatedHandler.ServeHTTP(w, r, agent)
		log.Printf("Agent with id %s authenticated", agent.ID())
	} else {
		h.NewUnauthorizedHandler("Not authenticated").ServeHTTP(w, r)
	}
}

// AuthenticationHandler }}}

// AuthenticatedHandlerFunc {{{

type AuthenticatedHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, data.Identifiable)
}

type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, data.Identifiable)

func (f AuthenticatedHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, a data.Identifiable) {
	f(w, r, a)
}

type AgentHandler struct {
	Agent data.Identifiable
	Fn    AuthenticatedHandlerFunc
}

func (h *AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Fn(w, r, h.Agent)
}

func NewAgentHandler(agent data.Identifiable, fn AuthenticatedHandlerFunc) http.Handler {
	return &AgentHandler{
		Agent: agent,
		Fn:    fn,
	}
}

// AuthenticatedHandlerFunc }}}
