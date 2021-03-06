package handles

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

type AccessHandle func(http.ResponseWriter, *http.Request, httprouter.Params, data.Access)

func Auth(h AccessHandle, auther transfer.Authenticator, s data.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		client, authenticated, err := auther(s, r)
		if err != nil {
			log.Printf("An error occurred during authentication, err: %s", err)
			log.Printf("%+v", r)
			Error(err)(w, r, ps)
			return
		}

		if authenticated {
			access := data.NewAccess(client, s)
			h(w, r, ps, access)
			log.Printf("Client with id %s authenticated", client.ID())
		} else {
			http.Redirect(w, r, "/sign-in", 402)
		}

	}
}

func Access(h AccessHandle, a data.Access) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps, a)
	}
}
