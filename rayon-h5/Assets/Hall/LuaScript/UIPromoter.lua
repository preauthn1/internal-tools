require("Module")
mmodule(...)local G = require "Global";
local groupName = "UIPromoter";

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local userdata = G.User.GetCurrentUser().UserData
local subList = {}
subList[1] = {}
subList[1].myName = "QRCode"
subList[2] = {}
subList[2].myName = "Overview"
subList[3] = {}
subList[3].myName = "Receive"

local mQRCodeUrl
local mAgentLevel
local mUpID = ""

local function CreateClone(prefab)
	local clone = CS.UnityEngine.GameObject.Instantiate(prefab)
	if(clone == nil)then return nil end
	clone.transform:SetParent(prefab.transform.parent)
	clone.transform.localPosition = CS.UnityEngine.Vector3.zero
	clone.transform.localScale = CS.UnityEngine.Vector3.one
	
	return clone	
end

local function CreatePrefabTax(index)
	if(poolTax[index] == nil)then
		local clone = CreateClone(prefabTax)
		if(clone ~= nil)then
			poolTax[index] = {}
			poolTax[index].root = clone.transform:GetComponent("RectTransform")
			poolTax[index].txtPlayerId = clone.transform:Find("TextPlayerId"):GetComponent("Text")
			poolTax[index].txtTaxBase = clone.transform:Find("TextTaxBase"):GetComponent("Text")
			poolTax[index].txtTaxOne = clone.transform:Find("TextTaxOne"):GetComponent("Text")
			poolTax[index].txtTaxGet = clone.transform:Find("TextTaxGet"):GetComponent("Text")
		end
	end
end

