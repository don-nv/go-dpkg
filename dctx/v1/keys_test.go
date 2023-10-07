package dctx_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewID_WithGoID_WithXRequestID_GoID_XRequestID(t *testing.T) {
	var ctx = context.Background()

	var goID = dctx.NewID()
	ctx = dctx.WithGoID(ctx, goID)

	var xReqID = dctx.NewID()
	ctx = dctx.WithXRequestID(ctx, xReqID)

	require.NotEmpty(t, goID)
	require.NotEmpty(t, xReqID)
	require.NotEqualValues(t, goID, xReqID)

	require.EqualValues(t, dctx.GoID(ctx), goID)
	require.EqualValues(t, dctx.XRequestID(ctx), xReqID)
}
