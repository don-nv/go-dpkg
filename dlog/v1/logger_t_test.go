package dlog_test

import (
	"errors"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

var ctx = dctx.New(dctx.OptionWithNewXRequestID())

func TestCommon(t *testing.T) {
	t.Log("ReadScope(): enabled")
	Log(dlog.New())

	t.Log("ReadScope(): disabled")
	Log(dlog.New(dlog.OptionLoggerWithReadScopeDisabled()))
}

func Log(l dlog.Logger) {
	var err error

	l = l.With().Name("a", "b").Build().WithDuration()
	defer l.CatchED(&err)

	l.E().Write("msg")

	l = l.With().Name("c").Build()
	l.W().Write("msg")

	l.I().Name("d", "e").Write("msg")
	l.D().Name("d", "e").Any("key", "value").Write("msg")

	l = l.With().Name("d", "e").Any("key", "value").Build()

	l.E().Scope(ctx).Write("msg")
	l = l.With().Scope(ctx).Build()
	l.W().Writef("%s", "msg")

	err = errors.New("error")
}

func TestLogger_ObjectMarshallerJSON(t *testing.T) {
	var consumer = NewConsumerLogger()

	err := consumer.LogObjectMarshallerJSON(ctx)
	require.Error(t, err)
}