local function CreatePrefabIncome(index)
	if(poolIncome[index] == nil)then
		local clone = CreateClone(prefabIncome)
		if(clone ~= nil)then
			poolIncome[index] = {}
			poolIncome[index].root = clone.transform:GetComponent("RectTransform")
			poolIncome[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolIncome[index].txtTaxOne = clone.transform:Find("TextTaxOne"):GetComponent("Text")
			poolIncome[index].txtTaxTwo = clone.transform:Find("TextTaxTwo"):GetComponent("Text")
			poolIncome[index].txtTaxGet = clone.transform:Find("TextTaxGet"):GetComponent("Text")
		end
	end
end

local function CreatePrefabReceive(index)
	if(poolReceive[index] == nil)then
		local clone = CreateClone(prefabReceive)
		if(clone ~= nil)then
			poolReceive[index] = {}
			poolReceive[index].root = clone.transform:GetComponent("RectTransform")
			poolReceive[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
			poolReceive[index].txtCoin = clone.transform:Find("TextCoin"):GetComponent("Text")
			poolReceive[index].txtState = clone.transform:Find("TextState"):GetComponent("Text")
		end
	end
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
	
	btnCopyQRCode = root.transform:Find("UIBody/DetailQRCode/ImageQRCode"):GetComponent("Button")
	btnCopyQRCode:AddClickListener(function()
        CS.Utility.CopyText(mQRCodeUrl)
        CS.UIBubblePrompt.Show("复制成功", groupName)
	end)	
	
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
		
	-- 我的推广码
	imgQRCode = root.transform:Find("UIBody/DetailQRCode/ImageQRCode"):GetComponent("Image")
	
	------- 我的推广码(QRCode)
	rootQRCodeExpand = root.transform:Find("RootQRCodeExpand"):GetComponent("RectTransform")
	rootQRCodeExpand.gameObject:SetActive(false)
	imgQRCodeExpandBg = root.transform:Find("RootQRCodeExpand/ImageBg"):GetComponent("RectTransform")
	imgQRCodeExpand = root.transform:Find("RootQRCodeExpand/ImageQRCode"):GetComponent("Image")
	
	btnQRCodeExpand = root.transform:Find("UIBody/DetailQRCode/ButtonExpand"):GetComponent("Button")
	btnQRCodeExpand:AddClickListener(function()	rootQRCodeExpand.gameObject:SetActive(true) end)
	btnCloseQRCodeExpand = root.transform:Find("RootQRCodeExpand"):GetComponent("Button")
	btnCloseQRCodeExpand:AddClickListener(function()	rootQRCodeExpand.gameObject:SetActive(false) end)
	
	--邀请码
	txtInviteCode = root.transform:Find("UIBody/DetailQRCode/InviteCode/Code"):GetComponent("Text")
	btnInviteCodeCopy = root.transform:Find("UIBody/DetailQRCode/InviteCode/BtnCopy"):GetComponent("Button")
	btnInviteCodeCopy:AddClickListener(function() 
		CS.Utility.CopyText(txtInviteCode.text)
		CS.UIBubblePrompt.Show("复制成功", groupName)
	end)
	--绑定
	rootSendInviteCode = root.transform:Find("UIBody/DetailQRCode/SendInviteCode"):GetComponent("RectTransform")
	ipSendInviteCode = root.transform:Find("UIBody/DetailQRCode/SendInviteCode/InputFieldCode"):GetComponent("InputField")
	btnSendInviteCode = root.transform:Find("UIBody/DetailQRCode/SendInviteCode/BtnSend"):GetComponent("Button")
	btnSendInviteCode:AddClickListener(function()
		SendSplitBindingInvitedCodeID()
	end)
	--上级ID
	rootUpID = root.transform:Find("UIBody/DetailQRCode/UpID"):GetComponent("RectTransform")
	txtUpID = root.transform:Find("UIBody/DetailQRCode/UpID/Code"):GetComponent("Text")
	rootUpID.gameObject:SetActive(false)
	
	
	------- 我的推广码(收益规则)
	rootRule = root.transform:Find("RootRule"):GetComponent("RectTransform")
	rootRule.gameObject:SetActive(false)
	btnRule = root.transform:Find("UIBody/DetailQRCode/ImageBuoy/ButtonRule"):GetComponent("Button")
	btnRule:AddClickListener(function()	rootRule.gameObject:SetActive(true) end)
	btnCloseRule = root.transform:Find("RootRule/ButtonClose"):GetComponent("Button")
	btnCloseRule:AddClickListener(function()	rootRule.gameObject:SetActive(false) end)
	
	poolRuleCard = {}
	for i=1, 6 do
		poolRuleCard[i] = {}
		poolRuleCard[i].txtStep = {}
		poolRuleCard[i].txtStep[1] = root.transform:Find("RootRule/AttachCard" .. i .. "/Text1"):GetComponent("Text")
		poolRuleCard[i].txtStep[2] = root.transform:Find("RootRule/AttachCard" .. i .. "/Text2"):GetComponent("Text")
		poolRuleCard[i].txtStep[3] = root.transform:Find("RootRule/AttachCard" .. i .. "/Text3"):GetComponent("Text")
	end
	
	------- 我的推广码(联系客服)
	btnService = root.transform:Find("UIBody/DetailQRCode/ImageBuoy/ButtonService"):GetComponent("Button")
	btnService:AddClickListener(function() 
		local openparam = CS.WindowNodeInitParam("UISendMsgToGM", CS.AppDefine.HallNumber)
        CS.WindowManager.Instance:OpenWindow(openparam)
	end)
	------- 我的推广码(推广引导)
	imgGuide = root.transform:Find("UIBody/DetailQRCode/ButtonCloseGuide/Image"):GetComponent("Image")
	btnGuide = root.transform:Find("UIBody/DetailQRCode/ImageBuoy/ButtonGuide"):GetComponent("Button")
	btnGuide:AddClickListener(function()
		imgGuide.transform.localScale = CS.UnityEngine.Vector3(0,1,1)
			btnCloseGuide.gameObject:SetActive(true)
		local sequence = CS.DG.Tweening.DOTween.Sequence()
		sequence:Append(imgGuide.transform:DOScaleX(1,0.2))
	end)
	btnCloseGuide = root.transform:Find("UIBody/DetailQRCode/ButtonCloseGuide"):GetComponent("Button")
	btnCloseGuide:AddClickListener(function()	
		local sequence = CS.DG.Tweening.DOTween.Sequence()
		sequence:Append(imgGuide.transform:DOScaleX(0,0.2))
		sequence:OnComplete(function()
			btnCloseGuide.gameObject:SetActive(false)
		end) 
	end)
	btnCloseGuide.gameObject:SetActive(false)
	
	-- 推广概况
	imgRankIcon = root.transform:Find("UIBody/DetailOverview/ButtonRank/ImageIcon"):GetComponent("Image")
	imgRankName = root.transform:Find("UIBody/DetailOverview/ButtonRank/ImageName"):GetComponent("Image")
	txtRankExp = root.transform:Find("UIBody/DetailOverview/ButtonRank/TextExp"):GetComponent("Text")
	-- 推广概况(推广等级)
	rootRank = root.transform:Find("RootRank"):GetComponent("RectTransform")
	btnRank = root.transform:Find("UIBody/DetailOverview/ButtonRank"):GetComponent("Button")
	btnRank:AddClickListener(function()	
		rootRank.gameObject:SetActive(true) 
		scrollRankCard:SetPageCount(#poolRankCard)

		
		mAgentLevel = mAgentLevel or 1
		scrollRankCard:ToFirstPage();
		for i=2, mAgentLevel do
			scrollRankCard:ToNextPage()
		end
	end)
	btnCloseRank = root.transform:Find("RootRank/ButtonClose"):GetComponent("Button")
	btnCloseRank:AddClickListener(function()	rootRank.gameObject:SetActive(false) end)
	
	rootShowExpValue = root.transform:Find("UIBody/DetailOverview/ButtonRank/ImageExp/Value"):GetComponent("RectTransform")
	txtShowExpValue = root.transform:Find("UIBody/DetailOverview/ButtonRank/TextExp"):GetComponent("Text")
	imgIconShow = root.transform:Find("UIBody/DetailOverview/ButtonRank/ImageIcon"):GetComponent("Image")
	imgNameShow = root.transform:Find("UIBody/DetailOverview/ButtonRank/ImageName"):GetComponent("Image")
	
	-- 推广概况(今日税收)
	rootTax = root.transform:Find("RootTax"):GetComponent("RectTransform")
	rootTax.gameObject:SetActive(false)
	btnTax = root.transform:Find("UIBody/DetailOverview/ButtonTax"):GetComponent("Button")
	btnTax:AddClickListener(function() rootTax.gameObject:SetActive(true) end)
	btnCloseTax = root.transform:Find("RootTax/ButtonClose"):GetComponent("Button")
	btnCloseTax:AddClickListener(function()	rootTax.gameObject:SetActive(false) end)
	
	txtTaxCount = root.transform:Find("UIBody/DetailOverview/ButtonTax/TextCount"):GetComponent("Text")
	
	-- 推广概况(今日收益)
	rootIncome = root.transform:Find("RootIncome"):GetComponent("RectTransform")
	rootIncome.gameObject:SetActive(false)
	btnIncome = root.transform:Find("UIBody/DetailOverview/ButtonIncome"):GetComponent("Button")
	btnIncome:AddClickListener(function()	rootIncome.gameObject:SetActive(true) end)
	btnCloseIncome = root.transform:Find("RootIncome/ButtonClose"):GetComponent("Button")
	btnCloseIncome:AddClickListener(function()	rootIncome.gameObject:SetActive(false) end)
	
	txtIncomeCount = root.transform:Find("UIBody/DetailOverview/ButtonIncome/TextCount"):GetComponent("Text")
	
	-- 推广概况(推广人数)
	rootPeople = root.transform:Find("RootPeople"):GetComponent("RectTransform")
	rootPeople.gameObject:SetActive(false)
	btnPeople = root.transform:Find("UIBody/DetailOverview/ButtonPeople"):GetComponent("Button")
	btnPeople:AddClickListener(function()	rootPeople.gameObject:SetActive(true) end)
	btnClosePeople = root.transform:Find("RootPeople/ButtonClose"):GetComponent("Button")
	btnClosePeople:AddClickListener(function()	rootPeople.gameObject:SetActive(false) end)
	
	txtPeopleCount = root.transform:Find("UIBody/DetailOverview/ButtonPeople/TextCount"):GetComponent("Text")
	
	-- 领取收益
	txtReceiveTotal = root.transform:Find("UIBody/DetailReceive/ImageState/TextTotal"):GetComponent("Text")
	txtReceiveAlready = root.transform:Find("UIBody/DetailReceive/ImageState/TextAlready"):GetComponent("Text")
	txtReceiveToday = root.transform:Find("UIBody/DetailReceive/ImageState/TextToday"):GetComponent("Text")
	btnReceiveGet = root.transform:Find("UIBody/DetailReceive/ImageState/ButtonGet"):GetComponent("Button")
	btnReceiveGet:AddClickListener(function()	rootGetCoin.gameObject:SetActive(true)	end)
	
	prefabReceive = root.transform:Find("UIBody/DetailReceive/ScrollView/Content/Prefab")
	prefabReceive.gameObject:SetActive(false)
	poolReceive = {}
	
	rootReceiveNoData = root.transform:Find("UIBody/DetailReceive/RootNoData"):GetComponent("RectTransform")
	
	------- 推广概况(推广等级)
	txtExpMax = root.transform:Find("RootRank/TextExpMax"):GetComponent("Text")
	imgIconCurrent = root.transform:Find("RootRank/ImageIconCurrent"):GetComponent("Image")
	imgNameCurrent = root.transform:Find("RootRank/ImageNameCurrent"):GetComponent("Image")
	imgIconNext = root.transform:Find("RootRank/ImageIconNext"):GetComponent("Image")
	imgNameNext = root.transform:Find("RootRank/ImageNameNext"):GetComponent("Image")
	rootExpValue = root.transform:Find("RootRank/ImageExp/Value"):GetComponent("RectTransform")
	txtExpValue = root.transform:Find("RootRank/TextExp"):GetComponent("Text")
	
	poolRankCard = {}
	for i=1, 6 do
		poolRankCard[i] = {}
		poolRankCard[i].root = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i):GetComponent("RectTransform")
		poolRankCard[i].imgIcon = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i .. "/ImageIcon"):GetComponent("Image")
		poolRankCard[i].imgName = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i .. "/ImageName"):GetComponent("Image")
		poolRankCard[i].txtStep = {}
		poolRankCard[i].txtStep[1] = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i .. "/TextStep1"):GetComponent("Text")
		poolRankCard[i].txtStep[2] = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i .. "/TextStep2"):GetComponent("Text")
		poolRankCard[i].txtStep[3] = root.transform:Find("RootRank/ScrollView/Content/ImageCard" .. i .. "/TextStep3"):GetComponent("Text")
		poolRankCard[i].expMax = 0
	end
	
	scrollRankCard = root.transform:Find("RootRank/ScrollView"):GetComponent("ScrollPage")
	btnRankLeft = root.transform:Find("RootRank/ButtonLeft"):GetComponent("Button")
	btnRankLeft:AddClickListener(function()	scrollRankCard:ToPrePage() end)	
	btnRankRight = root.transform:Find("RootRank/ButtonRight"):GetComponent("Button")
	btnRankRight:AddClickListener(function() scrollRankCard:ToNextPage() print("right") end)
	
	------- 推广概况(今日税收)
	rootTaxWeekList = root.transform:Find("RootTax/WeekList"):GetComponent("RectTransform")
	rootTaxWeekList.transform.localScale = CS.UnityEngine.Vector3(1,0,1)
	rootTaxWeekDown = root.transform:Find("RootTax/ButtonWeek/ImageDown"):GetComponent("RectTransform")
	rootTaxWeekDown.gameObject:SetActive(false)
	rootTaxWeekUp = root.transform:Find("RootTax/ButtonWeek/ImageUp"):GetComponent("RectTransform")
	rootTaxWeekUp.gameObject:SetActive(true)
	btnTaxWeekList = root.transform:Find("RootTax/ButtonWeek"):GetComponent("Button")
	btnTaxWeekList:AddClickListener(function()	ShowTaxWeekList() end)	
	txtTaxWeekList = root.transform:Find("RootTax/ButtonWeek/Text"):GetComponent("Text")
	txtTaxWeekList.text = os.date("%m / %d")
	poolTaxWeek = {}
	for i=1,7 do
		local tempIndex = i
		poolTaxWeek[i] = {}
		poolTaxWeek[i].txt = root.transform:Find("RootTax/WeekList/Button" .. i .. "/Text"):GetComponent("Text")
		poolTaxWeek[i].dateTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")-i+1}))
		poolTaxWeek[i].txt.text = os.date("%m / %d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")-i+1}))
		poolTaxWeek[i].btn = root.transform:Find("RootTax/WeekList/Button" .. i):GetComponent("Button")
		poolTaxWeek[i].btn:AddClickListener(function()	
			txtTaxWeekList.text = poolTaxWeek[tempIndex].txt.text
			ShowTaxWeekList()
			SendGetUserAgentTaxLogID(poolTaxWeek[tempIndex].dateTime)
		end)
	end
	
	txtTaxTotalCoin = root.transform:Find("RootTax/TextTotalCoin"):GetComponent("Text")
	txtTaxTotalReward = root.transform:Find("RootTax/TextTotalReward"):GetComponent("Text")
	txtTaxTotalCount = root.transform:Find("RootTax/TextTotalCount"):GetComponent("Text")
	
	rootTaxNoData = root.transform:Find("RootTax/RootNoData"):GetComponent("RectTransform")
	
	prefabTax = root.transform:Find("RootTax/ScrollView/Content/Prefab")
	prefabTax.gameObject:SetActive(false)
	poolTax = {}
	
	------- 推广概况(今日收益)
	prefabIncome = root.transform:Find("RootIncome/ScrollView/Content/Prefab")
	prefabIncome.gameObject:SetActive(false)
	poolIncome = {}
	
	rootIncomeNoData = root.transform:Find("RootIncome/RootNoData"):GetComponent("RectTransform")
	
	------- 推广概况(推广人数)
	txtTeam = root.transform:Find("RootPeople/TextTeam"):GetComponent("Text")
	txtDirectly = root.transform:Find("RootPeople/TextDirectly"):GetComponent("Text")
	txtIncrease = root.transform:Find("RootPeople/TextIncrease"):GetComponent("Text")
	txtLogin = root.transform:Find("RootPeople/TextLogin"):GetComponent("Text")
	
	------- 代理提現(金額)
	rootGetCoin = root.transform:Find("RootGetCoin"):GetComponent("RectTransform")
	rootGetCoin.gameObject:SetActive(false)
	btnCloseGetCoin = root.transform:Find("RootGetCoin/ButtonClose"):GetComponent("Button")
	btnCloseGetCoin:AddClickListener(function()	rootGetCoin.gameObject:SetActive(false) end)
	txtGetCoinOver = root.transform:Find("RootGetCoin/ImageOver/Text"):GetComponent("Text")
	inputFieldGetCoinApply = root.transform:Find("RootGetCoin/ImageApply/InputField"):GetComponent("InputField")
	btnGetCoinMax = root.transform:Find("RootGetCoin/ButtonMax"):GetComponent("Button")
	btnGetCoinMax:AddClickListener(function()	inputFieldGetCoinApply.text = txtGetCoinOver.text	end)
	btnGetCoinApply = root.transform:Find("RootGetCoin/ButtonApply"):GetComponent("Button")
	btnGetCoinApply:AddClickListener(function()		SendGetWithdrawOrderID(0) end)
	
	------- 代理提現(帳號)
	rootAccount = root.transform:Find("RootAccount"):GetComponent("RectTransform")
	rootAccount.gameObject:SetActive(false)
	btnCloseAccount = root.transform:Find("RootAccount/ButtonClose"):GetComponent("Button")
	btnCloseAccount:AddClickListener(function()	
		rootAccount.gameObject:SetActive(false) 
		--rootGetCoin.gameObject:SetActive(true)
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
	
	--------------------------------------------
	SetParentTopLayer(rootQRCodeExpand)
	SetParentTopLayer(rootRule)
	SetParentTopLayer(rootRank)
	SetParentTopLayer(rootTax)
	SetParentTopLayer(rootIncome)
	SetParentTopLayer(rootPeople)
	SetParentTopLayer(rootGetCoin)
	SetParentTopLayer(rootAccount)
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
	userdata = G.User.GetCurrentUser().UserData
	
	if(groupName ~= rootName)then
		animopen.enabled = false
	end	
	SendGetUserQACodeID()
	SendGetUserAgentExecID()
	SendGetAgentTaxBonusLogID()
	SendGetAgentTaxListID()
	SendGetUserAgentListID()
	SendGetWithdrawListID()
	SendGetUserAgentTaxLogID(os.date("%Y-%m-%d"))
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

-- 税收查看一周明细
function ShowTaxWeekList()
	rootTaxWeekUp.gameObject:SetActive(rootTaxWeekDown.gameObject.activeSelf == true)
	rootTaxWeekDown.gameObject:SetActive(rootTaxWeekDown.gameObject.activeSelf  == false)
	if(rootTaxWeekDown.gameObject.activeSelf == true)then
		local sequence = CS.DG.Tweening.DOTween.Sequence()
		sequence:Append(rootTaxWeekList:DOScaleY(1,0.2))
	else
		local sequence = CS.DG.Tweening.DOTween.Sequence()
		sequence:Append(rootTaxWeekList:DOScaleY(0,0.2))
	end
end

--查詢qrcode 1068
function SendGetUserQACodeID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetUserQACodeID, param)
end
-- region 我的推广码 QRCode 1069
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetUserQACodeRetID, function(receiveMsg)

	print("***************** receiveMsg.QRCode " .. receiveMsg.QRCode)

	if(receiveMsg.QRCode ~= "")then
		mQRCodeUrl = receiveMsg.QRCode
		local sprite =CS.Utility.GenQRCode(receiveMsg.QRCode)
		
		imgQRCode.sprite = sprite
		imgQRCodeExpand.sprite = sprite
	--[[
		G.Util.coroutine_call
		(
			function()
				local www = CS.UnityEngine.WWW(receiveMsg.QRCode)
				yield_return(www);
				 if www.error ~= nil then
					print(www.error)
				else
				
					if www.isDone and www.texture ~= nil and www.texture.width ~= 8   then
						local rect = {x=1,y = 1 ,width = www.texture.width,height = www.texture.height}
						imgQRCode.sprite=CS.UnityEngine.Sprite.Create(www.texture, CS.UnityEngine.Rect(0,0,www.texture.width,www.texture.height), CS.UnityEngine.Vector2(0, 0));
						imgQRCodeExpand.sprite=CS.UnityEngine.Sprite.Create(www.texture, CS.UnityEngine.Rect(0,0,www.texture.width,www.texture.height), CS.UnityEngine.Vector2(0, 0));
					else
					print("error")
					
					end
				end
			end
		)()
		]]--
	end
	
	txtInviteCode.text = tostring(receiveMsg.InviteCode) --邀请码
	
	--上级ID
	if receiveMsg.UpperUserID ~= 0 then
		mUpID = tostring(receiveMsg.UpperUserID)
	else
		mUpID ="";	
	end
	
	UpdateUpID()
