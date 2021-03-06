package api

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/config_util"
)

const configPrefix = "go-pttbbs:api"

func InitConfig() error {
	config()
	postInitConfig()
	return nil
}

func setStringConfig(idx string, orig string) string {
	return config_util.SetStringConfig(configPrefix, idx, orig)
}

func setBytesConfig(idx string, orig []byte) []byte {
	return config_util.SetBytesConfig(configPrefix, idx, orig)
}

func setIntConfig(idx string, orig int) int {
	return config_util.SetIntConfig(configPrefix, idx, orig)
}

func postInitConfig() {
	_ = setJwtTokenExpireTS(JWT_TOKEN_EXPIRE_TS)
	_ = setEmailJwtTokenExpireTS(EMAIL_JWT_TOKEN_EXPIRE_TS)
}

func setJwtTokenExpireTS(JwtTokenExpireTS int) (origJwtTokenExpireTS int) {
	origJwtTokenExpireTS = JWT_TOKEN_EXPIRE_TS

	JWT_TOKEN_EXPIRE_TS = JwtTokenExpireTS
	JWT_TOKEN_EXPIRE_DURATION = time.Duration(JWT_TOKEN_EXPIRE_TS) * time.Second

	return origJwtTokenExpireTS
}

func setEmailJwtTokenExpireTS(EmailJwtTokenExpireTS int) (origEmailJwtTokenExpireTS int) {
	origEmailJwtTokenExpireTS = EMAIL_JWT_TOKEN_EXPIRE_TS

	EMAIL_JWT_TOKEN_EXPIRE_TS = EmailJwtTokenExpireTS
	EMAIL_JWT_TOKEN_EXPIRE_DURATION = time.Duration(EMAIL_JWT_TOKEN_EXPIRE_TS) * time.Second

	return origEmailJwtTokenExpireTS
}
