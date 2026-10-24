require("Module")
mmodule(...)
local G = require "Global";
local groupName = "UIWXD";

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local user = G.User.GetCurrentUser();
local userdata = user.UserData
local subList = {}
subList[1] = {}
subList[1].myName = "Promotion"
subList[2] = {}
subList[2].myName = "Team"
subList[3] = {}
subList[3].myName = "Result"
subList[4] = {}
subList[4].myName = "Reward"
subList[5] = {}
subList[5].myName = "Cash"
subList[6] = {}
subList[6].myName = "Loan"
subList[7] = {}
subList[7].myName = "Help"

local mQRCodeUrl
local mUpID = ""

local function CreateClone(prefab)
	local clone = CS.UnityEngine.GameObject.Instantiate(prefab)
	if(clone == nil)then return nil end
	clone.transform:SetParent(prefab.transform.parent)
	clone.transform.localPosition = CS.UnityEngine.Vector3.zero
	clone.transform.localScale = CS.UnityEngine.Vector3.one
	
	return clone	
end

local function SetParentTopLayer(object)
	object.transform:SetParent(rootTopLayer.transform)
	object.transform.localPosition = CS.UnityEngine.Vector3.zero
	object.transform.localScale = CS.UnityEngine.Vector3.one
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	Init(this,this)
end

-- 初始介面資料
function Init(root,topLayer)

	rootTopLayer = topLayer
	animopen = root.transform:Find("UIBody"):GetComponent("Animation")
	
	btnClose = root.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
    btnClose:AddClickListener(function() EscapeWindow() end)
	
	for i=1, #subList do
		local tempIndex = i
		local myName = subList[i].myName
		
		subList[i].rootDetail = root.transform:Find("UIBody/Detail" .. myName):GetComponent("RectTransform")
		subList[i].rootDetail.gameObject:SetActive(false)
		subList[i].btn = root.transform:Find("UIBody/SubList/Grid/Button" .. myName):GetComponent("Button")
		subList[i].btn:AddClickListener(function() ChangeSubList(tempIndex) end)
		subList[i].imgSelect = root.transform:Find("UIBody/SubList/Grid/Button" .. myName .. "/ImageSelect"):GetComponent("Image")
		subList[i].imgSelect.gameObject:SetActive(false)
	end
	
	ChangeSubList(1)
	
	txtUserID = root.transform:Find("UIBody/TextUserID"):GetComponent("Text")
	txtUserID.text = "ID:" .. tostring(user.UserID)
	
	txtDate = root.transform:Find("UIBody/TextDate"):GetComponent("Text")
	
	-- 我的推廣 Promotion
	InitMyPromotion(root)
	-- 我的玩家
	InitAgentGetWithAmount(root)
	-- 我的業績
	InitUserAgentResult(root)
	-- 我的獎勵
	InitUserAgentWinAmount(root)
	-- 我的提現
	InitUserAgentGetWithdrawAmount(root)
	-- 我的預提
	InitUserAgentBorrow(root)
	-- 獎勵說明
	InitHelp(root)
end
function SetParentRoot(t)
	ParentRoot=t;
end
	
-- 初始化函数，实例化之后调用一次
function Start()
	Show(groupName)
end
	
-- 顯示介面
function Show(rootName)
	if(groupName ~= rootName)then
		animopen.enabled = false
	end
	
	-- 我的推廣
	SendMyPromotion()
	-- 我的玩家
	SendAgentGetWithAmount()
	-- 我的業績
	SendUserAgentResult()
	-- 我的獎勵
	SendUserAgentWinAmount()
	-- 我的提現
	SendUserAgentGetWithdrawAmount()
	-- 我的預提
	SendUserAgentBorrow()
	-- 獎勵說明
	SendRewardTips()
end
 
-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

-- 離開介面
function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 切換介面
function ChangeSubList(index)
	for i=1, #subList do
		subList[i].rootDetail.gameObject:SetActive(index == i)
		subList[i].imgSelect.gameObject:SetActive(index == i)
	end
end


