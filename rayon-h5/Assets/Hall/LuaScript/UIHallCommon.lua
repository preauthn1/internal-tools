require("Module")
mmodule(...)
local G = require "Global"
local groupName = "UIHall";

local rootRect = nil
local rootGameBg = {}

-----button HeadIcon
local callback_HeadIcon = nil
function SetCallbackHeadIcon(_callback)
	callback_HeadIcon = _callback
end	
-----button Recharge
local callback_Recharge = nil
function SetCallbackRecharge(_callback)
	callback_Recharge = _callback
end	


function Init(root)
	rootRect = root
	
	rootGameBgDefault = root.transform:Find("Background"):GetComponent("RectTransform")
	rootGameBg[1] = root.transform:Find("Background1"):GetComponent("RectTransform")
	rootGameBg[14] = root.transform:Find("Background14"):GetComponent("RectTransform")
	
	--ani_shark = root.transform:Find("Background1/_shark"):GetComponent("SkeletonGraphic");	
	--ani_kani = root.transform:Find("Background1/_kani"):GetComponent("SkeletonGraphic");	
	
	txt_Money = root.transform:Find("_Money"):GetComponent("Text")
	txt_PlayerID = root.transform:Find("_PlayerID"):GetComponent("Text")
	txt_GiftMoney = root.transform:Find("_GiftChange/_GiftMoney"):GetComponent("Text")
	
	img_FrameIcon = root.transform:Find("_HeadIcon/boarder"):GetComponent("Image")
	img_HeadIcon = root.transform:Find("_HeadIcon"):GetComponent("Image")
	G.Service.Vip.SaveHeadData(img_HeadIcon)
	img__VipIcon = root.transform:Find("_HeadIcon/_VipIcon"):GetComponent("Image")
	btn_HeadIcon = root.transform:Find("_HeadIcon"):GetComponent("Button")
	btn_HeadIcon:AddClickListener(function()
		if callback_HeadIcon ~= nil then
			callback_HeadIcon()
		end
	end)
	
	btn_Recharge = root.transform:Find("_Recharge"):GetComponent("Button")
	btn_Recharge:AddClickListener(function()
		if callback_Recharge ~= nil then
			callback_Recharge()
		end
	end)
	
	btn_GiftChange = root.transform:Find("_GiftChange"):GetComponent("Button")
	btn_GiftChange:AddClickListener(function()
        local openparam = CS.WindowNodeInitParam("UIGiftCode", CS.AppDefine.HallNumber)
		CS.WindowManager.Instance:OpenWindow(openparam)
	end)
	if G.User.payinfo.UIClub2AgentMode6 == false then
		btn_GiftChange.gameObject:SetActive(false)
	end
end

--切換至遊戲背景
function ChangeBackground(gameId)
	rootGameBgDefault.gameObject:SetActive(gameId ~= nil and gameId > 0)
	
	for k,v in pairs(rootGameBg) do		
		rootGameBg[k].gameObject:SetActive(k == gameId)		
	end
		
	HideAnimBackground()
	ShowAnimBackground(gameId)
end

--顯示背景動畫
function ShowAnimBackground(gameId)
--[[
	if(gameId==1)then
		timer = 1
		target = ani_kani.transform.position
		sequencekani = CS.DG.Tweening.DOTween.Sequence()
		sequencekani:SetLoops(-1,CS.DG.Tweening.LoopType.Restart)
		sequencekani:AppendCallback(function()
			ani_kani.AnimationState:SetAnimation(0, "walk", true);
			ani_kani.AnimationState:AddAnimation(0, "idle", true , 2);
		end)
		sequencekani:Append(ani_kani.transform:DOMove(CS.UnityEngine.Vector3(307,target.y,target.z), timer))
		sequencekani:Append(ani_kani.transform:DOMove(CS.UnityEngine.Vector3(-207,target.y,target.z), timer)) 
		sequencekani:AppendInterval(5)
		
		ani_shark.AnimationState:SetAnimation(0, "start", false);
		ani_shark.AnimationState:AddAnimation(0, "idle", true , 2);		
    end
	]]--
end
--關閉背景動畫
function HideAnimBackground()
	if(sequencekani ~= nil)then
		sequencekani:Kill()
	end
end

--更新金額
function SetMoney(value)
	txt_Money.text = value
end

--更新禮金額
function SetGiftMoney(value)
	txt_GiftMoney.text = value
end

--顯示充值
function ShowRecharge(active)
	btn_Recharge.gameObject:SetActive(active)
end
--設置頭像圖片
function SetHeadIcon(value)
	img_HeadIcon:ResetSpriteByName(CS.AppDefine.HallNumber, tostring(value))
	
	if value == 0 then
		img__VipIcon.gameObject:SetActive(false)
	else
		img__VipIcon.gameObject:SetActive(true)
	end	
end

--設置頭像外框圖片
function SetFrameIcon(value)
	G.Service.Vip.SetFrame(img_FrameIcon, value)
end

--設置VIP等级
function SetVipLevel(value)	
	img__VipIcon:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(value))
end

function SetPlayerID(value)
	txt_PlayerID.text = value
end