-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成

require("Module")
mmodule(...)

local G = require "Global"
local str_StopCo = "StopCoroutine"  -- 退出协程函数关键字符串

local co_receiveGameVerMsg          -- 接收游戏本版协程引用
local co_connectionGame             -- 连接游戏服务器协程引用
local co_confirmReturnHall          -- 确认是否返回大厅协程引用
local co_waitReConnectGame          -- 等待重连游戏协程

local co_confirmExitApp             -- 确认是否退出程序协程
local co_waitReConnectHall          -- 等待重连大厅协程


local RequestGameID                 -- 请求的游戏ID
local GameSkin                      -- 请求的游戏 Skin

local yield_return_HallConnect = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback("HallConnect", to_to_yield, cb)
end
)
local yield_return_RoomConnect = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback("RoomConnect", to_to_yield, cb)
end
)
local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)


local updateGameUiLoadingComplatedCallBack = nil
function SetUpdateGameUiLoadingComplatedCallBack(complatedCallBack)
	updateGameUiLoadingComplatedCallBack = complatedCallBack
end

local updateGameComplatedCallBack = nil
function SetUpdateGameComplatedCallBack(complatedCallBack)
	updateGameComplatedCallBack = complatedCallBack
end

local updateGameDownloadNextCallBack = nil
function SetUpdateGameDownloadNextCallBack(downloadNextCallBack)
	updateGameDownloadNextCallBack = downloadNextCallBack
end

local updateGameErrorCallBack = nil
function SetUpdateGameErrorCallBack(errorCallBack)
	updateGameErrorCallBack = errorCallBack
end

local updateGameFaildCallBack = nil
function SetUpdateGameFaildCallBack(callback)
	updateGameFaildCallBack = callback
end

local isAutoLoginGame = true
function SetIsAutoLoginGame(flag)
	isAutoLoginGame = flag
end

local ReConnectHall = true
function SetNotReConnectHall()
	ReConnectHall = false
end