-- *************************** 我的推廣 ***************************
function InitMyPromotion(root)
	imgQRCode = root.transform:Find("UIBody/DetailPromotion/ImageQRCode/Image"):GetComponent("Image")
	txtPromotionLink = root.transform:Find("UIBody/DetailPromotion/TextLink"):GetComponent("Text")
	txtPromotionUpId = root.transform:Find("UIBody/DetailPromotion/TextUpId"):GetComponent("Text")
	txtPromotionUserID = root.transform:Find("UIBody/DetailPromotion/TextUserID"):GetComponent("Text")
	txtPromotionGroupNum = root.transform:Find("UIBody/DetailPromotion/TextGroupNum"):GetComponent("Text")
	txtPromotionDirectly = root.transform:Find("UIBody/DetailPromotion/TextDirectly"):GetComponent("Text")
	txtPromotionNew = root.transform:Find("UIBody/DetailPromotion/TextNew"):GetComponent("Text")
	txtPromotionActive = root.transform:Find("UIBody/DetailPromotion/TextActive"):GetComponent("Text")
	
	btnPromotionCopy = root.transform:Find("UIBody/DetailPromotion/BtnCopy"):GetComponent("Button")
	btnPromotionCopy:AddClickListener(function()
        CS.Utility.CopyText(mQRCodeUrl)
        CS.UIBubblePrompt.Show("复制成功", groupName)
	end)
	btnPromotionCopy.gameObject:SetActive(false)
	
	btnPromotionShare = root.transform:Find("UIBody/DetailPromotion/BtnShare"):GetComponent("Button")
	btnPromotionShare:AddClickListener(function()
        CS.Utility.CopyText(mQRCodeUrl)
        CS.UIBubblePrompt.Show("已复制资讯,3秒后自动打开微信", groupName)
        local sequence = CS.DG.Tweening.DOTween.Sequence()
        sequence:AppendInterval(2)
        sequence:AppendCallback( function()
            CS.UnityEngine.Application.OpenURL("weixin://")
        end )
        sequence:AppendInterval(2)
        sequence:AppendCallback( function()
             if(pause) then
                pause=false;
             else   
                CS.UnityEngine.Application.OpenURL("https://example.invalid/")
             end
        end )
	end)
	btnPromotionShare.gameObject:SetActive(false)
	
	--邀请码
	txtInviteCode = root.transform:Find("UIBody/DetailPromotion/InviteCode/Code"):GetComponent("Text")
	btnInviteCodeCopy = root.transform:Find("UIBody/DetailPromotion/InviteCode/BtnCopy"):GetComponent("Button")
	btnInviteCodeCopy:AddClickListener(function() 
		CS.Utility.CopyText(txtInviteCode.text)
		CS.UIBubblePrompt.Show("复制成功", groupName)
	end)
	--绑定
	rootSendInviteCode = root.transform:Find("UIBody/DetailPromotion/SendInviteCode"):GetComponent("RectTransform")
	ipSendInviteCode = root.transform:Find("UIBody/DetailPromotion/SendInviteCode/InputFieldCode"):GetComponent("InputField")
	btnSendInviteCode = root.transform:Find("UIBody/DetailPromotion/SendInviteCode/BtnSend"):GetComponent("Button")
	btnSendInviteCode:AddClickListener(function()
		SendSplitBindingInvitedCodeID()
	end)
	
	rootQRCodeExpand = root.transform:Find("RootQRCodeExpand"):GetComponent("RectTransform")
	rootQRCodeExpand.gameObject:SetActive(false)
	imgQRCodeExpandBg = root.transform:Find("RootQRCodeExpand/ImageBg"):GetComponent("RectTransform")
	imgQRCodeExpand = root.transform:Find("RootQRCodeExpand/ImageQRCode"):GetComponent("Image")
	
	btnQRCodeExpand = root.transform:Find("UIBody/DetailPromotion/ImageQRCode/ButtonExpand"):GetComponent("Button")
	btnQRCodeExpand:AddClickListener(function()	
	
	print("1111111")
	rootQRCodeExpand.gameObject:SetActive(true) end)
	btnCloseQRCodeExpand = root.transform:Find("RootQRCodeExpand"):GetComponent("Button")
	btnCloseQRCodeExpand:AddClickListener(function()	rootQRCodeExpand.gameObject:SetActive(false) end)
	
	SetParentTopLayer(rootQRCodeExpand)
end
-- 查詢
function SendMyPromotion()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetMyPromotionID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetMyPromotionRetID, function(receiveMsg)
	if(receiveMsg.QrURL ~= "")then		
		mQRCodeUrl = receiveMsg.QrURL
		local sprite = CS.Utility.GenQRCode(mQRCodeUrl)		
		imgQRCode.sprite = sprite
		imgQRCodeExpand.sprite = sprite
		
		btnPromotionCopy.gameObject:SetActive(true)
		--btnPromotionShare.gameObject:SetActive(true)		
	end
	
	txtDate.text 				= receiveMsg.SystemTime
	txtPromotionLink.text 		= receiveMsg.QrURL--推廣連結
	if receiveMsg.BossID ~= 0 then
		txtPromotionUpId.text 	= tostring(receiveMsg.BossID)--上級ID
		mUpID = tostring(receiveMsg.BossID)
	else
		txtPromotionUpId.text 	= ""
		mUpID = ""
	end
	txtPromotionUserID.text 	= tostring(user.UserID)
	txtPromotionGroupNum.text 	= tostring(receiveMsg.TeamMember)--團隊人數
	txtPromotionDirectly.text 	= tostring(receiveMsg.MyMember)--直屬玩家
	txtPromotionNew.text 		= tostring(receiveMsg.TodayNewMember)--今日新增
	txtPromotionActive.text 	= tostring(receiveMsg.TodayActiveMember)--今日活躍
	txtInviteCode.text			= tostring(receiveMsg.InviteCode)--邀请码
	
	UpdateUpID()
end , groupName);

-- *************************** 我的玩家 ***************************
local WithAmountPlayerID = 0	-- 记录最后查询玩家

local function CreatePrefabTeam(index)
	if(poolTeam[index] == nil)then
		local clone = CreateClone(prefabTeam)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolTeam[index] = {}
			poolTeam[index].root = clone.transform:GetComponent("RectTransform")
			poolTeam[index].txtUserID = clone.transform:Find("TextUserID"):GetComponent("Text")
			poolTeam[index].txtGroupGet = clone.transform:Find("TextGroupGet"):GetComponent("Text")
			poolTeam[index].txtPlayerGet = clone.transform:Find("TextPlayerGet"):GetComponent("Text")
			poolTeam[index].txtGroupNum = clone.transform:Find("TextGroupNum"):GetComponent("Text")
			poolTeam[index].txtDirectly = clone.transform:Find("TextDirectly"):GetComponent("Text")
			poolTeam[index].txtType = clone.transform:Find("TextType"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitAgentGetWithAmount(root)
	btnTeamFont = root.transform:Find("UIBody/DetailTeam/ButtonFont"):GetComponent("Button")
	btnTeamFont:AddClickListener(function()	
		pageTeamCurrent = pageTeamCurrent -1
		ChangePageTeam()
		SendAgentGetWithAmount()
	end)
	btnTeamNext = root.transform:Find("UIBody/DetailTeam/ButtonNext"):GetComponent("Button")
	btnTeamNext:AddClickListener(function()	
		pageTeamCurrent = pageTeamCurrent +1
		ChangePageTeam()
		SendAgentGetWithAmount()
	end)
	btnTeamFind = root.transform:Find("UIBody/DetailTeam/ButtonFind"):GetComponent("Button")
	btnTeamFind:AddClickListener(function()
		SendAgentGetWithAmountByID()
	end)
	
	txtTeamFindId = root.transform:Find("UIBody/DetailTeam/InputUserId/Text"):GetComponent("Text")
	txtTeamPageCurrent = root.transform:Find("UIBody/DetailTeam/Page/TextCurrent"):GetComponent("Text")
	txtTeamPageTotal = root.transform:Find("UIBody/DetailTeam/Page/TextTotal"):GetComponent("Text")
	
	prefabTeam = root.transform:Find("UIBody/DetailTeam/ScrollView/Content/Prefab")
	prefabTeam.gameObject:SetActive(false)
	pageTeamPoolCount = 7 --單頁數量
	pageTeamCurrent = 1	--目前所在頁數
	poolTeam = {}
	teamList = {}
	for i=1, pageTeamPoolCount do
		CreatePrefabTeam(i)
	end
end
--查詢
function SendAgentGetWithAmount()
	WithAmountPlayerID = 0
	
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.TargetID = WithAmountPlayerID
	param.PageIndex = pageTeamCurrent
	param.PageSize = pageTeamPoolCount	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentGetWithAmountListID, param)
