local RecTran_Main = this.transform:Find("_Main"):GetComponent("RectTransform")
local RecTran_GameList = this.transform:Find("_Main/Panel/_GameList"):GetComponent("RectTransform")
local Scrr_GameScrollView = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView"):GetComponent("ScrollRect")
local RecTran_GameScrollView = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView"):GetComponent("RectTransform")
local RecTran_CarouselRoot = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView/GameListContent/_Carousel"):GetComponent("RectTransform")
local AdScrr_GameScrollView = RecTran_CarouselRoot.transform:Find("_offset/_AddScrollView"):GetComponent("ScrollRect")
local AdScrr_GameListContent = RecTran_CarouselRoot.transform:Find("_offset/_AddScrollView/GameListContent"):GetComponent("RectTransform")
local RecTran_Carousel = RecTran_CarouselRoot.transform:Find("_offset"):GetComponent("RectTransform")

local Btn_AdScrr = RecTran_CarouselRoot.transform:Find("_offset"):GetComponent("Button")
local RecTran_Panel = this.transform:Find("_Main/Panel"):GetComponent("RectTransform")
local Layout_Panel = this.transform:Find("_Main/Panel"):GetComponent("HorizontalLayoutGroup")
local RecTran_GameAttach = this.transform:Find("_Main/Panel/GameAttach"):GetComponent("RectTransform")

local RecTran_GameRootTemp = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView/GameListContent/_GameRootTemp"):GetComponent("RectTransform")
local GridLayoutGroup_GameRootTemp = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView/GameListContent/_GameRootTemp"):GetComponent("GridLayoutGroup")
local RecTran_GameAttachList = this.transform:Find("_Main/Panel/GameAttach/_GameScrollView/GameListContent"):GetComponent("RectTransform")
local Btn_LeftPage = this.transform:Find("_Main/Panel/GameAttach/_LeftPage"):GetComponent("Button")
local Btn_RightPage = this.transform:Find("_Main/Panel/GameAttach/_RightPage"):GetComponent("Button")
local Txt_LoadingGameTips = this.transform:Find("_Main/_LoadingGameTips"):GetComponent("Text")

local Img_BottomAttach = this.transform:Find("_Main/Bottom/ImageBg"):GetComponent("Image")
local Root_BottomAttach = this.transform:Find("_Main/Bottom/Attach"):GetComponent("RectTransform")
local Btn_agent = this.transform:Find("_Main/Bottom/Attach/_Agent"):GetComponent("Button")
--local Btn_Proxy = this.transform:Find("_Main/Bottom/Attach/_Proxy"):GetComponent("Button")
local Btn_Mail = this.transform:Find("_Main/Bottom/Attach/_Mail"):GetComponent("Button")
local Img_Mail1 = this.transform:Find("_Main/Bottom/Attach/_Mail/ImageState1"):GetComponent("Image")
local Img_Mail2 = this.transform:Find("_Main/Bottom/Attach/_Mail/ImageState2"):GetComponent("Image")
local Txt_Mail = this.transform:Find("_Main/Bottom/Attach/_Mail/ImageState2/Text"):GetComponent("Text")
local Img_CS1 = this.transform:Find("_Main/TOP/_CustomService/ImageState1"):GetComponent("Image")
local Img_CS2 = this.transform:Find("_Main/TOP/_CustomService/ImageState2"):GetComponent("Image")
local Txt_CS = this.transform:Find("_Main/TOP/_CustomService/ImageState2/Text"):GetComponent("Text")
local Btn_Coffer = this.transform:Find("_Main/Bottom/Attach/_Coffer"):GetComponent("Button")
local Img_Coffer1 = this.transform:Find("_Main/Bottom/Attach/_Coffer/ImageState1"):GetComponent("Image")
local Img_Coffer2 = this.transform:Find("_Main/Bottom/Attach/_Coffer/ImageState2"):GetComponent("Image")
local Btn_Exchange = this.transform:Find("_Main/Bottom/Attach/_Exchange"):GetComponent("Button")
--local Img_Exchange = this.transform:Find("_Main/Bottom/Attach/_Exchange/ImageState2"):GetComponent("Image")
local Btn_ExchangeFont1 = this.transform:Find("_Main/Bottom/Attach/_Exchange/ImageState1"):GetComponent("Image")
local Btn_ExchangeFont2 = this.transform:Find("_Main/Bottom/Attach/_Exchange/ImageState2"):GetComponent("Image")
local Btn_ExchangeFont3 = this.transform:Find("_Main/Bottom/Attach/_Exchange/ImageState3"):GetComponent("Image")
local Btn_ExchangeFont4 = this.transform:Find("_Main/Bottom/Attach/_Exchange/ImageState4"):GetComponent("Image")
local img_Bonus = this.transform:Find("_Main/Bottom/Attach/_Exchange/Bonus"):GetComponent("Image")
local txt__NowText = this.transform:Find("_Main/Bottom/Attach/_Exchange/Bonus/_NowText"):GetComponent("Text")
local txt__LimitText = this.transform:Find("_Main/Bottom/Attach/_Exchange/Bonus/_LimitText"):GetComponent("Text")


local Btn_Recharge = this.transform:Find("_Main/Bottom/_Recharge"):GetComponent("Button")
local Img_RechargeFont1 = this.transform:Find("_Main/Bottom/_Recharge/ImageState1"):GetComponent("Image")
local Img_RechargeFont2 = this.transform:Find("_Main/Bottom/_Recharge/ImageState2"):GetComponent("Image")
local Img_RechargeFont3 = this.transform:Find("_Main/Bottom/_Recharge/ImageState3"):GetComponent("Image")
local Img_RechargeFont4 = this.transform:Find("_Main/Bottom/_Recharge/ImageState4"):GetComponent("Image")
--[[
local Btn_RoomRecharge = this.transform:Find("_Room/_Recharge"):GetComponent("Button")'
]]
local Btn_CustomService = this.transform:Find("_Main/TOP/_CustomService"):GetComponent("Button")
local Img_CustomService = this.transform:Find("_Main/TOP/_CustomService/ImageState2"):GetComponent("Image")
local Btn_Notice = this.transform:Find("_Main/Bottom/Attach/_Notice"):GetComponent("Button")
local Img_Notice1 = this.transform:Find("_Main/Bottom/Attach/_Notice/ImageState1"):GetComponent("Image")
local Img_Notice2 = this.transform:Find("_Main/Bottom/Attach/_Notice/ImageState2"):GetComponent("Image")
local Img_Logo = this.transform:Find("_Main/TOP/_Logo"):GetComponent("Image")
local Img_LogoSprite = {}
local Btn_Promo = this.transform:Find("_Main/TOP/_Promo"):GetComponent("Button")
local Btn_AgentTop = this.transform:Find("_Main/TOP/_AgentTop"):GetComponent("Button")
--[[
local ani_shark = this.transform:Find("_Room/Background1/_shark"):GetComponent("SkeletonGraphic");
local ani_kani = this.transform:Find("_Room/Background1/_kani"):GetComponent("SkeletonGraphic");
]]
local Btn_Setting = this.transform:Find("_Main/TOP/_Setting"):GetComponent("Button")
local Btn_Upgrade = this.transform:Find("_Main/Bottom/Attach/_Upgrade"):GetComponent("Button")
--[[
local Btn_HeadIcon = this.transform:Find("_Main/TOP/_HeadIcon"):GetComponent("Button")
local Txt_Money = this.transform:Find("_Main/TOP/_Money"):GetComponent("Text")
local Btn_GetMoney = this.transform:Find("_Main/TOP/_GetMoney"):GetComponent("Button")
local RecTran_RoomMoney = this.transform:Find("_Room/_Money"):GetComponent("RectTransform")
local Txt_RoomMoney = this.transform:Find("_Room/_Money/Text"):GetComponent("Text")

local Txt_PlayerID = this.transform:Find("_Main/TOP/_PlayerID"):GetComponent("Text")
]]
local RecTran_Ann = this.transform:Find("_Main/TOP/_Ann"):GetComponent("RectTransform")
local Txt_NoticeContent = this.transform:Find("_Main/TOP/_Ann/Ann/_NoticeContent"):GetComponent("Text")
local RecTran_Room = this.transform:Find("_Room"):GetComponent("RectTransform")
local Btn_RoomlistBack = this.transform:Find("_Room/_RoomlistBack"):GetComponent("Button")
local Img_RoomTitle = this.transform:Find("_Room/title/_RoomTitle"):GetComponent("Image")
local RecTran_Rooms = this.transform:Find("_Room/_Rooms"):GetComponent("RectTransform")
local RecTran_RoomsBy = this.transform:Find("_Room/_RoomsBy"):GetComponent("RectTransform")
local RecTran_RoomsDdz = this.transform:Find("_Room/_RoomsDdz"):GetComponent("RectTransform")
local RecTran_RoomList = this.transform:Find("_Room/_RoomList"):GetComponent("RectTransform")
local RecTran_RoomRow1 = this.transform:Find("_Room/_RoomList/_RoomRow1"):GetComponent("RectTransform")
local RecTran_RoomRow2 = this.transform:Find("_Room/_RoomList/_RoomRow2"):GetComponent("RectTransform")
local RecTran_RoomRow2 = this.transform:Find("_Room/_RoomList/_RoomRow2"):GetComponent("RectTransform")
local Txt_RoomVersion = this.transform:Find("_Room/_Version"):GetComponent("Text")
local Txt_MaxOdds = this.transform:Find("_Room/_MaxOdds"):GetComponent("Text")
local RecTran_NoDataGame = this.transform:Find("_Main/Panel/_GameList/_NoDataGame"):GetComponent("RectTransform")
--local RecTran_RoomPerson = this.transform:Find("_Room/Background/Image"):GetComponent("RectTransform")
local RecTran_RoomTitle = this.transform:Find("_Room/title"):GetComponent("RectTransform")
local recTran__flagRoot = this.transform:Find("_Main/Panel/_flag"):GetComponent("RectTransform")

local ImgShader = this.transform:Find("_Main/Panel/_shader"):GetComponent("Image")
--[[
local RecTran_Background14 = this.transform:Find("_Room/Background14"):GetComponent("RectTransform")
local RecTran_Background1 = this.transform:Find("_Room/Background1"):GetComponent("RectTransform")
]]
local RecTran_EF = this.transform:Find("EF"):GetComponent("RectTransform")

local btn__MoneyPool = this.transform:Find("_Main/TOP/_MoneyPool"):GetComponent("Button")
local txt__Money = this.transform:Find("_Main/TOP/_MoneyPool/_Money"):GetComponent("Text")

local RecTran_CellPhone = this.transform:Find("CellPhone"):GetComponent("RectTransform")
local CellPhoneAni = this.transform:Find("CellPhone/CellphoneAni")

local btn__VipIcon = this.transform:Find("_Common/_HeadIcon/_VipIcon"):GetComponent("Button")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local RoomList = {};
local GameLogoList = {};
local GameVersion = {};
local GameT = {};
local gamelist = {};
local autolist = {};
local ranklist = {};

local iconNameList = {};

local moneypool = 0
local seqJackpotPoolRet = nil

local G = require "Global";
local user = G.User.GetCurrentUser();
local userdata = user.UserData
local CheckIsGetData = G.RequireG("Tools/CheckIsGetData")
local uiRank = G.RequireG("UIHallRank")
uiRank.Init(this.transform:Find("_Rank"):GetComponent("RectTransform"))
local uiCommon = G.RequireG("UIHallCommon")
uiCommon.Init(this.transform:Find("_Common"):GetComponent("RectTransform"))

local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

local groupName = "UIHall";

-- 当前请求进入的游戏ID
local RequestGameID = 0;
local isFisrtRequestServerList = true
local mailList = {};
local CurrentUrl = "";
local adtype = 0;
local activeGameNumber = 0;
local GameListWidth = 900;
local nowGameId = 0
local IsGetServerList = false

-- 游戏翻页时触发，用于显示或隐藏左右翻页按钮
local function onGamePageChanged(pageCount, currentIndex)
    --[[
    if pageCount == 1 then
        Btn_LeftPage.gameObject:SetActive(false);
        Btn_RightPage.gameObject:SetActive(false);
    elseif currentIndex == 0 then
        Btn_LeftPage.gameObject:SetActive(false);
        Btn_RightPage.gameObject:SetActive(true);
    elseif pageCount == currentIndex + 1 then
        Btn_LeftPage.gameObject:SetActive(true);
        Btn_RightPage.gameObject:SetActive(false);
    else
        Btn_LeftPage.gameObject:SetActive(true);
        Btn_RightPage.gameObject:SetActive(true);
    end
    ]] --
end

local RefreshServerListTime = 60
local RefresButtonTime = 60

local AdPageCount = 0
local AdCurrentIndex = 0
local RefreshAdTime = 6
local AdStartIndex = 0
local AdEndIndex = 0
local AdLeftSize = 0
local AdRightSize = 0
--翻页动作：0下一页 1前一页
local AdNextPageBehavior = 0
local function onAdPageChanged(pageCount, currentIndex)

    AdPageCount = pageCount
    AdNextPageBehavior = (AdCurrentIndex == currentIndex) and 0 or 1
    AdCurrentIndex = currentIndex
    local adl = G.User.GetAdList()
    if adl == nil or #adl == 0 then
        local AD1 = AdScrr_GameScrollView.transform:Find("GameListContent/_AD1")
        AD1.gameObject:SetActive(false)
        return
    end

    --设定对应广告行为与连结
    CurrentUrl = adl[currentIndex + 1].Link
    adtype = adl[currentIndex + 1].Behavior
    -- 切换页签
    if (currentIndex == 1) then
        AdRightSize = 1;
    end
    -- print("<color=yellow>onAdPageChanged</color>")

    if (pageCount > 1) then
        PageSwitch(pageCount, currentIndex);
    end
end

local function DoRefreshAD()
    if (AdPageCount <= 1) then
        return;
    end
    local idx = AdCurrentIndex + 1
    local ScrollPage = AdScrr_GameScrollView:GetComponent("ScrollPage")
    if idx >= AdPageCount then
        for i = 0, AdPageCount - 1 do
            ScrollPage:ToPrePage()
        end
    else
        ScrollPage:ToNextPage()
    end
end

local function RefreshAD()

    this:CancelInvokeRepeat(DoRefreshAD)
    this:InvokeRepeating(DoRefreshAD, 0, RefreshAdTime)
end

