package main

import (
	rtctokenbuilder "RtcTokenBuilder"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var appID = ""
var appCertificate = ""

func generate_RtcToken() (token string) {

	appID := appID
	//	print(appID)
	appCertificate := appCertificate
	channelName := "cloudRecording"
	uid := uint32(666666)
	//uidStr := "666666"
	expireTimeInSeconds := uint32(86400)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Printf("Token with uid: %s\n", result)
		return result
	}

	//	result, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, rtctokenbuilder.RoleAttendee, expireTimestamp)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Printf("Token with userAccount: %s\n", result)
	//		return result
	//	}
	return
}

//这个函数继续调试, 返回string类型的 resourceid
func generate_Resourceid() (resourceid string) {
	//	url := "https://api.agora.io/v1/apps/Appid/cloud_recording/acquire"
	url := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/acquire", "Appid", appID, -1)

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
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", "Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA=")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	//fmt.Println(ParseResponse(res)["resourceId"])

	//return string(ParseResponse(res)["resourceId"].(string))
	return ParseResponse(res)["resourceId"]
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Println(string(body))

	//var returnMap map[string] string
	//returnMap = make(map[string]string)

	//fmt.Println(ParseResponse()["resourceId"].(string))
	//return
	//return ParseResponse(res)["resourceId"].(string)

}

func start_Recording(token string, resourceid string) (sid string) {
	//url := "https://api.agora.io/v1/apps/Resourceid/mode/mix/start"

	url_pre := strings.Replace("https://api.agora.io/v1/apps/Appid/cloud_recording/resourceid/Resourceid/mode/mix/start", "Appid", appID, -1) //小于0代表全部替换

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
            "maxIdleTime": 60,
            "streamTypes": 2,
            "channelType": 1,
            "subscribeVideoUids": ["#allstream#"],
            "subscribeAudioUids": ["#allstream#"],
            "subscribeUidGroup": 0
        },
        "recordingFileConfig": {
            "avFileType": [
                "hls"
            ]
        },
        "storageConfig": {
            "accessKey": "？",
            "region": 1,
            "bucket": "cloudrecording-lpj",
            "secretKey": "？",
            "vendor": 2,
            "fileNamePrefix": [
                "Recording1",
                "Recording2"
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
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", "Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	return ParseResponse(res)["sid"]
}

//这里要传两个参数
func query_RecordingFile(resourceid string, sid string) {
	//url := "https://api.agora.io/v1/apps/9a3b95751df14cebaf8f155448ee0e36/cloud_recording/resourceid/Resourceid/mode/mix/query"
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

func update_RecordingCogfig(resourceid string, sid string) {
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
                "unSubscribeVideoUids": ["999999"]
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

func stop_Recording(resourceid string, sid string) {
	//	url := "https://api.agora.io/v1/apps/9a3b95751df14cebaf8f155448ee0e36/cloud_recording/resourceid/Etkl6g-zSB7EpP-Da1zN63dq8zdhhwI1ia8deuCa71RRwWU9aHZSVwb9aIQPgDscqbVpoCCGL4bDhDVp4HVVrirV-1dzw8wMir6ox9xYvcn2R2Ctage_cH6DQEu4-5UTIjnoRr2V3KmzVQtY19OljdjJF4HUjq5gLmm4Zdh2n8ZCx9AZgROLy_C8ZAwYdhp1w11hFj_dHTAIfhl-kFv89adQgfs3eXyQ7yX4KAF45p3XBpgF3txvWpzvZmoxEGtqAGc209G6r8CiLnTWErO7PW6_-a-nyVs5U-uW65dZAsgBs0Cj2qZSbkXc-4Yc3-GUrpVP1BciekbXiKZYXgmqOff-qmW1MJmTH06t6UrcNOY/sid/441eb5c17d46ba0347d767b4761fe905/mode/mix/stop"
	//url2 := "https://api.agora.io/v1/apps/9a3b95751df14cebaf8f155448ee0e36/cloud_recording/resourceid/" + resourceid + "/sid/" + sid + "/mode/mix/stop"

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
func writeFile_resourceid(resourceid string) {
	// 创建文件
	filePtr, err := os.Create("resourceid_info.json")
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

func writeFile_sid(sid string) {
	// 创建文件
	filePtr, err := os.Create("sid_info.json")
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
func readFile_resourceid() (result string) {
	file, err := os.Open("/Users/lpj/Desktop/Cloud Recording/AgoraDemo/Cloud Recording/go/sample/RtcTokenBuilder/resourceid_info.json") //打开
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

func readFile_sid() (result string) {
	file, err := os.Open("/Users/lpj/Desktop/Cloud Recording/AgoraDemo/Cloud Recording/go/sample/RtcTokenBuilder/sid_info.json") //打开
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

func main() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)

	//	var token string
	//	var resourceid string
	//	var sid string

	if s == "start" {
		token := generate_RtcToken()
		print(token)
		print("\n\n")
		resourceid := generate_Resourceid()
		print(resourceid)
		print("\n\n")

		//start recording
		sid := start_Recording(token, resourceid)
		print(sid)
		//分别将resourceid和sid写在两个文件中

		writeFile_resourceid(resourceid)
		writeFile_sid(sid)

	}

	if s == "query" {
		//定义一个查询参数
		//删掉json文件中字符串的前后引号
		resourceid := strings.Trim(readFile_resourceid(), "\"")
		sid := strings.Trim(readFile_sid(), "\"")
		query_RecordingFile(resourceid, sid)
	}

	//定义一个更新的参数
	//update
	if s == "update" {
		//定义一个查询参数
		//删掉json文件中字符串的前后引号
		resourceid := strings.Trim(readFile_resourceid(), "\"")
		sid := strings.Trim(readFile_sid(), "\"")
		update_RecordingCogfig(resourceid, sid)
	}

	//定义一个stop的参数
	if s == "stop" {
		//resourceid := readFile_resourceid()
		//sid := readFile_sid()
		//删掉json文件中字符串的前后引号
		resourceid := strings.Trim(readFile_resourceid(), "\"")
		sid := strings.Trim(readFile_sid(), "\"")
		stop_Recording(resourceid, sid)
	}
}
