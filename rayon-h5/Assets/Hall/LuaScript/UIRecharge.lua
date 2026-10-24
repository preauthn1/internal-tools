local Btn_Close = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local RecTran_ToggleGroup = this.transform:Find("UIBody/Toggle/ToggleGroup"):GetComponent("RectTransform")
local Tog_AgentPay = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay"):GetComponent("Toggle")
local Img_AgentPayMask = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay/Mask"):GetComponent("Image")
local Tog_AgentRecruit = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentRecruit"):GetComponent("Toggle")
local Tog_JiuLiaoPay = this.transform:Find("UIBody/Toggle/ToggleGroup/_JiuLiaoPay"):GetComponent("Toggle")
local Tog_Henry = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry"):GetComponent("Toggle")
local Tog_HenryMask = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/Mask"):GetComponent("Image")
local Tog_Henry_Gift = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/tag_ad")
local Tog_Henry_OSSM = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/tag_ad_1")
local Tog_HenryFont2 = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/ImageFont2"):GetComponent("Image")
local Tog_HenryMask2 = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/ImageMask2"):GetComponent("Image")

local Tog_888Pay = this.transform:Find("UIBody/Toggle/ToggleGroup/_888Pay"):GetComponent("Toggle")
local Tog_CFPay = this.transform:Find("UIBody/Toggle/ToggleGroup/_CFPay"):GetComponent("Toggle")
local Tog_BaoHePay = this.transform:Find("UIBody/Toggle/ToggleGroup/_BaoHePay"):GetComponent("Toggle")
local Tog_DongGuaZFPay = this.transform:Find("UIBody/Toggle/ToggleGroup/_DongGuaZFPay"):GetComponent("Toggle")
local Tog_DongGuaWeChatPay = this.transform:Find("UIBody/Toggle/ToggleGroup/_DongGuaWeChatPay"):GetComponent("Toggle")

local RecTran_OnlinePayDetail = this.transform:Find("UIBody/_OnlinePayDetail"):GetComponent("RectTransform")
local Txt_CashAmount = this.transform:Find("UIBody/_OnlinePayDetail/_CashAmount"):GetComponent("Text")
local Txt_SafeAmount = this.transform:Find("UIBody/_OnlinePayDetail/_SafeAmount"):GetComponent("Text")
local Btn_ConfirmOnlinePay = this.transform:Find("UIBody/_OnlinePayDetail/_ConfirmOnlinePay"):GetComponent("Button")
local Btn_Clear = this.transform:Find("UIBody/_OnlinePayDetail/_Clear"):GetComponent("Button")
local Btn_Amount = this.transform:Find("UIBody/_OnlinePayDetail/AmountList/_Amount")
local Btn_NoAmount = this.transform:Find("UIBody/_OnlinePayDetail/AmountList/_NoAmount")
local Txt_Note = this.transform:Find("UIBody/_OnlinePayDetail/_Note"):GetComponent("Text")
local Tog_AgentPay_Ad = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay/tag_ad")
local Tog_AgentPay_Ad_1 = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay/tag_ad_1")
--久聊
local img_JiuLiaoPayDisable = this.transform:Find("UIBody/Toggle/ToggleGroup/_JiuLiaoPay/Disable"):GetComponent("Image")
local img_JiuLiaoPay_tag_Close = this.transform:Find("UIBody/Toggle/ToggleGroup/_JiuLiaoPay/tag_Close"):GetComponent("Image")
local Ipt_RechageAmount = this.transform:Find("UIBody/_OnlinePayDetail/_RechageAmount"):GetComponent("InputField")
local txt_RechagePlaceholder = this.transform:Find("UIBody/_OnlinePayDetail/_RechageAmount/Placeholder"):GetComponent("Text")

local RecTran_AgentPayDetail = this.transform:Find("UIBody/_AgentPayDetail"):GetComponent("RectTransform")
local Txt_AgentPayTips = this.transform:Find("UIBody/_AgentPayDetail/_AgentPayTips"):GetComponent("Text")
local Txt_PlayerID = this.transform:Find("UIBody/_AgentPayDetail/_CopyMyID/_PlayerID"):GetComponent("Text")
local Btn_CopyMyID = this.transform:Find("UIBody/_AgentPayDetail/_CopyMyID"):GetComponent("Button")
local RecTran_AgentList = this.transform:Find("UIBody/_AgentPayDetail/AgentList/Grid"):GetComponent("RectTransform")
local Btn_AgentItem = this.transform:Find("UIBody/_AgentPayDetail/_AgentItem"):GetComponent("Button")
local Btn_ReportAgent = this.transform:Find("UIBody/_AgentPayDetail/_ReportAgent"):GetComponent("Button")
local Btn_RefreshAgentList = this.transform:Find("UIBody/_AgentPayDetail/_RefreshAgentList"):GetComponent("Button")
local RecTran_AgentRecuit = this.transform:Find("UIBody/_AgentRecuit"):GetComponent("RectTransform")
local Btn_ConfirmRecuit = this.transform:Find("UIBody/_AgentRecuit/_ConfirmRecuit"):GetComponent("Button")
local Txt_NoDataTips = this.transform:Find("UIBody/_AgentPayDetail/_NoDataTips"):GetComponent("Text")

local Ipt_AgentNote = this.transform:Find("UIBody/_AgentRecuit/_AgentNote"):GetComponent("InputField")
local Ipt_AgentWX = this.transform:Find("UIBody/_AgentRecuit/_AgentWX"):GetComponent("InputField")
local Ipt_AgentQQ = this.transform:Find("UIBody/_AgentRecuit/_AgentQQ"):GetComponent("InputField")
local Ipt_AgentAlipay = this.transform:Find("UIBody/_AgentRecuit/_AgentAlipay"):GetComponent("InputField")
local Ipt_AgentPhone = this.transform:Find("UIBody/_AgentRecuit/_AgentPhone"):GetComponent("InputField")
local Ipt_AgentName = this.transform:Find("UIBody/_AgentRecuit/_AgentName"):GetComponent("InputField")

local RecTran_reportAgent = this.transform:Find("UIBody/_reportAgent"):GetComponent("RectTransform")
local Btn_CloseReport = this.transform:Find("UIBody/_reportAgent/_CloseReport"):GetComponent("Button")
local Ipt_ReportAgentWX = this.transform:Find("UIBody/_reportAgent/_ReportAgentWX"):GetComponent("InputField")
local Ipt_ReportText = this.transform:Find("UIBody/_reportAgent/_ReportText"):GetComponent("InputField")
local Btn_ReportAgentSubmit = this.transform:Find("UIBody/_reportAgent/_ReportAgentSubmit"):GetComponent("Button")
local RecTran_ContactAgent = this.transform:Find("UIBody/_ContactAgent"):GetComponent("RectTransform")
local Btn_CloseContactAgent = this.transform:Find("UIBody/_ContactAgent/_CloseContactAgent"):GetComponent("Button")
local Txt_ContactAgentName = this.transform:Find("UIBody/_ContactAgent/_ContactAgentName"):GetComponent("Text")
local Btn_ShowReportAgent = this.transform:Find("UIBody/_ContactAgent/_ShowReportAgent"):GetComponent("Button")

local Txt_ContactAgentRootQQ = this.transform:Find("UIBody/_ContactAgent/scrolldata/data/QQ")
local Txt_ContactAgentQQ = Txt_ContactAgentRootQQ.transform:Find("_ContactAgent"):GetComponent("Text")
local Btn_ContactAgentQQCopy = Txt_ContactAgentRootQQ.transform:Find("_ContactAgentCopy"):GetComponent("Button")

local Txt_ContactAgentRootWX = this.transform:Find("UIBody/_ContactAgent/scrolldata/data/WX")
local Txt_ContactAgentWX = Txt_ContactAgentRootWX.transform:Find("_ContactAgent"):GetComponent("Text")
local Btn_ContactAgentWXCopy = Txt_ContactAgentRootWX.transform:Find("_ContactAgentCopy"):GetComponent("Button")

local Txt_ContactAgentRootAlipay = this.transform:Find("UIBody/_ContactAgent/scrolldata/data/Alipay")
local Txt_ContactAgentAlipay = Txt_ContactAgentRootAlipay.transform:Find("_ContactAgent"):GetComponent("Text")
local Btn_ContactAgentAlipayCopy = Txt_ContactAgentRootAlipay.transform:Find("_ContactAgentCopy"):GetComponent("Button")

local Txt_ContactAgentRootJiuliao = this.transform:Find("UIBody/_ContactAgent/scrolldata/data/Jiuliao")
local Txt_ContactAgentJiuliao = Txt_ContactAgentRootJiuliao.transform:Find("_ContactAgent"):GetComponent("Text")
local Btn_ContactAgentJiuliaoCopy = Txt_ContactAgentRootJiuliao.transform:Find("_ContactAgentCopy"):GetComponent("Button")

local RecTran_WebPanel = this.transform:Find("UIBody/_WebPanel"):GetComponent("RectTransform")
local RecTran_WebView = this.transform:Find("UIBody/_WebPanel/_View"):GetComponent("RectTransform")
local RecTran_WebGO = this.transform:Find("UIBody/_WebPanel/_WebGO"):GetComponent("RectTransform")

local Btn_CloseWeb = this.transform:Find("UIBody/_WebPanel/_CloseWeb"):GetComponent("Button")

local RecTran_BubblePrompt = this.transform:Find("UIBody/_BubblePrompt"):GetComponent("RectTransform")
local Img_BubblePrompt = this.transform:Find("UIBody/_BubblePrompt/UIBubblePrompt"):GetComponent("Image")
local txt_BubblePrompt = this.transform:Find("UIBody/_BubblePrompt/Date_2"):GetComponent("Text")

local recTran__Loading = this.transform:Find("UIBody/_Loading"):GetComponent("RectTransform")

local Img_AgentPayFont2 = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay/ImageFont2"):GetComponent("Image")
local Img_AgentPayMask2 = this.transform:Find("UIBody/Toggle/ToggleGroup/_AgentPay/ImageMask2"):GetComponent("Image")

