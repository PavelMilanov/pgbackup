package main

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logrus.Info("Сервер запущен")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(grace time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), grace*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		logrus.Info("Сервер остановлен")
	}
	return s.httpServer.Shutdown(ctx)
}
