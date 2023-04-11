package rtmpconverter

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var appID = "9a3b9XXXXXXXXXXXXXXX36"
var appCertificate = "965XXXXXXXXXXXXXXXd19"

// 客户 ID
var customerKey = "444aXXXXXXXXXXXXXXXad7"

// 客户密钥
var customerSecret = "4c69XXXXXXXXXXXXXXXbd9"

// 拼接客户 ID 和客户密钥并使用 base64 进行编码
var plainCredentials = customerKey + ":" + customerSecret
var base64Credentials = base64.StdEncoding.EncodeToString([]byte(plainCredentials))

// All ApI:  1 Creat  2  Delete    3 Update    4 Get    5  List
//1. Creat Converter   no转码推流  converterId string
func CreatConverter() (converterId string, err error) {
	url := "https://api.agora.io/cn/v1/projects/9a3XXXXXXXXXXXXXXX36/rtmp-converters"
	fmt.Println(url)

	method := "POST"
	payload := strings.NewReader(`{
    "converter": {
      "name": "rtmp8",
      "rawOptions": {
        "rtcChannel": "rtmp",
        "rtcStreamUid": 201
      },
      "rtmpUrl": "rtmp://vid-218.push.chinanetcenter.broadcastapp.agora.io/live/testzzg"
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
	//test
	fmt.Println(req.Header)

	res, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err)
		return "", err2
	}

	defer res.Body.Close()
	//这里读取多层次的json,还需要修改调整
	converterid := ParseResponse(res)["converter"].(map[string]interface{})["id"]

	//???zhelizenmebiablidao xiayiji??
	value, ok := converterid.(string)
	if !ok {
		fmt.Println("It's not ok for type string")
		return "", nil
	}

	return string(value), nil
}

//2. Delete Creat Converter
func DeleteConverter(converterId string) {
	//url := strings.Replace("https://api.agora.io/cn/v1/projects/9a3bXXXXXXXXXXXXXXX0e36/rtmp-converters/ConverterId","ConverterId",converterId,-1)
	s1 := "https://api.agora.io/cn/v1/projects/9a3b9XXXXXXXXXXXXXXX36/rtmp-converters/"
	s2 := converterId
	fmt.Println(s2)
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(s2)
	url := build.String()
	//url := "https://api.agora.io/cn/v1/projects/9a3bXXXXXXXXXXXXXXXe36/rtmp-converters/" + converterId
	fmt.Println(url)
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

//调整混音用户的音量大小
func UpdateConverter(converterId string) (string, error) {
	url := strings.Replace("https://api.agora.io/cn/v1/projects/9a3XXXXXXXXXXXXXXXe36/rtmp-converters/ConverterId", "ConverterId", converterId, -1)
	fmt.Println(url)
	method := "PATCH"
	payload := strings.NewReader(`{
    "converter": {
    "transcodeOptions": {
      "audioOptions": {
        "volumes": [
          {
            "volume": 50,
            "rtcStreamUid": 201
          },
          {
            "volume": 150,
            "rtcStreamUid": 201
          }
        ]
      }
    }
  },
  "fields": "transcodeOptions.audioOptions.volumes"
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
	//test
	fmt.Println(req.Header)

	res, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err2)
		return "", err2
	}

	defer res.Body.Close()

	converterid := ParseResponse(res)["converter"].(map[string]interface{})["id"]

	return converterid.(string), nil

}

func GetConverter(converterId string) (string, error) {
	url := strings.Replace("https://api.agora.io/cn/v1/projects/9a3bXXXXXXXXXXXXXXX0e36/rtmp-converters/ConverterId", "ConverterId", converterId, -1)
	//url := "https://api.agora.io/cn/v1/projects/9aXXXXXXXXXXXXXXX0e36/rtmp-converters/20A4291CEA470E041D0E349C5F811EE1"

	fmt.Println(url)
	method := "GET"
	//payload := ""

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+base64Credentials)
	req.Header.Add("Content-Type", "application/json")
	//test
	fmt.Println(req.Header)

	res, err2 := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err2
	}

	defer res.Body.Close()
	convertername := ParseResponse(res)["converter"].(map[string]interface{})["name"]
	value, ok := convertername.(string)
	if !ok {
		fmt.Println("It's not ok for type string")
		return "", nil
	}

	return string(value), nil
}

func ListConverter() (string, error) {
	url := "https://api.agora.io/cn/v1/projects/9a3bXXXXXXXXXXXXXXX36/rtmp-converters"
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
	convertercount := ParseResponse(res)["data"].(map[string]interface{})["total_count"]

	value, ok := convertercount.(string)
	if !ok {
		fmt.Println("It's not ok for type string")
		return "", nil
	}

	return string(value), nil
}

func ParseResponse(response *http.Response) map[string]interface{} {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
		fmt.Println(err)
	}

	return result
}

func WriteFile_converterId(converterId string) {
	filePtr, err := os.Create("converterId_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(converterId)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}

func ReadFile_converterId() (result string) {
	//这里记得修改路径
	file, err := os.Open("/Users/lpj/go/src/demo/rtmpconverter/converterId_info.json") 
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close() 

	line := bufio.NewReader(file)

	content, _, _ := line.ReadLine()
	fmt.Println(string(content))
	return string(content)
}
