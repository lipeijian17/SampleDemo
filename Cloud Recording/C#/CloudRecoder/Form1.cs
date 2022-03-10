using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Net;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;
using Newtonsoft.Json.Linq;
using RestSharp;

namespace CloudRecoder
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
        }
        //todo
        //your app id
        static string AppId = "xxxxx";
        //your Authorization;
        static string Authorization = "Basic xxxx";
        //your channel id
        static string cname = "123123";
        //robot id
        static string uid = "527834";
        //oss accessKey
        static string accessKey = "xxxxxx-xxx";
        //oss secretKey
        static string secretKey = "xxxx-xxxxx";
        //oss bucket
        static string bucket = "xxxxx";
        //oss vendor
        static int vendor = 0; // Number 类型，第三方云存储平台 0：七牛云 2：阿里云
        //oss region 不同的oss 能存储的区域不同
        static int region = 0;  //0：华东 1：华北 2：华南 ..

        //recoding sid
        string sid = "";

        string resId = "";

        //{
        //	"cname": "httpClient463224", channel id
        //		"uid" : "527841",        rebot id  
        //		"clientRequest" : {
        //		"resourceExpiredHour": 24,
        //		"scene" : 0  //如果要使用延时转码，需要将 scene 设置为 2。
        //	}
        //}
        private string GetResId()
        {
            var url = @"https://api.agora.io/v1/apps/";
            var method = @"/cloud_recording/acquire";
            var resUrl = url + AppId + method;
            string resId = "";
            try
            {
                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.POST);
                
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");
                var jsonBody = new JObject();
                var jsonRes = new JObject();
                
                jsonBody.Add(new JProperty("cname", cname));
                
                jsonBody.Add(new JProperty("uid", uid));
                jsonRes.Add(new JProperty("resourceExpiredHour", 24));
                jsonRes.Add(new JProperty("scene", 0));
                jsonBody.Add(new JProperty("clientRequest", jsonRes));
                var strJson = jsonBody.ToString();

                request.AddParameter("application/json", strJson, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        Console.WriteLine($"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["reason"].ToString()}");
                    }
                    return "";
                }
                resId = jsonResponse["resourceId"].ToString();
                
            }
            catch (Exception e1)
            {
                Console.WriteLine(e1.Message);
            }
            return resId;
        }

        //单流录制
        //https://api.agora.io/v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/mode/individual/start
        //        {
        //    "uid": "527841",
        //    "cname": "httpClient463224",
        //    "clientRequest": {
        //        "token": "<token if any>",
        //        "appsCollection": {
        //            "combinationPolicy": "postpone_transcoding"     //延时转码
        //        },
        //        "recordingConfig": {
        //            "maxIdleTime": 30,
        //            "streamTypes": 2,
        //            "channelType": 0,                //0:通信场景（默认）1：直播场景
        //            "videoStreamType": 1,           //如果频道中有用户开启了双流模式 设置订阅的视频流类型0：视频大流（默认）1：视频小流
        //            "subscribeVideoUids": [
        //                "123",                               //recoded user id
        //                "456"
        //            ],
        //            "subscribeAudioUids": [
        //                "123",
        //                "456"
        //            ],
        //            "subscribeUidGroup": 0          //预估的订阅人数峰值。在单流模式下，为必填参数[0~5] 0：1 到 2 个 UID,1：3 到 7 个 UID
        //        },
        //        "storageConfig": {
        //            "accessKey": "xxxxxxf",
        //            "region": 3,
        //            "bucket": "xxxxx",
        //            "secretKey": "xxxxx",
        //            "vendor": 2,
        //            "fileNamePrefix": [
        //                "directory1",
        //                "directory2"
        //            ]
        //}
        //    }
        //}
        private string StartRecoding(string resId)
        {
            if (resId == "")
            {
                Console.Write("resId is empty");
                return "";
            }
            string sidId = "";
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/mode/individual/start";
            var resUrl = url + AppId + method1 + resId + method2;
            try
            {
                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.POST);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");

                var jsonBody = new JObject();
                jsonBody.Add(new JProperty("cname", cname));
                jsonBody.Add(new JProperty("uid", uid));
                //clientRequest
                
                var jsonClientRequest = new JObject();
                jsonClientRequest.Add(new JProperty("token", ""));
                //recordingConfig
                var jsonRecordingConfig = new JObject();
                jsonRecordingConfig.Add(new JProperty("maxIdleTime", 30));
                jsonRecordingConfig.Add(new JProperty("streamTypes", 2));
                jsonRecordingConfig.Add(new JProperty("channelType", 0));
                jsonRecordingConfig.Add(new JProperty("videoStreamType", 1));
                //need Recording uid  todo
                jsonRecordingConfig.Add(new JProperty("subscribeVideoUids", new JArray("1499286043","1122")));
                jsonRecordingConfig.Add(new JProperty("subscribeAudioUids", new JArray("1499286043","1122")));
                jsonRecordingConfig.Add(new JProperty("subscribeUidGroup", 0));

                jsonClientRequest.Add(new JProperty("recordingConfig", jsonRecordingConfig));

                //storageConfig
                var jsonStorageConfig = new JObject();
                jsonStorageConfig.Add(new JProperty("accessKey", accessKey));
                jsonStorageConfig.Add(new JProperty("region", region));
                jsonStorageConfig.Add(new JProperty("bucket", bucket));
                jsonStorageConfig.Add(new JProperty("secretKey", secretKey));
                jsonStorageConfig.Add(new JProperty("vendor", vendor));
                jsonStorageConfig.Add(new JProperty("fileNamePrefix", new JArray("directory1", "directory2")));

                jsonClientRequest.Add(new JProperty("storageConfig", jsonStorageConfig));

                jsonBody.Add(new JProperty("clientRequest", jsonClientRequest));
                var reqBody = jsonBody.ToString();

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        Console.WriteLine($"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["reason"].ToString()}");
                        textBox1.Text = $"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["reason"].ToString()}";
                    }
                }
                else
                {
                    sidId = jsonResponse["sid"].ToString();
                    textBox1.Text = $"sidId:{sidId}\r\n";
                }
                
            }
            catch (Exception e)
            {
                Console.Write(e.Message);
            }
            return sidId;
        }

        
        //start
        private void button1_Click(object sender, EventArgs e)
        {

            resId = GetResId();

            sid = StartRecoding(resId);
        }

        //https://api.agora.io/v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/sid/<sid>/mode/individual/stop
        // {
        //    "cname": "httpClient463224",
        //    "uid": "527841",
        //    "clientRequest": {}
        //}
        //stop
        private void button2_Click(object sender, EventArgs e)
        {
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/sid/";
            var method3 = @"/mode/individual/stop";
            var resUrl = url + AppId + method1 + resId + method2 + sid + method3;
            try
            {
                var jsonBody = new JObject();
                jsonBody.Add(new JProperty("cname", cname));
                jsonBody.Add(new JProperty("uid", uid));
                jsonBody.Add(new JProperty("clientRequest", new JObject()));

                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.POST);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");
                var reqBody = jsonBody.ToString();

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        var serverResponse = jsonResponse.Property("serverResponse");
                        if (serverResponse != null)
                        {
                            Console.WriteLine($"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["serverResponse"].ToString()}");
                            textBox1.Text = $"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["serverResponse"].ToString()}";
                        }
                        else
                        {
                            textBox1.Text = $"code: {jsonResponse["code"].ToString()}";
                        }
                    }
                }
                else
                {
                    textBox1.Text = "stop recoding";
                }
            }
            catch(Exception e1)
            {
                Console.Write(e1.Message);
            }
        }


        //https://api.agora.io/v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/sid/<sid>/mode/individual/query
        //query
        private void button3_Click(object sender, EventArgs e)
        {
            //jsonPa();
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/sid/";
            var method3 = @"/mode/individual/query";
            var resUrl = url + AppId + method1 + resId + method2 + sid + method3;
            try
            {
                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.GET);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");
                var reqBody = "";

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        Console.WriteLine($"code: {jsonResponse["code"].ToString()}");
                    }
                    return;
                }
                var res = jsonResponse["serverResponse"];
                JArray jArray = JArray.Parse(res["fileList"].ToString());
                foreach (var arr in jArray)
                {
                    JObject jObj = JObject.Parse(arr.ToString());
                    string fileName = jObj["fileName"].ToString();
                    string trackType = jObj["trackType"].ToString();
                    textBox1.Text += $"fileName: {fileName}\r\n";
                    textBox1.Text += $"trackType: {trackType}\r\n"; 
                }

            }
            catch(Exception e1)
            {
                Console.Write(e1.Message);
            }
        }


        //https://api.agora.io/v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/mode/mix/start
        //
        //        {
        //    "uid": "527841",
        //    "cname": "httpClient463224",
        //    "clientRequest": {
        //        "token": "<token if any>",
        //        "recordingConfig": {
        //            "maxIdleTime": 30,
        //            "streamTypes": 2,
        //            "audioProfile": 1,     //设置输出音频的采样率[0-2]
        //            "channelType": 0,
        //            "videoStreamType": 0,
        //            "transcodingConfig": {
        //                "height": 640,
        //                "width": 360,
        //                "bitrate": 500,
        //                "fps": 15,
        //                "mixedVideoLayout": 1,     //设置视频合流布局 0、1、2 为预设的合流布局，3 为自定义合流布局。该参数设为 3 时必须设置 layoutConfig 参数。
        //                "backgroundColor": "#FF0000"   //视频画布的背景颜色
        //            },
        //            "subscribeVideoUids": [
        //                "123",
        //                "456"
        //            ],
        //            "subscribeAudioUids": [
        //                "123",
        //                "456"
        //            ],
        //            "subscribeUidGroup": 0
        //        },
        //        "storageConfig": {
        //            "accessKey": "xxxxxxf",
        //            "region": 3,
        //            "bucket": "xxxxx",
        //            "secretKey": "xxxxx",
        //            "vendor": 2,
        //            "fileNamePrefix": [
        //                "directory1",
        //                "directory2"
        //            ]
        //}
        //    }
        //}
        private string StartMulitRecoding(string resId)
        {
            if (resId == "")
            {
                Console.Write("resId is empty");
                return "";
            }
            string sidId = "";
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/mode/mix/start";
            var resUrl = url + AppId + method1 + resId + method2;

            try
            {
                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.POST);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");

                var jsonBody = new JObject();
                jsonBody.Add(new JProperty("cname", cname));
                jsonBody.Add(new JProperty("uid", uid));
                //clientRequest

                var jsonClientRequest = new JObject();
                jsonClientRequest.Add(new JProperty("token", ""));
                //recordingConfig
                var jsonRecordingConfig = new JObject();
                jsonRecordingConfig.Add(new JProperty("maxIdleTime", 30));
                jsonRecordingConfig.Add(new JProperty("streamTypes", 2));
                jsonRecordingConfig.Add(new JProperty("audioProfile", 1));
                jsonRecordingConfig.Add(new JProperty("channelType", 0));
                jsonRecordingConfig.Add(new JProperty("videoStreamType", 1));
                //transcodingConfig
                var jsonTranscodingConfig = new JObject();
                jsonTranscodingConfig.Add(new JProperty("height", 640));
                jsonTranscodingConfig.Add(new JProperty("width", 360));
                jsonTranscodingConfig.Add(new JProperty("bitrate", 500));
                jsonTranscodingConfig.Add(new JProperty("fps", 15));
                jsonTranscodingConfig.Add(new JProperty("mixedVideoLayout", 1));
                jsonTranscodingConfig.Add(new JProperty("backgroundColor", "#FF0000"));

                jsonRecordingConfig.Add(new JProperty("transcodingConfig", jsonTranscodingConfig));
                //need Recording uid  todo
                jsonRecordingConfig.Add(new JProperty("subscribeVideoUids", new JArray("1499286043", "1122")));
                jsonRecordingConfig.Add(new JProperty("subscribeAudioUids", new JArray("1499286043", "1122")));
                jsonRecordingConfig.Add(new JProperty("subscribeUidGroup", 0));

                jsonClientRequest.Add(new JProperty("recordingConfig", jsonRecordingConfig));

                //storageConfig
                var jsonStorageConfig = new JObject();
                jsonStorageConfig.Add(new JProperty("accessKey", accessKey));
                jsonStorageConfig.Add(new JProperty("region", region));
                jsonStorageConfig.Add(new JProperty("bucket", bucket));
                jsonStorageConfig.Add(new JProperty("secretKey", secretKey));
                jsonStorageConfig.Add(new JProperty("vendor", vendor));
                jsonStorageConfig.Add(new JProperty("fileNamePrefix", new JArray("directory1", "directory2")));

                jsonClientRequest.Add(new JProperty("storageConfig", jsonStorageConfig));

                jsonBody.Add(new JProperty("clientRequest", jsonClientRequest));
                var reqBody = jsonBody.ToString();

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        Console.WriteLine($"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["reason"].ToString()}");
                        textBox1.Text = $"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["reason"].ToString()}";
                    }
                }
                else
                {
                    sidId = jsonResponse["sid"].ToString();
                    textBox1.Text = $"sidId:{sidId}\r\n";
                }

            }
            catch (Exception e)
            {
                Console.Write(e.Message);
            }
            return sidId;
        }

        private void button4_Click(object sender, EventArgs e)
        {
            resId = GetResId();

            sid = StartMulitRecoding(resId);
        }

        //https://api.agora.io/v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/sid/<sid>/mode/mix/stop
        //stop
        //
        //{
        //    "cname": "httpClient463224",
        //    "uid": "527841",
        //    "clientRequest": {}
        //}
        private void button5_Click(object sender, EventArgs e)
        {
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/sid/";
            var method3 = @"/mode/mix/stop";
            var resUrl = url + AppId + method1 + resId + method2 + sid + method3;
            try
            {
                var jsonBody = new JObject();
                jsonBody.Add(new JProperty("cname", cname));
                jsonBody.Add(new JProperty("uid", uid));
                jsonBody.Add(new JProperty("clientRequest", new JObject()));

                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.POST);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");
                var reqBody = jsonBody.ToString();

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        var serverResponse = jsonResponse.Property("serverResponse");
                        if (serverResponse != null)
                        {
                            Console.WriteLine($"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["serverResponse"].ToString()}");
                            textBox1.Text = $"code: {jsonResponse["code"].ToString()} ,reason: {jsonResponse["serverResponse"].ToString()}";
                        }
                        else
                        {
                            textBox1.Text = $"code: {jsonResponse["code"].ToString()}";
                        }
                    }
                }
                else
                {
                    textBox1.Text = "stop recoding";
                }
            }
            catch (Exception e1)
            {
                Console.Write(e1.Message);
            }
        }

        //
        private void button6_Click(object sender, EventArgs e)
        {
            var url = @"https://api.agora.io/v1/apps/";
            var method1 = @"/cloud_recording/resourceid/";
            var method2 = @"/sid/";
            var method3 = @"/mode/mix/query";
            var resUrl = url + AppId + method1 + resId + method2 + sid + method3;
            try
            {
                var client = new RestClient(resUrl);
                client.Timeout = -1;
                var request = new RestRequest(Method.GET);
                //your Authorization;
                request.AddHeader("Authorization", Authorization);
                request.AddHeader("Content-Type", "application/json");
                var reqBody = "";

                request.AddParameter("application/json", reqBody, ParameterType.RequestBody);
                var response = client.Execute(request);

                JObject jsonResponse = JObject.Parse(response.Content);

                //bad request
                if (response.StatusCode != HttpStatusCode.OK)
                {
                    var codeProperty = jsonResponse.Property("code");
                    if (codeProperty != null)
                    {
                        Console.WriteLine($"code: {jsonResponse["code"].ToString()}");
                    }
                    return;
                }
                var res = jsonResponse["serverResponse"];
                var fileList = res["fileList"].ToString();
                textBox1.Text = $"fileName: {fileList}\r\n";
                //JArray jArray = JArray.Parse(res["fileList"].ToString());
                //foreach (var arr in jArray)
                //{
                //    JObject jObj = JObject.Parse(arr.ToString());
                //    string fileName = jObj["fileName"].ToString();
                //    string trackType = jObj["trackType"].ToString();
                //    textBox1.Text += $"fileName: {fileName}\r\n";
                //    textBox1.Text += $"trackType: {trackType}\r\n";
                //}

            }
            catch (Exception e1)
            {
                Console.Write(e1.Message);
            }
        }
    }
}
