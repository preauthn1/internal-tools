require("Module")
mmodule(...)

local G = require "Global"
local Download = require("Tools/DownloadTxtFile")

local isSimulation = false
local simulationInitRes = 0
local simulationGetIpRes = 1
local simulationGetIpIp = "0.0.0.0"
local simulationGetIpPort = 10000

local isEncode = false
local encodeFileName = "encodeFile.txt"

local useNewEncode = true

-- 版本号
local ver = 1

-- 連接 1游戏盾，2传统网络模式
local use = 1

-- 熱更 1游戏盾，2传统网络模式
local useCdn = 2

-- app key
local key = "NGhod3ZrMFlLRXpiTFdlY3hxZEh3clFzYlN6Ukl0Vnk4WTlQaU5iNzhEWXhibkxmUmM0amdsMGFPYVp5OGJrZTBMeWxfMWs3b0F6Sjd5WUdSNytxWkRoelM4UHE0c0ZmVjc2YTVDaXg2VE90UGlTLTNrSXR1cENvWlVuSThBOFpsYVhReDNOOXVDdGlWc2NjRnZPTjZpbjJ1YS1fU0NWa3Z5MnEtNWxxa0NxOHJsRlBmUUM2VElPdUV4MytBdDhKRHo1MjdzdklqdF9naHJrTmlnSnMrVHpTakNZb0tJWU1hM3I0Q0dETlgxcVMzNklLRHNBZlZKdEpNdGhXa19WOE1XalJsbGhQTnpQbUZPMzRXU3JQLUpRQXRJamNBYTQ4WG1aSHFSRmtJczRraWlKcFdGRUNVTU12RTlSc3FKQ3Z6b2wxQTFMbmJVZ0dpTWhyUlMxRGNaa1RCWTh1YXRRNXRjNEVuMjdYZXVCdjRocWJLbThHCg=="

-- 一般节点组标识
local groupnameNormal = "example.invalid";
-- Vip节点组标识
local groupnameVip = "example.invalid";

local port = 10000
-- 防护目标标识 大廳
local mLobby = "forward"
-- stage
local mLobbyStage = "servertest"

local portCdn = 8090
-- 防护目标标识 熱更
local mCdn = "update"

-- 檔案位置
local CnFileName ="aliyun_Cn.txt"
local CnFileNameNew ="aliyun_NewCn.txt"
local CnFileListDomain = 
{
	"aHR0cHM6Ly9zMy1hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3RlcnByb2NsaWVudC8=",
	"aHR0cDovL2tlZnUuemxndWFuZ29uZy5jb20vcnQv",
	"aHR0cHM6Ly9yYXcuZ2l0aHVidXNlcmNvbnRlbnQuY29tL1RlcnByb0NsaWVudC9hcmsvbWFzdGVyLw==",
	"aHR0cDovL2RsLmxtZmM4OC5jb20vZmlsZS8=",
}

local StageFileName ="aliyun_Stage.txt"
local StageFileNameNew ="aliyun_NewStage.txt"
local StageFileListDomain = 
{
	"aHR0cHM6Ly9zMy1hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3RlcnByb2NsaWVudC8=",
	"aHR0cDovL2tlZnUuemxndWFuZ29uZy5jb20vcnQv",
	"aHR0cHM6Ly9yYXcuZ2l0aHVidXNlcmNvbnRlbnQuY29tL1RlcnByb0NsaWVudC9hcmsvbWFzdGVyLw==",
	"aHR0cDovL2RsLmxtZmM4OC5jb20vZmlsZS8=",
}

local AnFileName ="aliyun_An.txt"
local AnFileNameNew ="aliyun_NewAn.txt"
local AnFileListDomain = 
{
	"aHR0cHM6Ly9zMy1hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3RlcnByb2NsaWVudC8=",
	"aHR0cDovL2tlZnUuemxndWFuZ29uZy5jb20vcnQv",
	"aHR0cHM6Ly9yYXcuZ2l0aHVidXNlcmNvbnRlbnQuY29tL1RlcnByb0NsaWVudC9hcmsvbWFzdGVyLw==",
	"aHR0cDovL2RsLmxtZmM4OC5jb20vZmlsZS8=",
}

local QaFileName ="aliyun_Qa.txt"
local QaFileNameNew ="aliyun_NewQa.txt"
local QaFileListDomain = 
{
	"aHR0cDovLzYxLjIxOC40NC4yMzE6MTQwLw==",
}

