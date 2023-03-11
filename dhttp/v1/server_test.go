package dhttp_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/don-nv/go-dpkg/dmath/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestServer_Running(t *testing.T) {
	var ctx, cancel = dctx.NewTTL(dctx.OptionTTLWithOSCancel(), dctx.OptionTTLWithTimeout(3*time.Second))
	defer cancel()

	var group = dsync.NewOneTimeGroup(ctx)
	defer func() { require.NoError(t, group.Wait()) }()

	group.GoUntilWait("test_http_server_running", func(ctx context.Context) error {
		return dhttp.NewServer(
			dhttp.ServerConfig{
				Address:              ":8082",
				RequestReadHeaderTTL: time.Second,
				RequestReadTTL:       time.Second,
			},
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				if dmath.Is50x50() {
					time.Sleep(time.Second)
				}

				defer dlog.New().D().Scope(r.Context()).Write("served")

				if err := r.Context().Err(); err != nil {
					_, err := rw.Write([]byte(err.Error()))
					derr.PanicOnE(err)

					return
				}

				_, err := rw.Write([]byte("ok"))
				derr.PanicOnE(err)
			}),
			dlog.New(),
			dhttp.OptionServerWithMiddleware(
				dhttp.OptionHandlerWithDefaults,
				dhttp.OptionHandlerWithTTL(time.Millisecond),
			),
		).
			Running(ctx)
	})
}
