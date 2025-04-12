package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fuckbug/api/internal/modules/app"
	"github.com/fuckbug/api/internal/modules/log"
	"github.com/fuckbug/api/internal/server/http/handlers"
)

type Server struct {
	server *http.Server
	logger handlers.Logger
}

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

func New(
	logger handlers.Logger,
	appService app.Service,
	logService log.Service,
	host string,
	port int,
) *Server {
	servers := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      loggingMiddleware(logger, NewHandler(logger, appService, logService)),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &Server{
		server: servers,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
