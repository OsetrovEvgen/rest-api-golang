package server_test

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCreateTask(t *testing.T) {
	var jsonStr = []byte(`{"column_id":"def_id", "name":"test_name", "description":"test_description"}`)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetTask(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/tasks/def_id", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPatchTask(t *testing.T) {
	var jsonStr = []byte(`{"description":"non_def_name"}`)
	req, _ := http.NewRequest("PATCH", "/api/v1/tasks/def_id", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLeftTask(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/tasks/def_id/left", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestRightTask(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/tasks/def_id/right", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUpTask(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/tasks/def_id/up", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestDownTask(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/tasks/def_id/down", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteTask(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)
}
