package roundbreaker

import (
	"context"
	"fmt"
	"github.com/don-nv/go-dpkg"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"net/http"
	"net/url"
	"time"
)

func (b Breaker) logInfoDoRequest(req *http.Request, body []byte) {
	dhttp.LogWithClientRequest(req, body, b.log.With().Scope(req.Context()), b.config.Logger).
		Build().
		I().
		Write("(request)")
}

// logInfoDoResponse - preserves body.
func (b Breaker) logInfoDoResponse(ctx context.Context, resp *http.Response) error {
	with, err := dhttp.LogWithClientResponse(resp, b.log.With(), b.config.Logger)
	if err != nil {
		return fmt.Errorf("creating log with client reponse: %w", err)
	}

	with.Build().
		I().
		Scope(ctx).
		Write("(response)")

	return nil
}

func (b Breaker) logDebugDoDone(ctx context.Context, startedAt time.Time) {
	if !dpkg.DebugEnabled() {
		return
	}

	b.log.
		With().
		Scope(ctx).
		String("duration", time.Since(startedAt).String()).
		Build().
		D().
		Write("done")
}

// logErrorSendOverHostsFailedAwaitingNextHost - see logErrorResponse.
func (b Breaker) logErrorSendOverHostsFailedAwaitingNextHost(
	err error,
	attempts *dstruct.AttemptsSync,
	req *http.Request,
	resp *http.Response,
	startedAt time.Time,
) {
	var (
		attemptN  = attempts.AttemptN()
		attemptsN = attempts.AttemptsN()
	)

	b.logErrorResponse(err, req, resp, startedAt, attemptN, attemptsN, "failed: awaiting next host...")
}

// logErrorHostAttemptAttemptsExceeded - see logErrorResponse.
func (b Breaker) logErrorHostAttemptAttemptsExceeded(
	err error,
	attempts *dstruct.AttemptsSync,
	req *http.Request,
	resp *http.Response,
	startedAt time.Time,
) {
	var (
		attemptN  = attempts.AttemptN()
		attemptsN = attempts.AttemptsN()
	)

	if attemptN == attemptsN {
		b.logErrorResponse(err, req, resp, startedAt, attemptN, attemptsN, "host attempts exceeded")

		return
	}

	b.logErrorResponse(
		err, req, resp, startedAt, attemptN, attemptsN, "host attempts exceeded (concurrently)",
	)
}

/*
logErrorResponse - logs message at dlog.LevelError. 'resp' body if not nil, gets closed. 'err' and 'resp' may be nil. In
this case corresponding log fields are omitted.
*/
func (b Breaker) logErrorResponse(
	err error,
	req *http.Request,
	resp *http.Response,
	startedAt time.Time,
	attemptN int,
	attemptsN int,
	msg string,
) {
	var with = b.log.With()

	if err != nil {
		with = with.Any("error", err)
	}

	if resp != nil {
		defer resp.Body.Close()

		withResp, err := dhttp.LogWithClientResponse(resp, with, dhttp.LoggerConfig{})
		if err != nil {
			with.Error("response", fmt.Errorf("creating log with client respone: %w", err))
		} else {
			with = withResp
		}
	}

	with.Scope(req.Context()).
		Any("url", req.URL.String()).
		Any("duration", time.Since(startedAt).String()).
		Any("attempt_n", attemptN).
		Any("attempts_n", attemptsN).
		Build().
		E().
		Write(msg)
}

func (b Breaker) logInfoHostAttemptOK(ctx context.Context, url *url.URL, startedAt time.Time) {
	b.log.I().
		Scope(ctx).
		Any("url", url.String()).
		Any("duration", time.Since(startedAt).String()).
		Write("ok")
}

func (b Breaker) logDebugSendOverHostsHostChosen(ctx context.Context, url *Host, awaitStartedAt time.Time) {
	if !dpkg.DebugEnabled() {
		return // Logger implementation adds data to logging context despite debug level state.
	}

	b.log.D().
		Scope(ctx).
		Any("host", url.hostPort).
		Any("duration", time.Since(awaitStartedAt).String()).
		Write("chosen")
}
