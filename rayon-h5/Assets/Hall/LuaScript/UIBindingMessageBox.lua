local Btn_OK = this.transform:Find("UIBody/Buttons/_OK"):GetComponent("Button")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"

-- 初始化函数，实例化之后调用一次
function Start()
	Btn_OK:AddClickListener
    (
    function()
		CS.UnityEngine.Application.Quit()
    end
    )
end