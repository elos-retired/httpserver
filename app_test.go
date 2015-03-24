package httpserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/elos/data"
	"github.com/elos/models/persistence"
	"github.com/elos/models/user"
	"gopkg.in/mgo.v2/bson"
)

func testStore() data.Store {
	return persistence.Store(persistence.MongoMemoryDB())
}

func testServerWithStore(s data.Store) *HTTPServer {
	server := NewHTTPServer("not used", 0, s)

	setupRoutes(server)
	setupAPI(server)
	return server
}

func testServer() *HTTPServer {
	return testServerWithStore(testStore())
}

func newRequest(action string, route string, body io.Reader, t *testing.T) *http.Request {
	req, err := http.NewRequest(action, route, body)
	if err != nil {
		t.Fatalf("Error while creating request: %s", err)
	}

	return req
}

func expectSuccess(w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != 200 {
		t.Log(w.Body.String())
		t.Errorf("Expected 200 status code, got %d", w.Code)
	}
}

func expectBodyContains(w *httptest.ResponseRecorder, s string, t *testing.T) {
	if !strings.Contains(w.Body.String(), s) {
		t.Errorf("Expected response body to contain %s", s)
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	s := testServer()
	w := httptest.NewRecorder()

	s.Router.ServeHTTP(w, newRequest("GET", "/", nil, t))
	expectSuccess(w, t)
	expectBodyContains(w, "elos", t)
}

func TestSignInGet(t *testing.T) {
	t.Parallel()

	s := testServer()
	w := httptest.NewRecorder()

	s.Router.ServeHTTP(w, newRequest("GET", "/sign-in", nil, t))
	expectSuccess(w, t)
	expectBodyContains(w, "Sign In", t)
}

func TestSignInPost(t *testing.T) {
	t.Parallel()
	store := testStore()
	// server := testServerWithStore(store)

	u, err := user.New(store)
	if err != nil {
		t.Fatalf("Error while creating user: %s", err)
	}

	u.SetKey(user.NewKey())

	if err = store.Save(u); err != nil {
		t.Fatalf("Error while saving user: %s", err)
	}

	/*
		w := httptest.NewRecorder()

		req := newRequest("POST", "/sign-in", nil, t)
		authValues := url.Values{"id": {u.ID().(bson.ObjectId).Hex()}, "key": {u.Key()}}
		req.PostForm = authValues

		server.ServeHTTP(w, req)

		if w.Code != 302 {
			t.Errorf("Expected redirect")
		}
	*/

	client := &http.Client{}

	s := NewHTTPServer("localhost", 8001, store)
	go s.Start()
	s.WaitStart()

	resp, err := client.PostForm("http://localhost:8001/sign-in", url.Values{"id": {u.ID().(bson.ObjectId).Hex()}, "key": {u.Key()}})

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 302 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		t.Log(buf.String())
		t.Errorf("expect redirect")
	}

	resp, err = client.Get("http://localhost:8001/sign-in")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expect success")
	}

	go s.Stop()
	s.WaitStop()

	/*
		req = newRequest("GET", "/user/calendar", nil, t)
		w = httptest.NewRecorder()

		server.ServeHTTP(w, req)
		expectSuccess(w, t)
	*/
}
