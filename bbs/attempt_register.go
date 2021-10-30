package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func AttemptRegister() (url string, verify string, err error) {
	return ptt.AttemptRemoteCaptcha()
}
