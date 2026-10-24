require("Module")
mmodule(...)

local G = require "Global";

local groupName = "UIHall";

-- 是否自動連線遊戲
local isTestConnect = false

-- 大廳到登入遊戲間隔時間
local loginTime = 2
-- 進入遊戲成功到離開遊戲間隔時間
local exitTime = 5

local gamelist = nil
local autolist = nil

local hallGameClick
local hallRoomClick

local nowGame

local testList = {}

local isInit = false
local isLoginFinish = false

function init(TestConnect, GameClick, RoomClick)
	isTestConnect = TestConnect
	hallGameClick = GameClick
	hallRoomClick = RoomClick	
end

function setServerList(game, auto)
	-- 關掉登入跳出的視窗
	if isTestConnect == false then
		return
	end
	
	-- 關掉登入跳出的視窗
	CS.WindowManager.Instance:CloseWindow("UIBindingGift", false);
	CS.WindowManager.Instance:CloseWindow("UIMail", false);

	gamelist = game
	autolist = auto
	
	if isInit == false then
		isInit = true
		
		login()	
	end
end

function getOneData()
	-- 沒資料
	if #gamelist == 0 then
		return nil
	end

	for gameId, subList in pairs(gamelist) do
		-- 選廳遊戲
		if (not autolist[gameId]) then

		-- 直連遊戲
		else
			for id, data in pairs(subList) do
				if tostring(data.GameID) ~= "23" then
					if #testList == 0 then
						return data
					end
					
					--確認是否連過
					local isFind = false
					for v, testData in pairs(testList) do
						if tostring(data.GameID) == tostring(testData.GameID) and tostring(data.Flag) == tostring(testData.Flag) then						
							isFind = true
							break
						end
					end
					
					if isFind == false then
						return data
					end
				end
			end
		end
	end
	
	return nil
end

function login()
	CS.TimeManager.Instance:DelayInvoke(loginTime,
	function()
		doLogin();
	end )
end

function doLogin()
	if isTestConnect == false then
		return
	end
	
	isLoginFinish = false
	
	nowGame = getOneData()
	
	if nowGame ~= nil then	
		ShowMessage("連線到 " .. nowGame.GameName .. " " .. tostring(nowGame.GameID) .. " - " .. tostring(nowGame.Flag) );
	
		table.insert(testList, nowGame)
		hallRoomClick(autolist[nowGame.GameID], 1);
	-- 連完一輪或完全沒資料
	else
		ShowMessage("清除資料");
		
		testList = {}
	
		login()		
	end
end

function loginFinish()
	if isTestConnect == false then
		return
	end
	
	if G.IsGameStart() == false then
		return
	end
	
	if isLoginFinish == true then
		return
	end
	
	isLoginFinish = true
	
	ShowMessage(nowGame.GameName .. " 連線成功");
	
	CS.TimeManager.Instance:DelayInvoke(exitTime,
	function()
		isInit = false
	
		ShowMessage("離開" .. nowGame.GameName);
		G.Service.Game.ExitGame();
	end )
end

function loginError(message)
	if isTestConnect == false then
		return
	end
	
	G.Tool.Debug.SendTelegram(message);
	--G.Tool.Debug.SendTelegram("hello from unity","b");
	--G.Tool.Debug.SendTelegram("hello from unity","c");
	--G.Tool.Debug.SendTelegram("hello from unity","r");
	
	ShowMessage(nowGame.GameName .. " " .. message);

	login()
end

function ShowMessage(message)
	CS.UnityEngine.Debug.LogWarning("@@@@@@@@@@@@@@ ShowMessage " .. message)
	CS.UIBubblePrompt.Show(message, groupName);
end