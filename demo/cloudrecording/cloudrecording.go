package cloudrecording

import (
	rtctokenbuilder "RtcTokenBuilder"
	"bufio"
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

var appID = "9a3b95751df14cebaf8f155448ee0e36"
var appCertificate = "9654e7de17154329a55e716830e07d19"

// 客户 ID  客户密钥
var customerKey = "33ebb114b3284064bf7d049e16f4d696"
var customerSecret = "c80110826fc74e45a9559518176a2f9a"

// 拼接客户 ID 和客户密钥并使用 base64 进行编码
var plainCredentials = customerKey + ":" + customerSecret
var base64Credentials = base64.StdEncoding.EncodeToString([]byte(plainCredentials))

var Wgcloudrecording sync.WaitGroup

//定义回调函数，回调的参数类型都一样的话，用一个回调就可以
//type callback_resourceid func(resourceid string, err error)
//type callback_sid func(sid string, err error)
type callback func(sid string, err error)

func Generate_RtcToken() (token string) {
	appID := appID
	appCertificate := appCertificate
	channelName := "cloudRecording"
	uid := uint32(666666)
	expireTimeInSeconds := uint32(86400)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds
	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		return result
	}
	return
}

func Generate_Resourceid() (resourceid string) {
	print(1111111111111111111)
	print("\n\n")

	url := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/acquire", "Appid", appID, -1)
	print(url)
	print("\n\n")
	method := "POST"
	payload := strings.NewReader(`{
          "cname": "cloudRecording",
	      "uid": "666666",
	      "clientRequest": {
     	    "region": "CN",
		    "resourceExpiredHour": 24
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

	defer res.Body.Close()
	return ParseResponse(res)["resourceId"]
}

func Start_Recording(token string, resourceid string, c callback) {
	go func() {
		url_pre := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/resourceid/Resourceid/mode/mix/start", "Appid", appID, -1) //小于0代表全部替换
		print(url_pre)
		print("\n\n")

		url := strings.Replace(url_pre, "Resourceid", resourceid, -1) //小于0代表全部替换

		print(url)
		print("\n\n")

		method := "POST"
		payload := strings.NewReader(strings.Replace(`{
      "uid": "666666",    
      "cname": "cloudRecording",
      "clientRequest": {
        "token": "Token",
        "recordingConfig": {
			"maxIdleTime":30,
			"streamTypes":2,
			"audioProfile":1,
			"channelType":1,
			"videoStreamType":0,
		    "transcodingConfig":{
				"width":640,
				"height":360,
				"fps":15,
				"bitrate":1200,
				"mixedVideoLayout":1,
				"backgroundColor":"#FF0000",
				"backgroundImage":"https://rela-livetest.oss-cn-zhangjiakou.aliyuncs.com/agoraVideo/%E8%A7%86%E9%A2%91%E5%9B%9B%E4%BA%BA%E5%8D%A0%E4%BD%8D%E5%9B%BE.png"
				}
        },

        "storageConfig": {
            "accessKey": "LTAI5t9AGYGk783aiEiXdsKy",
            "region": 1,
            "bucket": "cloudrecording-lpj",
            "secretKey": "NQoWsQKKr971e1nqGHyAoqOoE3qyFk",
            "vendor": 2,
            "fileNamePrefix": [
                "Recording1",
                "Recording21"
            ]
        }
      }
    }`, "Token", token, -1))

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

		//fmt.Println(res)
		defer res.Body.Close()

		time.Sleep(2 * time.Second)

		c(ParseResponse(res)["sid"], nil)
		//等待异步执行完毕
		Wgcloudrecording.Done()
	}()
}

//这里要传两个参数
func Query_RecordingFile(resourceid string, sid string) {
	url_pre := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/resourceid/Resourceid/sid/Sid/mode/mix/query", "Resourceid", resourceid, -1) //小于0代表全部替换
	url_pre2 := strings.Replace(url_pre, "Sid", sid, -1)
	url := strings.Replace(url_pre2, "Appid", appID, -1)

	method := "GET"

	payload := strings.NewReader(``)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Update_RecordingCogfig(resourceid string, sid string) {
	url_pre := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/resourceid/Resourceid/sid/Sid/mode/mix/update", "Resourceid", resourceid, -1) //小于0代表全部替换

	url_pre2 := strings.Replace(url_pre, "Sid", sid, -1)
	url := strings.Replace(url_pre2, "Appid", appID, -1)
	method := "POST"

	payload := strings.NewReader(`{
	  "uid": "666666",
	  "cname": "cloudRecording",
	  "clientRequest": {
		  "streamSubscribe": {
			  "audioUidList": {
				  "subscribeAudioUids": [
					"#allstream#"
				  ]
			  },
			  "videoUidList": {
                "subscribeVideoUids": ["888888"]
            }
		  }
	  }
  }                      `)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func UpdateLayout_RecordingCogfig(resourceid string, sid string) {
	url_pre := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/resourceid/Resourceid/sid/Sid/mode/mix/update", "Resourceid", resourceid, -1) //小于0代表全部替换

	url_pre2 := strings.Replace(url_pre, "Sid", sid, -1)
	url := strings.Replace(url_pre2, "Appid", appID, -1)
	method := "POST"

	payload := strings.NewReader(`{
	  "uid": "666666",
	  "cname": "cloudRecording",
	  "clientRequest": {
		"backgroundColor":"#000000"
		"backgroundImage":"https://rela-livetest.oss-cn-zhangjiakou.aliyuncs.com/agoraVideo/%E8%A7%86%E9%A2%91%E5%9B%9B%E4%BA%BA%E5%8D%A0%E4%BD%8D%E5%9B%BE.png"

	  }
  }                      `)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func Stop_Recording(resourceid string, sid string) {

	url_pre := strings.Replace("https://api.agora.io/v1/apps/9a3b95751df14cebaf8f155448ee0e36/cloud_recording/resourceid/Resourceid/sid/Sid/mode/mix/stop", "Resourceid", resourceid, -1) //小于0代表全部替换
	url := strings.Replace(url_pre, "Sid", sid, -1)

	fmt.Println(url)
	method := "POST"

	payload := strings.NewReader(`{
	  "cname": "cloudRecording",
	  "uid": "666666",  
	  "clientRequest":{
		  "async_stop": false   
	  }
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", "Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

//对Http请求的返回结果转map处理
func ParseResponse(response *http.Response) map[string]string {
	var result map[string]string
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}
	return result
}

//定义golang的json读写函数

func WriteFile_resourceid(resourceid string) {
	// 创建文件
	filePtr, err := os.Create("/Users/lpj/go/src/demo/cloudrecording/resourceid_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(resourceid)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}

func WriteFile_sid(sid string) {
	// 创建文件
	filePtr, err := os.Create("/Users/lpj/go/src/demo/cloudrecording/sid_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(sid)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}

//定义读取文件
func ReadFile_resourceid() (result string) {
	file, err := os.Open("/Users/lpj/go/src/demo/cloudrecording/resourceid_info.json") //打开
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close() //关闭

	line := bufio.NewReader(file)

	content, _, _ := line.ReadLine()
	fmt.Println(string(content))
	return string(content)
}

func ReadFile_sid() (result string) {
	file, err := os.Open("/Users/lpj/go/src/demo/cloudrecording/sid_info.json") //打开
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close() //关闭

	line := bufio.NewReader(file)

	content, _, _ := line.ReadLine()
	fmt.Println(string(content))
	return string(content)
}
