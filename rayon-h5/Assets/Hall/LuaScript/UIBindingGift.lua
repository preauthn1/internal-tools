local Btn_BingdingGift = this.transform:Find("UIBody/_BingdingGift"):GetComponent("Button")
local Btn_CloseGift = this.transform:Find("UIBody/_CloseGift"):GetComponent("Button")
-- local PoolStar = this.transform:Find("UIBody/PoolStar"):GetComponent("RectTransform")

local Img_Girl = this.transform:Find("UIBody/girl"):GetComponent("Image")
local Img_Girl1 = this.transform:Find("UIBody/girl1"):GetComponent("Image")
local Txt_Bonus = this.transform:Find("UIBody/BonusText"):GetComponent("Text")

local G = require "Global";
local user = G.User.GetCurrentUser();
local userdata = user.UserData
local PoolSequence = {}

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

-- 初始化函数，实例化之后调用一次
function Start()

	if G.Platform.getJPYL() == false then
		Img_Girl.gameObject:SetActive(true)
		Img_Girl1.gameObject:SetActive(false)
	else
		Img_Girl.gameObject:SetActive(false)
		Img_Girl1.gameObject:SetActive(true)
	end

    Btn_CloseGift:AddClickListener
    (
    function()
        EscapeWindow()
    end
    )

    Btn_BingdingGift:AddClickListener
    (
    function()
        local openparam = CS.WindowNodeInitParam("UIUserInfo", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
		openparam.WindowCloseCallBack = (function()
			EscapeWindow()
		end)
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )
	-- 綁定送金幣
	if(userdata.GiftMoney == 0)then
		return;
	end
	local value = userdata.GiftMoney % 100
	local money = userdata.GiftMoney / 100
	if(value > 0)then
		Txt_Bonus.text = tostring(money)
	else
		Txt_Bonus.text = math.floor(money)
	end
	
	
	-- for i=0,PoolStar.transform.childCount-1 do
	-- 	local timer = 0.3
	-- 	local object = PoolStar.transform:GetChild(i)
	-- 	local image = object.transform:GetComponent("Image")
	-- 	PoolSequence[i] = CS.DG.Tweening.DOTween.Sequence()
	-- 	PoolSequence[i]:AppendCallback(function() 
	-- 		object.transform.localScale = CS.UnityEngine.Vector3.zero
	-- 		image.color = CS.UnityEngine.Color(image.color.r,image.color.g,image.color.b,0)
	-- 	end)
	-- 	PoolSequence[i]:AppendInterval(CS.UnityEngine.Random.Range(0, 5))
	-- 	PoolSequence[i]:SetLoops(-1,CS.DG.Tweening.LoopType.Restart)
	-- 	PoolSequence[i]:Append(object.transform:DOScale(CS.UnityEngine.Vector3.one,timer))
	-- 	PoolSequence[i]:Join(object.transform:DOLocalRotate(CS.UnityEngine.Vector3(0,0,30),timer))
	-- 	PoolSequence[i]:Join(image:DOColor(CS.UnityEngine.Color(image.color.r,image.color.g,image.color.b,1), timer))
	-- 	PoolSequence[i]:AppendInterval(0.1)
	-- 	PoolSequence[i]:Append(object.transform:DOScale(CS.UnityEngine.Vector3.zero,0.05))
	-- end
end 


function EscapeWindow()
	-- for i=0,#PoolSequence do
	-- 	PoolSequence[i]:Kill()
	-- end
	
    CS.WindowManager.Instance:CloseWindow("UIBindingGift", false);
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	-- for i=0,PoolStar.transform.childCount-1 do
	-- 	local object = PoolStar.transform:GetChild(i)
	-- 	object.transform.localScale = CS.UnityEngine.Vector3.zero
	-- end
end