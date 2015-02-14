package httpserver

import (
	"log"
	"net/http"

	"github.com/elos/agents"
	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

func WebSocket(u transfer.WebSocketUpgrader, connMan autonomous.Manager) AccessHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, a *data.Access) {
		conn, err := u.Upgrade(w, r, a)

		if err != nil {
			log.Printf("An error occurred while upgrading to the websocket protocol, err: %s", err)
			// gorilla.websocket will handle response to client
			return
		}

		log.Printf("Agent with id %s just connected over websocket", a.ID())

		agent := agents.NewClientDataAgent(conn, a)
		connMan.StartAgent(agent)
	}
}
