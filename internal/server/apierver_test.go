package server_test

import (
	"bytes"
	"log"
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

///////////////////////////////////////// Project
func TestCreateProject(t *testing.T) {
	var jsonStr = []byte(`{"name":"test name", "description": "description"}`)
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetProject(t *testing.T) {
	if _, err := s.Store.DB.Exec(
		`INSERT INTO projects (id, name, description)
			VALUES ('test_id1', 'test_name', 'test_description')`,
	); err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/projects/test_id1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProject(t *testing.T) {

	if _, err := s.Store.DB.Exec(
		`INSERT INTO projects (id, name, description)
			VALUES ('test_id2', 'test_name', 'test_description')`,
	); err != nil {
		log.Fatal(err)
	}

	var jsonStr = []byte(`{"name":"new test name", "description": "new description"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/projects/test_id2", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPatchProject(t *testing.T) {

	if _, err := s.Store.DB.Exec(
		`INSERT INTO projects (id, name, description)
			VALUES ('test_id3', 'test_name', 'test_description')`,
	); err != nil {
		log.Fatal(err)
	}

	var jsonStr = []byte(`{"name":"new test name"}`)
	req, _ := http.NewRequest("PATCH", "/api/v1/projects/test_id3", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
func TestDeleteProject(t *testing.T) {
	if _, err := s.Store.DB.Exec(
		`INSERT INTO projects (id, name, description)
			VALUES ('test_id4', 'test_name', 'test_description')`,
	); err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("DELETE", "/api/v1/projects/test_id4", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)
}