function PageSwitch(pageCount, currentIndex)
    -- print("<color=yellow>PageSwitch</color>")

    if ((DoRefreshAD ~= nil)) then
        -- print("<color=yellow>"..tostring(DoRefreshAD==nil).."</color>") 
        this:CancelInvokeRepeat(DoRefreshAD)
        this:InvokeRepeating(DoRefreshAD, 6, RefreshAdTime)
    end

    local sequence = CS.DG.Tweening.DOTween.Sequence()
    if (pageCount <= 6) then
        --差异值！调整广告起使位置的差异数值
        local v = AdStartIndex - 1
        for i = AdStartIndex, AdEndIndex, 1 do
            RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i) .. "/cur"):GetComponent("Image").gameObject:SetActive((i == (currentIndex + 1 + v)))
        end
    else
        -------------------------------------------------- [总数大于6笔滑动处理]-------Start
        -------------- [纪录原始页签座标]-------Start
        local nativepos = {}
        for i = 1, 6, 1 do
            local o = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
            nativepos[i] = CS.UnityEngine.Vector3(o.transform.position.x, o.transform.position.y, o.transform.position.z)
        end
        -------------- [纪录原始页签座标]-------End
        -------------- [LOG]
        -- print("<color=purple>now page-------->"..tostring(currentIndex+1).."</color>")
        -------------- [当往下滑动一页等于最后一颗页签]-------Start 
        if ((currentIndex + 1) == AdEndIndex) then


            -- print("<color=yellow>Last Page</color>")


            -------------- [当往下滑动一页等于最后一页]-------Start 
            if ((currentIndex + 1) == AdPageCount) then
                RecTran_CarouselRoot.transform:Find("_offset/p5/cur"):GetComponent("Image").gameObject:SetActive(false)
                RecTran_CarouselRoot.transform:Find("_offset/p6/cur"):GetComponent("Image").gameObject:SetActive(true)
                -- print("<color=yellow>Final Page</color>")
                -------------- [当往下滑动一页等于最后一页]-------END 
            else
                -------------- [当往下滑动一页不是最后一页]-------Start
                -- print("<color=yellow>Need Chage Page</color>",AdStartIndex,AdEndIndex)
                AdEndIndex = AdEndIndex + 1
                AdStartIndex = AdStartIndex + 1
                ---------------------- [开始同步往左移动]-------------START
                local RecTran_CarouselRoot = this.transform:Find("_offset/pr");
                pr.gameObject:SetActive(true);

                pr.transform.localScale = CS.UnityEngine.Vector3(0.1, 0.1, 1)
                for i = 1, 6, 1 do
                    local o = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
                    local target = CS.UnityEngine.Vector3(o.transform.position.x - 34, o.transform.position.y, o.transform.position.z)
                    if (i == 1) then
                        sequence:Append(o.transform:DOMove(target, 0.2))
                    else
                        sequence:Join(o.transform:DOMove(target, 0.2))
                    end
                end
                --当最后一颗是最后一页时尺寸回覆1
                if (AdEndIndex == AdPageCount) then
                    sequence:Join(pr.transform:DOScale(1, 0.2))
                else
                    sequence:Join(pr.transform:DOScale(0.5, 0.2))
                end
                ---------------------- 开始同步往左移动-------------END                

                --隐藏演示的最后一颗与设定为远本座标 AppendCallback伫列（LIFO）储存 
                ---------------------- 移动结束,设定回为真实页签-------------START  
                sequence:AppendCallback(function()

                    RecTran_CarouselRoot.transform:Find("_offset/p6/cur").gameObject:SetActive(false);
                    RecTran_CarouselRoot.transform:Find("_offset/p5/cur").gameObject:SetActive(true);
                    --回归座标
                    for i = 1, 6, 1 do
                        local o = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
                        o.transform.position = CS.UnityEngine.Vector3(nativepos[i].x, nativepos[i].y, nativepos[i].z)
                    end
                    --关闭演示用页签
                    pr.gameObject:SetActive(false);
                    --当最后一颗是最后一页时尺寸回覆1

                    if (AdStartIndex ~= 1) then
                        RecTran_CarouselRoot.transform:Find("_offset/p1").transform.localScale = CS.UnityEngine.Vector3(0.5, 0.5, 1)
                    end


                    if (AdEndIndex == AdPageCount) then
                        RecTran_CarouselRoot.transform:Find("_offset/p6").transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                    else
                        RecTran_CarouselRoot.transform:Find("_offset/p6").transform.localScale = CS.UnityEngine.Vector3(0.5, 0.5, 1)
                    end

                    -- print("<color=yellow>reset2</color>")
                end)
                ---------------------- 移动结束,设定回为真实页签-------------END  

                --指定第六颗亮起,第五颗暗下
                sequence:AppendCallback(function()
                    RecTran_CarouselRoot.transform:Find("_offset/p5/cur"):GetComponent("Image").gameObject:SetActive(false)
                    RecTran_CarouselRoot.transform:Find("_offset/p6/cur"):GetComponent("Image").gameObject:SetActive(true)
                    -- print("<color=yellow>reset1</color>")
                end)
                -------------- [当往下滑动一页不是最后一页]-------END   
            end
            -------------- [当往下滑动一页等于最后一颗页签]-------END 
            -------------- [当往前滑动一页等于第一颗页签]-------START 
        elseif ((currentIndex + 1) == AdStartIndex) then

            -- print("<color=yellow>First Page</color>")
            if ((currentIndex + 1) == 1) then
                -- print("<color=yellow>Top Page</color>")
                sequence:AppendInterval(0.3)
                sequence:AppendCallback(function()
                    -- print("<color=yellow>Top Page------RUN</color>")
                    RecTran_CarouselRoot.transform:Find("_offset/p2/cur"):GetComponent("Image").gameObject:SetActive(false)
                    RecTran_CarouselRoot.transform:Find("_offset/p1/cur"):GetComponent("Image").gameObject:SetActive(true)
                    RecTran_CarouselRoot.transform:Find("p6/cur"):GetComponent("Image").gameObject:SetActive(false)
                end)
            else
                -- print("<color=yellow>Need Chage Page</color>")               
                AdEndIndex = AdEndIndex - 1
                AdStartIndex = AdStartIndex - 1

                ---------------------- [开始同步往左移动]-------------START
                local pl = RecTran_CarouselRoot.transform:Find("_offset/pl");
                pl.gameObject:SetActive(true);

                pl.transform.localScale = CS.UnityEngine.Vector3(0.1, 0.1, 1)
                for i = 1, 6, 1 do
                    local o = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
                    local target = CS.UnityEngine.Vector3(o.transform.position.x + 34, o.transform.position.y, o.transform.position.z)
                    if (i == 1) then
                        sequence:Append(o.transform:DOMove(target, 0.2))
                    else
                        sequence:Join(o.transform:DOMove(target, 0.2))
                    end
                end
                --当最后一颗是最后一页时尺寸回覆1
                if (AdStartIndex == 1) then
                    sequence:Join(pl.transform:DOScale(1, 0.2))
                else
                    sequence:Join(pl.transform:DOScale(0.5, 0.2))
                end
                ---------------------- 开始同步往左移动-------------END                

                --隐藏演示的第一颗与设定为远本座标 AppendCallback伫列（LIFO）储存 
                ---------------------- 移动结束,设定回为真实页签-------------START  
                sequence:AppendCallback(function()

                    RecTran_CarouselRoot.transform:Find("_offset/p1/cur").gameObject:SetActive(false);
                    RecTran_CarouselRoot.transform:Find("_offset/p2/cur").gameObject:SetActive(true);
                    --回归座标
                    for i = 1, 6, 1 do
                        local o = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
                        o.transform.position = CS.UnityEngine.Vector3(nativepos[i].x, nativepos[i].y, nativepos[i].z)
                    end
                    --关闭演示用页签
                    pl.gameObject:SetActive(false);

                    if (AdEndIndex ~= AdPageCount) then
                        RecTran_CarouselRoot.transform:Find("_offset/p6").transform.localScale = CS.UnityEngine.Vector3(0.5, 0.5, 1)
                    end

                    --当第一颗是第一页时尺寸回覆1
                    if (AdStartIndex == 1) then
                        RecTran_CarouselRoot.transform:Find("_offset/p1").transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                    else
                        RecTran_CarouselRoot.transform:Find("_offset/p1").transform.localScale = CS.UnityEngine.Vector3(0.5, 0.5, 1)
                    end

                    -- print("<color=yellow>reset2</color>")
                end)
                ---------------------- 移动结束,设定回为真实页签-------------END  

                --指定第1颗亮起,第5颗暗下
                sequence:AppendCallback(function()
                    RecTran_CarouselRoot.transform:Find("_offset/p2/cur"):GetComponent("Image").gameObject:SetActive(false)
                    RecTran_CarouselRoot.transform:Find("_offset/p1/cur"):GetComponent("Image").gameObject:SetActive(true)
                    -- print("<color=yellow>reset1</color>")
                end)
                -------------- [当往前滑动一页不是第一页]-------END
            end
            -------------- [当往前滑动一页等于第一颗页签]-------END
        else
            -------------- [当往前往后滑动一页不是最前或最后时]-------START
            sequence:AppendCallback(function()
                --重新设定选择页签
                local j = 1;
                for i = AdStartIndex, AdEndIndex, 1 do
                    RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(j) .. "/cur"):GetComponent("Image").gameObject:SetActive((i == (currentIndex + 1)))
                    j = j + 1;
                end
            end)
            -------------- [当往前往后滑动一页不是最前或最后时]-------END    
        end
    end
end

function InitPager(pagecount)
    -- print("<color=yellow>InitPager</color>")
    AdStartIndex = (pagecount <= 4) and 2 or 1;

    if (pagecount <= 4) then
        AdEndIndex = pagecount + 1;
    elseif (pagecount == 5) then
        AdEndIndex = 5;
    else
        AdEndIndex = 6;
    end

    for i = AdStartIndex, AdEndIndex do
        local p = RecTran_CarouselRoot.transform:Find("_offset/p" .. tostring(i))
        p.gameObject:SetActive(true);
    end
end

local annShadeWidth = Txt_NoticeContent.transform.parent:GetComponent("RectTransform").rect.width -- 公告文本遮罩宽度 
local annInitPos = Txt_NoticeContent.transform.localPosition
local RectTran_NoticeContent = Txt_NoticeContent.transform:GetComponent("RectTransform")

function GetNoticeContentWidth()
    return CS.UnityEngine.UI.LayoutUtility.GetPreferredWidth(RectTran_NoticeContent)
end

local Check_UserData = nil
local Check_Mail = nil
local Check_ExchangeLog = nil
local Check_LoginBonus = nil
local Check_CSData = nil
local Check_VipData = nil

-- 拦截邮件消息
G.Service.Mail.RegMailRetCallBack(function()
    if Check_Mail ~= nil then
        Check_Mail.SetGetDataFinish()
    end
    UpdateElement()
end)
-- 拦截邮件消息
G.Service.CS.RegCSRetCallBack(function()
    if Check_CSData ~= nil then
        Check_CSData.SetGetDataFinish()
    end
    UpdateCS()
end)

G.Service.ExchangeLog.RegLogRetCallBack(function()
    if Check_ExchangeLog ~= nil then
        Check_ExchangeLog.SetGetDataFinish()
    end
end)

G.Service.Vip.RegRetCallBack(function()
    if Check_VipData ~= nil then
        Check_VipData.SetGetDataFinish()
		uiCommon.SetFrameIcon(G.User.GetCurrentUser().UserData.AvatarFrame)
    end
end)

function SendGetLoginBonus()
    G.Service.LoginBonus.SendRequest(function()
        if Check_LoginBonus ~= nil then
            Check_LoginBonus.SetGetDataFinish()
        end

        if G.Service.LoginBonus.IsOpenLoginBonus() == false then
            Btn_Notice.gameObject:SetActive(false)
        else
            Btn_Notice.gameObject:SetActive(true)
        end

        RefreshBottomAttach()
    end)
end

local isGetAllDataFinish = false
function CheckGetData()
    if isGetAllDataFinish == false then
        if Check_UserData ~= nil then
            Check_UserData.UpdateGetData()
        end

        if Check_ExchangeLog ~= nil then
            Check_ExchangeLog.UpdateGetData()
        end

        if Check_Mail ~= nil then
            Check_Mail.UpdateGetData()
        end
        if Check_CSData ~= nil then
            Check_CSData.UpdateGetData()
        end

        if Check_LoginBonus ~= nil then
            Check_LoginBonus.UpdateGetData()
        end
		
        if Check_VipData ~= nil then
            Check_VipData.UpdateGetData()
        end		

        if (Check_UserData.IsGetDataFinish() == true and
                Check_ExchangeLog.IsGetDataFinish() == true and
                Check_Mail.IsGetDataFinish() == true and
                Check_LoginBonus.IsGetDataFinish() == true and
                Check_CSData.IsGetDataFinish() == true and
                Check_VipData.IsGetDataFinish() == true) then

			isGetAllDataFinish = true

            -- 關閉回報客服視窗
            UICustomer = CS.WindowManager.Instance:FindWindowNodeByName("UICustomer")
            if UICustomer ~= nil then
                CS.UnityEngine.GameObject.Destroy(UICustomer.WindowGameObject)
            end

            --CS.UILoading.Hide()
            CheckShowMail()
        end
	elseif(IsGetServerList == true)then
		CS.UILoading.Hide()
    end
end

function DoCheckNavigation()
    panelwidth = Scrr_GameScrollView:GetComponent("RectTransform").rect.width;

    nx = RecTran_GameAttachList.localPosition.x;
    p = math.ceil(activeGameNumber / 6)
    x1 = panelwidth - GameListWidth;

    t1, t2 = math.modf(nx / (GameListWidth - (200 + x1)));
    nowp = t1

    --Btn_LeftPage.gameObject:SetActive(nowp~=0)		
    --Btn_RightPage.gameObject:SetActive((nowp > -(p-1)))
end

