package server_test

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCreateProject(t *testing.T) {
	var jsonStr = []byte(`{"name":"test name", "description": "description"}`)
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetProject(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/projects/def_id", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProject(t *testing.T) {
	var jsonStr = []byte(`{"name":"non_def_name", "description": "non_def_description"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/projects/def_id", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPatchProject(t *testing.T) {
	var jsonStr = []byte(`{"name":"non_def_name"}`)
	req, _ := http.NewRequest("PATCH", "/api/v1/projects/def_id", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
func TestDeleteProject(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)
}
