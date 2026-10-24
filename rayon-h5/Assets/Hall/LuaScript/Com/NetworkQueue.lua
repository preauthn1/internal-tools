require("Module")
mmodule(...)
local proto = require("Com/Proto")
local protobuf = require "pb"
local G = require "Global"

-- 傳送消息Log開關
local boolean showSendLog = true
function SetShowSendLog(IsShow)
	showSendLog = IsShow
end

-- 消息解鎖狀態Log開關
local boolean showSendPauseLog = false
function SetShowSendPauseLog(IsShow)
	showSendPauseLog = IsShow
end

-- 接收消息Log開關
local boolean showReceiveLog = true
function SetShowReceiveLog(IsShow)
	showReceiveLog = IsShow
end

-- 刪除消息Log開關
local boolean showDeleteLog = false
function SetShowDeleteLog(IsShow)
	showDeleteLog = IsShow
end

-- 消息拦截集合
local msgHookList = { } 
-- 消息结构映射
local msgModelNames = { }


-- 大厅连接断开回调
OnHallSocketDisconnectCallback = nil

-- 大厅未登录时连接断开重连
OnHallReConnectCheck = nil

-- 游戏连接断开回调
OnGameSocketDisconnectCallback = nil


-- 获取消息协议映射的key值
local function getMsgModelKey(mainId, secondId)
    return mainId .. "." .. secondId
end

-- 解析网络消息
local function decodeNetMsgData(mainId, secondId, data)
    if data == nil then
        return nil
    end

    -- 获取协议映射到的结构
    local mkey = getMsgModelKey(mainId, secondId)
    local mm = msgModelNames[mkey]

    --[[
	if (G.GetRunGameID() ~= 14) then
        print("解析实体", mm)
    end
	]]--

    if mm ~= nil then
        local decodeData, errmsg = protobuf.decode(mm, data)

        if decodeData == false and errmsg ~= nil then
            print("<color=#ff0000ff>!!!!!!!!!!!!!!!!!!!!!!!!!!! 解析协议结构错误</color>：" .. errmsg)
            return nil
        else
            return decodeData
        end
	else	
		if data == nil or data == "" then
	
		else
			print("<color=#ff0000ff>!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! 消息结构映射错误,未加AddMsgModelName</color>," .. 
					" Type:<color=#ff00ffff>" .. mainId .. " " ..CS.Protocol.GetTypeString(mainId).."</color>"..
					" ProtoId:<color=#ff00ffff>" .. secondId .. " " ..CS.Protocol.GetProtocolString(mainId, secondId).."</color>")
		end
    end

    return data
end

-- 当收到网络消息时被调用，arg1 为主ID，arg2为副ID, arg3为消息内容
local function onReceiveNetworkMessage(mainId, secondId, content)

    --[[
	if (G.GetRunGameID() ~= 14) then
        print("recv msg." .. mainId .. "." .. secondId)
    end
	]]--

    -- 大厅断开
    if mainId == 0 and secondId == 0 then
		--3秒後重連
		CS.TimeManager.Instance:DelayInvoke(3,
		function()
			G.Service.NetworkService.ReLogonHall()
		end )
        return
    end

    -- 游戏断开
    if mainId == 0 and secondId == 1 then
        if G.IsGameStart() then
            G.Service.NetworkService.ReLoginRoom()
        end

        return
    end

    local onceHookList = { }

    -- for k, v in pairs(msgHookList) do
    for k = #msgHookList, 1, -1 do
        local v = msgHookList[k]

        -- 触发符合协议的接收对象
        if (v.mainId == mainId and v.secondId == secondId) then
            -- print("消息", v.mainId, v.secondId, v.isOne)
            if (v.isOnce == true) then
                table.insert(onceHookList, v.callfun)
            end
            if (type(content) ~= "table") then
                content = decodeNetMsgData(mainId, secondId, content)
                if (showReceiveLog == true) then
					if content ~= nil then
						ShowLog("收到网络消息", "<color=#ff0000ff>", mainId, secondId, content)
					else
						ShowLog("收到网络消息Error!!", "<color=#ff0000ff>", mainId, secondId, nil)					
					end
                end
            end
            v.callfun(content)
        end
    end

    -- 删除侦听一次的消息
    RemoveMsgHookByCallFuncList(onceHookList)
end