function Update()

    --開關巡航頁按鈕
    DoCheckNavigation()
    -- 更新熱更百分比
    DoHotUpdate()
	
    -- 檢查所有資料是否下載完畢
    CheckGetData()

    -- 檢查是否要過serverlist
    CheckServerList()

    -- 送審階段關掉跑馬燈 By Rayon
    if G.GameConfig.IsSubmissionCtx() == true then
        return
    end

    if G.IsGameStart() ~= true then
        if (#Txt_NoticeContent.text > 0 and Txt_NoticeContent.transform.localPosition.x > annInitPos.x - GetNoticeContentWidth() - annShadeWidth) then

            Txt_NoticeContent.transform.localPosition = Txt_NoticeContent.transform.localPosition - CS.UnityEngine.Vector3(CS.UnityEngine.Time.deltaTime * 50, Txt_NoticeContent.transform.localPosition.y, Txt_NoticeContent.transform.localPosition.z)
        else
            local message = G.Service.Ann.GetOneAnnMeg()

            if message == nil then
                RecTran_Ann.gameObject:SetActive(false)
                Txt_NoticeContent.text = ""
            else
                message = string.gsub(message, "%s", " ")
                Txt_NoticeContent.transform.localPosition = annInitPos
                Txt_NoticeContent.text = message
                RecTran_Ann.gameObject:SetActive(true)
            end
        end
    end
end

function Start()




    uiRank.SetCallbackBack(EscapeWindow)
    uiRank.SetCallbackRecharge(OpenRechargeUI)

    uiCommon.SetCallbackHeadIcon(OpenUserInfo)
    uiCommon.SetCallbackRecharge(OpenRechargeUI)
    --G.Service.Msg.Init()

    CalculateGameScrollView()
	Img_LogoSprite[2] = this.transform:Find("_Main/TOP/_Logo/Image2"):GetComponent("Image")
    Img_LogoSprite[170] = this.transform:Find("_Main/TOP/_Logo/Image170"):GetComponent("Image")
    Img_LogoSprite[172] = this.transform:Find("_Main/TOP/_Logo/Image172"):GetComponent("Image")
    Img_LogoSprite[174] = this.transform:Find("_Main/TOP/_Logo/Image174"):GetComponent("Image")

    -- 送審階段不開啟 Add By Rayon
    if G.GameConfig.IsSubmissionCtx() == true then
        -- 隱藏下方按鈕
        local _Btn = this.transform:Find("_Main/Bottom")
        if _Btn ~= nil then
            _Btn.gameObject:SetActive(false)
        end

        Btn_Coffer.gameObject:SetActive(false)
        Btn_Exchange.gameObject:SetActive(false)
        Btn_Upgrade.gameObject:SetActive(false)
    end
    -- 兌換文字 ExchangeFontType
    Btn_ExchangeFont1.gameObject:SetActive(tonumber(G.GameDefine.ExchangeFontType) == 0)
    Btn_ExchangeFont2.gameObject:SetActive(tonumber(G.GameDefine.ExchangeFontType) == 1)
    Btn_ExchangeFont3.gameObject:SetActive(tonumber(G.GameDefine.ExchangeFontType) == 2)
    Btn_ExchangeFont4.gameObject:SetActive(tonumber(G.GameDefine.ExchangeFontType) == 3)

    -- 新代理
    if G.Platform.getJPYL() == false then
        Btn_Promo.gameObject:SetActive(G.GameDefine.Club2Agent == 1)
        Btn_AgentTop.gameObject:SetActive(false)

        Img_Logo.gameObject:SetActive(false)
    else
        Btn_Promo.gameObject:SetActive(false)
        Btn_AgentTop.gameObject:SetActive(G.GameDefine.Club2Agent == 1)

        Img_Logo.gameObject:SetActive(true)
    end
    --Btn_Promo.gameObject:SetActive(G.GameDefine.Club2Agent == 1)
    --Btn_Proxy.gameObject:SetActive(G.GameDefine.Club2Agent == 1)
    local logoSprite = G.Platform.getHallLogo()
    if(logoSprite > 0)then
        if(Img_LogoSprite[logoSprite] ~= nil)then
            Img_Logo.gameObject:SetActive(true)
            Img_Logo.sprite = Img_LogoSprite[logoSprite].sprite
        else
            SetLogo()
        end
	else
		Img_Logo.gameObject:SetActive(false)
    end
    
    -- 商城文字
    local shopNameID = userdata.ShopNameID
    Img_RechargeFont1.gameObject:SetActive(shopNameID == 1)
    Img_RechargeFont2.gameObject:SetActive(shopNameID == 2)
    Img_RechargeFont3.gameObject:SetActive(shopNameID == 3)
    Img_RechargeFont4.gameObject:SetActive(shopNameID == 4)

    --彩金池是否開啟
    btn__MoneyPool.gameObject:SetActive(G.GameDefine.ShowCaiJinchi == 1)

    -- 充值開關 Add By Rayon
    --[[
    Btn_GetMoney.gameObject:SetActive(G.GameConfig.IsRechargeCtx())
    ]]
    uiCommon.ShowRecharge(G.GameConfig.IsRechargeCtx())
    Btn_Recharge.gameObject:SetActive(G.GameConfig.IsRechargeCtx())
    --[[
    Btn_RoomRecharge.gameObject:SetActive(G.GameConfig.IsRechargeCtx())
    ]]

    UpdateElement()
    UpdateCS()
    --todo
    --onGamePageChanged(1,0)
    --local scrollPage = Scrr_GameScrollView:GetComponent("ScrollPage")
    --scrollPage.OnPageChanged = scrollPage.OnPageChanged + onGamePageChanged
    onAdPageChanged(1, 0)
    local AdscrollPage = AdScrr_GameScrollView:GetComponent("ScrollPage")
    AdscrollPage.OnPageChanged = AdscrollPage.OnPageChanged + onAdPageChanged

    for i = 0, RecTran_Rooms.childCount - 1, 1 do
        local child = RecTran_Rooms:GetChild(i)
        local index = tonumber(child.name)
        RoomList[index] = RoomList[index] or {}
        RoomList[index].base = child
        RoomList[index].base:Find("UIBody/Icon").gameObject:SetActive(false)
    end

    for i = 0, RecTran_RoomsBy.childCount - 1, 1 do
        local child = RecTran_RoomsBy:GetChild(i)
        local index = tonumber(child.name)
        RoomList[index] = RoomList[index] or {}
        RoomList[index].by = child
    end
    for i = 0, RecTran_RoomsDdz.childCount - 1, 1 do
        local child = RecTran_RoomsDdz:GetChild(i)
        local index = tonumber(child.name)
        RoomList[index] = RoomList[index] or {}
        RoomList[index].ddz = child
    end

    RecTran_GameList.gameObject:SetActive(false)

    --[[for i = 0, RecTran_GameList.childCount - 2, 1 do
        GameLogoList[i + 1] = RecTran_GameList:GetChild(i)
    end]]
    for i = 1, G.Game.GetGameIDCount(), 1 do
        GameLogoList[i] = RecTran_GameList.transform:Find("GameLogo_" .. i):GetComponent("RectTransform")
    end

    for i = 1, 50, 1 do
        table.insert(iconNameList, "")
    end

    this:InvokeRepeating(RefreshServerList, 0, RefreshServerListTime)
    this:InvokeRepeating(UpdateButton, 0, RefresButtonTime)

    G.Service.ReCharge.GetRechargeList()

    Btn_agent.gameObject:SetActive(G.GameDefine.Agent == 1)

    RefreshBottomAttach()

    Btn_agent:AddClickListener(function()

        if (userdata.ClubOpen) then
            local openparam = CS.WindowNodeInitParam("UIAgent", CS.AppDefine.HallNumber)
            CS.WindowManager.Instance:OpenWindow(openparam)
        else
            CS.UIBubblePrompt.Show("代理功能近期上线,敬请期待", groupName);
        end
    end)

	local function OpenPromoUI()
        if (userdata.ClubOpen) then
			CS.UILoadingData.Show("载入中,请稍后...", 5)

            local openparam = CS.WindowNodeInitParam("UIProxy", CS.AppDefine.HallNumber)
            openparam.WindowCloseCallBack = (function()
                RefreshUserInfo()
            end)
            CS.WindowManager.Instance:OpenWindow(openparam)
        else
            CS.UIBubblePrompt.Show("代理功能近期上线,敬请期待", groupName);
        end	
	end
	
	btn__VipIcon:AddClickListener(function()
        local openparam = CS.WindowNodeInitParam("UIVipInfo", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
	end)
	
    Btn_Promo:AddClickListener(function()
        OpenPromoUI()
    end)

    Btn_AgentTop:AddClickListener(function()
		OpenPromoUI()
    end)

    Btn_RoomlistBack:AddClickListener(function()
        EscapeWindow()
    end)
    Btn_AdScrr:AddClickListener(function()
				
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 31,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        -- 0:複製,1:直連該網頁,2:一鍵入廳,3:甚麼都不做,4:前往充值,5:评分(IOS),7:開啟IOS儲存網頁 8:活動 9：首充返利
        if (adtype == 0) then
            CS.Utility.CopyText(CurrentUrl)
            CS.UIBubblePrompt.Show("复制成功", groupName);
        elseif (adtype == 1) then
            CS.UnityEngine.Application.OpenURL(CurrentUrl)
        elseif (adtype == 2) then
            local gameId = tonumber(CurrentUrl);
            local nversion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId))
            local rversion = GameVersion[gameId]
            --当尚未下载先下载游戏
            if (CheckGameIsDownLoad(gameId) and (VersionToNum(nversion, gameId) < VersionToNum(rversion, gameId))) then
                StartDownloadGame(gamelist[gameId], gameId)
            else

                if (not autolist[gameId]) then
                    local name = "GameTitle_" .. gameId
                    -- 一鍵入廳設定切換入廳文字
                    local userGameList = nil;
                    if ((G.User.GetUserGameList2() ~= nil)) then
                        userGameList = G.User.GetUserGameList2()
                    else
                        userGameList = G.User.GetUserGameList()
                    end

                    local fontVer = tonumber(G.Tool.TableHelper.Find(userGameList, function(v) return v.GameID == gameId end)[1].IconName)

                    if fontVer ~= 0 then
                        name = name .. "_" .. fontVer
                    end
                    G.Tool.Debug.RecordGet({
                        content = {
                            ContentType = 1,
                            Flag = 32,
                            Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
                        },
                        lt = 1,
                        sceneName = "UIHALL",
                        uid = G.User.GetCurrentUser().UserID
                    })
                    Img_RoomTitle:ResetSpriteByName(CS.AppDefine.HallNumber, name, true);
                    GameClick(gameId, gamelist[gameId]);
                else
                    G.Tool.Debug.RecordGet({
                        content = {
                            ContentType = 1,
                            Flag = 33,
                            Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
                        },
                        lt = 1,
                        sceneName = "UIHALL",
                        uid = G.User.GetCurrentUser().UserID
                    })
                    -- 自動入桌
                    RoomClick(autolist[gameId], 1);
                end
            end



        elseif (adtype == 3) then
            --甚麼多不做
        elseif (adtype == 4) then
            --Go Recharge
            OpenRechargeUI();

        elseif (adtype == 5) then
            if CurrentUrl ~= nil and CurrentUrl ~= '' then
                CS.UnityEngine.Application.OpenURL(CurrentUrl)
            end
        elseif (adtype == 6) then
            CS.UILoadingData.Show("正在请求资料");
            G.Service.Agent.GetAgentData(function(content)
                CS.UILoadingData.Hide();
                CS.Utility.CopyText(content.AgentDetailData.OfficeUser1)
                CS.UIBubblePrompt.Show("复制成功", groupName)
                CS.UnityEngine.Application.OpenURL("weixin://")
            end)
        elseif (adtype == 7) then
            if CurrentUrl ~= nil and CurrentUrl ~= '' then
                PlayCellPhoneAni()
            end
		elseif (adtype == 8)then 
			CheckShowActivity (CurrentUrl, nil)
		elseif (adtype == 9)then 
			ShowFirstRechargeRebateUI()
        end
    end)
    --[[
        Btn_HeadIcon:AddClickListener
        (
        function()
            G.Tool.Debug.RecordGet({
                                content= {  
                                            ContentType=1,
                                            Flag=20,
                                            Data=(string.format("%.3f",  CS.UnityEngine.Input.mousePosition.x)..","..string.format("%.3f",CS.UnityEngine.Input.mousePosition.y))
                                         },
                                lt=1,
                                sceneName="UIHALL",
                                uid=G.User.GetCurrentUser().UserID
                                })
            
            local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
    
            openparam.WindowCloseCallBack =
            (
            function()
                RefreshUserInfo()
            end
            )
    
            CS.WindowManager.Instance:OpenWindow(openparam)
        end
        )
    ]]
    Btn_Recharge:AddClickListener(function()
        --纪录LOG ONLINE
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 21,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        OpenRechargeUI()
    end)
    --[[
    Btn_RoomRecharge:AddClickListener
    (
    function()
         --纪录LOG ONLINE
        G.Tool.Debug.RecordGet({
                            content= {  
                                        ContentType=1,
                                        Flag=22,
                                        Data=(string.format("%.3f",  CS.UnityEngine.Input.mousePosition.x)..","..string.format("%.3f",CS.UnityEngine.Input.mousePosition.y))
                                     },
                            lt=1,
                            sceneName="UIHALL",
                            uid=G.User.GetCurrentUser().UserID
                            })
        OpenRechargeUI()
    end
    )
]]
    btn__MoneyPool:AddClickListener(function()
        local openparam = CS.WindowNodeInitParam("UIMoneyPool", CS.AppDefine.HallNumber)
        openparam.WindowData = moneypool;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end)

    Btn_Exchange:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 23,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
		
		-- 还没到达显示条件 跳出讯息
		local IsShow = IsShowExchangeButton()
		if IsShow == false then
			if userdata.ConvertPopupText ~= nil then
				G.Util.coroutine_call(function()
					showMessageBox(CS.MessageBoxStyle.OK, "", userdata.ConvertPopupText, "确定");
				end)()
			end
			return
		end
		
        if (user.FlowInfo.Status == 1 and user.FlowInfo.UnionPlatform == 1) then
            if (user.FlowInfo.Achievement == 0) then
                img_Bonus.gameObject:SetActive(true);
                local sequence = CS.DG.Tweening.DOTween.Sequence()
                img_Bonus.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)

                sequence:AppendInterval(1.2)
                sequence:Append(img_Bonus.transform:DOScale(0.1, 0.3))

                sequence:AppendCallback(function()
                    img_Bonus.gameObject:SetActive(false);
                end)
                return
            end
        end



        if (G.Service.BindingTest.IsBindingTestMode()) then --add by Bob

            if G.Service.BindingTest.IsBingdingPhone() ~= true then

                G.Util.coroutine_call(function()
                    local showMessage = "兑换需绑定用户信息，确认绑定吗?"
                    if tonumber(G.GameDefine.ExchangeFontType) == 1 then
                        showMessage = "福利需绑定用户信息，确认绑定吗？"
                    elseif tonumber(G.GameDefine.ExchangeFontType) == 2 then
                        showMessage = "清算需绑定用户信息，确认绑定吗？"
                    elseif tonumber(G.GameDefine.ExchangeFontType) == 3 then
                        showMessage = "提现需绑定用户信息，确认绑定吗？"
                    end

                    local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", showMessage, "绑定", "取消");
                    if boxResult == CS.MessageBoxResult.OK then
                        local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
                        openparam.WindowData = 1;
                        CS.WindowManager.Instance:OpenWindow(openparam)
                    end
                end)()
            else
                -- CS.UIBubblePrompt.Show("功能尚未开放，敬请期待。", "UIHall")

                local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)

                openparam.WindowCloseCallBack =
                (function()
                    G.Service.Login.RequestUserInfoHall()
                end)

                CS.WindowManager.Instance:OpenWindow(openparam)
            end


        else
            if userdata.IsBindGuest ~= true then

                G.Util.coroutine_call(function()
                    local showMessage = "需要绑定用户信息，确认绑定吗？"
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", showMessage, "绑定", "取消");
                    if boxResult == CS.MessageBoxResult.OK then
                        local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
                        openparam.WindowData = 1;
                        CS.WindowManager.Instance:OpenWindow(openparam)
                    end
                end)()
            else
                -- CS.UIBubblePrompt.Show("功能尚未开放，敬请期待。", "UIHall")

                local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)

                openparam.WindowCloseCallBack =
                (function()
                    G.Service.Login.RequestUserInfoHall()
                end)

                CS.WindowManager.Instance:OpenWindow(openparam)
            end
        end
    end)

    Btn_Mail:AddClickListener(function()

        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 25,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        local openparam = CS.WindowNodeInitParam("UIMail", CS.AppDefine.HallNumber)

        openparam.WindowCloseCallBack =
        (function()
            RefreshUserInfo()
            UpdateElement()
        end)


        CS.WindowManager.Instance:OpenWindow(openparam)
    end)

    Btn_Coffer:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 26,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        local openparam = CS.WindowNodeInitParam("UIBank", CS.AppDefine.HallNumber)

        openparam.WindowCloseCallBack =
        (function()
            RefreshUserInfo()
        end)


        CS.WindowManager.Instance:OpenWindow(openparam)
    end)


    Btn_Setting:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 27,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        local openparam = CS.WindowNodeInitParam("UISetting", CS.AppDefine.HallNumber)
        CS.WindowManager.Instance:OpenWindow(openparam)
    end)
    --[[
    Btn_GetMoney:AddClickListener
    (
    function()
        G.Tool.Debug.RecordGet({
                            content= {  
                                        ContentType=1,
                                        Flag=28,
                                        Data=(string.format("%.3f",  CS.UnityEngine.Input.mousePosition.x)..","..string.format("%.3f",CS.UnityEngine.Input.mousePosition.y))
                                     },
                            lt=1,
                            sceneName="UIHALL",
                            uid=G.User.GetCurrentUser().UserID
                            })
        OpenRechargeUI()
    end
    )
    ]]
    --Btn_Notice.gameObject:SetActive(false)
    Btn_Notice:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 29,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        G.Service.LoginBonus.SendRequest(function()
            local CashAmount = userdata.CashAmount
            local openparam = CS.WindowNodeInitParam("UINotice", CS.AppDefine.HallNumber)
            openparam.WindowCloseCallBack = (function()
                UpdateElement()

                -- 充值頁面關閉時重刷大廳資訊
                if userdata.CashAmount > CashAmount then
                    G.Service.Login.RequestUserInfoHall()
                else
                    RefreshUserInfo()
                end
            end)
            CS.WindowManager.Instance:OpenWindow(openparam)
        end)
    end)

    Btn_Upgrade:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 19,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        CheckShowBingdingGift(true)
        --[[local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)

        openparam.WindowData = 1;

        openparam.WindowCloseCallBack =
        (
        function()
            RefreshUserInfo()
        end
        )

        CS.WindowManager.Instance:OpenWindow(openparam)]]
    end)

    UpdateButton()
    Btn_CustomService:AddClickListener(function()
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 1,
                Flag = 18,
                Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        local openparam = CS.WindowNodeInitParam("UISendMsgToGM", CS.AppDefine.HallNumber)
        CS.WindowManager.Instance:OpenWindow(openparam)
    end)



    Btn_LeftPage:AddClickListener(function()

        nx = RecTran_GameAttachList.localPosition.x;
        p = math.ceil(activeGameNumber / 6)
        t1, t2 = math.modf(nx / GameListWidth);
        nowp = t1
        nextp = (nowp == 0) and 0 or (nowp + 1);
        -- print("<color=yellow>".. nx .."</color>")
        -- print("<color=yellow>activeGameNumber".. activeGameNumber .."</color>")
        -- print("<color=yellow>p".. p .."</color>")
        -- print("<color=yellow>nowp ".. nowp .."</color>")
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        --第一頁是0第二頁是900所以p-1
        if (nx < 0) and ((nextp * GameListWidth) <= 0) then

            target = CS.UnityEngine.Vector3(nextp * GameListWidth, RecTran_GameAttachList.localPosition.y, RecTran_GameAttachList.localPosition.z)
            sequence:Append(RecTran_GameAttachList.transform:DOLocalMove(target, 0.2))
            --RecTran_GameAttachList.localPosition=CS.UnityEngine.Vector3(nextp*900,RecTran_GameAttachList.localPosition.y,RecTran_GameAttachList.localPosition.z)
        else
            target = CS.UnityEngine.Vector3(0, RecTran_GameAttachList.localPosition.y, RecTran_GameAttachList.localPosition.z)
            sequence:Append(RecTran_GameAttachList.transform:DOLocalMove(target, 0.2))
            --RecTran_GameAttachList.localPosition=CS.UnityEngine.Vector3(0,RecTran_GameAttachList.localPosition.y,RecTran_GameAttachList.localPosition.z)
        end

        --todo
        -- Scrr_GameScrollView:GetComponent("ScrollPage"):ToPrePage();
    end)

    Btn_RightPage:AddClickListener(function()
        nx = RecTran_GameAttachList.localPosition.x;
        p = math.ceil(activeGameNumber / 6)

        t1, t2 = math.modf(nx / GameListWidth);
        nowp = t1
        nextp = (nowp == 0) and -1 or (nowp - 1);
        -- print("<color=yellow>".. nx .."</color>")
        -- print("<color=yellow>activeGameNumber".. activeGameNumber .."</color>")
        -- print("<color=yellow>p".. p .."</color>")
        -- print("<color=yellow>nowp ".. nowp .."</color>")
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        --第一頁是0第二頁是900所以p-1
        if (nx > ((p - 1) * -GameListWidth)) and ((nextp * GameListWidth) >= ((p - 1) * -GameListWidth)) then
            target = CS.UnityEngine.Vector3(nextp * GameListWidth, RecTran_GameAttachList.localPosition.y, RecTran_GameAttachList.localPosition.z)
            sequence:Append(RecTran_GameAttachList.transform:DOLocalMove(target, 0.2))
            --RecTran_GameAttachList.localPosition=CS.UnityEngine.Vector3(nextp*900,RecTran_GameAttachList.localPosition.y,RecTran_GameAttachList.localPosition.z)
        else
            target = CS.UnityEngine.Vector3(((p - 1) * -GameListWidth), RecTran_GameAttachList.localPosition.y, RecTran_GameAttachList.localPosition.z)
            sequence:Append(RecTran_GameAttachList.transform:DOLocalMove(target, 0.2))
            --RecTran_GameAttachList.localPosition=CS.UnityEngine.Vector3(((p-1) * -900),RecTran_GameAttachList.localPosition.y,RecTran_GameAttachList.localPosition.z)
        end

        --Scrr_GameScrollView:GetComponent("ScrollPage"):ToNextPage();
    end)

    -- 皜豢??剔瑪???斗
    if userdata.IsGaming and userdata.LockGameID ~= 0 and userdata.ServerAddr ~= "" then
        G.Util.coroutine_call(function()
            G.Tool.Debug.RecordGet({
                content = {
                    ContentType = 1,
                    Flag = 17,
                    Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
                },
                lt = 1,
                sceneName = "UIHALL",
                uid = G.User.GetCurrentUser().UserID
            })
            local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "你正在房间[" .. G.Game.ReplaceGameName(userdata.ServerName) .. "]中进行游戏，是否进入房间继续游戏？", "继续", "取消");
            if boxResult == CS.MessageBoxResult.OK then
                -- 紀錄BUG，有時候userdata.ServerAddr拆分數字會有問題
                -- 原因LoginRetID中有ServerAddr值，後來傳來的UserHallInfoID中卻無ServerAddr值
                if userdata.ServerAddr ~= "" then
                    local addressParm = G.Tool.StringHelper.Split(userdata.ServerAddr, ":");
                    G.Server.GetGameServerAddrCallback(addressParm[1], addressParm[2], function(res, ipAddress, port)
                        if res == 0 then
                            RequestGameID = userdata.LockGameID
                            G.SetRunGameIPAndPort(ipAddress, port)

                            --寫入SKIN
                            local userGameList = nil;
                            if ((G.User.GetUserGameList2() ~= nil)) then
                                userGameList = G.User.GetUserGameList2()
                            else
                                userGameList = G.User.GetUserGameList()
                            end
                            local skin = G.Tool.TableHelper.Find(userGameList, function(v) return v.GameID == RequestGameID end)[1].Skin

                            print("去检查游戏版本------------------------------------")
                            G.GameConfig.RequestGameRemoteVersion = GameVersion[RequestGameID]

                            G.Service.NetworkService.CheckGameVersion(RequestGameID, skin)
                        else
                            CS.UIBubblePrompt.Show("游戏连接失败(1)", "UIHall")
                            require("Tools/TestConnect").loginError("遊戲盾取得遊戲IP失敗")
                        end
                    end);
                else
                    CS.UIBubblePrompt.Show("游戏连接失败", "UIHall")
                end
            end
        end)()
    end
