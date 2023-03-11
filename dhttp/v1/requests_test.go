package dhttp_test

import (
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestRequestContextWithMaxTTL(t *testing.T) {
	var ins = []struct {
		RequestContextInitialTTL time.Duration
		RequestMaxTTL            time.Duration
		RequestTTLWantedResult   time.Duration
	}{
		{
			RequestContextInitialTTL: time.Second,
			RequestMaxTTL:            time.Second,
			RequestTTLWantedResult:   time.Second,
		},
		{
			RequestContextInitialTTL: time.Second,
			RequestMaxTTL:            500 * time.Millisecond,
			RequestTTLWantedResult:   500 * time.Millisecond,
		},
		{
			RequestContextInitialTTL: 0,
			RequestMaxTTL:            500 * time.Millisecond,
			RequestTTLWantedResult:   0,
		},
	}

	for i, in := range ins {
		var in = in
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var ctx, cancel = dctx.WithTTLTimeout(dctx.New(), in.RequestContextInitialTTL)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			require.NoError(t, err)

			req, cancel2 := dhttp.RequestContextWithMaxTTL(req, in.RequestMaxTTL)
			defer cancel2()

			require.Less(t, in.RequestTTLWantedResult-150*time.Microsecond, dhttp.RequestTTL(req))
			require.Greater(t, in.RequestTTLWantedResult, dhttp.RequestTTL(req))
		})
	}
}