local RecTran_HenryDetail = this.transform:Find("UIBody/_HenryDetail"):GetComponent("RectTransform")
local RecTran_HenryList = this.transform:Find("UIBody/_HenryDetail/HenryList/Grid"):GetComponent("RectTransform")
local RecTran_HenryItem = this.transform:Find("UIBody/_HenryDetail/_HenryItem"):GetComponent("RectTransform")
local Txt_HenryNoDataTips = this.transform:Find("UIBody/_HenryDetail/_NoDataTips"):GetComponent("Text")
local Btn_HenryLog = this.transform:Find("UIBody/_HenryDetail/_BtnLog"):GetComponent("Button")

local RecTran_CellPhone = this.transform:Find("UIBody/_HenryDetail/CellPhone"):GetComponent("RectTransform")
local CellPhoneAni = this.transform:Find("UIBody/_HenryDetail/CellPhone/CellphoneAni")

local img_tag_big = this.transform:Find("UIBody/Toggle/ToggleGroup/_Henry/tag_big"):GetComponent("Image")

local img_Title_Pic = this.transform:Find("UIBody/_HenryDetail/Title_Pic"):GetComponent("Image")
local img_Title_Pic1 = this.transform:Find("UIBody/_HenryDetail/Title_Pic1"):GetComponent("Image")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"
local _PayCallBack
local CallPayCount = 0
local DirPayUrl = ""
local ProxyRequestRetCallfun = nil
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback(to_to_yield, cb)
end)

local Json = require("Tools/Json")

