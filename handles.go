package httpserver

import (
	"log"
	"net/http"

	"github.com/elos/data"
	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

func Null() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	}
}

func Error(err error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		transfer.NewHTTPConnection(w, r, nil).ServerError(err)
	}
}

func BadMethod() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		transfer.NewHTTPConnection(w, r, nil).InvalidMethod()
	}
}

func BadAuth(reason string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		transfer.NewHTTPConnection(w, r, nil).Unauthorized()
	}
}

type AuthHandle func(http.ResponseWriter, *http.Request, httprouter.Params, *data.Access)

func Auth(h AuthHandle, s data.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		client, authenticated, err := transfer.AuthenticateRequest(s, r)
		if err != nil {
			log.Printf("An error occurred during authentication, err: %s", err)
			Error(err)(w, r, ps)
			return
		}

		if authenticated {
			access := data.NewAccess(client, s)
			h(w, r, ps, access)
			log.Printf("Client with id %s authenticated", client.ID())
		} else {
			log.Printf("Client with id %s authenticated", client.ID())
			BadAuth("Not authenticated")(w, r, ps)
		}

	}
}

func Post(k data.Kind, params []string) AuthHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, access *data.Access) {
		attrs := make(data.AttrMap)

		for _, k := range params {
			attrs[k] = r.FormValue(k)
		}

		c := transfer.NewHTTPConnection(w, r, access.Client)
		e := transfer.New(c, transfer.POST, k, attrs)
		go transfer.Route(e, access)
	}
}
