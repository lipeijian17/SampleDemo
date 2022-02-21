#! /usr/bin/python
# ! -*- coding: utf-8 -*-

import sys
import os
import requests
import time
import re
import json
from random import randint
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from src.RtcTokenBuilder import RtcTokenBuilder,Role_Attendee

appID = "9a3b95751df14cebaf8f155448ee0e36"
appCertificate = "9654e7de17154329a55e716830e07d19"

def generate_RtcToken():
  #appID = appID
  #appCertificate = appCertificate
  channelName = "cloudRecording"
  uid = 666666
  userAccount = "666666"
  expireTimeInSeconds = 86400
  currentTimestamp = int(time.time())
  privilegeExpiredTs = currentTimestamp + expireTimeInSeconds

  token = RtcTokenBuilder.buildTokenWithUid(appID, appCertificate, channelName, uid, Role_Attendee, privilegeExpiredTs)
  print("Token with int uid: {}".format(token))
  token = RtcTokenBuilder.buildTokenWithAccount(appID, appCertificate, channelName, userAccount, Role_Attendee, privilegeExpiredTs)
  print("Token with user account: {}".format(token))
  return token

def generate_Resourceid():
  url = "https://api.agora.io/v1/apps/{AppID}/cloud_recording/acquire".format(AppID=appID)
  #degug
  print("1---------------------------1")
  print(url)
  print("1---------------------------1")

  payload = "{\n    \"cname\": \"cloudRecording\",\n    \"uid\": \"666666\",\n    \"clientRequest\": {\n        \"region\": \"CN\",\n        \"resourceExpiredHour\": 24\n    }\n}"
  
  
  headers = {
    'Content-type': 'application/json;charset=utf-8',
    'Authorization': 'Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA='
  }

  response = requests.request("POST", url, headers=headers, data=payload)
  
  resourceId_json=eval(response.text)["resourceId"]
  #degug
  print(resourceId_json)
  return resourceId_json
  
  
def start_Recording(token,resourceid):
  url = "https://api.agora.io/v1/apps/{AppID}/cloud_recording/resourceid/{Resourceid}/mode/mix/start".format(AppID=appID,Resourceid=resourceid)
  
  #degug
  print("2---------------------------2")
  print(url)
  print("2---------------------------2")
    
   #这里的token需要传参
  payload = "{\n    \"uid\": \"666666\",\n    \"cname\": \"cloudRecording\",\n    \"clientRequest\": {\n        \"token\": \"0069a3b95751df14cebaf8f155448ee0e36IABlyKXKXdHoGT3Up4ls2DGJ6Ftw6xDF5V/F0K1yesCXkQykWmeGMoD+IgCzspQCUsAUYgQAAQBSwBRiAgBSwBRiAwBSwBRiBABSwBRi\",\n        \"recordingConfig\": {\n            \"maxIdleTime\": 30,\n            \"streamTypes\": 2,\n            \"channelType\": 1,\n            \"transcodingConfig\": {\n                \"height\": 640,\n                \"width\": 360,\n                \"bitrate\": 500,\n                \"fps\": 15,\n                \"mixedVideoLayout\": 1,\n                \"backgroundColor\": \"#FF0000\"\n            },\n            \"subscribeVideoUids\": [\"#allstream#\"],\n            \"subscribeAudioUids\": [\"#allstream#\"],\n            \"subscribeUidGroup\": 0\n        },\n        \"recordingFileConfig\": {\n            \"avFileType\": [\n                \"hls\"\n            ]\n        },\n        \"storageConfig\": {\n            \"accessKey\": \"LTAI5t7wpCo6Ng8pkh5DYFDb\",\n            \"region\": 1,\n            \"bucket\": \"cloudrecording-lpj\",\n            \"secretKey\": \"ed4U00yl7QWSA2ba5LuByb68dDIslM\",\n            \"vendor\": 2,\n            \"fileNamePrefix\": [\n                \"Recording1\",\n                \"Recording2\"\n            ]\n        }\n    }\n}"
  
  #degug
  print(token)
  
  headers = {
    'Content-type': 'application/json;charset=utf-8',
    'Authorization': 'Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA='
  }

  response = requests.request("POST", url, headers=headers, data=payload)

  sid_json = eval(response.text)["sid"]
  #degug
  print("3---------------------------3")
  print(response.text)
  print("3---------------------------3")
  
  #degug
  print(sid_json)
  return sid_json
  
#更新参数示例：更新订阅名单
def update_RecordingCogfig(resourceid,sid):
  url = "https://api.agora.io/v1/apps/{AppID}}/cloud_recording/resourceid/{Resourceid}/sid/{Sid}/mode/mix/update".format(AppID=appID,Resourceid=resourceid,Sid=sid)

  payload = "{\n    \"uid\": \"666666\",\n    \"cname\": \"cloudRecording\",\n    \"clientRequest\": {\n        \"streamSubscribe\": {\n            \"audioUidList\": {\n                \"subscribeAudioUids\": [\n                    \"#allstream#\"\n                ]\n            },\n            \"videoUidList\": {\n                \"unSubscribeVideoUids\": \"999999\"\n            }\n        }\n    }\n}                      "
  headers = {
    'Content-type': 'application/json;charset=utf-8',
    'Authorization': 'Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA='
  }

  response = requests.request("POST", url, headers=headers, data=payload)
#degug
  print(response.text)


def query_RecordingFile(resourceid,sid):
  url = "https://api.agora.io/v1/apps/{AppID}/cloud_recording/resourceid/{Resourceid}/sid/{Sid}/mode/mix/query".format(AppID=appID,Resourceid=resourceid,Sid=sid)

  payload = ""
  headers = {
    'Content-type': 'application/json;charset=utf-8',
    'Authorization': 'Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA='
  }

  response = requests.request("GET", url, headers=headers, data=payload)
  #degug
  print(response.text)


def stop_RecordingCogfig(resourceid,sid):
  url = "https://api.agora.io/v1/apps/9a3b95751df14cebaf8f155448ee0e36/cloud_recording/resourceid/{Resourceid}/sid/{Sid}/mode/mix/stop".format(AppID=appID,Resourceid=resourceid,Sid=sid)

  payload = "{\n    \"cname\": \"cloudRecording\",\n    \"uid\": \"666666\",  \n    \"clientRequest\":{\n        \"async_stop\": false   \n    }\n}"
  headers = {
    'Content-type': 'application/json;charset=utf-8',
    'Authorization': 'Basic YzVhNmE2ZTRkYmRjNGVkNjhkYTcyNDRhMWU3ODgyNDQ6NGE3Mzk2ODQ2NzcyNDdmMGFlODc4MWNhMGEyYTllZDA='
  }

  response = requests.request("POST", url, headers=headers, data=payload)
  #degug
  print(response.text)



def main():
  token = generate_RtcToken();
  resourceid = generate_Resourceid();
  sid = start_Recording(token,resourceid);
  
  #sleep 20
  #query_RecordingFile(resourceid,sid)
  
  #update
  #update_RecordingCogfig(resourceid,sid)
  
  #stop
  #stop_RecordingCogfig(resourceid,sid)

if __name__ == "__main__":
  main()