-- 在線支付
-- index 方舟索引
-- payType 統一支付索引
local onlineSetingData =
{
    { index = 1, payType = 1, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_Alipay"):GetComponent("Toggle") }, -- 支付宝
    { index = 21, payType = 7, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_WeixinQR"):GetComponent("Toggle") }, -- 微信掃碼
    { index = 11, payType = 6, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_AlipayQR"):GetComponent("Toggle") }, -- 支付寶掃碼

    { index = 19, payType = 106, gb = nil }, -- 冬瓜維信支付
    { index = 18, payType = 105, gb = nil }, -- 冬瓜支付保支付
    { index = 17, payType = 104, gb = nil }, -- 超凡支付
    { index = 16, payType = 103, gb = nil }, -- 超凡支付
    { index = 15, payType = 102, gb = nil }, -- 888支付
    { index = 12, payType = 8, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_AlipayFixed"):GetComponent("Toggle") }, -- 支付寶定額
    { index = 4, payType = 3, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_QQ"):GetComponent("Toggle") }, -- QQ
    { index = 2, payType = 2, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_Weixin"):GetComponent("Toggle") }, -- 微信
    { index = 3, payType = 4, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_Unionpay"):GetComponent("Toggle") }, -- 银联	
    { index = 5, payType = 5, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_JD"):GetComponent("Toggle") }, -- 京東
    { index = 22, payType = 9, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_WeixinFixed"):GetComponent("Toggle") }, -- 微信固额
    { index = 13, payType = 12, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_RedEnvelope"):GetComponent("Toggle") }, -- 紅包支付
    { index = 6, payType = 101, gb = nil }, -- 久聊支付
    { index = 7, payType = 11, gb = this.transform:Find("UIBody/Toggle/ToggleGroup/_YunShan"):GetComponent("Toggle") }, -- 闪云支付
}

-- 外部协议白名单
schemes = {
    "weixin", -- 微信
    "alipays", -- 支付宝
    "alipay", -- 支付宝
    "mqqwpa", -- QQ
    "mqqapi", -- QQ
    "jiuliao", -- 久聊
}

local apiaconfig = {}

-- api table 對應固定金額
local staticpayment = { 8, 9, 101, 102, 105, 106 }
-- api table 對應开新画面不走WEBVIEW
local redirecpayment = { 9, 7, 6 }
local GoShield = { 4 }
local LastPayType = nil
local LastPayGWID = nil
local LastAmount = nil
local userdata = G.User.GetCurrentUser().UserData
local user = G.User.GetCurrentUser()
local groupName = "UIRecharge"

local payTypeTogs = {}
local amountBtns = {}
local isloadpayurl = false -- 是否已经载入过充值页面，每次重新充值是设置为false
local isAppError = false -- 是否已經收到過自動開啟支付失敗訊息，每次重新充值是设置为false
local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local curWebView
local pause = false;

local accTime = 0
local lastGetTime = -10

local currentPayType = 0
local currentGwid = 0
function SetCurrentPayType(Type, gwid)
    print(Type, gwid)
    currentPayType = Type
    currentGwid = gwid
end

--当下选择付费方式
function GetCurrentPayType()
    return currentPayType, currentGwid;
end

--方舟支付編號
function GetCurrentPayIndex()
    for _, data in pairs(onlineSetingData) do
        if data.payType == currentPayType then
            return data.index
        end
    end
end


-- 加载充值金额
function loadChargeAmountUI()
    --移除金额按钮
    for _, btn in pairs(amountBtns) do
        if (btn ~= nil) then
            CS.UnityEngine.GameObject.Destroy(btn)
        end
    end

    Ipt_RechageAmount.text = ''
    amountBtns = {}

    local _payType, gwid = GetCurrentPayType()
    local _IsFixedAmt = (G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1][8] == 0)

    -- 是否為固額

    Ipt_RechageAmount.interactable = not _IsFixedAmt;

    -- [1]是使用find皆會給一組流水號,[3]为统一支付给予的金额清单
    G.Tool.TableHelper.Print(apiaconfig)

    local aconfig = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1][3]

    local aconfig1 = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1];
    local onlinepaymin = 0
    if (aconfig1 ~= nil) then
        onlinepaymin = math.modf(aconfig1[1] / 100)
    end

    local checkInput = false
    for _, v in pairs(aconfig) do
        --提示预设讯息或最低金额
        if checkInput == false then
            if onlinepaymin > 0 then
                if G.Platform.IsChangeCoin() == true then
                    txt_RechagePlaceholder.text = (_IsFixedAmt and '请选择' or '请输入') .. '充值金额，最低充值金额为' .. tostring(onlinepaymin) .. '币';
                else
                    txt_RechagePlaceholder.text = (_IsFixedAmt and '请选择' or '请输入') .. '充值金额，最低充值金额为' .. tostring(onlinepaymin) .. '元';
                end
            else
                txt_RechagePlaceholder.text = (_IsFixedAmt and '请选择' or '请输入') .. '充值金额';
            end
            checkInput = true
        end

        local btn = CS.UnityEngine.GameObject.Instantiate(Btn_Amount.gameObject)
        btn.transform:SetParent(Btn_Amount.gameObject.transform.parent.transform)
        btn.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
        btn.transform.localPosition = CS.UnityEngine.Vector3(0, 0, 0)

        local btnamt = math.modf(v / 100)
        btn.transform:Find("txt"):GetComponent("Text").text = (_IsFixedAmt and '' or '+') .. tostring(btnamt)

        btn:GetComponent("Button"):AddClickListener(function()
            AddMoneyAmount(btnamt)
        end)
        btn:SetActive(true)
        table.insert(amountBtns, btn)
    end
end

function CloseWebView()
    print("CloseWebView")
    local wb = RecTran_WebGO:GetComponent("UniWebView")
    if (wb ~= nil) then
        wb:Hide()
        CS.UnityEngine.GameObject.Destroy(wb)
    end
    curWebView = nil
    RecTran_WebPanel.gameObject:SetActive(false)
end

-- 初始化函数，实例化之后调用一次
function Start()
    -- 闪付广告
    if G.Platform.GetHenryAD() == true then
        img_Title_Pic.gameObject:SetActive(false)
        img_Title_Pic1.gameObject:SetActive(true)
    else
        img_Title_Pic.gameObject:SetActive(true)
        img_Title_Pic1.gameObject:SetActive(false)
    end

    -- 加贈2% by AZLDC
    if G.Platform.getNoTwoPa() == false then
        print("AZLDC : 關閉加贈!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
        Tog_Henry_Gift.gameObject:SetActive(false)
    else
        print("AZLDC : 開啟加贈!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
        Tog_Henry_Gift.gameObject:SetActive(true)
    end

    -- 推薦標籤 by AZLDC
    -- Tog_Henry_OSSM.gameObject:SetActive( G.Platform.getOSuSuMe() )
    if G.Platform.getOSuSuMe() == false then
        print("AZLDC : 關閉推薦!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
        Tog_Henry_OSSM.gameObject:SetActive(false)
    else
        print("AZLDC : 開啟推薦!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
        Tog_Henry_OSSM.gameObject:SetActive(true)
    end

    -- 替換官方支付顯示為VIP支付 by AZLDC
    if G.Platform.getHenryToVIP() == true then
        Tog_Henry.image.sprite = Tog_HenryFont2.sprite
        Tog_HenryMask.sprite = Tog_HenryMask2.sprite
    end

    -- 986顯示1%
    if G.Platform.GetBranch() == 1 then
        Tog_AgentPay_Ad_1.gameObject:SetActive(true)
    else
        Tog_AgentPay_Ad_1.gameObject:SetActive(false)
    end

    -- 蓝钻闪付支付显示大额支付
    if G.Platform.GetBranch() == 2 then
        img_tag_big.gameObject:SetActive(true)
    else
        img_tag_big.gameObject:SetActive(false)
    end

    CS.AppManager.Instance.OnAppPause = CS.AppManager.Instance.OnAppPause + AppPause

    local agentNameID = userdata.AgentNameID
    if agentNameID == 2 then
        Tog_AgentPay.image.sprite = Img_AgentPayFont2.sprite
        Img_AgentPayMask.sprite = Img_AgentPayMask2.sprite
    end

    if G.Service.UIFlag.IsShowJiuLiao() == false then
        Tog_JiuLiaoPay.gameObject:SetActive(false);
    end

    --強制關閉
    if CS.Utility.GetbundleIdentifier() == "com.stage65.test" then
        Tog_JiuLiaoPay.gameObject:SetActive(false);
    end

    RecTran_BubblePrompt.gameObject:SetActive(false)
    Btn_AgentItem.gameObject:SetActive(false)

    Btn_CloseWeb:AddClickListener(CloseWebView)
    Btn_Close:AddClickListener(function()
        G.Tool.UILoading.Hide()
        CS.UILoadingData.Hide()
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
    end)
    --處理內頁金額
    for _, data in pairs(onlineSetingData) do
        if data.gb ~= nil then
            --  print("onValueChanged",data.payType)
            data.gb.onValueChanged:AddListener(function(isOn)

                -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
                if isOn == false then
                    RecTran_OnlinePayDetail.gameObject:SetActive(false)
                    return
                end

                InitPayConfig()

                -- 内页点选时重新抽选
                local paylist = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == data.index end);
                local len = tablelength(paylist);
                local assignnum = math.random(1, len);
                SetCurrentPayType(data.index, paylist[assignnum][7])
                ClearRechargeBtn()
                if data.index == 7 then
                    Txt_Note.gameObject:SetActive(true)
                else
                    Txt_Note.gameObject:SetActive(false)
                end

                local url = G.GameConfig.RechargeListDomain .. "/PayCenter/amountlist/" .. tostring(G.GameConfig.GamePlatformID)
                G.Tool.Debug.RecordGet({
                    content = {
                        ContentType = 2,
                        Flag = 1,
                        Data = tostring(url)
                    },
                    lt = 1,
                    sceneName = "UIHALL",
                    uid = G.User.GetCurrentUser().UserID
                })

                payConfig = G.Service.ReCharge.GetPayConfig()
                ProcessAmtList(payConfig)
                RecTran_OnlinePayDetail.gameObject:SetActive(true)
            end)
        end
    end

    Tog_JiuLiaoPay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(6, 6)
        apiaconfig = G.Service.ReCharge.GetPaymentInfoList(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)
    --冬瓜维信支付侧边列點擊
    Tog_DongGuaWeChatPay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(19, 19)
        apiaconfig = G.Service.ReCharge.GetDongGuaWeChatPayInfo(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)
    --冬瓜支付宝支付侧边列點擊
    Tog_DongGuaZFPay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(18, 19)

        apiaconfig = G.Service.ReCharge.GetDongGuaZFPayInfo(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)
    --宝盒支付侧边列點擊
    Tog_BaoHePay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(17, 16)
        apiaconfig = G.Service.ReCharge.GetBaoHePayInfo(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)

    --超凡支付侧边列點擊
    Tog_CFPay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(16, 11)
        apiaconfig = G.Service.ReCharge.GetCFPayInfo(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)
    --888支付侧边列點擊
    Tog_888Pay.onValueChanged:AddListener(function(isOn)

        -- isOn是關掉的情況不必重新抓資料和顯示 by Rayon
        if isOn == false then
            RecTran_OnlinePayDetail.gameObject:SetActive(false)
            return
        end

        ClearRechargeBtn()
        SetCurrentPayType(15, 10)
        apiaconfig = G.Service.ReCharge.Get888PayInfo(1);
        CS.UILoadingData.Show("读取中....,请稍后...", 5)
        ProcessPayList()
        RecTran_OnlinePayDetail.gameObject:SetActive(isOn)
    end)

    Tog_Henry.onValueChanged:AddListener(function(isOn)

        RecTran_HenryDetail.gameObject:SetActive(isOn)
        if isOn == true then
            ClickHenry()
        end
    end)


    function ShowErrorMessage()
        G.Util.coroutine_call(function()
            local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "支付渠道忙碌中,请再试一次...", "确定");
            if boxResult == CS.MessageBoxResult.OK then
                CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
            end
        end)()
    end

    --確認支付點擊
    Btn_ConfirmOnlinePay:AddClickListener(function()
        local _payType, gwid = GetCurrentPayType()
        if G.GameConfig.GetServerOption() == 1 or G.GameConfig.GetServerOption() == 2 then
            ShowBubblePrompt("方舟:" .. _payType .. " 第三方:" .. gwid)
        end


        -- 久聊
        if _payType == "101" then
            nowPayConfig = apiaconfig
            --888    
        elseif _payType == "102" then
            nowPayConfig = apiaconfig
            -- 超凡CF    
        elseif _payType == "103" then
            nowPayConfig = apiaconfig
            -- BaoHe
        elseif _payType == "104" then
            nowPayConfig = apiaconfig
        elseif _payType == "105" then
            nowPayConfig = apiaconfig
        elseif _payType == "106" then
            nowPayConfig = apiaconfig
            -- 其他
        else
            nowPayConfig = GetNowPayConfig()
        end

        if nowPayConfig == nil then
            -- 检查充值资料有错误
            ShowErrorMessage()
            return
        end

        local aconfig = G.Tool.TableHelper.Find(nowPayConfig, function(v) return v[6] == _payType and v[7] == gwid end)[1];
        if (aconfig == nil) then
            -- 检查充值资料有错误
            ShowErrorMessage()
            return
        end

        local onlinepaymin = math.modf(aconfig[1] / 100)
        local onlinepaymax = math.modf(aconfig[2] / 100)
        local payAmountMin = onlinepaymin
        local payAmountMax = onlinepaymax
        local currentMoneyAmount = tonumber(G.Tool.StringHelper.IsEmpty(Ipt_RechageAmount.text) and 0 or Ipt_RechageAmount.text)

        if (currentMoneyAmount > payAmountMax) then
            Ipt_RechageAmount.text = tostring(payAmountMax)
            CS.UIBubblePrompt.Show("最高充值金额为" .. payAmountMax, groupName)
            return
        elseif (currentMoneyAmount < payAmountMin) then
            Ipt_RechageAmount.text = tostring(payAmountMin)
            CS.UIBubblePrompt.Show("最少充值金额为" .. payAmountMin, groupName)
            return
        end

        -- 若為固額	檢查輸入金額是否相符
        local _IsFixedAmt = (G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1][8] == 0)
        if _IsFixedAmt == true then
            local isFind = false
            local inputValue = tonumber(Ipt_RechageAmount.text)
            local aconfigValue = G.Tool.TableHelper.Find(nowPayConfig, function(v) return v[6] == _payType and v[7] == gwid end)[1][3]
            for _, v in pairs(aconfigValue) do
                local btnamt = math.modf(v / 100)
                if inputValue == btnamt then
                    isFind = true
                    break
                end
            end

            -- 检查充值资料有错误
            if isFind == false then
                -- 检查充值资料有错误
                ShowErrorMessage()
                return
            end
        end

        local sendMsg = {}
        -- local pt = tostring(GetCurrentPayIndex())
        local pt = _payType

        sendMsg.PayTypeID = tonumber(pt)
        sendMsg.Amount = tonumber(Ipt_RechageAmount.text)
        LastPayType = tonumber(pt)
        LastPayGWID = tonumber(gwid)
        LastAmount = tonumber(Ipt_RechageAmount.text)

        CS.UILoadingData.Show("正在创建订单,请稍后...", 10,
            (function()
                CS.UIBubblePrompt.Show("创建订单失败", groupName)
            end))
        G.Tool.Debug.RecordGet({
            content = {
                ContentType = 2,
                Flag = 3,
                Data = "PayType:" .. tostring(pt),
                Amt = tostring(LastAmount),
            },
            lt = 1,
            sceneName = "UIHALL",
            uid = G.User.GetCurrentUser().UserID
        })
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.CreatePayOrderID, sendMsg)
    end)
    --清除金額欄位
    Btn_Clear:AddClickListener(function()
        local _payType, gwid = GetCurrentPayType()
        local aconfig = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1];
        local endeditmin = math.modf(aconfig[1] / 100)
        G.Tool.TableHelper.Print(apiaconfig)
        AddMoneyAmount(-1)
    end)

    --代理支付
    Tog_AgentPay.onValueChanged:AddListener(function(isOn)
        RecTran_AgentPayDetail.gameObject:SetActive(isOn)
    end)

    Btn_CopyMyID:AddClickListener(function()
        local user = G.User.GetCurrentUser()
        local msg = "我是「" .. G.GameDefine.AppName .. "遊戲」的玩家，"
        msg = msg .. "我的ID：" .. Txt_PlayerID.text
        msg = msg .. " 充值金額："

        CS.Utility.CopyText(msg)
        CS.UIBubblePrompt.Show("复制成功", groupName)
        -- CS.UnityEngine.Application.OpenURL("weixin://")
    end)
    --舉報代理
    Btn_ReportAgent:AddClickListener(function()
        RecTran_reportAgent.gameObject:SetActive(true)
        Ipt_ReportAgentWX.text = ""
    end)
    --刷新代理清單
    Btn_RefreshAgentList:AddClickListener(function()
        RefreshAgentList(true)
    end)
    --招募代理
    Tog_AgentRecruit.onValueChanged:AddListener(function(isOn)
        RecTran_AgentRecuit.gameObject:SetActive(isOn)
    end)
    --送出招募代理
    Btn_ConfirmRecuit:AddClickListener(function()

        if G.Tool.StringHelper.IsEmpty(Ipt_AgentName.text) then
            CS.UIBubblePrompt.Show("请输入姓名", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_AgentPhone.text) then
            CS.UIBubblePrompt.Show("请输入电话号码", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_AgentQQ.text) then
            CS.UIBubblePrompt.Show("请输入QQ号码", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_AgentWX.text) then
            CS.UIBubblePrompt.Show("请输入微信号", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_AgentAlipay.text) then
            CS.UIBubblePrompt.Show("请输入支付宝号", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_AgentNote.text) then
            CS.UIBubblePrompt.Show("请输入说明信息", groupName)
            return
        end

        local sendMsg = {}

        sendMsg.AreaID = G.GameConfig.GameSiteID
        sendMsg.Name = Ipt_AgentName.text
        sendMsg.Tel = Ipt_AgentPhone.text
        sendMsg.QQ = Ipt_AgentQQ.text
        sendMsg.WXNo = Ipt_AgentWX.text
        sendMsg.Alipay = Ipt_AgentAlipay.text
        sendMsg.Memo = Ipt_AgentNote.text

        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ApplyAgentID, sendMsg)
    end)
    --關閉回報頁面
    Btn_CloseReport:AddClickListener(function()
        RecTran_reportAgent.gameObject:SetActive(false)
    end)
    --送出舉報內容
    Btn_ReportAgentSubmit:AddClickListener(function()

        if G.Tool.StringHelper.IsEmpty(Ipt_ReportText.text) then
            CS.UIBubblePrompt.Show("请输入举报内容", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_ReportAgentWX.text) then
            CS.UIBubblePrompt.Show("请输入代理微信号", groupName)
            return
        end

        local sendMsg = {}

        sendMsg.AreaID = G.GameConfig.GameSiteID
        sendMsg.AgentName = Ipt_ReportAgentWX.text
        sendMsg.Content = Ipt_ReportText.text

        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ReportAgentID, sendMsg)
    end)
    --複製QQ
    Btn_ContactAgentQQCopy:AddClickListener(function()
        CS.Utility.CopyText(Txt_ContactAgentQQ.text)
        CS.UIBubblePrompt.Show("已复制代理资讯,3秒后自动打开QQ", groupName)
        local sequence = CS.DG.Tweening.DOTween.Sequence()

        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            CS.UnityEngine.Application.OpenURL("mqqwpa://im/chat?chat_type=wpa&uin=" .. Txt_ContactAgentQQ.text)
        end)
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            if (pause) then
                pause = false;
            else
                CS.UnityEngine.Application.OpenURL("https://example.invalid/download/")
            end
        end)
    end)
    --複製維信
    Btn_ContactAgentWXCopy:AddClickListener(function()
        CS.Utility.CopyText(Txt_ContactAgentWX.text)
        CS.UIBubblePrompt.Show("已复制代理资讯,3秒后自动打开微信", groupName)
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            CS.UnityEngine.Application.OpenURL("weixin://")
        end)
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            if (pause) then
                pause = false;
            else
                CS.UnityEngine.Application.OpenURL("https://example.invalid/")
            end
        end)
    end)

    Btn_ContactAgentAlipayCopy:AddClickListener(function()
        CS.Utility.CopyText(Txt_ContactAgentAlipay.text)
        CS.UIBubblePrompt.Show("已复制代理资讯,3秒后自动打开支付宝", groupName)
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            CS.UnityEngine.Application.OpenURL("alipays://")
        end)
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            if (pause) then
                pause = false;
            else
                CS.UnityEngine.Application.OpenURL("https://example.invalid/login/index.htm?goto=https%3A%2F%example.invalid%2Fportal%2Fi.htm")
            end
        end)
    end)

    Btn_ContactAgentJiuliaoCopy:AddClickListener(function()

        CS.Utility.CopyText(Txt_ContactAgentJiuliao.text)
        if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
            CS.UIBubblePrompt.Show("已复制代理资讯,请开启久聊添加好友", groupName)
        elseif CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
            CS.UIBubblePrompt.Show("已复制代理资讯,3秒后自动打开久聊,若无法自动打开时，请手动打开", groupName)
        end
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        sequence:AppendInterval(2)
        sequence:AppendCallback(function()
            CS.UnityEngine.Application.OpenURL("jiuliao://")
        end)

        -- IOS沒安裝久聊跳轉到官網
        if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
            sequence:AppendInterval(2)
            sequence:AppendCallback(function()
                if (pause) then
                    pause = false;
                else
                    CS.UnityEngine.Application.OpenURL("http://example.invalid")
                end
            end)
        end
    end)

    Btn_CloseContactAgent:AddClickListener(function()
        RecTran_ContactAgent.gameObject:SetActive(false)
    end)

    Btn_ShowReportAgent:AddClickListener(function()
        RecTran_ContactAgent.gameObject:SetActive(false)
        RecTran_reportAgent.gameObject:SetActive(true)
    end)
    --輸入金額驗證
    Ipt_RechageAmount.onEndEdit:AddListener(function()
        if G.Tool.StringHelper.IsEmpty(Ipt_RechageAmount.text) then
            Ipt_RechageAmount.text = "0"
        else
            local current = G.Tool.Math.ParseInt(Ipt_RechageAmount.text)
            -- local aconfig = G.Service.ReCharge.GetPayAmountRange()
            -- add by bob
            local _payType, gwid = GetCurrentPayType()
            local aconfig = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1];
            if (aconfig ~= nil) then

                local endeditmin = math.modf(aconfig[1] / 100)
                local endeditmax = math.modf(aconfig[2] / 100)
                local payAmountMin = endeditmin
                local payAmountMax = endeditmax
                if (current > payAmountMax) then
                    current = payAmountMax
                end
            end

            Ipt_RechageAmount.text = tostring(current)
        end
    end)

    Btn_HenryLog:AddClickListener(function()
        PlayCellPhoneAni(-1)
    end)
