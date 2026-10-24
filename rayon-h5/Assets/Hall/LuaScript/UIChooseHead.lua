
local Btn_Head = this.transform:Find("UIBody/_Btn/_BtnHead"):GetComponent("Button")
local HeadSelect = this.transform:Find("UIBody/_Btn/_BtnHead/ImageSelect"):GetComponent("Image")
local Btn_Frame = this.transform:Find("UIBody/_Btn/_BtnFrame"):GetComponent("Button")
local FrameSelect = this.transform:Find("UIBody/_Btn/_BtnFrame/ImageSelect"):GetComponent("Image")
local Btn_Cancel = this.transform:Find("UIBody/_Cancel"):GetComponent("Button")

local Img_Frame = this.transform:Find("UIBody/_frame"):GetComponent("Image")
local Img_Head = this.transform:Find("UIBody/_head"):GetComponent("Image")
local Img_Vip = this.transform:Find("UIBody/_vip"):GetComponent("Image")

local HeadPanel = this.transform:Find("UIBody/_HeadScr"):GetComponent("RectTransform")
local HeadAttach = this.transform:Find("UIBody/_HeadScr/Content"):GetComponent("RectTransform")
local HeadPrefab = this.transform:Find("UIBody/_HeadScr/Content/Item"):GetComponent("RectTransform")

local FramePanel = this.transform:Find("UIBody/_FrameScr"):GetComponent("RectTransform")
local FrameAttach = this.transform:Find("UIBody/_FrameScr/Content"):GetComponent("RectTransform")
local FramePrefab = this.transform:Find("UIBody/_FrameScr/Content/Item"):GetComponent("RectTransform")

local btn__vip = this.transform:Find("UIBody/_vip"):GetComponent("Button")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";

local userdata = G.User.GetCurrentUser().UserData

local groupName = "UIChooseHead";

local HeadData = {}
local FrameData = {}

local currentHeadIndex
local currentFrameIndex = 0
local currentVipLevel = 0

-- 初始化函数，实例化之后调用一次
function Start()
    Btn_Head:AddClickListener
    (
    function()
		ShowType(1)
    end
    )
	
    Btn_Frame:AddClickListener
    (
    function()
		ShowType(2)
    end
    )	

    Btn_Cancel:AddClickListener
    (
    function()
		EscapeWindow()
    end
    )
	
	btn__vip:AddClickListener
    (
    function()
		EscapeWindow()
		
        local openparam = CS.WindowNodeInitParam("UIVipInfo", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)		
    end
    )
	
	-- 頭相框是0表示是舊機制,關掉VIP圖是和換框按鈕
	if userdata.AvatarFrame == 0 then
		Btn_Frame.gameObject:SetActive(false)
		Img_Vip.gameObject:SetActive(false)
	end
end

function ShowType(Type)
	if Type == 1 then
		HeadPanel.gameObject:SetActive(true)
		FramePanel.gameObject:SetActive(false)
	else
		HeadPanel.gameObject:SetActive(false)
		FramePanel.gameObject:SetActive(true)
	end
	
	HeadSelect.gameObject:SetActive(Type == 1)
	FrameSelect.gameObject:SetActive(Type == 2)	
end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)
end 

