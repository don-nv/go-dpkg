package dhttp_test

import (
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/derr/v1"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkClient_POST(b *testing.B) {
	const requestsTimeout = 5 * time.Millisecond

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"code":"OK"}`))
		derr.PanicOnE(err)
	}))

	client := dhttp.MustNewClient(
		dhttp.ClientConfig{
			RequestsDefaultTTL: requestsTimeout,
		},
		dlog.New(),
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.POST(dctx.New(), srv.URL, []byte(`{"key": "value"}`))
		derr.PanicOnE(err)
		defer resp.Body.Close()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.POST(dctx.New(), srv.URL, []byte(`{"key": "value"}`))
		derr.PanicOnE(err)
		defer resp.Body.Close()
	}
}
