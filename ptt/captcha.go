package ptt

import (
	"fmt"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/utils"
	"github.com/gin-gonic/gin"
)

type CaptchaInsertParams struct {
	Secret string `json:"secret"`
	Handle string `json:"handle"`
	Verify string `json:"verify"`
}

type CaptchaResult struct{}

func AttemptRemoteCaptcha(c *gin.Context) (url string, verify string, err error) {
	handle, err := cmsys.RandomTextCode()
	if err != nil {
		return "", "", err
	}
	verify, err = cmsys.RandomTextCode()
	if err != nil {
		return "", "", err
	}

	url, err = captchaInsertRemote(handle, verify, c)
	if err != nil {
		return "", "", err
	}

	return url, verify, nil
}

func captchaInsertRemote(handle string, verify string, c *gin.Context) (url string, err error) {
	params := &CaptchaInsertParams{
		Secret: ptttype.CAPTCHA_INSERT_SECRET,
		Handle: handle,
		Verify: verify,
	}
	var result_b *CaptchaResult
	statusCode, err := utils.BackendGet(c, ptttype.CAPTCHA_INSERT_URI, params, nil, &result_b)
	if err != nil {
		return "", err
	}
	if statusCode != 200 {
		return "", ErrCaptcha
	}

	url = fmt.Sprintf("%v?handle=%v", ptttype.CAPTCHA_URL_PREFIX, handle)

	return url, nil
}
