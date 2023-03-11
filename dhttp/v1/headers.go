package dhttp

import (
	"errors"
	"fmt"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderKeyContentType              = "Content-Type"
	HeaderValueContentTypeRaw         = "application/raw-request"
	HeaderValueContentTypeJSON        = "application/json"
	HeaderValueContentTypeOctetStream = "application/octet-stream"
	HeaderKeyXRequestID               = "X-Request-Id"
	HeaderKeyContentLength            = "Content-Length"
	HeaderKeyAuthorization            = "Authorization"
)

/*
HeadersRemove - removes any header found by 'keys' from 'headers'. Returned headers is a modified or unchanged instance
pointing to the same memory location.
*/
func HeadersRemove(headers http.Header, keys ...string) http.Header {
	for _, key := range keys {
		delete(headers, key)
	}

	return headers
}

// HeadersCloneAndHideValues - returns modified 'headers' clone where each key value found by 'keys' is set to '-'.
func HeadersCloneAndHideValues(headers http.Header, keys ...string) http.Header {
	var clone = headers.Clone()

	for _, key := range keys {
		if clone.Get(key) != "" {
			clone.Set(key, "-")
		}
	}

	return clone
}

/*
HeaderBearerTokenGet - looks up respective header, validates authorization scheme and returns token value. Empty token
is considered to be invalid, in that case returned values would be "" and err != nil.
*/
func HeaderBearerTokenGet(headers http.Header) (string, error) {
	var (
		header      = headers.Get(HeaderKeyAuthorization)
		schemeToken = strings.Split(header, " ")
	)

	const partsN = 2
	if n := len(schemeToken); n < partsN {
		return "", fmt.Errorf("invalid header parts, %d != %d", n, partsN)
	}

	const scheme = "Bearer"
	if schemeToken[0] != scheme {
		return "", fmt.Errorf("invalid scheme, not %q", scheme)
	}

	token := schemeToken[1]
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}

// HeaderContentLengthGet - returns content length provided in HeaderKeyContentLength.
func HeaderContentLengthGet(headers http.Header) int {
	s := headers.Get(HeaderKeyContentLength)
	if s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			dlog.E().Stack().Writef("%s", err)
		}

		return n
	}

	return 0
}