local DevFileName ="aliyun_Dev.txt"
local DevFileNameNew ="aliyun_NewDev.txt"
local DevFileListDomain =
{
	"aHR0cDovLzYxLjIxOC40NC4yMzE6MTQwLw==",
}

local ListDomain
local FileName

local tryCount = 3
local accTryCount = 0
local initSuccess = false

function init()
	if CS.UnityEngine.PlayerPrefs.GetInt("aliyun_version", 0) == 0 then
		Save(ver, use, key, groupnameNormal, groupnameVip, port, mLobby, useCdn, portCdn, mCdn)
	end
	
	if isSimulation == true then
		Save(ver, use, key, groupnameNormal, groupnameVip, port, mLobby, useCdn, portCdn, mCdn)	
	end	
end

function Save(verData, useData, keyData, groupnameNormalData, groupnameVipData, portData, mLobbyData, useCdnData, portCdnData, mCdnData)
	print("verData", verData)
	print("keyData", keyData)
	print("groupnameNormalData", groupnameNormalData)
	print("groupnameVipData", groupnameVipData)
	print("useData", useData)
	print("portData", portData)
	print("mLobbyData", mLobbyData)
	print("useCdnData", useCdnData)
	print("portCdnData", portCdnData)
	print("mCdnData", mCdnData)
	CS.UnityEngine.PlayerPrefs.SetInt("aliyun_version", verData)
	CS.UnityEngine.PlayerPrefs.SetString("aliyun_appkey", keyData)
	
	CS.UnityEngine.PlayerPrefs.SetString("aliyun_groupnameNormal", groupnameNormalData)
	CS.UnityEngine.PlayerPrefs.SetString("aliyun_groupnameVip", groupnameVipData)

	CS.UnityEngine.PlayerPrefs.SetInt("aliyun_type", useData)
	CS.UnityEngine.PlayerPrefs.SetInt("aliyun_port", portData)
	CS.UnityEngine.PlayerPrefs.SetString("aliyun_mark_lobby", mLobbyData)

	CS.UnityEngine.PlayerPrefs.SetInt("aliyun_type_cdn", useCdnData)
	CS.UnityEngine.PlayerPrefs.SetInt("aliyun_port_cdn", portCdnData)
	CS.UnityEngine.PlayerPrefs.SetString("aliyun_mark_cdn", mCdnData)

	CS.UnityEngine.PlayerPrefs.Save()
end

-- 是否支援遊戲盾
function IsAliyunSDK()
	if isSimulation == true then
		return true
	end

	if G.GameConfig.GetServerOption() == 0 or
		G.GameConfig.GetServerOption() == 1 then
		return false
	end
		
	
	-- 編輯器下讀設定設定
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor or 
		CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor	 then
		return false
	else
		if CS.AliyunSDK.s_IsAliyunSDK == true then
			return true
		else
			return false
		end
	end
end

function GetVer()
	return CS.UnityEngine.PlayerPrefs.GetInt("aliyun_version")
end

function GetType()
	return CS.UnityEngine.PlayerPrefs.GetInt("aliyun_type")
end

-- 目前連線方式
local nowUse
function GetNowUse()
	return nowUse
end
function SetNowUse(use)
	nowUse = use
end

-- 下載的設定連線方式
local downloadUse
function GetDownloadUse()
	return downloadUse
end

function GetUseCdn()
    if IsAliyunSDK() ~= true then
		return 2
	-- 編輯器下讀設定
	elseif CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor or 
		CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor	 then
		return 2
	-- statge、美服和國服讀設定
	elseif G.GameConfig.GetServerOption() == 3 or
			G.GameConfig.GetServerOption() == 4 or 
			G.GameConfig.GetServerOption() == 5 then
		return CS.UnityEngine.PlayerPrefs.GetInt("aliyun_type_cdn")
	-- 預設走TF
	else
		return 2
	end	
end

function GetKey()
	return CS.UnityEngine.PlayerPrefs.GetString("aliyun_appkey")
end

function GetMLobby()
	-- 美服和國服
	if G.GameConfig.GetServerOption() == 3 or 
		G.GameConfig.GetServerOption() == 4 then
		return CS.UnityEngine.PlayerPrefs.GetString("aliyun_mark_lobby")
	-- stage
	else
		return mLobbyStage	
	end