end
--查詢玩家
function SendAgentGetWithAmountByID()
	-- 没输入 查询所有玩家
	if txtTeamFindId.text == "" then
		pageTeamCurrent = 1
		SendAgentGetWithAmount()
		return
	end
	
	local TargetID = tonumber(txtTeamFindId.text)
	-- 和上一次查询不同
	if WithAmountPlayerID ~= TargetID then
		WithAmountPlayerID = TargetID
		pageTeamCurrent = 1
	end
	
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.TargetID = WithAmountPlayerID
	param.PageIndex = pageTeamCurrent
	param.PageSize = pageTeamPoolCount	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentGetWithAmountListID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentGetWithAmountListRetID, function(receiveMsg)

	teamList = receiveMsg.UserAgentGetWithAmountList or {}	
	--pageTeamCurrent = 1
	pageTeamMax =math.ceil(receiveMsg.PageCount/7) 
	ChangePageTeam()
end , groupName);
-- 切換分頁
function ChangePageTeam()
	if pageTeamMax == 0 then
		txtTeamPageCurrent.text = "1"
		txtTeamPageTotal.text = "1"
	else
		txtTeamPageCurrent.text = tostring(pageTeamCurrent)
		txtTeamPageTotal.text = tostring(pageTeamMax )
	end
	
	btnTeamFont.interactable = (pageTeamCurrent > 1)
	btnTeamNext.interactable = (pageTeamCurrent < pageTeamMax)
		
	for i=1, pageTeamPoolCount do
		poolTeam[i].root.gameObject:SetActive(true)
		poolTeam[i].txtUserID.text 		= ""
		poolTeam[i].txtGroupGet.text 	= ""
		poolTeam[i].txtPlayerGet.text 	= ""
		poolTeam[i].txtGroupNum.text 	= ""
		poolTeam[i].txtDirectly.text	= ""
		poolTeam[i].txtType.text 		= ""
	end
		
	local list = teamList or {}
	for i=1, #list do
		poolTeam[i].root.gameObject:SetActive(true)
		poolTeam[i].txtUserID.text 		= list[i].Userid
		poolTeam[i].txtGroupGet.text 	= tostring(list[i].Groupoffer *0.01)--本周團隊貢獻值
		poolTeam[i].txtPlayerGet.text 	= tostring(list[i].Playeroffer *0.01)--本周玩家貢獻值
		poolTeam[i].txtGroupNum.text 	= tostring(list[i].Membernumbers)--團隊人數
		poolTeam[i].txtDirectly.text 	= tostring(list[i].Belongs)--直屬玩家
		
		--類型 1玩家 2代理
		if list[i].Type == 1 then
			poolTeam[i].txtType.text = "玩家"
		elseif list[i].Type == 2 then
			poolTeam[i].txtType.text = "代理"
		end
	end
end

