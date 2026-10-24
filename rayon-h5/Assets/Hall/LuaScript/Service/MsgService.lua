require("Module")
mmodule(...)

local G = require "Global"

local groupName = "MsgService"

local MsgList = { }; 		-- 所有訊息列表
local ShowMsgList = { }; 	-- 顯示訊息列表

local MsgKey = "Msg"
local isShow = false
local RefreshTime = 5
local accIndex = -1

function ShowLog(data)
	--print(data)
	--CS.UnityEngine.Debug.LogWarning(data)
end

function Init()
	--CS.UnityEngine.PlayerPrefs.DeleteKey(MsgKey)
	
	CS.TimeManager.Instance:CancelInvokeRepeat(GetMsgFromServer)	
	CS.TimeManager.Instance:InvokeRepeating(GetMsgFromServer, 0, RefreshTime)
end

function GetMsgFromServer()
	print("$$$$$$$$$$$$$$$$$$ GetMsgFromServer")
end

function ShowMsg(Type, DirectShow)
	ShowLog("##################### ShowMsg " .. Type)
	
	if MsgList[Type] == nil then
		return
	end
	
	if DirectShow == true then
		for k, v in pairs(MsgList[Type]) do
			if v.isAddShow == false then		
				if v.isDirectShow == true then		
					v.isAddShow = true
					table.insert(ShowMsgList, v)
				end
			end
		end
	else
		for k, v in pairs(MsgList[Type]) do
			if v.isAddShow == false then		
				v.isAddShow = true
				table.insert(ShowMsgList, v)
			end
		end
	end
	
	DoShow()
end

function DoShow()
	ShowLog("##################### DoShow")

	if isShow == true then
		ShowLog("##################### isShow true")
		return
	end
	
    local message = GetOneMsg()
	if message == nil then
		ShowLog("##################### no data")
		return
	end
	
	ShowLog("show message index " .. tostring(message.index))
		
	isShow = true
	
	SaveMsg(message.type, message.index, message.userID)
	
	local openparam = CS.WindowNodeInitParam("UIMessageBoxNew", CS.AppDefine.HallNumber)
	openparam.NodeType = CS.BaseNodeType.Propmt;
	openparam.LoadComplatedCallBack = function(windowNode)
		local txt = windowNode.WindowGameObject.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")
		txt.text =  message.msg	
	end
	openparam.WindowCloseCallBack =
	(
	function()
		isShow = false
			
		DoShow()
	end
	)
	CS.WindowManager.Instance:OpenWindow(openparam)
end

function GetOneMsg()
    for k, v in pairs(ShowMsgList) do
		if v.isShow == false then
			v.isShow = true
			return v
		end
    end
	return nil
end

function GetMsgIndex(Type, UserID)
	local key = MsgKey .. "_" .. tostring(Type) ..  "_" ..tostring(UserID)
	local index = CS.UnityEngine.PlayerPrefs.GetInt(key,0)
	return index
end

function SaveMsg(Type, Index, UserID)
	ShowLog("##################### SaveMsg " .. Index)
	
	local nowSaveIndex = GetMsgIndex()
	if Index <= nowSaveIndex then
		ShowLog("##################### not SaveMsg, nowSaveIndex " .. nowSaveIndex .. " save" .. Index)	
		return
	end	
	
	-- 紀錄
	local key = MsgKey .. "_" .. tostring(Type) ..  "_" ..tostring(UserID)
	CS.UnityEngine.PlayerPrefs.SetInt(key, Index)
	CS.UnityEngine.PlayerPrefs.Save()
	
	-- 告知已顯示
	local sendMsg = { }
	sendMsg.Idx = Index
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.NotifyServer.UpdateHallUserListID, sendMsg)
end

