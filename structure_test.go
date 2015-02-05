package httpserver_test

import (
	. "github.com/elos/httpserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var _ = Describe("Structure", func() {

	// HTTPMethods {{{
	Describe("HTTPMethods", func() {
		Context("Defines the generic http methods", func() {
			It("Define POST", func() {
				Expect(POST).ToNot(BeNil())
			})

			It("Defines GET", func() {
				Expect(GET).ToNot(BeNil())
			})
		})

		Context("Defines a map of acceptable HTTPMethods", func() {
			It("Defines the map", func() {
				Expect(HTTPMethods).ToNot(BeNil())
				Expect(HTTPMethods).To(BeAssignableToTypeOf(make(map[string]bool)))
			})

			Context("Access", func() {
				var (
					val bool
					ok  bool
				)

				methods := [4]string{POST, GET, "POST", "GET"}

				for i := range methods {
					m := methods[i]
					It(fmt.Sprintf("Accepts %s accessed using variable", m), func() {
						val, ok = HTTPMethods[m]
						Expect(ok).To(BeTrue())
						Expect(val).To(BeTrue())
					})
				}

			})

		})

	})
	// HTTPMethods }}}

	// HandlerMap {{{
	Describe("HandlerMap", func() {
		var (
			h HandlerMap
			m HandlerMap
		)

		route := "route"
		nullHandler := NewNullHandler()

		h = HandlerMap{}

		m = HandlerMap{
			route: nullHandler,
		}

		h[route] = m

		It("Can be recursively defined", func() {
			Expect(h[route]).To(Equal(m))
		})

		It("Implements http.Handler", func() {
			h.ServeHTTP(httptest.NewRecorder(), &http.Request{})
		})
	})
	// HandlerMap }}}

	// HTTPMethodHandler {{{
	Describe("HTTPMethodHandler", func() {

		sharedInvalidMethodHandler := NewNullHandler()

		var testInvalidMethod = func(r *http.Request) http.Handler {
			return sharedInvalidMethodHandler
		}

		h := NewHTTPMethodHandler(testInvalidMethod)
		d := NewHTTPMethodHandler(testInvalidMethod)
		Describe("Creation", func() {

			It("instantiates and creates Methods map", func() {
				Expect(h).ToNot(BeNil())
				Expect(h.Methods).ToNot(BeNil())

				Expect(d).ToNot(BeNil())
				Expect(d.Methods).ToNot(BeNil())

			})

			It("Defines its invalid method handler", func() {
				Expect(h.NewBadMethodHandler).ToNot(BeNil())
				Expect(d.NewBadMethodHandler).ToNot(BeNil())
			})

			It("instantiates a new method handler each time", func() {
				Expect(h).ToNot(Equal(d))
				Expect(h).To(BeAssignableToTypeOf(d))
			})
		})

		n1 := NewNullHandler()
		n2 := NewNullHandler()
		n3 := NewNullHandler()

		Describe("Specifying which methods to handle", func() {

			h.Handle(POST, n1)
			h.Handle(GET, n2)
			h.Handle("RANDO", n3)

			Context("Adds the methods it will handle", func() {
				var (
					val interface{}
					ok  bool
				)

				methods := [3]string{POST, GET, "RANDO"}

				for i := range methods {
					m := methods[i]
					It(fmt.Sprintf("had defined %s", m), func() {
						val, ok = h.Methods[POST]
						Expect(val).ToNot(BeNil())
						Expect(ok).To(BeTrue())
					})
				}

			})

			It("Doesn't define methods not h.Handle'd", func() {
				_, ok := h.Methods["ASF"]
				Expect(ok).ToNot(BeTrue())
			})

			It("Will override if told to handle method twice", func() {
				// Hack to make n1 and n2 different beyond pointers....sigh...
				n1.Handled[&http.Request{Method: GET}] = true
				n2.Handled[&http.Request{Method: POST}] = true

				h.Handle(POST, n1)
				h.Handle(POST, n2)

				val, ok := h.Methods[POST]

				Expect(ok).To(BeTrue())
				Expect(val).ToNot(BeNil())

				Expect(val).To(Equal(n2))
				Expect(val).ToNot(Equal(n1))

				h.Handle(POST, n1)
				val, ok = h.Methods[POST]
				Expect(ok).To(BeTrue())
				Expect(val).To(Equal(n1))

				// Hack to make n1 and n2 different clean up
				n1.Handled = make(map[*http.Request]bool)
				n2.Handled = make(map[*http.Request]bool)
			})
		})

		Describe("Serving requests", func() {
			w := httptest.NewRecorder()

			Context("Serving methods that were defined", func() {
				AfterEach(func() {
				})

				methods := map[string]*NullHandler{
					POST:    n1,
					GET:     n2,
					"RANDO": n3,
				}

				sharedInvalidMethodHandler.Handled = make(map[*http.Request]bool)

				It("Calls to ServeHTTP should dispatch correctly", func() {
					for method, handler := range methods {
						r := &http.Request{Method: method}

						By(fmt.Sprintf("Correctly dispatching %s", method))
						h.ServeHTTP(w, r)
						sharedInvalidMethodHandler.Handled = make(map[*http.Request]bool)
						Expect(handler.Handled).To(HaveLen(1))
						Expect(handler.Handled[r]).To(BeTrue())

						h.ServeHTTP(w, r)
						Expect(handler.Handled[r]).To(BeTrue())
						Expect(sharedInvalidMethodHandler.Handled).To(HaveLen(0))
					}
				})

			})

			Context("Serving methods that were not defined", func() {
				It("Passes unknown methods to a invalidMethodHandler", func() {
					sharedInvalidMethodHandler.Handled = make(map[*http.Request]bool)
					r := &http.Request{Method: "ADSF"}
					h.ServeHTTP(w, r)
					Expect(sharedInvalidMethodHandler.Handled).To(HaveLen(1))
					Expect(sharedInvalidMethodHandler.Handled[r]).To(BeTrue())
				})
			})
		})
	})
	// HTTPMethodHandler }}}

	// Routes Setup {{{
	Describe("Routes Setup", func() {
		n1 := NewNullHandler()
		n2 := NewNullHandler()
		n3 := NewNullHandler()

		hm := HandlerMap{
			"v1": HandlerMap{
				GET: n1,
				"users": HandlerMap{
					POST: n2,
				},
				"events": HandlerMap{
					POST: n3,
				},
			},
		}

		mux := http.NewServeMux()
		SetupRoutes(hm, mux, "")

		u1, _ := url.Parse("/v1")
		u2, _ := url.Parse("/v1/users")
		u3, _ := url.Parse("/v1/events")

		r1 := &http.Request{
			Method: GET,
			URL:    u1,
		}

		r2 := &http.Request{
			Method: POST,
			URL:    u2,
		}

		r3 := &http.Request{
			Method: POST,
			URL:    u3,
		}

		h1, _ := mux.Handler(r1)
		h2, _ := mux.Handler(r2)
		h3, _ := mux.Handler(r3)

		It("Registered GET to /v1/", func() {
			Expect(h1.(*HTTPMethodHandler).Methods[GET]).To(Equal(n1))
		})
		It("Registered POST to /v1/users", func() {
			Expect(h2.(*HTTPMethodHandler).Methods[POST]).To(Equal(n2))
		})
		It("Registered POST to /v1/events", func() {
			Expect(h3.(*HTTPMethodHandler).Methods[POST]).To(Equal(n3))
		})
	})
	// Routes Setup }}}

})
