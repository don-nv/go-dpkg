package dgin

import (
	"github.com/don-nv/go-dpkg/dctx/v1"
	dhttp "github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/gin-gonic/gin"
)

/*
OptionHandlerWithDefaults

Default preset:
  - dhttp.OptionRequestContextWithNewGoID;
  - dhttp.OptionRequestContextWithXRequestID;
  - dhttp.OptionResponseWriterHeaderWithXRequestID;
*/
func OptionHandlerWithDefaults(c *gin.Context) {
	var req = c.Request

	req = dhttp.OptionRequestContextWithNewGoID()(req)
	req = dhttp.OptionRequestContextWithXRequestID()(req)

	var id = dctx.XRequestID(req.Context())
	dhttp.OptionResponseWriterHeaderWithXRequestID(c.Writer, id)

	c.Request = req

	c.Next()
}
