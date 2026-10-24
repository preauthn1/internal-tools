local btnClose = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local btnChoose = this.transform:Find("UIBody/ButtonChoose"):GetComponent("Button")
local btnDelete = this.transform:Find("UIBody/ButtonDelete"):GetComponent("Button")
--local imgTitle = this.transform:Find("UIBody/Background/Title/ImageTitle"):GetComponent("Image")

local btnMail = this.transform:Find("UIBody/SubList/Grid/ButtonMail"):GetComponent("Button")
local imgMailSelect = this.transform:Find("UIBody/SubList/Grid/ButtonMail/ImageSelect"):GetComponent("Image")
local imgMailState1 = this.transform:Find("UIBody/SubList/Grid/ButtonMail/ImageState1"):GetComponent("Image")
local imgMailState2 = this.transform:Find("UIBody/SubList/Grid/ButtonMail/ImageState2"):GetComponent("Image")
local imgMailTitle = this.transform:Find("UIBody/SubList/Grid/ButtonMail/ImageTitle"):GetComponent("Image")
local rectMailDetail = this.transform:Find("UIBody/DetailMail"):GetComponent("RectTransform")
local rectMailGrid = this.transform:Find("UIBody/DetailMail/Grid"):GetComponent("RectTransform")
local prefabMail = this.transform:Find("UIBody/DetailMail/Grid/Button").gameObject
local txtMailNote = this.transform:Find("UIBody/DetailMail/TextNote"):GetComponent("Text")

local btnPost = this.transform:Find("UIBody/SubList/Grid/ButtonPost"):GetComponent("Button")
local imgPostSelect = this.transform:Find("UIBody/SubList/Grid/ButtonPost/ImageSelect"):GetComponent("Image")
local imgPostState1 = this.transform:Find("UIBody/SubList/Grid/ButtonPost/ImageState1"):GetComponent("Image")
local imgPostState2 = this.transform:Find("UIBody/SubList/Grid/ButtonPost/ImageState2"):GetComponent("Image")
local imgPostTitle = this.transform:Find("UIBody/SubList/Grid/ButtonPost/ImageTitle"):GetComponent("Image")
local rectPostDetail = this.transform:Find("UIBody/DetailPost"):GetComponent("RectTransform")
local txtPostContent = this.transform:Find("UIBody/DetailPost/Content"):GetComponent("Text")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local groupName = "UIMail";

local userdata = G.User.GetCurrentUser().UserData
local sequenceScale
local pool = {}

-- 初始化函数，实例化之后调用一次
function Start()

    btnClose:AddClickListener(function() EscapeWindow() end)
	btnMail:AddClickListener(function() ChangeSubList(1) end)
	btnPost:AddClickListener(function() ChangeSubList(2) end)
	btnDelete:AddClickListener(function()  end)
	btnChoose:AddClickListener(function() 
		for i=1, #pool do
			pool[i].imgSelect.gameObject:SetActive(true) 
		end
	end)
end
 
-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	windowData = windowData or 1
	ChangeSubList(windowData)
	
	--imgTitle.sprite = imgMailTitle.sprite
end
 

function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

function CreatPrefabMail()
	local clone = CS.UnityEngine.GameObject.Instantiate(prefabMail)
	if(clone == nil)then return nil end
	clone.transform:SetParent(prefabMail.transform.parent)
	clone.transform.localPosition = CS.UnityEngine.Vector3.zero
	clone.transform.localScale = CS.UnityEngine.Vector3.one
	
	local index = #pool +1
	pool[index] = {}
	pool[index].root = clone.transform:GetComponent("RectTransform")
	pool[index].btn = clone.transform:GetComponent("Button")
	pool[index].state1 = clone.transform:Find("ImageState1"):GetComponent("Image")
	pool[index].state2 = clone.transform:Find("ImageState2"):GetComponent("Image")
	pool[index].title = clone.transform:Find("TextTitle"):GetComponent("Text")
	pool[index].snedtime = clone.transform:Find("TextTime"):GetComponent("Text")
	pool[index].line = clone.transform:Find("ImageLine"):GetComponent("RectTransform")
	pool[index].content = clone.transform:Find("TextContent"):GetComponent("Text")
	pool[index].contentrect = clone.transform:Find("TextContent"):GetComponent("RectTransform")
	pool[index].contentfitter = clone.transform:Find("TextContent"):GetComponent("ContentSizeFitter")
	pool[index].btnSelect =  clone.transform:Find("Button"):GetComponent("Button")
	pool[index].imgSelect =  clone.transform:Find("Button/ImageSelect"):GetComponent("Image")
	pool[index].baseSizeDelta = pool[index].root.sizeDelta
end