end

function OpenRechargeUI()
    --[[
    if G.Service.UIFlag.IsIOSRecharge() then
        local openparam = CS.WindowNodeInitParam("UIRecharge_IOS", CS.AppDefine.HallNumber)
        openparam.WindowCloseCallBack =
        (
        function()
            RefreshUserInfo()
        end
        )
        CS.WindowManager.Instance:OpenWindow(openparam)
    else
    ]] --

    local openparam = nil;
    local CashAmount = userdata.CashAmount

    if G.GameConfig.IsSubmissionCtx() == true then
        openparam = CS.WindowNodeInitParam("UIRecharge_IOS", CS.AppDefine.HallNumber)
    else
        openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
    end
    openparam.WindowCloseCallBack =
    (function()
        -- 充值頁面關閉時重刷大廳資訊
        if userdata.CashAmount > CashAmount then
            G.Service.Login.RequestUserInfoHall()
        else
            RefreshUserInfo()
        end
    end)
    CS.WindowManager.Instance:OpenWindow(openparam)
    --end
end

function OpenUserInfo()
    G.Tool.Debug.RecordGet({
        content = {
            ContentType = 1,
            Flag = 20,
            Data = (string.format("%.3f", CS.UnityEngine.Input.mousePosition.x) .. "," .. string.format("%.3f", CS.UnityEngine.Input.mousePosition.y))
        },
        lt = 1,
        sceneName = "UIHALL",
        uid = G.User.GetCurrentUser().UserID
    })

    local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)

    openparam.WindowCloseCallBack =
    (function()
        RefreshUserInfo()
    end)

    CS.WindowManager.Instance:OpenWindow(openparam)
end

local CompetioniInfoRuleMsg = nil
function SendGetCompetioniInfoRule(GameID)
    local parm = {}
    parm.GameID = GameID
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetioniInfoRuleID, parm)
end

-- 比賽活動公告
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetioniInfoRuleRetID, function(receiveMsg)
    CompetioniInfoRuleMsg = receiveMsg
end, groupName);

function OpenRankHelpUI(GameID)
    local openparam = CS.WindowNodeInitParam("UIRankHelp", CS.AppDefine.HallNumber)
    openparam.WindowData = {}
    openparam.WindowData.GameID = GameID
    openparam.WindowData.receiveMsg = CompetioniInfoRuleMsg
    CS.WindowManager.Instance:OpenWindow(openparam)
end

function CheckShowMail()
    -- 判断是否弹出公告信息
    if userdata.AnnMsgShow ~= 1 then
        CheckMailEnd()
        return
    end

    local openparam = CS.WindowNodeInitParam("UIMail", CS.AppDefine.HallNumber)
    openparam.WindowData = 2
    openparam.WindowCloseCallBack = (function()
        RefreshUserInfo()
        UpdateElement()
        CheckMailEnd()
    end)

    --print("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",openparam.WindowData)
    CS.WindowManager.Instance:OpenWindow(openparam)
end

function CheckMailEnd()
	-- 有绑定 顯示每日登入獎勵
	if G.User.IsCurrentUserGuest() == true then
		CheckShowBingdingGift(false)
	-- 顯示充值反利
	else
		CheckShowHDRechageBonus()
	end
end

function CheckShowBingdingGift(showBinding)
    -- 送審階段不開啟 Add By Rayon
    if G.GameConfig.IsSubmissionCtx() == true then
        return
    end
	
	if(userdata.GiftMoney == 0)then
		-- 按鈕
		if showBinding == true then
			local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
			openparam.WindowData = 1;
			openparam.WindowCloseCallBack = (function()
				CheckShowHDRechageBonus()
			end)
			CS.WindowManager.Instance:OpenWindow(openparam)	
		-- 登入遊戲
		else
			CheckShowHDRechageBonus()
		end
		return
	end

    local openparam = CS.WindowNodeInitParam("UIBindingGift", CS.AppDefine.HallNumber)
    openparam.WindowCloseCallBack = (function()
		CheckShowHDRechageBonus()
    end)
    CS.WindowManager.Instance:OpenWindow(openparam)
    return
end

function CheckShowHDRechageBonus()
	--G.Tool.TableHelper.Print(user.FlowInfo)
	-- add by rayon
    if (user.FlowInfo.Status == 1 and user.FlowInfo.UnionPlatform == 1 and  user.FlowInfo.BonusRate > 0  ) then
		local openparam = CS.WindowNodeInitParam("UIExchangeRebate", CS.AppDefine.HallNumber)
		CS.WindowManager.Instance:OpenWindow(openparam)
	end
	
	--顯示活動
	if (userdata.ShowCampaignAtLogin == 1) then
		CheckShowActivity (nil, CheckShowFirstRechargeRebate)
	else
		CheckShowFirstRechargeRebate()
	end
end

function CheckShowFirstRechargeRebate()
	-- 顯示首充返利
	if (userdata.ShowFirstDepositAtLogin == 1) then
		ShowFirstRechargeRebateUI()
	end
end

function ShowFirstRechargeRebateUI()
	CS.UILoadingData.Show("载入中,请稍后...", 10,
	(function()
		CS.UIBubblePrompt.Show("取得充值返利资讯失败", groupName)
	end))
	
	G.Service.FirstRechargeRebate.SendRequest(function()
		CS.UILoadingData.Hide()
		local openparam = CS.WindowNodeInitParam("UIActivity2", CS.AppDefine.HallNumber)
		CS.WindowManager.Instance:OpenWindow(openparam)
	end)
end

function CheckShowReward()
    if G.Service.LoginBonus.RetCode == -4 then
        CheckShowLoginGift()
        return false
    elseif G.Service.LoginBonus.RetCode == -3 then
        CheckShowLoginGift()
        return false
    elseif G.Service.LoginBonus.RetCode ~= 2 then
        return false
    end

    local openparam = CS.WindowNodeInitParam("UINotice", CS.AppDefine.HallNumber)
    openparam.WindowCloseCallBack = RefreshUserInfo
    CS.WindowManager.Instance:OpenWindow(openparam)
end

function CheckShowLoginGift()
    local openparam = CS.WindowNodeInitParam("UILoginGift", CS.AppDefine.HallNumber)
    openparam.WindowCloseCallBack = nil
    CS.WindowManager.Instance:OpenWindow(openparam)
end

function CheckDailyGiftInGame()
    G.Service.LoginBonus.SendRequest(UpdateElement)
end

function CheckShowActivity (url, closeCallback) 
    local openparam = CS.WindowNodeInitParam("UIActivity",CS.AppDefine.HallNumber)
    if url ~= nil then openparam.WindowData = url end
	if closeCallback ~= nil then
		openparam.WindowCloseCallBack = closeCallback
	end
    CS.WindowManager.Instance:OpenWindow(openparam)
 end

function CheckShowEvaluate()
    --未储值
    if userdata.SumCharge == 0 then
        return false
    end

    --玩家金额不足
    if userdata.CashAmount < 1 then
        return false
    end

    local play_times = CS.UnityEngine.PlayerPrefs.GetInt("play_times", 0)
    local evaluate_times = CS.UnityEngine.PlayerPrefs.GetInt("evaluate_times", 3 - 20)

    --print("~~~~~~~~~~~~~~~~~~~~~~~~~~play_times",play_times)
    --print("~~~~~~~~~~~~~~~~~~~~~~~~~~evaluate_times",evaluate_times)
    if play_times - evaluate_times < 20 then
        return false
    end

    CS.UnityEngine.PlayerPrefs.SetInt("evaluate_times", play_times)
    local openparam = CS.WindowNodeInitParam("UIEvaluate", CS.AppDefine.HallNumber)
    --openparam.WindowCloseCallBack = callback
    CS.WindowManager.Instance:OpenWindow(openparam)
end

--计算游戏列表在画面显示范围
function CalculateGameScrollView()
    --计算全萤幕跟最小尺寸左边宽度
    local scrollWidth = (RecTran_Main.rect.width - RecTran_Panel.rect.width) * 0.5
    scrollWidth = RecTran_Main.rect.width - scrollWidth
    --扣掉左方显示的子物件
    for i = 0, RecTran_Panel.transform.childCount do
        local child = RecTran_Panel.transform:GetChild(i)
        --子物件不存在
        if child == nil then break end
        --子物件是要计算的物件停止扣除
        if child.gameObject.name == RecTran_GameAttach.gameObject.name then break end
        --子物件有显示才需计算
        if child.gameObject.activeSelf == true then
            child = child.transform:GetComponent("RectTransform")
            if child == nil then break end
            scrollWidth = scrollWidth - child.sizeDelta.x
        end
    end
    scrollWidth = scrollWidth - Layout_Panel.padding.left - Layout_Panel.padding.right
    RecTran_GameAttach.sizeDelta = CS.UnityEngine.Vector2(scrollWidth, RecTran_GameAttach.sizeDelta.y)
end

local CheckServerListTime = 1
local AccCheckServerListTime = 0
function CheckServerList()
    if IsGetServerList == true then
        return
    end

    AccCheckServerListTime = AccCheckServerListTime + CS.UnityEngine.Time.deltaTime
    if AccCheckServerListTime < CheckServerListTime then
        return
    else
        AccCheckServerListTime = 0
    end

    isFisrtRequestServerList = true

    RefreshServerList()
end

function GetUserInfoData()

    if Check_UserData == nil then
        Check_UserData = CheckIsGetData.Create()
        Check_UserData.Init(G.Service.Login.RequestUserInfoHall, "UserData", false)
    else
        Check_UserData.StartGet()
    end

    if Check_ExchangeLog == nil then
        Check_ExchangeLog = CheckIsGetData.Create()
        G.Service.ExchangeLog.ClearLoadLogList()
        Check_ExchangeLog.Init(G.Service.ExchangeLog.InitLoadLogList, "ExchangeLog", false)
    else
        Check_ExchangeLog.StartGet()
    end

    if Check_Mail == nil then
        Check_Mail = CheckIsGetData.Create()
        G.Service.Mail.ClearMailList()
        Check_Mail.Init(G.Service.Mail.InitLoadMailList, "Check_Mail", false)
    else
        Check_Mail.StartGet()
    end

    if Check_CSData == nil then
        Check_CSData = CheckIsGetData.Create()
        G.Service.CS.ClearCSList()
        Check_CSData.Init(G.Service.CS.InitLoadCS, "Check_CSData", false)
    else
        Check_CSData.StartGet()
    end

    if Check_LoginBonus == nil then
        Check_LoginBonus = CheckIsGetData.Create()
        Check_LoginBonus.Init(SendGetLoginBonus, "LoginBonus", false)
    else
        Check_LoginBonus.StartGet()
    end
	
    if Check_VipData == nil then
        Check_VipData = CheckIsGetData.Create()
        Check_VipData.Init(G.Service.Vip.GetVipData, "Vip", false)
    else
        Check_VipData.StartGet()
    end	
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
    -- 依照平台決定背景音樂(特製) by AZLDC >>>>>>>>>>>>>>
    --if math.random(0, 10) > 5 then
    --        G.Tool.Sound.PlayHallBGM("bg_music")
    --else
    --        G.Tool.Sound.PlayHallBGM("bg_music_1")
    --end
    if G.Platform.getJPYL() == false then
        if math.random(0, 10) > 5 then
            G.Tool.Sound.PlayHallBGM("bg_music")
        else
            G.Tool.Sound.PlayHallBGM("bg_music_1")
        end
    else
        if math.random(0, 10) > 5 then
            G.Tool.Sound.PlayHallBGM("bg_music_2")
        else
            G.Tool.Sound.PlayHallBGM("bg_music_3")
        end
    end
    -- 依照平台決定背景音樂(特製) by AZLDC <<<<<<<<<<<<<<
    G.Tool.Sound.MusicStatus(G.GameConfig.MusicStatus)
    G.Tool.Sound.SoundStatus(G.GameConfig.SoundStatus)
    -- 開啟大廳
    if RecTran_Room.gameObject.activeSelf == false and uiRank.GetIsShow() == false then
        RecTran_Main.gameObject:SetActive(true);
    end
    require("Tools/TestConnect").init(false, GameClick, RoomClick)

    GetUserInfoData()
    --重新載入排行榜資料
    uiRank.RefreshShow()

    if windowData ~= nil and windowData.OpenRechargeUI then
        OpenRechargeUI()
    end
