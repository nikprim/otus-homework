package internalhttp

import (
	"context"
	golog "log"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type zerologWriter struct {
	zerolog zerolog.Logger
}

func (zlw zerologWriter) Write(p []byte) (n int, err error) {
	zlw.zerolog.Error().Msg(string(p))

	return len(p), nil
}

type Server struct {
	host string
	port int
	app  Application

	server *http.Server
}

type Application interface{}

func NewServer(host string, port int, app Application) *Server {
	return &Server{host, port, app, nil}
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         s.host + ":" + strconv.Itoa(s.port),
		Handler:      MakeRouter(s.app),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     golog.New(zerologWriter{log.Logger}, "", golog.LstdFlags),
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
