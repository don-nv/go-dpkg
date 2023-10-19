package dctx_test

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewID_AddGoID_AddXRequestID_GoID_XRequestID(t *testing.T) {
	var ctx = context.Background()

	var goID = uuid.NewString()
	ctx = dctx.AddGoID(ctx, goID)

	var xReqID = uuid.NewString()
	ctx = dctx.AddXRequestID(ctx, xReqID)

	require.NotEmpty(t, goID)
	require.NotEmpty(t, xReqID)
	require.NotEqualValues(t, goID, xReqID)

	require.EqualValues(t, dctx.GoID(ctx), goID)
	require.EqualValues(t, dctx.XRequestID(ctx), xReqID)
}
