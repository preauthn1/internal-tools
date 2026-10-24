local btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local imgTitle = this.transform:Find("UIBody/Background/ImageTitle"):GetComponent("Image")
local txtClose = this.transform:Find("UIBody/TextClose"):GetComponent("Text")

local btnSign = this.transform:Find("UIBody/SubList/Grid/ButtonSign"):GetComponent("Button")
local imgSignSelect = this.transform:Find("UIBody/SubList/Grid/ButtonSign/ImageSelect"):GetComponent("Image")
local scrollSign = this.transform:Find("UIBody/DetailSign/ScrollView"):GetComponent("ScrollRect")
local rectSign = this.transform:Find("UIBody/DetailSign"):GetComponent("RectTransform")
local txtSignLimit = this.transform:Find("UIBody/DetailSign/chatbg/TextLimit"):GetComponent("Text")
local txtSignReceive = this.transform:Find("UIBody/DetailSign/TextReceive"):GetComponent("Text")
--local imgElement = this.transform:Find("UIBody/DetailSign/ImageElement"):GetComponent("Image")

local rectAnimGetReward = this.transform:Find("UIBody/Animation"):GetComponent("RectTransform")
local perfabRewardStart = this.transform:Find("UIBody/Animation/AniStart").gameObject
local perfabRewardIdle = this.transform:Find("UIBody/Animation/AniIdle").gameObject

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local userdata = G.User.GetCurrentUser().UserData
local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

local dailyPool = {}
local dailyList = {}
local imgCount = {}

local groupName = "UINotice";
-- 初始化函数，实例化之后调用一次
function Start()

	btnClose:AddClickListener(function()
		CS.WindowManager.Instance:CloseWindow(groupName, false)
	end)
	btnSign:AddClickListener(function() ChangeSubList(1) end)
	
	txtSignLimit.text = "游戏签到免费为您提供金币，签到为"
	txtSignLimit.text = txtSignLimit.text .. G.Service.LoginBonus.DayMax .. "天签到，请先绑定真实有效帐号"
	if G.Service.LoginBonus.GiftCash > 0 then
		txtSignLimit.text = txtSignLimit.text .. "，并充值".. G.Service.LoginBonus.GiftCashToCash .."金币"
	end
	txtSignLimit.text = txtSignLimit.text .. "，完成后即可领取免费金币"
end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
	-- 回收当前页面注册的消息
	G.NetQ.RemoveMsgHookByGroup(groupName)

	-- 清空当前界面的浮动提示消息
	CS.UIBubblePrompt.ClearPrompt(groupName)
end  

function EscapeWindow()
	CS.WindowManager.Instance:CloseWindow(groupName, false)
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	for i=1,10 do
		imgCount[i] = this.transform:Find("UIBody/DetailSign/ImageCount"..i):GetComponent("Image")
	end
		imgCount[0.5] = this.transform:Find("UIBody/DetailSign/ImageCount0.5"):GetComponent("Image")
	
	local rowCount = 4
	local childCount = this.transform:Find("UIBody/DetailSign/ScrollView/Grid").transform.childCount
	for i=1,childCount do
		for j=1,rowCount do
			local target = this.transform:Find("UIBody/DetailSign/ScrollView/Grid/Content"..i.."/Item"..j)
			local index = #dailyPool +1
			dailyPool[index] = {}
			dailyPool[index].root = target.transform:GetComponent("RectTransform")
			dailyPool[index].btn = target.transform:GetComponent("Button")
			dailyPool[index].img = target.transform:GetComponent("Image")
			dailyPool[index].day = target.transform:Find("ImageDay"):GetComponent("Image")
			dailyPool[index].receive = target.transform:Find("ImageReceive"):GetComponent("Image")
			dailyPool[index].unreceive = target.transform:Find("ImageUnreceive"):GetComponent("Image")
			dailyPool[index].count = target.transform:Find("ImageCount"):GetComponent("Image")
			dailyPool[index].root.gameObject:SetActive(false)
		end
	end
	
	--scrollSign.enabled = G.Service.LoginBonus.DayMax > rowCount*2
	btnSign.gameObject:SetActive(G.Service.LoginBonus.RetCode ~= -1)
	for i=1,G.Service.LoginBonus.DayMax do
		local index = i
		UpdateSignItem(index)
		dailyPool[i].root.gameObject:SetActive(true)
		dailyPool[i].btn:AddClickListener(function() ClickApply(index) end)
		local count = G.Service.LoginBonus.ContentList[i] *0.01
		if(imgCount[count] ~= nil)then
			dailyPool[i].count.sprite = imgCount[count].sprite
		else
			dailyPool[i].count.gameObject:SetActive(false)
		end
	end
	
	local startSubIndex = GetStartSubIndex()
	ChangeSubList(startSubIndex)
    --企劃說應該永遠都會存在,所以不會關閉,也沒有該美術視窗
	--txtClose.gameObject:SetActive(startSubIndex == 0)
end

function GetStartSubIndex()
	local index = 0
	if G.Service.LoginBonus.RetCode ~= -1 then
		index = 1
	end
	
	return index
end

function GetSubCount()
	
	childCount = this.transform:Find("UIBody/SubList/Grid").transform.childCount
	
	return childCount
end

-- 切换
function ChangeSubList(index)
	imgSignSelect.gameObject:SetActive(index==1)
	rectSign.gameObject:SetActive(index==1)
	
	if index == 1 then
		ShowSign()
	end
end


