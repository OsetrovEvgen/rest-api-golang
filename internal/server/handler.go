package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *APIServer) errorresp(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	s.Config.Logger.Info(r.Method + " " + r.URL.Path + " [" + strconv.FormatInt(int64(code), 10) + "]")
	json.NewEncoder(w).Encode(map[string]string{"detail": err.Error()})
}

func (s *APIServer) successresp(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	s.Config.Logger.Info(r.Method + " " + r.URL.Path + " [" + strconv.FormatInt(int64(code), 10) + "]")
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
