require("Module")
mmodule(...)

local G = require("Global")

-- region 游戏ID配置
GameIDs = {};
GameIDs["1"] = "ddz" -- 斗地主
GameIDs["2"] = "qznn" -- 抢庄牛牛
GameIDs["3"] = "zjh" -- 诈金花
GameIDs["4"] = "dzpk" -- 德州扑克
GameIDs["5"] = "wzp" -- 梭哈
GameIDs["6"] = "qzesyd" -- 21点
GameIDs["7"] = "tbnn" -- 经典牛牛
GameIDs["8"] = "sangong" -- 三公
GameIDs["9"] = "sgj" -- 水果机
GameIDs["10"] = "sss" -- 十三水
GameIDs["11"] = "lhd" -- 龙虎斗
GameIDs["12"] = "ttz" -- 推筒子
GameIDs["13"] = "bjl" -- 百家乐
GameIDs["14"] = "by" -- 捕鱼
GameIDs["15"] = "scmj" -- 四川麻将
GameIDs["16"] = "srnn" -- 双人牛牛
GameIDs["17"] = "ermj" -- 兩人麻将
GameIDs["18"] = "fkeb" -- 疯狂二八
GameIDs["19"] = "brnn" -- 百人牛牛
GameIDs["20"] = "hhdz" -- 紅黑大戰
GameIDs["21"] = "hhnn" -- 紅黑牛牛
GameIDs["22"] = "codyqzesyd" -- 21
GameIDs["23"] = "lhj" -- 老虎机
GameIDs["24"] = "bynew" -- 新捕鱼
GameIDs["25"] = "dfdc" -- 多福多财
GameIDs["26"] = "tgw" -- 糖果屋
GameIDs["27"] = "ald" -- 阿拉丁
-- endregion


-- 遊戲名稱
GameNames = {};
GameNames["1"] = { [0] = "斗地主", [1] = "迪祥斗地主", [2] = "天津随心玩", [3] = "小拾拼十", [4] = "10086游戏中心斗地主" }
GameNames["2"] = { [0] = "抢庄牛牛", [1] = "欢乐抢庄", [2] = "四人抢庄", [3] = "快来抢庄", [4] = "PK游戏", [5] = "抢牛大战", [6] = "201牛牛", [7] = "拼五张", [8] = "最爱斗牛" }
GameNames["3"] = { [0] = "炸金花", [1] = "火拼三牌", [2] = "飘三叶", [3] = "赢三张", [4] = "三张牌", [5] = "摸三张", [6] = "勾紧花", [7] = "闷金花", [8] = "土豪赢三张" }
GameNames["4"] = { [0] = "德州扑克" }
GameNames["5"] = { [0] = "梭哈" }
GameNames["6"] = { [0] = "21点", [1] = "Black Jack" }
GameNames["7"] = { [0] = "经典牛牛", [1] = "经典斗牌", [2] = "通比牛牛", [3] = "通比拼十" }
GameNames["8"] = { [0] = "三公" }
GameNames["9"] = { [0] = "股市大亨", [1] = "大平台", [2] = "BTA大战", [3] = "股市老虎机" }
GameNames["10"] = { [0] = "十三水" }
GameNames["11"] = { [0] = "龙虎斗", [1] = "龙争虎斗", [2] = "龙虎争霸", [3] = "龙与虎", [4] = "争强好胜" }
GameNames["12"] = { [0] = "推筒子" }
GameNames["13"] = { [0] = "百家乐", [1] = "百家欢乐牌", [2] = "千人欢乐牌", [3] = "欢乐牌", [4] = "庄闲大战", [5] = "百人大战" }
GameNames["14"] = { [0] = "捕鱼" }
GameNames["15"] = { [0] = "四川麻将", [1] = "奥云黑龙江麻将", [2] = "湘约麻将", [3] = "翱翔红中麻将" }
GameNames["16"] = { [0] = "双人牛牛", [1] = "双人拼十" }
GameNames["17"] = { [0] = "双人麻将" }
GameNames["18"] = { [0] = "二八杠" }
GameNames["19"] = { [0] = "百人牛牛", [1] = "百战牛魔王", [2] = "百人大战" }
GameNames["20"] = { [0] = "红黑大战" }
GameNames["21"] = { [0] = "红黑牛牛", [1] = "红黑牛魔王", [2] = "红黑斗" }
GameNames["22"] = { [0] = "21点" }
GameNames["23"] = { [0] = "多福多财" }
GameNames["24"] = { [0] = "捕鱼大玩咖" }
GameNames["25"] = { [0] = "多福多财", [1] = "皇家老虎机" }
GameNames["26"] = { [0] = "糖果屋" }
GameNames["27"] = { [0] = "阿拉丁" }

