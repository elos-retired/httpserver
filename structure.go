package httpserver

import (
	"fmt"
	"log"
	"net/http"
)

// HTTPMethods {{{

const POST string = "POST"
const GET string = "GET"

var HTTPMethods = map[string]bool{
	POST: true,
	GET:  true,
}

// HTTPMethods }}}

// HandlerMap {{{

type HandlerMap map[string]http.Handler

func (h HandlerMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// HandlerMap }}}

// HTTPMethodHandler {{{

// Redirects http requests based on the request's HTTP method
type HTTPMethodHandler struct {
	NewBadMethodHandler BadMethodHandlerConstructor
	Methods             map[string]http.Handler
}

// Satisfies http.Handler interface, will dispatch ServeHTTP to
// one of it's method handlers, if one doesn't exist for the
// specified method then it handles the response with a BadMethodHandler
func (h *HTTPMethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := h.Methods[r.Method]
	if ok {
		handler.ServeHTTP(w, r)
	} else {
		h.NewBadMethodHandler(r).ServeHTTP(w, r)
	}
}

// Registers a handler for a method
func (h *HTTPMethodHandler) Handle(method string, handler http.Handler) {
	h.Methods[method] = handler
}

// Creates a new httpMethodHandler
func NewHTTPMethodHandler(c BadMethodHandlerConstructor) *HTTPMethodHandler {
	return &HTTPMethodHandler{
		Methods:             make(map[string]http.Handler),
		NewBadMethodHandler: c,
	}
}

// HTTPMethodHandler }}}

// Routes Setup {{{

// joins a route prefix with the route
// e.g., join("/hey", "ho") => "/hey/ho"
func join(prefix string, route string) string {
	return fmt.Sprintf("%s/%s", prefix, route)
}

// Calls the recursively defined SetupRoutes
func SetupHTTPRoutes(hm HandlerMap) {
	SetupRoutes(hm, http.DefaultServeMux, "")
}

// recursively sets up the routes
func SetupRoutes(hm HandlerMap, mux *http.ServeMux, prefix string) {
	methodHandler := NewHTTPMethodHandler(NewBadMethodHandler)
	for routeName, handler := range hm {
		// type assert
		subHM, ok := handler.(HandlerMap)

		// We are being pointed to another handler map
		if ok {
			SetupRoutes(subHM, mux, join(prefix, routeName))
		} else { // this is a handler
			_, isHTTPMethod := HTTPMethods[routeName]
			if isHTTPMethod {
				methodHandler.Handle(routeName, handler)
			} else {
				// if you implement this, test it
				log.Print("this functionality is not defined")
			}
		}
	}
	if prefix != "" {
		mux.Handle(prefix, methodHandler)
	}
}

// Routes Setup }}}