--支援版本[CallbackHotUpdateFail]，下載失敗時是否跳過詢問視窗
local function SetHotUpdateFaildPassMessageBox(yes_or_no)
	if (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail")) == false then
		return
	end
	
	if (CS.HotUpdate.HotUpdateManager.s_HaveHotUpdateFaild) ~= true then
		return
	end
	
	CS.HotUpdate.HotUpdateManager.Instance:SetHotUpdateFaildPassMessageBox(yes_or_no)
end

--支援版本[CallbackHotUpdateFail]，下載失敗的事件
local function SetCallbackHotUpdateFaild(callback)
	if (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail")) == false then
		return
	end
	
	if (CS.HotUpdate.HotUpdateManager.s_HaveHotUpdateFaild) ~= true then
		return
	end
	
	CS.HotUpdate.HotUpdateManager.Instance:SetCallbackHotUpdateFaild(function()
		CS.HotUpdate.HotUpdateManager.Instance:SetCallbackHotUpdateFaild(nil);
		
		local pass = (CS.NetworkManager.Instance.RoomSocket.MyConnectState == CS.ConnectState.Connect)
		SetHotUpdateFaildPassMessageBox(pass)
		
		if(callback ~= nil) then
			callback()
		end
		
		if(pass == true)then					
			if updateGameDownloadNextCallBack ~= nil then
				updateGameDownloadNextCallBack()
			end			
		end		
	end)
end

--支援版本[CallbackHotUpdateFail]，下載失敗，詢問視窗點取消按鈕的事件
local function SetCallbackHotUpdateFailCancel(callback)
	if (G.Tool.SupportVersion.CheckHall("CallbackHotUpdateFail")) == false then
		return
	end

	if (CS.HotUpdate.HotUpdateManager.s_HaveHotUpdateFaild) ~= true then
		return
	end
	
	CS.HotUpdate.HotUpdateManager.Instance:SetCallbackHotUpdateFailCancel(function()
		CS.HotUpdate.HotUpdateManager.Instance:SetCallbackHotUpdateFailCancel(nil);
		
		if(callback ~= nil) then
			callback()
		end
	end)
end

-- 关闭大厅连接
function CloseHallSocket()
    CS.UINetworkMessage:Hide()
    CS.NetworkManager.Instance.HallSocket:Close()
    -- 停掉所有大厅连接相关的协程
    CS.TimeManager.Instance:StopAllCoroutines("HallConnect")
    CS.NetworkManager.Instance:RemoveEvents( function(s) return s.MainID ~= G.Proto.BClass.Game end)
end   

-- 关闭游戏连接
function CloseRoomSocket()
    G.Room.GameLoginSuccess(false)
    CS.UINetworkMessage:Hide()
    CS.NetworkManager.Instance.RoomSocket:Close()
    -- 停掉所有房间连接相关的协程 
    CS.TimeManager.Instance:StopAllCoroutines("RoomConnect")
    CS.NetworkManager.Instance:RemoveEvents( function(s) return s.MainID == G.Proto.BClass.Game end)
end   

-- 发送Token给遊戲服务器
function SendTokenMsg()
	local TokenData = G.User.GetFullToken()
	G.NetQ.SendGameMsg(G.Proto.BClass.Common, G.Proto.Common.ConnTokenID, { token = [REDACTED] })
end

-- 连接游戏服务器
local function connectionGame(call_success, call_error)
    local ipAddress, port = G.CurrentGameIPAndPort()
    print("function Start连接游戏服务器",ipAddress, port)
    if (co_connectionGame ~= nil) then
    
     coroutine.resume(co_connectionGame, str_StopCo)
      end
    co_connectionGame = coroutine.create( function()
        CS.DG.Tweening.DOTween.KillAll()
        print("尝试连接游戏服务器" .. ipAddress .. ":" .. port)
        local co_val = yield_return_RoomConnect(CS.NetworkManager.Instance:SendRoomConnect(ipAddress, port))
        if (co_val == str_StopCo) then return end
        if CS.NetworkManager.Instance.RoomSocket.MyConnectState == CS.ConnectState.Connect then
            G.NetQ.Resume()
			
			SendTokenMsg()
			
            print("启动游戏心跳机制")
            G.HeartBeat.StartGameServerHeartBeat()

            -- 发送本地IP给服务器
            G.Server.SendLocalIPMsg(1)

            if (type(call_success) == "function") then
                print("连接成功,回调成功函数")
                G.Room.GameLoginSuccess(false)
				call_success()
            end
        else
            print("连接游戏失败")
            if (type(call_error) == "function") then
                call_error()
            end
        end
    end )
    coroutine.resume(co_connectionGame)
end

local function SendLoginRoomMsg()

    print("发送登录房间消息")
    local LoginGameRoomInfo = { }

    -- 用户ID
    LoginGameRoomInfo.UserID = G.User.GetCurrentUser().UserID
    -- 证书
    LoginGameRoomInfo.Cer = G.User.GetCurrentUser().Cer
    -- 设备码
    LoginGameRoomInfo.HDCode = G.Env.GetHDCode()
    -- 设备类型
    LoginGameRoomInfo.HDType = G.Env.GetHDType()

    print(G.Proto.BClass.GameRoom, G.Proto.GameRoom.LoginRoomID)


    -- 发送登录游戏服务消息
    G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.GameRoom.LoginRoomID, LoginGameRoomInfo)
end 

-- 连接游戏成功处理
local function connectionGameSuccess()


    CS.UINetworkMessage.Hide()

    print("加载游戏场景")
    if G.IsGameStart() == false then
        G.Util.coroutine_call
        (
        function()
            -- 切换场景
            G.Util.async_to_sync(CS.GameManager.Instance.EnterGame)(G.GetGamePName(RequestGameID)) 
            print("当前加载场景数量 " .. CS.UnityEngine.SceneManagement.SceneManager.sceneCount)
            SendLoginRoomMsg()
        end
        )()
    else
        SendLoginRoomMsg()
    end

end

local function reconnectionCheck()
	local gameData = G.GetRunGameData()

	if(RequestGameID==24)then
		return
	end

	G.Tool.TableHelper.Print(gameData)
	if(gameData.GameType == 1) then
		local parm = {}
		parm.GameID = gameData.GameID
		parm.FlagID = tonumber(gameData.Flag)
		parm.UserID = G.User.GetCurrentUserID()
		
		-- 
		G.NetQ.RegOnceMsgHook(G.Proto.BClass.GameRoom, G.Proto.Hall.GameCompetitionRetID, function(receiveMsg)
			-- 比賽還沒有結束
			if receiveMsg.PiSaiTime > 0 then
				connectionGameSuccess()
			else
				G.Service.Game.ExitGame();
				G.IsGameing(false)
			end
		end );

		-- 
		G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.Hall.GameCompetitionID, parm)
	else
		connectionGameSuccess()
	end
