package dhttp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"io"
	"net/http"
	"time"
)

/*
RequestTTL - returns 'req' TTL. Currently, TTL is considered to be found in requests' context only.
  - < 0: expired
  - = 0: no TTL
  - > 0: not expired
*/
func RequestTTL(req *http.Request) time.Duration {
	t, ok := req.Context().Deadline()
	if !ok {
		return 0
	}

	d := time.Until(t)
	if d == 0 {
		return -1
	}

	return d
}

type RequestOption func(req *http.Request) *http.Request

func OptionRequest(req *http.Request, options ...RequestOption) *http.Request {
	for _, option := range options {
		req = option(req)
	}

	return req
}

// OptionRequestHeaderWith - appends 'value' to the 'key' header. 'key' gets canonically formatted.
func OptionRequestHeaderWith(key, value string) RequestOption {
	return func(req *http.Request) *http.Request {
		req.Header.Add(key, value)

		return req
	}
}

func OptionRequestHeaderWithAuthorization(schemeToken string) RequestOption {
	return func(req *http.Request) *http.Request {
		return OptionRequestHeaderWith(HeaderKeyAuthorization, schemeToken)(req)
	}
}

func OptionRequestHeaderWithContentType(typ string) RequestOption {
	return func(req *http.Request) *http.Request {
		return OptionRequestHeaderWith(HeaderKeyContentType, typ)(req)
	}
}

/*
OptionRequestHeaderWithXRequestID - is used to populate request respective header with dctx.XRequestID(). If id wasn't
set before, gets populated with '-' value instead as an indicator of incomplete result.
*/
func OptionRequestHeaderWithXRequestID() RequestOption {
	return func(req *http.Request) *http.Request {
		var ctx = req.Context()

		id := dctx.XRequestID(ctx)
		if id == "" {
			id = dlog.HiddenValueString
		}

		return OptionRequestHeaderWith(HeaderKeyXRequestID, id)(req)
	}
}

/*
OptionRequestContextWithNewGoID - is used to populate request context with dctx.WithNewGoID(). Returned request is a
shallow copy of 'req'.
*/
func OptionRequestContextWithNewGoID() RequestOption {
	return func(req *http.Request) *http.Request {
		var ctx = dctx.WithNewGoID(
			req.Context(),
		)

		return req.WithContext(ctx)
	}
}

/*
OptionRequestContextWithXRequestID - is used to populate request context with dctx.WithXRequestID() from respective
header. If id is not present, context gets populated with a new one using dctx.WithNewXRequestID(). Returned request is
a shallow copy of 'req'.
*/
func OptionRequestContextWithXRequestID() RequestOption {
	return func(req *http.Request) *http.Request {
		var ctx = req.Context()

		id := req.Header.Get(HeaderKeyXRequestID)
		if id != "" {
			ctx = dctx.WithXRequestID(ctx, id)
		} else {
			ctx = dctx.WithNewXRequestID(ctx)
		}

		return req.WithContext(ctx)
	}
}

// RequestBodyAppendAndKeep - reads 'req' body into 'buffer' keeping 'req' body io.ReadCloser unread.
func RequestBodyAppendAndKeep(req *http.Request, body []byte) ([]byte, error) {
	var buff = bytes.NewBuffer(body)

	_, err := io.Copy(buff, req.Body)
	if err != nil {
		return nil, fmt.Errorf("io.Copy: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing read request body: %w", err)
	}

	req.Body = io.NopCloser(buff)

	return buff.Bytes(), nil
}

// RequestContextWithTTL - sets request context ttl == [d]. Returned request is a shallow copy.
func RequestContextWithTTL(req *http.Request, d time.Duration) (*http.Request, context.CancelFunc) {
	var ctx, cancel = dctx.WithTTLTimeout(req.Context(), d)

	return req.WithContext(ctx), cancel
}

/*
RequestContextWithMaxTTL - if [RequestTTL] returns no ttl, then sets [d] as context ttl via [RequestContextWithTTL].
Returned context cancel func is either context cancelling or no-op.
*/
func RequestContextWithMaxTTL(req *http.Request, d time.Duration) (*http.Request, context.CancelFunc) {
	ttl := RequestTTL(req)
	if ttl > 0 && ttl > d {
		return RequestContextWithTTL(req, d)
	}

	return req, func() {}
}