end , groupName);

--查詢推廣員 1070
function SendGetUserAgentExecID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentExecID, param)
end

-- region 回傳推廣員 1071
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentExecRetID, function(receiveMsg)
	txtTaxCount.text = (receiveMsg.AgentTaxAmount *0.01)  .. "元" -- 今日税收
	txtIncomeCount.text = (receiveMsg.AgentTaxBonus *0.01) .. "元" --今日收益
	txtPeopleCount.text = receiveMsg.AgentTotalUser .. "人" -- 推广人数

	txtExpMax.text = "再获得".. tostring((receiveMsg.AgentAmount - receiveMsg.AgentUserAmount) *0.01) .."元收益即可升级"
	
	mAgentLevel = receiveMsg.AgentLevel -- 我的等级
	-- 目前的等级
	if(#poolRankCard >= receiveMsg.AgentLevel and receiveMsg.AgentLevel > 0)then
		imgIconCurrent.sprite = poolRankCard[receiveMsg.AgentLevel].imgIcon.sprite
		imgNameCurrent.sprite = poolRankCard[receiveMsg.AgentLevel].imgName.sprite
		imgIconShow.sprite = poolRankCard[receiveMsg.AgentLevel].imgIcon.sprite
		imgNameShow.sprite = poolRankCard[receiveMsg.AgentLevel].imgName.sprite
	end
		
	-- 下个的等级
	imgIconNext.gameObject:SetActive(receiveMsg.AgentLevel +1 <= #poolRankCard)
	imgNameNext.gameObject:SetActive(receiveMsg.AgentLevel +1 <= #poolRankCard)
	if(receiveMsg.AgentLevel +1 <= #poolRankCard)then
		imgIconNext.sprite = poolRankCard[receiveMsg.AgentLevel +1].imgIcon.sprite
		imgNameNext.sprite = poolRankCard[receiveMsg.AgentLevel +1].imgName.sprite
	end
	
	txtExpValue.text = (receiveMsg.AgentUserAmount *0.01).. "/" .. (receiveMsg.AgentAmount *0.01)
	txtShowExpValue.text = (receiveMsg.AgentUserAmount *0.01).. "/" .. (receiveMsg.AgentAmount *0.01)
	
	if(receiveMsg.AgentAmount > 0)then
		local _x = (receiveMsg.AgentUserAmount/receiveMsg.AgentAmount)
		rootShowExpValue.transform.localScale = CS.UnityEngine.Vector3(_x, 1, 1)
		rootExpValue.transform.localScale = CS.UnityEngine.Vector3(_x, 1, 1)
	end
end , groupName);

-- 查詢收益規則 1078
function SendGetAgentTaxListID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentTaxListID, nil)
end
-- region 收益規則 1079
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentTaxListRetID, function(receiveMsg)
	local list = receiveMsg.TaxList or {}
	for i=1, #list do
		local level = list[i].Level
		poolRuleCard[level].txtStep[1].text = (list[i].Stage1 *0.1) .. "%" -- 一级税收
		poolRuleCard[level].txtStep[2].text = (list[i].Stage2 *0.1) .. "%" -- 二级税收
		poolRuleCard[level].txtStep[3].text = (list[i].Stage3 *0.1) .. "%" -- 三级税收
		
		poolRankCard[level].txtStep[1].text = (list[i].Stage1 *0.1) .. "%" -- 一级税收
		poolRankCard[level].txtStep[2].text = (list[i].Stage2 *0.1) .. "%" -- 二级税收
		poolRankCard[level].txtStep[3].text = (list[i].Stage3 *0.1) .. "%" -- 三级税收
	end
end , groupName);

-- 查詢今日稅收 1072
function SendGetUserAgentTaxLogID(dateTime)
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	param.DateTime = dateTime
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentTaxLogID, param)
end
-- region 今日稅收 1073
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentTaxLogRetID, function(receiveMsg)
	txtTaxTotalCoin.text = (receiveMsg.TotalTax *0.01)
	txtTaxTotalReward.text = (receiveMsg.TotalTaxBonux *0.01)
	txtTaxTotalCount.text = receiveMsg.TotalPeople
	
	local list = receiveMsg.UserAgentTaxLog or {}
	rootTaxNoData.gameObject:SetActive(#list == 0)
	
	for i=1, #poolTax do
		poolTax[i].root.gameObject:SetActive(false)
	end
	for i=1, #list do
		CreatePrefabTax(i)
		poolTax[i].root.gameObject:SetActive(true)
		poolTax[i].txtPlayerId.text = list[i].UserID -- 玩家ID
		poolTax[i].txtTaxBase.text = (list[i].btaxamount *0.01) -- 产生税收
		poolTax[i].txtTaxOne.text = (list[i].cdtaxamount *0.01) -- 一级+二级产生税收
		poolTax[i].txtTaxGet.text = (list[i].totalTaxAmount *0.01) -- 税收收益
	end
end , groupName);

