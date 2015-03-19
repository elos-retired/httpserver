package templates

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/transfer"
)

type (
	NotFoundError string

	RenderError struct {
		err error
	}

	ServerError struct {
		err error
	}
)

func NewNotFoundError(n Name) *NotFoundError {
	e := NotFoundError(n)
	return &e
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("templates error: could not find %s", string(*n))
}

func NewRenderError(err error) *RenderError {
	return &RenderError{err}
}

func (r RenderError) Error() string {
	return fmt.Sprintf("templates error: rendering failed %s", r.err)
}

func (r RenderError) Err() error {
	return r.err
}

func NewServerError(err error) *ServerError {
	return &ServerError{err}
}

func (s ServerError) Error() string {
	return fmt.Sprintf("templates error: server error %s", s.err)
}

func (s ServerError) Err() error {
	return s.err
}

const (
	RenderErrorResponseString   = "We had trouble rendering this screen, if the problem persists contact support"
	NotFoundErrorResponseString = RenderErrorResponseString
	ServerErrorResponseString   = RenderErrorResponseString
)

func CatchError(c *transfer.HTTPConnection, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *NotFoundError:
		c.ResponseWriter().Write([]byte(NotFoundErrorResponseString))
	case *RenderError:
		c.ResponseWriter().Write([]byte(RenderErrorResponseString))
	case *ServerError:
		c.ResponseWriter().Write([]byte(ServerErrorResponseString))
	default:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	}
	log.Printf("Templates package catch error caught %s", err)
}
