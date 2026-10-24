require("Module")
mmodule(...)local G = require "Global";
local groupName = "UILB";

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local userdata = G.User.GetCurrentUser().UserData
local user = G.User.GetCurrentUser();
local subList = {}
subList[1] = {}
subList[1].myName = "QRCode"
subList[2] = {}
subList[2].myName = "Reward"
subList[3] = {}
subList[3].myName = "Cash"
subList[4] = {}
subList[4].myName = "Help"

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
	
	btnCopyQRCode = root.transform:Find("UIBody/DetailQRCode/ImgQRCodeBG"):GetComponent("Button")
	btnCopyQRCode:AddClickListener(function()
        --CS.Utility.CopyText(mQRCodeUrl)
        --CS.UIBubblePrompt.Show("复制成功", groupName)
		rootQRCodeExpand.gameObject:SetActive(true)			
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
	txtUserID = root.transform:Find("UIBody/TextUserID"):GetComponent("Text")
	txtUserID.text = "ID:" .. tostring(user.UserID)
		
	-- 我的推广码
	imgQRCode = root.transform:Find("UIBody/DetailQRCode/ImgQRCodeBG/ImageQRCode"):GetComponent("Image")
	
	-- ------- 我的推广码(QRCode)
	 rootQRCodeExpand = root.transform:Find("RootQRCodeExpand"):GetComponent("RectTransform")
	 rootQRCodeExpand.gameObject:SetActive(false)
	 imgQRCodeExpandBg = root.transform:Find("RootQRCodeExpand/ImageBg"):GetComponent("RectTransform")
	 imgQRCodeExpand = root.transform:Find("RootQRCodeExpand/ImageQRCode"):GetComponent("Image")
	
	 btnQRCodeExpand = root.transform:Find("UIBody/DetailQRCode/ButtonExpand"):GetComponent("Button")
	 btnQRCodeExpand:AddClickListener(function()	rootQRCodeExpand.gameObject:SetActive(true) end)
	 txtQRCodeAbleCashOut = root.transform:Find("UIBody/DetailQRCode/info1/AbleCashOut"):GetComponent("Text")
	 txtQRCodeLBCash = root.transform:Find("UIBody/DetailQRCode/info1/LBCash"):GetComponent("Text")
	 txtQRCodeAllSuccess = root.transform:Find("UIBody/DetailQRCode/info2/TxtAllSuccess"):GetComponent("Text")
	 txtQRCodeYSuccess = root.transform:Find("UIBody/DetailQRCode/info2/TxtYSuccess"):GetComponent("Text")
	 txtQRCodeTSuccess = root.transform:Find("UIBody/DetailQRCode/info2/TxtTSuccess"):GetComponent("Text")
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

	-- 提现纪录
	txtRewardTotal = root.transform:Find("UIBody/DetailReward/ImageBg/TextTotal"):GetComponent("Text")
	prefabReward = root.transform:Find("UIBody/DetailReward/ScrollView/Content/Prefab")
	prefabReward.gameObject:SetActive(false)
	poolReward = {}	
	rootRewardNoData = root.transform:Find("UIBody/DetailReward/RootNoData"):GetComponent("RectTransform")



	------- 代理提現(金額)
	ipfCash = root.transform:Find("UIBody/DetailCash/InputFieldCash"):GetComponent("InputField")
	--txtCashAble = root.transform:Find("UIBody/DetailCash/InputFieldCash/Text"):GetComponent("Text")
	btnCashApply = root.transform:Find("UIBody/DetailCash/ButtonApply"):GetComponent("Button")
	btnCashApply:AddClickListener(function() SendSplitConvertOrderID()	end)
	txttime = root.transform:Find("UIBody/DetailCash/CountTime/TxtResult"):GetComponent("Text")

	
	------- 说明
	poolpyramid = {}
    for i=1, 5 do
		poolpyramid[i] = {}
		poolpyramid[i].bperson = root.transform:Find("UIBody/DetailHelp/Description1/B"..(1+((i-1)*2))):GetComponent("Text")
		poolpyramid[i].cperson = root.transform:Find("UIBody/DetailHelp/Description1/C"..(1+((i-1)*2))):GetComponent("Text")
		poolpyramid[i].dperson = root.transform:Find("UIBody/DetailHelp/Description1/D"..(1+((i-1)*2))):GetComponent("Text")	
		poolpyramid[i].bcoin = root.transform:Find("UIBody/DetailHelp/Description1/B"..(i*2)):GetComponent("Text")
		poolpyramid[i].ccoin = root.transform:Find("UIBody/DetailHelp/Description1/C"..(i*2)):GetComponent("Text")
		poolpyramid[i].dcoin = root.transform:Find("UIBody/DetailHelp/Description1/D"..(i*2)):GetComponent("Text")		
	end
    txtHelperTotal = root.transform:Find("UIBody/DetailHelp/Description1/TextTotal"):GetComponent("Text")
	
	
	--------------------------------------------
	SetParentTopLayer(rootQRCodeExpand)
	
	
end
	
-- 初始化函数，实例化之后调用一次
function Start()
	Show(groupName)
end
function SetParentRoot(t)
	ParentRoot=t;
end
-- 顯示介面
function Show(rootName)
	userdata = G.User.GetCurrentUser().UserData
	
	if(groupName ~= rootName)then
		animopen.enabled = false
	end


	   SendGetMyInfoID()
	   SendGetMyConvertID()
	   SendSplitGetRewardTipsID()
	   SendGetMyRewardID()
	
	-- SendGetUserAgentTaxLogID(os.date("%Y-%m-%d"))
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


function SendGetMyInfoID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()	
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyInfoID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyInfoRetID, function(receiveMsg)
	if(receiveMsg.QrURL ~= "")then		
		mQRCodeUrl = receiveMsg.QrURL
		local sprite = CS.Utility.GenQRCode(mQRCodeUrl)		
		imgQRCode.sprite = sprite
		imgQRCodeExpand.sprite = sprite		
		btnCopyQRCode.gameObject:SetActive(true)
			
	end
	
	 txtQRCodeAbleCashOut.text  = tostring(receiveMsg.RewardAmount*0.01)--可提現金額 
	 txtQRCodeLBCash.text  =tostring(receiveMsg.SplitAmount*0.01)--裂變用戶額
	 txtQRCodeAllSuccess.text  =tostring(receiveMsg.SuccessCount)--所有成功人數
	 txtQRCodeYSuccess.text  =tostring(receiveMsg.SuccessLastday)--昨日成功人數
	 txtQRCodeTSuccess.text  =  tostring(receiveMsg.SuccessToday)--今日成功人數
	 ipfCash.text  = tostring(receiveMsg.RewardAmount)--可提現金額
	 txtInviteCode.text =tostring(receiveMsg.Invited)--邀请码
    
	--上级ID
	if receiveMsg.UpperUserID ~= 0 then
		mUpID = tostring(receiveMsg.UpperUserID)
	else
		mUpID ="";
	end
	
	--mUpID = 上级ID
	UpdateUpID()
end , groupName);

