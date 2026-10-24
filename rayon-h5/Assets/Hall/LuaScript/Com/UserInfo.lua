require("Module")
mmodule(...)

local G = require("Global")
local currentuser = nil
local currentUserState = nil
local payinfo=nil

local userStateValue_Free = 0      -- 站立
local userStateValue_Queue = 1     -- 排队
local userStateValue_Sit = 2       -- 坐下
local userStateValue_Ready = 3     -- 准备
local userStateValue_Game = 4      -- 游戏中
local userStateValue_Offline = 5   -- 断线

function Clear()
    CS.PrefsUtility.DeleteKey("UserID")
    CS.PrefsUtility.DeleteKey("Cer")
    CS.PrefsUtility.DeleteKey("LevelKey")
    CS.PrefsUtility.DeleteKey("UserLevel")
    CS.PrefsUtility.DeleteKey("Token")
    currentuser = nil
end

function table.contains(table, element)
  for _, value in pairs(table) do
    if value == element then
      return true
    end
  end
  return false
end

-- 设置本地用户信息
function SetLocalUserInfo(user)
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsPlayer then
		if (CS.AppDefine.useDeviceUniqueIdentifier== true) then	
			return
		end
	end

    CS.PrefsUtility.SetInt("UserID", user.UserID)
    CS.PrefsUtility.SetString("Cer", user.Cer)
    CS.PrefsUtility.SetInt("UserLevel", user.UserData.UserLevel)
    CS.PrefsUtility.SetString("LevelKey", user.UserData.LevelKey)
end

-- 设置用户信息
function SetCurrentUserInfo(user)
	for i=#user.UserData.Banner,1,-1 do
		--0:複製,1:直連該網頁,2:一鍵入廳,3:甚麼都不做,4:前往充值,5:评分
		if user.UserData.Banner[i].Behavior == 5 then
			--获取设备类型 1.Windows,2.Android,3.IOS
			if G.Env.GetHDType() ~= 3 then
				table.remove(user.UserData.Banner,i)
			end
		end
	end
    currentuser = user
    SetLocalUserInfo(user)
end

-- 获取用户信息
function GetCurrentUser()
    return currentuser
end

-- 获取当前用户ID
function GetCurrentUserID()
    return currentuser.UserID
end

-- 获取当前用户证书
function GetCurrentUserCer()
    return currentuser.Cer
end

-- 获取当前用户类型 1.正式 2.機器人 3.遊客
function GetCurrentUserType()
    return currentuser.UserData.UserType
end

-- 当前用户是否游客
function IsCurrentUserGuest()
    return GetCurrentUserType() == 3
end

-- 获取本地用户信息
function GetLocalUserInfo()
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsPlayer then
		if (CS.AppDefine.useDeviceUniqueIdentifier == true) then
			return nil
		end
	end
	
    local userID = CS.PrefsUtility.GetInt("UserID")
    local cer = CS.PrefsUtility.GetString("Cer")
    local userLevel = CS.PrefsUtility.GetInt("UserLevel")
    local levelKey = CS.PrefsUtility.GetString("LevelKey")
	
    if userID == 0 or cer == "" then
        return nil
    end

    return { UserID = userID, Cer = cer, UserLevel = userLevel, LevelKey = levelKey }
end

-- 用户是否已经登录
function IsLogon()
    return currentuser ~= nil
end


-- 检查是否显示指定的游戏
function CheckUserGame(gameid)
    if currentuser == nil then
        return false
    end

    if currentuser.UserData == nil then
        return false
    end

    local gamelist = currentuser.UserData.GameList

    for _, gid in pairs(gamelist) do
        if gid == gameid then
            return true
        end
    end

    return false
end

function GetUserGameList()
    return currentuser.UserData.GameList
end
function GetUserGameList2()
   -- return nil;
     return currentuser.UserData.GameList2
end

function GetAdList()
   -- return nil;
     return currentuser.UserData.Banner
end

-- 设置用户状态
function SetUserState(state)
    currentUserState = state
end

-- 当前用户是游戏状态
function IsCurrentUserGameing()
    if currentUserState == nil then
        return false
    end

    return currentUserState.StateValue == userStateValue_Game or currentUserState.StateValue == userStateValue_Offline
end

-- 当前用户是否旁观状态
function IsCurrentUserLookOn()
    return not IsCurrentUserGameing()
end

-- 当前用户是否准备状态
function IsCurrentUserReady()
    if currentUserState == nil then
        return false
    end

    return currentUserState.StateValue == userStateValue_Ready
end

-- 当期用户是否坐下状态
function IsCurrentUserSit()
    if currentUserState == nil then
        return false
    end

    return currentUserState.StateValue == userStateValue_Sit
end

-- 当前用户桌号
function CurrentUserTableNo()
    if currentUserState == nil then
        return -1
    end

    return currentUserState.TableNo
end

-- 当前用户椅子号
function CurrentUserSeatNo()
    if currentUserState == nil then
        return -1
    end

    return currentUserState.SeatNo
end

-- 设置充值金额配置 
function SetPayAmountConfig(data)
    local list = G.Tool.StringHelper.Split(data, "_")
	if list ~= nil then
		print("+++++++++++++++++++", list[1], tostring(G.GameConfig.GamePlatformID))
		if (list[1] == tostring(G.GameConfig.GamePlatformID)) then
			currentuser.UserData.PayAmountConfig = list[2]
		end
	else
		print("+++++++++++++++++++ SetPayAmountConfig list is nil")
	end
end

function GetToken()
    local token = [REDACTED]"Token")	
	local tokenData = {}
	if token ~= "" and token ~= nil then
		local jasnData = G.Tool.Aes.Decode(token)
		tokenObject = jasnData["token"]
		tokenData.value = tokenObject["value"]
		tokenData.time = tokenObject["time"]
		--print(" Token value " .. tokenData.value)
		--print(" Token time " .. tokenData.time)
	end
	return tokenData
end

function GetFullToken()
    local token = [REDACTED]"Token")
	return token
end


function CheckOpenData(isAnd, ShowSumCharge, ShowGamePlayTime, ShowBindPhone)
	local user = GetCurrentUser();
	local userdata = user.UserData
	local GamePlayTime = G.GetGamePlayTime()
	local ShowTime = ShowGamePlayTime * 60
	
	-- AND
	if isAnd == true then
		-- 充值金額达到指定
		if userdata.SumCharge < ShowSumCharge then
			return false
		end
		
		-- 時間達到指定
		if GamePlayTime < ShowTime then
			return false
		end
		
		-- 绑定
		if ShowBindPhone == 1 then
			if userdata.IsBindGuest == false then
				return false
			end
		end		
		
		return true
	
	-- OR
	else
		-- 充值金額达到指定
		if userdata.SumCharge >= ShowSumCharge then
			return true
		end		
		
		-- 時間達到指定
		if GamePlayTime >= ShowTime then
			return true
		end

		-- 绑定
		if ShowBindPhone == 1 then
			if userdata.IsBindGuest == true then
				return true
			end
		end
		
		return false
	end
end

