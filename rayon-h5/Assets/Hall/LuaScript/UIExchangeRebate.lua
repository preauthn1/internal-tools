local Btn_Close = this.transform:Find("UIBody/_Close"):GetComponent("Button")
local txt__Bonus = this.transform:Find("UIBody/_Bonus"):GetComponent("Text")
local img__Reach = this.transform:Find("UIBody/_Reach"):GetComponent("Image")
local txt__NowText = this.transform:Find("UIBody/_NowText"):GetComponent("Text")
local txt__LimitText = this.transform:Find("UIBody/_LimitText"):GetComponent("Text")
local txt__FlowRate = this.transform:Find("UIBody/_FlowRate"):GetComponent("Text")


-- local PoolStar = this.transform:Find("UIBody/PoolStar"):GetComponent("RectTransform")

local G = require "Global";
local PoolSequence = {}

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

-- 初始化函数，实例化之后调用一次
function Start()
    
    Btn_Close:AddClickListener
    (
    function()
        EscapeWindow()
    end
    )

   
	
end 


function EscapeWindow()
	-- for i=0,#PoolSequence do
	-- 	PoolSequence[i]:Kill()
	-- end	

    CS.WindowManager.Instance:CloseWindow("UIExchangeRebate", false);
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	user = G.User.GetCurrentUser();
	
	-- G.Tool.TableHelper.Print(user.FlowInfo )
	txt__Bonus.text=math.modf(user.FlowInfo.BonusRate*0.01);
	txt__NowText.text=tostring(user.FlowInfo.FlowAmountNow*0.01);
	txt__LimitText.text=tostring(user.FlowInfo.FlowAmountLimit*0.01);
	img__Reach.gameObject:SetActive(user.FlowInfo.Achievement == 1)
	txt__FlowRate.text=tostring(user.FlowInfo.FlowAmountRate);

	--FlowAmountRate
	-- for i=0,PoolStar.transform.childCount-1 do
	-- 	local object = PoolStar.transform:GetChild(i)
	-- 	object.transform.localScale = CS.UnityEngine.Vector3.zero
	-- end
end