require("Module")
mmodule(...)

local G = require("Global")

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

-- 0 方舟
-- 1 986
-- 2 藍鑽
local branch = 0
-- Check branch data
function GetBranch()
    return branch
end

-- 110 77棋牌
-- 111 迪祥斗地主
-- 112 奥云黑龙江麻将
-- 113 天天卡五星
-- 114 201牛牛
-- 115 湘约麻将
-- 116 翱翔红中麻将
-- 122 土豪赢三张
-- 124 最爱斗牛

-- 要過濾的遊戲名稱
-- true 走流程1
-- false 走流程2
local SubmitForReview = {};
SubmitForReview["111"] = false
SubmitForReview["112"] = false
--SubmitForReview["115"] = true
SubmitForReview["116"] = true

-- 是否在送審階段
function IsSubmitForReview()
    local PlatformID = G.GameConfig.GamePlatformID
    if SubmitForReview[tostring(PlatformID)] == nil then
        return false
    else
        return SubmitForReview[tostring(PlatformID)]
    end
end

-- 版權
local copyright = {}
copyright["111"] = "上海迪祥信息科技有限公司"
copyright["113"] = "襄阳狂野怪兽信息科技有限公司"
copyright["115"] = "关注财产安全 拒绝赌博"
copyright["116"] = "广东翱翔网络科技有限公司"

-- 取得版權
function GetCopyright()
    local PlatformID = G.GameConfig.GamePlatformID
    if copyright[tostring(PlatformID)] == nil then
        return ""
    else
        return copyright[tostring(PlatformID)]
    end
end

-- 把“元”字改成 “币”
local changeCoin = {};
changeCoin["111"] = true
changeCoin["114"] = true
changeCoin["115"] = true
changeCoin["116"] = true
changeCoin["122"] = true
changeCoin["124"] = true
changeCoin["127"] = true
changeCoin["135"] = true
changeCoin["136"] = true
changeCoin["138"] = true
changeCoin["140"] = true
changeCoin["148"] = true
changeCoin["149"] = true
changeCoin["161"] = true

-- 是否把“元”字改成 “币”
function IsChangeCoin()
    local PlatformID = G.GameConfig.GamePlatformID
    if changeCoin[tostring(PlatformID)] == nil then
        return false
    else
        return changeCoin[tostring(PlatformID)]
    end
end

-- 唯一顯示的遊戲
local onlyShowGameId = {};
onlyShowGameId["111"] = 1
onlyShowGameId["112"] = 15
--onlyShowGameId["115"] = 15
onlyShowGameId["116"] = 15

-- 取得唯一顯示的遊戲
function GetOnlyShowGameId()
    local PlatformID = G.GameConfig.GamePlatformID
    if onlyShowGameId[tostring(PlatformID)] == nil then
        return 0
    else
        return onlyShowGameId[tostring(PlatformID)]
    end
end

-- 载入背景
local changeLoadingBG = {};
changeLoadingBG["122"] = true
changeLoadingBG["124"] = true
changeLoadingBG["135"] = true
changeLoadingBG["136"] = true
changeLoadingBG["143"] = true
changeLoadingBG["144"] = true
changeLoadingBG["145"] = true
changeLoadingBG["172"] = true
changeLoadingBG["174"] = true
changeLoadingBG["178"] = true

-- 是否修改载入背景
function isChangeLoadingBG()
    local PlatformID = G.GameConfig.GamePlatformID
    if changeLoadingBG[tostring(PlatformID)] == nil then
        return false
    else
        return changeLoadingBG[tostring(PlatformID)]
    end
end

-- 载入LOGO
local changeLoadingLogo = {};
changeLoadingLogo["122"] = true
changeLoadingLogo["124"] = true
changeLoadingLogo["135"] = true
changeLoadingLogo["136"] = true
changeLoadingLogo["143"] = true
changeLoadingLogo["144"] = true
changeLoadingLogo["145"] = true
changeLoadingLogo["172"] = true
changeLoadingLogo["174"] = true
changeLoadingLogo["178"] = true

-- 是否修改载入LOGO
function isChangeLoadingLogo()
    local PlatformID = G.GameConfig.GamePlatformID
    if changeLoadingLogo[tostring(PlatformID)] == nil then
        return false
    else
        return changeLoadingLogo[tostring(PlatformID)]
    end
end

function IsUpdateAllGame()
    local isUpdateAllGame = CS.UnityEngine.PlayerPrefs.GetInt("IsShowAllGame", 0)
    return isUpdateAllGame
end