function SendGetMyConvertID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()	
	param.PageIndex = 1
	param.PageSize = 30
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyConvertID, param)
end
-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyConvertRetID, function(receiveMsg)
	txtRewardTotal.text= (receiveMsg.RewardAmount*0.01)
	poolResult = receiveMsg.SplitGetMyConvert or {}	
	ChangeReward()
	
end , groupName);
--切換分頁
function ChangeReward()
	local function CreatePrefabReward(index)
		if(poolReward[index] == nil)then
			local clone = CreateClone(prefabReward)
			if(clone ~= nil)then
				clone.gameObject:SetActive(true)
				poolReward[index] = {}
				poolReward[index].root = clone.transform:GetComponent("RectTransform")
				poolReward[index].txtTime = clone.transform:Find("TextTime"):GetComponent("Text")
				poolReward[index].txtTotal = clone.transform:Find("TextTotal"):GetComponent("Text")
				
			end
		end
	end
	
		
	for i=1, #poolResult do
		CreatePrefabReward(i)
	end
	
	local list = poolResult or {}
	for i=1, #list do
		poolReward[i].txtTime.text 		= list[i].LogDate--時間
		poolReward[i].txtTotal.text 	= tostring(list[i].ConvertAmount *0.01)--日總業績
		
	end
end

-- 申請提現查询
function SendGetMyRewardID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyRewardID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetMyRewardRetID, function(receiveMsg)
    local ct = receiveMsg.MyWait
	
	if(receiveMsg.MyWait==0)then
		txttime.text="当前可提现"
		btnCashApply.gameObject:SetActive(true)
	else
        btnCashApply.gameObject:SetActive(false)
		local function showt()
			if(ct==0)then
				ParentRoot:CancelInvokeRepeat(showt)
				txttime.text="当前可提现"
				btnCashApply.gameObject:SetActive(true)
			else
				ct= ct-1
				txttime.text=G.Tool.UserHelper.sec2time(ct)
			end
		end
		ParentRoot:InvokeRepeating(showt,0 , 1)
	end
	 txtQRCodeAbleCashOut.text  = tostring(receiveMsg.MyReward*0.01)--可提現金額 
	 ipfCash.text  = tostring(receiveMsg.MyReward*0.01)--可提現金額
	
end , groupName);


-- 申請提現查询
function SendSplitConvertOrderID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.lngAmount = tonumber(ipfCash.text)
		param.ChvIPaddress = G.Tool.Debug.GetPublicIP()	
		
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitConvertOrderID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 查詢結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitConvertOrderRetID, function(receiveMsg)
	
	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)
	SendGetMyRewardID()
end , groupName);

--说明
function SendSplitGetRewardTipsID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetRewardTipsID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitGetRewardTipsRetID, function(receiveMsg)
	txtHelperTotal.text=(receiveMsg.ReturenAmount*0.01)
	poolTip = receiveMsg.SplitGetRewardTips or {}	
	ChangePyramid()
	
end , groupName);

function ChangePyramid()
	
	local list = poolTip or {}
	for i=1, #list do

		print(list[i].RecruitB,list[i].RecruitC,list[i].RecruitD)
		print(list[i].RewardB,list[i].RewardC,list[i].RewardD)
		poolpyramid[i].bperson.text		= tostring(list[i].RecruitB )--B人數
		poolpyramid[i].cperson.text		= tostring(list[i].RecruitC )--C人數
		poolpyramid[i].dperson.text		= tostring(list[i].RecruitD )--D人數
		poolpyramid[i].bcoin.text 		= tostring(list[i].RewardB *0.01)--B 金額
		poolpyramid[i].ccoin.text 		= tostring(list[i].RewardC *0.01)--C 金額
		poolpyramid[i].dcoin.text 		= tostring(list[i].RewardD *0.01)--D 金額

		
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
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SplitAddUserRelationID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 绑定邀请码結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SplitAddUserRelationRetID, function(receiveMsg)
	
	if receiveMsg.RetCode == 1 then
		SendGetMyInfoID()
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


