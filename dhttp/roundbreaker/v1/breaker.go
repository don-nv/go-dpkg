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
	requestsOptions    []dhttp.RequestOption
	rounds             *Rounds
	client             dhttp.IClient
	requestsDefaultTTL time.Duration
	log                dlog.Logger
	logConfig          dhttp.LoggerConfig
}

/*
MustNewBreaker - creates new Breaker. 'requestsDefaultTTL' is required and must be > 0. It is applied to any request
sent if its own ttl has not been set before Do(). This requirement allows broken or long-without deadline requests
interruption.

Logging. If dpkg.DebugEnabled(), then logs additional data at debug level. Each error, but last (which is returned) is
logged between attempts on fail. In this case, response (if exists) is logged despite any response log configuration.
Request (once before attempting) and response (once after attempts) are always logged and configurable.
*/
func MustNewBreaker(
	requestsOptions []dhttp.RequestOption,
	client dhttp.IClient,
	rounds *Rounds,
	requestsDefaultTTL time.Duration,
	log dlog.Logger,
	defaultLogConfig dhttp.LoggerConfig,
) Breaker {
	if requestsDefaultTTL < 0 {
		panic(fmt.Errorf("non-positive %q requests default ttl", requestsDefaultTTL))
	}

	return Breaker{
		requestsOptions:    requestsOptions,
		requestsDefaultTTL: requestsDefaultTTL,
		rounds:             rounds,
		client:             client,
		log:                log.With().Name("round_breaker").Build(),
		logConfig:          defaultLogConfig,
	}
}

// SetLoggerConfig - replaces default logger config and returns Breaker with modified config.
func (b Breaker) SetLoggerConfig(config dhttp.LoggerConfig) Breaker {
	b.logConfig = config

	return b
}

// Send - is the same as Do(), but creates request by itself.
func (b Breaker) Send(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return b.Do(req)
}

/*
Do - is a potentially long-running task:
  - If 'req' has no ttl (see dhttp.RequestTTL()), then default request ttl is applied before loop;
  - 'req' gets send to the next host in rounds;
  - If an error occurred during sending or response status code > 500, sending is considered to be failed;
  - If sending was failed, sending is retried until all attempts are not exceeded;
  - If all attempts were exceeded, host gets suspended for a while and 'req' gets send to the next host in rounds;
  - If sending was succeeded, host attempts gets reset;
  - For each new 'req' starting host is the next in rounds;
  - If 'req' scheme and host gets populated with current host values;
  - 'req' context controls sending duration;
  - 'req' may be reused after return, but scheme and host changes made should be taken into account;
*/
func (b Breaker) Do(req *http.Request) (*http.Response, error) {
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
		var r, cancel = dhttp.RequestContextWithTTL(req, b.requestsDefaultTTL)
		defer cancel()

		req = r
	}

	for _, option := range b.requestsOptions {
		req = option(req)
	}

	b.logDoRequest(req, body)

	var startedAt = time.Now()
	defer b.logDoDone(req.Context(), startedAt)

	resp, err := b.sendOverHosts(req, getBody)
	if err != nil {
		return nil, fmt.Errorf("sending over hosts: %w", err)
	}

	err = b.logDoResponse(req.Context(), resp)
	if err != nil {
		return nil, fmt.Errorf("logging do response: %w", err)
	}

	return resp, nil
}

func (b Breaker) sendOverHosts(req *http.Request, getBody func() io.ReadCloser) (*http.Response, error) {
	for errLast := error(nil); ; {
		var nextStartedAt = time.Now()

		host, err := b.rounds.nextHost(req.Context())
		if err != nil {
			return nil, errors.Join(errLast, fmt.Errorf("getting next host: %w", err))
		}

		b.logSendOverHostsHostChosen(req.Context(), host, nextStartedAt)

		resp, err := b.attemptHost(req, getBody, host)
		if err == nil {
			// TODO: This should cancel pending enabling goroutine and enable host if it was disabled (concurrently).
			host.Attempts.Reset()

			return resp, nil
		}

		host.disableFor()
		errLast = err
	}
}

func (b Breaker) attemptHost(req *http.Request, getBody func() io.ReadCloser, host *Host) (*http.Response, error) {
	defer func() { req.Body = getBody() }() // Recover request body for possible future sends.

	req.URL.Scheme = host.Scheme
	req.URL.Host = host.HostPort

	for host.Attempts.Next() {
		// Disabled concurrently.
		if !host.enabled() {
			return nil, derr.ErrExceeded
		}

		req.Body = getBody()

		// Cannot be moved in outer scope, because next attempt awaiting may last long.
		var startedAt = time.Now()

		resp, err := b.client.Do(req)
		if isSendNoErrLess500(err, resp) {
			b.logAttemptHostOK(req.Context(), req.URL, startedAt)

			return resp, nil
		}

		b.logAttemptHostAwaitingRetry(err, &host.Attempts, req, resp, startedAt)

		aerr := host.Attempts.AwaitDelay(req.Context())
		if aerr != nil {
			return nil, errors.Join(err, aerr)
		}
	}

	// Disabled concurrently.
	if !host.enabled() {
		return nil, derr.ErrExceeded
	}

	req.Body = getBody()

	// Last attempt || at least one attempt if all previous attempts were failed (on host re-enabling).
	var startedAt = time.Now()

	resp, err := b.client.Do(req)
	if isSendNoErrLess500(err, resp) {
		b.logAttemptHostOK(req.Context(), req.URL, startedAt)

		return resp, nil
	}

	b.logAttemptHostAttemptsExceeded(err, &host.Attempts, req, resp, startedAt)

	return nil, derr.ErrExceeded
}

func isSendNoErrLess500(err error, resp *http.Response) bool {
	return err == nil && resp.StatusCode < http.StatusInternalServerError
}
