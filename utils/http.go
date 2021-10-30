package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/mockhttp"
	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
)

//BackendPost
//
//Params
//  postData: http-post data
//  result: resp-data, requires pointer of pointer to malloc.
//
//Ex:
//    url := backend.LOGIN_R
//    postData := &backend.LoginParams{}
//    result := &backend.LoginResult{}
//    BackendPost(c, url, postData, nil, &result)
func BackendPost(url string, postData interface{}, headers map[string]string, result interface{}) (statusCode int, err error) {
	if isTest {
		return mockhttp.HTTPPost(url, postData, result)
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	jsonBytes, err := json.Marshal(postData)
	if err != nil {
		return 500, err
	}

	buf := bytes.NewBuffer(jsonBytes)

	// req
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return 500, err
	}

	return httpProcess(req, headers, result)
}

func BackendGet(url string, params interface{}, headers map[string]string, result interface{}) (statusCode int, err error) {
	if isTest {
		return mockhttp.HTTPPost(url, params, result)
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	v, _ := query.Values(params)
	url = url + "?" + v.Encode()

	// req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 500, err
	}

	return httpProcess(req, headers, result)
}

func httpProcess(req *http.Request, headers map[string]string, result interface{}) (statusCode int, err error) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// send http
	resp, err := httpClient.Do(req)
	if err != nil {
		return 500, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 501, err
	}

	if statusCode != 200 {
		return statusCode, fmt.Errorf("http status-code err status-code: %v body: %v", statusCode, string(body))
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return 501, err
	}

	return 200, nil
}

func MergeURL(urlMap map[string]string, url string) string {
	urlList := strings.Split(url, "/")

	newURLList := make([]string, len(urlList))
	for idx, each := range urlList {
		if len(each) == 0 || each[0] != ':' {
			newURLList[idx] = each
			continue
		}

		theKey := each[1:]
		theVal := urlMap[theKey]

		newURLList[idx] = theVal
	}

	return strings.Join(newURLList, "/")
}

func GetCookie(c *gin.Context, name string) string {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return ""
	}
	if cookie == nil {
		return ""
	}

	return cookie.Value
}
