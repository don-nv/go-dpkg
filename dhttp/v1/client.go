package dhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"net"
	"net/http"
	"net/url"
	"time"
)

type IClient interface {
	// Do - executes 'req'. 'req' may be modified each time Do() is called applying some sort of middleware.
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http                   *http.Client
	requestsDefaultOptions []RequestOption
	requestsDefaultTTL     time.Duration
	log                    dlog.Logger
	logConfig              LoggerConfig
}

func MustNewClient(config ClientConfig, log dlog.Logger) Client {
	client, err := NewClient(config, log)
	derr.PanicOnE(err)

	return client
}

func NewClient(config ClientConfig, log dlog.Logger) (Client, error) {
	err := config.validate()
	if err != nil {
		return Client{}, fmt.Errorf("validating configuration: %w", err)
	}

	// Is a copy of http.DefaultTransport.
	var transport = &http.Transport{
		Proxy: config.Proxy,
		DialContext: func(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
			return dialer.DialContext
		}(
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if config.TLS != nil {
		transport.TLSClientConfig = config.TLS
	}

	client := Client{
		http:                   &http.Client{Transport: transport},
		requestsDefaultTTL:     config.RequestsDefaultTTL,
		requestsDefaultOptions: config.RequestsOptions,
		log:                    log.With().Name("http_client").Build(),
		logConfig:              config.Logger,
	}

	return client, nil
}

func (c Client) ResetLogging(config LoggerConfig) Client {
	c.logConfig = config

	return c
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

func (c Client) POST(
	ctx context.Context, url string, body []byte, options ...RequestOption,
) (
	*http.Response, error,
) {
	var log = c.log.With().Scope(ctx).Name("posting").Build()

	return c.processRequest(ctx, http.MethodPost, url, body, log, options...)
}

func (c Client) PATCH(
	ctx context.Context, url string, options ...RequestOption,
) (
	*http.Response, error,
) {
	var log = c.log.With().Scope(ctx).Name("patching").Build()

	return c.processRequest(ctx, http.MethodPatch, url, nil, log, options...)
}

func (c Client) PUT(
	ctx context.Context, url string, options ...RequestOption,
) (
	*http.Response, error,
) {
	var log = c.log.With().Scope(ctx).Name("putting").Build()

	return c.processRequest(ctx, http.MethodPut, url, nil, log, options...)
}

func (c Client) GET(
	ctx context.Context, url string, options ...RequestOption,
) (
	*http.Response, error,
) {
	var log = c.log.With().Scope(ctx).Name("getting").Build()

	return c.processRequest(ctx, http.MethodGet, url, nil, log, options...)
}

/*
Any
  - 'body' is optional and may be nil;
*/
func (c Client) Any(
	ctx context.Context, method, url string, body []byte, options ...RequestOption,
) (
	*http.Response, error,
) {
	var log = c.log.With().Scope(ctx).Name("sending").Build()

	return c.processRequest(ctx, method, url, body, log, options...)
}

/*
processRequest
  - 'body' is optional and may be nil;
*/
func (c Client) processRequest(
	ctx context.Context, method, url string, body []byte, log dlog.Logger, options ...RequestOption,
) (
	buffer *http.Response, err error,
) {
	if d := c.requestsDefaultTTL; d > 0 {
		_, ok := ctx.Deadline()
		if !ok {
			c, cancel := dctx.WithTTLTimeout(ctx, d)
			defer cancel()

			ctx = c
		}
	}

	req, err := c.newRequest(ctx, method, url, body, options...)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return c.sendRequest(req, body, log)
}

// newRequest - creates new request and OptionRequest().
func (c Client) newRequest(
	ctx context.Context, method, url string, body []byte, options ...RequestOption,
) (
	*http.Request, error,
) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("NewRequestWithContext: %w", err)
	}

	req = OptionRequest(req, c.requestsDefaultOptions...)
	return OptionRequest(req, options...), nil
}

// sendRequest - sends 'req' and logs request and response according to LoggerConfig{}.
func (c Client) sendRequest(req *http.Request, body []byte, log dlog.Logger) (*http.Response, error) {
	var logData = LogWithClientRequest(req, body, log.With(), c.logConfig)

	var err error
	defer func() { l := logData.Build(); l.CatchED(&err) }()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("retrier.Send: %w", err)
	}

	logData, err = LogWithClientResponse(resp, logData, c.logConfig)
	if err != nil {
		return nil, fmt.Errorf("creating log with client response: %w", err)
	}

	return resp, nil
}

type ClientConfig struct {
	TLS             *tls.Config
	RequestsOptions []RequestOption
	// RequestsDefaultTTL - is applied for each request context if it has no TTL.
	RequestsDefaultTTL time.Duration
	Logger             LoggerConfig
	Proxy              func(*http.Request) (*url.URL, error)
}

func (c ClientConfig) validate() error {
	const minReqsTTL = 5 * time.Millisecond
	if v := c.RequestsDefaultTTL; v > 0 {
		if v < minReqsTTL {
			return fmt.Errorf("too short requests default timeout, %q < %q", v, minReqsTTL)
		}
	}

	return nil
}