-- 查詢今日收益 1082
function SendGetAgentTaxBonusLogID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentTaxBonusLogID, param)
end
-- region 今日收益 1083
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetAgentTaxBonusLogRetID, function(receiveMsg)
	
	local list = receiveMsg.TaxBonusLogList or {}
	rootIncomeNoData.gameObject:SetActive(#list == 0)
	
	for i=1, #poolIncome do
		poolIncome[i].root.gameObject:SetActive(false)
	end
	for i=1, #list do
		CreatePrefabIncome(i)
		poolIncome[i].root.gameObject:SetActive(true)
		local LogDate = string.sub(list[i].LogDate,1,10)
		LogDate = string.gsub(LogDate,'-','/')
		poolIncome[i].txtTime.text = LogDate -- 时间
		poolIncome[i].txtTaxOne.text = (list[i].BTaxBonusAmount *0.01) -- 一级推广收益
		poolIncome[i].txtTaxTwo.text = (list[i].CDTaxBonusAmount *0.01) -- 二级+三级推广收益
		poolIncome[i].txtTaxGet.text = (list[i].TotalTaxBonusAmount *0.01) -- 收益
	end
end , groupName);

-- 查詢領取收益 1074
function SendGetWithdrawListID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetWithdrawListID, param)
end
-- region 領取收益 1075
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetWithdrawListRetID, function(receiveMsg)
	txtReceiveTotal.text = (receiveMsg.TotalTaxAmount *0.01) -- 总收益
	txtReceiveAlready.text = (receiveMsg.ConvertAmount *0.01) -- 已提现收益
	txtReceiveToday.text = (receiveMsg.PerformedAmount *0.01) -- 今日可提现
	
	txtGetCoinOver.text = (receiveMsg.RemainAmount *0.01)
	inputFieldGetCoinApply.text = (receiveMsg.RemainAmount *0.01)
	
	local list = receiveMsg.WithdrawList or {}	
	rootReceiveNoData.gameObject:SetActive(#list == 0)
	
	for i=1, #poolReceive do
		poolReceive[i].root.gameObject:SetActive(false)
	end
	for i=1, #list do
		CreatePrefabReceive(i)
		poolReceive[i].root.gameObject:SetActive(true)
		local LogDate = string.sub(list[i].LogDate,1,10)
		LogDate = string.gsub(LogDate,'-','/')
		poolReceive[i].txtTime.text = LogDate -- 时间
		poolReceive[i].txtCoin.text = (list[i].WithdrawAmount *0.01) -- 提现金额
		poolReceive[i].txtState.text = list[i].Status -- 提现状态
	end
end , groupName);

