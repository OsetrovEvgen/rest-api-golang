package main

import (
	"github.com/osetr/rest-api/api/v1/internal/server"
	"github.com/osetr/rest-api/api/v1/internal/store"
	"github.com/sirupsen/logrus"
)

func main() {
	s := server.NewAPIServer(server.NewConfig())
	if err := s.SetStore(store.NewConfig()); err != nil {
		logrus.Fatal(err)
	}
	s.SetRouter()

	if err := s.Start(); err != nil {
		logrus.Fatal(err)
	}
}