function ShowUpdateMessage(OkCallback)
    local PlatformID = G.GameConfig.GamePlatformID
    if SubmitForReview[tostring(PlatformID)] == nil then
        return
    end

    if SubmitForReview[tostring(PlatformID)] == true then
        return
    end

    local isUpdateAllGame = IsUpdateAllGame()
    if isUpdateAllGame ~= 0 then
        return
    end

    G.Util.coroutine_call(function()
        local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "是否更新游戏?", "是", "否");
        if boxResult == CS.MessageBoxResult.OK then
            CS.UnityEngine.PlayerPrefs.SetInt("IsShowAllGame", 1)
            CS.UnityEngine.PlayerPrefs.Save()
            OkCallback()
        end
    end)()
end

-- 預設關閉客服清單
local defaultclose = {}
defaultclose[1] = 122
defaultclose[2] = 124
defaultclose[3] = 127
defaultclose[4] = 135
defaultclose[5] = 136
defaultclose[6] = 114
defaultclose[7] = 138
defaultclose[8] = 140
defaultclose[9] = 148
defaultclose[10] = 149
defaultclose[11] = 161
function IsShowStartUpCS(pid)
    local closecs = false
    -- 是否為預設關閉清單


    for k, v in pairs(defaultclose) do
        if (pid == v) then
            closecs = true;
        end
    end
    --如果符合預設關閉清單在判斷是否符合二階段
    if (closecs) then
        --二階段驗證,確認使用者是否已經符合儲值滿300以及遊玩時間滿足條件
        local closecs2 = CS.UnityEngine.PlayerPrefs.GetInt("IsStartUpCS")
        closecs = (closecs2 == 0) and true or false;

        return (closecs ~= true)
    else
        return true;
    end
end

-- 客服按钮开关运算方式
local CustomServiceOperationAnd = {};
CustomServiceOperationAnd["122"] = false
CustomServiceOperationAnd["124"] = false
CustomServiceOperationAnd["127"] = false
CustomServiceOperationAnd["135"] = false
CustomServiceOperationAnd["136"] = false
CustomServiceOperationAnd["140"] = false
CustomServiceOperationAnd["148"] = false
CustomServiceOperationAnd["149"] = false
CustomServiceOperationAnd["161"] = false

-- 取得客服按钮开关运算方式
function getCustomServiceOperationAnd()
    if G.User.andor.CSAndOr ~= nil then
        return G.User.andor.CSAndOr
    else
        print("### G.User.andor.CSAndOr = nil")
    end

    local PlatformID = G.GameConfig.GamePlatformID
    if CustomServiceOperationAnd[tostring(PlatformID)] == nil then
        return true
    else
        return CustomServiceOperationAnd[tostring(PlatformID)]
    end
end

-- 兑换按钮开关运算方式
local ExchangeOperationAnd = {};
ExchangeOperationAnd["122"] = false
ExchangeOperationAnd["124"] = false
ExchangeOperationAnd["127"] = false
ExchangeOperationAnd["135"] = false
ExchangeOperationAnd["136"] = false
ExchangeOperationAnd["114"] = false
ExchangeOperationAnd["138"] = false
ExchangeOperationAnd["140"] = false
ExchangeOperationAnd["148"] = false
ExchangeOperationAnd["149"] = false
ExchangeOperationAnd["161"] = false

-- 取得兑换按钮开关运算方式
function getExchangeOperationAnd()
    if G.User.andor.ConvertAndOr ~= nil then
        return G.User.andor.ConvertAndOr
    else
        print("### G.User.andor.ConvertAndOr = nil")
    end

    local PlatformID = G.GameConfig.GamePlatformID
    if ExchangeOperationAnd[tostring(PlatformID)] == nil then
        return true
    else
        return ExchangeOperationAnd[tostring(PlatformID)]
    end
end

-- VIP充值和闪付按钮开关运算方式
local RechargeOperationAnd = {};
RechargeOperationAnd["122"] = false
RechargeOperationAnd["124"] = false
RechargeOperationAnd["127"] = false
RechargeOperationAnd["135"] = false
RechargeOperationAnd["136"] = false
RechargeOperationAnd["148"] = false
RechargeOperationAnd["149"] = false

-- 取得VIP充值和闪付按钮开关运算方式
function getRechargeOperationAnd()
    if G.User.andor.VIPAgentAndOr ~= nil then
        return G.User.andor.VIPAgentAndOr
    else
        print("### G.User.andor.VIPAgentAndOr = nil")
    end

    local PlatformID = G.GameConfig.GamePlatformID
    if RechargeOperationAnd[tostring(PlatformID)] == nil then
        return true
    else
        return RechargeOperationAnd[tostring(PlatformID)]
    end
end

-- 个人讯息绑定手机显示支付宝和银行卡
local RoleBindZfiAndUp = {};
RoleBindZfiAndUp["122"] = true
RoleBindZfiAndUp["124"] = true
RoleBindZfiAndUp["135"] = true
RoleBindZfiAndUp["136"] = true

-- 取得个人讯息绑定手机显示支付宝和银行卡
function getRoleBindZfiAndUp()
    local PlatformID = G.GameConfig.GamePlatformID
    if RoleBindZfiAndUp[tostring(PlatformID)] == nil then
        return false
    else
        return RoleBindZfiAndUp[tostring(PlatformID)]
    end
