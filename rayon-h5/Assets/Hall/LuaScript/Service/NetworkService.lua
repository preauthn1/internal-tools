require("Module")
mmodule(...)


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---
 
local G = require "Global";


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

local isGameReconnecting = false;
local isHallReconnecting = false;

-- 断线重连次数
local hallRetryTimes = 10

-- 断线重连次数
local gameRetryTimes = 10

-- region 大厅断线重连

-- 重新登录大厅
local function onReLogonHallRet(recvMsg)
    if recvMsg.Code == 0 then
        CS.UIBubblePrompt.Show(receiveMsg.Message .. "请退出程序重新登录.", "UIHall")
    else
        if G.IsGameStart() == false then
            -- CS.UIBubblePrompt.Show("登录大厅成功", "UIHall")
        end
    end
end

local function doRelogonHall()
    G.Server.LoadServerAddr(
    ( function()
        isHallReconnecting = false
        CS.UINetworkMessage.Hide()

        if G.IsGameStart() == false then
            CS.UIBubblePrompt.Show("正在登录大厅服务器...", "UIHall")
        end

        G.Service.Login.AutoLogon(onReLogonHallRet)
    end ),
    ( function()
        if hallRetryTimes == 0 then
            print("connect error")

            isHallReconnecting = false
            G.Util.coroutine_call
            (
            function()
                -- 如果已经启动游戏，那么在自动重连之后失败不再处理
                if G.IsGameStart() == false then
                    -- 断开连接
                    G.Service.NetworkService.CloseHallSocket()
                    -- 断开连接
                    G.Service.NetworkService.CloseRoomSocket()

                    -- 确认
                    CS.UINetworkMessage.Hide()
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "服务器连接断开，请检查网络并尝试重新登录游戏。", "确定");

                    if boxResult == CS.MessageBoxResult.OK then
                        CS.UnityEngine.Application.Quit()
                    end
                end
            end
            )()
        else
            G.Util.coroutine_call
            (
            function()

                if G.IsGameStart() == false then
                    CS.UINetworkMessage.Show("网络不稳定，正在重新连接(" .. hallRetryTimes .. ")...")
                end

                hallRetryTimes = hallRetryTimes - 1
                yield_return_HallConnect(CS.UnityEngine.WaitForSeconds(1))

                doRelogonHall()
            end
            )()
        end
    end ))
end
  
local function HallSocketDisconnetHandler()
    if isHallReconnecting == false then
        isHallReconnecting = true

        print("hall socket disconnect...", G.IsGameStart())
        CS.TimeManager.Instance:DelayInvoke(0.1,
        function()
            hallRetryTimes = 10
            doRelogonHall()
        end )
    end
end;

-- endregion

-- region 游戏断线重连

local function doReLogonGame()

    -- 获取游戏服务器地址
    local ipAddress, port = G.CurrentGameIPAndPort()

    print("尝试重连游戏服务器" .. ipAddress .. ":" .. port)
    G.Util.coroutine_call( function()

        -- 协程连接游戏服务器
        yield_return_RoomConnect(CS.NetworkManager.Instance:SendRoomConnect(ipAddress, port))

        if CS.NetworkManager.Instance.RoomSocket.MyConnectState == CS.ConnectState.Connect then            
            isGameReconnecting = false;

             -- 设置断线重连时候使用的地址和端口
            G.SetRunGameIPAndPort(ipAddress, port)

            print("重连到游戏服务器成功")
            CS.UINetworkMessage.Hide()
            G.NetQ.Resume()
            -- 重连成功，做一次消息解锁

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

            -- 发送登录游戏服务消息
            G.NetQ.SendMsg(G.Proto.BClass.GameRoom, G.Proto.GameRoom.LoginRoomID, LoginGameRoomInfo)
        else
            if gameRetryTimes == 0 then
                print("connect error")
                isGameReconnecting = false;

                G.Util.coroutine_call
                (
                function()
                    -- 断开连接
                    G.Service.NetworkService.CloseRoomSocket()

                    -- 确认
                    CS.UINetworkMessage.Hide()
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "游戏服务器连接断开，点击确定返回大厅。", "确定");

                    if boxResult == CS.MessageBoxResult.OK then
                        G.Service.Game.ExitGame()
                    end
                end
                )()
            else
                G.Util.coroutine_call
                (
                function()
                    CS.UINetworkMessage.Show("网络不稳定，正在重新连接(" .. gameRetryTimes .. ")...")

                    gameRetryTimes = gameRetryTimes - 1
                    yield_return_RoomConnect(CS.UnityEngine.WaitForSeconds(1))

                    doReLogonGame()
                end
                )()
            end
        end
    end )();