end

-- region 刷新用户个人信息
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserHallInfoID, function(receiveMsg)
    NeedCheckUserInfo = false
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetJackpotPoolID, nil)



    if Check_UserData ~= nil then
        Check_UserData.SetGetDataFinish()
    end

    if receiveMsg.Code == 1 then



        G.User.SetCurrentUserInfo(receiveMsg)
        RefreshUserInfo()

        -- 除了第一次進入遊戲，之後回到大廳取得玩家資訊後，再更新一次遊戲列表 by Rayon
        if isFisrtRequestServerList == false then
            RefreshAD()
            RefreshServerList()
            this:CancelInvokeRepeat(RefreshServerList)
            this:InvokeRepeating(RefreshServerList, 60, RefreshServerListTime)

            UpdateButton()
            this:CancelInvokeRepeat(UpdateButton)
            this:InvokeRepeating(UpdateButton, 60, RefresButtonTime)

            CheckShowEvaluate()
            CheckDailyGiftInGame()
        end

		-- 頭相框是0表示是舊機制,關掉取值
		if Check_VipData ~= nil then		
			print("AvatarFrame ", G.User.GetCurrentUser().UserData.AvatarFrame)
			print("UserLevel ", G.User.GetCurrentUser().UserData.UserLevel)
			if G.User.GetCurrentUser().UserData.AvatarFrame == 0 then
				print("###################### Check_VipData SetGetDataFinish ", G.User.GetCurrentUser().UserData.AvatarFrame)
				Check_VipData.SetGetDataFinish()
			end
		end
		
        UpdateButton()
        RefreshDownloadGame()
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end, groupName);

-- endregion

function RefreshUserInfo()
    user = G.User.GetCurrentUser();
    userdata = user.UserData

    local playerID = "ID:" .. user.UserID

    if userdata.IsBindGuest == false then
        playerID = playerID .. "(游客)"
    end
    uiCommon.SetPlayerID(playerID)

    -- 送審階段不開啟 Add By Rayon
    if G.GameConfig.IsSubmissionCtx() == false then
        Btn_Upgrade.gameObject:SetActive(userdata.IsBindGuest == false)
    else
        Btn_Upgrade.gameObject:SetActive(false)
    end
    RefreshBottomAttach()
    --[[
    Btn_HeadIcon:GetComponent("Image"):ResetSpriteByName(CS.AppDefine.HallNumber, tostring(userdata.HeadID))
    Txt_Money.text = G.Tool.UserHelper.ToCash(userdata.CashAmount);
    Txt_RoomMoney.text = G.Tool.UserHelper.ToCash(userdata.CashAmount);
    ]]
    uiCommon.SetHeadIcon(userdata.HeadID)
    uiCommon.SetFrameIcon(userdata.AvatarFrame)
    uiCommon.SetVipLevel(userdata.UserLevel)
    uiCommon.SetMoney(G.Tool.UserHelper.ToCash(userdata.CashAmount))
    uiCommon.SetGiftMoney(G.Tool.UserHelper.ToCash(userdata.CouponAmount))
    uiRank.SetMoney(G.Tool.UserHelper.ToCash(userdata.CashAmount))
end

function RefreshServerList()
    if G.IsGameStart() ~= true then
        -- 熱更期間不重整遊戲介面
        if IsHotUpdate == true then
            return
        end

        if isFisrtRequestServerList then
            CS.UILoadingData.Show("正在请求游戏列表");
        end

        local parm = {}
        parm.Area_ID = G.GameConfig.GameSiteID
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestServerListID, parm)
    end
end

function SetGameIcon(GameObject, Path, IconName)
    local Image = GameObject:Find(Path .. "/ImageTitle_0/name"):GetComponent("Image")
    if Image ~= nil then
        Image:ResetSpriteByName(CS.AppDefine.HallNumber, IconName, true);
    end

	local GiftObject = GameObject:Find(Path .. "/img_gifticon")
	if GiftObject ~= nil then
		local GiftImage = GiftObject:GetComponent("Image")
		if GiftImage ~= nil then
			GiftImage.gameObject:SetActive(G.User.payinfo.UIClub2AgentMode6)	
		end
	end
	

    -- 修改shaher，避免黑邊
    --[[
    local ImageTitleG = GameObject:Find(Path .. "/Ani_Icon")
    if ImageTitleG ~= nil then
        local MaskableGraphic = ImageTitleG:GetComponent("MaskableGraphic")
        if MaskableGraphic ~= nil then
            MaskableGraphic.material = ImgShader.material
        end
    end
    ]] --
    --[[
    local ImageTitleG = GameObject:Find(Path .. "/ImageTitle_0")
    if ImageTitleG.transform.childCount > 0 then
        SkeletonGraphic = ImageTitleG.transform:GetChild(0)
        local MaskableGraphic = SkeletonGraphic:GetComponent("MaskableGraphic")
        if MaskableGraphic ~= nil then
            MaskableGraphic.material = ImgShader.material
        end
    end	
    ]] --
end

