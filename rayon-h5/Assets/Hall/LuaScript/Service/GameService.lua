require("Module")
mmodule(...)

local G = require "Global"

-- 场景初始化
local onGameSceneInit = nil

-- 桌子号改变事件(主要用于水果机之类的百人游戏)
local tableNoChange = nil

-- region 公共协议拦截

-- 用户入桌消息
function RegUserEnterMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.UserEnterID, function(data)
        --G.Tool.TableHelper.Print(data)
        callfun(data)
        if (data.User.UserID == G.User.GetCurrentUser().UserID and data.User.StateValue == 2) then
            --print("发送准备请求")
            ReadyGame()
        end
    end , G.GetGamePName())
end

-- 用户状态改变消息
function RegUserStateChangeMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.UserStateChangeID, function(data)
        if (data.User.UserID == G.User.GetCurrentUser().UserID) then
			
			require("Tools/TestConnect").loginFinish()
			
			-- 同步每一回合玩家身上的錢
			G.User.GetCurrentUser().UserData.CashAmount = data.User.Score
			-- 用戶錢改變
			local node = CS.WindowManager.Instance:FindWindowNodeByName("UIRecharge")
			if node ~= nil then
				local Txt_CashAmount = node.WindowGameObject.transform:Find("UIBody/_OnlinePayDetail/_CashAmount"):GetComponent("Text")
				if Txt_CashAmount ~= nil then
					Txt_CashAmount.text =  G.Tool.UserHelper.ToCash(data.User.Score)
				end
			end
			
            -- 处理用户桌子号改变时，修改顶部房间信息
            --G.Tool.TableHelper.Print(data)
            if (tableNoChange ~= nil) then
                if (data.User.TableNo >= 0) then
                    tableNoChange(data.User.TableNo)
                end
            end

            -- 设置全局用户状态值
            G.User.SetUserState(data.User)

            if (data.User.StateValue > 3) then
                G.IsGameing(true)
            else
                G.IsGameing(false)
            end
        end
        callfun(data);
    end , G.GetGamePName())
end

-- 游戏开始消息
function RegGameStartMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.GameStartID, function(data)
        callfun(data)
		
		-- 百人場顯示充值訊息
		--if G.Game.IsMultiGame(G.GetRunGameID()) == true then
			G.Service.Msg.ShowMsg(1, false)
		--end
		
		--[[开启五星功能，记录游玩次数
		local key = "play_times"
		local times = CS.UnityEngine.PlayerPrefs.GetInt(key,0)
		CS.UnityEngine.PlayerPrefs.SetInt(key,times +1)
		CS.UnityEngine.PlayerPrefs.Save()]]
		
    end , G.GetGamePName())
end

-- 游戏结束消息
function RegGameEndMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.GameEndID, function(data)
        callfun(data);
    end , G.GetGamePName())
end

-- 用户离开消息
function RegUserLeaveMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.UserLeaveID, function(data)
        callfun(data);
    end , G.GetGamePName())
end

-- 截拦当前桌所有玩家动画完成消息
function RegAllUserActionCommitedMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.AllUserActionCommitID, function(data)
        --print("收到动画通知", data.ActionName)
        callfun(data.ActionName);
    end , G.GetGamePName())
end

-- 旁观用户消息通知
function RegLookOnUserNotify(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.LookOnUserNotifyID, function(data)
        --print("收到旁观用户通知")
        callfun(data)
    end , G.GetGamePName())
end

-- 遊戲倒數通知
function RegGameBeginCountDownMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.GameBeginCountDownID, function(data)
        callfun(data);
    end , G.GetGamePName())
end

-- 遊戲停止倒數通知
function RegGameBeginCountDownStopMsg(callfun)
    G.NetQ.RegMsgHook(G.Proto.BClass.Game, G.Proto.Game.GameBeginCountDownStopID, function(data)
        callfun(data);
    end , G.GetGamePName())
end

-- endregion