end


-- 发送登录游戏消息,自动连接一次游戏
local function sendLoginRoomMessage()
	
	 if(RequestGameID == 24) then
        EnterNewBy()
        return
     end
    

    CS.UILoadingData.Show("游戏加载中...", 10,
    (
    function()
        CS.UIBubblePrompt.Show("服务器连接失败，请稍后再试", "UIHall")
    end
    ))

    connectionGame(connectionGameSuccess)
end

-- 游戏更新完成事件
local function onGameUpdateSucess()
	if updateGameComplatedCallBack ~= nil then
		updateGameComplatedCallBack()
	end

	if isAutoLoginGame == true then
		sendLoginRoomMessage()
	end
	isAutoLoginGame = true
	    	
	SetCallbackHotUpdateFaild(nil);
	SetCallbackHotUpdateFailCancel(nil);
	if updateGameDownloadNextCallBack ~= nil then
		updateGameDownloadNextCallBack()
	end
end

--local function UpdateError()
local function ConnectionGameError()
	--CS.UnityEngine.Debug.LogError("ConnectionGameError")
    CloseRoomSocket()

	--[[if updateGameErrorCallBack ~= nil then
		updateGameErrorCallBack()
		updateGameErrorCallBack = nil
	end
	
	if updateGameDownloadNextCallBack ~= nil then
		updateGameDownloadNextCallBack()
	end]]
end