end

function ProcessPayList()
    Btn_ConfirmOnlinePay.gameObject:SetActive(true)
    recTran__Loading.gameObject:SetActive(false);
    CS.UILoadingData.Hide()
    G.Tool.UILoading.Hide()
    loadChargeAmountUI()
end

function ProcessAmtList(data)
    if (data == "" or data == nil) then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end

    if (data.data == nil or data.data == "") then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end

    recTran__Loading.gameObject:SetActive(false);
    CS.UILoadingData.Hide()

    -- 只有需要显示充值才加载充值类型
    if (G.Service.UIFlag.IsShowOnlinePay() == true) then
        --  local types = G.Tool.StringHelper.Split(userdata.PaySort, ",")
        local types = apiaconfig
        G.Tool.UILoading.Hide()
        recTran__Loading.gameObject:SetActive(false);

        local index = 1;

        -- 修正當開啟的號碼非連號會掛掉的問題 by Bob
        -- for _, v in pairs(types) do
        for i = 1, tablelength(types), 1 do
            if index == 1 then
                loadChargeAmountUI()
            end
            index = index + 1;
        end
    end
    Btn_ConfirmOnlinePay.gameObject:SetActive(true)
end

--指定使用内建浏览器开启
function ProcessOrderRedirctRet(data)

    if (data == "" or data == nil) then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end

    if (data.data == nil or data.data == "") then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end

    G.Tool.Debug.RecordGet({
        content = {
            ContentType = 2,
            Flag = 6,
            Data = tostring(data.data)
        },
        lt = 51,
        sceneName = "UIHALL",
        uid = G.User.GetCurrentUser().UserID
    })

    -- Android 用储存本地
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor then

        if (string.find(data.data, 'n.href="', 1, true) ~= nil) then
            print("用解析URL ")
            local startindex = string.find(data.data, 'n.href="', 1, true) + 8;
            local tmpstr = G.Tool.StringHelper.SubString(data.data, startindex, G.Tool.StringHelper.length(data.data) - startindex);
            --雙引號計算多一個   
            local Len = string.find(tmpstr, '";</script>', 1, true) - 1;
            local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
            print("DirPayUrl", url)
            CS.UnityEngine.Application.OpenURL(url)
            return
        elseif (string.find(data.data, 'n.replace(\'', 1, true) ~= nil) then
            print("Replace用解析URL ");
            local startindex = string.find(data.data, 'n.replace(\'', 1, true) + 11;
            local tmpstr = G.Tool.StringHelper.SubString(data.data, startindex, G.Tool.StringHelper.length(data.data) - startindex);
            local Len = string.find(tmpstr, '\');</script>', 1, true) - 1;
            local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
            print("DirPayUrl", url)
            CS.UnityEngine.Application.OpenURL(url)
            return
        else
            print("NO IP Redirect ")
            ProcessOrderRet(data);
        end

        return
    end

    -- IOS 用解析URL

    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor then
        print("IOS~~~~ ")
        if (string.find(data.data, 'n.href="', 1, true) ~= nil) then
            print("用解析URL ")
            local startindex = string.find(data.data, 'n.href="', 1, true) + 8;
            local tmpstr = G.Tool.StringHelper.SubString(data.data, startindex, G.Tool.StringHelper.length(data.data) - startindex);
            --雙引號計算多一個   
            local Len = string.find(tmpstr, '";</script>', 1, true) - 1;
            local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
            print("DirPayUrl", url)
            CS.UnityEngine.Application.OpenURL(url)
            return
        elseif (string.find(data.data, 'n.replace(\'', 1, true) ~= nil) then
            print("Replace用解析URL ");
            local startindex = string.find(data.data, 'n.replace(\'', 1, true) + 11;
            local tmpstr = G.Tool.StringHelper.SubString(data.data, startindex, G.Tool.StringHelper.length(data.data) - startindex);
            local Len = string.find(tmpstr, '\');</script>', 1, true) - 1;
            local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
            print("DirPayUrl", url)
            CS.UnityEngine.Application.OpenURL(url)
            return

        else
            ProcessOrderRet(data);
        end
    end
end