end

-- 金牌娛樂特製
local JPYL = {};
JPYL["143"] = true
JPYL["144"] = true
JPYL["145"] = true
JPYL["172"] = true
JPYL["174"] = true
JPYL["178"] = true

-- 取得是否為金牌娛樂特製
function getJPYL()
    local PlatformID = G.GameConfig.GamePlatformID
    if JPYL[tostring(PlatformID)] == nil then
        return false
    else
        return JPYL[tostring(PlatformID)]
    end
end

-- 沒有加贈2%
local NoTwoPa = {};
NoTwoPa["119"] = false
NoTwoPa["128"] = false
NoTwoPa["129"] = false
NoTwoPa["130"] = false
NoTwoPa["142"] = false
NoTwoPa["143"] = false
NoTwoPa["144"] = false
NoTwoPa["145"] = false
NoTwoPa["146"] = false
NoTwoPa["147"] = false
-- NoTwoPa["986"] = false
-- NoTwoPa["987"] = false
-- NoTwoPa["988"] = false
-- NoTwoPa["989"] = false
-- NoTwoPa["990"] = false
-- NoTwoPa["991"] = false
-- NoTwoPa["992"] = false
-- NoTwoPa["993"] = false
-- NoTwoPa["994"] = false
-- NoTwoPa["995"] = false
NoTwoPa["110"] = false
NoTwoPa["126"] = false
NoTwoPa["172"] = false
NoTwoPa["174"] = false
NoTwoPa["178"] = false

-- 取得是否加贈2%
function getNoTwoPa()
    local PlatformID = G.GameConfig.GamePlatformID
    if NoTwoPa[tostring(PlatformID)] == nil then
        return true
    else
        return NoTwoPa[tostring(PlatformID)]
    end
end

--是否推薦 (OSSM = おすすめ = 推薦) by AZLDC
local OSSM = {};
OSSM["174"] = true
function getOSuSuMe()
    local PlatformID = G.GameConfig.GamePlatformID
    if OSSM[tostring(PlatformID)] == nil then
        return false
    else
        return OSSM[tostring(PlatformID)]
    end
end

--是否替換Henry成VIP顯示 by AZLDC
local VIPType = {};
VIPType["174"] = true
function getHenryToVIP()
    local PlatformID = G.GameConfig.GamePlatformID
    if VIPType[tostring(PlatformID)] == nil then
        return false
    else
        return VIPType[tostring(PlatformID)]
    end
end


-- 大廳上方logo (1.金牌, 2.轟達)
local HallLogo = {}
HallLogo["119"] = 2
HallLogo["128"] = 2
HallLogo["129"] = 2
HallLogo["130"] = 2
HallLogo["170"] = 170
HallLogo["172"] = 172
HallLogo["174"] = 174
HallLogo["178"] = 178

function getHallLogo()
    local PlatformID = G.GameConfig.GamePlatformID
    if HallLogo[tostring(PlatformID)] == nil then
        return 0
    else
        return HallLogo[tostring(PlatformID)]
    end
end

-- 綁定按鈕
local BindingGift = {}
BindingGift["108"] = false
BindingGift["141"] = false
BindingGift["151"] = false

function getBindingGift()
    local PlatformID = G.GameConfig.GamePlatformID
    if BindingGift[tostring(PlatformID)] == nil then
        return true
    else
        return BindingGift[tostring(PlatformID)]
    end
end

-- 改變兌換按鈕順序
local ChangeExchangeOrder = {}
ChangeExchangeOrder["174"] = true

-- 是否改變兌換按鈕順序
function IsChangeExchangeOrder()
    -- 986全部移到最上方
    if G.Platform.GetBranch() == 1 then
        return true
    else
        local PlatformID = G.GameConfig.GamePlatformID
        if ChangeExchangeOrder[tostring(PlatformID)] == nil then
            return false
        else
            return ChangeExchangeOrder[tostring(PlatformID)]
        end
    end
end

-- 顯示兌換銀聯按鈕標籤
local ExchangeMaskUP = {}
ExchangeMaskUP["174"] = true

-- 是否顯示兌換銀聯按鈕標籤
function IsSExchangeMaskUP()
    local PlatformID = G.GameConfig.GamePlatformID
    if ExchangeMaskUP[tostring(PlatformID)] == nil then
        return false
    else
        return ExchangeMaskUP[tostring(PlatformID)]
    end
end


-- 闪付广告
local HenryAD = {}
HenryAD["174"] = true

-- 取得闪付广告
function GetHenryAD()
    local PlatformID = G.GameConfig.GamePlatformID
    if HenryAD[tostring(PlatformID)] == nil then
        return false
    else
        return HenryAD[tostring(PlatformID)]
    end
end