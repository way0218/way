package way

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/way0218/way/logger"
	"github.com/way0218/way/transport"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	option  options
	context context.Context
	cancel  func()
	logger  *zap.Logger
}

func New(opts ...Option) *App {
	option := options{
		ctx:     context.Background(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		logger:  logger.NewLogger(),
		servers: nil,
	}
	for _, o := range opts {
		o(&option)
	}
	ctx, cancel := context.WithCancel(option.ctx)
	return &App{
		option:  option,
		context: ctx,
		cancel:  cancel,
		logger:  option.logger,
	}
}

func (a *App) Logger() *zap.Logger {
	return a.option.logger
}

func (a *App) Server() []transport.Server {
	return a.option.servers
}

func (a *App) Run() error {
	a.logger.Info("app run", zap.String("service_name", a.option.name), zap.String("version", a.option.version))
	g, ctx := errgroup.WithContext(a.context)
	for _, srv := range a.option.servers {
		server := srv
		g.Go(func() error {
			<-ctx.Done()
			return server.Stop()
		})
		g.Go(func() error {
			return server.Start()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.option.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

type Option func(opt *options)

type options struct {
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	logger  *zap.Logger
	servers []transport.Server
}

func Name(name string) Option {
	return func(option *options) {
		option.name = name
	}
}

func Version(version string) Option {
	return func(option *options) {
		option.version = version
	}
}

func Server(server ...transport.Server) Option {
	return func(option *options) {
		option.servers = server
	}
}

func Logger(logger *zap.Logger) Option {
	return func(option *options) {
		option.logger = logger
	}
}
