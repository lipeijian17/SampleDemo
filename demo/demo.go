package main

import (
	"demo/channelmanager"
	"demo/cloudplayer"
	"demo/cloudrecording"
	"demo/cloudtranscoding"
	"demo/mutibitrate"
	"demo/rtmpconverter"
	"fmt"
	"os"
	"strings"
)

func main() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)

	//******云端录制接口******
	if s == "startCloudRecording" {
		token := cloudrecording.Generate_RtcToken()
		print(token)
		print("\n\n")
		resourceid := cloudrecording.Generate_Resourceid()
		print("-------------------------------------------")
		print(resourceid)
		print("-------------------------------------------")
		print("\n\n")
		cloudrecording.WriteFile_resourceid(resourceid)

		cloudrecording.Wgcloudrecording.Add(1)
		cloudrecording.Start_Recording(token, resourceid, func(sid string, err error) {
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			print("\n\n")
			print("-------------------------------------------")
			print(sid)
			print("-------------------------------------------")
			print("\n\n")

			cloudrecording.WriteFile_sid(sid)
			print(sid)
			print("\n\n")
		})
		print("\n\n")
		fmt.Println("Waiting for response...")
		cloudrecording.Wgcloudrecording.Wait()
	}

	if s == "queryCloudRecording" {
		resourceid := strings.Trim(cloudrecording.ReadFile_resourceid(), "\"")
		sid := strings.Trim(cloudrecording.ReadFile_sid(), "\"")
		cloudrecording.Query_RecordingFile(resourceid, sid)
	}

	//update
	if s == "updateCloudRecording" {
		//定义一个查询参数
		resourceid := strings.Trim(cloudrecording.ReadFile_resourceid(), "\"")
		sid := strings.Trim(cloudrecording.ReadFile_sid(), "\"")
		cloudrecording.Update_RecordingCogfig(resourceid, sid)
	}

	if s == "updatelayoutCloudRecording" {
		//定义一个查询参数
		resourceid := strings.Trim(cloudrecording.ReadFile_resourceid(), "\"")
		sid := strings.Trim(cloudrecording.ReadFile_sid(), "\"")
		cloudrecording.UpdateLayout_RecordingCogfig(resourceid, sid)
	}

	//定义一个stop的参数
	if s == "stopCloudRecording" {
		resourceid := strings.Trim(cloudrecording.ReadFile_resourceid(), "\"")
		sid := strings.Trim(cloudrecording.ReadFile_sid(), "\"")
		cloudrecording.Stop_Recording(resourceid, sid)
	}

	//******Rtmp Converter接口******
	if s == "creatRtmpConverter" {
		//creat converter

		converterId, _ := rtmpconverter.CreatConverter()
		//var converterI
		print(converterId)
		print("\n\n")
		//将converterId写在文件中，真实开发环境可以存在数据库
		rtmpconverter.WriteFile_converterId(converterId)
	}

	if s == "deleteRtmpConverter" {
		converterId := strings.Trim(rtmpconverter.ReadFile_converterId(), "\"")
		print(converterId)
		print("\n\n")
		rtmpconverter.DeleteConverter(converterId)
	}

	if s == "getRtmpConverter" {
		print("\n\n")
		converterId := strings.Trim(rtmpconverter.ReadFile_converterId(), "\"")
		print("\n\n")
		print(converterId)
		rtmpconverter.GetConverter(converterId)
	}

	if s == "listRtmpConverter" {
		//list converter

		convertercount, _ := rtmpconverter.ListConverter()
		print(convertercount)
		print("\n\n")
	}

	//******Cloud Player接口******
	if s == "creatCloudPlayer" {
		//creat converter
		token := cloudplayer.Generate_RtcToken()
		print(token)
		print("\n\n")
		playerid, _ := cloudplayer.CreatCloudPlayer(token)
		cloudplayer.WriteFile_playerid(playerid)
		print(playerid)
		print("\n\n")
	}

	if s == "deleteCloudPlayer" {
		playerId := strings.Trim(cloudplayer.ReadFile_playerid(), "\"")
		print(playerId)
		print("\n\n")
		cloudplayer.DeleteCloudPlayer(playerId)
	}

	if s == "updateCloudPlayer" {
		//delete converter
		playerId := strings.Trim(cloudplayer.ReadFile_playerid(), "\"")
		print(playerId)
		print("\n\n")

		cloudplayer.UpdateCloudPlayer(playerId)
	}

	if s == "listCloudPlayer" {
		//delete converter
		playerId := strings.Trim(cloudplayer.ReadFile_playerid(), "\"")
		print(playerId)
		print("\n\n")

		cloudplayer.ListCloudPlayer()
	}

	//******云端转码/合图接口******
	//开始转码
	if s == "startCloudTranscoder" {
		tokenName := cloudtranscoding.Acquire_BuildToken()
		print("\n\n")
		print("-------------------------------------------")
		print("\n\n")
		print(tokenName) //
		print("\n\n")
		print("-------------------------------------------")
		print("\n\n")

		cloudtranscoding.Wgcloudtranscoding.Add(1)
		cloudtranscoding.Create_CloudTranscoder(tokenName, func(taskid string, err error) {
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			print("\n\n")
			print("-------------------------------------------")
			print(taskid)
			print("-------------------------------------------")
			print("\n\n")
			//调试看到 taskid没有生成
			cloudtranscoding.WriteFile_taskid(taskid)
			print("\n\n")
			print(taskid)
			print("\n\n")

		})

		fmt.Println("Waiting for response...")
		cloudtranscoding.Wgcloudtranscoding.Wait()

	}

	if s == "updateCloudTranscoder" {
		//待实现
	}

	if s == "queryCloudTranscoder" {
		//待实现
	}

	if s == "deleteCloudTranscoder" {
		//待实现
	}

	//******极速直播多码率//******
	//启用多码率功能
	if s == "enableMutibitrate" {
		status := mutibitrate.EnableMutibitrate()
		print(status)
	}
	//关闭多码率功能
	if s == "disableMutibitrate" {
		mutibitrate.DisableMutibitrate()
	}

	if s == "updateMutibitrate_480p" {
		mutibitrate.UpdateMutibitrate_480p()
	}

	if s == "updateMutibitrate_360p" {
		mutibitrate.UpdateMutibitrate_360p()
	}

	if s == "queryMutibitrate" {
		mutibitrate.QueryMutibitrate()
	}

	if s == "queryMutibitrateChannels" {
		mutibitrate.QueryMutibitrateChannels()
	}

	//******频道管理接口//******

	//踢出某个人
	if s == "kickoffuser" {
		channelmanager.Kickoffuser()
	}

	//解散整个频道
	if s == "kickoffchannel" {
		channelmanager.Kickoffchannel()
	}

}
