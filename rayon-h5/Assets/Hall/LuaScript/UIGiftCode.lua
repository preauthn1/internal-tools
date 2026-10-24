local btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local Btn_ExchangeButton = this.transform:Find("UIBody/ExchangeButton"):GetComponent("Button")
local Ipt_Code = this.transform:Find("UIBody/InputField"):GetComponent("InputField")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local userdata = G.User.GetCurrentUser().UserData

local groupName = "UIGiftCode";

-- 初始化函数，实例化之后调用一次
function Start()

    btnClose:AddClickListener(function() EscapeWindow() end)
	
    Btn_ExchangeButton:AddClickListener(function()
		if G.Tool.StringHelper.IsEmpty(Ipt_Code.text) then
			CS.UIBubblePrompt.Show("内容不得为空白", groupName)
			return
		end   		
		
        local sendMsg = { }
		sendMsg.ExchangeSequence = Ipt_Code.text
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ExchangeCouponID, sendMsg)			
	end)
	
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

end

-- region 輸入兌換返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ExchangeCouponRetID, function(receiveMsg)
	if receiveMsg.RetCode == 1 then
		G.Service.Login.RequestUserInfoHall()	
	end

	CS.UIBubblePrompt.Show(receiveMsg.ChvMsg, groupName)	
	
end , groupName);