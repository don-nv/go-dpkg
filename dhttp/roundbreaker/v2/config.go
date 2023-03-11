package roundbreaker

import (
	"errors"
	"github.com/don-nv/go-dpkg/dhttp/v1"
	"time"
)

type Config struct {
	RequestsOptions []dhttp.RequestOption
	// RequestsDefaultTTL (required) - is applied to request context if it has no TTL.
	RequestsDefaultTTL time.Duration
	/*
		Logger - configures request (once before sending) and response (once after sending) logging. May log additional
		data at debug level if dpkg.DebugEnabled(). Each error, but last, is logged if host attempt fails, then response
		(if exists) is logged despite any response log configuration.
	*/
	Logger dhttp.LoggerConfig
}

func (c Config) validate() error {
	if c.RequestsDefaultTTL < 1 {
		return errors.New("requests default ttl < 1")
	}

	return nil
}
