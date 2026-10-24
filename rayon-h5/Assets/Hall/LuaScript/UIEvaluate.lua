local btnClose = this.transform:Find("BtnClose"):GetComponent("Button")
local btnWeb = this.transform:Find("BtnWeb"):GetComponent("Button")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local userdata = G.User.GetCurrentUser().UserData
local statusList = {}
local mailPool = {}
local mailList = G.Service.Mail.GetMailList()
local dailyPool = {}
local dailyList = {}

local groupName = "UIEvaluate";
-- 初始化函数，实例化之后调用一次
function Start()

	btnClose:AddClickListener(function()
		EscapeWindow()
	end)
	
	btnWeb:AddClickListener(function()
		local adl=G.User.GetAdList()
		for i=1,#adl do
		--0:複製,1:直連該網頁,2:一鍵入廳,3:甚麼都不做,4:前往充值,5:评分
			if adl[i].Behavior == 5 then
				CS.UnityEngine.Application.OpenURL(adl[i].Link)
				return
			end
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
	
end


function EscapeWindow()
	
	local play_times = CS.UnityEngine.PlayerPrefs.GetInt("play_times",0)
	CS.UnityEngine.PlayerPrefs.SetInt("evaluate_times",play_times)

	CS.WindowManager.Instance:CloseWindow(groupName, false)
end