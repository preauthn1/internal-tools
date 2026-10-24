require("Module")
mmodule(...)local G = require "Global";
local groupName = "UINetProfit";

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local userdata = G.User.GetCurrentUser().UserData
local subList = {}
subList[1] = {}
subList[1].myName = "QRCode"

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
	
	rootUserInfo = root.transform:Find("RootQRCodeExpand"):GetComponent("RectTransform")
	--------------------------------------------
	SetParentTopLayer(rootUserInfo)
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
	
	SendGetUserQACodeID()
	
end
 
 --查詢qrcode
function SendGetUserQACodeID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.NetProfitAgentGetAgentInfoID, param)
end
-- region 我的推广码 QRCode
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.NetProfitAgentGetAgentInfoRetID, function(receiveMsg)

	print("***************** receiveMsg.QRCode " .. receiveMsg.QRCode)

	if(receiveMsg.QRCode ~= "")then
		mQRCodeUrl = receiveMsg.QRCode
		local sprite =CS.Utility.GenQRCode(receiveMsg.QRCode)
		
		imgQRCode.sprite = sprite
		imgQRCodeExpand.sprite = sprite
	end
	
	txtInviteCode.text = tostring(receiveMsg.InviteCode) --邀请码
	
	--上级ID
	if receiveMsg.UpperUserID ~= 0 then
		mUpID = tostring(receiveMsg.UpperUserID)
	else
		mUpID =""	
	end
	
	UpdateUpID()
end , groupName);
 
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


-- 绑定邀请码
function SendSplitBindingInvitedCodeID()
	
		local param = {}
		param.UserID = G.User.GetCurrentUserID()
		param.InviteCode = ipSendInviteCode.text
		
		if G.Tool.StringHelper.IsEmpty(ipSendInviteCode.text) then
			CS.UIBubblePrompt.Show("请输入邀请码", groupName)
			return
		end
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.NetProfitAgentAddUserRelationID, param)
		--CS.UILoadingData.Show("送出申请，等待回应", 10)
end

-- 绑定邀请码結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.NetProfitAgentAddUserRelationRetID, function(receiveMsg)
	
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