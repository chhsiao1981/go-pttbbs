package mockhttp

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

func HTTPPost(url string, data interface{}, result interface{}) (statusCode int, err error) {
	log.Infof("HTTPPost: url: %v", url)
	switch url {
	default:
		return 500, ErrURL
	}
}

func parseResult(backendResult interface{}, httpResult interface{}) (statusCode int, err error) {
	convert(backendResult, httpResult)

	return 200, nil
}

func convert(dataBackend interface{}, dataResult interface{}) {
	valueBackend := reflect.ValueOf(dataResult)
	valuePttbbs := reflect.ValueOf(dataBackend)
	valueBackend.Elem().Set(valuePttbbs)
}