-- 注册消息拦截(新,檢查是否重複加入)
local function regMsgHookInternalNew(mainId, secondId, callfun, isOnce, groupName)
    -- 消息拦截信息
	if isOnce == true then
	    for k, v in pairs(msgHookList) do
			if v.mainId == mainId and v.secondId == secondId then
				v.callfun = callfun
				return v
			end
		end
	end
	
	local msgHookInfo = { mainId = mainId, secondId = secondId, callfun = callfun, isOnce = isOnce, groupName = groupName }
	table.insert(msgHookList, msgHookInfo)	
	
    return msgHookInfo
end

-- 注册消息拦截
local function regMsgHookInternal(mainId, secondId, callfun, isOnce, groupName)
    -- 消息拦截信息
    local msgHookInfo = { mainId = mainId, secondId = secondId, callfun = callfun, isOnce = isOnce, groupName = groupName }
	
    table.insert(msgHookList, msgHookInfo)
	
    return msgHookInfo
end

-- 注册多次侦听的消息拦截
function RegMsgHook(mainId, secondId, callfunc, groupName)
    regMsgHookInternal(mainId, secondId, callfunc, false, groupName)
end

-- 注册单次侦听的消息拦截
function RegOnceMsgHook(mainId, secondId, callfunc, groupName)
    regMsgHookInternal(mainId, secondId, callfunc, true, groupName)
end

-- 注册单次侦听的消息拦截(新,檢查是否重複加入)
function RegOnceMsgHookNew(mainId, secondId, callfunc, groupName)
    regMsgHookInternalNew(mainId, secondId, callfunc, true, groupName)
end

-- 移除所有游戏消息拦截
function RemoveGameMsgHook()
    local gameHookList = { }

    for k = #msgHookList, 1, -1 do
        v = msgHookList[k]
        if (v.mainId == proto.BClass.Game) then
            table.insert(gameHookList, v.callfun)
        end
    end

    RemoveMsgHookByCallFuncList(gameHookList)
end