-- 接收游戏版本号处理
local function receiveGameVerMsg(data,gameId,gameSkin)
    if (co_receiveGameVerMsg ~= nil) then coroutine.resume(co_receiveGameVerMsg, str_StopCo) end
    CS.UIMessageBox.Hide()

    local newVersion;
    print("检查游戏ID", gameId,G.GetGamePName(gameId))
    print("游戏Skin", GameSkin)
	
	local localVersion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName(gameId));

    if G.Env.IsIOS() then
        newVersion = data.IOSVersion
    elseif G.Env.IsAndroid() then
        newVersion = data.AndroidVersion
    elseif G.Env.IsWindows() then
        newVersion = data.AndroidVersion
    end
    print("VERSION CHECK:",newVersion,localVersion)
	if(newVersion == nil or newVersion == "")then		
		CS.UIBubblePrompt.Show("游戏版号错误", "UIHall")
		if (updateGameFaildCallBack ~= nil) then
			updateGameFaildCallBack()
		end
		if (updateGameDownloadNextCallBack ~= nil) then
			updateGameDownloadNextCallBack()
		end
		return
	end
    local checkReslut = CS.HotUpdate.HotUpdateManager.CheckVersion(newVersion, localVersion)

    if checkReslut ~= 0 and checkReslut ~= -1 then
        -- 有版本更新，将Socket断开
        --CloseRoomSocket()
        CS.UILoadingData.Hide()
        co_receiveGameVerMsg = coroutine.create( function()
            local message = "";

            if localVersion == CS.HotUpdate.HotUpdateManager.DEFAULT_VERSION then
                message = "游戏需要下载，点击确定下载游戏";
            else
                message = "游戏需要更新，点击确定更新游戏";
            end

           -- local co_val = showMessageBox(CS.MessageBoxStyle.OKCancel, "", message, "确定", "取消");
           -- if (co_val == str_StopCo) then return end
           -- local boxResult = co_val

           -- if boxResult == CS.MessageBoxResult.OK then
                -- 更新或者下载游戏
				
				local UpdateDomain = G.GameConfig.UpdateDomain
				local AliyunDomain = G.Tool.Aliyun.GetCdnDomain()
				if AliyunDomain ~= "" then
					UpdateDomain = AliyunDomain
				end

				-- 依造版本號取對應檔案 by Rayon
                if(gameSkin=="0" or gameSkin=="") then
			        CS.AppDefine.UpdateFileUrl = UpdateDomain .. "/" .. G.GetGamePName(gameId) .. "/" .. newVersion
                else                    
                    CS.AppDefine.UpdateFileUrl = UpdateDomain .. "/"..gameSkin.."/" .. G.GetGamePName(gameId) .. "/" .. newVersion
                end                
                print("<color=yellow>GAME URL :"..CS.AppDefine.UpdateFileUrl.."</color>"); 
                print("更新游戏ID", gameId,G.GetGamePName(gameId))
				
				CS.HotUpdate.HotUpdateManager.Instance:SetUpdateProjectName(G.GetGamePName(gameId))
                CS.HotUpdate.HotUpdateManager.Instance:TryStartUpdateRemoteAssetsToLocal()

				if (RegisterComplate == true) then
					CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated = CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated - onGameUpdateSucess
					RegisterComplate = nil
				end
				
				if updateGameUiLoadingComplatedCallBack == nil then
					CS.UIUpdateGame.Show();
				else
					CS.UIUpdateGame.Show(updateGameUiLoadingComplatedCallBack);
				end
				
				if (RegisterComplate == nil) then
					CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated = CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated + onGameUpdateSucess
					RegisterComplate = true
				end
				
				SetCallbackHotUpdateFaild(function()
					if (updateGameFaildCallBack ~= nil) then
						updateGameFaildCallBack()
					end
				end);
				--下載失敗，點取消按鈕				
				SetCallbackHotUpdateFailCancel(function()
					if updateGameDownloadNextCallBack ~= nil then
						updateGameDownloadNextCallBack()
					end
				end);
		  --  end
        end )
        assert(coroutine.resume(co_receiveGameVerMsg), "确定更新协程报错")
    else
        sendLoginRoomMessage()
    end
end

function RequestHotUpdate(serverinfo,skin)
	print("RequestHotUpdate",serverinfo.GameID,G.GetGamePName(serverinfo.GameID))
	receiveGameVerMsg(serverinfo,serverinfo.GameID,skin)
end

