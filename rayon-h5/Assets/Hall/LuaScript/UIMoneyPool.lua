local Txt_Money = this.transform:Find("UIBody/_Money"):GetComponent("Text")
local Btn_CloseGift = this.transform:Find("UIBody/_CloseGift"):GetComponent("Button")
local btn__CloseRule = this.transform:Find("UIBody/Rule/_CloseRule"):GetComponent("Button")
local img_Rule = this.transform:Find("UIBody/Rule"):GetComponent("RawImage")
local btn__OpenRule = this.transform:Find("UIBody/_OpenRule"):GetComponent("Button")



--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---
local G = require "Global";
local PoolSequence = {}
local sequence=nil

-- 初始化函数，实例化之后调用一次
function Start()

  
    Btn_CloseGift:AddClickListener
    (
    function()
    	EscapeWindow()
    end
    )
    btn__CloseRule:AddClickListener
    (
    function()
       img_Rule.gameObject:SetActive(false);
    end
    )
    btn__OpenRule:AddClickListener
    (
    function()
       img_Rule.gameObject:SetActive(true);
    end
    )


end

function DestoryDotween()
        if (sequence ~= nil) then sequence:Kill() end
end

-- 每帧调用一次
function Update()

end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    DestoryDotween()

end  

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)

	

	moneypool=windowData
     sequence= CS.DG.Tweening.DOTween.Sequence()
     sequence:SetLoops(-1)
   
     sequence:AppendCallback( function()
        moneypool= moneypool+10
        Txt_Money.text=tostring(G.Tool.UserHelper.ToCashDetail(moneypool))
     end)

end

-- 在Window下按下Esc或在Android 下按下Back
function EscapeWindow()
	 CS.WindowManager.Instance:CloseWindow("UIMoneyPool", false);
end