-- 查詢推廣人數 1084
function SendGetUserAgentListID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentListID, param)
end
-- region　推廣人數 1085
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetUserAgentListRetID, function(receiveMsg)
	txtTeam.text = receiveMsg.TotalUser -- 团队人数
	txtDirectly.text = receiveMsg.TempBUser -- 直属下级
	txtIncrease.text = receiveMsg.TodayNewUser -- 今日新增人数
	txtLogin.text = receiveMsg.TodayLogin -- 今日登录人数
end , groupName);

-- 申請提現
function SendGetWithdrawOrderID(uiexchange)
	if(userdata.IsBindGuest ~= true)then
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
	elseif(userdata.IsBindBankCard == false and uiexchange == 1)then
        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = uiexchange;
        CS.WindowManager.Instance:OpenWindow(openparam)
	elseif(userdata.IsBindZhifubao == false and uiexchange == 2)then
        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = uiexchange;
        CS.WindowManager.Instance:OpenWindow(openparam)	
	else
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.Amount = inputFieldGetCoinApply.text
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetWithdrawOrderID, param)
		CS.UILoadingData.Show("送出申请，等待回应",500)
	end
end
-- region　申請提現
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetWithdrawOrderRetID, function(receiveMsg)
	CS.UILoadingData.Show(receiveMsg.ChvMsg,3)
end , groupName);


-- 绑定邀请码
function SendSplitBindingInvitedCodeID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.InviteCode = ipSendInviteCode.text
		
		if G.Tool.StringHelper.IsEmpty(ipSendInviteCode.text) then
			CS.UIBubblePrompt.Show("请输入邀请码", groupName)
			return
		end
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentAddUserRelationID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 绑定邀请码結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentAddUserRelationRetID, function(receiveMsg)
	
	if receiveMsg.RetCode == 1 then
		SendGetUserQACodeID()
		--mUpID = ipSendInviteCode.text
		--UpdateUpID()
	end
	
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
	
end , groupName);

function UpdateUpID()
	if mUpID == nil or mUpID == "" then
		rootSendInviteCode.gameObject:SetActive(true)
		rootUpID.gameObject:SetActive(false)
	else
		rootSendInviteCode.gameObject:SetActive(false)
		rootUpID.gameObject:SetActive(true)	
		txtUpID.text = "上级ID: " .. mUpID	
	end
end