function AddMsg(Type, Index, Msg, UserID, DirectShow)
	if G.User.GetCurrentUser() == nil then
		return
	end

	-- 判斷訊息是不是自己的
	if tostring(UserID) ~= tostring(G.User.GetCurrentUser().UserID) then
		ShowLog("!!!!!!!!!!! drop not self data")
		return
	end

	local nowIndex = GetMsgIndex(Type, UserID)
		
	-- 丟棄舊資料
	if nowIndex >= Index then
		ShowLog("!!!!!!!!!!! drop old data, now:" .. nowIndex .. ", new:" .. Index)
		return
	end	
	
	-- 避免重複加入
	if MsgList[Type] ~= nil then
		for k, v in pairs(MsgList[Type]) do
			if v.index == Index then
				ShowLog("##################### drop the same data" .. Index)
				return
			end			
		end
	end
	
	ShowLog("------------------------------------------------ AddMsg " .. nowIndex .. " " .. Index)
	
	local data = {}
	data.index = Index
	data.type = Type
	data.msg = Msg
	data.userID = UserID
	data.isAddShow = false
	data.isShow = false
	data.isDirectShow = DirectShow
		
	if MsgList[Type] == nil then
		MsgList[Type] = {}
	end
	
	table.insert(MsgList[Type], data)
end

function TestAddMsg()
	if accIndex == -1 then
		accIndex = GetMsgIndex(1, G.User.GetCurrentUser().UserID) + 1
	end
	
	print("############ TestAddMsg")
	AddMsg(1, accIndex, "Test\n" .. "index:" .. accIndex .. "\n", tostring(G.User.GetCurrentUser().UserID), true)
	
	accIndex = accIndex + 1
	
	-- 遊戲中不立刻顯示
	if G.IsGameStart() == true then
	
	-- 在大廳立刻顯示
	else
		ShowMsg(1, false)
	end	
end

G.NetQ.RegMsgHook(G.Proto.BClass.NotifyServer, G.Proto.NotifyServer.UserNotifyRetH2CID, function(receiveMsg)
	-- 在大廳立刻顯示
	if G.IsGameStart() == false then
		AddMsg(receiveMsg.Type, receiveMsg.Idx, receiveMsg.Message, receiveMsg.UserID, true)
        if(tostring(receiveMsg.Type)=="1")then
            G.Tool.Debug.RecordGet({
                                content= {  
                                            ContentType=2,
                                            Flag=7,
                                            Data= tostring(receiveMsg.Message),
                                            id=receiveMsg.Idx
                                         },
                                lt=1,
                                sceneName="UIHALL",
                                uid=G.User.GetCurrentUser().UserID
                                })
		end
		ShowMsg(1, true)
	else
		local TimeSeq_Msg = receiveMsg.TimeSeq
		local TimeSeq_Game = G.Room.GetUserData().TimeSeq
		ShowLog("#################### TimeSeq_Msg " .. TimeSeq_Msg)
		ShowLog("#################### TimeSeq_Game " .. TimeSeq_Game)
		
		-- 讯息时间比进入游戏时间小,就直接显示
		if TimeSeq_Msg < TimeSeq_Game then
			ShowLog("@@@@@@@@@@@@@@@@@@@ Direct Show ")
			AddMsg(receiveMsg.Type, receiveMsg.Idx, receiveMsg.Message, receiveMsg.UserID, true)	
             if(tostring(receiveMsg.Type)=="1")then
            G.Tool.Debug.RecordGet({
                                content= {  
                                            ContentType=2,
                                            Flag=7,
                                            Data= tostring(receiveMsg.Message),
                                            id=receiveMsg.Idx
                                         },
                                lt=1,
                                sceneName="UIHALL",
                                uid=G.User.GetCurrentUser().UserID
                                })
		end
			ShowMsg(1, true)
		else
			ShowLog("@@@@@@@@@@@@@@@@@@@ Not Show")
			AddMsg(receiveMsg.Type, receiveMsg.Idx, receiveMsg.Message, receiveMsg.UserID, false)
		end
	end
	
	
end , groupName);