-- region 请求服务器列表返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ServerListDataID, function(receiveMsg)

    IsGetServerList = true

    if Txt_LoadingGameTips.gameObject.activeSelf then
        Txt_LoadingGameTips.gameObject:SetActive(false)
    end

    if isFisrtRequestServerList then
        isFisrtRequestServerList = false
        CS.UILoadingData.Hide()
    end

    --G.Tool.TableHelper.Print(receiveMsg)

    -- 熱更期間不重整遊戲介面
    if IsHotUpdate == true then
        return
    end

    gamelist = {};
    autolist = {};
    for i, gameInfo in ipairs(receiveMsg.ServerList) do

        if (not gamelist[gameInfo.GameID]) then
            gamelist[gameInfo.GameID] = {};
        end
        if (gameInfo.Single == 1) then
            autolist[gameInfo.GameID] = {};
            autolist[gameInfo.GameID] = gameInfo;
        end


        gamelist[gameInfo.GameID][gameInfo.Flag] = gameInfo;
    end

    -- 临时变量，用于确定是否需要显示游戏Icon
    for i = 1, #GameLogoList, 1 do
        GameLogoList[i].gameObject:SetActive(false)
        GameLogoList[i]:SetParent(RecTran_GameList)
    end

    RecTran_CarouselRoot.gameObject:SetActive(false)
    RecTran_GameRootTemp.transform.parent:DestoryAllChildren(false);
    RecTran_CarouselRoot.gameObject:SetActive(true)

    -- 修正刷新列表記憶體一直增加的問題
    CS.Utility.CallGC();

    activeGameNumber = 0
    local logoParent;

    local userGameList = nil
    if ((G.User.GetUserGameList2() ~= nil)) then
        userGameList = G.User.GetUserGameList2()

    else
        userGameList = G.User.GetUserGameList()
    end

    -- 只顯示的遊戲
    local onlyShowId = G.Platform.GetOnlyShowGameId()
    local isShowAll = G.Game.IsShowAllGame()

    for _, newgid in pairs(userGameList) do
        local gid = nil
        local tv = nil
        if (tonumber(newgid)) then
            gid = newgid
            tv = 0
        else
            gid = newgid.GameID
            tv = newgid.Popular
            GameT[gid] = tv
        end

        for gameId, subList in pairs(gamelist) do

            -- 依造充值金額開啟，鬥地主預設開啟
            local isShow = false
            if gameId == onlyShowId then
                isShow = true
            else
                isShow = isShowAll
            end

            local GameObject = GameLogoList[gameId]

            if isShow == true and gameId == gid and GameObject ~= nil then


                t1, t2 = math.modf(activeGameNumber / 8);
                -- 8個小遊戲為一單位，取餘數 == 0
                if t2 == 0 then
                    logoParent = CS.UnityEngine.GameObject.Instantiate(RecTran_GameRootTemp.gameObject);
                    logoParent.transform:SetParent(RecTran_GameRootTemp.parent);
                    logoParent.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                    logoParent.transform.localPosition = CS.UnityEngine.Vector3.zero
                    logoParent:SetActive(true)
                end
                activeGameNumber = activeGameNumber + 1

                GameObject:SetParent(logoParent.transform);
                GameObject.gameObject:SetActive(true);
                GameObject.transform.localScale = CS.UnityEngine.Vector3.one

                GameObject:GetComponent("Button").onClick:RemoveAllListeners()



                --游戏入口图标题 设定文字
                local name = "GameIcon_" .. gameId .. "_" .. tostring(newgid.IconName)
                --名称一样不再重新载入
                if iconNameList[gameId] ~= name then
                    iconNameList[gameId] = name

                    SetGameIcon(GameObject, "gameInfo", name)
                    SetGameIcon(GameObject, "Pressed", name)
                end

                local imgt_game = GameObject:Find("gameInfo/t")
                local imgt_press = GameObject:Find("Pressed/t")
                local remoteversion = "0.0.0"

                if (G.Env.IsIOS()) then
                    remoteversion = GetIOSVersion(subList);
                    GameVersion[gid] = GetIOSVersion(subList)
                elseif (G.Env.IsAndroid()) then
                    remoteversion = GetAndroidVersion(subList);
                    GameVersion[gid] = GetAndroidVersion(subList);
                else
                    remoteversion = GetAndroidVersion(subList);
                    GameVersion[gid] = GetAndroidVersion(subList)
                end

                local nversion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gid))

                ResetRibbon(imgt_game, gameId)
                ResetRibbon(imgt_press, gameId)

                GameObject:GetComponent("Button"):AddClickListener(function()
                    local nowversion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gid))

                    --判斷新變數single 0則維持舊機制顯示列表 1則自動入桌,缺產生第一個房間v資料
                    if (not autolist[gameId]) then
                        local name = "GameTitle_" .. gameId
                        -- 设置选厅标题
                        if tonumber(newgid.IconName) ~= 0 then
                            --name = name .. "_" ..G.GameDefine.Version
                            name = name .. "_" .. tonumber(newgid.IconName)
                        end
                        Img_RoomTitle:ResetSpriteByName(CS.AppDefine.HallNumber, name, true);
                        GameClick(gameId, subList);
                    else
                        if (VersionToNum(nowversion, gameId) < VersionToNum(remoteversion, gameId)) then

                            --StartDownloadGame(subList,gameId)
                        else
                            -- 自動入桌
                            RoomClick(autolist[gameId], 1);
                        end
                    end
                end);

                -- 判斷是否下載過
                if RecTran_NoDataGame ~= nil then
                    local IsDownload = CheckGameIsDownLoad(gameId)
                    -- 已經下載過
                    if IsDownload == 0 then
                        -- 需要更新
                        if (VersionToNum(nversion, gameId) < VersionToNum(remoteversion, gameId)) then
                            IsDownload = 2
                        end
                    end

                    -- 0不需要 1第一次下載 2更新
                    if IsDownload > 0 then
                        TempNoDataGame = GameObject.transform:Find("NoDataGame")
                        if TempNoDataGame == nil then
                            TempNoDataGame = CS.UnityEngine.GameObject.Instantiate(RecTran_NoDataGame);
                            TempNoDataGame.transform:SetParent(GameObject.transform);
                            TempNoDataGame.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                            TempNoDataGame.transform.localPosition = CS.UnityEngine.Vector3.zero
                            TempNoDataGame.name = "NoDataGame"
                        end

                        -- 顯示未安裝按鈕
                        InitNoDataButtonIcon(GameObject, TempNoDataGame, IsDownload, gameId)

                        TempNoDataGame:GetComponent("Button").onClick:RemoveAllListeners()
                        TempNoDataGame:GetComponent("Button"):AddClickListener(function()
                            StartDownloadGame(subList, gameId)
                        end);
                    else
                        ShowIcon(gameId)
                    end
                end
            end
        end
    end
    -- todo 
    -- Scrr_GameScrollView:GetComponent("ScrollPage"):SetPageCount(math.ceil(activeGameNumber / 6.0));

    local adl = G.User.GetAdList()
    print("<color=red>ad count:" .. #adl .. "</color>")
    if (AdPageCount <= 1) then
        InitPager(#adl);
        PageSwitch(#adl, 0);
    end
    RefreshAD()
    -- 移除广告上的动画
    RemoveAdAni()
    for i = 1, #adl, 1 do
        if (AdScrr_GameScrollView.transform:Find("GameListContent/_AD" .. tostring(i) .. "/_Content") == nil) then
            local prev_obj = AdScrr_GameScrollView.transform:Find("GameListContent/_AD" .. tostring(i - 1)).gameObject;
            local ad = CS.UnityEngine.GameObject.Instantiate(prev_obj);
            ad.transform:SetParent(prev_obj.transform.parent)
            ad.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
            ad.transform.localPosition = CS.UnityEngine.Vector3.zero
            ad.name = "_AD" .. tostring(i)
        end

        local ImgContent = AdScrr_GameScrollView.transform:Find("GameListContent/_AD" .. tostring(i) .. "/_Content"):GetComponent("Image");
        local ImgRoot = AdScrr_GameScrollView.transform:Find("GameListContent/_AD" .. tostring(i)):GetComponent("Image");
        ImgRoot.gameObject:SetActive(true);
        SetAd(adl[i].Picture, ImgContent, i)
        -- 指定动画
        SetAdAni(adl, i)
    end

    AdScrr_GameScrollView:GetComponent("ScrollPage"):SetPageCount(#adl);

    -- 依造小遊戲多寡調整顯示位置和大小
    UpdateGridLayoutGroup(logoParent, activeGameNumber)

    --CS.UILoading.Hide()

    -- 选厅介面开启时也要同时间刷新房间资讯
    if RecTran_Room.gameObject.activeSelf == true then
        local roomData = nil
        for gameId, subList in pairs(gamelist) do
            if gameId == nowGameId then
                roomData = subList
                break
            end
        end

        -- 有游戏厅资料 刷新
        if roomData ~= nil then
            ShowGameRoom(nowGameId, roomData, true)
            -- 没游戏厅资料 显示大厅
        else
            RecTran_Main.gameObject:SetActive(true);
            RecTran_Room.gameObject:SetActive(false)
        end
    end

    require("Tools/TestConnect").setServerList(gamelist, autolist)
end, groupName);

function RemoveAdAni()
    for i = 0, AdScrr_GameListContent.transform.childCount - 1, 1 do
        local ad = AdScrr_GameListContent.transform:GetChild(i)
        local aniAttach = ad.transform:Find("_Ani")
        if aniAttach ~= nil then
            for j = 0, aniAttach.transform.childCount - 1, 1 do
                local ani = aniAttach.transform:GetChild(j)
                ani.transform:SetParent(RecTran_Carousel.transform)
                ani.gameObject:SetActive(false)
            end
        end
    end
end

function SetAdAni(adl, i)
    local ani = RecTran_Carousel.transform:Find("_aniAd_" .. tostring(adl[i].Behavior))
    if ani == nil then
        return
    end

    local ad = AdScrr_GameScrollView.transform:Find("GameListContent/_AD" .. tostring(i))
    if ad == nil then
        return
    end

    local aniAttach = ad.transform:Find("_Ani")
    if aniAttach == nul then
        return
    end

    ani.transform:SetParent(aniAttach.transform)
    ani.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
    ani.transform.localPosition = CS.UnityEngine.Vector3.zero
    ani.gameObject:SetActive(true)
end

function SetLogo()
    G.Util.coroutine_call(function()
        local url = CS.Utility.GetFullPathOfStreamingAsset("HallLogo.png");

        --url = "http://example.invalid/Upload/AdImage/ad_p7_hall.png"
        local www = CS.UnityEngine.WWW(url)
        yield_return(www);
        if www.error ~= nil then
            print(www.error)
            Img_Logo.gameObject:SetActive(false)
        else
            if www.isDone and www.texture ~= nil and www.texture.width ~= 8 then
                local rect = { x = 1, y = 1, width = www.texture.width, height = www.texture.height }
                Img_Logo.gameObject:SetActive(true)
                Img_Logo.sprite = CS.UnityEngine.Sprite.Create(www.texture, CS.UnityEngine.Rect(0, 0, www.texture.width, www.texture.height), CS.UnityEngine.Vector2(0, 0));
                Img_Logo:SetNativeSize()
                Img_Logo.transform.localScale = CS.UnityEngine.Vector3.one * 0.5
            end
        end
    end)()
end

local Ad_URL = {}
function SetAd(url, Img_Content, index)
	-- 避免重複下載
	if Ad_URL[index] == url then
		return
	end
		
    Img_Content.gameObject:SetActive(false)


    G.Util.coroutine_call(function()
        --url = "http://example.invalid/Upload/AdImage/ad_p7_hall.png"
		local www = CS.UnityEngine.Networking.UnityWebRequest.Get(url)
		local texDl = CS.UnityEngine.Networking.DownloadHandlerTexture(true);		
		www.downloadHandler = texDl;
		yield_return(www:SendWebRequest())
        if www.error ~= nil then
            print(www.error)
        else
            if www.isDone and texDl.texture ~= nil and texDl.texture.width ~= 8 then
                local rect = { x = 1, y = 1, width = texDl.texture.width, height = texDl.texture.height }
                Img_Content.gameObject:SetActive(true)
                Img_Content.sprite = CS.UnityEngine.Sprite.Create(texDl.texture, CS.UnityEngine.Rect(0, 0, texDl.texture.width, texDl.texture.height), CS.UnityEngine.Vector2(0, 0));

				Ad_URL[index] = url
            end
        end
    end)()
end

function InitNoDataButtonIcon(scoureData, cloneData, downloadtype, gameId)

    --scoureData.transform:Find("gameInfo").gameObject:SetActive(false)
    scoureData:GetComponent("Animator").enabled = false
    cloneData.gameObject:SetActive(true)
    cloneData.transform:Find("_Loading"):GetComponent("Image").gameObject:SetActive(false)
    cloneData.transform:Find("_Percentage"):GetComponent("Text").gameObject:SetActive(false)
    cloneData.transform:Find("_WaitText").gameObject:SetActive(false)

    -- SetSprite(scoureData, "Pressed/LoadingIcon"..tostring(G.GameDefine.Version), cloneData, "Icon")
    --SetSprite(scoureData, "Pressed/LoadingIcon"..tostring(G.GameDefine.Version), cloneData, "_Icon")

    -- for i = 0, cloneData.transform.childCount -1, 1 do
    -- 	local child = cloneData.transform:GetChild(i)
    -- 	if child.name == "_Image_" .. tostring(gameId) then
    -- 		child.gameObject:SetActive(true)
    -- 	end
    -- end

    -- -- 1第一次下載
    -- if downloadtype == 1 then
    -- 	cloneData.transform:Find("_Download").gameObject:SetActive(true)
    -- -- 2更新
    -- elseif downloadtype == 2 then
    -- 	cloneData.transform:Find("_Update").gameObject:SetActive(true)	
    -- end
end

function SetSprite(scoureData, scoureName, replaceData, replaceName)
    local _replaceData = replaceData.transform:Find(replaceName):GetComponent("Image")
    if _replaceData == nil then
        return
    end

    local _scoreData = scoureData.transform:Find(scoureName):GetComponent("Image")
    if _scoreData == nil then
        return
    end

    _replaceData.sprite = _scoreData.sprite
end

IsHotUpdate = false
WindowNode_HotUpdateGameObject = nil
Slider_HotUpdate = nil
Img_Icon = nil
Img_Circle = nil
Img_Circle1 = nil
Txt_Persent = nil
DownloadGameList = {}
DownloadSubList = {}
DownloadingGameId = 0

function SetDownloadGame(subList, gameId)
    --CS.UnityEngine.Debug.LogError("--------SetDownloadGame :　" .. tostring(gameId).. "/" .. tostring(G.GetGamePName(gameId)) .. "/" .. tostring(#DownloadGameList))
    local isSave = true
    for i = 1, #DownloadGameList do
        if DownloadGameList[i] == gameId then
            --CS.UnityEngine.Debug.LogError("--------get some gameId :　" .. tostring(gameId).. "/" .. tostring(G.GetGamePName(gameId)) .. "/" .. tostring(i))
            isSave = false
        end
    end

    if isSave == true then
        --CS.UnityEngine.Debug.LogError("--------SetDownloadGame save")
        DownloadGameList[#DownloadGameList + 1] = gameId
        DownloadSubList[#DownloadSubList + 1] = subList
        print("SetDownloadGame", gameId)
    end

    local GameObject = GameLogoList[gameId]
    local NoDataGame = GameObject.transform:Find("NoDataGame")

    if NoDataGame ~= nil then
        if DownloadingGameId ~= gameId then
            local temp_WaitText = NoDataGame.transform:Find("_WaitText")
            temp_WaitText.gameObject:SetActive(true)

            local temp_Circle = NoDataGame.transform:Find("_Loading"):GetComponent("Image")
            temp_Circle.gameObject:SetActive(false)

            local temp_Persent = NoDataGame.transform:Find("_Percentage"):GetComponent("Text")
            temp_Persent.gameObject:SetActive(false)
        end
    end
end

function RefreshDownloadGame()
    local supportVersion = (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail"))

    --CS.UnityEngine.Debug.LogError("--------RefreshDownloadGame")
    for i = 1, #DownloadGameList do
        local gameId = DownloadGameList[i]
        local GameObject = GameLogoList[gameId]
        local NoDataGame = GameObject.transform:Find("NoDataGame")
        if NoDataGame ~= nil then
            if DownloadingGameId ~= gameId then
                local temp_WaitText = NoDataGame.transform:Find("_WaitText")
                temp_WaitText.gameObject:SetActive(supportVersion == true)

                local temp_Circle = NoDataGame.transform:Find("_Loading"):GetComponent("Image")
                temp_Circle.gameObject:SetActive(false)

                local temp_Persent = NoDataGame.transform:Find("_Percentage"):GetComponent("Text")
                temp_Persent.gameObject:SetActive(false)
            end
        end
    end

    if (supportVersion == false) then
        --CS.UnityEngine.Debug.LogError("--------RefreshDownloadGame : supportVersion false")
        if (DownloadingGameId ~= nil and DownloadingGameId ~= 0) then
            --CS.UnityEngine.Debug.LogError("--------RefreshDownloadGame :　--CS.UIUpdateGame.Hide")
            CS.UIUpdateGame.Hide();
        end
        HotUpdateFaild(DownloadingGameId)

        DownloadGameList = {}
        return
    end

    DownloadNextGame()
end

function StartDownloadGame(subList, gameId)
	if subList == nil then
		return
	end
	
    print("StartDownloadGame", gameId, IsHotUpdate)
    --CS.UnityEngine.Debug.LogError("---------------StartDownloadGame :　" .. tostring(gameId).. "/" .. tostring(G.GetGamePName(gameId)))
    --當該遊戲也在下載中時,不繼續加入下載列表
    if (DownloadingGameId == gameId) then
        return
    end

    if IsHotUpdate == true then
        SetDownloadGame(subList, gameId)
        return
    end

    IsHotUpdate = true;
    --設定該ID為正在下載中的遊戲ID
    DownloadingGameId = gameId

    for k, v in pairs(subList) do

        G.Service.NetworkService.SetUpdateGameUiLoadingComplatedCallBack(function(WindowNode)
            BeginHotUpdate(gameId, WindowNode)
        end)

        G.Service.NetworkService.SetUpdateGameComplatedCallBack(function()
            HotUpdateEnd(gameId)
        end)

        G.Service.NetworkService.SetUpdateGameDownloadNextCallBack(function()
            DownloadNextGame()
        end)

        G.Service.NetworkService.SetUpdateGameErrorCallBack(function()
            HotUpdateError(gameId)
        end)

        G.Service.NetworkService.SetUpdateGameFaildCallBack(function()
            HotUpdateFaild(gameId)
        end)

        -- 指定熱更完成是否自動入桌
        --判斷是否自動進入,是則設為true否則設為false---  start
        --[[
        if (not autolist[gameId]) then
            G.Service.NetworkService.SetIsAutoLoginGame(false);
        else
            G.Service.NetworkService.SetIsAutoLoginGame(true);
        end
        ]] --
        -- 下載完不要自動入桌 By Rayon
        G.Service.NetworkService.SetIsAutoLoginGame(false);

        RoomClick(v, 0, false);

        break
    end
end

function BeginHotUpdate(gameId, WindowNode)
    print("BeginHotUpdate", gameId, G.GetGamePName(gameId))
    --CS.UnityEngine.Debug.LogError("--------BeginHotUpdate :　" .. tostring(gameId).. "/" .. tostring(G.GetGamePName(gameId)))

    WindowNode_HotUpdateGameObject = WindowNode.WindowGameObject
    WindowNode_HotUpdateGameObject.transform.localPosition = CS.UnityEngine.Vector3(0, 9999, 0)

    local GameObject = GameLogoList[gameId]
    UpdateTempNoDataGame = GameObject.transform:Find("NoDataGame")
    if UpdateTempNoDataGame ~= nil then
        --Img_Icon = UpdateTempNoDataGame.transform:Find("_Icon"):GetComponent("Image")
        --Img_Icon.gameObject:SetActive(true)

        local temp_WaitText = UpdateTempNoDataGame.transform:Find("_WaitText")
        temp_WaitText.gameObject:SetActive(false)

        Img_Circle = UpdateTempNoDataGame.transform:Find("_Loading"):GetComponent("Image")
        Img_Circle.gameObject:SetActive(true)
        Img_Circle1 = UpdateTempNoDataGame.transform:Find("_Loading/_Loading"):GetComponent("Image")
        Img_Circle1.fillAmount = 1

        Txt_Persent = UpdateTempNoDataGame.transform:Find("_Percentage"):GetComponent("Text")
        Txt_Persent.gameObject:SetActive(true)
        Txt_Persent.text = "0%"
    end
end

function DoHotUpdate()
    if IsHotUpdate == false then
        return
    end

    if WindowNode_HotUpdateGameObject == nil then
        return;
    end

    if Slider_HotUpdate == nil then
        Slider_HotUpdate = WindowNode_HotUpdateGameObject.transform:Find("Slider"):GetComponent("Slider")
    end

    if Slider_HotUpdate == nil then
        return
    end

    if Img_Icon ~= nil then
        Img_Icon.fillAmount = 1 - Slider_HotUpdate.value
    end

    if Img_Circle1 ~= nil then
        Img_Circle1.fillAmount = 1 - Slider_HotUpdate.value
    end

    if Txt_Persent ~= nil then
        Txt_Persent.text = string.format("%.2f", Slider_HotUpdate.value * 100)
    end
end

function HotUpdateEnd(gameId)
    local supportVersion = (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail"))
    if (supportVersion == false) then
        --CS.UnityEngine.Debug.LogError("--------HotUpdateEnd :　--CS.UIUpdateGame.Hide")
        CS.UIUpdateGame.Hide();
    end
    print("HotUpdateEnd", gameId, G.GetGamePName(gameId))
    --CS.UnityEngine.Debug.LogError("--------HotUpdateEnd :　" .. tostring(gameId).. "/" .. tostring(G.GetGamePName(gameId)))
    WindowNode_HotUpdateGameObject = nil
    Slider_HotUpdate = nil

    IsHotUpdate = false
    Img_Icon = nil
    Slider_HotUpdate = nil
    Img_Circle = nil
    Txt_Persent = nil

    DownloadingGameId = 0

    ShowIcon(gameId)

    local openparam = CS.WindowNodeInitParam("UIMessage", CS.AppDefine.HallNumber)
    openparam.NodeType = CS.BaseNodeType.AboveNormal;
    openparam.LoadComplatedCallBack = function(windowNode)
        local txt = windowNode.WindowGameObject.transform:Find("Date_2"):GetComponent("Text")
        print("gameIdgameIdgameIdgameIdgameIdgameId " .. tostring(gameId))
        txt.text = G.Game.GetGameName(tostring(gameId)) .. " 下载完成"
    end
    windowNode = CS.WindowManager.Instance:OpenWindow(openparam)
end

function DownloadNextGame()
    if #DownloadGameList > 0 then
        --CS.UnityEngine.Debug.LogError("--------DownloadNextGame")

        local gameId = DownloadGameList[1]
        local subList = DownloadSubList[1]
        print("DownloadNextGame", #DownloadGameList, DownloadGameList[1])
        table.remove(DownloadGameList, 1)
        table.remove(DownloadSubList, 1)
        StartDownloadGame(subList, gameId)
    else
        --CS.UnityEngine.Debug.LogError("--------DownloadNextGame list is null")
    end
end

function HotUpdateError()
    local supportVersion = (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail"))
    if (supportVersion == false) then
        --CS.UnityEngine.Debug.LogError("--------HotUpdateError :　--CS.UIUpdateGame.Hide")
        CS.UIUpdateGame.Hide();
    end

    print("HotUpdateError", DownloadingGameId, G.GetGamePName(DownloadingGameId))
    --CS.UnityEngine.Debug.LogError("--------HotUpdateError :　" .. tostring(DownloadingGameId).. "/" .. tostring(G.GetGamePName(DownloadingGameId)))
    local GameObject = GameLogoList[DownloadingGameId]
    UpdateTempNoDataGame = GameObject.transform:Find("NoDataGame")
    if UpdateTempNoDataGame ~= nil then
        local temp_WaitText = UpdateTempNoDataGame.transform:Find("_WaitText")
        temp_WaitText.gameObject:SetActive(false)
    end

    DownloadingGameId = 0
end

function HotUpdateFaild()
    print("HotUpdateFaild", DownloadingGameId, G.GetGamePName(DownloadingGameId))
    --CS.UnityEngine.Debug.LogError("--------HotUpdateError :　" .. tostring(DownloadingGameId).. "/" .. tostring(G.GetGamePName(DownloadingGameId)))
    if (DownloadingGameId == nil) then
        return
    elseif (GameLogoList[DownloadingGameId] == nil) then
        return
    end

    local NowDownloadingGameId = DownloadingGameId
    local openparam = CS.WindowNodeInitParam("UIMessage", CS.AppDefine.HallNumber)
    openparam.NodeType = CS.BaseNodeType.AboveNormal;
    openparam.LoadComplatedCallBack = function(windowNode)
        local txt = windowNode.WindowGameObject.transform:Find("Date_2"):GetComponent("Text")
        txt.text = G.Game.GetGameName(tostring(NowDownloadingGameId)) .. " 下载失败"
    end
    windowNode = CS.WindowManager.Instance:OpenWindow(openparam)

    local GameObject = GameLogoList[DownloadingGameId]
    UpdateTempNoDataGame = GameObject.transform:Find("NoDataGame")
    if UpdateTempNoDataGame ~= nil then
        local temp_WaitText = UpdateTempNoDataGame.transform:Find("_WaitText")
        temp_WaitText.gameObject:SetActive(false)

        local temp_Circle = UpdateTempNoDataGame.transform:Find("_Loading"):GetComponent("Image")
        temp_Circle.gameObject:SetActive(false)

        local temp_Persent = UpdateTempNoDataGame.transform:Find("_Percentage"):GetComponent("Text")
        temp_Persent.gameObject:SetActive(false)
    end

    WindowNode_HotUpdateGameObject = nil

    IsHotUpdate = false
    Img_Icon = nil
    Slider_HotUpdate = nil
    Img_Circle = nil
    Txt_Persent = nil

    DownloadingGameId = 0
end

function ShowIcon(gameId)
    local GameObject = GameLogoList[gameId]
    if GameObject ~= nil then
        GameObject:GetComponent("Animator").enabled = true

        TempNoDataGame = GameObject.transform:Find("NoDataGame")
        if TempNoDataGame ~= nil then
            TempNoDataGame.gameObject:SetActive(false)
        end

        GameInfo = GameObject.transform:Find("gameInfo")
        GamePress = GameObject.transform:Find("Pressed")

        if GameInfo ~= nil then
            GameInfo.gameObject:SetActive(true)
            local imgt = GameInfo:Find("t")
            ResetRibbon(imgt, gameId)
        end
        if GamePress ~= nil then
            local imgt = GamePress:Find("t")
            ResetRibbon(imgt, gameId)
        end
    end
end

function Destroy(Go)
    CS.UnityEngine.GameObject.Destroy(Go.gameObject)
end

function DestroyChild(Go)
    for i = 0, Go.transform.childCount - 1 do
        node = Go.transform:GetChild(i)
        if (node ~= nil) then Destroy(node) end
    end
end

--設定遊戲彩帶狀態
function ResetRibbon(imgt, gameId)
    --當下遊戲版本號
    local nversion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId))
    --遠端版本號
    local rversion = GameVersion[gameId]
    --當版本號不足時
    if (VersionToNum(nversion, gameId) < VersionToNum(rversion, gameId)) then
        imgt.gameObject:SetActive(true);
        local num = imgt.transform.childCount
        if (num > 0) then
            DestroyChild(imgt)
        end

        local RibbonItem = CS.UnityEngine.GameObject.Instantiate(recTran__flagRoot.transform:Find("t6"));
        RibbonItem.transform:SetParent(imgt.transform)
        RibbonItem.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
        RibbonItem.transform.localPosition = CS.UnityEngine.Vector3.zero

        -- imgt:ResetSpriteByName(CS.AppDefine.HallNumber, "t_6", true);
    elseif (GameT[gameId] == 0) then
        imgt.gameObject:SetActive(false);

    else

        imgt.gameObject:SetActive(true);
        local num = imgt.transform.childCount
        if (num > 0) then
            DestroyChild(imgt)
        end

        local RibbonItem = CS.UnityEngine.GameObject.Instantiate(recTran__flagRoot.transform:Find("t" .. GameT[gameId]));
        RibbonItem.transform:SetParent(imgt.transform)
        RibbonItem.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
        RibbonItem.transform.localPosition = CS.UnityEngine.Vector3.zero

        --imgt:ResetSpriteByName(CS.AppDefine.HallNumber, "t_"..GameT[gameId], true);      
    end
end

function InitUpdateGridLayoutGroup(logoParent)
    if logoParent == nil then
        return
    end
    local _grid = logoParent:GetComponent("GridLayoutGroup");
    if _grid ~= nil then
        _grid.padding.left = GridLayoutGroup_GameRootTemp.padding.left
        _grid.padding.right = GridLayoutGroup_GameRootTemp.padding.right
        _grid.example.invalid = GridLayoutGroup_GameRootTemp.example.invalid
        _grid.padding.bottom = GridLayoutGroup_GameRootTemp.padding.bottom
        _grid.spacing = GridLayoutGroup_GameRootTemp.spacing
        _grid.constraintCount = GridLayoutGroup_GameRootTemp.constraintCount
    end
end

function UpdateGridLayoutGroup(logoParent, activeGameNumber)
    if G.GameConfig.IsSubmissionCtx() == false then
        --return
    end

    if G.Game.IsShowAllGame() == true then
        InitUpdateGridLayoutGroup(logoParent)
        return
    end

    if (activeGameNumber == 1) then
        SetGridLayoutGroup(logoParent, 300, 200, 0, 0, 1.5, 4)
    elseif (activeGameNumber == 2) then
        SetGridLayoutGroup(logoParent, 125, 100, 100, 0, 1.2, 4)
    elseif (activeGameNumber == 3) then
        SetGridLayoutGroup(logoParent, 50, 100, 25, 0, 1, 4)
    elseif (activeGameNumber == 4) then
        SetGridLayoutGroup(logoParent, 100, 0, 100, 0, 1.2, 2)
    end
end

function SetGridLayoutGroup(_parent, _left, _top, _spacingX, _spacingY, _childSize, _constraintCount)
    --print("SetGridLayoutGroup",_left,_top,_spacingX,_spacingY,_childSize)
    if _parent == nil then
        return
    end

    local _grid = _parent:GetComponent("GridLayoutGroup");
    if _grid ~= nil then
        _grid.padding.left = _left
        _grid.padding.right = 0
        _grid.example.invalid = _top
        _grid.padding.bottom = 0
        _grid.spacing = CS.UnityEngine.Vector2(_spacingX, _spacingY)
        _grid.constraintCount = _constraintCount
    end

    for i = 0, _parent.transform.childCount - 1, 1 do
        _parent.transform:GetChild(i).localScale = CS.UnityEngine.Vector3.one * _childSize
    end
end

function VersionToNum(v, gameId)
    if CS.AppDefine.OpenHotUpdate == false then
        return 0
    end

    if v == nil or v == "" then
        print("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! parse VersionToNum error " .. gameId)
        return 100000
    end

    local strarr = G.Tool.StringHelper.Split(v, ".")

    return (tonumber(strarr[1]) * 100000) + (tonumber(strarr[2]) * 1000) + (tonumber(strarr[3]) * 1)
end


function GetAndroidVersion(subList)

    for k, v in pairs(subList) do
        return v.AndroidVersion
    end
end

function GetIOSVersion(subList)

    for k, v in pairs(subList) do
        return v.IOSVersion
    end
end

function GameClick(gameId, subList)
    local supportVersion = (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail"))
    if IsHotUpdate == true and supportVersion == false then
        CS.UIBubblePrompt.Show("您正在下载游戏，请稍后", groupName);
        return
    end

    G.Platform.ShowUpdateMessage(RefreshServerList)
    ShowGameRoom(gameId, subList, false)

    RecTran_Main.gameObject:SetActive(false);
    RecTran_Room.gameObject:SetActive(true);
    RecTran_EF.gameObject:SetActive(false);
    AnimRoom(gameId)
end

function ShowGameRoom(gameId, subList, interactable)
    nowGameId = gameId

    subList = G.Tool.TableHelper.Sort(subList, function(a, b)
        return a.Flag < b.Flag
    end)

    --[[捕鱼特制，預載比賽規則
    if gameId == 14 then
        SendGetCompetioniInfoRule(gameId)
    end]]
    -- 临时变量，用于确定是否需要显示房间Icon
    local isShowIcon = false
    --[[
    --捕鱼特制，更换背景与关闭角色图
    RecTran_Background14.gameObject:SetActive(gameId == 14)
    RecTran_Background1.gameObject:SetActive(gameId == 1)
    ]]
    uiCommon.ChangeBackground(gameId)

    --RecTran_RoomList.transform.position= CS.UnityEngine.Vector3(RecTran_RoomList.transform.position.x,((gameId == 14 or gameId == 24 or gameId == 1) and -90 or (0)   ),RecTran_RoomList.transform.position.z)  
    RecTran_RoomList.transform.position = CS.UnityEngine.Vector3(RecTran_RoomList.transform.position.x, -90, RecTran_RoomList.transform.position.z)
    for index, roomdata in pairs(RoomList) do
        if roomdata.base ~= nil then
            roomdata.base:SetParent(RecTran_Rooms)
            roomdata.base.gameObject:SetActive(false)
        end
        if roomdata.by ~= nil then
            roomdata.by:SetParent(RecTran_RoomsBy)
            roomdata.by.gameObject:SetActive(false)
        end
        if roomdata.ddz ~= nil then
            roomdata.ddz:SetParent(RecTran_RoomsDdz)
            roomdata.ddz.gameObject:SetActive(false)
        end
    end

    --G.Tool.TableHelper.Print("@@@@@@@@@@@@@@@@@@@", subList)

    local Oddsetting = false;
    Txt_MaxOdds.gameObject:SetActive(false);

    for index, roomdata in pairs(RoomList) do
        local gameRoomIns = roomdata.base

        if gameId == 14 and roomdata.by ~= nil then
            gameRoomIns = roomdata.by
        elseif gameId == 24 and roomdata.by ~= nil then
            gameRoomIns = roomdata.by
        elseif gameId == 1 and roomdata.ddz ~= nil then
            gameRoomIns = roomdata.ddz
        end

        local localVersion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId));
        Txt_RoomVersion.text = "V." .. tostring(localVersion)

        for k, v in pairs(subList) do
            if v.Flag == tostring(index) then
                if ((v.MaxOdd ~= nil) and (v.MaxOdd ~= 0) and (Oddsetting ~= true)) then
                    Oddsetting = true;
                    Txt_MaxOdds.text = string.format("所有场最大倍率%s倍", v.MaxOdd)
                    Txt_MaxOdds.gameObject:SetActive(true);
                end

                if (v.GameID == 1) then
                    gameRoomIns:SetParent(G.CaseWhen(tonumber(v.Flag) <= 6, RecTran_RoomRow1.gameObject.transform, RecTran_RoomRow2.gameObject.transform))
                else
                    gameRoomIns:SetParent(G.CaseWhen(k <= #RoomList, RecTran_RoomRow1.gameObject.transform, RecTran_RoomRow2.gameObject.transform))
                end
                gameRoomIns.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                gameRoomIns.localPosition = CS.UnityEngine.Vector3(0, 0, 0)
                gameRoomIns.gameObject:SetActive(true)

                --体验场要放第一个
                if index == 6 then
                    gameRoomIns.transform:SetAsFirstSibling()
                end
                --在動畫演示前先將按鈕觸控
                gameRoomIns:GetComponent("Button").interactable = interactable;

                gameRoomIns:Find("UIBody/Info1"):GetComponent("Text").text = string.format(G.Tool.UserHelper.ToCashUnit(v.LoginMoney) .. "入场");

                if v.GameID == 14 or v.GameID == 24 then
                    if index == 6 then
                        gameRoomIns:Find("UIBody/Info2"):GetComponent("Text").text = "输赢不计"
                        gameRoomIns:Find("UIBody/Info1"):GetComponent("Text").text = "免费体验"
                    else
                        gameRoomIns:Find("UIBody/Info2"):GetComponent("Text").text = string.format("倍率：" .. G.Tool.UserHelper.ToCashUnit(v.BaseScore));
                    end
                elseif v.GameID == 1 then
                    if index == 6 then
                        gameRoomIns:Find("UIBody/Info2"):GetComponent("Text").text = string.format("底分：" .. G.Tool.UserHelper.ToCashUnit(v.BaseScore));
                        gameRoomIns:Find("UIBody/Info1"):GetComponent("Text").text = "免费体验"
                    else
                        gameRoomIns:Find("UIBody/Info2"):GetComponent("Text").text = string.format("底分：" .. G.Tool.UserHelper.ToCashUnit(v.BaseScore));
                    end
                else
                    gameRoomIns:Find("UIBody/Info2"):GetComponent("Text").text = string.format("底分：" .. G.Tool.UserHelper.ToCashUnit(v.BaseScore));
                end

                local RankHelp = gameRoomIns:Find("UIBody/_RankHelp")
                if RankHelp ~= nil then
                    RankHelp:GetComponent("Button").onClick:RemoveAllListeners()
                    RankHelp:GetComponent("Button"):AddClickListener(function()
                        OpenRankHelpUI(v.GameID)
                    end);
                end

                if v.GameType == 1 then
                    SendGetCompetioniInfoRule(v.GameID)
                end

                gameRoomIns:GetComponent("Button").onClick:RemoveAllListeners()
                gameRoomIns:GetComponent("Button"):AddClickListener(function()
                    if v.GameType == 1 then
                        uiRank.Show(v.GameID, v.Flag)
                    else
                        RoomClick(v, 1);
                    end

                    uiRank.SetCallbackStart(function()
                        RoomClick(v, 1);
                    end)
                end);

                uiRank.SetCallbackRankHelp(function()
                    OpenRankHelpUI(v.GameID)
                end)

                if v.GameID == 1 then
                    gameRoomIns:Find("UIBody/Icon").gameObject:SetActive(true)

                else
                    gameRoomIns:Find("UIBody/Icon").gameObject:SetActive(false)
                end

                --local roomEf = gameRoomIns.transform:Find("UIBody/RoomEF").gameObject
                --CS.Utility.SetUISortingOrder(roomEf)
            end
        end
    end
end

function AnimRoom(gameId)
    local sequence = CS.DG.Tweening.DOTween.Sequence()

    local timer = 0.3
    --[[
    local target = Btn_RoomlistBack.transform.position
    Btn_RoomlistBack.transform.position = CS.UnityEngine.Vector3(-1500,target.y,target.z)
    sequence:Append(Btn_RoomlistBack.transform:DOMove(target, timer))
    
    timer = 0.25
    target = RecTran_RoomMoney.transform.position
    RecTran_RoomMoney.transform.position = CS.UnityEngine.Vector3(-1000,target.y,target.z)
    sequence:Join(RecTran_RoomMoney.transform:DOMove(target, timer))
]]
    timer = 0.25
    target = RecTran_RoomTitle.transform.position
    RecTran_RoomTitle.transform.position = CS.UnityEngine.Vector3(1000, target.y, target.z)
    sequence:Join(RecTran_RoomTitle.transform:DOMove(target, timer))
    --[[
    timer = 0.25
    target = Btn_RoomRecharge.transform.position
    Btn_RoomRecharge.transform.position = CS.UnityEngine.Vector3(target.x,-500,target.z)
    sequence:Join(Btn_RoomRecharge.transform:DOMove(target, timer))
       ]]

    timer = 0.25
    target = RecTran_RoomList.transform.position
    if (gameId == 1) then
        target = CS.UnityEngine.Vector3(target.x, -90, target.z)
    elseif (gameId == 14 or gameId == 24) then
        target = CS.UnityEngine.Vector3(target.x, -90, target.z)
    end
    RecTran_RoomList.transform.position = CS.UnityEngine.Vector3(1000, target.y, target.z)
    sequence:Join(RecTran_RoomList.transform:DOMove(target, timer))

    if (gameId == 14 or gameId == 24) then

        math.randomseed(4)


        for i = 1, #RoomList do


            if (RoomList[i].by.transform:Find("UIBody") ~= nil) then

                RoomList[i].by.transform:Find("UIBody").transform.localScale = CS.UnityEngine.Vector3.zero
            end
        end
        rnd = G.Tool.Math.MakeRandom(1, RecTran_RoomRow1.childCount);

        --動畫播放
        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                G.Tool.Ani.BubbleScaleS2B1(RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)).transform:Find("UIBody"), 1.3, 0.1, nil)
            end)
            sequence:AppendInterval(0.1)
        end
        --啟用按鈕觸動
        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)):GetComponent("Button").interactable = true;
            end)
            sequence:AppendInterval(0.1)
        end



    elseif (gameId == 1) then
        for i = 1, #RoomList do
            if (RoomList[i].ddz.transform:Find("UIBody") ~= nil) then
                RoomList[i].ddz.transform:Find("UIBody").transform.localScale = CS.UnityEngine.Vector3.zero
            end
        end
        --提供不重複亂數排序
        rnd = G.Tool.Math.MakeRandom(1, RecTran_RoomRow1.childCount);
        rnd1 = G.Tool.Math.MakeRandom(1, RecTran_RoomRow2.childCount);


        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                G.Tool.Ani.LandingScale(RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)).transform:Find("UIBody"), 3, 0.15, nil)
            end)
            sequence:AppendInterval(0.1)
        end

        for i = 0, RecTran_RoomRow2.childCount - 1, 1 do
            sequence:AppendCallback(function()
                G.Tool.Ani.LandingScale(RecTran_RoomRow2:GetChild((rnd1[i + 1] - 1)).transform:Find("UIBody"), 3, 0.15, nil)
            end)
            sequence:AppendInterval(0.1)
        end

        --啟用按鈕觸動
        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)):GetComponent("Button").interactable = true;
            end)
            sequence:AppendInterval(0.1)
        end
        --啟用按鈕觸動
        for i = 0, RecTran_RoomRow2.childCount - 1, 1 do
            sequence:AppendCallback(function()
                RecTran_RoomRow2:GetChild((rnd[i + 1] - 1)):GetComponent("Button").interactable = true;
            end)
            sequence:AppendInterval(0.1)
        end

    else
        for i = 1, #RoomList do
            if (RoomList[i].base.transform:Find("UIBody") ~= nil) then
                RoomList[i].base.transform:Find("UIBody").transform.localScale = CS.UnityEngine.Vector3.zero
            end
        end
        --提供不重複亂數排序
        rnd = G.Tool.Math.MakeRandom(1, RecTran_RoomRow1.childCount);

        for i = 1, #rnd, 1 do
            print("<color=yellow>" .. rnd[i] .. "</color>");
        end

        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                G.Tool.Ani.LandingScale(RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)).transform:Find("UIBody"), 3, 0.15, nil)
            end)
            sequence:AppendInterval(0.1)
        end
        --啟用按鈕觸動
        for i = 0, RecTran_RoomRow1.childCount - 1, 1 do
            sequence:AppendCallback(function()
                RecTran_RoomRow1:GetChild((rnd[i + 1] - 1)):GetComponent("Button").interactable = true;
            end)
            sequence:AppendInterval(0.1)
        end
    end
    --[[
    timer = 1
    target = ani_kani.transform.position
    sequence:Join(ani_kani.transform:DOMove(CS.UnityEngine.Vector3(307,target.y,target.z)   , timer))
    sequence:Append(ani_kani.transform:DOMove(CS.UnityEngine.Vector3(-207,target.y,target.z), timer))  
    ]]
    --timer = 0.5
    --target = RecTran_RoomPerson.transform.position
    --RecTran_RoomPerson.transform.position = CS.UnityEngine.Vector3(1000,target.y,target.z)
    --sequence:Join(RecTran_RoomPerson.transform:DOMove(target, timer))    
    --[[
    if(gameId==1)then
     ani_shark.AnimationState:SetAnimation(0, "start", false);
     ani_shark.AnimationState:AddAnimation(0, "idle", true , 2);
     ani_kani.AnimationState:SetAnimation(0, "walk", true);
     ani_kani.AnimationState:AddAnimation(0, "idle", true , 2);
    end
	]]
