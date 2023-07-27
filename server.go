package wordy

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prplx/wordy/pkg/jsonlog"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, listener ...net.Listener) error {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		ErrorLog:       log.New(logger, "", 0),
	}

	if len(listener) > 0 {
		return s.httpServer.Serve(listener[0])
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