--指定使用WEBVIEW开启
function ProcessOrderRet(data)
    if (data == "" or data == nil) then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end

    if (data.data == nil or data.data == "") then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
        return
    end
    G.Tool.Debug.RecordGet({
        content = {
            ContentType = 2,
            Flag = 6,
            Data = tostring(data.data)
        },
        lt = 50,
        sceneName = "UIHALL",
        uid = G.User.GetCurrentUser().UserID
    })

    local function GoWebView(payhtml)
        CS.UILoadingData.Hide()
        RecTran_WebPanel.gameObject:SetActive(true)

        local webView = RecTran_WebGO:GetComponent("UniWebView")
        if (webView == nil) then
            webView = RecTran_WebGO.gameObject:AddComponent(typeof(CS.UniWebView))
        end
        webView:CleanCache()
        CS.UILoadingData.Show("正在启动支付,请稍后...", 5)

        -- 4个参数 分别为 距离上边距离，左边距离，下边距离，右边距离
        local screenWidth = CS.UniWebViewHelper.screenWidth * 0.5
        local screenHeight = CS.UniWebViewHelper.screenHeight * 0.5
        local rectWidth = RecTran_WebPanel.rect.width * 0.5
        local rectHeight = RecTran_WebPanel.rect.height * 0.5
        local scale = screenHeight / rectHeight
        local top = math.modf((rectHeight - (RecTran_WebView.rect.height * 0.5 + RecTran_WebView.anchoredPosition.y)) * scale)
        local bottom = math.modf((rectHeight - (RecTran_WebView.rect.height * 0.5 - RecTran_WebView.anchoredPosition.y)) * scale)

        scale = screenWidth / rectWidth
        local left = math.modf((rectWidth - (RecTran_WebView.rect.width * 0.5 - RecTran_WebView.anchoredPosition.x)) * scale)
        local right = math.modf((rectWidth - (RecTran_WebView.rect.width * 0.5 + RecTran_WebView.anchoredPosition.x)) * scale)
        -- print(top,left,bottom,right)
        webView.insets = CS.UniWebViewEdgeInsets(top, left, bottom, right)
        webView:OnLoadComplete('+', OnLoadComplete)
        webView:OnReceivedMessage('+', OnReceivedMessage)
        -- 按下返回鍵不關掉webview，\Assets\FrameWork\XLuaWebViewExport中加上typeof(UniWebView.WebViewShouldCloseDelegate)
        --webView:OnWebViewShouldClose('+', OnWebViewShouldClose)

        -- 加入外部协议白名单
        for i = 1, #schemes do
            webView:AddUrlScheme("xy:" .. schemes[i])
        end

        -- 载入页面时不显示loading
        --webView:SetShowSpinnerWhenLoading(false)
        webView.toolBarShow = false

        webView:LoadHTMLString(payhtml, nil)
        curWebView = webView
        isloadpayurl = false
        isAppError = false
    end

    -- 避免解析錯誤
    data.data = string.gsub(data.data, "\"", "'")

    -- Android 用储存本地
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor then
        print("ANDROID")
        --      local function WritePay(payhtml)
        --      CS.UILoadingData.Hide()
        --      path = example.invalid.Path.Combine(CS.UnityEngine.Application.dataPath, tostring(os.time()) .. ".html");
        --      fileStream_Open = example.invalid.FileStream(path, example.invalid.FileMode.OpenOrCreate);
        --      streamWriter2 = example.invalid.StreamWriter(fileStream_Open, CS.System.Text.UTF8Encoding(true));
        --      streamWriter2:Write(payhtml);
        --      streamWriter2:Flush();
        --      streamWriter2:Close();			
        --      w = CS.UnityEngine.WWW(path);
        --      CS.UnityEngine.Application.OpenURL(w.url);
        --      end 
        _PayCallBack = GoWebView;
        PaymentPass(data.data, GoWebView);

        return
    end

    -- IOS 用解析URL
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor then
        print("IOS~~~~ ")
        _PayCallBack = GoWebView;
        PaymentPass(data.data, GoWebView);
        return;
    end

    CS.UIBubblePrompt.Show("支付渠道忙碌中请稍后再试...", groupName)
end

function PaymentPass(data, p)

    if (string.find(data, 'example.invalid', 1, true) ~= nil) then
        local startindex = string.find(data, 'example.invalid', 1, true) - 8;
        local tmpstr = G.Tool.StringHelper.SubString(data, startindex, G.Tool.StringHelper.length(data) - startindex);
        local Len = string.find(tmpstr, '\';</script>', 1, true) - 1;
        local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
        print("ali url", url)
        CS.UnityEngine.Application.OpenURL(url);
    elseif (string.find(data, 'tiaozhuan', 1, true) ~= nil) then
        local startindex = string.find(data, 'example.invalid', 1, true) - 8;
        local tmpstr = G.Tool.StringHelper.SubString(data, startindex, G.Tool.StringHelper.length(data) - startindex);
        --雙引號計算多一個
        local Len = string.find(tmpstr, '";', 1, true) - 1;
        local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
        print("tiaozhuan url", url)
        if (CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor) then
            DirPayUrl = url;
            GetProxyRequest(url, ProcessPayRet);
        else
            CS.UnityEngine.Application.OpenURL(url);
        end

    elseif (string.find(data, 'n.href="', 1, true) ~= nil) then
        print("用解析URL ")
        local startindex = string.find(data, 'n.href="', 1, true) + 8;
        local tmpstr = G.Tool.StringHelper.SubString(data, startindex, G.Tool.StringHelper.length(data) - startindex);
        --雙引號計算多一個   
        local Len = string.find(tmpstr, '";</script>', 1, true) - 1;
        local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
        print("DirPayUrl", url)
        CS.UnityEngine.Application.OpenURL(url)
        return
    elseif (string.find(data, 'n.replace(\'', 1, true) ~= nil) then
        print("Replace用解析URL ");
        local startindex = string.find(data, 'n.replace(\'', 1, true) + 11;
        local tmpstr = G.Tool.StringHelper.SubString(data, startindex, G.Tool.StringHelper.length(data) - startindex);
        local Len = string.find(tmpstr, '\');</script>', 1, true) - 1;
        local url = G.Tool.StringHelper.SubString(tmpstr, 1, Len);
        print("DirPayUrl", url)
        CS.UnityEngine.Application.OpenURL(url)
        return
    else
        print("NO IP Process ")
        p(data);
    end
end

--循环解析HTML
function ProcessPayRet(data)
    print("<color=yellow>ProcessPayRet</color>")
    print(url)

    if (CallPayCount > 5) then
        _PayCallBack = nil;
        CallPayCount = 0;
        DirPayUrl = "";
        CS.UILoadingData.Hide()
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        return
    end

    if (data == "" or data == nil) then
        CallPayCount = CallPayCount + 1
        GetProxyRequest(DirPayUrl, ProcessPayRet)
        return
    end

    if (data.data == nil or data.data == "") then
        CallPayCount = CallPayCount + 1
        GetProxyRequest(DirPayUrl, ProcessPayRet)
        return
    end
    print("<color=yellow>Call PayCallBack </color>")
    print(data.data)
    _PayCallBack(data.data);
    _PayCallBack = nil;
    CallPayCount = 0;
    DirPayUrl = "";
end

function RefreshAgentList(refresh)
    local sendMsg = {}
    sendMsg.AreaID = G.GameConfig.GameSiteID

    if G.Agent.GetAgentInfo() == nil or refresh then
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentListID, sendMsg)
    else
        CreateAgentList()
    end
end

function CreateAgentList()
    local agentList = G.Agent.GetAgentInfo()

    if #agentList > 0 then
        Txt_NoDataTips.gameObject:SetActive(false)
    else
        Txt_NoDataTips.gameObject:SetActive(true)
    end

    for i = 1, #agentList, 1 do
        local agentItem = CS.UnityEngine.GameObject.Instantiate(Btn_AgentItem.gameObject)
        agentItem.transform:SetParent(RecTran_AgentList.transform)
        agentItem.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
        agentItem.transform.localPosition = CS.UnityEngine.Vector3.zero

        agentItem.transform.gameObject:SetActive(true)
        local txtName = agentItem.transform:Find("Name"):GetComponent("Text")
        txtName.text = agentList[i].Name

        local imgQQ = agentItem.transform:Find("Grid/ImageQQ")
        imgQQ.gameObject:SetActive(agentList[i].QQ ~= "" and agentList[i].QQ ~= "0")
        local imgWXNo = agentItem.transform:Find("Grid/ImageWXNo")
        imgWXNo.gameObject:SetActive(agentList[i].WXNo ~= "" and agentList[i].WXNo ~= "0")
        local imgAlipay = agentItem.transform:Find("Grid/ImageAlipay")
        imgAlipay.gameObject:SetActive(agentList[i].Alipay ~= "" and agentList[i].Alipay ~= "0")
        local imgJiuliao = agentItem.transform:Find("Grid/ImageJiuliao")
        imgJiuliao.gameObject:SetActive(agentList[i].Jiuliao ~= "" and agentList[i].Jiuliao ~= "0")

        openText = agentItem.transform:Find("OpenText")
        openText.gameObject:SetActive(false)
        infoText = agentItem.transform:Find("InfoText")
        infoText.gameObject:SetActive(true)
        if G.Platform.GetBranch() == 2 then
            if imgJiuliao.gameObject.activeSelf == true then
                imgQQ.gameObject:SetActive(false)
                imgWXNo.gameObject:SetActive(false)
                imgAlipay.gameObject:SetActive(false)
                imgJiuliao.gameObject:SetActive(false)
                openText.gameObject:SetActive(true)
                infoText.gameObject:SetActive(false)
            end
        end

        agentItem.transform:GetComponent("Button"):AddClickListener(function()
            if G.Platform.GetBranch() == 2 then
                if agentList[i].Jiuliao ~= "" and agentList[i].Jiuliao ~= "0" then
                    CS.UnityEngine.Application.OpenURL(agentList[i].Jiuliao)
                    return
                end
            end

            RecTran_ContactAgent.gameObject:SetActive(true)
            Ipt_ReportAgentWX.text = txtName.text

            Txt_ContactAgentQQ.text = agentList[i].QQ
            Txt_ContactAgentName.text = agentList[i].Name
            Txt_ContactAgentWX.text = agentList[i].WXNo
            Txt_ContactAgentAlipay.text = agentList[i].Alipay
            Txt_ContactAgentJiuliao.text = agentList[i].Jiuliao

            if agentList[i].QQ == "" or agentList[i].QQ == "0" then
                Txt_ContactAgentRootQQ.gameObject:SetActive(false)
            else
                Txt_ContactAgentRootQQ.gameObject:SetActive(true)
            end

            if agentList[i].WXNo == "" or agentList[i].WXNo == "0" then
                Txt_ContactAgentRootWX.gameObject:SetActive(false)
            else
                Txt_ContactAgentRootWX.gameObject:SetActive(true)
            end

            if agentList[i].Alipay == "" or agentList[i].Alipay == "0" then
                Txt_ContactAgentRootAlipay.gameObject:SetActive(false)
            else
                Txt_ContactAgentRootAlipay.gameObject:SetActive(true)
            end

            if agentList[i].Jiuliao == "" or agentList[i].Jiuliao == "0" then
                Txt_ContactAgentRootJiuliao.gameObject:SetActive(false)
            else
                Txt_ContactAgentRootJiuliao.gameObject:SetActive(true)
            end
        end)
    end
