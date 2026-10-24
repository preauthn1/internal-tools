local btnClose = this.transform:Find("UIBody/_Close"):GetComponent("Button")
local txtContent = this.transform:Find("UIBody/ScrollView/Content"):GetComponent("Text")

local G = require "Global";
local groupName = "UIRankHelp";

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

-- 初始化函数，实例化之后调用一次
function Start()

    btnClose:AddClickListener(function()
        EscapeWindow()
    end)

end 


function EscapeWindow()	
    CS.WindowManager.Instance:CloseWindow(groupName, false);
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	if windowData ~= nil then
		ShowContent(windowData.receiveMsg)
		SendGetRule(windowData.GameID)
	else
		EscapeWindow()
	end
end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

-- 
function SendGetRule(GameID)
	local parm = {}
	parm.GameID = GameID
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetioniInfoRuleID, parm)
end

-- 比賽活動公告
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GameCompetioniInfoRuleRetID, function(receiveMsg)	
	ShowContent(receiveMsg)
end , groupName);

function ShowContent(receiveMsg)
	if receiveMsg == nil then
		return
	end
	local content = "1.比赛采单人单房制,入场费用为"
	content = content .. SetFormat(G.Tool.UserHelper.ToCash(receiveMsg.Charge))
	content = content .. "金币,玩家可以在一个固定房间来创造最好的比赛积分"
	
	local minute = 0
	local sec = receiveMsg.PlayTime
	while(sec >= 60) do
		sec = sec - 60
		minute = minute +1
	end
	
	content = content .. "\n\n2.在限定的时间（"
	content = content .. SetFormat(minute)
	content = content .. "分"
	if sec > 0 then 
		content = content .. SetFormat(sec) .. "秒"
	else 
		content = content .. "钟"
	end
	content = content .. "）与固定比赛币（"
	content = content .. SetFormat(G.Tool.UserHelper.ToCash(receiveMsg.Coin))
	content = content .. "币）条件下,利用技巧获得高比赛积分"
	
	content = content .. "\n\n3.充值优惠：比赛期间,充值到一定数量的玩家可以每日获得不同次数的免费比赛机会"
	
	local CompareTicket = function (a,b)
		return a.FreeTicket<b.FreeTicket
	end
	
	table.sort(receiveMsg.OrderTicket,CompareTicket)
	for k,v in pairs(receiveMsg.OrderTicket) do
		content = content .. "\n       充值"
		content = content .. SetFormat(G.Tool.UserHelper.ToCash(v.OrderReach))
		content = content .. "以上在比赛期间每日免费"
		content = content .. SetFormat(v.FreeTicket)
		content = content .. "次"
	end
	
	content = content .. "\n\n4.比赛开始到结束时间由官方于比赛前发布"	
	content = content .. "\n\n5.入场后所发射子弹皆为比赛币消耗,与玩家本身金币无关,玩家可以尽情发射,鱼击杀后所获得分数为比赛积分,也与玩家金币无关"
	content = content .. "\n\n6.玩家名次显示在排行榜,最终排名结果以官方公布为主"
	content = content .. "\n\n7.奖励说明："
		
	local CompareBonus = function (a,b)
		return a>b
	end
	
	for k,v in pairs(receiveMsg.RewardInfo) do
		content = content .. "\n\n"
		if v.Type == 1 then
			content = content .. "\n<color=#8219B4FF>- 每日赛 -"
		else
			content = content .. "\n<color=#EE5800FF>- 累积赛 -"
		end
		--table.sort(v.Bonus,CompareBonus)
		for i=1,#v.Bonus do			
			content = content .. "\n       第" .. tostring(i) .. "名"
			content = content .. SetFormat(G.Tool.UserHelper.ToCash(v.Bonus[i]))
			content = content .. "金币"
		end
		content = content .. "</color>"
		
		if v.Type == 1 then
			content = content .. "\n\n每日赛截止时间为每日<color=#FF0053FF>" .. SetFormat(v.DeadLine).. "</color>"
		else
			content = content .. "\n\n累积赛截止时间为<color=#FF0053FF>".. SetFormat(os.date("%m月%d日",StringToTime(receiveMsg.EndTime)))
			content = content .. " " .. SetFormat(v.DeadLine).. "</color>"
		end
	end	

	txtContent.text = content
end

function SetFormat(val)
	--return "<color=#FF0000FF>" ..tostring(val).."</color>"
	return tostring(val)
end

function StringToTime(timeString)
    if type(timeString) ~= 'string' then error('string2time: timeString is not a string') return 0 end
    local fun = string.gmatch( timeString, "%d+")
    local y = fun() or 0
    if y == 0 then error('timeString is a invalid time string') return 0 end
    local m = fun() or 0
    if m == 0 then error('timeString is a invalid time string') return 0 end
    local d = fun() or 0
    if d == 0 then error('timeString is a invalid time string') return 0 end
    local H = fun() or 0
    if H == 0 then error('timeString is a invalid time string') return 0 end
    local M = fun() or 0
    if M == 0 then error('timeString is a invalid time string') return 0 end
    local S = fun() or 0
    if S == 0 then error('timeString is a invalid time string') return 0 end
    return os.time({year=y, month=m, day=d, hour=H,min=M,sec=S})
end