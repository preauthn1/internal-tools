local btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local btnSubmit = this.transform:Find("UIBody/ButtonSubmit"):GetComponent("Button")
local iptContent = this.transform:Find("UIBody/InputField"):GetComponent("InputField")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"
local groupName = "UISendMailToGM";

-- 初始化函数，实例化之后调用一次
function Start()	
	
    btnClose:AddClickListener(function() EscapeWindow() end)
	btnSubmit:AddClickListener(function() SendSubmit() end)
	
	iptContent.placeholder.text = G.User.GetCurrentUser().UserData.QuestionContactTip
end


function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

function SendSubmit()
	if G.Tool.StringHelper.IsEmpty(iptContent.text) then
		CS.UIBubblePrompt.Show("请输入消息内容", groupName)
		return
	end
	
	local sendMsg = { }
	sendMsg.Msg = iptContent.text
	sendMsg.QuestionType = 3
	sendMsg.PlatformID = G.GameDefine.GamePlatformID 
	sendMsg.AreaID = G.GameDefine.GameSiteID
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ContactServiceID, sendMsg)
end


-- region 聯絡客服結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ContactServiceRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("讯息已发送完成", "UIHall");
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
	
	EscapeWindow()
	
end , groupName);
