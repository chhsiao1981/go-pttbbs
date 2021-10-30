package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const ATTEMPT_REGISTER_R = "/attempt_register"

type AttemptRegisterResult struct {
	URL    string `json:"url"`
	Verify string `json:"verify"`
}

func AttemptRegisterWrapper(c *gin.Context) {
	JSON(AttemptRegister, nil, c)
}

func AttemptRegister(remoteAddr string, params interface{}) (result interface{}, err error) {
	url, verify, err := bbs.AttemptRegister()
	if err != nil {
		return nil, err
	}

	return &AttemptRegisterResult{
		URL:    url,
		Verify: verify,
	}, nil
}