-- *************************** 我的業績 ***************************
local function CreatePrefabResult(index)
	if(poolResult[index] == nil)then
		local clone = CreateClone(prefabResult)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolResult[index] = {}
			poolResult[index].root = clone.transform:GetComponent("RectTransform")
			poolResult[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolResult[index].txtTotal = clone.transform:Find("TextTotal"):GetComponent("Text")
			poolResult[index].txtGroup = clone.transform:Find("TextGroup"):GetComponent("Text")
			poolResult[index].txtDirectly = clone.transform:Find("TextDirectly"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitUserAgentResult(root)
	btnResultFont = root.transform:Find("UIBody/DetailResult/ButtonFont"):GetComponent("Button")
	btnResultFont:AddClickListener(function()	
		pageResultCurrent = pageResultCurrent -1
		ChangePageResult()
		SendUserAgentResult()
	end)
	btnResultNext = root.transform:Find("UIBody/DetailResult/ButtonNext"):GetComponent("Button")
	btnResultNext:AddClickListener(function()	
		pageResultCurrent = pageResultCurrent +1
		ChangePageResult()
		SendUserAgentResult()
	end)
	
	
	txtResultWeek = root.transform:Find("UIBody/DetailResult/TextWeek"):GetComponent("Text")
	txtResultGroup = root.transform:Find("UIBody/DetailResult/TextGroup"):GetComponent("Text")
	txtResultDirectresult = root.transform:Find("UIBody/DetailResult/TextDirectresult"):GetComponent("Text")
	txtResultBonus = root.transform:Find("UIBody/DetailResult/TextBonus"):GetComponent("Text")
	
	txtResultPageCurrent = root.transform:Find("UIBody/DetailResult/Page/TextCurrent"):GetComponent("Text")
	txtResultPageTotal = root.transform:Find("UIBody/DetailResult/Page/TextTotal"):GetComponent("Text")
	
	prefabResult = root.transform:Find("UIBody/DetailResult/ScrollView/Content/Prefab")
	prefabResult.gameObject:SetActive(false)
	pageResultPoolCount = 7 --單頁數量
	pageResultCurrent = 1	--目前所在頁數
	poolResult = {}
	resultList = {}
	for i=1, pageResultPoolCount do
		CreatePrefabResult(i)
	end
end
-- 查詢
function SendUserAgentResult()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.PageIndex = pageResultCurrent
	param.PageSize = pageResultPoolCount
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentResultListID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentResultListRetID, function(receiveMsg)
	txtResultWeek.text = tostring(receiveMsg.Result *0.01)--週總業績
	txtResultGroup.text = tostring(receiveMsg.Groupresult *0.01)--本週團隊業機績
	txtResultDirectresult.text = tostring(receiveMsg.Directresult *0.01)--本週直屬總業績
	txtResultBonus.text = tostring(receiveMsg.Bonus *0.01)--金現獎勵

	resultList = receiveMsg.UserAgentResultList or {}	
	pageResultMax = math.ceil(receiveMsg.PageCount/pageResultPoolCount)
	ChangePageResult()
end , groupName);
--切換分頁
function ChangePageResult()
	if pageResultMax == 0 then
		txtResultPageCurrent.text = "1"	
		txtResultPageTotal.text = "1"	
	else
		txtResultPageCurrent.text = tostring(pageResultCurrent)
		txtResultPageTotal.text = tostring(pageResultMax)
	end
	
	btnResultFont.interactable = (pageResultCurrent > 1)
	btnResultNext.interactable = (pageResultCurrent < pageResultMax)
		
	for i=1, pageResultPoolCount do
		poolResult[i].txtTime.text 		= ""
		poolResult[i].txtTotal.text 	= ""
		poolResult[i].txtGroup.text 	= ""
		poolResult[i].txtDirectly.text 	= ""
	end
	
	local list = resultList or {}
	for i=1, #list do
		poolResult[i].txtTime.text 		= list[i].Time--時間
		poolResult[i].txtTotal.text 	= tostring(list[i].Dayresult *0.01)--日總業績
		poolResult[i].txtGroup.text 	= tostring(list[i].Daygroupresult *0.01)--日團隊業績
		poolResult[i].txtDirectly.text 	= tostring(list[i].Daydirectresult *0.01)--日直屬業績
	end
end

-- *************************** 我的獎勵 ***************************
local function CreatePrefabReward(index)
	if(poolReward[index] == nil)then
		local clone = CreateClone(prefabReward)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolReward[index] = {}
			poolReward[index].root = clone.transform:GetComponent("RectTransform")
			poolReward[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolReward[index].txtTotal = clone.transform:Find("TextTotal"):GetComponent("Text")
			poolReward[index].txtGroup = clone.transform:Find("TextGroup"):GetComponent("Text")
			poolReward[index].txtDirectly = clone.transform:Find("TextDirectly"):GetComponent("Text")
			poolReward[index].txtReward = clone.transform:Find("TextReward"):GetComponent("Text")
			poolReward[index].txtState = clone.transform:Find("TextState"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitUserAgentWinAmount(root)
	btnRewardFont = root.transform:Find("UIBody/DetailReward/ButtonFont"):GetComponent("Button")
	btnRewardFont:AddClickListener(function()	
		pageRewardCurrent = pageRewardCurrent -1
		ChangePageReward()
	end)
	btnRewardNext = root.transform:Find("UIBody/DetailReward/ButtonNext"):GetComponent("Button")
	btnRewardNext:AddClickListener(function()	
		pageRewardCurrent = pageRewardCurrent +1
		ChangePageReward()

	end)
	
	txtRewardPageCurrent = root.transform:Find("UIBody/DetailReward/Page/TextCurrent"):GetComponent("Text")
	txtRewardPageTotal = root.transform:Find("UIBody/DetailReward/Page/TextTotal"):GetComponent("Text")
	
	prefabReward = root.transform:Find("UIBody/DetailReward/ScrollView/Content/Prefab")
	prefabReward.gameObject:SetActive(false)
	pageRewardPoolCount = 8 --單頁數量
	pageRewardCurrent = 1	--目前所在頁數
	poolReward = {}
	rewardList = {}
	for i=1, pageRewardPoolCount do
		CreatePrefabReward(i)
	end
end
-- 查詢
function SendUserAgentWinAmount()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentWinAmountListID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentWinAmountListRetID, function(receiveMsg)
	rewardList = receiveMsg.UserAgentWinAmount or {}	
	pageRewardCurrent = 1
	pageRewardMax = math.ceil(#rewardList / pageRewardPoolCount)
	ChangePageReward()
end , groupName);
-- 切換分頁
function ChangePageReward()
	if pageRewardMax == 0 then
		txtRewardPageCurrent.text = "1"	
		txtRewardPageTotal.text = "1"	
	else
		txtRewardPageCurrent.text = tostring(pageRewardCurrent)
		txtRewardPageTotal.text = tostring(pageRewardMax)
	end
	
	btnRewardFont.interactable = (pageRewardCurrent > 1)
	btnRewardNext.interactable = (pageRewardCurrent < pageRewardMax)
		
	for i=1, #poolReward do
		poolReward[i].txtTime.text 		= ""
		poolReward[i].txtTotal.text 	= ""
		poolReward[i].txtGroup.text 	= ""
		poolReward[i].txtDirectly.text 	= ""
		poolReward[i].txtReward.text 	= ""
		poolReward[i].txtState.text 	= ""
	end

	local list = rewardList or {}
	for i=1, #list do
		if i > pageRewardPoolCount then
			break
		end
		poolReward[i].txtTime.text 		= list[i].SendTime--時間
		poolReward[i].txtTotal.text 	= tostring(list[i].Weekresult *0.01)--周總業績
		poolReward[i].txtGroup.text 	= tostring(list[i].Weekgroupresult *0.01)--周團隊業績
		poolReward[i].txtDirectly.text 	= tostring(list[i].Weekdirectresult *0.01)--周直屬業績
		poolReward[i].txtReward.text 	= tostring(list[i].bonus *0.01)--獎勵
				
		--狀態 回傳1是可提現，其他時就回傳0
		if list[i].Status == 1 then
			poolReward[i].txtState.text = "可提现"
		else
			poolReward[i].txtState.text = "不可提现"			
		end		
	end
end

-- *************************** 我的提現 ***************************
local function CreatePrefabCash(index)
	if(poolCash[index] == nil)then
		local clone = CreateClone(prefabCash)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolCash[index] = {}
			poolCash[index].root = clone.transform:GetComponent("RectTransform")
			poolCash[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolCash[index].txtIndex = clone.transform:Find("TextIndex"):GetComponent("Text")
			poolCash[index].txtCoin = clone.transform:Find("TextCoin"):GetComponent("Text")
			poolCash[index].txtMode = clone.transform:Find("TextMode"):GetComponent("Text")
			poolCash[index].txtAccount = clone.transform:Find("TextAccount"):GetComponent("Text")
			poolCash[index].txtState = clone.transform:Find("TextState"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitUserAgentGetWithdrawAmount(root)
	btnCashFont = root.transform:Find("UIBody/DetailCash/ButtonFont"):GetComponent("Button")
	btnCashFont:AddClickListener(function()	
		pageCashCurrent = pageCashCurrent -1
		ChangePageCash()
		SendUserAgentGetWithdrawAmount()
	end)
	btnCashNext = root.transform:Find("UIBody/DetailCash/ButtonNext"):GetComponent("Button")
	btnCashNext:AddClickListener(function()	
		pageCashCurrent = pageCashCurrent +1
		ChangePageCash()
		SendUserAgentGetWithdrawAmount()
	end)
	btnCashApply = root.transform:Find("UIBody/DetailCash/ButtonApply"):GetComponent("Button")
	btnCashApply:AddClickListener(function()		
		SendGetWithdrawOrderID(-1)
	end)
	
	-- 提领输入金钱
	rootGetCoin = root.transform:Find("RootGetCoin"):GetComponent("RectTransform")
	rootGetCoin.gameObject:SetActive(false)	
	btnCloseGetCoin = root.transform:Find("RootGetCoin/ButtonClose"):GetComponent("Button")
	btnCloseGetCoin:AddClickListener(function()	rootGetCoin.gameObject:SetActive(false) end)
	txtGetCoinOver = root.transform:Find("RootGetCoin/ImageOver/Text"):GetComponent("Text")
	inputFieldGetCoinApply = root.transform:Find("RootGetCoin/ImageApply/InputField"):GetComponent("InputField")
	btnGetCoinMax = root.transform:Find("RootGetCoin/ButtonMax"):GetComponent("Button")
	btnGetCoinMax:AddClickListener(function()
	
	
		local maxValue = math.modf(tonumber(txtGetCoinOver.text))
	
		inputFieldGetCoinApply.text = tostring(maxValue)
		

			end)
	btnGetCoinApply = root.transform:Find("RootGetCoin/ButtonApply"):GetComponent("Button")
	btnGetCoinApply:AddClickListener(function()		
		if inputFieldGetCoinApply.text == "" then
			CS.UIBubblePrompt.Show("请输入金额", groupName)
			return
		end
		
		local input = tonumber(inputFieldGetCoinApply.text)
		local maxValue = tonumber(txtGetCoinOver.text)
		if input == 0 or input > maxValue then
			CS.UIBubblePrompt.Show("金额输入错误", groupName)
			return
		end
		
		--SendGetWithdrawOrderID(-2)
		SendGetWithdrawOrderID()
	end)
	btnCloseCoin = root.transform:Find("RootGetCoin/ButtonClose"):GetComponent("Button")
	btnCloseCoin:AddClickListener(function()	
		rootGetCoin.gameObject:SetActive(false)
	end)	
	SetParentTopLayer(rootGetCoin)
	
	-- 提领方式
	rootAccount = root.transform:Find("RootAccount"):GetComponent("RectTransform")
	rootAccount.gameObject:SetActive(false)
	btnCloseAccount = root.transform:Find("RootAccount/ButtonClose"):GetComponent("Button")
	btnCloseAccount:AddClickListener(function()	
		rootAccount.gameObject:SetActive(false) 
	end)
	btnAccountAlipay = root.transform:Find("RootAccount/Option/ButtonAlipay"):GetComponent("Button")
	btnAccountAlipay:AddClickListener(function()	SendGetWithdrawOrderID(1) end)
	btnAccountAlipay.gameObject:SetActive(userdata.ConvertAlipayShow == 1)
	btnAccountUnionpay = root.transform:Find("RootAccount/Option/ButtonUnionpay"):GetComponent("Button")
	btnAccountUnionpay:AddClickListener(function()	SendGetWithdrawOrderID(2) end)
	btnAccountUnionpay.gameObject:SetActive(userdata.ConvertBankShow == 1)
	btnAccountWeixin = root.transform:Find("RootAccount/Option/ButtonWeixin"):GetComponent("Button")
	btnAccountWeixin:AddClickListener(function()	SendGetWithdrawOrderID() end)
	btnAccountWeixin.gameObject:SetActive(false)	
	SetParentTopLayer(rootAccount)	
	
	txtCashUserID = root.transform:Find("UIBody/DetailCash/TextID"):GetComponent("Text")
	txtCashTotal = root.transform:Find("UIBody/DetailCash/TextTotal"):GetComponent("Text")
	txtCashGet = root.transform:Find("UIBody/DetailCash/TextGet"):GetComponent("Text")
	txtCashCan = root.transform:Find("UIBody/DetailCash/TextCan"):GetComponent("Text")
	
	txtCashPageCurrent = root.transform:Find("UIBody/DetailCash/Page/TextCurrent"):GetComponent("Text")
	txtCashPageTotal = root.transform:Find("UIBody/DetailCash/Page/TextTotal"):GetComponent("Text")
	
	prefabCash = root.transform:Find("UIBody/DetailCash/ScrollView/Content/Prefab")
	prefabCash.gameObject:SetActive(false)
	pageCashPoolCount = 4 --單頁數量
	pageCashCurrent = 1	--目前所在頁數
	poolCash = {}
	cashList = {}
	for i=1, pageCashPoolCount do
		CreatePrefabCash(i)
	end
end
-- 查詢
function SendUserAgentGetWithdrawAmount()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.PageIndex = pageCashCurrent
	param.PageSize = pageCashPoolCount	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentGetWithdrawAmountListID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentGetWithdrawAmountListRetID, function(receiveMsg)
	txtCashUserID.text = tostring(receiveMsg.Userid)
	txtCashTotal.text = tostring(receiveMsg.Bonus *0.01)--總獎勵金額
	txtCashGet.text = tostring(receiveMsg.Rest *0.01)--已提領金額
	txtCashCan.text = tostring(receiveMsg.Allow *0.01)--可提領金額
	txtGetCoinOver.text = tostring(receiveMsg.Allow *0.01)--可提領金額
	
	cashList = receiveMsg.UserAgentGetWithdrawAmountList or {}	
	--pageCashCurrent = 1
	pageCashMax = math.ceil(receiveMsg.Pagecount/pageCashPoolCount)
	ChangePageCash()
end , groupName);
-- 切換分頁
function ChangePageCash()
	if pageCashMax == 0 then
		txtCashPageCurrent.text = "1"
		txtCashPageTotal.text = "1"
	else
		txtCashPageCurrent.text = tostring(pageCashCurrent)
		txtCashPageTotal.text = tostring(pageCashMax)
	end
	btnCashFont.interactable = (pageCashCurrent > 1)
	btnCashNext.interactable = (pageCashCurrent < pageCashMax)
		
	for i=1, #poolCash do
		poolCash[i].root.gameObject:SetActive(true)
		poolCash[i].txtTime.text 	= ""
		poolCash[i].txtIndex.text 	= ""
		poolCash[i].txtCoin.text 	= ""
		poolCash[i].txtMode.text 	= ""
		poolCash[i].txtAccount.text = ""
		poolCash[i].txtState.text 	= ""
	end
		
	local list = cashList or {}
	for i=1, #list do
		poolCash[i].root.gameObject:SetActive(true)
		poolCash[i].txtTime.text 	= list[i].Datetime--時間
		poolCash[i].txtIndex.text 	= tostring(list[i].Noid)--單號
		poolCash[i].txtCoin.text 	= tostring(list[i].Amount *0.01)--金額
		poolCash[i].txtAccount.text = tostring(list[i].Account)--提現帳號
		
		-- 提现类型 1支付宝 2银联
	--	if list[i].Type == 1 then
			poolCash[i].txtMode.text = list[i].Type
		--elseif list[i].Type == 2 then
		--	poolCash[i].txtMode.text = "银联"
		--end		
		
		--狀態 0未發放 1發放 2拒絕
		if list[i].Status == 0 then
			poolCash[i].txtState.text = "未发放"
		elseif list[i].Status == 1 then
			poolCash[i].txtState.text = "已发放"
		elseif list[i].Status == 2 then
			poolCash[i].txtState.text = "拒绝"
		else
			poolCash[i].txtState.text = ""
		end
	end
end
-- 申請提現
function SendGetWithdrawOrderID(uiexchange)
	rootGetCoin.gameObject:SetActive(false)
	rootAccount.gameObject:SetActive(false)
	

	-- 未绑定
	--[[if(userdata.IsBindGuest ~= true)then

		G.Util.coroutine_call(function()
			local showMessage = "需要绑定用户信息，确认绑定吗？"					
			local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", showMessage, "绑定", "取消");
			if boxResult == CS.MessageBoxResult.OK then
				local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
				openparam.WindowData = 1;
				openparam.WindowCloseCallBack =(function()
					userdata = G.User.GetCurrentUser().UserData
				end)
				CS.WindowManager.Instance:OpenWindow(openparam)
			end
		end)()
	-- 绑定银联
	elseif(userdata.IsBindBankCard == false and uiexchange == 2)then
        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = uiexchange;
        CS.WindowManager.Instance:OpenWindow(openparam)
	-- 绑定支付宝
	elseif(userdata.IsBindZhifubao == false and uiexchange == 1)then
        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = uiexchange;
        CS.WindowManager.Instance:OpenWindow(openparam)
	-- 输入金额
	else]]if uiexchange == -1 then
		inputFieldGetCoinApply.text = ""
		rootGetCoin.gameObject:SetActive(true)
	-- 选择领取方式
	elseif uiexchange == -2 then
		rootGetCoin.gameObject:SetActive(false)
		rootAccount.gameObject:SetActive(true)
	-- 送出提领
	else
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.LngAmount = tonumber(inputFieldGetCoinApply.text)
		param.ChvIPaddress = G.Tool.Debug.GetPublicIP()	
		param.Type=	1
		-- -- 银联
		-- if uiexchange == 1 then
		-- 	param.Type = 2
		-- -- 支付宝
		-- elseif uiexchange == 2 then
		-- 	param.Type = 1
		-- end
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitConvertOrderID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
	end
end
-- region　申請提現
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitConvertOrderRetID, function(receiveMsg)
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
end , groupName);

-- region　申請预提現
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitPreConvertOrderRetID, function(receiveMsg)
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
end , groupName);


-- *************************** 我的預提 ***************************
local function CreatePrefabLoan(index)
	if(poolLoan[index] == nil)then
		local clone = CreateClone(prefabLoan)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolLoan[index] = {}
			poolLoan[index].root = clone.transform:GetComponent("RectTransform")
			poolLoan[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolLoan[index].txtIndex = clone.transform:Find("TextIndex"):GetComponent("Text")
			poolLoan[index].txtCoin = clone.transform:Find("TextCoin"):GetComponent("Text")
			poolLoan[index].txtMode = clone.transform:Find("TextMode"):GetComponent("Text")
			poolLoan[index].txtAccount = clone.transform:Find("TextAccount"):GetComponent("Text")
			poolLoan[index].txtState = clone.transform:Find("TextState"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitUserAgentBorrow(root)
	btnLoanApply = root.transform:Find("UIBody/DetailLoan/ButtonApply"):GetComponent("Button")
	btnLoanApply:AddClickListener(function()
		inputFieldLoanGetCoinApply.text = ""
		rootLoanGetCoin.gameObject:SetActive(true)
	end)
	
	txtLoanUserID = root.transform:Find("UIBody/DetailLoan/TextID"):GetComponent("Text")
	txtLoanReward = root.transform:Find("UIBody/DetailLoan/TextReward"):GetComponent("Text")
	txtLoanGet = root.transform:Find("UIBody/DetailLoan/TextGet"):GetComponent("Text")
	txtLoanCan = root.transform:Find("UIBody/DetailLoan/TextCan"):GetComponent("Text")
	txtLoanOver = root.transform:Find("UIBody/DetailLoan/TextOver"):GetComponent("Text")
	txtLoanLast = root.transform:Find("UIBody/DetailLoan/TextLast"):GetComponent("Text")
	
	txtLoanPageCurrent = root.transform:Find("UIBody/DetailLoan/Page/TextCurrent"):GetComponent("Text")
	txtLoanPageTotal = root.transform:Find("UIBody/DetailLoan/Page/TextTotal"):GetComponent("Text")
	
	-- 帮助
	rootLoanHelp = root.transform:Find("UIBody/DetailLoan/LoanHelp"):GetComponent("RectTransform")
	rootLoanHelp.gameObject:SetActive(false)	
	btnCloseLoanHelp = root.transform:Find("UIBody/DetailLoan/LoanHelp/ButtonClose"):GetComponent("Button")
	btnCloseLoanHelp:AddClickListener(function()	rootLoanHelp.gameObject:SetActive(false) end)
	SetParentTopLayer(rootLoanHelp)		
	
	btnLoanHelp = root.transform:Find("UIBody/DetailLoan/ButtonHelp"):GetComponent("Button")
	btnLoanHelp:AddClickListener(function()	rootLoanHelp.gameObject:SetActive(true) end)
	
	-- 输入预提金钱
	rootLoanGetCoin = root.transform:Find("UIBody/DetailLoan/RootGetCoin"):GetComponent("RectTransform")
	rootLoanGetCoin.gameObject:SetActive(false)	
	btnLoanCloseGetCoin = root.transform:Find("UIBody/DetailLoan/RootGetCoin/ButtonClose"):GetComponent("Button")
	btnLoanCloseGetCoin:AddClickListener(function()	rootLoanGetCoin.gameObject:SetActive(false) end)
	txtLoanGetCoinOver = root.transform:Find("UIBody/DetailLoan/RootGetCoin/ImageOver/Text"):GetComponent("Text")
	inputFieldLoanGetCoinApply = root.transform:Find("UIBody/DetailLoan/RootGetCoin/ImageApply/InputField"):GetComponent("InputField")
	btnLoanGetCoinMax = root.transform:Find("UIBody/DetailLoan/RootGetCoin/ButtonMax"):GetComponent("Button")
	btnLoanGetCoinMax:AddClickListener(function()	
	
		local maxValue = math.modf(tonumber(txtLoanGetCoinOver.text))
		
		inputFieldLoanGetCoinApply.text = tostring(maxValue)


		end)
	btnLoanGetCoinApply = root.transform:Find("UIBody/DetailLoan/RootGetCoin/ButtonApply"):GetComponent("Button")
	btnLoanGetCoinApply:AddClickListener(function() SendLoanGetCoin() end)
	SetParentTopLayer(rootLoanGetCoin)

	prefabLoan = root.transform:Find("UIBody/DetailLoan/ScrollView/Content/Prefab")
	prefabLoan.gameObject:SetActive(false)
	pageLoanPoolCount = 6 --單頁數量
	pageLoanCurrent = 1	--目前所在頁數
	poolLoan = {}
	loanList = {}
	for i=1, pageLoanPoolCount do
		CreatePrefabLoan(i)
	end
end
-- 查詢
function SendUserAgentBorrow()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.PageIndex = pageLoanCurrent
	param.PageSize = pageLoanPoolCount
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentBorrowListID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentBorrowListRetID, function(receiveMsg)
	txtLoanUserID.text 		= G.User.GetCurrentUserID()
	txtLoanReward.text 		= tostring(receiveMsg.bonusamount *0.01)--截止本日零點預估現金獎勵
	txtLoanCan.text 		= tostring(receiveMsg.AmountLimit *0.01)--可提領金額
	txtLoanGet.text 		= tostring(receiveMsg.PreConvertAmount *0.01)--已預提金額
	txtLoanOver.text 		= tostring(receiveMsg.RemainAmount *0.01)--剩餘可預提金額
	txtLoanLast.text 		= "上周已预提金额: " .. tostring(receiveMsg.withdrawlastweek *0.01)--上周已预提金额
	txtLoanGetCoinOver.text = tostring(receiveMsg.AmountLimit*0.01)
	
	loanList = receiveMsg.UserAgentBorrowList or {}	
	pageLoanCurrent = 1
	pageLoanMax = receiveMsg.PageCount
	ChangePageLoan()
end , groupName);
-- 切換分頁
function ChangePageLoan()		
	for i=1, #poolLoan do
		poolLoan[i].root.gameObject:SetActive(true)
		poolLoan[i].txtTime.text 	= ""
		poolLoan[i].txtIndex.text 	= ""
		poolLoan[i].txtCoin.text 	= ""
		poolLoan[i].txtMode.text 	= ""
		poolLoan[i].txtAccount.text = ""
		poolLoan[i].txtState.text 	= ""
	end
	
	local list = loanList or {}
	for i=1, #list do
		if i > pageLoanPoolCount then
			break
		end
	
		poolLoan[i].root.gameObject:SetActive(true)		
		poolLoan[i].txtTime.text 	= list[i].LogDate--時間
		poolLoan[i].txtIndex.text 	= tostring(list[i].ID)--單號
		poolLoan[i].txtCoin.text 	= tostring(list[i].Amount *0.01)--金額
		poolLoan[i].txtMode.text 	= "游戏"--提現類型
		poolLoan[i].txtAccount.text = tostring(list[i].Account)--提現帳號

		poolLoan[i].txtState.text = list[i].Status
		--回傳1是已發放，0是無發放
		-- if list[i].Status == 1 then
			-- poolLoan[i].txtState.text = "已发放"
		-- else
			-- poolLoan[i].txtState.text = "未发放"
		-- end
	end
end
-- 送出提现
function SendLoanGetCoin()
	if inputFieldLoanGetCoinApply.text == "" then
		CS.UIBubblePrompt.Show("请输入金额", groupName)
		return
	end
	
	local input = tonumber(inputFieldLoanGetCoinApply.text)
	local maxValue = tonumber(txtLoanGetCoinOver.text)
	if input == 0 or input > maxValue then
		CS.UIBubblePrompt.Show("金额输入错误", groupName)
		return
	end
	
	local param = {}
	 
	 param.UserID = G.User.GetCurrentUserID()	
	 param.Amount = tonumber(inputFieldLoanGetCoinApply.text)
	      
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitPreConvertOrderID, param)
	--CS.UILoadingData.Show("送出申请，等待回应", 10)
end
-- region　申請提現
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitPreConvertOrderRet, function(receiveMsg)
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
end , groupName);

-- *************************** 獎勵說明 ***************************
local function CreatePrefabHelp(index)
	if(poolHelp[index] == nil)then
		local clone = CreateClone(prefabHelp)
		if(clone ~= nil)then
			clone.gameObject:SetActive(true)
			poolHelp[index] = {}
			poolHelp[index].root = clone.transform:GetComponent("RectTransform")
			poolHelp[index].txtLevel = clone.transform:Find("TextLevel"):GetComponent("Text")
			poolHelp[index].txtPerformance = clone.transform:Find("TextPerformance"):GetComponent("Text")
			poolHelp[index].txtCommission = clone.transform:Find("TextCommission"):GetComponent("Text")
		end
	end
end
-- 初始化
function InitHelp(root)
	txtHelpText = {}
	rootHelpDescription1 = root.transform:Find("UIBody/DetailHelp/ScrollView/Content/Description1")
	
	for i=1, 5 do
		data = rootHelpDescription1.transform:Find("Text" .. tostring(i)):GetComponent("Text")
		table.insert(txtHelpText, data)
	end
	txtHelpTotal = rootHelpDescription1.transform:Find("TextTotal"):GetComponent("Text")
	
	prefabHelp = root.transform:Find("UIBody/DetailHelp/ScrollView/Content/Description2/Content/Prefab")	
	prefabHelp.gameObject:SetActive(false)
	pageHelpPoolCount = 15 --單頁數量
	poolHelp = {}
	helpList = {}
	for i=1, pageHelpPoolCount do
		CreatePrefabHelp(i)
	end
end
-- 查詢
function SendRewardTips()
	local param = {}
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetRewardTipsID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetRewardTipsRetID, function(receiveMsg)		
	txtHelpText[1].text = string.gsub(receiveMsg.Self, "<BR>", "\n")
	txtHelpText[2].text = string.gsub(receiveMsg.AgentA, "<BR>", "\n")
	txtHelpText[3].text = string.gsub(receiveMsg.AgentB, "<BR>", "\n")
	txtHelpText[4].text = string.gsub(receiveMsg.AgentC, "<BR>", "\n")
	txtHelpText[5].text = string.gsub(receiveMsg.AgentD, "<BR>", "\n")
	txtHelpTotal.text = receiveMsg.AgentBonus
	
	helpList = receiveMsg.GetRewardTipsRetData or {}	
	
	ChangePageHelp()
end , groupName);
-- 切換分頁
function ChangePageHelp()
	for i=1, #poolHelp do
		poolHelp[i].txtLevel.text 		= ""
		poolHelp[i].txtPerformance.text = ""
		poolHelp[i].txtCommission.text 	= ""
	end
	
	local list = helpList or {}
	for i=1, #list do
		if i > pageHelpPoolCount then
			break
		end
		
		poolHelp[i].txtLevel.text = list[i].name
		poolHelp[i].txtCommission.text = "每万" .. list[i].BonusRate .. "元"
		
		local rangText
		if list[i].minamount >= 10000 then
			t1, t2 = math.modf(list[i].minamount/10000)
			rangText = tostring(t1) .. "万"
		else
			rangText = tostring(list[i].minamount)		
		end
		
		if list[i].maxamount > 0 then
			rangText = rangText .. "到"
			if list[i].maxamount >= 10000 then
				t1, t2 = math.modf(list[i].maxamount/10000)
				rangText = rangText .. tostring(t1) .. "万"
			else
				rangText = rangText .. tostring(list[i].maxamount)			
			end
		else
			rangText = rangText .. "以上"
		end
		poolHelp[i].txtPerformance.text = rangText
	end
end

-- 绑定邀请码
function SendSplitBindingInvitedCodeID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.InviteCode = ipSendInviteCode.text
		
		if G.Tool.StringHelper.IsEmpty(ipSendInviteCode.text) then
			CS.UIBubblePrompt.Show("请输入邀请码", groupName)
			return
		end
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitAddUserRelationID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 绑定邀请码結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UnlimitAddUserRelationRetID, function(receiveMsg)
	
	if receiveMsg.RetCode == 1 then
		SendMyPromotion()
		--mUpID = ipSendInviteCode.text
		--UpdateUpID()
	end
	
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
	
end , groupName);

function UpdateUpID()

	if mUpID == nil or mUpID == "" then
		rootSendInviteCode.gameObject:SetActive(true)
	else
		rootSendInviteCode.gameObject:SetActive(false)	
	end
end