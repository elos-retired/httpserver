package httpserver_test

import (
	. "github.com/elos/httpserver"
	"github.com/elos/transfer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

var _ = Describe("Handlers", func() {

	// NullHandler {{{
	Describe("NullHandler", func() {
		var (
			h *NullHandler
		)

		BeforeEach(func() {
			h = NewNullHandler()
		})

		Describe("NewHullHandler", func() {
			It("Allocates and returns a new NullHandler", func() {
				Expect(h).ToNot(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&NullHandler{}))
				Expect(h.Handled).ToNot(BeNil())
				Expect(h.Handled).To(BeEmpty())
			})
		})

		Describe("ServeHTTP", func() {
			It("Adds a request to it's handled map when asked to serve it", func() {
				r := &http.Request{Method: "FOOBAR"}
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				Expect(h.Handled).To(HaveKeyWithValue(r, true))
			})
		})

		Describe("Reset()", func() {
			It("Wipes it Handled map", func() {
				h.Handled[&http.Request{}] = true
				h.Reset()
				Expect(h.Handled).To(BeEmpty())
			})
		})
	})
	// NullHandler }}}

	// ErrorHandler {{{
	Describe("ErrorHandler", func() {

		var (
			err error
			h   *ErrorHandler
		)

		BeforeEach(func() {
			err = errors.New("This is a test error")
			h = NewErrorHandler(err).(*ErrorHandler)
		})

		Describe("NewErrorHandler", func() {
			It("Allocates and returns a new ErrorHandler", func() {
				Expect(h).ToNot(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&ErrorHandler{}))
				Expect(h.Err).To(Equal(err))
			})
		})

		Describe("ServeHTTP", func() {
			var (
				w1 *httptest.ResponseRecorder
				w2 *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				w1 = httptest.NewRecorder()
				w2 = httptest.NewRecorder()
				transfer.ServerError(w1, err)
				h.ServeHTTP(w2, &http.Request{})
			})

			It("Uses transfer.to write the error response", func() {
				Expect(w1.Body).To(Equal(w2.Body))
			})
		})

	})
	// ErrorHandler }}}

	// ResourceHandler {{{
	Describe("ResoureHandler", func() {
		var (
			resource map[string]interface{}
			code     int
			h        *ResourceHandler
		)

		BeforeEach(func() {
			resource = map[string]interface{}{
				"hey": "ho",
			}

			code = 200

			h = NewResourceHandler(code, resource).(*ResourceHandler)
		})

		Describe("NewResourceHandler", func() {
			It("Allocates and returns a new ResourceHandler", func() {
				Expect(h).ToNot(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&ResourceHandler{}))
			})

			It("Sets up necessary information", func() {
				By("Setting status code")
				Expect(h.Code).To(Equal(code))
				By("Setting resource")
				Expect(h.Resource).To(Equal(resource))
			})
		})

		Describe("ServeHTTP", func() {
			var (
				w1 *httptest.ResponseRecorder
				w2 *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				w1 = httptest.NewRecorder()
				w2 = httptest.NewRecorder()
				transfer.WriteResourceResponse(w1, code, resource)
				h.ServeHTTP(w2, &http.Request{})
			})

			It("Uses transfer.WriteResourceResponse", func() {
				Expect(w1.Body).To(Equal(w2.Body))
				Expect(w1.Code).To(Equal(w2.Code))
			})
		})

	})
	// ResourceHandler }}}

	// BadMethodHandler {{{
	Describe("BadMethodHandler", func() {
		var (
			r *http.Request
			h *BadMethodHandler
		)

		BeforeEach(func() {
			r = &http.Request{Method: "BOOP"}
			h = NewBadMethodHandler(r).(*BadMethodHandler)
		})

		Describe("NewBadMethodHandler", func() {
			It("Should allocate and return a new BadMethodHandler", func() {
				Expect(h).ToNot(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&BadMethodHandler{}))
			})

			It("Should set it's RequestedMethod field", func() {
				Expect(h.RequestedMethod).To(Equal(r.Method))
			})
		})

		Describe("ServeHTTP", func() {
			var (
				w1 *httptest.ResponseRecorder
				w2 *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				w1 = httptest.NewRecorder()
				w2 = httptest.NewRecorder()
				transfer.InvalidMethod(w1)
				h.ServeHTTP(w2, &http.Request{})
			})

			It("Uses transfer.InvalidMethod", func() {
				Expect(w1.Body).To(Equal(w2.Body))
				Expect(w1.Code).To(Equal(w2.Code))
			})
		})
	})
	// BadMethodHandler }}}

	// UnauthorizedHandler {{{
	Describe("UnauthorizedHandler", func() {
		var (
			reason string
			h      *UnauthorizedHandler
		)
		BeforeEach(func() {
			reason = "asdf"
			h = NewUnauthorizedHandler(reason).(*UnauthorizedHandler)
		})

		Describe("NewUnauthorizedHandler", func() {
			It("Allocates and returns a new UnauthorizedHandler", func() {
				Expect(h).ToNot(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&UnauthorizedHandler{}))
			})

			It("Sets the reason field", func() {
				Expect(h.Reason).To(Equal(reason))
			})
		})

		Describe("ServeHTTP", func() {
			var (
				w1 *httptest.ResponseRecorder
				w2 *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				w1 = httptest.NewRecorder()
				w2 = httptest.NewRecorder()
				transfer.Unauthorized(w1)
				h.ServeHTTP(w2, &http.Request{})
			})

			It("Uses transfer.Unauthorized", func() {
				Expect(w1.Body).To(Equal(w2.Body))
				Expect(w1.Code).To(Equal(w2.Code))
			})
		})

	})

	// UnauthorizedHandler {{{

	// AuthenticationHandler }}}
	// AuthenticationHandler }}}

	// Authenticators {{{
	Describe("Authenticators", func() {
		It("Defines DefaultAuthenticator", func() {
			Expect(DefaultAuthenticator).ToNot(BeNil())
		})
	})
	// Authenticators }}}

	// AuthenticationHandler {{{
	Describe("AuthenticationHandler", func() {
		s := data.NewNullStoreWithType("mongo")
		a, _ := user.New(s)
		a.SetID(s.NewObjectID())
		authed := true
		var err error = nil

		var authenticator transfer.RequestAuthenticator = func(r *http.Request) (data.Identifiable, bool, error) {
			return a, authed, err
		}

		n1 := NewNullHandler()

		var errHandlerC = func(e error) http.Handler {
			return n1
		}

		n2 := NewNullHandler()

		var unauthHandlerC = func(r string) http.Handler {
			return n2
		}

		transferCalledCount := 0

		var transferFunc AuthenticatedHandlerFunc = func(w http.ResponseWriter, r *http.Request, a data.Agent) {
			transferCalledCount = transferCalledCount + 1

		}

		h := NewAuthenticationHandler(authenticator, errHandlerC, unauthHandlerC, transferFunc)

		Describe("NewAuthenticationHandler", func() {
			It("Allocates and returns a new AuthenticationHandler", func() {
				Expect(h).NotTo(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&AuthenticationHandler{}))
			})

			It("Sets fields", func() {
				By("sets Authenticator")
				h := h.(*AuthenticationHandler)
				Expect(h.Authenticator).NotTo(BeNil())
				Expect(h.NewErrorHandler).NotTo(BeNil())
				Expect(h.NewUnauthorizedHandler).NotTo(BeNil())
				Expect(h.AuthenticatedHandler).NotTo(BeNil())
			})
		})

		Describe("ServeHTTP", func() {
			w := httptest.NewRecorder()
			r := &http.Request{}

			// Contexts would be more appropiate here
			// but have to do it this way because of the order
			// Ginkgo evaluates shit
			It("Handles Authentication Error", func() {
				transferCalledCount = 0
				err = errors.New("This is an error during the authentication process")

				h.ServeHTTP(w, r)
				By("Uses the ErrorHandlerConstructor to respond")
				Expect(n1.Handled).To(HaveKeyWithValue(r, true))

				By("Does not use an unauthorized handler")
				Expect(n2.Handled).To(BeEmpty())
				By("Doesn't allow transfer")
				Expect(transferCalledCount).To(BeZero())

				n1.Reset()
				err = nil
			})

			It("Handles Authentication Failure", func() {
				transferCalledCount = 0
				authed = false
				h.ServeHTTP(w, r)

				By("Uses the UnauthorizedHandlerConstructor to response")
				Expect(n2.Handled).To(HaveKeyWithValue(r, true))

				By("Does not touch the error handler")
				Expect(n1.Handled).To(BeEmpty())
				By("Doesn't allow transfer")
				Expect(transferCalledCount).To(BeZero())

				n2.Reset()
				authed = true
			})

			It("Handles Authentication Successful", func() {
				transferCalledCount = 0
				h.ServeHTTP(w, r)

				By("Uses the TransferFunction")
				Expect(transferCalledCount).To(Equal(1))

				By("Does not touch an error handler")
				Expect(n1.Handled).To(BeEmpty())
				By("Does not touch an unauthorized handler")
				Expect(n2.Handled).To(BeEmpty())
			})
		})

	})
	// AuthenticationHandler }}}

	// AuthenticatedHandler {{{
	Describe("AuthenticatedHandler", func() {

		agent := user.New()
		agent.SetID(data.NewObjectID())

		var (
			wr http.ResponseWriter
			rr *http.Request
			ra data.Agent
		)

		var tf = func(w http.ResponseWriter, r *http.Request, a data.Agent) {
			wr = w
			rr = r
			ra = a
		}

		h := NewAgentHandler(agent, tf)

		Describe("NewAuthenticatedHandler", func() {
			It("Allocates and returns a new AuthenticatedHandler", func() {
				Expect(h).NotTo(BeNil())
				Expect(h).To(BeAssignableToTypeOf(&AgentHandler{}))
			})

			It("Sets fields", func() {
				h := h.(*AgentHandler)
				Expect(h.Agent).To(Equal(agent))
				// Still can't figure out the function equality thing
				Expect(h.Fn).NotTo(BeNil())
			})
		})

		Describe("ServeHTTP", func() {
			It("calls Fn with writer, request, and agent", func() {
				w := httptest.NewRecorder()
				r := &http.Request{}
				h.ServeHTTP(w, r)
				Expect(wr).To(Equal(w))
				Expect(rr).To(Equal(r))
				Expect(ra).To(Equal(agent))
			})
		})
	})
	// AuthenticatedHandler }}}
})