end

-- EnterRoom
function RoomClick(serverinfo, gohotupdate, formHotUpdate)
    local supportVersion = (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail"))
    if IsHotUpdate == true and formHotUpdate == nil and supportVersion == false then
        CS.UIBubblePrompt.Show("您正在下载游戏，请稍后", groupName);
        return
    end

    G.Platform.ShowUpdateMessage(RefreshServerList)

    if ((userdata.CashAmount + userdata.BankAmount) > serverinfo.KickOutMoney and serverinfo.KickOutMoney ~= 0 and serverinfo.KickOutMoney ~= nil and CheckGameIsDownLoad(serverinfo.GameID) and gohotupdate ~= 0) then
        CS.UIBubblePrompt.Show("您的金币超出房间限制" .. G.Tool.UserHelper.ToCashDetail(serverinfo.KickOutMoney) .. "!请前往菁英场进行游戏", groupName);
        return
    end


    RequestGameID = serverinfo.GameID

    local userGameList = nil;
    if ((G.User.GetUserGameList2() ~= nil)) then
        userGameList = G.User.GetUserGameList2()
    else
        userGameList = G.User.GetUserGameList()
    end
    local skin = G.Tool.TableHelper.Find(userGameList, function(v) return v.GameID == RequestGameID end)[1].Skin
    -- 强制解锁
    G.NetQ.Resume()
    local requestGameServerAddrInfo = {};
    requestGameServerAddrInfo.ServerFlag = serverinfo.Flag;
    requestGameServerAddrInfo.GameID = serverinfo.GameID;
    requestGameServerAddrInfo.IP = serverinfo.IPAddress
    requestGameServerAddrInfo.Port = serverinfo.Port

    G.GameConfig.RequestGameRemoteVersion = GameVersion[serverinfo.GameID]
    G.GameConfig.RequestGameServerFlag = serverinfo.Flag
    if (formHotUpdate == nil) then
        G.Service.NetworkService.RequestGameGameServerAddr(requestGameServerAddrInfo, skin)
    else
        G.Service.NetworkService.RequestHotUpdate(serverinfo, skin)
    end
end

-- endregion

-- region 登录房间消息处理
G.NetQ.RegMsgHook(G.Proto.BClass.GameRoom, G.Proto.GameRoom.LoginRoomRetID, function(receiveMsg)
    CS.UILoadingData.Hide()
    --Login Room 重整
    --RefreshServerList()

    local RequestGameID = receiveMsg.RoomInfo.GameID
    local RequestFlag = receiveMsg.RoomInfo.Flag
    local obj = GameLogoList[RequestGameID]
    local imgt_game = obj:Find("gameInfo/t")
    local imgt_press = obj:Find("Pressed/t")
    ResetRibbon(imgt_game, RequestGameID)
    ResetRibbon(imgt_press, RequestGameID)




    if receiveMsg.Code ~= 1 then
        if (G.IsGameStart() == true) then
            G.Service.Game.ExitGame()
        end


        -- 关闭连接
        if (receiveMsg.Code == -2) then
            G.Util.coroutine_call(function()
                -- 充值開關 Add By Rayon
                if G.GameConfig.IsRechargeCtx() == false then
                    showMessageBox(CS.MessageBoxStyle.OKCancel, "", receiveMsg.Message, "继续", "取消");
                else
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", receiveMsg.Message, "前往充值", "取消");
                    if boxResult == CS.MessageBoxResult.OK then
                        OpenRechargeUI()
                    end
                end
            end)()
        elseif (receiveMsg.Code == -3) then
            G.Util.coroutine_call(function()
                local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", receiveMsg.Message, "继续", "取消");
                if boxResult == CS.MessageBoxResult.OK then
                    local serverParm = G.Tool.StringHelper.Split(receiveMsg.UserData.ServerAddr, "|");
                    RequestGameID = tonumber(serverParm[1])
                    local addressParm = G.Tool.StringHelper.Split(serverParm[2], ":");
                    G.Server.GetGameServerAddrCallback(addressParm[1], addressParm[2], function(res, ipAddress, port)
                        if res == 0 then
                            G.SetRunGameIPAndPort(ipAddress, port)

                            --寫入SKIN
                            local userGameList = nil;
                            if ((G.User.GetUserGameList2() ~= nil)) then
                                userGameList = G.User.GetUserGameList2()
                            else
                                userGameList = G.User.GetUserGameList()
                            end
                            local skin = G.Tool.TableHelper.Find(userGameList, function(v) return v.GameID == RequestGameID end)[1].Skin

                            G.Service.NetworkService.CheckGameVersion(RequestGameID, skin)
                        else
                            CS.UIBubblePrompt.Show("游戏连接失败(1)", "UIHall")
                            require("Tools/TestConnect").loginError("遊戲盾取得遊戲IP失敗")
                        end
                    end);
                end
            end)()
        else
            -- 提示失败原因
            CS.UIBubblePrompt.Show(receiveMsg.Message, "UIHall")
        end
    else

        RecTran_Main.gameObject:SetActive(false);

        print("进入房间成功")
        G.NetQ.Pause()
        G.Room.SetRoomInfo(receiveMsg.RoomInfo)
        G.Room.SetUserData(receiveMsg.UserData)


        if G.IsGameStart() == true then
            CS.FrameWork.ObjectPoolManager.Instance:Retrieve();
            print("reload游戏场景")
            G.Service.Game.ReloadGame(function()
                G.Room.GameLoginSuccess(true)
                CS.GameManager.Instance:StartGame()
            end)
        else
            print("关闭大厅界面")
            G.SetRunGame(RequestGameID)
            G.SetRunGameData(gamelist[RequestGameID][RequestFlag])
            G.Tool.Sound.CloseBGMSound()
            CS.WindowManager.Instance:CloseWindow("UIHall", true);
            CS.Utility.CallGC();
            G.Room.GameLoginSuccess(true)
            CS.GameManager.Instance:StartGame();
        end

        print("当前加载场景数量 " .. CS.UnityEngine.SceneManagement.SceneManager.sceneCount)
    end
end, groupName);
-- endregion

-- region 玩家属性变化
G.NetQ.RegMsgHook(G.Proto.BClass.NotifyServer, G.Proto.NotifyServer.AttrChangeID, function(receiveMsg)
    -- 循环变化属性，做相应操作
    for i, Attr in ipairs(receiveMsg.Attrs) do
        -- 0 金币变化 1 新邮件
        if Attr.TypeID == 0 then
            if this.gameObject.activeSelf == true then
                -- 金額變多，有可能是充值，刷新介面
                if Attr.Param > userdata.CashAmount then
                    G.Service.Login.RequestUserInfoHall()
                else
                    userdata.CashAmount = Attr.Param
                    RefreshUserInfo()
                end

                -- 在充值介面錢變動時刷新介面
                local node = CS.WindowManager.Instance:FindWindowNodeByName("UIRecharge")
                if node ~= nil then
                    local Txt_CashAmount = node.WindowGameObject.transform:Find("UIBody/_OnlinePayDetail/_CashAmount"):GetComponent("Text")
                    if Txt_CashAmount ~= nil then
                        Txt_CashAmount.text = G.Tool.UserHelper.ToCash(Attr.Param)
                    end
                end
            else
                print("@@@@@@@@@@@@ 大廳介面不顯示,不刷新大廳資訊")
            end
        elseif Attr.TypeID == 1 then
            G.Tool.Sound.PlayHallSound("NewMail")
            G.Service.Mail.RefreshMailList();
        elseif Attr.TypeID == 3 then
            G.User.SetPayAmountConfig(Attr.SParam)
        elseif Attr.TypeID == 4 then
            G.Tool.Sound.PlayHallSound("NewMail")
            G.Service.CS.RefreshCSList();
        end
    end
end, groupName);

-- region 彩池接收
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetJackpotPoolRetID, function(receiveMsg)
    local randomnum = math.random(-500000, 500000)

    moneypool = receiveMsg.amount + randomnum
    print("彩池接收", moneypool, randomnum)
    if (seqJackpotPoolRet ~= nil) then
        seqJackpotPoolRet:Kill()
    end
    seqJackpotPoolRet = CS.DG.Tweening.DOTween.Sequence()
    seqJackpotPoolRet:SetLoops(-1)

    seqJackpotPoolRet:AppendCallback(function()
        moneypool = moneypool + 10
        txt__Money.text = tostring(G.Tool.UserHelper.ToCashDetail(moneypool))
    end)
end, groupName);
-- endregion

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    if CS.SoundManager.Instance then
        G.Tool.Sound.CloseBGMSound()
    end

    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)

    -- 移除邮件列表消息接拦截
    G.Service.Mail.RemoveRegMailMsg()

    G.Service.CS.RemoveRegCSMsg()
    --todo
    --local scrollPage = Scrr_GameScrollView:GetComponent("ScrollPage")
    --scrollPage.OnPageChanged = scrollPage.OnPageChanged - onGamePageChanged
    local AdscrollPage = AdScrr_GameScrollView:GetComponent("ScrollPage")
    AdscrollPage.OnPageChanged = AdscrollPage.OnPageChanged - onAdPageChanged

    if (seqJackpotPoolRet ~= nil) then
        seqJackpotPoolRet:Kill()
    end
