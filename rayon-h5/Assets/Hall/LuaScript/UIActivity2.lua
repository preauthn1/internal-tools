local recTran__TopNumber = this.transform:Find("UIBody/_TopNumber"):GetComponent("RectTransform")
local txt__Number0_1 = this.transform:Find("UIBody/_TopNumber/_Number0_1"):GetComponent("Text")
local txt__Number0_2 = this.transform:Find("UIBody/_TopNumber/_Number0_2"):GetComponent("Text")

local txt__Number2_1 = this.transform:Find("UIBody/_Number2_1"):GetComponent("Text")
local txt__Number2_2 = this.transform:Find("UIBody/_Number2_2"):GetComponent("Text")
local txt__Number2_4 = this.transform:Find("UIBody/_Number2_4"):GetComponent("Text")
local txt__Number2_3 = this.transform:Find("UIBody/_Number2_3"):GetComponent("Text")

local img__BoxOpen = this.transform:Find("UIBody/_BoxOpen"):GetComponent("Image")
local btn__BoxNotOpen = this.transform:Find("UIBody/_BoxNotOpen"):GetComponent("Button")
local txt__Time = this.transform:Find("UIBody/_Time"):GetComponent("Text")
local btn__CloseBtn = this.transform:Find("UIBody/_CloseBtn"):GetComponent("Button")
local btn__JoinBtn = this.transform:Find("UIBody/_JoinBtn"):GetComponent("Button")

local recTran__Bar = this.transform:Find("UIBody/_Acc/ImageExp/_Bar"):GetComponent("RectTransform")
local txt__Number = this.transform:Find("UIBody/_Acc/_Number"):GetComponent("Text")

local txt__JoinNumber = this.transform:Find("UIBody/_Number/_JoinNumber"):GetComponent("Text")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local groupName = "UIActivity2"

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

local seq
local data = nil

-- 初始化函数，实例化之后调用一次
function Start()

    btn__CloseBtn:AddClickListener(EscapeWindow)

    btn__JoinBtn:AddClickListener
    (
    function()
        EscapeWindow()

        local openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )
	
	btn__BoxNotOpen:AddClickListener
    (
    function()
		CS.UILoadingData.Show("领取中,请稍后...", 10,
		(function()
			CS.UIBubblePrompt.Show("领取失败", groupName)
		end))
		
        local sendMsg = {}
		sendMsg.UserID = G.User.GetCurrentUser().UserID
		sendMsg.EventID = data.FirstDepositUserRecord.EventID
		sendMsg.IPAddress = G.GetLocalIP()
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.FirstDepositCollectID, sendMsg)
    end
    )
	
end 

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	recTran__TopNumber.gameObject:SetActive(false)
	
	data = G.Service.FirstRechargeRebate.GetData()
	
	-- 單筆充值限制 (單位:分)
	local DepositLimit = math.modf(data.FirstDepositInfo.DepositLimit / 100)
	-- 返利金額 (單位:分)	 
	local ReturnAmount = math.modf(data.FirstDepositInfo.ReturnAmount / 100)
	-- 需贏得金額 (單位:分)
	local WinAmount = math.modf(data.FirstDepositInfo.WinAmount / 100)
	
	txt__Number0_1.text = tostring(DepositLimit)
	txt__Number0_2.text = tostring(ReturnAmount)
	txt__Number2_1.text = tostring(DepositLimit)
	txt__Number2_2.text = tostring(ReturnAmount)
	txt__Number2_3.text = tostring(WinAmount)
	txt__Number2_4.text = tostring(ReturnAmount)
	txt__Time.text = data.FirstDepositInfo.EndTime
	txt__JoinNumber.text = tostring(0)
	
	recTran__TopNumber.gameObject:SetActive(true)
	
	now = math.modf(data.FirstDepositUserRecord.WinAmount / 100)
	target = WinAmount
	txt__Number.text = tostring(now) .. "/" .. tostring(target)
	
	-- bar
	local value = now / target
	local width = recTran__Bar.rect.width
	if value < 1 then
		recTran__Bar.transform.localPosition = CS.UnityEngine.Vector3(-1.5 * width, 0, 0) + CS.UnityEngine.Vector3(value * width , 0, 0)
	else
		recTran__Bar.transform.localPosition = CS.UnityEngine.Vector3(-width * 0.5, 0, 0)	
	end
		
	local setAmount = function(temp)
		txt__JoinNumber.text = math.modf(temp)
	end	
	seq = CS.DG.Tweening.DOTween.Sequence()
	seq:Append(CS.DG.Tweening.DOTween.To(setAmount, 100, data.FirstDepositInfo.CurrentUserCount, 0.5));
	
	-- 还不可以领
	if data.FirstDepositUserRecord.Flag  == 0 then
		btn__BoxNotOpen.gameObject:SetActive(true)
		img__BoxOpen.gameObject:SetActive(false)
		btn__BoxNotOpen.interactable = false
	-- 可以领
	elseif data.FirstDepositUserRecord.Flag  == 1 then
		btn__BoxNotOpen.gameObject:SetActive(true)
		img__BoxOpen.gameObject:SetActive(false)
		btn__BoxNotOpen.interactable = true
	-- 已領取
	else
		btn__BoxNotOpen.gameObject:SetActive(false)
		img__BoxOpen.gameObject:SetActive(true)
	end	
end

-- 领取结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.FirstDepositCollectRetID, function(receiveMsg)	
	CS.UILoadingData.Hide()
	
	if receiveMsg.RetCode == 1 then
		local ReturnAmount = math.modf(data.FirstDepositInfo.ReturnAmount / 100)
		
		G.Util.coroutine_call(function()
			CS.UIMessageBox.Show(CS.MessageBoxStyle.OK, "", "恭喜您获得返利" .. tostring(ReturnAmount) .. "元!!", "确定");
		end	)()	
		
		img__BoxOpen.gameObject:SetActive(true)
		btn__BoxNotOpen.gameObject:SetActive(false)
	else
		CS.UIBubblePrompt.Show(receiveMsg.RetMessage, groupName)	
	end
end , groupName);

function OnDestroy()
	if seq ~= nil then
		seq:Kill()
	end
	
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(windowName)
end

	