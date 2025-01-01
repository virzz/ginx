package ginx

import (
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func handleUpgrade(c *gin.Context) {
	if v := c.Param("token"); v != Conf.Upgrade {
		c.AbortWithStatus(404)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	newFile := exePath + "_upgrade"
	err = c.SaveUploadedFile(file, newFile)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = os.Chmod(newFile, 0o755)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = os.Rename(newFile, exePath)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	buf, err := exec.Command(exePath, "restart").CombinedOutput()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.String(200, string(buf))
}
