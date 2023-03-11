package dhttp

import (
	"bytes"
	"fmt"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"io"
	"net/http"
)

type LoggerConfig struct {
	// RequestHeadersHidden - each request header value found is replaced with dlog.HiddenValueString.
	RequestHeadersHidden []string
	// RequestSkip - skips request logging entirely.
	RequestSkip bool
	// RequestBodyOmitted - logs request body as OmittedValue.
	RequestBodyOmitted bool
	// ResponseSkip - skips response logging entirely.
	ResponseSkip bool
	// ResponseBodyOmitted - logs response body as OmittedValue.
	ResponseBodyOmitted bool
}

/*
LogWithClientRequest - adds request data to 'data' according to 'config'. HeaderKeyAuthorization is omitted by default
(see HeadersCloneAndHideValues()).
*/
func LogWithClientRequest(req *http.Request, body []byte, data dlog.Data, config LoggerConfig) dlog.Data {
	if config.RequestSkip {
		return data
	}

	var infoBody = []byte{dlog.HiddenValueByte}

	if !config.RequestBodyOmitted {
		infoBody = body
	}

	var info = clientRequestInfo{
		Method: req.Method,
		Path:   req.URL.Path,
		Headers: HeadersCloneAndHideValues(
			req.Header, append(config.RequestHeadersHidden, HeaderKeyAuthorization)...,
		),
		Body: infoBody,
	}

	reqJSON, err := info.MarshalJSON()
	if err != nil {
		dlog.E().Stack().Writef("%s", err)
		// 'info' bytes escapes here only, because 'Any' accepts 'any' parameter.
		return data.Any("request", info)
	}

	return data.Bytes("request", reqJSON)
}

// LogWithClientResponse - preserves response body.
func LogWithClientResponse(resp *http.Response, data dlog.Data, config LoggerConfig) (dlog.Data, error) {
	if config.ResponseSkip {
		return data, nil
	}

	var infoBody = []byte{dlog.HiddenValueByte}

	if !config.ResponseBodyOmitted {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return dlog.Data{}, fmt.Errorf("reading response body: %w", err)
		}
		_ = resp.Body.Close()

		resp.Body = io.NopCloser(bytes.NewReader(b))

		infoBody = b
	}

	var info = clientResponseInfo{
		Code:    resp.StatusCode,
		Headers: resp.Header,
		Body:    infoBody,
	}

	respJSON, err := info.MarshalJSON()
	if err != nil {
		dlog.E().Stack().Writef("%s", err)
		// 'info' escapes here only, because 'Any' accepts 'any'.
		return data.Any("response", info), nil
	}

	return data.Bytes("response", respJSON), nil
}
