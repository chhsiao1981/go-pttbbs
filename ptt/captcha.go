package ptt

import (
	"fmt"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/utils"
)

type CaptchaInsertParams struct {
	Secret string `json:"secret"`
	Handle string `json:"handle"`
	Verify string `json:"verify"`
}

type CaptchaResult struct{}

//AttemptRemoteCaptcha
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/captcha.c#L181
//Let middleware handle the verify. (requiring middleware store session-cookie and verify in redis)
func AttemptRemoteCaptcha() (url string, verify string, err error) {
	handle, err := cmsys.RandomTextCode()
	if err != nil {
		return "", "", err
	}
	verify, err = cmsys.RandomTextCode()
	if err != nil {
		return "", "", err
	}

	url, err = captchaInsertRemote(handle, verify)
	if err != nil {
		return "", "", err
	}

	return url, verify, nil
}

func captchaInsertRemote(handle string, verify string) (url string, err error) {
	params := &CaptchaInsertParams{
		Secret: ptttype.CAPTCHA_INSERT_SECRET,
		Handle: handle,
		Verify: verify,
	}
	var result_b *CaptchaResult
	statusCode, err := utils.BackendGet(ptttype.CAPTCHA_INSERT_URI, params, nil, &result_b)
	if err != nil {
		return "", err
	}
	if statusCode != 200 {
		return "", ErrCaptcha
	}

	url = fmt.Sprintf("%v?handle=%v", ptttype.CAPTCHA_URL_PREFIX, handle)

	return url, nil
}
