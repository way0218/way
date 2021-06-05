package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/way0218/way/logger"
	"github.com/way0218/way/transport"
	"go.uber.org/zap"
)

const (
	DefaultTimeout = 30 * time.Second
)

type Server struct {
	*http.Server
	listener net.Listener
	network  string
	address  string
	timeout  time.Duration
	router   *http.ServeMux
	logger   *zap.Logger
}

var _ transport.Server = (*Server)(nil)

type ServerOption func(server *Server)

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: DefaultTimeout,
		router:  http.NewServeMux(),
		logger:  logger.NewLogger(),
	}
	for _, opt := range opts {
		opt(srv)
	}
	srv.Server = &http.Server{Handler: srv}
	return srv
}

func Network(network string) ServerOption {
	return func(server *Server) {
		server.network = network
	}
}

func Address(addr string) ServerOption {
	return func(server *Server) {
		server.address = addr
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(server *Server) {
		server.timeout = timeout
	}
}

func Logger(logger *zap.Logger) ServerOption {
	return func(server *Server) {
		server.logger = logger
	}
}

func Router(router *http.ServeMux) ServerOption {
	return func(server *Server) {
		server.router = router
	}
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.timeout)
	defer cancel()
	h := func(ctx context.Context, req *http.Request) {
		s.router.ServeHTTP(res, req)
	}
	h(ctx, req.WithContext(ctx))
}

func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.listener = lis
	s.logger.Info(fmt.Sprintf("[HTTP] server listening on: %s", lis.Addr().String()))
	if err := s.Serve(lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("[HTTP] server stopping")
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.Shutdown(ctx)
}

func (s *Server) Router() *http.ServeMux {
	return s.router
}
