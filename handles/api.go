package handles

import (
	"net/http"

	"github.com/elos/data"
	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

type ParamsList []string

func Params(v ...string) ParamsList {
	return ParamsList(v)
}

func Post(k data.Kind, params ParamsList) AccessHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, access data.Access) {
		attrs := make(data.AttrMap)

		for _, k := range params {
			attrs[k] = r.FormValue(k)
		}

		c := transfer.NewHTTPConnection(w, r, access.Client())
		e := transfer.New(c, transfer.POST, k, attrs)
		go transfer.Route(e, access)
	}
}

func Get(k data.Kind, params ParamsList) AccessHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, access data.Access) {
		attrs := make(data.AttrMap)

		for _, k := range params {
			attrs[k] = r.FormValue(k)
		}

		c := transfer.NewHTTPConnection(w, r, access.Client())
		e := transfer.New(c, transfer.GET, k, attrs)
		go transfer.Route(e, access)
	}
}
