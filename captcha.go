package ginx

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/captcha"
	"github.com/virzz/vlog"

	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/req"
	"github.com/virzz/ginx/rsp"
)

func init() {
	captcha.NewCaptchaEquation(1)
}

func CaptchaHandler(c *gin.Context) {
	id, _code, data, err := captcha.CreateB64()
	if err != nil {
		vlog.Error("Failed to create base64 captcha", "err", err.Error())
		c.AbortWithStatusJSON(200, rsp.C(code.CaptchaGenerate))
		return
	}
	if c.GetHeader("X-Debug-Captcha") != "" {
		c.Header("Captcha", _code)
	}
	c.JSON(200, rsp.S(req.Captcha{UUID: id, Code: data}))
}

func CaptchaCheck(uuid, code string) bool {
	if Conf != nil && Conf.Captcha {
		return true
	}
	return captcha.CheckOk(uuid, code)
}