end

function ClearRechargeBtn()
    --移除金额按钮
    for _, btn in pairs(amountBtns) do
        if (btn ~= nil) then
            CS.UnityEngine.GameObject.Destroy(btn)
        end
    end

    --清除上一次的输入
    Ipt_RechageAmount.text = ''

    Btn_ConfirmOnlinePay.gameObject:SetActive(false)

    Txt_Note.gameObject:SetActive(false)
end

function DisableJiuLiao()
    img_JiuLiaoPay_tag_Close.gameObject:SetActive(true)
    Tog_JiuLiaoPay.interactable = false
end

-- 拦截代理信息
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ProxyRequestRetID, function(data)
    if data.type == "2" then
        --模拟SERVER传来资料是空的,随机指定
        --    local ccc=math.random(70,80)
        --    local x,y=math.modf( ccc/2)
        --    print("",ccc,y)
        --      if(y==0)then
        --	    data=""
        --      end
        ProxyRequestRetCallfun(data)
    end
end, groupName)



function GetProxyRequest(url, callfun)
    ProxyRequestRetCallfun = callfun

    sendata = {};
    sendata.data = url;
    sendata.type = 2;
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ProxyRequestID, sendata)
end

function ShowBubblePrompt(text)
    RecTran_BubblePrompt.gameObject:SetActive(true)

    local fadeTime = 0.2
    local showTime = 1.0
    local objColor

    objColor = Img_BubblePrompt.color
    CS.UnityEngine.Color(objColor.r, objColor.g, objColor.b, 0)
    if sqs ~= nil then sqs:Kill() end
    sqs = CS.DG.Tweening.DOTween.Sequence()
    sqs:Append(Img_BubblePrompt:DOColor(CS.UnityEngine.Color(objColor.r, objColor.g, objColor.g, 1), fadeTime))
    sqs:AppendInterval(showTime)
    sqs:Append(Img_BubblePrompt:DOColor(CS.UnityEngine.Color(objColor.r, objColor.g, objColor.g, 0), fadeTime))

    txt_BubblePrompt.text = text
    objColor = txt_BubblePrompt.color
    CS.UnityEngine.Color(objColor.r, objColor.g, objColor.b, 0)
    if sqs1 ~= nil then sqs1:Kill() end
    sqs1 = CS.DG.Tweening.DOTween.Sequence()
    sqs1:Append(txt_BubblePrompt:DOColor(CS.UnityEngine.Color(objColor.r, objColor.g, objColor.g, 1), fadeTime))
    sqs1:AppendInterval(showTime)
    sqs1:Append(txt_BubblePrompt:DOColor(CS.UnityEngine.Color(objColor.r, objColor.g, objColor.g, 0), fadeTime))
    sqs1:AppendCallback(function()
        RecTran_BubblePrompt.gameObject:SetActive(false)
    end)
end

function Update()
    if RecTran_WebPanel.gameObject.activeSelf == true then
        local wv = RecTran_WebGO.gameObject:GetComponent(typeof(CS.UniWebView))
        if wv == nil then
            EscapeWindow()
        end
    end

    accTime = accTime + CS.UnityEngine.Time.deltaTime
end

function OnWebViewShouldClose(webView)
    return false
end

function OnReceivedMessage(webView, message)
    CS.UILoadingData.Hide()
    if (message.path == "close") then
        CS.UnityEngine.Object.Destroy(webView)
        if (message:GetArgs("message") ~= nil) then
            CS.UIBubblePrompt.Show(message:GetArgs("message"), groupName)
        end
    end
end

function OnLoadComplete(webView, success, errorMessage)
    print("OnLoadComplete ...................................", success, errorMessage)

    if (isloadpayurl == false) then
        isloadpayurl = true
        if (success) then
            CS.UILoadingData.Hide()
            webView:Show()
        else
            CS.UILoadingData.Hide()
            print("加载webview失败:", errorMessage)
            webView:Hide()
            CS.UnityEngine.Object.Destroy(webView)
        end
    end

    -- 顯示錯誤訊息
    if success == false then
        _, q = string.find(errorMessage, "启动第三方app失败")
        -- 開啟APP錯誤
        if q ~= nil then
            -- 第一次為自動開啟失敗
            if isAppError == false then
                isAppError = true
                -- 玩家點擊開啟失敗
            else
                ShowBubblePrompt("启动支付APP失败，请确保您手机上安装了相应的客户端。")
            end
            -- 其他錯誤
        else
            print("OnLoadComplete error:", errorMessage)
        end
    end
end

function IsUrlThirdScheme(url)
    for i = 1, #schemes do
        if (string.find(url, schemes[i])) then
            return true
        end
    end
    return false
end

function AppPause(isPause)
    local wb = RecTran_WebGO:GetComponent("UniWebView")
    if (wb ~= nil and wb.currentUrl) then
        if IsUrlThirdScheme(wb.currentUrl) then
            CloseWebView()
        end
    end

    if isPause then
        pause = true;
    end
end

function AddMoneyAmount(amount)

    local currentMoneyAmount = tonumber(Ipt_RechageAmount.text)
    local _payType, gwid = GetCurrentPayType()
    if currentMoneyAmount == nil then
        currentMoneyAmount = 0
    end
    local _IsFixedAmt = (G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1][8] == 0)
    -- add by bob 當使用固定金額時直接寫入固定金額不做後面運算
    if (_IsFixedAmt) then
        if amount == -1 then
            Ipt_RechageAmount.text = ''
        else
            Ipt_RechageAmount.text = tostring(amount)
        end
        return
    end
    if amount == -1 then
        currentMoneyAmount = 0
    else
        currentMoneyAmount = currentMoneyAmount + tonumber(amount)
    end

    -- add by bob

    local aconfig = G.Tool.TableHelper.Find(apiaconfig, function(v) return v[6] == _payType and v[7] == gwid end)[1];

    if (aconfig ~= nil) then

        min, mint2 = math.modf(aconfig[1] / 100)
        max, maxt2 = math.modf(aconfig[2] / 100)
        local payAmountMin = min
        local payAmountMax = max
        if (currentMoneyAmount > payAmountMax) then
            currentMoneyAmount = payAmountMax
        end
    end

    Ipt_RechageAmount.text = (currentMoneyAmount == 0) and '' or tostring(currentMoneyAmount)
end


-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    if (CS.AppManager.Instance) then
        CS.AppManager.Instance.OnAppPause = CS.AppManager.Instance.OnAppPause - AppPause
    end

    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

function desConfig(configData)
    local aconfig = {}
    local pi = 0
    -- 資料
    local results = nil
    local allpaymentdata = nil

    local myTable = CS.MiniJSON.Json.Deserialize(configData.data)
    if myTable == nil then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        return pi, nil
    end

    local iter = myTable:GetEnumerator()
    while iter:MoveNext() do
        local v = iter.Current
        if v.Key == "data" then
            results = v.Value:GetEnumerator()
        end
    end

    while results:MoveNext() do
        local v = results.Current
        if v.Key == "payAmountInfo" then
            allpaymentdata = v.Value:GetEnumerator()
        end
    end

    -- #region insert payment
    while allpaymentdata:MoveNext() do
        local v = allpaymentdata.Current
        local paymentdata
        paymentdata = v:GetEnumerator()
        local payid
        local isAdd = true
        local checkLimit = true
        if G.Platform.GetBranch() == 1 then
            checkLimit = false
        elseif G.Platform.GetBranch() == 2 then
            checkLimit = false
        end
        while paymentdata:MoveNext() do
            local p = paymentdata.Current
            if p.Key == "payType" and aconfig[pi] and not aconfig[pi][5] then
                payid = tostring(p.Value)
                aconfig[pi][5] = payid;
            end
            if p.Key == "minAmount" and aconfig[pi] and not aconfig[pi][1] then
                aconfig[pi][1] = p.Value;
            end
            if p.Key == "maxAmount" and aconfig[pi] and not aconfig[pi][2] then
                aconfig[pi][2] = p.Value;
            end
            if p.Key == "status" and aconfig[pi] and not aconfig[pi][4] then
                aconfig[pi][4] = p.Value;
            end
            if p.Key == "amounts" and aconfig[pi] and not aconfig[pi][3] then
                aconfig[pi][3] = {};
                aconfig[pi][3] = G.Tool.StringHelper.Split(p.Value, ",");
            end

            if p.Key == "typeID" then

                if isAdd == true then
                    pi = pi + 1
                    if (aconfig[pi] == nil) then
                        aconfig[pi] = {};
                        aconfig[pi][6] = p.Value;
                    end
                end
            end
            if p.Key == "gwid" and aconfig[pi] and not aconfig[pi][7] then
                aconfig[pi][7] = p.Value;
            end
            if p.Key == "isRange" and aconfig[pi] and not aconfig[pi][8] then
                aconfig[pi][8] = p.Value;
            end
        end
    end

    return pi, aconfig
end

local isInitPayConfig = false
local tempApiaconfig = {}
function InitPayConfig()
    -- 第二次直接拿暫存的資料
    if isInitPayConfig == true then
        apiaconfig = tempApiaconfig
        return
    end

    data = G.Service.ReCharge.GetPayConfig()
    if (data == "" or data == nil) then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        return
    end

    if (data.data == nil or data.data == "") then
        CS.UILoadingData.Show("支付渠道忙碌中请稍后再试...", 5)
        return
    end

    local pi = 0
    pi, apiaconfig = desConfig(data)

    --- endregion
    if (pi == 0) then
        CS.UILoadingData.Show("无支付清单,请稍后再试...", 1)
        return
    end

    -- 紀錄第一次的資料
    isInitPayConfig = true
    tempApiaconfig = apiaconfig

    G.Tool.TableHelper.Print("apiaconfig", apiaconfig)
