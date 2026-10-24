require("Module")
mmodule(...)
local G = require "Global"
local groupName = "UIHall";

local root = nil
local pool = {}
local ranklist = {}

local refreshGameID = nil
local refreshFlagID = nil
local receiveMsgCharge = 0


-----button back
local callbackBack = nil
function SetCallbackBack(_callback)
	callbackBack = _callback
end	

-----button start
local callbackStart = nil
function SetCallbackStart(_callback)
	callbackStart = _callback
end	

-----button rechange
local callbackRecharge = nil
function SetCallbackRecharge(_callback)
	callbackRecharge = _callback
end

-----button rank help
local callbackRankHelp = nil
function SetCallbackRankHelp(_callback)
	callbackRankHelp = _callback
end	

function Init(_root)
	root = _root
    perfab = root.transform:Find("ScrollView/Content/Item"):GetComponent("RectTransform").gameObject	
	btnBack = root.transform:Find("_ButtonBack"):GetComponent("Button")
	btnRankHelp = root.transform:Find("_ButtonHelp"):GetComponent("Button")
	btnStart = root.transform:Find("_ButtonStart"):GetComponent("Button")
	btnContinue = root.transform:Find("_ButtonContinue"):GetComponent("Button")
	btnRankList = {}
	btnRankList[1] = {}
	btnRankList[1].btn = root.transform:Find("AttachRank/_ButtonToday"):GetComponent("Button")
	btnRankList[1].btn.enabled = false
	btnRankList[1].normal = root.transform:Find("AttachRank/_ButtonToday/ImageNormal"):GetComponent("Image")
	btnRankList[1].normal.gameObject:SetActive(false)
	btnRankList[1].choose = root.transform:Find("AttachRank/_ButtonToday/ImageChoose"):GetComponent("Image")
	btnRankList[1].choose.gameObject:SetActive(false)
	
	btnRankList[2] = {}
	btnRankList[2].btn = root.transform:Find("AttachRank/_ButtonYesterday"):GetComponent("Button")
	btnRankList[2].btn.enabled = false
	btnRankList[2].normal = root.transform:Find("AttachRank/_ButtonYesterday/ImageNormal"):GetComponent("Image")
	btnRankList[2].normal.gameObject:SetActive(false)
	btnRankList[2].choose = root.transform:Find("AttachRank/_ButtonYesterday/ImageChoose"):GetComponent("Image")
	btnRankList[2].choose.gameObject:SetActive(false)
		
	btnRankList[3] = {}
	btnRankList[3].btn = root.transform:Find("AttachRank/_ButtonWeek"):GetComponent("Button")
	btnRankList[3].btn.enabled = false
	btnRankList[3].normal = root.transform:Find("AttachRank/_ButtonWeek/ImageNormal"):GetComponent("Image")
	btnRankList[3].normal.gameObject:SetActive(false)
	btnRankList[3].choose = root.transform:Find("AttachRank/_ButtonWeek/ImageChoose"):GetComponent("Image")
	btnRankList[3].choose.gameObject:SetActive(false)
	
	btnRecharge = root.transform:Find("_Recharge"):GetComponent("Button")
	txtVersion = root.transform:Find("_Version"):GetComponent("Text")
	txtMoney = root.transform:Find("_Money/Text"):GetComponent("Text")
	txtScore = root.transform:Find("_Score/Text"):GetComponent("Text")
	txtRank = root.transform:Find("_Rank/Text"):GetComponent("Text")
	txtFree = root.transform:Find("_Free/Text"):GetComponent("Text")
	txtNote = {}
	txtNote[1] = root.transform:Find("_Note/Text1"):GetComponent("Text")
	txtNote[2] = root.transform:Find("_Note/Text2"):GetComponent("Text")
	txtNote[3] = root.transform:Find("_Again/Text"):GetComponent("Text")
	rootAngin = root.transform:Find("_Again")
		
	-----button back
	btnBack:AddClickListener(function()
		if callbackBack ~= nil then
			callbackBack()
		end
	end)
		
	-----button start	
	btnStart:AddClickListener(function()
		if G.Platform.IsChangeCoin() == true then
			local content = "若没有免费次数，需入场费"..G.Tool.UserHelper.ToCash(receiveMsgCharge).."金币，确定要进入游戏吗?"
		else
			local content = "若没有免费次数，需入场费"..G.Tool.UserHelper.ToCash(receiveMsgCharge).."元金币，确定要进入游戏吗?"
		end
		CS.UIMessageBox.Show(CS.MessageBoxStyle.OKCancel, "", content, "确定", "取消", function(rst)
			if rst == CS.MessageBoxResult.OK then		
				CS.UILoadingData.Show("等待进入游戏中....,请稍后...",5)	
				local parm = {}
				parm.GameID = refreshGameID
				parm.FlagID = tonumber(refreshFlagID)
				parm.UserID = G.User.GetCurrentUserID()
				G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.BuyTicketID, parm)
            end
        end)

	end)
	-----button continue	
	btnContinue:AddClickListener(function()
		if callbackStart ~= nil then
			callbackStart()
		end

	end)
	
	-----button rechange
	btnRecharge:AddClickListener(function()
		if callbackRecharge ~= nil then
			callbackRecharge()
		end
	end)	
	-----button rank help
	btnRankHelp:AddClickListener(function()
		if callbackRankHelp ~= nil then
			callbackRankHelp()
		end
	end)	
	
	-----button today
	btnRankList[1].btn:AddClickListener(function()
		ShowRankList(1)
	end)
	
	-----button yesterday
	btnRankList[2].btn:AddClickListener(function()
		ShowRankList(2)
	end)
	
	-----button week
	btnRankList[3].btn:AddClickListener(function()
		ShowRankList(3)
	end)	
