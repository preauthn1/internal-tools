require("Module")
mmodule(...)

local G = require "Global"
local groupName = "AnnSerice";

local AnnList = { };            -- 公告队列

function IsShow(message)
	
	local user = G.User.GetCurrentUser()
	local userdata = user.UserData
	
	local isShowAll = G.Game.IsShowAllGame()
	
	-- 只顯示的遊戲
	local onlyShowId = G.Platform.GetOnlyShowGameId()
	local gameNames = G.Game.GetFilterGameNamess()
    for k, v in pairs(gameNames) do		
		-- 找關鍵字
		q = string.match(message, "%>" .. v .. "%(")
		if q ~= nil then
			-- 是否顯示遊戲
			if isShowAll == true then
				if G.IsShowGame(tonumber(k)) == false then
					return nil
				end
			else
				if onlyShowId > 0 then
					-- 顯示只顯示的遊戲
					if tonumber(k) ~= onlyShowId then
						return nil
					end
				end
			end

			local iconName = G.Game.GetGameIconName(k)
			if iconName > 0 then
				local name = G.Game.GetGameNameByIconName(tostring(k), iconName)
				local newMessage = string.gsub(message, v, name)
				return newMessage
			end
		end	
	end

	return message
end

function GetOneAnnMeg()

    for k, v in pairs(AnnList) do
		--if v ~= nil then
		if v.ID > 0 then
			--AnnList[v] = nil
			
			local msg = IsShow(v.MsgContent)
			v.ID = 0
			if msg ~= nil then
				return msg
			end
		end
	end
	
	return nil
	--[[
    local message = G.Tool.TableHelper.Pull(AnnList)

    --    if message == nil then
    --        return "123123123123123sdsdf";
    --    else
    --        return message
    --    end

    return message
	]]--
end

function AddOneAnn(ann)
	-- 全平台
	if ann.PlatformID == 0 then
		table.insert(AnnList, ann)
	-- 根据平台做筛选
	elseif ann.PlatformID == G.GameConfig.GamePlatformID then
		table.insert(AnnList, ann)
	end
end

-- 是否过滤系统讯息 并优先显示
local isFilterSystemAnn = true
-- 小馀这个值表示为系统讯息跑马灯
local systemIndex = 100000
		
G.NetQ.RegMsgHook(G.Proto.BClass.NotifyServer, G.Proto.NotifyServer.AnnID, function(receiveMsg)
	local user = G.User.GetCurrentUser()
	if user == nil then
		return
	end

	if G.Game.IsShowAllGame() == false then
		return
	end
	
	local TempAnnList = AnnList
	
	-- 清空跑马灯资料,避免过多
	AnnList = {}
	
	-- 旧方法 不区分系统讯息
	if isFilterSystemAnn == false then
		for i, ann in ipairs(receiveMsg.Ann) do
			AddOneAnn(ann)
		end
	-- 过滤系统讯息 并优先显示
	else		
		-- 不清掉系统讯息跑马灯
		for i, ann in ipairs(TempAnnList) do
			if ann.ID > 0 then
				if ann.ID < systemIndex then
					AddOneAnn(ann)
				end
			end
		end
			
		-- 优先加入系统讯息跑马灯
		for i, ann in ipairs(receiveMsg.Ann) do
			if ann.ID < systemIndex then
				AddOneAnn(ann)
			end
		end	
		
		-- 再加入其他讯息跑马灯
		for i, ann in ipairs(receiveMsg.Ann) do
			if ann.ID > systemIndex then
				AddOneAnn(ann)
			end
		end
	end
end , groupName);