end

function GetNowPayConfig()
    local dataNow = G.Service.ReCharge.GetPayConfig()
    if (dataNow == "" or dataNow == nil) then
        return nil
    end

    if (dataNow.data == nil or dataNow.data == "") then
        return nil
    end

    local pi = 0
    pi, apiaconfigNow = desConfig(dataNow)

    --- endregion
    if (pi == 0) then
        return nil
    end

    G.Tool.TableHelper.Print("now apiaconfig", apiaconfigNow)

    return apiaconfigNow
end


-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)

    InitPayConfig()

    --- -預設加仔原本內容
    local leftMenu = {}

    -- 設置資料
    local setingData = {
        {}, -- IOS
        { idx = 1, isshow = G.Service.UIFlag.IsShowOnlinePay(), detailgb = RecTran_OnlinePayDetail }, -- 在線支付
        { idx = 2, isshow = G.Service.UIFlag.IsShowAgentPay(), gb = Tog_AgentPay, detailgb = RecTran_AgentPayDetail }, -- 代理支付
        { idx = 3, isshow = G.Service.UIFlag.IsShowAgentApply(), gb = Tog_AgentRecruit, detailgb = RecTran_AgentRecuit }, -- 代理招商
        { idx = 4, isshow = G.Service.UIFlag.IsShowJiuLiao(), gb = Tog_JiuLiaoPay, detailgb = RecTran_OnlinePayDetail }, -- 久聊支付
        { idx = 5, isshow = G.Service.UIFlag.IsShowShanFu(), gb = Tog_Henry, detailgb = RecTran_HenryDetail }, -- 閃付支付
        { idx = 6, isshow = G.Service.UIFlag.IsShow888(), gb = Tog_888Pay, detailgb = RecTran_OnlinePayDetail }, -- 888支付
        { idx = 7, isshow = G.Service.UIFlag.IsShowCF(), gb = Tog_CFPay, detailgb = RecTran_OnlinePayDetail }, -- CF支付
        { idx = 8, isshow = G.Service.UIFlag.IsShowBaoHe(), gb = Tog_BaoHePay, detailgb = RecTran_OnlinePayDetail }, -- 宝盒支付
        { idx = 9, isshow = G.Service.UIFlag.IsShowDongGuaZF(), gb = Tog_DongGuaZFPay, detailgb = RecTran_OnlinePayDetail }, -- 冬瓜支付保支付
        { idx = 10, isshow = G.Service.UIFlag.IsShowDongGuaWechat(), gb = Tog_DongGuaWeChatPay, detailgb = RecTran_OnlinePayDetail }, -- 冬瓜維信支付
    }

    RecTran_OnlinePayDetail.gameObject:SetActive(false)
    RecTran_AgentPayDetail.gameObject:SetActive(false)
    RecTran_AgentRecuit.gameObject:SetActive(false)

    --关闭久聊
    if ((#G.Service.ReCharge.GetPaymentInfoList(1) == 0)) then
        DisableJiuLiao();
    end
    -- 排序		
    local order = G.User.GetCurrentUser().UserData.UISort
    if order == nil or order == "" then
        print("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! no UISort data")

        order = { 1, 3, 2, 4, 5, 6, 7, 8, 9, 10, 11 }
    else
        order = {}
        value_1 = G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.UISort, 1, 1)
        value_2 = G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.UISort, 3, 1)
        value_3 = G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.UISort, 5, 1)
        value_4 = G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.UISort, 7, 1)

        table.insert(order, value_1) -- 不可以動        
        table.insert(order, 6)
        table.insert(order, 7)
        table.insert(order, 8)
        table.insert(order, 9)
        table.insert(order, 10)
        table.insert(order, 11)
        table.insert(order, value_2)
        table.insert(order, value_3)
        table.insert(order, value_4)
        table.insert(order, 5)
    end

    local isFirst = true
    for i = 2, #order, 1 do
        menu = setingData[tonumber(order[i])]

        --將在線支付拆出到左列
        if menu.idx == 1 then
            if (menu.isshow == true) then
                --支付类型清单：loop -- typeid
                for _, data in pairs(onlineSetingData) do
                    if data.gb ~= nil then
                        gb = data.gb

                        local isFind = false

                        if tempApiaconfig ~= nil then
                            --支付总清单涵盖重复的支付宝：loop -- payType
                            G.Tool.TableHelper.Print(tempApiaconfig)
                            for _, configData in pairs(tempApiaconfig) do

                                if tostring(configData[6]) == tostring(data.index) then
                                    -- 
                                    isFind = true
                                    break
                                end
                            end
                        end


                        if isFind == true then
                            gb.gameObject:SetActive(true)
                            gb.gameObject.transform:SetParent(this.transform)
                            gb.gameObject.transform:SetParent(RecTran_ToggleGroup)
                            local pos = gb.gameObject.transform.localPosition
                            pos.z = -1
                            gb.gameObject.transform.localPosition = pos

                            if isFirst == true then
                                data.gb.isOn = true
                                data.gb.gameObject.transform.parent:GetComponent("ToggleGroup"):NotifyToggleOn(data.gb)

                                isFirst = false
                                menu.detailgb.gameObject:SetActive(true)
                                -- 随机抽选机制（将在点选时做现有清单随机抽选送出）
                                local paylist = G.Tool.TableHelper.Find(tempApiaconfig, function(v) return tostring(v[6]) == tostring(data.index) end);
                                local len = tablelength(paylist);
                                --抽选号码                             
                                local assignnum = math.random(1, len);
                                --设定支付方式（PaytypeID）与厂商ID（GWID）                              
                                SetCurrentPayType(data.index, paylist[assignnum][7])
                                payConfig = G.Service.ReCharge.GetPayConfig()
                                ProcessAmtList(payConfig)
                            end
                        else
                            gb.gameObject:SetActive(false)
                        end
                    end
                end
            else
                for _, data in pairs(onlineSetingData) do
                    if data.gb ~= nil then
                        data.gb.gameObject:SetActive(false)
                    end
                end
            end
        else
            menu.gb.gameObject.transform:SetParent(this.transform)
            menu.gb.gameObject.transform:SetParent(RecTran_ToggleGroup)
            local pos = menu.gb.gameObject.transform.localPosition
            pos.z = -1
            menu.gb.gameObject.transform.localPosition = pos

            if (menu.isshow == true) then
                -- print("idx:"..menu.idx,"true")
                menu.gb.gameObject:SetActive(true)
                if (isFirst) then
                    menu.gb.isOn = true
                    menu.gb.gameObject.transform.parent:GetComponent("ToggleGroup"):NotifyToggleOn(menu.gb)
                    menu.detailgb.gameObject:SetActive(true)
                    isFirst = false

                    if menu.idx == 4 then
                        SetCurrentPayType(101)
                        apiaconfig = G.Service.ReCharge.GetPaymentInfoList(1);
                        ProcessPayList()
                    elseif menu.idx == 5 then
                        RecTran_HenryDetail.gameObject:SetActive(true)
                        ClickHenry()
                    elseif menu.idx == 6 then
                        SetCurrentPayType(102)
                        apiaconfig = G.Service.ReCharge.Get888PayInfo(1);
                        ProcessPayList()
                    elseif menu.idx == 7 then
                        SetCurrentPayType(103)
                        apiaconfig = G.Service.ReCharge.GetCFPayInfo(1);
                        ProcessPayList()
                    elseif menu.idx == 8 then
                        SetCurrentPayType(104)
                        apiaconfig = G.Service.ReCharge.GetBaoHePayInfo(1);
                        ProcessPayList()
                    elseif menu.idx == 9 then
                        SetCurrentPayType(18, 19)
                        apiaconfig = G.Service.ReCharge.GetDongGuaZFPayInfo(1);
                        ProcessPayList()
                    elseif menu.idx == 10 then
                        SetCurrentPayType(19, 19)
                        apiaconfig = G.Service.ReCharge.GetDongGuaWeChatPayInfo(1);
                        ProcessPayList()
                    end
                end
            else
                print("idx:" .. menu.idx, "false")
                menu.gb.gameObject:SetActive(false)
            end
        end
    end

    Txt_PlayerID.text = tostring(user.UserID)
    Txt_CashAmount.text = G.Tool.UserHelper.ToCashDetail(userdata.CashAmount)
    Txt_SafeAmount.text = G.Tool.UserHelper.ToCashDetail(userdata.BankAmount)
    RefreshAgentList(false)


    --久聊初始化
    if CS.JiuLiaoSDK.s_IsJiuLiaoSDK == true then
        if G.GameDefine.JiuLiaoId ~= "" and G.GameDefine.JiuLiaoSecrety ~= "" then
            CS.JiuLiaoSDK.init(G.GameDefine.JiuLiaoId, G.GameDefine.JiuLiaoSecrety)
        end
    end
end



function EscapeWindow()
    if RecTran_WebPanel.gameObject.activeSelf then
        CloseWebView()
    elseif RecTran_reportAgent.gameObject.activeSelf then
        RecTran_reportAgent.gameObject:SetActive(false)
    elseif RecTran_ContactAgent.gameObject.activeSelf then
        RecTran_ContactAgent.gameObject:SetActive(false)
    else
        CS.WindowManager.Instance:CloseWindow("UIRecharge", false)
    end
end


function tablelength(T)
    local count = 0
    for _ in pairs(T) do count = count + 1 end
    return count
end


-- 获取代理回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentListRetID, function(receiveMsg)
    RecTran_AgentList.transform:DestoryAllChildren(false)
    G.Agent.SetAgentInfo(receiveMsg.AgentList)
    CreateAgentList()
end, groupName)

