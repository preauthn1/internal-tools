local Btn_Recharge = this.transform:Find("UIBody/_Recharge"):GetComponent("Button")
local Btn_Bingding = this.transform:Find("UIBody/_Bingding"):GetComponent("Button")
local Btn_Close = this.transform:Find("UIBody/_Close"):GetComponent("Button")
local Text_Content = this.transform:Find("UIBody/_Content"):GetComponent("Text")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local groupName = "UILoginGift"
local user = G.User.GetCurrentUser();
local userdata = user.UserData


-- 初始化函数，实例化之后调用一次
function Start()

    Btn_Close:AddClickListener(EscapeWindow)

    Btn_Recharge:AddClickListener
    (
    function()
        EscapeWindow()

        local openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )
	
	Btn_Bingding:AddClickListener
    (
    function()
        EscapeWindow()

        local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )
end 


function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	local temp = "亲爱的玩家:\n      活动期间只要每天登入游戏,\n就可以获得金币的奖励,最多可以领<size=40><color=#FFDF52FF>7天</color></size>。\n\n<color=#FFDF52FF>本活动需绑定帐号"

	--G.Service.LoginBonus.GiftCash = 0
	if G.Service.LoginBonus.GiftCash > 0 then		
		if G.Platform.IsChangeCoin() == true then
			temp = temp .. ",充值" .. tostring(G.Service.LoginBonus.GiftCashToCash) .. "币"
		else
			temp = temp .. ",充值" .. tostring(G.Service.LoginBonus.GiftCashToCash) .. "元"
		end
	end
	temp = temp .. "。</color>"
	
	Text_Content.text = temp
	
	local showRecharge = false
	
	if G.Service.LoginBonus.GiftCash > userdata.SumCharge then
		showRecharge = true
	end
	
	Btn_Bingding.gameObject:SetActive(showRecharge == false)
	Btn_Recharge.gameObject:SetActive(showRecharge == true)
end