end

local function GameSocketDisconnetHandler()
    if isGameReconnecting == false then
        isGameReconnecting = true
        CS.DG.Tweening.DOTween.KillAll()
        print("game socket disconnect...")
        CS.TimeManager.Instance:DelayInvoke(0.1,
        function()
            gameRetryTimes = 10
            doReLogonGame()
        end
        )
    end
end
  
-- endregion

function InitHallNetworkMonitor()
    G.NetQ.OnHallSocketDisconnectCallback = nil
    -- 连接断开检查
    G.NetQ.OnHallReConnectCheck = function(callback)

        CS.UIBubblePrompt.Show("正在连接服务器", "startup")

        G.Server.LoadServerAddr(( function()
            CS.UIBubblePrompt.Show("连接服务器成功", "startup")

            if callback ~= nil then
                callback()
            end

        end ),( function()
            CS.UIBubblePrompt.Show("无法连接到服务器,请检查您的网络是否正常.", "startup")
        end ))
    end
end

-- 开启监听大厅重连
function RunHallNetworkMonitor()
    isHallReconnecting = false
    G.NetQ.OnHallSocketDisconnectCallback = HallSocketDisconnetHandler
    G.NetQ.OnHallReConnectCheck = HallSocketDisconnetHandler
end

-- 开启监听游戏重连
function RunGameNetworkMonitor()
    isGameReconnecting = false;
    G.NetQ.OnGameSocketDisconnectCallback = GameSocketDisconnetHandler
end
 
-- 关闭大厅连接
function CloseHallSocket()
    CS.UINetworkMessage:Hide()
    isHallReconnecting = false
    CS.NetworkManager.Instance.HallSocket:Close()
    G.Server.ResetIsLoading()
    -- 停掉所有大厅连接相关的协程
    CS.TimeManager.Instance:StopAllCoroutines("HallConnect")
    CS.NetworkManager.Instance:RemoveEvents( function(s) return s.MainID ~= G.Proto.BClass.Game end)
end   

-- 关闭游戏连接
function CloseRoomSocket()
    CS.UINetworkMessage:Hide()
    isGameReconnecting = false
    CS.NetworkManager.Instance.RoomSocket:Close()
    -- 停掉所有房间连接相关的协程
    CS.TimeManager.Instance:StopAllCoroutines("RoomConnect")
    CS.NetworkManager.Instance:RemoveEvents( function(s) return s.MainID == G.Proto.BClass.Game end)
end   

-- region 服务器主动断开连接提示
G.NetQ.RegMsgHook(G.Proto.BClass.Common, G.Proto.Common.TipMessageID, function(receiveMsg)

    -- CS.UIMessageBox.Hide()
    CS.UILoadingData.Hide()
    CS.UINetworkMessage.Hide()

    -- MsgType 1 房间主动踢出  2 大厅主动踢出 3 普通提示
    if receiveMsg.MsgType == 1 then
        -- 断开连接
        G.Service.NetworkService.CloseRoomSocket()

        G.Util.coroutine_call
        (
        function()
            if (string.find(receiveMsg.Message, "金币不足") ~= nil) then
                G.Util.coroutine_call
                (
                function()
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
                )()
            else
                local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", receiveMsg.Message, "确定");
                if boxResult == CS.MessageBoxResult.OK and G.IsGameStart() then
                    G.Service.Game.ExitGame();
                end
            end
        end
        )()
    elseif receiveMsg.MsgType == 2 then
        -- 断开连接
        G.Service.NetworkService.CloseHallSocket()
        -- 断开连接
        G.Service.NetworkService.CloseRoomSocket()

        G.Util.coroutine_call
        (
        function()
            local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", receiveMsg.Message, "确定");

            if boxResult == CS.MessageBoxResult.OK then
                CS.UnityEngine.Application.Quit()
            end
        end
        )()
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
 