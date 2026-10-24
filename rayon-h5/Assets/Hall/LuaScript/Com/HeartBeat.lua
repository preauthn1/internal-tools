require("Module")
mmodule(...)

--协议号
Proto = require("Com/Proto")
--网络模块
NetQ = require("Com/NetworkQueue")

local beatInterval = 0
local connLoseTimeSpan = 0

--开始心跳消息
local function beginHeartBeatMsg(hbtype)
    if beatInterval > 0 and connLoseTimeSpan > 0 then
        print("lua通知开始心跳.")
        CS.HeartBeatManager.Instance:StartHeartBeat(hbtype, beatInterval, connLoseTimeSpan)

        return true
    end

    return false
end

--发送请求心跳配置 
--hbtype:0 大厅服务器心跳 1 游戏服务器心跳
local function sendRequestHeartBeatCfg(hbtype)
    if hbtype == 0 then
        print("向大厅服务器请求心跳配置")
        NetQ.SendHallMsg(Proto.BClass.Common, Proto.Common.RequestHeartBeatConfigID)
    else
        print("向游戏服务器请求心跳配置")
        NetQ.SendGameMsg(Proto.BClass.Common, Proto.Common.RequestHeartBeatConfigID)
    end        
end

--启动心跳机制
--hbtype:0 大厅服务器心跳 1 游戏服务器心跳
local function startHeartBeat(hbtype)
   print("启动心跳," .. hbtype)
   if beginHeartBeatMsg(hbtype) == true then
        return
   end

   --接收心跳配置信息
   local recvFunc = function (beatcfgMsg)
         beatInterval = beatcfgMsg.BeatInterval
         connLoseTimeSpan = beatcfgMsg.ConnLoseTimeSpan

         print("接收到心跳配置." .. beatInterval .. "," .. connLoseTimeSpan)

         --开始发送心跳数据包
         beginHeartBeatMsg(hbtype)
   end

   NetQ.RegOnceMsgHook(Proto.BClass.Common, Proto.Common.HeartBeatConfigID, recvFunc)

   --请求心跳配置
   sendRequestHeartBeatCfg(hbtype)
end



--开始大厅服务器心跳包
function StartHallServerHeartBeat()
    startHeartBeat(0)
end

--开始游戏服务器心跳包
function StartGameServerHeartBeat()
    startHeartBeat(1)
end
