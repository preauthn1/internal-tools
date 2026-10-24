require("Module")
mmodule(...)

local G = require "Global"

-- 是否支援
function IsHenrySDK()
	-- 編輯器下
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor or 
		CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor	 then
		return false
	else
		if CS.HenrySDK.s_IsHenrySDK == true then
			return true
		else
			return false
		end
	end
end

function showMessageChat(url, appKey, appId, userId, userName, deviceId, timestamp, agentId, sign, key)
	if IsHenrySDK() == false then
		print("############## Henry not suppost Henry")
		return
	end

	print("############## Henry showMessageChat")
	CS.HenrySDK.showMessageChat(url, appKey, appId, userId, userName, deviceId, timestamp, agentId, sign, "0")
	
end