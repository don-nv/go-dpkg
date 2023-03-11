package dgin

import (
	dhttp "github.com/don-nv/go-dpkg/dhttp/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func R200Empty(c *gin.Context) {
	R200(c, nil)
}

func R200(c *gin.Context, v any) {
	c.JSON(http.StatusOK, v)
}

func R202Empty(c *gin.Context) {
	R202(c, nil)
}

func R202(c *gin.Context, v any) {
	c.JSON(http.StatusAccepted, v)
}

func R400(c *gin.Context, code dhttp.CodeError, err error) {
	Abort(c, http.StatusBadRequest, code, err)
}

func R403(c *gin.Context, code dhttp.CodeError, err error) {
	Abort(c, http.StatusForbidden, code, err)
}

func R404(c *gin.Context, code dhttp.CodeError, err error) {
	Abort(c, http.StatusNotFound, code, err)
}

func R500(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, dhttp.Code500General, nil)
}

func Abort(c *gin.Context, httpCode int, codeError dhttp.CodeError, err error) {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	resp := dhttp.ResponseError{
		Code:    codeError,
		Message: msg,
	}

	c.AbortWithStatusJSON(httpCode, resp)
}
