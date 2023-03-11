package dhttp

import (
	"encoding/json"
	"net/http"
)

//go:generate easyjson --all $GOFILE

type clientRequestInfo struct {
	Method  string          `json:"method,omitempty"`
	Path    string          `json:"path,omitempty"`
	Headers http.Header     `json:"headers,omitempty"`
	Body    json.RawMessage `json:"body,omitempty"`
}

type clientResponseInfo struct {
	Code    int             `json:"code,omitempty"`
	Headers http.Header     `json:"headers,omitempty"`
	Body    json.RawMessage `json:"body,omitempty"`
}
