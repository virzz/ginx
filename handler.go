package ginx

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(code.Codes)) }
