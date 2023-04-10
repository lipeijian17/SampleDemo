package cloudplayer

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
	"time"
)

var appID = "9a3b95751df14cebaf8f155448ee0e36"
var appCertificate = "9654e7de17154329a55e716830e07d19"

// 客户 ID  客户密钥
var customerKey = "33ebb114b3284064bf7d049e16f4d696"
var customerSecret = "c80110826fc74e45a9559518176a2f9a"
var timestamp = time.Now().UTC().String()

// 拼接客户 ID 和客户密钥并使用 base64 进行编码
var plainCredentials = customerKey + ":" + customerSecret
var base64Credentials = base64.StdEncoding.EncodeToString([]byte(plainCredentials))

//var Wgcloudplayer sync.WaitGroup

//type callback func(playerid string, err error)

func Generate_RtcToken() (token string) {
	appID := appID
	appCertificate := appCertificate
	channelName := "cloudPlayer"
	//这里暂时固定云播放器的uid位 ：777777
	uid := uint32(777777)
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

func CreatCloudPlayer(token string) (playerid string, err error) {
	url := strings.Replace("https://api.agora.io/cn/v1/projects/Appid/cloud-player/players", "Appid", appID, -1) //小于0代表全部替换
	print(url)
	print("\n\n")

	method := "POST"
	payload := strings.NewReader(strings.Replace(`{
			"player": {
				"streamUrl": "https://public-bucket-pro-1305122626.cos.ap-shanghai.myqcloud.com/defaultuserbackgroundimage.jpg",
				"channelName": "cloudPlayer",
				"token": "Token",
				"uid": 777777,
				"idleTimeout": 60,
				"playTs": 0,
				"name": "test123456"
				}

    }`, "Token", token, -1))

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
	req.Header.Add("X-Request-ID", "111111111111111")

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

	//这里需要修改
	playerid_pre := ParseResponse(res)["player"].(map[string]interface{})["id"]
	//converterid := ParseResponse(res)["converter"].(map[string]interface{})["id"]

	value, ok := playerid_pre.(string)
	if !ok {
		fmt.Println("It's not ok for type string")
		return "", nil
	}

	return string(value), nil
}

//这里读取文件的格式还需要更新
func DeleteCloudPlayer(playerid string) {
	url_pre := strings.Replace("https://api.agora.io/cn/v1/projects/Appid/cloud-player/players/playerId", "playerId", playerid, -1)
	url := strings.Replace(url_pre, "Appid", appID, -1) //小于0代表全部替换

	print(url)
	print("\n\n")

	method := "DELETE"
	//payload := ""

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-ID", "222222222222222")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	//打印response
	fmt.Println(res)
	defer res.Body.Close()
}

//调整混音用户的音量大小
func UpdateCloudPlayer(playerid string) (string, error) {
	url_pre1 := strings.Replace("https://api.agora.io/cn/v1/projects/Appid/cloud-player/players/playerId?sequence=Sequence", "Appid", appID, -1)
	url_pre2 := strings.Replace(url_pre1, "Sequence", "1", -1)
	url := strings.Replace(url_pre2, "playerId", playerid, -1)
	fmt.Println(url)
	method := "PATCH"
	payload := strings.NewReader(`{
		"player": {
			"audioOptions": {
				"volume": 10
			},
			"isPause": false,
			"streamUrl": "https://haokan.baidu.com/v?vid=9667043335917964919&",
			"seekPosition": 40
		 }
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-ID", timestamp)
	//test
	fmt.Println(req.Header)

	res, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err2)
		return "", err2
	}

	defer res.Body.Close()
	//???zhelizenmebiablidao xiayiji??
	return "", nil

}

func ListCloudPlayer() (string, error) {
	//本例不添加筛选条件，只查询当前所有云播放器
	url := strings.Replace("https://api.agora.io/v1/projects/Appid/cloud-player/players", "Appid", appID, -1)

	fmt.Println(url)
	method := "GET"
	//payload := ""

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 增加 Authorization header
	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")
	//test
	fmt.Println(req.Header)

	res, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err)
		return "", err2
	}

	defer res.Body.Close()
	//???zhelizenmebiablidao xiayiji??
	playercount := ParseResponse(res)["totalSize"]
	return playercount.(string), nil
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

func WriteFile_playerid(playerid string) {
	// 创建文件
	filePtr, err := os.Create("/Users/lpj/go/src/demo/cloudplayer/playerid_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(playerid)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}

//定义读取文件
func ReadFile_playerid() (result string) {
	file, err := os.Open("/Users/lpj/go/src/demo/cloudplayer/playerid_info.json") //打开
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
