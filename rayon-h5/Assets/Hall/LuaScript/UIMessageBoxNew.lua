local G = require "Global"

local btn_Close = this.transform:Find("UIBody/_Close"):GetComponent("Button")
local btn_OK = this.transform:Find("UIBody/Buttons/_OK"):GetComponent("Button")
local txt_Content = this.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")

local groupName = "UIMessageBoxNew";

function Start()
    btn_Close:AddClickListener(function() EscapeWindow() end)
    btn_OK:AddClickListener(function() EscapeWindow() end)
end

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

