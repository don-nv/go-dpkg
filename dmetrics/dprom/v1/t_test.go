package dprom_test

import (
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dmetrics/dprom/v1"
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var ctx, cancel = dctx.NewTTL(dctx.OptionTTLWithOSCancel())

func Test(t *testing.T) {
	t.SkipNow()

	go func() {
		<-ctx.Done()
		cancel()

		os.Exit(0)
	}()

	go func() {
		var v1 = dprom.NewEntry.
			WithEntityGroup("api").
			WithEntity("http").
			WithMethodGroup("v1")

		for ctx.Err() == nil {
			func() {
				var result string

				defer func() {
					v1.
						WithMethod("get_payment").
						WithResult(result).
						HistogramMs().
						Write()
				}()

				time.Sleep(time.Duration(rand.IntN(2000)) * time.Millisecond)

				result = strconv.Itoa(rand.IntN(10))
			}()
		}
	}()

	//nolint:gosec,G114
	err := http.ListenAndServe(
		"0.0.0.0:8081",
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			t.Log("metrics requested")
			dprom.ServeHTTP(writer, request)
		}),
	)
	require.NoError(t, err)
}