-- 检查游戏版本
function CheckGameVersion(gameid,skin)
    RequestGameID = gameid
    GameSkin=skin
    CS.UILoadingData.Show("游戏加载中...", 10,
    (
		function()
			CS.UIBubblePrompt.Show("服务器连接失败，请稍后再试", "UIHall")
			ConnectionGameError()
			require("Tools/TestConnect").loginError("服务器连接失败，请稍后再试")
		end
    ))

    if(gameid==24)then

    	local data={};
    	data.AndroidVersion=G.GameConfig.RequestGameRemoteVersion;
    	data.IOSVersion=G.GameConfig.RequestGameRemoteVersion;

		receiveGameVerMsg(data,RequestGameID,GameSkin)
	    return
    end	

    -- 连接游戏服务器
    connectionGame(
		-- 成功处理
		function()
			print("连接到游戏服务器成功")
			CS.UILoadingData.Show("游戏加载中...", 10,
			(
			function()
				print("版本验证错误")
				CS.UIBubblePrompt.Show("版本验证错误", "UIHall")
				ConnectionGameError()
				require("Tools/TestConnect").loginError("版本验证错误")
			end
			))


			--改为以大听取得游戏版本号做验证
		    local data={};
		    data.AndroidVersion=G.GameConfig.RequestGameRemoteVersion;
		    data.IOSVersion=G.GameConfig.RequestGameRemoteVersion;
			receiveGameVerMsg(data,RequestGameID,GameSkin)
			    
		    	

--			G.NetQ.RegOnceMsgHook(G.Proto.BClass.GameRoom, G.Proto.GameRoom.GameVerInfoID, function(data)
--				receiveGameVerMsg(data,RequestGameID,GameSkin)
--			end);
--			G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.GameRoom.RequestGameVerID)
		end
		,
		-- 失败处理
		function()
			print("连接到游戏服务器失败")
			CS.UILoadingData.Hide();
			CS.UIBubblePrompt.Show("服务器连接失败，请稍后再试", "UIHall")
			ConnectionGameError()
			require("Tools/TestConnect").loginError("服务器连接失败，请稍后再试")
		end
    )
end

function EnterNewBy()

	CS.UILoadingData.Show("游戏加载中...", 10,
	(
		function()
		print("连接伺服器失败....")
		CS.UIBubblePrompt.Show("请重新尝试连接", "UIHall")
				
		end
	)) 
    local ipAddress, port = G.CurrentGameIPAndPort()
    print(ipAddress,port)
    if(not G.IsGameStart())  then
	    G.Util.coroutine_call
	    (
	        function()
	            -- 切换场景
	            G.Util.async_to_sync(CS.GameManager.Instance.EnterGame)(G.GetGamePName(RequestGameID))
	            print("当前加载场景数量 " .. CS.UnityEngine.SceneManagement.SceneManager.sceneCount)
	            print("关闭大厅界面")
	            G.GameConfig.EnterGameID = RequestGameID
	            G.SetRunGame(RequestGameID)
	            G.Tool.Sound.CloseBGMSound()
	            CS.WindowManager.Instance:CloseWindow("UIHall", true);                
	            --设置登录游戏状态为成功
	            G.Room.GameLoginSuccess(true)

	             
	            print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"..G.CurrentGameIPAndPort())
	            CS.GameManager.Instance:StartGame();
	        end
	    )()
    else
        CS.FrameWork.ObjectPoolManager.Instance:Retrieve();
        G.Service.Game.ReloadGame( function()
        	--设置登录游戏状态为成功
        	G.Room.GameLoginSuccess(true)
        	print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"..G.CurrentGameIPAndPort())
        	CS.GameManager.Instance:StartGame()
        end )
    end
    print("当前加载场景数量 " .. CS.UnityEngine.SceneManagement.SceneManager.sceneCount)
    CS.Utility.CallGC();
    

    -- 发送请求游戏服务器地址
   -- G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestGameServerAddrID, gamedata)

    
  
end