end

function GetMCdn()
	return mCdn
end

function GetGroupnameNormal()
	return CS.UnityEngine.PlayerPrefs.GetString("aliyun_groupnameNormal")
end

function GetGroupnameVip()
	return CS.UnityEngine.PlayerPrefs.GetString("aliyun_groupnameVip")
end

function GetPort()
	return CS.UnityEngine.PlayerPrefs.GetInt("aliyun_port")
end

function GetPortCdn()
	return CS.UnityEngine.PlayerPrefs.GetInt("aliyun_port_cdn")
end

function GetToken()
	local HDCode = G.Env.GetHDCode()
	return HDCode
end

function GetGroupName()
	-- 依造level走不同groupname
	local userLevel = 1
	local localUserData = G.User.GetLocalUserInfo()
	if (localUserData ~= nil) then
		userLevel = localUserData.UserLevel
	end
	
	local groupName
	if userLevel > 1 then
		groupName = GetGroupnameVip()
	else
		groupName = GetGroupnameNormal()
	end
	
	return groupName
end

function ShowLog(str)
	print(str)
end

function RandomList(t)
    if type(t)~="table" then
        return
    end
	
    local temp = {}
	for _, data in pairs(t) do
		table.insert(temp, data)
	end
	t = temp
	
    local tab={}
    local index=1
    while #t~=0 do
        local n=math.random(0,#t)
        if t[n]~=nil then
            tab[index]=t[n]
            table.remove(t,n)
            index=index+1
        end
    end
    return tab
end

local downlodIndex = 1
-- ErrorCallback
-- 0 key不同,需強更
-- 1 全部檔案下載失敗
-- 2 Aes Decode失敗
-- 3 下載後Init失敗
-- 4 decodeBase64失敗
-- 5 無新資料
function DownloadAliyunFile(SuccessCallback, ErrorCallback)
	m_onSuccessCallback = SuccessCallback
	m_onErrorCallback = ErrorCallback
	
	downlodIndex = 1
	
	if G.GameConfig.GetServerOption() == 0 then
		ListDomain = DevFileListDomain
		FileName = DevFileName
	elseif G.GameConfig.GetServerOption() == 1 then
		ListDomain = QaFileListDomain
		FileName = QaFileName
	elseif G.GameConfig.GetServerOption() == 3 then
		ListDomain = AnFileListDomain
		FileName = AnFileName
	elseif G.GameConfig.GetServerOption() == 4 then
		ListDomain = CnFileListDomain
		FileName = CnFileName
	elseif G.GameConfig.GetServerOption() == 5 then
		ListDomain = StageFileListDomain
		FileName = StageFileName
	end
	
	if useNewEncode == true then
		if G.GameConfig.GetServerOption() == 0 then
			FileName = DevFileNameNew
		elseif G.GameConfig.GetServerOption() == 1 then
			FileName = QaFileNameNew
		elseif G.GameConfig.GetServerOption() == 3 then
			FileName = AnFileNameNew
		elseif G.GameConfig.GetServerOption() == 4 then
			FileName = CnFileNameNew
		elseif G.GameConfig.GetServerOption() == 5 then
			FileName = StageFileNameNew
		end
	end
	
	if isEncode == true then
		FileName = encodeFileName
	end
	
	
	-- 打亂順序
	ListDomain = RandomList(ListDomain)
	
	DoDownloadAliyunFile()
end

function DoDownloadAliyunFile()
	if downlodIndex <= #ListDomain then
		local url = G.Tool.Base64.decodeBase64(ListDomain[downlodIndex]) .. FileName
		downlodIndex = downlodIndex + 1
		Download.DownloadFile(url, onSuccessCallback, DoDownloadAliyunFile)	
	else
		ShowLog("全部檔案下載失敗")
		m_onErrorCallback(1)
	end
end

function onSuccessCallback(data)

	if isEncode == true then
		local jasnData = G.Tool.Aes.Encode(data)
		CS.UnityEngine.Debug.LogWarning("########################## Encode data\n" .. jasnData)
		--jasnData = G.Tool.Aes.Decode(jasnData)
		--print(jasnData)		
	end
	

	local dataList = G.Tool.StringHelper.NewSplit(data, "\n")
	if #dataList == 1 then
		if useNewEncode == false then
			data = G.Tool.Aes.Decode(data)
		else
			data = G.Tool.Aes.DecodeNoJson(data)
		end		
		if data == nil then
			m_onErrorCallback(2)
			return			
		end
		dataList = G.Tool.StringHelper.NewSplit(data, "\n")
	end
	
	local nowVer = GetVer()
	local _fileVer = tonumber(dataList[1])
	if _fileVer >= nowVer then
		downloadUse = tonumber(dataList[2])
		local _use = 1
		local _key = dataList[3]
		local _groupnameNormal = dataList[4]
		local _groupnameVip = dataList[5]
		local _port = tonumber(dataList[6])
		local _mLobby = dataList[7]
		local _useCdn = dataList[8]		
		local _portCdn = dataList[9]
		local _mCdn = dataList[10]
		
		Save(_fileVer, _use, _key, _groupnameNormal, _groupnameVip, _port, _mLobby, _useCdn,_portCdn, _mCdn)
		
		local isNeedRestart = false
		local _nowKey = GetKey()
		
		-- key有更新，必須重新啟動遊戲
		if _nowKey ~= _key then
			isNeedRestart = true
		end
		
		if isNeedRestart == false then
			if initSuccess == false then
				ShowLog("尚未初始化過,需初始化")
				accTryCount = 0
				InitEx(true)
			else
				ShowLog("已經初始化過")			
				m_onSuccessCallback()
			end
		else
			ShowLog("key不相同，必須重新啟動")
			m_onErrorCallback(0)
		end
	else
		downloadUse = 1
		ShowLog("無新資料")
		m_onSuccessCallback()
	end
end

function Init()
	init()
	
	InitEx(false)
end

-- SDK初始化
function InitEx(isNeedCallback)
	if accTryCount >= tryCount then
		if isNeedCallback == true then
			m_onErrorCallback(3)
		end
		return
	end
	
	accTryCount = accTryCount + 1
	
	local skey = GetKey()
	skey = G.Tool.Base64.decodeBase64(skey)
	if skey == "" then
		if isNeedCallback == true then
			m_onErrorCallback(4)
		end	
		return
	end
	
	initRes = CS.AliyunSDK.initEx(skey, GetToken())
	if isSimulation == true then
		initRes = simulationInitRes
		if isNeedCallback == true then
			initRes = 0
		end
	end
	
	if initRes ~= 0 then
		ShowLog("InitEx 失敗")
		
		InitEx(isNeedCallback)
	else
		initSuccess = true
		ShowLog("InitEx 成功")
		
		if isNeedCallback == true then
			m_onSuccessCallback()
		end
	end
end

function IsInitSuccess()
	return initSuccess
end

local isGet = false
function getProxyTcpByDomain(domain, port)
	print("######################## getProxyTcpByDomain Token", GetToken())
	print("######################## getProxyTcpByDomain GroupName", GetGroupName())
	print("######################## getProxyTcpByDomain domain", domain)
	print("######################## getProxyTcpByDomain port", port)

	if isSimulation == true then
		if isGet == false then
			isGet = true
			return simulationGetIpRes, simulationGetIpIp, simulationGetIpPort
		else
			return 0, simulationGetIpIp, simulationGetIpPort
		end
	end
	
	local data = CS.AliyunSDK.getProxyTcpByDomain(GetToken(), GetGroupName(), domain, tostring(port))
	local resList = G.Tool.StringHelper.Split(data, "\n")
	local res = tonumber(resList[1])
	if res ~= 0 then
		ShowLog("取得IP失敗" .. tostring(res))
		return 1, "", 0
	else
		local ip = resList[2]
		local port = tonumber(resList[3])
		ShowLog("取得IP成功 " .. ip .. ":" .. port)
		return res, ip, port
	end
end

function GetIP(port)
	local res, ip, port = getProxyTcpByDomain(GetMLobby(), tostring(port))
	ShowLog("取得IP " .. tostring(res) .. " " .. ip .. " " .. tostring(port))
	return res, ip, port
end

function GetCdnDomain()
	local UpdateDomain = ""
	if GetUseCdn() == 1 then
		local res, ip, port = getProxyTcpByDomain(GetMCdn(), tostring(GetPortCdn()))
		if res == 0 then
			UpdateDomain = "http://" .. ip .. ":" .. tostring(port) .. "/UpdateFiles"
		end
	end
	print("######################## GetCdnDomain", UpdateDomain)
	
	return UpdateDomain
end