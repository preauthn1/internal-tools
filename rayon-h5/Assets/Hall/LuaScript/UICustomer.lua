local Btn_Close = this.transform:Find("UIBody/_Close"):GetComponent("Button")
local Btn_Submit = this.transform:Find("UIBody/_Submit"):GetComponent("Button")
local Ipt_Content = this.transform:Find("UIBody/_Content"):GetComponent("InputField")
local Ipt_Email = this.transform:Find("UIBody/_Content"):GetComponent("InputField")
local Drop_QuestionType = this.transform:Find("UIBody/_QuestionType"):GetComponent("Dropdown")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"

-- 初始化函数，实例化之后调用一次
function Start()	
	Btn_Close:AddClickListener(function()
        CS.WindowManager.Instance:CloseWindow("UICustomer", false);
	end
	)	
	
	Drop_QuestionType.value = 3
	Btn_Submit:AddClickListener( function()
		if G.Tool.StringHelper.IsEmpty(Ipt_Content.text) then
            CS.UIBubblePrompt.Show("请输入消息内容", groupName)
            return
        end
		if G.Tool.StringHelper.IsEmpty(Ipt_Email.text) then
            CS.UIBubblePrompt.Show("请输入邮箱号码", groupName)
            return
        end
		
		StartUp.SendLoginCustomer(Ipt_Content.text,Ipt_Email.text ,Drop_QuestionType.value + 1)
        CS.WindowManager.Instance:CloseWindow("UICustomer", false);
    end
    )
	
end