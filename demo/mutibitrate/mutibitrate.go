package mutibitrate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//多码率观众端，本例需要把 authorization字段由Basic auth的鉴权改为hmac格式，还需继续修改，再次之前，可以先用node.js的demo验证功能
var appID = "9a3b9XXXXXXXXXXXXXXXXXXXXXee0e36"

//var appCertificate = "9654e7dXXXXXXXXXXXXXXXXXXXXX0e07d19"
// 客户 ID  客户密钥
var customerKey = "33eXXXXXXXXXXXXXXXXXXXXX696"
var customerSecret = "c801XXXXXXXXXXXXXXXXXXXXX2f9a"

var timestamp = time.Now().UTC().String()

// Sha256 Calculate the sha256 hash of a string
func Sha256(str string) []byte {
	h := sha256.New()
	_, _ = h.Write([]byte(str))
	return h.Sum(nil)
}

// HmacSha256 Calculate the sha256 hash of a string using the HMAC method
func HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))

	return mac.Sum(nil)
}

// HmacSha256ToString Calculate the sha256 hash of a string using the HMAC method, outputs lowercase hexits
func HmacSha256ToString(key string, data string) string {
	return hex.EncodeToString(HmacSha256(key, data))
}

func EnableMutibitrate() (status string) {
	url := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/config", "Appid", appID, -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "POST"

	payload := strings.NewReader(`{
			"enabled" : true
		 }`)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	//拼接生成需要加密的字符串    ==>签名的算法还是有问题，下一步如何推进，推进完部署docker
	request_line := "POST /v1/projects/9a3b9XXXXXXXXXXXXXXXXXXXXXe36/rtls/abr/config HTTP/2.0"
	print(request_line)
	signing_string_pre := "host: api.agora.io\ndate: Date\nPOST /v1/projects/9a3b95751XXXXXXXXXXXXXXXXXXXXXee0e36/rtls/abr/config HTTP/2.0"
	signing_string := strings.Replace(signing_string_pre, "Date", timestamp, -1)
	Signing := HmacSha256ToString(signing_string, customerSecret)

	authorization_pre := strings.Replace("hmac username=CustomerKey, algorithm='hmac-sha256', headers='host date request-line', signature='Signature'", "CustomerKey", customerKey, -1)
	authorization := strings.Replace(authorization_pre, "Signature", Signing, -1)

	bodySign := Sha256("")
	digest_pre := "SHA-256=" + string(bodySign)

	str, _ := base64.StdEncoding.DecodeString(digest_pre)
	realStr := strings.TrimSpace(string(str))
	str2 := []byte(realStr)
	digest := base64.StdEncoding.EncodeToString(str2)

	// 增加 Authorization header
	//req.Header.Add("x-date", timestamp.Format(time.RFC3339))
	//req.Header.Add("Date", time.Now().Format(time.RFC3339))
	req.Header.Add("x-date", timestamp)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Digest", digest)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")

	//print(signature)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", body)
	return
	//return ParseResponse(res)["status"]
}

func DisableMutibitrate() (status string) {
	url := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/config", "Appid", appID, -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "POST"

	payload := strings.NewReader(`{
			"enabled" : false
		 }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header   这里需要更新修改
	req.Header.Add("Authorization", "Basic "+"signature")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	return ParseResponse(res)["status"]
}

//创建或更新转码模板配置
func UpdateMutibitrate_480p() (status string) {

	url_pre := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/config/codecs/CodecId", "Appid", appID, -1)
	url := strings.Replace(url_pre, "CodecId", "480p", -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "POST"

	payload := strings.NewReader(`{
		"video": {
			"width": 854,
			"height": 480,
			"fps": 15
		  }
		 }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+"signature")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	return ParseResponse(res)["status"]
}

//创建或更新转码模板配置
func UpdateMutibitrate_360p() (status string) {

	url_pre := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/config/codecs/CodecId", "Appid", appID, -1)
	url := strings.Replace(url_pre, "CodecId", "360p", -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "POST"

	payload := strings.NewReader(`{
		"video": {
			"width": 640,
			"height": 360,
			"fps": 15
		  }
		 }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+"signature")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	return ParseResponse(res)["status"]
}

//查询转码模板配置
func QueryMutibitrate() (status string) {

	url_pre := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/config/codecs", "Appid", appID, -1)
	url := strings.Replace(url_pre, "CodecId", "360p", -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "POST"

	payload := strings.NewReader(`{"" : ""}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+"signature")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	return ParseResponse(res)["status"]
}

//观众端查询所有可选的频道名称&码率参数
func QueryMutibitrateChannels() (status string) {

	url_pre := strings.Replace("https://api.agora.io/v1/projects/Appid/rtls/abr/target-channels/MutibitrateChannels?uid=999999", "Appid", appID, -1)
	url := strings.Replace(url_pre, "CodecId", "360p", -1)

	print("\n\n")
	print(url)
	print("\n\n")
	method := "GET"

	payload := strings.NewReader(`{"" : ""}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+"signature")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	return ParseResponse(res)["status"]
}

func ParseResponse(response *http.Response) map[string]string {
	var result map[string]string
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}
	return result
}