-- 获取游戏服务器地址
function RequestGameGameServerAddr(gamedata,Skin)
    print("请求游戏服务器地址信息")
    CS.UILoadingData.Show("游戏加载中...", 10,
    (
    function()
        CS.UIBubblePrompt.Show("游戏大厅分配失败", "UIHall")
		ConnectionGameError()
		require("Tools/TestConnect").loginError("游戏大厅分配失败")
    end
    ))
    
    local GameIPAddress = gamedata.IP
    local GamePort = gamedata.Port
         

    

        G.Server.GetGameServerAddrCallback(GameIPAddress, GamePort, function(res, ipAddress, port)
			if res == 0 then
				if G.GameConfig.GetServerOption() == 4 and gamedata.GameID==24  then
         			   G.SetRunGameIPAndPort("0.0.0.0",GamePort)
         	    elseif G.GameConfig.GetServerOption() == 5 and gamedata.GameID==24  then
         	    	   G.SetRunGameIPAndPort("0.0.0.0",GamePort)
		        else
		        	   G.SetRunGameIPAndPort(ipAddress, port)
		        end
				
				CheckGameVersion(gamedata.GameID,Skin)
			else
				ConnectionGameError()
				CS.UIBubblePrompt.Show("游戏连接失败(1)", "UIHall")
				require("Tools/TestConnect").loginError("遊戲盾取得遊戲IP失敗")
			end
		end );

        return;





    -------------------------------------------------------------------------------


  --   -- 拦截游戏服务器地址消息
  --   G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GameServerAddrInfoID, function(receiveMsg)
  --       G.Server.GetGameServerAddrCallback(receiveMsg.IPAddress, receiveMsg.Port, function(res, ipAddress, port)
		-- 	if res == 0 then
		-- 		 if G.GameConfig.GetServerOption() == 5 and gamedata.GameID==24  then
  --        			 G.SetRunGameIPAndPort("0.0.0.0",receiveMsg.Port)
		--         elseif G.GameConfig.GetServerOption() == 4 and gamedata.GameID==24 then
		--           G.SetRunGameIPAndPort("0.0.0.0",receiveMsg.Port)
		--         elseif G.GameConfig.GetServerOption() == 1 and gamedata.GameID==24 then
		--         	print("Warring Assign IP")
		--         	G.SetRunGameIPAndPort("0.0.0.0",receiveMsg.Port)	
		--         else
		--         	G.SetRunGameIPAndPort(ipAddress, port)
		--         end
				
		-- 		CheckGameVersion(gamedata.GameID,Skin)
		-- 	else
		-- 		ConnectionGameError()
		-- 		CS.UIBubblePrompt.Show("游戏连接失败(1)", "UIHall")
		-- 		require("Tools/TestConnect").loginError("遊戲盾取得遊戲IP失敗")
		-- 	end
		-- end );
  --   end );

  --   -- 发送请求游戏服务器地址
  --   G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestGameServerAddrID, gamedata)
end

-- 重新进入游戏(如:断线重连),retryTimes　重连次数,首次调用不传参数
function ReLoginRoom(retryTimes)
    -- if (not G.IsGameing()) then return end
    print("重新连接游戏服务器", retryTimes)
    if (co_confirmReturnHall ~= nil) then coroutine.resume(co_confirmReturnHall, str_StopCo) end
    if (co_waitReConnectGame ~= nil) then coroutine.resume(co_waitReConnectGame, str_StopCo) end

    if (retryTimes == nil) then
        CS.UIMessageBox.Hide()
        CS.UINetworkMessage.Hide()
    end
    local gameRetryTimes = G.CaseWhen(retryTimes == nil, 10, retryTimes)
    --CS.UINetworkMessage.Show("网络不稳定，正在重新连接(" .. gameRetryTimes .. ")...")
    if(RequestGameID==24)then
    	return
    end
    connectionGame(reconnectionCheck, function()
        if (gameRetryTimes == 0) then
            print("connect error")
            CS.UINetworkMessage.Hide()
            co_confirmReturnHall = coroutine.create( function()
                local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "游戏服务器连接断开，点击确定返回大厅。", "确定");
                if (boxResult == str_StopCo) then return end
                if boxResult == CS.MessageBoxResult.OK then
                    G.Service.Game.ExitGame()
                end
            end )
            coroutine.resume(co_confirmReturnHall)
        else
            gameRetryTimes = gameRetryTimes - 1
            co_waitReConnectGame = coroutine.create( function()
                print("等待1秒后再连")
                local co_val = yield_return_RoomConnect(CS.UnityEngine.WaitForSeconds(1))
                if (co_val == str_StopCo) then return end
                print("等待1秒完成马上重连")
                ReLoginRoom(gameRetryTimes)
            end )
            coroutine.resume(co_waitReConnectGame)
        end
    end )