end

function EscapeWindow()

    if uiRank.GetIsShow() == true then
        uiRank.Hide();
        RecTran_Room.gameObject:SetActive(true);
    elseif RecTran_Room.gameObject.activeSelf then
        -- 删除所有显示出来的子节点
        RecTran_Main.gameObject:SetActive(true);
        RecTran_Room.gameObject:SetActive(false);
        RecTran_EF.gameObject:SetActive(true);
        uiCommon.ChangeBackground()

        G.Service.Game.CloseScene()
    else
        G.Util.coroutine_call(function()
            local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancelClose, "", "您确定要确认退出游戏吗？", "再玩一下", "残忍退出");

            if boxResult == CS.MessageBoxResult.Cancel then
                G.Service.Login.Logout()
                CS.UnityEngine.Application.Quit();
            end
        end)()
    end
end

function CheckGameIsDownLoad(gameId)
    if CS.AppDefine.OpenHotUpdate == false then
        return 0
    end

    local version = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId));
    if version == "0.0.0" then
        return 1
    else
        return 0
    end
end

function tablelength(T)
    local count = 0
    for _ in pairs(T) do count = count + 1 end
    return count
end

function UpdateElement()
    Img_Mail1.gameObject:SetActive(G.Service.Mail.GetUnReadMailCount() == 0)
    Img_Mail2.gameObject:SetActive(G.Service.Mail.GetUnReadMailCount() > 0)
    Txt_Mail.text = G.Service.Mail.GetUnReadMailCount()


    Img_Notice1.gameObject:SetActive(G.Service.LoginBonus.RetCode ~= 2)
    Img_Notice2.gameObject:SetActive(G.Service.LoginBonus.RetCode == 2)

    Img_Coffer1.gameObject:SetActive(userdata.SumCharge == 0)
    Img_Coffer2.gameObject:SetActive(userdata.SumCharge > 0)
end

function UpdateCS()
    Img_CS1.gameObject:SetActive(G.Service.CS.GetUnReadCSCount() == 0)
    Img_CS2.gameObject:SetActive(G.Service.CS.GetUnReadCSCount() > 0)
    Txt_CS.text = G.Service.CS.GetUnReadCSCount()
end

function RefreshBottomAttach()
    local _count = 0
    for i = 0, Root_BottomAttach.transform.childCount - 1 do
        local obj = Root_BottomAttach.transform:GetChild(i)
        if (obj.gameObject.activeSelf == true) then
            _count = _count + 1
        end
    end

    local _x = (_count - 1) * 123.5 + 103 + 219;
    local _y = Img_BottomAttach.rectTransform.sizeDelta.y

    Img_BottomAttach.rectTransform.sizeDelta = CS.UnityEngine.Vector2(_x, _y)
end

function UpdateButton()
    Btn_CustomService.gameObject:SetActive(IsShowCustomServiceButton())
	--Btn_Exchange.gameObject:SetActive(IsShowExchangeButton())
		
	local IsShow = IsShowExchangeButton()
	-- 还没到达显示条件
	if IsShow == false then
		-- 是否强制显示
		if userdata.ShowConvertPopup == 1 then
			Btn_Exchange.gameObject:SetActive(true)
		else
			Btn_Exchange.gameObject:SetActive(false)
		end
	else
		Btn_Exchange.gameObject:SetActive(IsShow)
	end
	
    local btnex = nil;
    if (tonumber(G.GameDefine.ExchangeFontType) == 0) then
        btnex = Btn_ExchangeFont1
    elseif (tonumber(G.GameDefine.ExchangeFontType) == 1) then
        btnex = Btn_ExchangeFont2
    elseif (tonumber(G.GameDefine.ExchangeFontType) == 2) then
        btnex = Btn_ExchangeFont3
    elseif (tonumber(G.GameDefine.ExchangeFontType) == 3) then
        btnex = Btn_ExchangeFont4
    end

    if (user.FlowInfo.Status == 1 and user.FlowInfo.UnionPlatform == 1) then

        if (IsExchangeAchievement()) then
            btnex:ResetSpriteByName(CS.AppDefine.HallNumber, "nexchange", true)
        else
            btnex:ResetSpriteByName(CS.AppDefine.HallNumber, "bexchange", true)
        end
        txt__NowText.text = tostring(user.FlowInfo.FlowAmountNow * 0.01);
        txt__LimitText.text = tostring(user.FlowInfo.FlowAmountLimit * 0.01);
    end

    RefreshBottomAttach()
end

-- 客服按钮
function IsShowCustomServiceButton()
    local isAnd = G.Platform.getCustomServiceOperationAnd()
    local result = G.User.CheckOpenData(isAnd, userdata.ShowCS, userdata.PlaytimeCS, userdata.ShowCSBindPhone)
    if result == true then
        CS.UnityEngine.PlayerPrefs.SetInt("IsStartUpCS", 1)
    end

    return result
end

-- 兑换按钮
function IsShowExchangeButton()
    -- 讀設定
    if G.GameDefine.FlagExchage == 0 then
        return false
    end

    local isAnd = G.Platform.getExchangeOperationAnd()
    local result = G.User.CheckOpenData(isAnd, userdata.ShowConvertLimit, userdata.PlaytimeConvert, userdata.ShowConvertBindPhone)
    return result
end

function IsExchangeAchievement()


    return (user.FlowInfo.Achievement == 1)
end

function PlayCellPhoneAni()
    RecTran_CellPhone.gameObject:SetActive(true)

    local cellPhone = CS.UnityEngine.GameObject.Instantiate(CellPhoneAni.gameObject)
    cellPhone.transform:SetParent(CellPhoneAni.parent.transform)
    cellPhone.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
    cellPhone.transform.localPosition = CS.UnityEngine.Vector3.zero
    cellPhone.transform.gameObject:SetActive(true)

    CS.TimeManager.Instance:DelayInvoke(1.5,
        function()
            RecTran_CellPhone.gameObject:SetActive(false)
            if cellPhone ~= nil then
                CS.UnityEngine.GameObject.Destroy(cellPhone)
                cellPhone = nil
            end

            local openUrl = CurrentUrl .. tostring(G.GameConfig.GamePlatformID)
            CS.UnityEngine.Application.OpenURL(openUrl)
        end)
end