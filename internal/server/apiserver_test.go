package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/osetr/rest-api/internal/server"
	"github.com/osetr/rest-api/internal/store"
)

var s *server.APIServer

func TestMain(m *testing.M) {
	s = server.NewAPIServer(server.NewConfig())
	sconf := store.NewConfig()
	sconf.DBName = "restapi_test"
	if err := s.SetStore(sconf); err != nil {
		logrus.Fatal(err)
	}
	s.SetRouter()
	s.Store.InitTestData()

	code := m.Run()

	s.Store.DB.Exec("DELETE FROM projects")
	s.Store.DB.Exec("DELETE FROM columns")
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