end


-- 重连登录大厅
function ReLogonHall(isAutoLogin, retryTimes)
	-- 被大廳強制關閉，不需要重連
	if ReConnectHall == false then
		return
	end

	if G.IsGameInitFinish() == false then
		return
	end
	
	-- 沒有驗證過版號，完整重連
	if StartUp.GetIsCheckVersion() == false then
		print("沒有驗證過版號，完整重連")
		CloseHallSocket()
		StartUp.BeginConnect()
		return
	end		

    isAutoLogin = G.CaseWhen(isAutoLogin == nil, true, isAutoLogin)
    if (co_confirmExitApp ~= nil) then coroutine.resume(co_confirmExitApp, str_StopCo) end
    if (co_waitReConnectHall ~= nil) then coroutine.resume(co_waitReConnectHall, str_StopCo) end
    if (retryTimes == nil) then
        CS.UIMessageBox.Hide()
    end
    local hallRetryTimes = G.CaseWhen(retryTimes == nil, 10, retryTimes)
    local isGameStart = G.IsGameStart()
    if isGameStart == false then
        CS.UINetworkMessage.Show("网络不稳定，正在重新连接(" .. hallRetryTimes .. ")...")
    end
    print("开始重连大厅", retryTimes, isGameStart)
    G.Server.LoadServerAddr(
    -- 成功处理
    function()
        if isGameStart == false then
            CS.UINetworkMessage.Hide()
            print("是否游戏中", isGameStart)
           -- CS.UIBubblePrompt.Show("登录大厅...", "UIHall")
        end

        G.Service.Login.AutoLogon( function(recvMsg)
			-- 記錄問題 By Rayon
			-- 在登入過程中重新連線後不會切到大廳
			-- 應該切回 StartUp.onLogonRet
			
			if isGameStart == false then
				if CS.WindowManager.Instance:FindWindowNodeByName("UIHall") == nil then
					CloseHallSocket()
					StartUp.BeginConnect()
				end
			end
			
			if recvMsg.Code == 0 then
                CS.UIBubblePrompt.Show(recvMsg.Message .. "请退出程序重新登录.", "UIHall")
			end
        end )
    end
    ,
    -- 失败处理
    function()
        if hallRetryTimes == 0 then
            print("connect error")
            co_confirmExitApp = coroutine.create( function()
                -- 如果已经启动游戏，那么在自动重连之后失败不再处理
                if isGameStart == false then
                    -- 断开连接
                    CloseHallSocket()
                    -- 断开连接
                    CloseRoomSocket()

                    -- 确认
                    CS.UINetworkMessage.Hide()
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "服务器连接断开，请检查网络并尝试重新登录游戏。", "确定");
                    if (boxResult == str_StopCo) then return end

                    if boxResult == CS.MessageBoxResult.OK then
                        CS.UnityEngine.Application.Quit()
                    end
                end
            end )
            coroutine.resume(co_confirmExitApp)
        else
            hallRetryTimes = hallRetryTimes - 1
            print(hallRetryTimes)
            co_waitReConnectHall = coroutine.create( function()
                print("等待1秒后再连")
                local co_val = yield_return_HallConnect(CS.UnityEngine.WaitForSeconds(1))
                if (co_val == str_StopCo) then return end
                print("等待1秒完成马上重连")
                ReLogonHall(true, hallRetryTimes)
            end )
            coroutine.resume(co_waitReConnectHall)
        end
    end
    )
end

