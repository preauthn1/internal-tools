local BgAttach = this.transform:Find("BG"):GetComponent("RectTransform")
local BgPic = this.transform:Find("BG/Pic1")

local Scrollbar = this.transform:Find("Scrollbar"):GetComponent("Scrollbar")
local PicMask = this.transform:Find("Scrollbar/Sliding Area/Handle/PicMask"):GetComponent("RectTransform")

local ScrollbarTemp = this.transform:Find("Scrollbar_Temp"):GetComponent("Scrollbar")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"
local groupName = "UIVerification"

local EventSystem = CS.UnityEngine.EventSystems.EventSystem
local offset = 0.03
local errorTimes = 3
local accErrorTimes = 0
local pos = 0

-- 初始化函数，实例化之后调用一次
function Start()
	Init()
end

function Init()
	pos = math.random(2,8) * 0.1
	Scrollbar.value = pos
	ScrollbarTemp.value = pos
	print("Imit",pos)

	BgPic:SetParent(BgAttach)
	BgPic.transform.position = CS.UnityEngine.Vector3(0,0,0)	
	BgPic:SetParent(PicMask)
	BgPic.transform.localPosition = CS.UnityEngine.Vector3(BgPic.localPosition.x,BgPic.localPosition.y,0)	

	Scrollbar.value = 0
end

-- 是否点中了UI
local isPointDownOnUI = 0

function Update()
    if G.Env.IsWindows() then
        if EventSystem.current:IsPointerOverGameObject() and CS.UnityEngine.Input:GetMouseButtonDown(0) then
            isPointDownOnUI = 1
            return
        end

        if isPointDownOnUI and CS.UnityEngine.Input.GetMouseButtonUp(0) then
            isPointDownOnUI = 2
        end
    else
        if CS.UnityEngine.Input.touchCount > 0 and EventSystem.current:IsPointerOverGameObject(CS.UnityEngine.Input.GetTouch(0).fingerId) and CS.UnityEngine.Input.GetTouch(0).phase == CS.UnityEngine.TouchPhase.Began then
            isPointDownOnUI = 1
            return
        end

        if isPointDownOnUI and CS.UnityEngine.Input.GetTouch(0).phase == CS.UnityEngine.TouchPhase.Ended then
            isPointDownOnUI = 2
        end
    end

	if isPointDownOnUI ~= 2 then
		return
	end
	
	isPointDownOnUI = 0
	
    if Scrollbar.value == 0 then
    	return
    end

	print("Scrollbar",Scrollbar.value)

	if Scrollbar.value > pos - offset and Scrollbar.value < pos + offset then
		CS.UIBubblePrompt.Show("输入正确", groupName);
	    CS.WindowManager.Instance:CloseWindow(groupName, false)
	else
		accErrorTimes = accErrorTimes + 1
		Scrollbar.value = 0
		
		if accErrorTimes < 3 then
			CS.UIBubblePrompt.Show("输入错误", groupName);
		else
			CS.UIBubblePrompt.Show("输入错误次数已达上限", groupName);
			CS.WindowManager.Instance:CloseWindow(groupName, false)
		end
	end
end