-- region 修改头像回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ModifyFaceRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        userdata.HeadID = currentHeadIndex
		
        CS.UIBubblePrompt.Show("修改头像成功", "UIChooseHead");
        --CS.WindowManager.Instance:CloseWindow("UIChooseHeadIcon", false);
		
		UpdateHead()
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- region 修改外框回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ModifyAvatarFrameRetID, function(receiveMsg)
    if receiveMsg.RetCode == 1 then
        userdata.AvatarFrame = currentFrameIndex
		
        CS.UIBubblePrompt.Show("修改外框成功", "UIChooseHead");
        --CS.WindowManager.Instance:CloseWindow("UIChooseHeadIcon", false);
		
		UpdateFrame()
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	
	-- Vip等级
	currentVipLevel = userdata.UserLevel	
	Img_Vip:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(currentVipLevel))
	
	-- 頭像
    currentHeadIndex = userdata.HeadID
	
    for i = 1, 16, 1 do	
		local clone = CS.UnityEngine.GameObject.Instantiate(HeadPrefab)
		clone.transform:SetParent(HeadAttach.transform)
		clone.transform.localPosition = CS.UnityEngine.Vector3.zero
		clone.transform.localScale = CS.UnityEngine.Vector3.one
		
		local data = {}
		data.index = i
		data.icon = clone.transform:Find("icon"):GetComponent("Image")
		data.use = clone.transform:Find("use"):GetComponent("Image")
        table.insert(HeadData, data)
		
		data.icon:ResetSpriteByName(CS.AppDefine.HallNumber, tostring(data.index))
		
        local btnIcon = clone:GetComponent("Button")
        btnIcon:AddClickListener
        (
        function()		
			if currentHeadIndex == i then
				return
			end

            currentHeadIndex = i;		

			local sendMsg = { }
			sendMsg.FaceID = currentHeadIndex
			G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ModifyFaceID, sendMsg)
        end
        )
    end
	
	HeadPrefab.gameObject:SetActive(false)
	
	UpdateHead()
	
	-- 外框
	currentFrameIndex = userdata.AvatarFrame
	
	local vipAllData = G.Service.Vip.GetVipAllData()	
	for i, vipData in ipairs(vipAllData) do
		local clone = CS.UnityEngine.GameObject.Instantiate(FramePrefab)
		clone.transform:SetParent(FrameAttach.transform)
		clone.transform.localPosition = CS.UnityEngine.Vector3.zero
		clone.transform.localScale = CS.UnityEngine.Vector3.one
		
		local data = {}
		data.index = vipData.lv
		data.icon = clone.transform:Find("icon"):GetComponent("Image")
		data.use = clone.transform:Find("use"):GetComponent("Image")
		data.lock = clone.transform:Find("lock"):GetComponent("Image")
		data.vip = clone.transform:Find("vip"):GetComponent("Image")
        table.insert(FrameData, data)
		
		G.Service.Vip.SetFrame(data.icon, data.index)
		data.vip:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(data.index))
		
        local btnIcon = clone:GetComponent("Button")
        btnIcon:AddClickListener
        (
        function()
			if data.index > currentVipLevel then
				CS.UIBubblePrompt.Show("VIP等级达到" .. tostring(data.index) .. "级可获得", "UIChooseHead");
				return
			end		
		
			if currentFrameIndex == i then
				return
			end
		
            currentFrameIndex = i;		

			local sendMsg = { }
			sendMsg.UserID = G.User.GetCurrentUserID()
			sendMsg.AvatarFrameID = currentFrameIndex
			G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ModifyAvatarFrameID, sendMsg)
        end
        )
    end
	
	FramePrefab.gameObject:SetActive(false)

	UpdateFrame()
	
	ShowType(1)
end

function UpdateHead()
	Img_Head:ResetSpriteByName(CS.AppDefine.HallNumber, tostring(currentHeadIndex))
	
    for i, data in ipairs(HeadData) do
		if data.index == currentHeadIndex then
			data.use.gameObject:SetActive(true)
		else
			data.use.gameObject:SetActive(false)
		end
	end
end

function UpdateFrame()
	G.Service.Vip.SetFrame(Img_Frame, currentFrameIndex)
	
    for i, data in ipairs(FrameData) do
		if data.index == currentFrameIndex then
			data.use.gameObject:SetActive(true)
		else
			data.use.gameObject:SetActive(false)
		end
		
		if data.index > currentVipLevel then
			data.lock.gameObject:SetActive(true)
		else
			data.lock.gameObject:SetActive(false)		
		end
	end
end

function SetCurrentHead(headIcon)
    Img_ChoosedMark.transform:SetParent(headIcon)
    Img_ChoosedMark.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
    Img_ChoosedMark.transform.localPosition = CS.UnityEngine.Vector3(0, 0, 0)
end

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIChooseHead", false);
end