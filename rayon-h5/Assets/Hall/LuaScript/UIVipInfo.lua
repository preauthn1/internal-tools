local Btn_RightMove = this.transform:Find("UIBody/_RightMove"):GetComponent("Button")
local Btn_LeftMove = this.transform:Find("UIBody/_LeftMove"):GetComponent("Button")
local Btn_Cancel = this.transform:Find("UIBody/_Cancel"):GetComponent("Button")
local Btn_Recharge = this.transform:Find("UIBody/_Recharge"):GetComponent("Button")

local Img_VipNext = this.transform:Find("UIBody/_VipNext"):GetComponent("Image")
local Img_VipNow = this.transform:Find("UIBody/_VipNow"):GetComponent("Image")
local RecTran_Need = this.transform:Find("UIBody/Need"):GetComponent("RectTransform")
local Txt_NeedValue = this.transform:Find("UIBody/Need/_NeedValue"):GetComponent("Text")
local Txt_ExpValue = this.transform:Find("UIBody/_Exp/_ExpValue"):GetComponent("Text")
local RecTran_Exp = this.transform:Find("UIBody/_Exp/ImageExp/Value"):GetComponent("RectTransform")

local Img_Frame = this.transform:Find("UIBody/_Frame"):GetComponent("Image")

local Img_NowVip = this.transform:Find("UIBody/bgSide/_NowVip"):GetComponent("Image")
local Txt_Inner = this.transform:Find("UIBody/bgSide/_Inner"):GetComponent("Text")
local Txt_AccValue = this.transform:Find("UIBody/bgSide/_AccValue"):GetComponent("Text")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";

local userdata = G.User.GetCurrentUser().UserData

local groupName = "UIVipInfo";

local maxLevel = 0
local nowLevel = 0
local nowValue = 0

-- 初始化函数，实例化之后调用一次
function Start()
    Btn_Recharge:AddClickListener
    (
    function()
		EscapeWindow()
		
        local openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )

    Btn_RightMove:AddClickListener
    (
    function()
		nowLevel = nowLevel + 1
		
		UpdateData()
    end
    )

    Btn_LeftMove:AddClickListener
    (
    function()
		nowLevel = nowLevel - 1
		
		UpdateData()
    end
    )

    Btn_Cancel:AddClickListener
    (
    function()
        EscapeWindow()
    end
    )
end

function UpdateData()
	if nowLevel == maxLevel then
		Btn_RightMove.gameObject:SetActive(false)
	else
		Btn_RightMove.gameObject:SetActive(true)	
	end
	
	if nowLevel <= 1 then
		Btn_LeftMove.gameObject:SetActive(false)
	else
		Btn_LeftMove.gameObject:SetActive(true)	
	end
	
	G.Service.Vip.SetFrame(Img_Frame, nowLevel)
	Img_Frame:ResetSpriteByName(CS.AppDefine.HallNumber, "f" .. tostring(nowLevel))
	Img_VipNow:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(nowLevel))
	
	local nowData = G.Service.Vip.GetVipDataByLevel(nowLevel)
	local nextData = G.Service.Vip.GetVipDataByLevel(nowLevel + 1)

	-- 有下一級資料
	if nextData ~= nil then
		RecTran_Need.gameObject:SetActive(true)
		Img_VipNext.gameObject:SetActive(true)
		Img_NowVip:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(nowLevel + 1))
		Img_VipNext:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(nowLevel + 1))	
	
		Txt_AccValue.text = string.format("%d", math.modf(nextData.nextexp / 100))
		
		local nextexp = (nextData.nextexp - nowData.nextexp) / 100
		
		local inner = nextData.inner
		if nextData.add > 0 then
			inner = inner .. "\n" .. nextData.inner1 .. string.format("%d", math.modf(nextData.add / 100))
		end
		Txt_Inner.text = inner

		local nowLevelValue = nowValue - nowData.nextexp / 100
		if nowLevelValue < 0 then
			nowLevelValue = 0
		end
		
		local width = RecTran_Exp.rect.width
		if nowLevel < maxLevel then
			local value = 0
			-- 玩家比目前等級還高
			if userdata.UserLevel >= nowLevel + 1 then
				RecTran_Need.gameObject:SetActive(false)	
				--Txt_NeedValue.text = "0"
				Txt_ExpValue.text = string.format("%d", math.modf(nextexp)) .. "/" .. string.format("%d", math.modf(nextexp))
				value = 1
			-- 玩家小於 等於 目前等級 上方狀態
			else
				-- 防呆
				if nextexp > 0 then
					Txt_ExpValue.text = string.format("%d", math.modf(nowLevelValue)) .. "/" .. string.format("%d", math.modf(nextexp))
					Txt_NeedValue.text = string.format("%d", math.modf(nextexp - nowLevelValue))
					value = nowLevelValue / nextexp
				else
					Txt_ExpValue.text = string.format("%d", math.modf(nowLevelValue)) .. "/" .. tostring("-")
					Txt_NeedValue.text = tostring("-")
					value = 0			
				end
			end
			
			if value > 1 then
				value = 1
			end
			
			RecTran_Exp.transform.localPosition = CS.UnityEngine.Vector3(-1.5 * width, 0, 0) + CS.UnityEngine.Vector3(value * width , 0, 0)
		else
			Txt_ExpValue.text = tostring(nowLevelValue) .. "/-"		
			Txt_NeedValue.text = "-"		
			RecTran_Exp.transform.localPosition = CS.UnityEngine.Vector3(-width * 0.5, 0, 0)
		end
	-- 無下一級資料
	else
		Txt_NeedValue.text = "-"
		Txt_ExpValue.text = "-/-"
		RecTran_Need.gameObject:SetActive(false)
		Img_VipNext.gameObject:SetActive(false)	
	
		Img_NowVip:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(nowLevel))	
		Txt_AccValue.text = string.format("%d", math.modf(nowData.nextexp / 100))

		local inner = nowData.inner
		if nowData.add > 0 then
			inner = inner .. "\n" .. nowData.inner1 .. string.format("%d", math.modf(nowData.add / 100))
		end
		Txt_Inner.text = inner
	
		RecTran_Exp.transform.localPosition = CS.UnityEngine.Vector3(-1000, 0, 0)	
	end
end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)
end 

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	nowLevel = userdata.UserLevel
	maxLevel = G.Service.Vip.GetVipDataCount()
	nowValue = userdata.SumCharge / 100
	
	UpdateData()
end

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIVipInfo", false);
end