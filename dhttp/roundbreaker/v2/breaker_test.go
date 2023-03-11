package roundbreaker_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dhttp/roundbreaker/v2"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

//nolint:funlen
func TestBreaker_Sending(t *testing.T) {
	const (
		requestPath = "/some/path?and=some&query=here"
		// Regulates first N responses having 500 status code.
		responsesWith500 = 13
	)
	var (
		requestBody = []byte("?")
		responsesN  int
		// Contains order in which servers were called.
		serversCallOrder []int

		newHandler = func(name string, serverN int) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				defer t.Logf("served: %q", name)

				require.EqualValues(t, requestPath, r.URL.Path+"?"+r.URL.Query().Encode())

				responsesN++
				serversCallOrder = append(serversCallOrder, serverN)

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.EqualValues(t, string(requestBody), string(body))

				if responsesN < responsesWith500 {
					rw.WriteHeader(http.StatusInternalServerError)
				}

				_, err = rw.Write([]byte(fmt.Sprintf("{\"key\":%d}", serverN)))
				require.NoError(t, err)
			}
		}
	)

	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		srv1 = httptest.NewServer(nil)
		srv2 = httptest.NewServer(nil)
		srv3 = httptest.NewServer(nil)
	)
	defer srv1.Close()
	defer srv2.Close()
	defer srv3.Close()

	srv1.Config.Handler = newHandler(srv1.URL, 1)
	srv2.Config.Handler = newHandler(srv2.URL, 2)
	srv3.Config.Handler = newHandler(srv3.URL, 3)

	url1, err := url.Parse(srv1.URL)
	require.NoError(t, err)
	url2, err := url.Parse(srv2.URL)
	require.NoError(t, err)
	url3, err := url.Parse(srv3.URL)
	require.NoError(t, err)

	var (
		disabledFor = 50 * time.Millisecond
		rounds      = roundbreaker.NewRounds(
			(&roundbreaker.HostsEnvConfigsYAML{
				{
					Scheme:   url1.Scheme,
					HostPort: url1.Host,
					Attempts: dstruct.Attempts{
						MaxN: 3,
					},
					DisabledFor: disabledFor,
				},
				{
					Scheme:   url2.Scheme,
					HostPort: url2.Host,
					Attempts: dstruct.Attempts{
						MaxN: 3,
					},
					DisabledFor: disabledFor,
				},
				{
					Scheme:   url3.Scheme,
					HostPort: url3.Host,
					Attempts: dstruct.Attempts{
						MaxN: 3,
					},
					DisabledFor: disabledFor,
				},
			}).Hosts(),
		)
	)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, requestPath, bytes.NewReader(requestBody),
	)
	require.NoError(t, err)

	var (
		startedAt = time.Now()
		// A bit bigger than estimated test execution.
		defaultTTL = 70 * time.Millisecond
		breaker    = roundbreaker.MustNewBreaker(
			http.DefaultClient,
			rounds,
			roundbreaker.Config{
				RequestsDefaultTTL: defaultTTL,
			},
			dlog.New(),
		)
	)

	// Receives 500 firs, all servers get disabled, first enables and responses with 200.
	resp1, err := breaker.Do(req)
	require.NoError(t, err)
	defer resp1.Body.Close()

	body, err := io.ReadAll(resp1.Body)
	require.NoError(t, err)
	require.EqualValues(t, "{\"key\":1}", string(body))

	// Next server in a sequence (2) responses with 200 (was enabled after disabledFor).
	resp2, err := breaker.Do(req)
	require.NoError(t, err)
	defer resp2.Body.Close()

	body, err = io.ReadAll(resp2.Body)
	require.NoError(t, err)
	require.EqualValues(t, "{\"key\":2}", string(body))

	require.EqualValues(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2}, serversCallOrder)
	// Check if host enabling awaiting was.
	require.Greater(t, time.Since(startedAt), disabledFor)
	require.Less(t, time.Since(startedAt), defaultTTL)
}

//nolint:funlen
func TestBreaker_SendingConcurrency(t *testing.T) {
	const (
		responsesN  = 1500
		requestPath = "/some/path?and=some&query=here"
	)

	var (
		requestBody = []byte("?")

		newHandler = func(serverN int, responseBody string) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				defer t.Logf("served: %d", serverN)

				require.EqualValues(t, requestPath, r.URL.Path+"?"+r.URL.Query().Encode())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.EqualValues(t, string(requestBody), string(body))

				if time.Now().UnixNano()%5 == 1 {
					rw.WriteHeader(http.StatusInternalServerError)
				}

				_, err = rw.Write([]byte(responseBody))
				require.NoError(t, err)
			}
		}
	)

	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		srv1 = httptest.NewServer(newHandler(1, "1"))
		srv2 = httptest.NewServer(newHandler(2, "2"))
		srv3 = httptest.NewServer(newHandler(3, "3"))
	)
	defer srv1.Close()
	defer srv2.Close()
	defer srv3.Close()

	url1, err := url.Parse(srv1.URL)
	require.NoError(t, err)
	url2, err := url.Parse(srv2.URL)
	require.NoError(t, err)
	url3, err := url.Parse(srv3.URL)
	require.NoError(t, err)

	var (
		disabledFor = 50 * time.Millisecond
		rounds      = roundbreaker.NewRounds(
			(&roundbreaker.HostsEnvConfigsYAML{
				{
					Scheme:   url1.Scheme,
					HostPort: url1.Host,
					Attempts: dstruct.Attempts{
						MaxN: 1,
					},
					DisabledFor: disabledFor,
				},
				{
					Scheme:   url2.Scheme,
					HostPort: url2.Host,
					Attempts: dstruct.Attempts{
						MaxN: 5,
					},
					DisabledFor: disabledFor,
				},
				{
					Scheme:   url3.Scheme,
					HostPort: url3.Host,
					Attempts: dstruct.Attempts{
						MaxN: 10,
					},
					DisabledFor: disabledFor,
				},
			}).Hosts(),
		)
	)

	var (
		startedAt = time.Now()
		breaker   = roundbreaker.MustNewBreaker(
			dhttp.MustNewClient(dhttp.ClientConfig{}, dlog.New()),
			rounds,
			roundbreaker.Config{
				RequestsOptions: []dhttp.RequestOption{
					dhttp.OptionRequestHeaderWithAuthorization("Bearer jwt"),
					dhttp.OptionRequestHeaderWithContentType(dhttp.HeaderValueContentTypeJSON),
					dhttp.OptionRequestHeaderWithXRequestID(),
				},
				// Estimated max execution time of the test.
				RequestsDefaultTTL: 1000 * time.Millisecond,
				Logger:             dhttp.LoggerConfig{},
			},
			dlog.New(),
		)
		wg = dsync.NewGroup(ctx)
	)

	for i := 0; i < responsesN; i++ {
		wg.Go(func(ctx context.Context) error {
			resp, err := breaker.Send(dctx.WithNewXRequestID(ctx), http.MethodPost, requestPath, requestBody)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.NotZero(t, string(body))

			return nil
		})
	}

	err = wg.Wait()
	require.NoError(t, err)
	require.Greater(t, time.Since(startedAt), disabledFor)
}
