/*
//project: ccp
//file: ccp.go
//time: 2019/1/23 18:13:43
//author: ceephoen
//contact: ceephoen@163.com
//software: GoLand
//license: ceephoen@163.com
//desc:
*/
package sdk

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unsafe"
)

type CCP struct {
	// CCP struct like CCP class in Python SDK
	IP       string
	Port     string
	Version  string
	AppId    string
	AccSid   string
	AccToken string
}

func (ccp *CCP) Init(to, data []string, templateId string) (url, body string, headers map[string]string) {
	/*
		to: the number to send;
		data: the data to send;
		smsId: the ID of sms template.
	*/

	// format timestamp
	batch := time.Now().Format(TimeFormat)

	// sign
	sign := ccp.AccSid + ccp.AccToken + batch

	// md5
	MD5 := md5.New()
	MD5.Write([]byte(sign))
	lowerSign := hex.EncodeToString(MD5.Sum(nil))

	// lowerSign to upperSign
	upperSign := strings.ToUpper(lowerSign)

	// combine url
	url = strings.Join([]string{"https://", ccp.IP, ":", ccp.Port, "/", ccp.Version, "/Accounts/", ccp.AccSid, "/SMS/TemplateSMS?sig=", upperSign}, "")

	// auth
	src := ccp.AccSid + ":" + batch
	auth := base64.StdEncoding.EncodeToString([]byte(src))

	// body
	var b string
	for _, d := range data {
		b += strings.Join([]string{d, ","}, "")
	}
	b = "[" + b + "]"

	s := `{"to": "%s", "datas": %s, "templateId": "%s", "appId": "%s"}`

	// []string{"mobile", "mobile"} ---> "mobile,mobile"
	sl := strings.Join(to, ",")

	body = fmt.Sprintf(s, sl, b, templateId, ccp.AppId)
	headers = map[string]string{"Accept": "application/json", "Content-Type": "application/json;charset=utf-8", "Authorization": auth}

	return
}

func SendCode(to, data []string, templateId string) (r map[string]interface{}) {

	// instance
	var Ccp = CCP{ServerIp, ServerPort, SoftVersion, AppId, AccountSid, AccountToken}

	// url, body , headers
	var url, body string
	var headers map[string]string

	url, body, headers = Ccp.Init(to, data, templateId)

	// http-Client
	client := &http.Client{}

	// request
	request, _ := http.NewRequest("POST", url, strings.NewReader(body))

	// add headers
	//request.Header.Set("Accept", headers["Accept"])
	//request.Header.Set("Content-Type", headers["Content-Type"])
	//request.Header.Set("Authorization", headers["Authorization"])

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	// post-request
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	// read response body
	b, err := ioutil.ReadAll(resp.Body)

	// parse response body
	//StringBody
	sb := (*string)(unsafe.Pointer(&b))
	if err != nil {
		// handle error
		log.Fatal("json parse error ", err)
	}

	// result map
	r = make(map[string]interface{}, 0)

	// json deserialization
	e := json.Unmarshal([]byte(*sb), &r)

	// json deserialization error handler
	if e != nil {
		log.Fatal("json parse error ", e)
	}

	// return r
	log.Println("response data: ", r)
	return
}
