require("Module")
mmodule(...)

local G = require "Global"
local groupName = "LoginBonusService";

LoginTimes = 0 --己领取次数
RetCode=-1 --领取代码 0 本日己领取, 1 成功领取, 2可领取, -1 活动不存在, -2 已领超过七次限制, -3 充值金额不足, -4 帐号未绑定
GiftCash = 0
GiftCashToCash = 0
DayMax = 7
ContentList = {}
ContentList[1] = 50
ContentList[2] = 50
ContentList[3] = 50
ContentList[4] = 100
ContentList[5] = 100
ContentList[6] = 100
ContentList[7] = 200

--查询登入奖励
function SendRequest(callback)
	Callback = callback
	local param = {}
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestLoginBonusID, param)
end

--领取登入奖励
function SendApply(callback)
	Callback = callback
	local param = {}
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ApplyLoginBonusID, param)
end

function IsOpenLoginBonus()
	if ((RetCode == -1) or (RetCode == -2)) then
		return false
	else
		return true
	end
end

-- 领取代码 0本日己领取,1成功领取,2未領取,-1活动不存在,-2已领超过七次限制,-3充值金额不足,-4帐号未绑定
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.RequestLoginBonusRetID, function(receiveMsg)
	LoginTimes=receiveMsg.LoginTimes
    RetCode=receiveMsg.RetCode
	GiftCashToCash= G.Tool.UserHelper.ToCash(receiveMsg.GiftCash)
	GiftCash=receiveMsg.GiftCash
	if(receiveMsg.BonusDayData ~= nil and #receiveMsg.BonusDayData > 0)then
		ContentList = receiveMsg.BonusDayData
	end
    local userdata = G.User.GetCurrentUser()
	userdata.CashAmount = receiveMsg.CurrentAmount
	
	if Callback ~= nil then
		Callback()
		Callback = nil
	end
end , groupName);