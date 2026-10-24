local Btn_Common = this.transform:Find("UIBody/Buttons/_Common"):GetComponent("Button")
local Btn_OK = this.transform:Find("UIBody/Buttons/_OK"):GetComponent("Button")
local Btn_Cancel = this.transform:Find("UIBody/Buttons/_Cancel"):GetComponent("Button")
local Text_Content = this.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"

-- 初始化函数，实例化之后调用一次
function Start()
	-- 顯示錯誤訊息
	local text = "无法连接到服务器，点击确定重新连线."	
	Text_Content.text = text .. "\n" .. "(" .. G.Server.GetErrorCode() .. ")"
	
    Btn_Common:AddClickListener
    (
    function()
		if CS.WindowManager.Instance:FindWindowNodeByName("UICustomer") ~= nil then
			CS.WindowManager.Instance:CloseWindow("UICustomer", false);
		end

		local openparam = CS.WindowNodeInitParam("UICustomer", CS.AppDefine.HallNumber)
		openparam.NodeType = CS.BaseNodeType.Propmt;
		CS.WindowManager.Instance:OpenWindow(openparam)
	end
    )

	Btn_OK:AddClickListener
    (
    function()
        CS.WindowManager.Instance:CloseWindow("UILoginMessageBox", false);
		StartUp.BeginConnect()
    end
    )
	
    Btn_Cancel:AddClickListener
    (
    function()
        CS.WindowManager.Instance:CloseWindow("UILoginMessageBox", false);
		CS.UnityEngine.Application.Quit()
    end
    )
end