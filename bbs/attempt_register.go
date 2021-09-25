package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/gin-gonic/gin"
)

func AttemptRegister(c *gin.Context) (url string, verify string, err error) {
	return ptt.AttemptRemoteCaptcha(c)
}