end

local function CreateClone(index)
		local clone = CS.UnityEngine.GameObject.Instantiate(perfab)
		if(clone == nil)then return nil end
		clone.transform:SetParent(perfab.transform.parent)
		clone.transform.localPosition = CS.UnityEngine.Vector3.zero
		clone.transform.localScale = CS.UnityEngine.Vector3.one
			
		pool[index] = {}
		pool[index].root = clone.transform:GetComponent("RectTransform")
		pool[index].imgBg = {}
		pool[index].imgBg[1] = clone.transform:Find("ImageBg1"):GetComponent("Image")
		pool[index].imgBg[2] = clone.transform:Find("ImageBg2"):GetComponent("Image")
		pool[index].imgNo = {}
		pool[index].imgNo[1] = clone.transform:Find("ImageNo1"):GetComponent("Image")
		pool[index].imgNo[2] = clone.transform:Find("ImageNo2"):GetComponent("Image")
		pool[index].imgNo[3] = clone.transform:Find("ImageNo3"):GetComponent("Image")
		pool[index].txtNo = clone.transform:Find("TextNo"):GetComponent("Text")
		pool[index].txtID = clone.transform:Find("TextID"):GetComponent("Text")
		pool[index].txtRank = clone.transform:Find("TextRank"):GetComponent("Text")
		pool[index].imgMy = clone.transform:Find("ImageMy"):GetComponent("Image")
		
	end

function ShowRankList(index)
	for i=1,3 do
		btnRankList[i].choose.gameObject:SetActive(i == index)
	end

	local list = ranklist[index]	
	for i=1,#list do 
		if i > #pool then
			CreateClone(i) 
		end
		
		pool[i].root.gameObject:SetActive(true)
		pool[i].imgNo[1].gameObject:SetActive(list[i].SortID == 1)
		pool[i].imgNo[2].gameObject:SetActive(list[i].SortID == 2)
		pool[i].imgNo[3].gameObject:SetActive(list[i].SortID == 3)
		pool[i].txtNo.gameObject:SetActive(list[i].SortID > 3)
		if list[i].SortID <= 100 then
			pool[i].txtNo.text = tostring(list[i].SortID)
		elseif list[i].SortID > 100 then
			pool[i].txtNo.text = "100+"
		end
		pool[i].txtID.text = tostring(list[i].UserID)
		pool[i].txtRank.text = G.Tool.UserHelper.ToCash(list[i].Score)
		
		pool[i].imgMy.gameObject:SetActive(list[i].UserID == G.User.GetCurrentUserID())
	end	
	
	for i =#list +1 ,#pool do
		pool[i].root.gameObject:SetActive(false)
	end
