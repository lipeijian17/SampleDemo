package cloudtranscoding

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var appID = "9a3XXXXXXXXXXXXXXXXXXXXXe0e36"

var appCertificate = "9654e7XXXXXXXXXXXXXXXXXXXXX07d19"

// 客户 ID  客户密钥
var customerKey = "33ebb1XXXXXXXXXXXXXXXXXXXXXf4d696"
var customerSecret = "c80110XXXXXXXXXXXXXXXXXXXXX176a2f9a"

// 拼接客户 ID 和客户密钥并使用 base64 进行编码
var plainCredentials = customerKey + ":" + customerSecret
var base64Credentials = base64.StdEncoding.EncodeToString([]byte(plainCredentials))

var Wgcloudtranscoding sync.WaitGroup

type callback func(taskid string, err error)

func Acquire_BuildToken() (tokenName string) {
	url := strings.Replace("https://api.agora.io/v1/projects/Appid/rtsc/cloud-transcoder/builderTokens", "Appid", appID, -1)

	print("\n\n")
	print(url)
	print("\n\n")

	method := "POST"
	payload := strings.NewReader(`{
		"instanceId" : "abc123456"
     }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	print("\n\n")
	print(ParseResponse(res)["tokenName"])
	print("\n\n")
	return ParseResponse(res)["tokenName"]
}

//RTC真实用户的两个token还需动态生成替换进去
func Create_CloudTranscoder(tokenName string, c callback) {
	go func() {
		url_pre := strings.Replace("https://api.agora.io/v1/projects/Appid/rtsc/cloud-transcoder/tasks?builderToken=TokenName", "Appid", appID, -1)
		print(url_pre)
		print("\n\n")
		url := strings.Replace(url_pre, "TokenName", tokenName, -1)
		//url := "https://api.agora.io/v1/projects/9a3bXXXXXXXXXXXXXXXXXXXXX8ee0e36/rtsc/cloud-transcoder/tasks?builderToken=guxRPm8euZb3q2bkMEEaPvv26_uPR-NVFjr-BzpIgZPUGH4X1j-VF-UXZcow_QbtFE_mYkr4hJIyWrui5O2w8nYqHT3D2KeKV07umQXhD796Gn3cOP5uy8yTbOucmZxy5pQzeK1Ds-v4a2pj9y2K-V8HUbvKTCzsSE0srJtF_W0VU2es4biFOnSroiekup46VBiVbyjHm7vCWz4Zy4LnvQ"
		print(url)
		print("\n\n")

		method := "POST"
		payload := strings.NewReader(`{
		"services":{
			"cloudTranscoder":{
				"serviceType":"cloudTranscoderV2",
				"config":{
					"transcoder":{
						"idleTimeout":300,
						"audioInputs":[
							{
								"rtc":{
									"rtcChannel":"test01",
									"rtcUid": 123,
									//这里的token需要重新生成
									"rtcToken":"0069a3b95751df14cebaf8f155448ee0e36IACOfutrO514AdLt0d0KTeeSNaB70oL8OehE2McOnrSzWH7vSdTSY0iIIgAXlFsAAsQaZAQAAQD+wxpkAgD+wxpkAwD+wxpkBAD+wxpk"
								}
							},
							{
								"rtc":{
									"rtcChannel":"test01",
									"rtcUid": 456,
									//这里的token需要重新生成
									"rtcToken":"0069a3b95751df14cebaf8f155448ee0e36IACIjmRzArwQQQSE15kMUc8nUIFl7Jd7JkIfvZkwWEUw837vSdRxw6ixIgBp0p8BXcQaZAQAAQD+wxpkAgD+wxpkAwD+wxpkBAD+wxpk"
								}
							}
						],
						"canvas":{
							"width":960,
							"height":480,
							"color":0,
							"backgroundImage":"https://public-bucket-pro-1305122626.cos.ap-shanghai.myqcloud.com/defaultuserbackgroundimage.jpg"，
							"fillMode": "FIT"
						},
						"waterMarks":[
							{
								"imageUrl":"https://cdn.dev.utown.io/i/20220621/9/c/a/9ca603d2eac048359f586cbc9b054db7.png",
								"region":{
									"x":0,
									"y":0,
									"width":100,
									"height":100,
									"zOrder":50
								}
							}
						],
						"videoInputs":[
							{
								"rtc":{
									"rtcChannel":"test01",
									"rtcUid": 123,
									"rtcToken":"0069a3b95751df14cebaf8f155448ee0e36IACOfutrO514AdLt0d0KTeeSNaB70oL8OehE2McOnrSzWH7vSdTSY0iIIgAXlFsAAsQaZAQAAQD+wxpkAgD+wxpkAwD+wxpkBAD+wxpk"
								},
								"placeholderImageUrl":"https://public-bucket-pro-1305122626.cos.ap-shanghai.myqcloud.com/defaultuserbackgroundimage.jpg",
								"region":{
									"x":0,
									"y":0,
									"width":320,
									"height":360,
									"zOrder":1
								}
							},
							{
								"rtc":{
									"rtcChannel":"test01",
									"rtcUid": 456,
									"rtcToken":"0069a3b95751df14cebaf8f155448ee0e36IACIjmRzArwQQQSE15kMUc8nUIFl7Jd7JkIfvZkwWEUw837vSdRxw6ixIgBp0p8BXcQaZAQAAQD+wxpkAgD+wxpkAwD+wxpkBAD+wxpk"
								},
								"placeholderImageUrl":"https://public-bucket-pro-1305122626.cos.ap-shanghai.myqcloud.com/defaultuserbackgroundimage.jpg",
								"region":{
									"x":320,
									"y":0,
									"width":320,
									"height":320,
									"zOrder":1
								}
							}
						],
						"outputs":[
							{
								"rtc":{
									"rtcChannel":"test",
									"rtcUid":1000,
									//这里的token要动态生成传进来
									"rtcToken":"0069a3b95751df14cebaf8f155448ee0e36IABTlswZDxRJuGFdmG8Y3+tC62RtzzV8zaeISIJ0Y1P05wx+f9gXoye0IgAi4RcEXnc1ZAQAAQDQdDVkAgDQdDVkAwDQdDVkBADQdDVk"
								},
								"audioOption":{
									"profileType":"AUDIO_PROFILE_MUSIC_STANDARD"
								},
								"videoOption":{
									"fps":30,
									"codec":"H264",
									"bitrate":800,
									"width":960,
									"height":480,
									"lowBitrateHighQuality":false
								}
							}
						]
					}
				}
			}
		}

     }`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}

		// 增加 Authorization header
		req.Header.Add("Authorization", "Basic "+base64Credentials)
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)

		defer res.Body.Close()
		print("\n\n")
		time.Sleep(5 * time.Second)
		print("\n\n")
		print(ParseResponse(res)["taskid"])
		print("\n\n")

		c(ParseResponse(res)["taskid"], nil)
		Wgcloudtranscoding.Done()

	}()
}

func ParseResponse(response *http.Response) map[string]string {
	var result map[string]string
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}
	return result
}

func WriteFile_taskid(taskid string) {
	// 创建文件
	filePtr, err := os.Create("/Users/lpj/go/src/demo/cloudtranscoding/taskid_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(taskid)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}
