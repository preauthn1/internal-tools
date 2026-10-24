require("Module")
mmodule(...)

local roomInfo = nil
local roomLevelNames = {
    "初级场",
    "菁英场",
    "豪华场",
    "富豪场",
    "皇家场",
    "体验场"
}

local ddzLevelNames = {
    "高级场",
    "菁英场",
    "富豪场",
    "",
    "",
    "",
    "高级场",
    "菁英场",
    "富豪场",
    "",
    "",
    ""	
}

local multiRoomName = "万人场"

-- 游戏房间提示步骤定义
RoomTips = {
    Allot = "正在分配座位……",
    WaitOtherPlayer = "正在等待其他玩家进入……",
    ReAllot = "正在重新分配座位……",
    WaitGameStart = "等待游戏开始……",
    Restore = "正在恢复游戏场景……"
}
function Clear()
    roomInfo = nil
end
 
-- 设置代理信息
function SetRoomInfo(currentRoomInfo)
    roomInfo = currentRoomInfo
end

-- 获取代理信息
function GetRoomInfo()
    return roomInfo
end 


local userData = nil
-- 设置玩家登入游戏状态
function SetUserData(currentUserData)
    userData = currentUserData
end

-- 获取玩家登入游戏状态
function GetUserData()
    return userData
end 


-- 1 试炼场 2 精英场 豪华场 4富豪场 5皇家场
function GetRoomLevel()
    return roomInfo.Flag
end

-- 获取游戏等级场名称
function GetRoomLevelName()
    return roomLevelNames[tonumber(roomInfo.Flag)]
end

-- 获取游戏萬人場名稱
function GetMultiRoomName()
    return multiRoomName
end

-- 获取斗地主場名稱
function GetDDZRoomName()
    return ddzLevelNames[tonumber(roomInfo.Flag)]
end

local gameIsLoginSuccess = false
-- 获取/设置游戏是否登录成功
function GameLoginSuccess(val)
    if (val ~= nil) then
        gameIsLoginSuccess = val
    else
        return gameIsLoginSuccess
    end
end