end

-----text version
function SetVersion(_ver)
	txtVersion.text = _ver
end
		
-----text money
function SetMoney(_ver)
	txtMoney.text = _ver
end


function Show(gameID,flagID)
	refreshGameID = gameID
	refreshFlagID = flagID
	local parm = {}
	parm.GameID = gameID
	parm.FlagID = tonumber(flagID)
	parm.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetitionID, parm)
	
	RefreshFixedupdate()
end

function RefreshShow()
	if(refreshGameID == nil) then
		return
	end
	if(refreshFlagID == nil) then
		return
	end
	Show(refreshGameID,refreshFlagID)
end

--固定時間刷新
function RefreshFixedupdate()
	if sequenceRefresh ~= nil then
		return
	end
	sequenceRefresh = CS.DG.Tweening.DOTween.Sequence()
	sequenceRefresh:SetLoops(-1,CS.DG.Tweening.LoopType.Restart)	
	sequenceRefresh:AppendInterval(60)			
	sequenceRefresh:AppendCallback( function()
		RefreshShow()
	end)

end

function Hide()
    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
	
	root.gameObject:SetActive(false)

	refreshGameID = nil
	refreshFlagID = nil	
	if sequenceRefresh ~= nil then
		sequenceRefresh:Kill()
		sequenceRefresh = nil
	end

end

function GetIsShow()
	return root.gameObject.activeSelf
end

-- 比賽活動公告
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetitionRetID, function(receiveMsg)	
	root.gameObject:SetActive(true)
	txtScore.text = G.Tool.UserHelper.ToCash(receiveMsg.BestScore)
	txtRank.text = receiveMsg.UserRanking

	receiveMsgCharge = receiveMsg.Charge
	btnStart.gameObject:SetActive(receiveMsg.PiSaiTime <= 0)
	btnContinue.gameObject:SetActive(receiveMsg.PiSaiTime > 0)	
	
	rootAngin.gameObject:SetActive(receiveMsg.FreeTimes == 0)
	
	txtFree.text = receiveMsg.FreeTimes
	
	for i=1 ,#txtNote do
		txtNote[i].text = ""
	end
	for i=1 ,#receiveMsg.MsgInfo do
		txtNote[i].text = receiveMsg.MsgInfo[i]		
	end
	
	ranklist[1] = {}
	ranklist[2] = {}
	ranklist[3] = {}
	local showRankListIndex = false
	for k,v in pairs(receiveMsg.Leaderboard) do
		local index = nil
		if v.TitleName == "today" then
			index = 1
		elseif v.TitleName == "yesterday" then
			index = 2
		elseif v.TitleName == "week" then
			index = 3
		end
		
		if index ~= nil then
		ranklist[index] = v.LeaderboardList
		btnRankList[index].btn.enabled = true
		btnRankList[index].normal.gameObject:SetActive(true)
			if showRankListIndex == false then
				showRankListIndex = true
				ShowRankList(index)
			end
		end
	end
end , groupName);


-- region 回傳買票入場
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.BuyTicketRetID, function(receiveMsg)
	CS.UILoadingData.Hide()
	if receiveMsg.RetCode == 1 then
		if callbackStart ~= nil then
			callbackStart()
		end
	else
		CS.UIBubblePrompt.Show(receiveMsg.RetMsg, groupName)
	end
end , groupName);