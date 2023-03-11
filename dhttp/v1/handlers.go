package dhttp

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"net/http"
	"time"
)

/*
OptionHandlerWithDefaults
  - OptionHandlerWithRecover;
  - OptionHandlerWithXRequestID;
  - OptionHandlerWithGoID;
*/
func OptionHandlerWithDefaults(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			next = OptionHandlerWithRecover(next)
			next = OptionHandlerWithXRequestID(next)
			next = OptionHandlerWithGoID(next)

			next.ServeHTTP(resp, req)
		},
	)
}

// OptionHandlerWithGoID - see OptionRequestContextWithNewGoID().
func OptionHandlerWithGoID(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			req = OptionRequestContextWithNewGoID()(req)

			next.ServeHTTP(resp, req)
		},
	)
}

// OptionHandlerWithTTL - see OptionRequestContextWithTTL().
func OptionHandlerWithTTL(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(resp http.ResponseWriter, req *http.Request) {
				var cancel context.CancelFunc

				req, cancel = RequestContextWithTTL(req, d)
				defer cancel()

				next.ServeHTTP(resp, req)
			},
		)
	}
}

/*
OptionHandlerWithXRequestID - see OptionRequestContextWithXRequestID() and OptionResponseWriterHeaderWithXRequestID().
*/
func OptionHandlerWithXRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			req = OptionRequestContextWithXRequestID()(req)

			var id = dctx.XRequestID(
				req.Context(),
			)
			OptionResponseWriterHeaderWithXRequestID(resp, id)

			next.ServeHTTP(resp, req)
		},
	)
}

func OptionHandlerWithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			defer func() {
				v := recover()
				if v != nil {
					resp.WriteHeader(http.StatusInternalServerError)

					dlog.E().
						Scope(req.Context()).
						Any("panicked", true).
						Stack().
						Writef("recovered: %+v", v)
				}
			}()

			next.ServeHTTP(resp, req)
		},
	)
}