-- 初始化游戏顶部菜单界面
-- _parent顶部菜单附在指定的节点
-- isCustomExit是否需要自定义退出确认界面
-- isCustomSetUI是否需要自定义設定界面 By Rayon
function InitGameTopMenu(_parent, isCustomExitUI, isCustomSetUI,annPos)
    isCustomExitUI = G.CaseWhen(isCustomExitUI == nil, false, isCustomExitUI)
    isCustomSetUI = G.CaseWhen(isCustomSetUI == nil, false, isCustomSetUI)
    local initParas = {
        isCustomExitUI = isCustomExitUI;
		isCustomSetUI = isCustomSetUI;
		annPos = annPos;
    }
    CS.WindowManager.Instance:OpenGameWindow("UIGameTopMenu", CS.AppDefine.HallNumber, _parent, initParas);
	SetCallbackExit(nil)
end

local CallbackExit = nil
function SetCallbackExit(callback)
	CallbackExit = callback
end

function GetCallbackExit()
	return CallbackExit
end

-- 发送客户端动作完成消息
function SendUserActionCommit(data, isLookOn)
    isLookOn = G.CaseWhen(isLookOn == nil, false, isLookOn)
    if (not isLookOn) then
        --print("发送动画消息", data)
        G.NetQ.SendMsg(G.Proto.BClass.Game, G.Proto.Game.UserActionCommitID, { ActionName = data })
    else
        --print("旁观用户不发送动画消息", data)
    end
end

-- 发送准备消息
function ReadyGame()
	-- 牌桌型顯示充值訊息
	local GameID = G.GetRunGameID()
	--if GameID ~= 14 then
		--if G.Game.IsMultiGame(GameID) == false then
			G.Service.Msg.ShowMsg(1, false)
		--end
	--end

    G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.GameRoom.GameReadyID)
end

-- 发送换桌消息
function ChangeDesk()
    G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.GameRoom.UserQueueID)
end

function CloseScene()
	CS.GameManager.Instance.QuitGame()
end

-- 退出游戏场景
function ExitGame(callfun, hallData,callfunBeforeStop)
            if (type(callfunBeforeStop) == "function") then
                callfunBeforeStop()
            end
    G.Tool.Sound.CloseSound()
    RemoveLoadGameLua()
    if G.IsGameStart() then
        CS.WindowManager.Instance:CloseWindowByProjectName(G.GetGamePName())
    end
    roomNameChange = nil
    G.ClearRunGame()
    CS.FrameWork.ObjectPoolManager.Instance:ClearAllPool();
    G.ClearPreLoadGameResState()

    G.Service.NetworkService.CloseRoomSocket()
    onGameSceneInit = nil
    local exitGame = G.Util.async_to_sync(CS.GameManager.Instance.QuitGame)
    G.Util.coroutine_call( function()

		-- 離開遊戲立刻顯示顯示充值訊息
		G.Service.Msg.ShowMsg(1, false)
	
        --print("退出游戏场景")
        G.NetQ.Resume()
        G.Tool.UILoading.Hide()

		if CS.WindowManager.Instance:FindWindowNodeByName("UIHall") == nil then
			local openparam = CS.WindowNodeInitParam("UIHall", CS.AppDefine.HallNumber)
			-- 返回大厅时需要回传给大厅的数据
			openparam.WindowData = hallData
			openparam.NodeType = 0
			CS.WindowManager.Instance:OpenWindow(openparam)
		end
		
        exitGame( function()

            if (type(callfun) == "function") then
                callfun()
            end
        end )
    end )()
end

function ReloadGame(callfun)
    local reLoadGame = G.Util.async_to_sync(CS.GameManager.Instance.ReLoadGame)
    G.Util.coroutine_call( function()
        reLoadGame( function()
            if (type(callfun) == "function") then callfun() end
        end )
    end )()
end

-- 注册游戏场景初始化方法
function RegNotifyGameSceneInit(callfun)
    -- onGameSceneInit = callfun
end 

-- 注册桌子号改变事件
function RegTableNoChange(callfun)
    tableNoChange = callfun
end

function RemoveTableNoChange()
    tableNoChange = nil
end

function RemoveLoadGameLua()
    for k, _ in pairs(package.loaded) do
        if (string.find(k, string.format("/%s/", G.GetGamePName()))) then
            package.loaded[k] = nil
            package.preload[k] = nil
        end
    end
end
