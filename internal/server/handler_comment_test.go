package server_test

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCreateComment(t *testing.T) {
	var jsonStr = []byte(`{"task_id": "def_id","text":"test text"}`)
	req, _ := http.NewRequest("POST", "/api/v1/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetComment(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/comments/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPatchComment(t *testing.T) {
	var jsonStr = []byte(`{"Text":"non_def_text"}`)
	req, _ := http.NewRequest("PATCH", "/api/v1/comments/def_id", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteComment(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/comments/def_id", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)
}
