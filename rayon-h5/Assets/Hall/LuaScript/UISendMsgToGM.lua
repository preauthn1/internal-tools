--local Btn_SendMsg = this.transform:Find("UIBody/_SendMsg"):GetComponent("Button")
--local Ipt_Content = this.transform:Find("UIBody/_Content"):GetComponent("InputField")
--local Btn_CloseSendMsgTOGM = this.transform:Find("UIBody/_CloseSendMsgTOGM"):GetComponent("Button")
--local Txt_ContentLength = this.transform:Find("UIBody/_ContentLength"):GetComponent("Text")
local RecTran_Mask = this.transform:Find("Mask"):GetComponent("RectTransform")

local btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local imgTitle = this.transform:Find("UIBody/Background/ImageTitle"):GetComponent("Image")

local btnService = this.transform:Find("UIBody/SubList/Grid/ButtonService"):GetComponent("Button")
local imgServiceSelect = this.transform:Find("UIBody/SubList/Grid/ButtonService/ImageSelect"):GetComponent("Image")
local imgServiceTitle = this.transform:Find("UIBody/SubList/Grid/ButtonService/ImageTitle"):GetComponent("Image")
local btnService = this.transform:Find("UIBody/SubList/Grid/ButtonService"):GetComponent("Button")
local btnSendMailToGM = this.transform:Find("UIBody/SubList/ButtonSendMailToGM"):GetComponent("Button")

local RecTran_DetailService = this.transform:Find("UIBody/DetailService"):GetComponent("RectTransform")
local RecTran_WebGO = this.transform:Find("UIBody/DetailService/_WebGO"):GetComponent("RectTransform")

local btnHelp = this.transform:Find("UIBody/SubList/Grid/ButtonHelp"):GetComponent("Button")
local imgHelpSelect = this.transform:Find("UIBody/SubList/Grid/ButtonHelp/ImageSelect"):GetComponent("Image")
local imgHelpTitle = this.transform:Find("UIBody/SubList/Grid/ButtonHelp/ImageTitle"):GetComponent("Image")

local RecTran_DetailHelp = this.transform:Find("UIBody/DetailHelp"):GetComponent("RectTransform")
local txtHelpContent = this.transform:Find("UIBody/DetailHelp/Content"):GetComponent("Text")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---
local G = require "Global";
local gameDefine = require("Com/GameDefine")

local userdata = G.User.GetCurrentUser().UserData
local groupName = "UISendMsgToGM";
local contentMax = 256
local isloadpayurl = false          -- 是否已经载入过充值页面，每次重新充值是设置为false
local UniWebView
local SubIndex = 1 --分页索引

-- 初始化函数，实例化之后调用一次
function Start()
	btnClose:AddClickListener(function() EscapeWindow() end)
	
	btnService:AddClickListener(function() ChangeSub(1) end)
	btnHelp:AddClickListener(function() ChangeSub(2) end)
	
	btnSendMailToGM:AddClickListener(function() 
		 EscapeWindow()
		local openparam = CS.WindowNodeInitParam("UISendMailToGM", CS.AppDefine.HallNumber)
        CS.WindowManager.Instance:OpenWindow(openparam)
	end)
	
	ChangeSub(1)
end

function RefreshWindowData(windowData)
	LoadService()
	LoadHelp()
end
 
-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end  

function EscapeWindow()
    G.Service.CS.RefreshCSList();
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 切换分页
function ChangeSub(index)
	SubIndex = index
	imgServiceSelect.gameObject:SetActive(index == 1)
	RecTran_DetailService.gameObject:SetActive(index == 1)
	
	imgHelpSelect.gameObject:SetActive(index == 2)
	RecTran_DetailHelp.gameObject:SetActive(index == 2)
	
	if UniWebView ~= nil and index == 1 then
		UniWebView:Show()
	else
		if UniWebView ~= nil then
			UniWebView:Hide()
		end
	end
	
	if index == 1 then
		ShowService()
	else
		ShowHelp()
	end
end

-- 显示客服
function ShowService()
	imgTitle.sprite = imgServiceTitle.sprite
end

-- 显示帮助
function ShowHelp()
	imgTitle.sprite = imgHelpTitle.sprite
