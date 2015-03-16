package handles

import (
	"net/http"

	"github.com/elos/data"
	"github.com/elos/httpserver/templates"
	"github.com/elos/models"
	"github.com/elos/models/user"
	"github.com/elos/transfer"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

func SignIn(s sessions.Store) AccessHandle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
		session, _ := s.Get(r, transfer.AuthSession)

		session.Values[transfer.ID] = a.Client().ID().(bson.ObjectId).Hex()
		session.Values[transfer.Key] = a.Client().(models.User).Key()
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			http.Redirect(w, r, "/calendar", http.StatusFound)
		}
	}
}

func RegisterHandle(s data.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u, err := user.NewWithName(s, r.FormValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templates.Render(w, r, templates.AccountCreated, u)
	}
}