function ShowSign()
	-- 增加画面的丰富度
	--imgElement.gameObject:SetActive(GetSubCount() < 5)
	
    print("<color=yellow>retcode:</color>".. tostring(G.Service.LoginBonus.RetCode));
	scrollSign.gameObject:SetActive(G.Service.LoginBonus.RetCode ~= -1 and G.Service.LoginBonus.RetCode ~= -2 )
	--先暫時呈現7個領完,等美術想到解決方案會補入新圖
    --txtSignLimit.gameObject:SetActive(G.Service.LoginBonus.RetCode ~= -1)
	--txtSignReceive.gameObject:SetActive(G.Service.LoginBonus.RetCode == -1 or G.Service.LoginBonus.RetCode == -2 )
	
	if G.Service.LoginBonus.RetCode == -2 then
		txtSignReceive.text = "奖励已全数领取，敬请期待更多好礼活动唷~"
	end
end

function ClickApply(index)
	if G.Service.LoginBonus.LoginTimes +1 == index and G.Service.LoginBonus.RetCode == 2 then
		G.Service.LoginBonus.SendApply(ApplyCallback)
	elseif G.Service.LoginBonus.RetCode == 0 then
		CS.UIBubblePrompt.Show("今日奖励已领取，请隔日再来", groupName)
	-- 帐号未绑定
	elseif G.Service.LoginBonus.RetCode == -4 then
		ShowBindingAccount()
	-- 充值金额不足
	elseif G.Service.LoginBonus.RetCode == -3 then
		ShowRecharge()
	end
end

function ShowBindingAccount()
	G.Util.coroutine_call( function()
		local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "需要绑定用户信息，确认绑定吗?", "确定", "取消");
		if boxResult == CS.MessageBoxResult.OK then
			local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
			openparam.WindowData = 1
			CS.WindowManager.Instance:OpenWindow(openparam)
		end
	end )();
end

function ShowRecharge()
	G.Util.coroutine_call( function()
		local msg = "活动需要充值".. G.Service.LoginBonus.GiftCashToCash .."金币才可参加，要立即充值吗？"
		local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", msg, "确定", "取消");
		if boxResult == CS.MessageBoxResult.OK then
			local openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
			CS.WindowManager.Instance:OpenWindow(openparam)
		end
	end )();
end


function ApplyCallback()
	for i=1,#dailyPool do
		UpdateSignItem(i)
	end
	if G.Service.LoginBonus.RetCode == 1 then
		ShowAnimGetReward()
	end
end

function UpdateSignItem(index)
	
	dailyPool[index].receive.gameObject:SetActive(G.Service.LoginBonus.LoginTimes >= index)
	--dailyPool[index].count.gameObject:SetActive(G.Service.LoginBonus.LoginTimes < index)
	local showUnreceive = (G.Service.LoginBonus.LoginTimes +1 == index and G.Service.LoginBonus.RetCode == 2)
	dailyPool[index].unreceive.gameObject:SetActive(showUnreceive)
    
	dailyPool[index].btn.enabled = (G.Service.LoginBonus.LoginTimes < index)
	if showUnreceive == true then
		ShowAnimUnreceive(index)
	end
	
	if G.Service.LoginBonus.LoginTimes >= index then
		dailyPool[index].img.color = CS.UnityEngine.Color(0.45,0.45,0.45,1)
		dailyPool[index].day.color = CS.UnityEngine.Color(1,1,1,0.2)
	else
		dailyPool[index].img.color = CS.UnityEngine.Color(1,1,1,1)
		dailyPool[index].day.color = CS.UnityEngine.Color(1,1,1,1)
	end
end

function ShowAnimGetReward()
	rectAnimGetReward.gameObject:SetActive(true)
	local cloneStart = CS.UnityEngine.GameObject.Instantiate(perfabRewardStart).transform
	cloneStart:SetParent(rectAnimGetReward.transform)
	cloneStart.localScale = CS.UnityEngine.Vector3.one
	cloneStart.localPosition = CS.UnityEngine.Vector3(0,-50,0)
	cloneStart.gameObject:SetActive(true)
	local cloneIdle = CS.UnityEngine.GameObject.Instantiate(perfabRewardIdle).transform
	cloneIdle:SetParent(rectAnimGetReward.transform)
	cloneIdle.localScale = CS.UnityEngine.Vector3.one
	cloneIdle.localPosition = CS.UnityEngine.Vector3(0,-50,0)
	cloneIdle.gameObject:SetActive(false)
	
	--G.Tool.Sound.PlayGameSound("adVS")
	local sequence = CS.DG.Tweening.DOTween.Sequence()
    sequence:AppendInterval(1.0)
	sequence:AppendCallback( function()
        cloneStart.gameObject:SetActive(false)
		CS.UnityEngine.GameObject.Destroy(cloneStart.gameObject)
		cloneIdle.gameObject:SetActive(true)
    end)
	sequence:AppendInterval(1.0)
	sequence:AppendCallback( function()
		rectAnimGetReward.gameObject:SetActive(false)
        cloneIdle.gameObject:SetActive(false)
		CS.UnityEngine.GameObject.Destroy(cloneIdle.gameObject)
    end)
end

function ShowAnimUnreceive(index)
	local targetColor = dailyPool[index].unreceive.color
	sequenceUnreceive = CS.DG.Tweening.DOTween.Sequence()
	sequenceUnreceive:SetLoops(-1,CS.DG.Tweening.LoopType.Yoyo)
    sequenceUnreceive:Append(dailyPool[index].unreceive:DOColor(CS.UnityEngine.Color(targetColor.r, targetColor.g, targetColor.g, 0), 1))
end