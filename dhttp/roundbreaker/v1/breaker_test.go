package roundbreaker_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dhttp/roundbreaker/v1"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

//nolint:funlen
func TestBreaker_Sending(t *testing.T) {
	var (
		reqsN      int
		newHandler = func(response string) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.EqualValues(t, "?", string(body))

				if reqsN < 9 {
					reqsN++
					rw.WriteHeader(http.StatusInternalServerError)
				}

				d, err := json.Marshal(map[string]any{
					"key": response,
				})
				require.NoError(t, err)

				if len(response) > 0 {
					_, err = rw.Write(d)
					require.NoError(t, err)
				}

				t.Log("served: " + response)
			}
		}
	)

	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		srv1 = httptest.NewServer(newHandler("1"))
		srv2 = httptest.NewServer(newHandler("2"))
		srv3 = httptest.NewServer(newHandler(""))
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
		attemptsDelays = []time.Duration{
			500 * time.Millisecond,
			500 * time.Millisecond,
		}
		hostDisabledFor = time.Millisecond

		rounds = roundbreaker.NewRounds(
			[]roundbreaker.Host{
				{
					Scheme:   url1.Scheme,
					HostPort: url1.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
				{
					Scheme:   url2.Scheme,
					HostPort: url2.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
				{
					Scheme:   url3.Scheme,
					HostPort: url3.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
			},
		)
	)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, "/some/path?and=some&query=here#too", bytes.NewReader([]byte("?")),
	)
	require.NoError(t, err)

	var breaker = roundbreaker.MustNewBreaker(
		nil, http.DefaultClient, rounds, 3500*time.Millisecond, dlog.New(), dhttp.LoggerConfig{},
	)

	resp1, err := breaker.Do(req)
	require.NoError(t, err)
	defer resp1.Body.Close()

	body, err := io.ReadAll(resp1.Body)
	require.NoError(t, err)
	require.EqualValues(t, "{\"key\":\"1\"}", string(body))

	resp2, err := breaker.Do(req)
	require.NoError(t, err)
	defer resp2.Body.Close()

	body, err = io.ReadAll(resp2.Body)
	require.NoError(t, err)
	require.EqualValues(t, "{\"key\":\"2\"}", string(body))
}

//nolint:funlen
func TestBreaker_SendingConcurrency(t *testing.T) {
	const (
		responseOKAfterNRequests = 1000
		path                     = "/some/path?and=some&query=here"
	)

	var (
		reqsN      atomic.Int32
		newHandler = func(serverName string) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				require.EqualValues(t, path, r.URL.Path+"?"+r.URL.Query().Encode())

				for k, vs := range r.Header {
					require.EqualValuesf(t, 1, len(vs), "%q %v", k, vs)
				}

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.EqualValues(t, "?", string(body))

				if reqsN.Load() < responseOKAfterNRequests {
					reqsN.Add(1)

					rw.WriteHeader(http.StatusInternalServerError)
				}

				if len(serverName) > 0 {
					_, err = rw.Write([]byte(serverName))
					require.NoError(t, err)
				}

				t.Log("served: " + serverName)
			}
		}
	)

	var (
		ctx = dctx.New(dctx.OptionWithOSCancel())

		srv1 = httptest.NewServer(newHandler("1"))
		srv2 = httptest.NewServer(newHandler("2"))
		srv3 = httptest.NewServer(newHandler("3"))
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
		attemptsDelays = []time.Duration{
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
		}
		hostDisabledFor = 50 * time.Millisecond

		rounds = roundbreaker.NewRounds(
			[]roundbreaker.Host{
				{
					Scheme:   url1.Scheme,
					HostPort: url1.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
				{
					Scheme:   url2.Scheme,
					HostPort: url2.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
				{
					Scheme:   url3.Scheme,
					HostPort: url3.Host,
					Attempts: dstruct.NewAttemptsV1Sync(dstruct.AttemptsV1{
						Delays: attemptsDelays,
					}),
					DisabledFor: hostDisabledFor,
				},
			},
		)
	)

	var (
		breaker = roundbreaker.MustNewBreaker(
			[]dhttp.RequestOption{
				dhttp.OptionRequestHeaderWithAuthorization("Bearer jwt"),
				dhttp.OptionRequestHeaderWithContentType(dhttp.HeaderValueContentTypeJSON),
				dhttp.OptionRequestHeaderWithXRequestID(),
			},
			dhttp.MustNewClient(
				dhttp.ClientConfig{
					RequestsDefaultTTL: 5 * time.Millisecond,
					RequestsOptions: []dhttp.RequestOption{
						dhttp.OptionRequestHeaderWithAuthorization("Bearer jwt"),
						dhttp.OptionRequestHeaderWithContentType(dhttp.HeaderValueContentTypeJSON),
						dhttp.OptionRequestHeaderWithXRequestID(),
					},
				},
				dlog.New(),
			),
			rounds,
			1500*time.Millisecond,
			dlog.New(),
			dhttp.LoggerConfig{
				RequestBodyOmitted: true,
			},
		)
		wg = dsync.NewGroup(ctx)
	)

	for i := 0; i < responseOKAfterNRequests*2; i++ {
		wg.Go(func(ctx context.Context) error {
			resp, err := breaker.Send(dctx.WithNewXRequestID(ctx), http.MethodPost, path, []byte("?"))
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
}
