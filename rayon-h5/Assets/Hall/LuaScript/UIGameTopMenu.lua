local RectTran_root = this.transform:GetComponent("RectTransform")
local Btn_exit = this.transform:Find("menu/_menuGroup/_exit"):GetComponent("Button")
local Btn_help = this.transform:Find("menu/_menuGroup/_help"):GetComponent("Button")
local Btn_set = this.transform:Find("menu/_menuGroup/_set"):GetComponent("Button")
local Btn_recharge = this.transform:Find("_recharge"):GetComponent("Button")
local RecTran_ann = this.transform:Find("_ann"):GetComponent("RectTransform")
local Txt_anncontent = this.transform:Find("_ann/shade/_anncontent"):GetComponent("Text")
local Txt_time = this.transform:Find("_time"):GetComponent("Text")
local Txt_roomname = this.transform:Find("_roomname"):GetComponent("Text")
local Btn_menuShow = this.transform:Find("menu/_show"):GetComponent("Button")
local Btn_menuHide = this.transform:Find("menu/_menuGroup/_hide"):GetComponent("Button")
local RecTran_menuGroup = this.transform:Find("menu/_menuGroup"):GetComponent("RectTransform")
local Btn__menuGroupBG = this.transform:Find("menu/_menuGroup/_menuGroupBG"):GetComponent("Button")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global"

local event_exit_callback = function()
--    if (G.IsGameing()) then
--        if (isCustomExitUI == true) then
--            local openparam = CS.WindowNodeInitParam("UIExitGame", G.GetGamePName())
--            openparam.NodeType = 5
--            CS.WindowManager.Instance:OpenWindow(openparam)
--        else
--            local openparam = CS.WindowNodeInitParam("UIExitGame", CS.AppDefine.HallNumber)
--            openparam.NodeType = 5
--            CS.WindowManager.Instance:OpenWindow(openparam)
--        end
--    else
        G.Service.Game.ExitGame()
--    end
end

local isCustomExitUI = false
local event_exit = function()
	local callback = G.Service.Game.GetCallbackExit()
    if (callback ~= nil) then
		callback(event_exit_callback)
	else
		event_exit_callback()
    end
end

local event_help = function()
    local openparam = CS.WindowNodeInitParam("UIHelp", G.GetGamePName())
    openparam.NodeType = 5
    CS.WindowManager.Instance:OpenWindow(openparam)
end

local isCustomSetUI = false
local event_set = function()
	-- 改成和Exit一樣，可以外部指定使用哪一個介面 By Rayon
	if (isCustomSetUI == true) then
		local openparam = CS.WindowNodeInitParam("UISetGame", G.GetGamePName())
		openparam.NodeType = 5
		CS.WindowManager.Instance:OpenWindow(openparam)
	else
		local openparam = CS.WindowNodeInitParam("UISetGame", CS.AppDefine.HallNumber)
		openparam.NodeType = 5
		CS.WindowManager.Instance:OpenWindow(openparam)
	end
end

local event_recharge = function()
	local openparam =nil
	if G.GameConfig.IsSubmissionCtx() == true then
		openparam = CS.WindowNodeInitParam("UIRecharge_IOS", CS.AppDefine.HallNumber)     
	else
		openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)      
	end
    CS.WindowManager.Instance:OpenWindow(openparam)
end

local event_showMenu = function()
	RecTran_menuGroup.gameObject:SetActive(true)
end

local event_hideMenu = function()
	RecTran_menuGroup.gameObject:SetActive(false)
end

local annShadeWidth = Txt_anncontent.transform.parent:GetComponent("RectTransform").sizeDelta.x -- 公告文本遮罩宽度 
local annInitPos = Txt_anncontent.transform.localPosition
local RectTran_anncontent = Txt_anncontent.transform:GetComponent("RectTransform")

function GetNoticeContentWidth()
    return CS.UnityEngine.UI.LayoutUtility.GetPreferredWidth(RectTran_anncontent)
end 
 
-- 初始化函数，实例化之后调用一次
function Start()
    Btn_exit:AddClickListener(event_exit)
    Btn_help:AddClickListener(event_help)
    Btn_set:AddClickListener(event_set)
	Btn_recharge:AddClickListener(event_recharge)
	Btn_recharge.gameObject:SetActive(G.GameConfig.IsRechargeCtx())
	
    Btn_menuShow:AddClickListener(event_showMenu)
    Btn_menuHide:AddClickListener(event_hideMenu)
	Btn__menuGroupBG:AddClickListener(event_hideMenu)
	event_hideMenu()
	
	local GameId = G.GetRunGameID()

	if G.Game.IsMultiGame(GameId) then
		Txt_roomname.text = G.Room.GetMultiRoomName()
	elseif GameId == 1 then
		Txt_roomname.text = G.Room.GetDDZRoomName()
	else
	
		Txt_roomname.text = G.Room.GetRoomLevelName()
		if (G.Game.GetIsShowTableNo(GameId)) then
			G.Service.Game.RegTableNoChange( function(tableNo)
				tableNo = tableNo + 1
				Txt_roomname.text = string.format("%s %s", G.Room.GetRoomLevelName(), string.rep("0", 3 - #(tostring(tableNo))) .. tableNo)
			end )
		end
	end
	
	RecTran_ann.gameObject:SetActive(false)	
	
	local canvas = CS.UnityEngine.GameObject.Find("Canvas")
	if(canvas ~= nil) then
		local _x = canvas:GetComponent("RectTransform").sizeDelta.x
		RectTran_root.sizeDelta = CS.UnityEngine.Vector2(_x,0)
	end
end

function Update()
    if CS.UnityEngine.Input.GetKeyDown("escape") then
        if (CS.WindowManager.Instance:FindWindowNodeByName("UIExitGame") == nil) then
            event_exit()
        end
    end

    Txt_time.text = os.date("%H:%M")
	
	-- 送審階段關掉跑馬燈 By Rayon
	if G.GameConfig.IsSubmissionCtx() == true then
		return
	end
	
    if (#Txt_anncontent.text > 0 and Txt_anncontent.transform.localPosition.x > annInitPos.x - GetNoticeContentWidth() - annShadeWidth) then
        Txt_anncontent.transform.localPosition = Txt_anncontent.transform.localPosition - CS.UnityEngine.Vector3(CS.UnityEngine.Time.deltaTime * 50, Txt_anncontent.transform.localPosition.y, Txt_anncontent.transform.localPosition.z)
    else
        local message = G.Service.Ann.GetOneAnnMeg()

        if message == nil then
            RecTran_ann.gameObject:SetActive(false)
            Txt_anncontent.text = ""
        else
            message = string.gsub(message,"%s"," ")
            Txt_anncontent.transform.localPosition = annInitPos
            Txt_anncontent.text = message
            RecTran_ann.gameObject:SetActive(true)
        end
    end
end

function OnDestroy()
    G.Service.Game.RemoveTableNoChange()
    CS.WindowManager.Instance:CloseWindow("UIExitGame", false)
end

-- 接收公告通知
function ReceiveAnnNotify(data)

    Txt_anncontent.text = Txt_anncontent.text .. data
    annWidth = Txt_anncontent.transform:GetComponent("RectTransform").sizeDelta.x
    RecTran_ann.gameObject:SetActive(true)

end


 
-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
    if (windowData ~= nil) then
        isCustomExitUI = windowData.isCustomExitUI
		
		-- 改成和Exit一樣，可以外部指定使用哪一個介面 By Rayon
        isCustomSetUI = windowData.isCustomSetUI
		
		if windowData.annPos ~= nil then
			RecTran_ann.transform.position = windowData.annPos
		end
    end
end
 