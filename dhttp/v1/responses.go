package dhttp

import (
	"bytes"
	"fmt"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/don-nv/go-dpkg/dmath/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"net/http"
)

type ResponseBuffer struct {
	Headers    http.Header
	buffer     *bytes.Buffer
	StatusCode int
}

// NewClientResponseBuffer - closes response body.
func NewClientResponseBuffer(response *http.Response) (ResponseBuffer, error) {
	n, ok := dmath.Cast[int64, int](response.ContentLength)
	if n < 0 /* 0, because body may be empty. */ || !ok {
		n = dsync.PoolBufferSizeMax
	}

	var buff = dsync.PoolBufferGet(n)

	_, err := buff.ReadFrom(response.Body)
	if err != nil {
		return ResponseBuffer{}, fmt.Errorf("buff.ReadFrom: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return ResponseBuffer{}, fmt.Errorf("response.Body.Close: %w", err)
	}

	var respBuff = ResponseBuffer{
		Headers:    response.Header,
		buffer:     buff,
		StatusCode: response.StatusCode,
	}

	return respBuff, nil
}

func (r *ResponseBuffer) Body() []byte {
	if r.Released() {
		dlog.E().Stack().Write("released")

		return nil
	}

	return r.buffer.Bytes()
}

func (r *ResponseBuffer) Release() {
	if r.Released() {
		return
	}

	dsync.PoolBufferPut(r.buffer)
	r.buffer = nil
}

func (r *ResponseBuffer) Released() bool {
	return r.buffer == nil
}

type ResponseError struct {
	Code    CodeError `json:"code"`
	Message string    `json:"message,omitempty"`
}

func OptionResponseWriterHeaderWithXRequestID(resp http.ResponseWriter, id string) {
	resp.Header().Set(HeaderKeyXRequestID, id)
}
