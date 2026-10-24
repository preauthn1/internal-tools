require("Module")
mmodule(...)

local G = require "Global"
local groupName = "FirstRechargeRebateService";

local data = nil

-- 查询充值返利
function SendRequest(callback)
	Callback = callback
	
	local sendMsg = {}
	sendMsg.UserID = G.User.GetCurrentUser().UserID	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.FirstDepositGetInfoID, sendMsg)
end

-- 查询结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.FirstDepositGetInfoRetID, function(receiveMsg)	
	data = receiveMsg

	if Callback ~= nil then
		Callback()
		Callback = nil
	end
end , groupName);

function GetData()
	return data
end