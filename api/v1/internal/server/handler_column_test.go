package server_test

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCreateColumn(t *testing.T) {
	var jsonStr = []byte(`{"project_id":"def_id", "name":"test_name"}`)
	req, _ := http.NewRequest("POST", "/api/v1/columns", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetColumn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/columns/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPatchColumn(t *testing.T) {
	var jsonStr = []byte(`{"name":"non_def_name"}`)
	req, _ := http.NewRequest("PATCH", "/api/v1/columns/def_id", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLeftColumn(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/columns/def_id/left", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestRightColumn(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/columns/def_id/right", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteColumn(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/columns/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)
}