-- 代理招商申请回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ApplyAgentRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("申请成功", groupName)
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName)
    end
end, groupName)


-- 举报代理
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ReportAgentRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("举报成功", groupName)
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName)
    end
end, groupName)


local sqs
local sqs1



-- 创建订单返回
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.CreatePayOrderRetID, function(receiveMsg)
    if receiveMsg.code == 1 then
        if not G.Env.IsWindows() then
            if (LastPayType ~= 6) then
                local url = string.format(G.GameConfig.ChargeDomain, receiveMsg.OrderID, LastPayGWID)
                print("--------------- Normal url", url, G.GameConfig.ChargeDomain)

                -- 透过server取得网页结果

                local _payType = tostring(GetCurrentPayType())
                local IsRedirect = (G.Tool.TableHelper.Find(redirecpayment, function(p) return tostring(p) == _payType; end)[1] ~= nil)
                local IsGoShield = (G.Tool.TableHelper.Find(GoShield, function(p) return tostring(p) == _payType; end)[1] ~= nil)
                print("IsRedirect", IsRedirect, _payType)
                print("IsGoShield", IsGoShield, _payType)
                if (IsRedirect) then
                    GetProxyRequest(url, ProcessOrderRedirctRet)
                    --测试用资料,模拟金流来源程式码
                    --  GetProxyRequest("http://0.0.0.0:3333/ralipay.html", ProcessOrderRedirctRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/mapi.html", ProcessOrderRedirctRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/unipay.html", ProcessOrderRedirctRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/wechatfix.html", ProcessOrderRedirctRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/pay02.html", ProcessOrderRedirctRet)
                elseif (IsGoShield and (CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android)) then
                    print("---------------Go 盾")
                    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
                        -- 使用超级盾
                        if G.Platform.getUseQsu() == true then
                            if G.Tool.Qsu.IsQsuSDK() == true then
                                local res, ip, port = G.Tool.Qsu.getProxyTcpByDomain(G.Tool.Qsu.GetPay(), 9090)
                                if res == 0 then
                                    CS.UILoadingData.Hide()
                                    local url = "http://" .. ip .. ":" .. port .. "/paycenter/index?orderid=" .. receiveMsg.OrderID
                                    print("---------------Qsu url", url)
                                    CS.UnityEngine.Application.OpenURL(url)
                                    return
                                end
                            end
                            -- 使用游戏盾
                        elseif G.Tool.Aliyun.IsAliyunSDK() == true then
                            local res, ip, port = G.Tool.Aliyun.getProxyTcpByDomain("pay", 9090)
                            if res == 0 then
                                CS.UILoadingData.Hide()
                                local url = "http://" .. ip .. ":" .. port .. "/paycenter/index?orderid=" .. receiveMsg.OrderID
                                print("---------------Aliyun url", url)
                                CS.UnityEngine.Application.OpenURL(url)
                                return
                            end
                        end
                    end
                else
                    GetProxyRequest(url, ProcessOrderRet)
                    --测试用资料,模拟金流来源程式码
                    --GetProxyRequest("http://0.0.0.0:3333/ralipay.html", ProcessOrderRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/mapi.html", ProcessOrderRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/unipay.html", ProcessOrderRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/wechatfix.html", ProcessOrderRet)
                    --  GetProxyRequest("http://0.0.0.0:3333/pay02.html", ProcessOrderRet)
                end

                -- DO OPEN 久聊              
            else
                print(tostring(user.UserID) .. ":" .. tostring(LastAmount))
                CS.UILoadingData.Hide()

                if CS.JiuLiaoSDK.s_IsJiuLiaoSDK == true then
                    if G.GameDefine.JiuLiaoName ~= "" then
                        if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
                            CS.JiuLiaoSDK.shareBankTransfer(G.GameDefine.JiuLiaoName, user.UserID, (LastAmount * 100), receiveMsg.OrderID, "游戏币充值")
                            -- IOS版参数对调
                        elseif CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
                            CS.JiuLiaoSDK.shareBankTransfer(user.UserID, G.GameDefine.JiuLiaoName, (LastAmount * 100), receiveMsg.OrderID, "游戏币充值")
                        end
                    end
                end
            end
        else
            CS.UILoadingData.Hide()
            CS.UIBubblePrompt.Show("只有移动客户端支持支付功能.", groupName)
        end
    else
        CS.UILoadingData.Hide()
        CS.UIBubblePrompt.Show(receiveMsg.message, groupName)
    end
end, groupName)



-- *********************** 閃付 ***********************
function ClickHenry()
    if accTime - lastGetTime < 10 then
        return
    end
    lastGetTime = accTime

    if RecTran_HenryList.transform.childCount == 0 then
        Txt_HenryNoDataTips.gameObject:SetActive(true)
        Txt_HenryNoDataTips.text = "读取中....,请稍后..."
    end

    local param = {}
    param.UserID = G.User.GetCurrentUserID()
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetShanFuID, param)
end


--[[
支付宝	1	-
微信	2	-
银行卡	3	-
支付宝红包	4	支付宝
支付宝转账	5	支付宝
微信红包	6	微信
微信转账	7	微信
大额专属通道	8	-
微转银二维码	9	大额专属通道
支转银二维码	10	大额专属通道
云闪付二维码	11	大额专属通道
银行卡充值	12	银行卡
]] --

-- 查詢閃付返回
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetShanFuRetID, function(receiveMsg)
    CS.UILoadingData.Hide()
    print(receiveMsg.CodeID)
    print(receiveMsg.Data)

    RecTran_HenryList.transform:DestoryAllChildren(true)

    Txt_HenryNoDataTips.text = "暂无资料"

    local agentsDara = Json.decode(receiveMsg.Data)
    if agentsDara.code ~= 200 then
        Txt_HenryNoDataTips.gameObject:SetActive(true)
        return
    end

    if #agentsDara.data > 0 then
        Txt_HenryNoDataTips.gameObject:SetActive(false)
    else
        Txt_HenryNoDataTips.gameObject:SetActive(true)
    end
    i = 0
    for _, data1 in pairs(agentsDara.data) do
        local henryItem = CS.UnityEngine.GameObject.Instantiate(RecTran_HenryItem.gameObject)
        henryItem.transform:SetParent(RecTran_HenryList.transform)
        henryItem.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
        henryItem.transform.localPosition = CS.UnityEngine.Vector3.zero
        henryItem.transform.gameObject:SetActive(true)

        henryItem.transform:Find("FakeFace/Face" .. i).gameObject:SetActive(true)

        -- 名字
        local txtName = henryItem.transform:Find("Name"):GetComponent("Text")
        txtName.text = data1.shopName

        -- 評價
        local txtEvaluation = henryItem.transform:Find("EvaluationValue"):GetComponent("Text")
        txtEvaluation.text = tostring(data1.evaluateNum) .. "+"

        -- 星星數
        local txtStar = henryItem.transform:Find("StarValue"):GetComponent("Text")
        local value = tonumber(data1.evaluate)
        txtStar.text = tostring(value / 10)
        t1, t2 = math.modf(value / 10)
        for j = 1, t1 do
            henryItem.transform:Find("Star/Star" .. tostring(j) .. "/Star"):GetComponent("Image").fillAmount = 1
        end
        if t2 > 0 then
            henryItem.transform:Find("Star/Star" .. tostring(t1 + 1) .. "/Star"):GetComponent("Image").fillAmount = t2
        end

        -- 顯示支援充值類型
        local agentsDara1 = Json.decode(data1.payTypes)
        for index, tpyeData in pairs(agentsDara1) do
            henryItem.transform:Find("IconList/Grid/Image" .. index).gameObject:SetActive(true)
        end

        local Btn_Recharge = henryItem:GetComponent("Button")
        local agentId = data1.chatUserId
        Btn_Recharge:AddClickListener(function()
            PlayCellPhoneAni(agentId)
        end)
        i = i + 1
    end
end, groupName)

local cellPhone = nil
function PlayCellPhoneAni(agentId)
    RecTran_CellPhone.gameObject:SetActive(true)

    cellPhone = CS.UnityEngine.GameObject.Instantiate(CellPhoneAni.gameObject)
    cellPhone.transform:SetParent(CellPhoneAni.parent.transform)
    cellPhone.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
    cellPhone.transform.localPosition = CS.UnityEngine.Vector3.zero
    cellPhone.transform.gameObject:SetActive(true)

    CS.TimeManager.Instance:DelayInvoke(1.5,
        function()
            user = G.User.GetCurrentUser();

            local param = {}
            param.AgentID = tostring(agentId)
            param.AppID = tostring(G.GameConfig.GamePlatformID)
            param.DeviceID = G.Env.GetHDCode()
            param.UserID = tostring(user.UserID)
            param.UserName = tostring(user.UserID)

            if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer
                    or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXEditor
                    or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.OSXPlayer then
                param.Type = 1
            elseif CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
                param.Type = 2
            elseif CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor
                    or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsPlayer then
                param.Type = 2
            end

            G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetShanFuChatID, param)
        end)

    -- 避免發生錯誤時一直沒關起來
    CS.TimeManager.Instance:DelayInvoke(10,
        function()
            RecTran_CellPhone.gameObject:SetActive(false)
            if cellPhone ~= nil then
                CS.UnityEngine.GameObject.Destroy(cellPhone)
                cellPhone = nil
            end
        end)
end


-- 閃付加密返回
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetShanFuChatRetID, function(receiveMsg)
    RecTran_CellPhone.gameObject:SetActive(false)
    if cellPhone ~= nil then
        CS.UnityEngine.GameObject.Destroy(cellPhone)
        cellPhone = nil
    end

    G.Tool.Henry.showMessageChat(receiveMsg.ClientUrl,
        receiveMsg.AppKey,
        receiveMsg.AppID,
        receiveMsg.UserID,
        receiveMsg.UserName,
        receiveMsg.DeviceID,
        receiveMsg.TimeStamp,
        tonumber(receiveMsg.AgentID),
        receiveMsg.SecretID,
        "0")
end, groupName) 

