require("Module")
mmodule(...)

local G = require "Global"

local HallVer = {}
HallVer["CallbackHotUpdateFail"] = "10.3.0"

local GameVer = {}

-- 支援版號(checkVer => ex.10.0.0)
local function CompareHall(checkVer)
	local currentver = G.Env.GetHallVersion()
	local arr1 = G.Tool.StringHelper.Split(currentver,".");
	local arr2 = G.Tool.StringHelper.Split(checkVer,".");
	
	if (3 < #arr2 and #arr2 < 1) then
		return nil
	end
	
	if(arr1[1] > arr2[1])then
		return true	
	elseif (arr1[1] == arr2[1]) and (arr1[2] > arr2[2]) then
		return true
	elseif (arr1[1] == arr2[1]) and (arr1[2] == arr2[2]) and (arr1[3] >= arr2[3]) then
		return true
	end
	
	return false
end

function CheckHall(checkKey)
	local ver = HallVer[tostring(checkKey)]
	return CompareHall(ver)
end


-- 支援版號(checkVer => ex.1.0.0)
local function CompareGame(checkVer,gameId)
	local currentver = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId));
	local arr1 = G.Tool.StringHelper.Split(currentver,".");
	local arr2 = G.Tool.StringHelper.Split(checkVer,".");
	if (3 < #arr2 and #arr2 < 1) then
		return nil
	end
	
	if(arr1[1] > arr2[1])then
		return true
	elseif (arr1[2] > arr2[2]) then
		return true
	elseif (arr1[3] >= arr2[3]) then
		return true
	end
	return false
end

function CheckGame(checkKey,gameId)
	local ver = GameVer[tostring(checkKey)]
	return CompareGame(ver,gameId)
end