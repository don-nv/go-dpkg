package roundbreaker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"io"
	"net/http"
	"time"
)

type Breaker struct {
	client dhttp.IClient
	rounds *Rounds
	config Config
	log    dlog.Logger
}

func MustNewBreaker(client dhttp.IClient, rounds *Rounds, config Config, log dlog.Logger) Breaker {
	derr.PanicOnE(config.validate())

	return Breaker{
		client: client,
		rounds: rounds,
		config: config,
		log:    log.With().Name("round_breaker").Build(),
	}
}

// SetLoggerConfig - replaces default logger config and returns Breaker with modified config.
func (b Breaker) SetLoggerConfig(config dhttp.LoggerConfig) Breaker {
	b.config.Logger = config

	return b
}

/*
Send - is the same as Do(), but creates request by itself. 'ctx' controls request lifetime. 'body' may be nil. 'options'
get applied before Config.RequestsOptions.
*/
func (b Breaker) Send(
	ctx context.Context, method, path string, body []byte, options ...dhttp.RequestOption,
) (
	*http.Response, error,
) {
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	for _, option := range options {
		req = option(req)
	}

	return b.Do(req)
}

/*
Do - sends 'req' to the next host in rounds. If an error occurred during sending to the current host or response status
code >= 500, then moves to the next host and repeats sending. Once all host attempts were exceeded (sequentially), this
host is suspended for delay and not considered in rounds for a while. Request duration is controlled by its context. If
request has not TTL (see dhttp.RequestTTL()), then Config.RequestsDefaultTTL is applied.
*/
func (b Breaker) Do(req *http.Request) (*http.Response, error) {
	for _, option := range b.config.RequestsOptions {
		req = option(req)
	}

	if req.Body == nil {
		req.Body = io.NopCloser(bytes.NewReader(nil))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}

	var getBody = func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(body))
	}

	if dhttp.RequestTTL(req) == 0 {
		var r, cancel = dhttp.RequestContextWithTTL(req, b.config.RequestsDefaultTTL)
		defer cancel()

		req = r
	}

	b.logInfoDoRequest(req, body)
	defer b.logDebugDoDone(req.Context(), time.Now())

	resp, err := b.sendOverHosts(req, getBody)
	if err != nil {
		return nil, fmt.Errorf("sending over hosts: %w", err)
	}

	err = b.logInfoDoResponse(req.Context(), resp)
	if err != nil {
		return nil, fmt.Errorf("logging do response: %w", err)
	}

	return resp, nil
}

func (b Breaker) sendOverHosts(req *http.Request, getBody func() io.ReadCloser) (*http.Response, error) {
	defer func() { req.URL.Scheme = ""; req.URL.Host = ""; req.Body = getBody() }() // Recover request.

	for errLast := error(nil); ; {
		var nextStartedAt = time.Now()

		host, err := b.rounds.nextHost(req.Context())
		if err != nil {
			return nil, errors.Join(errLast, fmt.Errorf("getting next host: %w", err))
		}

		b.logDebugSendOverHostsHostChosen(req.Context(), host, nextStartedAt)

		var (
			resp         *http.Response
			reqStartedAt = time.Now()
		)
		err = host.attempt(func() error {
			req.URL.Scheme = host.scheme
			req.URL.Host = host.hostPort
			req.Body = getBody()

			resp, err = b.client.Do(req) //nolint:bodyclose // Is logged or returned outside the callback.
			if err != nil {
				return fmt.Errorf("doing via http client: %w", err)
			}
			if resp.StatusCode < http.StatusInternalServerError {
				b.logInfoHostAttemptOK(req.Context(), req.URL, reqStartedAt)

				return nil
			}

			return errors.New(resp.Status)
		})
		if err != nil {
			errLast = err

			if errors.Is(err, derr.ErrExceeded) {
				b.logErrorHostAttemptAttemptsExceeded(err, &host.attempts, req, resp, reqStartedAt)
			} else {
				b.logErrorSendOverHostsFailedAwaitingNextHost(err, &host.attempts, req, resp, reqStartedAt)
			}

			continue
		}

		return resp, nil
	}
}
