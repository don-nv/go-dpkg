package dhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
	log  dlog.Logger
}

func NewServer(config ServerConfig, router http.Handler, log dlog.Logger, options ...ServerOption) Server {
	srv := http.Server{
		Addr:              config.Address,
		Handler:           router,
		ReadHeaderTimeout: config.RequestReadHeaderTTL,
		ReadTimeout:       config.RequestReadTTL,
		IdleTimeout:       config.IdleTimeout,
	}

	for _, option := range options {
		option(&srv)
	}

	return Server{
		http: &srv,
		log:  log.With().Name("http_server").Build(),
	}
}

func (s Server) Router() http.Handler {
	return s.http.Handler
}

func (s Server) Running(ctx context.Context) (err error) {
	var group = dsync.NewOneTimeGroup(ctx)
	defer group.WaitE(&err)

	group.GoUntilWait("http_server_closing", func(ctx context.Context) error {
		<-ctx.Done()

		err := s.http.Shutdown(context.Background())
		if !derr.In(err, nil, http.ErrServerClosed) {
			return fmt.Errorf("shutting down: %w", err)
		}

		return nil
	})

	group.GoUntilWait("http_server_running", func(ctx context.Context) error {
		listener, err := listenTCP(s.http.Addr)
		if err != nil {
			return fmt.Errorf("listening %s: %w", networkTCP, err)
		}

		var I = s.log.I().Scope(ctx).Any("address", s.http.Addr).Any("network", networkTCP)

		if s.http.TLSConfig != nil {
			listener = tls.NewListener(listener, s.http.TLSConfig)

			I = I.Any("with_tls", true)
		}

		defer I.Write("running...").Write("...done")

		err = s.http.Serve(listener)
		if !derr.In(err, nil, http.ErrServerClosed) {
			return fmt.Errorf("serving: %w", err)
		}

		return nil
	})

	return nil
}

type ServerOption func(server *http.Server)

func OptionServerWithMiddleware(fs ...func(next http.Handler) http.Handler) ServerOption {
	return func(server *http.Server) {
		for _, f := range fs {
			server.Handler = f(server.Handler)
		}
	}
}

func OptionServerWithTLS(config *tls.Config) ServerOption {
	return func(server *http.Server) {
		server.TLSConfig = config
	}
}

type ServerConfig struct {
	Address              string
	RequestReadHeaderTTL time.Duration
	RequestReadTTL       time.Duration
	IdleTimeout          time.Duration
}
