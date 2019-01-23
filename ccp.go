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
package ccp

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"encoding/base64"
	"net/http"
	"io/ioutil"
	"unsafe"
	"fmt"
	"encoding/json"
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

func (ccp *CCP) Create(to, smsId string, data []string) (url, body string, headers map[string]string) {
	/*
	to: the number to send;
	data: the data to send;
	smsId: the ID of of template.
	 */

	// format timestamp
	batch := time.Now().Format(timeFormat)

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

	body = strings.Join([]string{"{to:", to, ",", "datas:", b, ",", "templateId:", smsId, ",", "appId:", ccp.AppId, "}"}, "")
	headers = map[string]string{"Accept": "application/json", "Content-Type": "application/json;charset=utf-8", "Authorization": auth}

	return
}

func SendCode(to string, data []string, smsId string) (cons map[string]interface{}) {

	// instance
	var Ccp = CCP{ServerIp, ServerPort, SoftVersion, AppId, AccountSid, AccountToken}

	// url, body , headers
	var url, body string
	var headers map[string]string
	smsId = SmsId

	url, body, headers = Ccp.Create(to, smsId, data)

	// http-Client
	client := &http.Client{}

	// request
	request, _ := http.NewRequest("POST", url, strings.NewReader(body))

	// headers
	request.Header.Set("Accept", headers["Accept"])
	request.Header.Set("Content-Type", headers["Content-Type"])
	request.Header.Set("Authorization", headers["Authorization"])

	// post-request
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	// parse response body
	stringBody := (*string)(unsafe.Pointer(&b))
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// result map
	cons = make(map[string]interface{})

	// json deserialization
	e := json.Unmarshal([]byte(*stringBody), &cons)

	// json deserialization error handler
	if err != nil {
		fmt.Println(e)
	} else {
		fmt.Println("statusCode--->", cons["statusCode"])         // statusCode---> 000000
		fmt.Printf("statusCode-Type--->%T\n", cons["statusCode"]) // statusCode-Type--->string
	}

	// return cons
	return
}