end
-- 载入客服
function LoadService()
	local user = G.User.GetLocalUserInfo()
	-- edit下固定帳號登入取出來的值會是nil
	if user == nil then
		user = G.User.GetCurrentUser()
	end
    local url="";
    local hd=0; 
    if( G.Env.IsAndroid()) then
      hd=2;
    elseif(G.Env.IsIOS()) then
      hd=1;
    end

	local CustomServiceDomain
	user1 = G.User.GetCurrentUser()
	if user1.FlowInfo.CustomerDomain == nil or user1.FlowInfo.CustomerDomain == "" then
		CustomServiceDomain = G.GameConfig.CustomServiceDomain
	else
		CustomServiceDomain = user1.FlowInfo.CustomerDomain
	end
	--986就无须串入参数
	if G.Platform.GetBranch() == 1 then
		url = CustomServiceDomain;
    else
		if(G.Tool.StringHelper.IsUpgrade_v1(G.Env.GetHallVersion(),hd)) then
			url = CustomServiceDomain.."/ServiceTalk?isUploadFile=true&cerCode="..user.Cer;
		else
			url = CustomServiceDomain.."/ServiceTalk?cerCode="..user.Cer;
		end
	end	
    print(url)
	
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor then
	
	else
		UniWebView = RecTran_WebGO:GetComponent("UniWebView")
			
		if (UniWebView == nil) then
			UniWebView = RecTran_WebGO.gameObject:AddComponent(typeof(CS.UniWebView))
		end
	end
	
	if UniWebView ~= nil then
		UniWebView:CleanCache()
	end
	
	-- 4个参数 分别为 距离上边距离，左边距离，下边距离，右边距离
	local screenWidth = CS.UniWebViewHelper.screenWidth * 0.5
	local screenHeight = CS.UniWebViewHelper.screenHeight * 0.5
	local rectWidth = RecTran_Mask.rect.width * 0.5
	local rectHeight = RecTran_Mask.rect.height * 0.5
	local scale = screenHeight/rectHeight
	top = math.modf((rectHeight - (RecTran_DetailService.rect.height * 0.5 + RecTran_DetailService.anchoredPosition.y))*scale)
	bottom = math.modf((rectHeight - (RecTran_DetailService.rect.height * 0.5 - RecTran_DetailService.anchoredPosition.y))*scale)
	
	scale = screenWidth/rectWidth
	left = math.modf((rectWidth - (RecTran_DetailService.rect.width * 0.5 - RecTran_DetailService.anchoredPosition.x))*scale)
	right = math.modf((rectWidth - (RecTran_DetailService.rect.width * 0.5 + RecTran_DetailService.anchoredPosition.x))*scale)
    -- 只有安卓手機才做內縮
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
		left = left + 80
		right = right + 80
    end	
	
	print(top,left,bottom,right)
	if UniWebView ~= nil then
		UniWebView.insets = CS.UniWebViewEdgeInsets(top,left,bottom,right)
		UniWebView:OnLoadComplete('+', OnLoadComplete)
		UniWebView:OnReceivedMessage('+', OnReceivedMessage)
		-- 按下返回鍵不關掉webview，\Assets\FrameWork\XLuaWebViewExport中加上typeof(UniWebView.WebViewShouldCloseDelegate)
		--UniWebView:OnWebViewShouldClose('+', OnWebViewShouldClose)

		-- 载入页面时不显示loading
		UniWebView:SetShowSpinnerWhenLoading(false)
		UniWebView.toolBarShow = false
		UniWebView.url = url

		UniWebView:Load()
	end
end

function Update()
    if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WindowsEditor then
		return
	end
	
	if RecTran_WebGO.gameObject.activeSelf == true then
		local wv = RecTran_WebGO.gameObject:GetComponent(typeof(CS.UniWebView))
		if wv == nil then
			EscapeWindow()
		end
	end
end

-- 载入帮助
function LoadHelp()
	G.Service.FAQ.GetFAQ(function(content)
        txtHelpContent.text = content
    end)
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


    if (isloadpayurl == false) then
        isloadpayurl = true
      
        if (success) then
            CS.UILoadingData.Hide()
			if SubIndex == 1 then
				webView:Show()
			end
        else
            CS.UILoadingData.Hide()
            print("加载webview失败:", errorMessage)
            webView:Hide()
            CS.UnityEngine.Object.Destroy(webView)

            CS.UIBubblePrompt.Show("启动客服失败。", groupName)
        end
    end
end