-- 要過濾的遊戲名稱
local FilterGameNames = {};
FilterGameNames["1"] = "斗地主" -- 3
FilterGameNames["2"] = "抢庄牛牛" -- 6
FilterGameNames["3"] = "炸金花" -- 2
FilterGameNames["4"] = "德州扑克" -- 10
FilterGameNames["5"] = "梭哈" -- 15
FilterGameNames["6"] = "21点"
FilterGameNames["7"] = "经典牛牛" -- 7
FilterGameNames["8"] = "三公" -- 9
FilterGameNames["9"] = "股市大亨" -- 22
FilterGameNames["10"] = "十三水" -- 13
FilterGameNames["11"] = "龙虎斗" -- 4
FilterGameNames["12"] = "推筒子" -- 12
FilterGameNames["13"] = "百家乐" -- 5
FilterGameNames["14"] = "捕鱼" -- 14
FilterGameNames["15"] = "四川麻将" -- 11
FilterGameNames["16"] = "双人牛牛" -- 16
FilterGameNames["17"] = "兩人麻将" -- 17
FilterGameNames["18"] = "二八杠" -- 23
FilterGameNames["19"] = "百人牛牛" -- 18
FilterGameNames["20"] = "红黑大战" -- 24
FilterGameNames["21"] = "红黑牛牛" -- 25
FilterGameNames["22"] = "21点" -- 8
FilterGameNames["23"] = "老虎机"
FilterGameNames["24"] = "捕鱼大玩咖" -- 26
FilterGameNames["25"] = "多福多财" -- 27
FilterGameNames["26"] = "糖果屋" -- 28
FilterGameNames["27"] = "阿拉丁" -- 29

function GetGameIDCount()
    local count = 0
    for k, v in pairs(GameIDs) do
        count = count + 1
    end
    return count
end

-- region 配置需要显示桌子号的游戏ID
local showTableNoGame = {}
showTableNoGame["9"] = true
showTableNoGame["11"] = true

-- 判断游戏是否需要显示桌子号
function GetIsShowTableNo(_gameid)
    return showTableNoGame[tostring(_gameid)] == true
end

-- 是否為百人場
function IsMultiGame(GameId)
    if GameId == 13 then
        return true
    elseif GameId == 11 then
        return true
    elseif GameId == 19 then
        return true
    elseif GameId == 20 then
        return true
    elseif GameId == 21 then
        return true
    elseif GameId == 9 then
        return true
    else
        return false
    end
end

-- 获取游戏IconName
function GetGameIconName(gameId)
    local userGameList = G.User.GetUserGameList2()
    for k, v in pairs(userGameList) do
        if tostring(v.GameID) == tostring(gameId) then
            return tonumber(v.IconName)
        end
    end
    print("GetGameIconName not find " .. gameId)
    return 0
end

-- 获取游戏项目名称
function GetGameName(GameId)
    local iconName = GetGameIconName(GameId)
    local gameName = GetGameNameByIconName(tostring(GameId), iconName)
    return gameName
end


-- 从gameID和iconName获取游戏项目名称
function GetGameNameByIconName(GameId, IconName)
    if GameId == nil or GameId == "nil" then
        print("gameID is nil")
        return ""
    elseif GameNames[tostring(GameId)] == nil then
        print("GameNames not find, gameID:" .. tostring(GameId))
        return ""
    else
        return GameNames[tostring(GameId)][IconName]
    end
end

-- 取得要過濾的名稱
function GetFilterGameNamess()
    return FilterGameNames
end

-- 替換訊息中遊戲名稱
function ReplaceGameName(message)
    for k, v in pairs(FilterGameNames) do
        -- 找關鍵字
        q = string.match(message, v)
        if q ~= nil then
            local iconName = GetGameIconName(k)
            if iconName > 0 then
                local name = GetGameNameByIconName(k, iconName)
                local newMessage = string.gsub(message, v, name)
                return newMessage
            end
        end
    end

    return message
end

-- 是否顯示所有遊戲
function IsShowAllGame()
    local user = G.User.GetCurrentUser()
    local userdata = user.UserData

    local isShowAll = false
    -- 累计充值
    if userdata.SumCharge >= userdata.ShowAllGameLimit then
        isShowAll = true
    end

    -- 绑定手机
    if isShowAll == true then
        if userdata.ShowAllGameBindPhone == 1 then
            if userdata.IsBindGuest == true then
                isShowAll = true
            end
        end
    end

    if G.Platform.GetOnlyShowGameId() == 0 then
        return isShowAll
    end

    if isShowAll == true then
        -- 送審階段只顯示單一遊戲
        if G.Platform.IsSubmitForReview() == true then
            --print("送審階段 " .. tostring(G.IsFirstTimePlay()))
            if G.IsFirstTimePlay() == true then
                --print("送審階段只顯示單一遊戲")
                isShowAll = false
            else
                --print("送審階段顯示所有遊戲")
                isShowAll = true
            end
        else
            local isUpdateAllGame = G.Platform.IsUpdateAllGame()
            --print("非送審階段 " .. tostring(isUpdateAllGame))
            if isUpdateAllGame == 0 then
                --print("非送審階段只顯示單一遊戲")
                isShowAll = false
            else
                --print("非送審階段顯示所有遊戲")
                isShowAll = true
            end
        end
    end

    return isShowAll
end