-- 根据处理函数移除消息拦截
function RemoveMsgHookByCallFunc(callfun)
    local callfunIdx = 0

    for k, v in pairs(msgHookList) do
        if (callfun == v.callfun) then
            callfunIdx = k
            if (showDeleteLog == true) then
                print("删除消息", v.mainId, v.secondId, "剩余消息数", #msgHookList)
            end
            break
        end
    end

    if callfunIdx > 0 then
        table.remove(msgHookList, callfunIdx)
    end
end

-- 移除一组处理函数
function RemoveMsgHookByCallFuncList(removelist)
    for _, v in pairs(removelist) do
        RemoveMsgHookByCallFunc(v)
    end
end

-- 根据组名移除消息拦截
function RemoveMsgHookByGroup(groupName)
    local hookGroupList = { }

    for k, v in pairs(msgHookList) do
        if v.groupName ~= nil and v.groupName ~= "" and v.groupName == groupName then
            table.insert(hookGroupList, v.callfun)
        end
    end


    RemoveMsgHookByCallFuncList(hookGroupList)
end


-- 暂停处理消息
function Pause()
    CS.NetworkManager.Instance:SetReceiveMessagePause(true)
end

-- 恢复处理消息
function Resume()
    CS.NetworkManager.Instance:SetReceiveMessagePause(false)
end

-- 获取游戏消息拦截数量
function GetGameMsgHookLen()
    local count = 0
    for _, v in pairs(msgHookList) do
        if (v.mainId == proto.BClass.Game) then
            count = count + 1
        end
    end
    return count
end

-- 网络消息编码
local function encodeNetMsgData(mainId, secondId, data)
    if data == nil then
        return nil
    end

	--[[
    if (G.GetRunGameID() ~= 14) then
        print("发送消息 " .. mainId .. "." .. secondId)
    end
	]]--

    local mkey = getMsgModelKey(mainId, secondId)

    local dataName = msgModelNames[mkey]
    if dataName ~= nil and #dataName > 0 then
        data = protobuf.encode(dataName, data)
	else
		print("<color=#ff0000ff>!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! 消息结构映射错误,未加AddMsgModelName</color>," .. 
				" Type:<color=#ff00ffff>" .. mainId .. " " ..CS.Protocol.GetTypeString(mainId).."</color>"..
				" ProtoId:<color=#ff00ffff>" .. secondId .. " " ..CS.Protocol.GetProtocolString(mainId, secondId).."</color>")	
    end

    return data
end

function ShowLog(text, color, mainId, secondId, data)
    local logText = ""..color..text.."</color>"..
					" Type:"..color..mainId.." " ..CS.Protocol.GetTypeString(mainId).."</color>"..
					" ProtoId:"..color..secondId.." " ..CS.Protocol.GetProtocolString(mainId, secondId).."</color>"
    if (data ~= nil) then
        if (type(data) == "table") then
            G.Tool.TableHelper.Print(logText, data)
        else
			print(logText.."\n".. data)
        end
    else
		print(logText.."\n nil")
    end
end

-- 发送网络消息
-- 自动判断发往大厅服务器还是游戏服务器
function SendMsg(mainId, secondId, data)

	if (showSendPauseLog == true) then
	    if (G.GetRunGameID() ~= 14) then
	        print("发消息锁状态", mainId, secondId, CS.NetworkManager.Instance:GetPause())
	    end
	end

    -- 发送大厅消息
    if (mainId == proto.BClass.Hall) then
		if data ~= nil and data.Area_ID == nil then
			data.Area_ID = G.GameConfig.GameSiteID
		end
	
        SendHallMsg(mainId, secondId, data)
        -- 发送游戏消息
    else
        SendGameMsg(mainId, secondId, data)
    end
end
 
-- 向大厅服务器发送消息
function SendHallMsg(mainId, secondId, data)
    if CS.NetworkManager.Instance.HallSocket.MyConnectState ~= CS.ConnectState.Connect then
        G.Service.NetworkService.ReLogonHall()
    else
        if (showSendLog == true) then
            ShowLog("向大厅服务器发送消息", "<color=#ff00ffff>", mainId, secondId, data)
        end

        data = encodeNetMsgData(mainId, secondId, data)
        CS.NetworkManager.Instance:SendHallMessage(mainId, secondId, data)
    end
end

-- 向游戏服务器发送消息
function SendGameMsg(mainId, secondId, data)
    if (showSendLog == true) then
        ShowLog("向游戏服务器发送消息", "<color=#ff00ffff>", mainId, secondId, data)
    end

    data = encodeNetMsgData(mainId, secondId, data)
    CS.NetworkManager.Instance:SendRoomMessage(mainId, secondId, data)
end




-- 添加返回消息协议解析结构
function AddReturnMsgModelName(mainId, secondId, name)
    AddMsgModelName(mainId, secondId, name)
end


-- 添加消息协议解析结构
-- 重新更贴切的命名函数
-- 由于AddReturnMsgModelName已在多处使用,继续保留,以后应该使用此方法来添加协议和实体的映射
function AddMsgModelName(mainId, secondId, name)
    local mkey = getMsgModelKey(mainId, secondId)
    msgModelNames[mkey] = "netproto." .. name
end

-- 初始化消息映射
local function initMsgModel()
    -- Common
    AddMsgModelName(proto.BClass.Common, proto.Common.IPReportID, "IPReport")
    AddMsgModelName(proto.BClass.Common, proto.Common.HeartBeatConfigID, "HeartBeatConfig")
    AddMsgModelName(proto.BClass.Common, proto.Common.TipMessageID, "TipMessage")
    AddMsgModelName(proto.BClass.Common, proto.Common.ConnTokenID, "ConnTokenMessage")
   
   
	
    -- 通知
    AddMsgModelName(proto.BClass.NotifyServer, proto.NotifyServer.AnnID, "AnnList")
    AddMsgModelName(proto.BClass.NotifyServer, proto.NotifyServer.AttrChangeID, "AttrChangeList")
    AddMsgModelName(proto.BClass.NotifyServer, proto.NotifyServer.UserNotifyRetH2CID, "UserNotifyInfo")

    --  Hall
	-- 訊息
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UpdateUserNotifyIdxC2HID, "ClientUpdateUserNotifyIdx")
	
     --代理訪問(防止網址被取得)
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ProxyRequestRetID, "ProxyData")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ProxyRequestID, "ProxyData")
     --代理訪問POST(防止網址被取得)
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ProxyPOSTRequestRetID, "ProxyData")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ProxyPOSTRequestID, "ProxyPOSTData")

	--比赛
	AddMsgModelName(proto.BClass.Hall, proto.Hall.GameCompetitionID, "GetGameCompetitionMsg")
	AddMsgModelName(proto.BClass.Hall, proto.Hall.GameCompetitionRetID, "GameCompetitionMsg")
	 
	AddMsgModelName(proto.BClass.Hall, proto.Hall.GameCompetioniInfoRuleID, "GetGameCompetitionRuleMsg")
	AddMsgModelName(proto.BClass.Hall, proto.Hall.GameCompetioniInfoRuleRetID, "GameCompetitionRuleMsg")	

	--buy ticket
	AddMsgModelName(proto.BClass.Hall, proto.Hall.BuyTicketID, "GetBuyTicketMsg")
	AddMsgModelName(proto.BClass.Hall, proto.Hall.BuyTicketRetID, "BuyTicketMsg")


    -- 登录
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FindSetPwdByPhoneID, "FindSetPwdByPhone")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FindSetPwdByPhoneRetID, "RetMessage")

    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestServerListID, "GetShowGameFlagByAreaMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ServerListDataID, "AllGameServerInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GameServerAddrInfoID, "GameServerAddr")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestHallVerID, "HallVerReq")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.HallVerInfoID, "HallVerInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GuestLoginID, "GuestLogin")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserLoginID, "UserLogin")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.LoginRetID, "UserLoginRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestGameServerAddrID, "RequestGameServerAddrInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserLogoutID, "UserLogout")

    -- 刷新用户信息
    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestUserHallInfoID, "RequestUserHallInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserHallInfoID, "UserLoginRet")

    -- 大厅->个人信息
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindZhifubaoID, "ZhifubaoInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindZhifubaoRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyPasswordID, "ModifyPassword")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyPasswordRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SendPhoneVCodeID, "SendPhoneVCode")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SendPhoneVCodeRetID, "SendPhoneVCodeRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindGuestAccountID, "BindGuestAccount")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindGuestAccountRetID, "RetMessage")

    -- 大厅->个人信息->修改头像
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyFaceID, "ModifyFace")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyFaceRetID, "RetMessage")

    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyAvatarFrameID, "ModifyAvatarFrameMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyAvatarFrameRetID, "ModifyAvatarFrameRetMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetGradeSettingDataID, "GradeSettingDataList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetGradeSettingDataRetID, "GradeSettingDataList")
	
    -- 大厅->银行
    AddMsgModelName(proto.BClass.Hall, proto.Hall.DepositMoneyID, "MoneyDeposit")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.DepositMoneyRetID, "MoneyDepositRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyBankPasswordID, "ModifyBankPassword")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ModifyBankPasswordRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FindSetBankPwdByPhoneID, "FindSetBankPwdByPhone")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FindSetBankPwdByPhoneRetID, "RetMessage")

    -- 大厅->邮件
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMailListRetID, "MailList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMailDetailID, "RequestMailDetail")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMailDetailRetID, "MailDetail")

   
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMailListRetID, "MailList")

    -- 大厅->联系客服
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ContactServiceID, "ContactService")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ContactServiceRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CheckCSNotifyRetID, "CSNotiftyCount")

    -- 大厅-游戏充值->订单申请
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CreatePayOrderID, "CreatePayOrder")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CreatePayOrderRetID, "CreatePayOrderRet")

    -- 大厅-游戏充值->代理列表
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetAgentListID, "RequestAgentList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetAgentListRetID, "AgentList")

    -- 大厅-游戏充值->代理申请
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ApplyAgentID, "ApplyAgent")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ApplyAgentRetID, "RetMessage")

    -- 大厅-游戏充值->举报代理
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ReportAgentID, "ReportAgent")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ReportAgentRetID, "RetMessage")

    -- 大厅->兑换
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ConvertMoneyID, "ConvertMoney")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ConvertMoneyRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindBankCardID, "BankCardInfo")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.BindBankCardRetID, "RetMessage")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetConvertLogID, "GetConvertLog")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetConvertLogRetID, "ConvertLogList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetBankListRetID, "BankList")

    -- 大厅->FAQ
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetFAQID, "GetFAQ")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetFAQRetID, "FAQDetail")

    -- 大厅->代理
    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestAgentListRetID, "AgentData")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.JoinAgentListID, "JoinAgentList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.JoinAgentListRetID, "JoinAgentList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ApplyAgentListRetID, "ApplyAgentList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CheckAgentRetID, "ApplyAgentList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ApplyAgentListID, "RequestApplyAgent")

    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserQACodeID, "GetUserQRCodeMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserQACodeRetID, "GetUserQRCodeRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentExecID, "GetUserAgentExecMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentExecRetID, "GetUserAgentExecRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentTaxLogID, "GetUserAgentTaxLogMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentTaxLogRetID, "GetUserAgentTaxLogRetList")
	
	-- 大厅->代理>Promoter 全民
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetWithdrawListID, "GetWithdrawMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetWithdrawListRetID, "GetWithdrawRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetWithdrawOrderID, "GetWithdrawOrderMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetWithdrawOrderRetID, "GetWithdrawOrderRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetAgentTaxListRetID, "GetAgentTaxListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetAgentTaxBonusLogID, "GetAgentTaxBonusLogMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetAgentTaxBonusLogRetID, "GetAgentTaxBonusLogRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentListID, "GetUserAgentListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetUserAgentListRetID, "GetUserAgentListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetShanFuID, "GetShanFuMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetShanFuRetID, "GetShanFuRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetShanFuChatID, "GetShanFuChatMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetShanFuChatRetID, "GetShanFuChatRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentAddUserRelationID, "AddUserRelationMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentAddUserRelationRetID, "AddUserRelationRet")
	
    ---- 大厅->代理>lb 裂變
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyInfoID, "SplitGetMyInfoMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyInfoRetID, "SplitGetMyInfoRet")
 
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyConvertID, "SplitGetMyConvertMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyConvertRetID, "SplitGetMyConvertRet")

    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetRewardTipsID, "SplitGetRewardTipsMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetRewardTipsRetID, "SplitGetRewardTipsRet")    
    
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyRewardID, "SplitGetMyRewardMsg")    
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitGetMyRewardRetID, "SplitGetMyRewardRet")   

    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitConvertOrderID, "SplitConvertOrderMsg")    
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitConvertOrderRetID, "SplitConvertOrderRet")

    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitBindingInvitedCodeID, "SplitBindingInvitedCodeMsg")    
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitBindingInvitedCodeRetID, "SplitBindingInvitedCodeRet")   

    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitAddUserRelationID, "AddUserRelationMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.SplitAddUserRelationRetID, "AddUserRelationRet")
	
	-- 大厅->代理>NetProfit 淨利潤
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentInfoID, "NetProfitAgentGetAgentInfoMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentInfoRetID, "NetProfitAgentGetAgentInfoRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetSubUserListID, "NetProfitAgentGetSubUserListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetSubUserListRetID, "NetProfitAgentGetSubUserListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentWithdrawAmountID, "NetProfitAgentGetAgentWithdrawAmountMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentWithdrawAmountRetID, "NetProfitAgentGetAgentWithdrawAmountRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetWithdrawAmountID, "NetProfitAgentGetWithdrawAmountMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetWithdrawAmountRetID, "NetProfitAgentGetWithdrawAmountRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetBonusAmountID, "NetProfitAgentGetBonusAmountMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetBonusAmountRetID, "NetProfitAgentGetBonusAmountRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentWithDrawOrderID, "NetProfitAgentWithDrawOrderMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentWithDrawOrderRetID, "NetProfitAgentWithDrawOrderMsgRet")
    --AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentApplyListID, "NetProfitAgentApplyListMsg")
    --AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentApplyListRetID, "NetProfitAgentApplyListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetWithdrawDetailID, "NetProfitAgentGetWithdrawDetailMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetWithdrawDetailRetID, "NetProfitAgentGetWithdrawDetailRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentDetailID, "NetProfitAgentGetAgentDetailMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentDetailRetID, "NetProfitAgentGetAgentDetailRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentReportLogID, "NetProfitAgentGetAgentReportLogMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentGetAgentReportLogRetID, "NetProfitAgentGetAgentReportLogRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentAddUserRelationID, "AddUserRelationMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.NetProfitAgentAddUserRelationRetID, "AddUserRelationRet")
	
	-- 大厅->代理>WXD 無限級
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithAmountID, "UserAgentGetWithAmountMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithAmountRetID, "UserAgentGetWithAmountRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithAmountListID, "UserAgentGetWithAmountListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithAmountListRetID, "UserAgentGetWithAmountListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithdrawAmountListID, "UserAgentGetWithdrawAmountListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentGetWithdrawAmountListRetID, "UserAgentGetWithdrawAmountListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentResultListID, "UserAgentResultListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentResultListRetID, "UserAgentResultListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentBorrowListID, "UserAgentBorrowListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentBorrowListRetID, "UserAgentBorrowListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentWinAmountListID, "UserAgentWinAmountListMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentWinAmountListRetID, "UserAgentWinAmountListRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMyPromotionID, "GetMyPromotionMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetMyPromotionRetID, "GetMyPromotionRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetRewardTipsID, "GetRewardTipsMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetRewardTipsRetID, "GetRewardTipsRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitConvertOrderID, "UnlimitConvertOrderMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitConvertOrderRetID, "UnlimitConvertOrderRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitPreConvertOrderID, "UnlimitPreConvertOrderMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitPreConvertOrderRetID, "UnlimitPreConvertOrderRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentUserModeGetID, "UserAgentUserModeGetMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UserAgentUserModeGetRetID, "UserAgentUserModeGetRet")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitAddUserRelationID, "AddUserRelationMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.UnlimitAddUserRelationRetID, "AddUserRelationRet")
	
    -- GameRoom
    AddMsgModelName(proto.BClass.GameRoom, proto.GameRoom.GameVerInfoID, "GameVerInfo")
    AddMsgModelName(proto.BClass.GameRoom, proto.GameRoom.LoginRoomID, "LoginGameRoomInfo")
    AddMsgModelName(proto.BClass.GameRoom, proto.GameRoom.LoginRoomRetID, "LoginGameRoomRet")
	AddMsgModelName(proto.BClass.GameRoom, proto.Hall.GameCompetitionID, "GetGameCompetitionMsg")
	AddMsgModelName(proto.BClass.GameRoom, proto.Hall.GameCompetitionRetID, "GameCompetitionMsg")

    -- Game
    AddMsgModelName(proto.BClass.Game, proto.Game.UserStateChangeID, "UserStateChange")
    AddMsgModelName(proto.BClass.Game, proto.Game.UserEnterID, "UserEnter")
    AddMsgModelName(proto.BClass.Game, proto.Game.UserLeaveID, "UserLeave")
    AddMsgModelName(proto.BClass.Game, proto.Game.UserPokerID, "UserPoker")
    AddMsgModelName(proto.BClass.Game, proto.Game.UserActionCommitID, "UserActionCommited")
    AddMsgModelName(proto.BClass.Game, proto.Game.AllUserActionCommitID, "AllUserActionCommited")
    AddMsgModelName(proto.BClass.Game, proto.Game.LookOnUserNotifyID, "LookOnUserNotify")	
	
    AddMsgModelName(proto.BClass.Hall, proto.Hall.RequestLoginBonusRetID, "LoginBonusData")
	AddMsgModelName(proto.BClass.Game, proto.Game.GameBeginCountDownID, "GameBeginCountDown")

    -- 奖池回覆
    AddMsgModelName(proto.BClass.Hall, proto.Hall.GetJackpotPoolRetID, "JackpotPoolRet")    

	-- 禮碼
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ExchangeCouponID, "ExchangeCouponMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.ExchangeCouponRetID, "ExchangeCouponRetMsg")
	
	--活動
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CampaignGetListID, "CampaignGetList")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.CampaignGetListRetID, "CampaignGetListRet")
	
	-- 首充返利
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FirstDepositGetInfoID, "FirstDepositGetInfoMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FirstDepositGetInfoRetID, "FirstDepositGetInfoRetMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FirstDepositCollectID, "FirstDepositCollectMsg")
    AddMsgModelName(proto.BClass.Hall, proto.Hall.FirstDepositCollectRetID, "FirstDepositCollectRetMsg")
	
end


-- 初始化消息拦截
local function initMsgHook()
    -- 注册网络消息
    CS.NetworkManager.Instance.OnReceiveNetworkMessage = CS.NetworkManager.Instance.OnReceiveNetworkMessage + onReceiveNetworkMessage

end



-- 初始化
function Init()
    -- 初始化消息协议结构映射
    initMsgModel()

    -- 初始化事件侦听
    initMsgHook()
end

function Clear()
    -- 移除网络消息侦听
    CS.NetworkManager.Instance.OnReceiveNetworkMessage = CS.NetworkManager.Instance.OnReceiveNetworkMessage - onReceiveNetworkMessage
end


-- return
-- {
--    OnReceiveNetworkMessage = OnReceiveNetworkMessage,
--    RegMsg = RegMsg,
--    Pause = Pause,
--    Resume = Resume,
--    Remove = Remove,
--    InitReturnMsgModel = InitReturnMsgModel,
--    AddReturnMsgModelName = AddReturnMsgModelName,
--    SendMsg = SendMsg,
--    RemoveGameMsg = RemoveGameMsg,
--    GetMsgLen = GetMsgLen,
--    RemoveList = RemoveList
-- }