function ShowContent(index,content)
	local show = pool[index].show
	if show == false then
		pool[index].line.gameObject:SetActive(true)
		pool[index].content.gameObject:SetActive(true)
		pool[index].content.transform.localScale = CS.UnityEngine.Vector3(1,0,1)
		pool[index].content.text = content
		pool[index].show = true
	end
	
	
	local timer = 0.1
	local sequence = CS.DG.Tweening.DOTween.Sequence()
	sequence:AppendInterval(0.1)
	sequence:AppendCallback(function()
		local targetScale = CS.UnityEngine.Vector3(1,0,1)
		-- 此则内容高度
		local targetSizeDelta = pool[index].baseSizeDelta.y
		--要拖曳到的位置
		local targetAnchored = rectMailGrid.anchoredPosition.y
		if show == false then
			targetScale = CS.UnityEngine.Vector3.one
			targetSizeDelta = targetSizeDelta + pool[index].line.sizeDelta.y
			targetSizeDelta = targetSizeDelta + pool[index].contentrect.sizeDelta.y
			targetSizeDelta = targetSizeDelta + 20 --预留空间
			
			local space = pool[index].root.anchoredPosition.y
			space = space + pool[index].baseSizeDelta.y * 0.5
			space = space - targetSizeDelta
			space = space + rectMailDetail.sizeDelta.y
			space = space + rectMailGrid.anchoredPosition.y
			-- 内文过长置顶
			if targetSizeDelta > rectMailDetail.sizeDelta.y then
				targetAnchored = pool[index].root.anchoredPosition.y + rectMailDetail.sizeDelta.y
				targetAnchored = rectMailDetail.sizeDelta.y - targetAnchored
				targetAnchored = targetAnchored - pool[index].baseSizeDelta.y * 0.5
			-- 置底
			elseif space < 0 then
				targetAnchored = targetAnchored + CS.UnityEngine.Mathf.Abs(space)
			end
		end
		targetSizeDelta = CS.UnityEngine.Vector2(pool[index].baseSizeDelta.x,targetSizeDelta)
				
		if sequenceScale ~= nil then
			sequenceScale:Kill()
		end
		sequenceScale = CS.DG.Tweening.DOTween.Sequence()
		sequenceScale:Append(pool[index].content.transform:DOScale(targetScale, timer))
		sequenceScale:Join(pool[index].root:DOSizeDelta(targetSizeDelta, timer))		
		sequenceScale:Append(rectMailGrid:DOAnchorPosY(targetAnchored, 0))
		sequenceScale:AppendCallback(function()
			if show == true then
				HideContent(index)
			end
		end)
	end)
end

function HideContent(index)
	pool[index].line.gameObject:SetActive(false)
	pool[index].content.gameObject:SetActive(false)
	pool[index].root.sizeDelta = pool[index].baseSizeDelta
	pool[index].show = false
end

function ChangeSubList(index)
	imgMailSelect.gameObject:SetActive(index == 1)
	imgPostSelect.gameObject:SetActive(index == 2)
	rectMailDetail.gameObject:SetActive(index == 1)
	rectPostDetail.gameObject:SetActive(index == 2)
	if index == 1 then
		ShowMail()
	else
		ShowPost()
	end
end

function ShowMail()
	--imgTitle.sprite = imgMailTitle.sprite
	mailList = G.Service.Mail.GetMailList()
	
	txtMailNote.gameObject:SetActive(#mailList == 0)
	
	for i=1, #mailList do
		if i > #pool then CreatPrefabMail() end
		
		local index = i
		pool[index].root.gameObject:SetActive(true)
		pool[index].btn.onClick:RemoveAllListeners()
		pool[index].btn:AddClickListener(function() SendMailDetail(index) end)
		pool[index].btnSelect.onClick:RemoveAllListeners()
		pool[index].btnSelect:AddClickListener(function() 
			pool[index].imgSelect.gameObject:SetActive(pool[index].imgSelect.gameObject.activeSelf == false) 
		end)
		pool[index].state1.sprite = imgMailState1.sprite
		pool[index].state2.sprite = imgMailState2.sprite
		pool[index].state1.gameObject:SetActive(mailList[index].IsRead == false)
		pool[index].state2.gameObject:SetActive(mailList[index].IsRead == true)
		pool[index].imgSelect.gameObject:SetActive(false) 
		pool[index].title.text = mailList[index].Title
		pool[index].snedtime.text = mailList[index].SendTime
		HideContent(index)
	end
	
	for i=#mailList +1, #pool do
		pool[i].root.gameObject:SetActive(false)
	end
end

function ShowPost()
	--imgTitle.sprite = imgPostTitle.sprite
	
	txtPostContent.text = userdata.AnnMsg
end

function SendMailDetail(index)
	--再次查询是因为取得信件列表没有信件内容
	local param = { }
	param.ID = mailList[index].ID
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetMailDetailID, param)
end

-- region 查看邮件详情返回
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetMailDetailRetID, function(receiveMsg)
	
	for i=1,#mailList do
		if mailList[i].ID == receiveMsg.ID then
			mailList[i].Title = receiveMsg.Title
			mailList[i].IsRead = receiveMsg.IsRead
			mailList[i].SendTime = receiveMsg.SendTime
			mailList[i].Content = receiveMsg.Content
			
			pool[i].state1.gameObject:SetActive(mailList[i].IsRead == false)
			pool[i].state2.gameObject:SetActive(mailList[i].IsRead == true)
			
			ShowContent(i,receiveMsg.Content)
		end
	end
	
end , groupName);
-- endregion
