package main

import (
	"github.com/osetr/rest-api/internal/server"
	"github.com/osetr/rest-api/internal/store"
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
