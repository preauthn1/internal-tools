local G = require "Global";
local groupName = "UIProxy";

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local userdata = G.User.GetCurrentUser().UserData
local subList = {}
subList[1] = {}
subList[1].myName = "Promoter"

subList[2] = {}
subList[2].myName = "2"

subList[3] = {}
subList[3].myName = "LB"

subList[4] = {}
--subList[4].myName = "4"
subList[4].myName = "NetProfit"

subList[5] = {}
subList[5].myName = "WXD"

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	CS.UILoadingData.Hide()

	btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
    btnClose:AddClickListener(function() EscapeWindow() end)
	
	rootTopLayer = this.transform:Find("TopLayer"):GetComponent("RectTransform")
	
	for i=1, #subList do
		local tempIndex = i
		local myName = subList[i].myName
		subList[i].rootDetail = this.transform:Find("UIBody/UI" .. myName):GetComponent("RectTransform")
		subList[i].rootDetail.gameObject:SetActive(false)
		subList[i].btn = this.transform:Find("UIBody/SubList/Grid/Button" .. myName):GetComponent("Button")
		subList[i].btn:AddClickListener(
			function()
				ChangeSubList(tempIndex) 
			end)
		subList[i].btn.gameObject:SetActive(false)
		subList[i].imgSelect = this.transform:Find("UIBody/SubList/Grid/Button" .. myName .. "/ImageSelect"):GetComponent("Image")
		subList[i].imgSelect.gameObject:SetActive(false)
		if(myName ~= tostring(i)) then
			-- if(myName~="NetProfit")then
				subList[i].uiscript = G.RequireG("UI"..myName)
				subList[i].uiscript.Init(subList[i].rootDetail,rootTopLayer)
			-- end
		end
	end
		
	if(G.User.payinfo.UIClub2AgentMode1 == true)then
		ChangeSubList(1)
	elseif(G.User.payinfo.UIClub2AgentMode2 == true)then
		ChangeSubList(2)
	elseif(G.User.payinfo.UIClub2AgentMode3 == true)then
		ChangeSubList(3)
	elseif(G.User.payinfo.UIClub2AgentMode4 == true)then
		ChangeSubList(4)		
	elseif(G.User.payinfo.UIClub2AgentMode5 == true)then
		ChangeSubList(5)
	end
	
	-- 進入動畫
	imgLogo = this.transform:Find("ButtonComeIn/ImageLogo"):GetComponent("Image")
	imgComeIn = this.transform:Find("ButtonComeIn/ImageComeIn"):GetComponent("Image")
	rootDoorL = this.transform:Find("ButtonComeIn/ImageDoorL"):GetComponent("RectTransform")
	rootDoorR = this.transform:Find("ButtonComeIn/ImageDoorR"):GetComponent("RectTransform")
	btnComeIn = this.transform:Find("ButtonComeIn"):GetComponent("Button")
	btnComeIn.gameObject:SetActive(true)
	btnComeIn:AddClickListener(function() 
		local sequence = CS.DG.Tweening.DOTween.Sequence()
		sequence:AppendCallback(function()
			imgComeIn.gameObject:SetActive(false)
		end)
		sequence:Append(imgLogo:DOColor(CS.UnityEngine.Color(imgLogo.color.r,imgLogo.color.g,imgLogo.color.b,0), 0.2))
		sequence:Append(rootDoorL:DOAnchorPosX(-1000, 0.5))
		sequence:Join(rootDoorR:DOAnchorPosX(1000, 0.5))
		sequence:OnComplete(function()
			btnComeIn.gameObject:SetActive(false)
		end)	
	end)
	
end
	
-- 初始化函数，实例化之后调用一次
function Start()
	SendUserAgentUserModeGetID()
end
 
-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
	
    CS.Utility.CallGC();
	--[[ 清除lua腳本
	for k, _ in pairs(package.loaded) do
		for i=1, #subList do
			if (string.find(k, subList[i].myName)) then
				package.loaded[k] = nil
				package.preload[k] = nil
			end
		end
    end]]
end

-- 離開介面
function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 切換介面
function ChangeSubList(index)
	-- 净利润暂时不要画面所以先将净利润导去全民
	--[[
	if(index==4)then
		subList[1].rootDetail.gameObject:SetActive(true)
		subList[4].rootDetail.gameObject:SetActive(false)
		subList[1].uiscript.Show()
	end
	]]--
	for i=1, #subList do
		-- 净利润暂时不要画面所以先将净利润导去全民
		--if((index==4)and ((i==1)or (i==4) ))then
		
		--else
			subList[i].rootDetail.gameObject:SetActive(index == i)
		--end		
		subList[i].imgSelect.gameObject:SetActive(index == i)
	end
	if(subList[index].uiscript ~= nil)then
		subList[index].uiscript.SetParentRoot(this)
		subList[index].uiscript.Show()
	end
end

--玩家資訊
function SendUserAgentUserModeGetID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentUserModeGetID, param)
end
-- region 玩家資訊
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserAgentUserModeGetRetID, function(receiveMsg)
	local mode = {}
	mode[1] = G.User.payinfo.UIClub2AgentMode1 == true or receiveMsg.Mode1 == 1
	mode[2] = G.User.payinfo.UIClub2AgentMode2 == true or receiveMsg.Mode2 == 1
	mode[3] = G.User.payinfo.UIClub2AgentMode3 == true or receiveMsg.Mode3 == 1
	mode[4] = G.User.payinfo.UIClub2AgentMode4 == true or receiveMsg.Mode4 == 1
	mode[5] = G.User.payinfo.UIClub2AgentMode5 == true or receiveMsg.Mode5 == 1
	
	local count = 0
	for i=1, #mode do
		if(mode[i] == true)then
			count = count +1
		end
	end
	
	for i=1, #subList do
	print("555555 ", i)
		if(count == 1)then
			subList[i].btn.gameObject:SetActive(false)
		else
			subList[i].btn.gameObject:SetActive(mode[i])
		end
	end
end , groupName);