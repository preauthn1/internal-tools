local Btn_Close = this.transform:Find("UIBody/Top/_Close"):GetComponent("Button")
local Btn_Recharge1 = this.transform:Find("UIBody/_Recharge1"):GetComponent("Button")
local Btn_Recharge2 = this.transform:Find("UIBody/_Recharge2"):GetComponent("Button")
local Btn_Recharge3 = this.transform:Find("UIBody/_Recharge3"):GetComponent("Button")
local Btn_Recharge4 = this.transform:Find("UIBody/_Recharge4"):GetComponent("Button")
local Btn_Recharge5 = this.transform:Find("UIBody/_Recharge5"):GetComponent("Button")
local Btn_Recharge6 = this.transform:Find("UIBody/_Recharge6"):GetComponent("Button")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

-- 初始化函数，实例化之后调用一次
function Start()
    Btn_Close:AddClickListener
    (
    function()
        CS.WindowManager.Instance:CloseWindow("UIRecharge_IOS", false);
    end
    )

    Btn_Recharge1:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )

    Btn_Recharge2:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )

    Btn_Recharge3:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )

    Btn_Recharge4:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )

    Btn_Recharge5:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )

    Btn_Recharge6:AddClickListener
    (
    function()
        CS.UIBubblePrompt.Show("功能还没实现", groupName);
    end
    )
end

-- 每帧调用一次
function Update()

end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()

end  

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)

end

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIRecharge_IOS", false);
end