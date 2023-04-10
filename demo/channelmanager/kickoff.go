package channelmanager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var appID = "9a3b95751df14cebaf8f155448ee0e36"
var appCertificate = "9654e7de17154329a55e716830e07d19"

// 客户 ID  客户密钥
var customerKey = "33ebb114b3284064bf7d049e16f4d696"
var customerSecret = "c80110826fc74e45a9559518176a2f9a"

// 拼接客户 ID 和客户密钥并使用 base64 进行编码
var plainCredentials = customerKey + ":" + customerSecret
var base64Credentials = base64.StdEncoding.EncodeToString([]byte(plainCredentials))

func Kickoffuser() {
	url := "https://api.agora.io/dev/v1/kicking-rule"
	print(url)
	print("\n\n")

	method := "POST"
	payload := strings.NewReader(strings.Replace(`{
		"appid": "Appid",
		"cname": "cloudPlayer",
		"uid": 777777,
		"ip": "",
		"time": 0,
		"privileges": [
		  "join_channel"
		]

    }`, "Appid", appID, -1))

	client := &http.Client{}
	//作为client端生成的request body，需不需要手动关闭呢，答案是不需要的
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")
	//	req.Header.Add("X-Request-ID", "111111111111111")

	res, err := client.Do(req)
	print(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	//打印response
	fmt.Println(res)

	//作为client端处理response的时候，有一点要注意的是，body一定要close，否则会造成GC回收不到，继而产生内存泄露；
	defer res.Body.Close()
}

//解散频道
func Kickoffchannel() {
	url := "https://api.agora.io/dev/v1/kicking-rule"
	print(url)
	print("\n\n")

	method := "POST"
	payload := strings.NewReader(strings.Replace(`{
		"appid": "Appid",
		"cname": "cloudPlayer",
		"time": 0,
		"privileges": [
		  "join_channel"
		]

    }`, "Appid", appID, -1))

	client := &http.Client{}
	//作为client端生成的request body，需不需要手动关闭呢，答案是不需要的
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")
	//	req.Header.Add("X-Request-ID", "111111111111111")

	res, err := client.Do(req)
	print(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	//打印response
	fmt.Println(res)

	//作为client端处理response的时候，有一点要注意的是，body一定要close，否则会造成GC回收不到，继而产生内存泄露；
	defer res.Body.Close()

}

//对Http请求的返回结果转map处理
func ParseResponse(response *http.Response) map[string]interface{} {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
		fmt.Println(err)
	}

	return result
}
