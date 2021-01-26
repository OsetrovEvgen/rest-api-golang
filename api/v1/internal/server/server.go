package server

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/osetr/rest-api/api/v1/internal/store"
	"github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	Config *Config
	Router *mux.Router
	Store  *store.Store
}

// NewAPIServer ...
func NewAPIServer(config *Config) *APIServer {
	return &APIServer{
		Config: config,
	}
}

// Start ...
func (s *APIServer) Start() error {
	if s.Router == nil {
		logrus.Fatal("set router before starting server")
	}
	if s.Store == nil {
		logrus.Fatal("set store before starting server")
	}
	logrus.Info("starting api server")
	return http.ListenAndServe(s.Config.BindAddr, s.Router)
}

// SetStore ...
func (s *APIServer) SetStore(config *store.Config) error {
	st := store.NewStore(config)
	if err := st.Open(); err != nil {
		return err
	}
	s.Store = st
	return nil
}

// SetRouter ...
func (s *APIServer) SetRouter() {
	s.Router = mux.NewRouter()

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	s.Router.Handle("/docs", sh)
	s.Router.Handle("/swagger.yaml", http.FileServer(http.Dir("./api/v1/internal/server/")))

	s.Router.HandleFunc("/api/v1/projects", s.createProject()).Methods("POST")
	s.Router.HandleFunc("/api/v1/projects", s.getProjects()).Methods("GET")
	s.Router.HandleFunc("/api/v1/projects/{id}", s.getProject()).Methods("GET")
	s.Router.HandleFunc("/api/v1/projects/{id}", s.updateProject()).Methods("PUT")
	s.Router.HandleFunc("/api/v1/projects/{id}", s.patchProject()).Methods("PATCH")
	s.Router.HandleFunc("/api/v1/projects/{id}", s.deleteProject()).Methods("DELETE")

	s.Router.HandleFunc("/api/v1/columns", s.createColumn()).Methods("POST")
	s.Router.HandleFunc("/api/v1/columns", s.getColumns()).Methods("GET")
	s.Router.HandleFunc("/api/v1/columns/{id}", s.getColumn()).Methods("GET")
	s.Router.HandleFunc("/api/v1/columns/{id}", s.patchColumn()).Methods("PATCH")
	s.Router.HandleFunc("/api/v1/columns/{id}", s.deleteColumn()).Methods("DELETE")
	s.Router.HandleFunc("/api/v1/columns/{id}/left", s.leftColumn()).Methods("POST")
	s.Router.HandleFunc("/api/v1/columns/{id}/right", s.rightColumn()).Methods("POST")

	s.Router.HandleFunc("/api/v1/tasks", s.createTask()).Methods("POST")
	s.Router.HandleFunc("/api/v1/tasks", s.getTasks()).Methods("GET")
	s.Router.HandleFunc("/api/v1/tasks/{id}", s.getTask()).Methods("GET")
	s.Router.HandleFunc("/api/v1/tasks/{id}", s.patchTask()).Methods("PATCH")
	s.Router.HandleFunc("/api/v1/tasks/{id}", s.deleteTask()).Methods("DELETE")
	s.Router.HandleFunc("/api/v1/tasks/{id}/down", s.downTask()).Methods("POST")
	s.Router.HandleFunc("/api/v1/tasks/{id}/up", s.upTask()).Methods("POST")
	s.Router.HandleFunc("/api/v1/tasks/{id}/left", s.leftTask()).Methods("POST")
	s.Router.HandleFunc("/api/v1/tasks/{id}/right", s.rightTask()).Methods("POST")

	s.Router.HandleFunc("/api/v1/comments", s.createComment()).Methods("POST")
	s.Router.HandleFunc("/api/v1/comments", s.getComments()).Methods("GET")
	s.Router.HandleFunc("/api/v1/comments/{id}", s.getComment()).Methods("GET")
	s.Router.HandleFunc("/api/v1/comments/{id}", s.patchComment()).Methods("PATCH")
	s.Router.HandleFunc("/api/v1/comments/{id}", s.deleteComment()).Methods("DELETE")
}