-- region 服务器主动断开连接提示
G.NetQ.RegMsgHook(G.Proto.BClass.Common, G.Proto.Common.TipMessageID, function(receiveMsg)

    CS.UILoadingData.Hide()
    CS.UINetworkMessage.Hide()

    G.Tool.TableHelper.Print(receiveMsg)

    -- MsgType 1 房间主动踢出  2 大厅主动踢出 3 普通提示
    if receiveMsg.MsgType == 1 then
        -- 断开连接
        CloseRoomSocket()
        -- 如果有消息才显示k
        if (#receiveMsg.Message > 0) then
            G.Util.coroutine_call
            (
            function()
                -- Param等于1表示金币不足
                if (receiveMsg.Param == 1) then
                    G.Util.coroutine_call
                    (
                    function()
						-- 充值開關 Add By Rayon
						if G.GameConfig.IsRechargeCtx() == false then
							showMessageBox(CS.MessageBoxStyle.OKCancel, "", receiveMsg.Message, "继续", "取消");
							
							if G.IsGameStart() then
								G.Service.Game.ExitGame();							
							end
						else
							local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", receiveMsg.Message, "前往充值", "取消");
							
							if G.IsGameStart() then
								if boxResult == CS.MessageBoxResult.OK then
									local hallParam = { };
									hallParam.OpenRechargeUI = true;
									G.Service.Game.ExitGame(nil, hallParam);
								else
									G.Service.Game.ExitGame();
								end
							else
								if boxResult == CS.MessageBoxResult.OK then
									local hallData = { };
									hallData.OpenRechargeUI = true;
									local node = CS.WindowManager.Instance:FindWindowNodeByName("UIHall")
									node.WindowData = hallData
								end
							end
						end
                    end
                    )()
				elseif (receiveMsg.Param == 4) then
					if G.IsGameStart() then
						-- 間隔一秒再踢出遊戲
						CS.TimeManager.Instance:DelayInvoke(1,
						function()
							G.Service.Game.ExitGame();
						end )					
					end
				-- 切换AB服,断线,指定IP,重新连线
				elseif (receiveMsg.Param == 6) then
					if G.IsGameStart() then
					--print("@@@@@@@@@@@@@@@@@@@@@@ Time " .. CS.System.DateTime.Now.Minute .. ":" .. CS.System.DateTime.Now.Second .. "." .. CS.System.DateTime.Now.Millisecond)
						G.Service.NetworkService.CloseRoomSocket()
						local addressParm = G.Tool.StringHelper.Split(receiveMsg.Message, ":");
						G.Server.GetGameServerAddrCallback(addressParm[1], addressParm[2], function(res, ipAddress, port)
							if res == 0 then
								G.SetRunGameIPAndPort(ipAddress, port)
								-- 間隔一秒再重新連線
								CS.TimeManager.Instance:DelayInvoke(1,
								function()
								--print("@@@@@@@@@@@@@@@@@@@@@@ Time " .. CS.System.DateTime.Now.Minute .. ":" .. CS.System.DateTime.Now.Second .. "." .. CS.System.DateTime.Now.Millisecond)
									G.Service.NetworkService.ReLoginRoom()
								end )
							else
								CS.UIBubblePrompt.Show("与游戏失去连线", "UIHall")
								require("Tools/TestConnect").loginError("遊戲盾取得遊戲IP失敗")
								G.Service.Game.ExitGame();
							end
						end );
					end				
                else
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", receiveMsg.Message, "确定");
                    if boxResult == CS.MessageBoxResult.OK and G.IsGameStart() then
                        G.Service.Game.ExitGame();
                    end
                end
            end
            )()
        end

    elseif receiveMsg.MsgType == 2 then
		ReConnectHall = false
        ReLoginRoom = function() end
        ReLogonHall = function() end
        -- 断开连接
        CloseHallSocket()
        -- 断开连接
        CloseRoomSocket()
        if (#receiveMsg.Message > 0) then
            G.Util.coroutine_call
            (
            function()
                local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", receiveMsg.Message, "确定");

                if boxResult == CS.MessageBoxResult.OK then
                    CS.UnityEngine.Application.Quit()
                end
            end
            )()
        end
	
    else
        G.Util.coroutine_call
        (
        function()
            local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", receiveMsg.Message, "确定");
        end
        )()
    end
end );

-- endregion



-- endregion
