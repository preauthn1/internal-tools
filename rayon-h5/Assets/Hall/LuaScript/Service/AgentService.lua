-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成
require("Module")
mmodule(...)

local G = require "Global"

local agentContent;

function GetAgentData(callfun)
    
        -- 拦截代理信息
        G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.RequestAgentListRetID, function(data)
          callfun(data)
        end )
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestAgentListID, nil)
end

function ApplyForDownline(callfun,senddata)
    
        -- 拦截代理信息
        G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.JoinAgentListRetID, function(data)
         callfun(data)
        end )        
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.JoinAgentListID, senddata )
end

function ApplyForAgent(callfun,sendMsg)
    
        -- 拦截代理信息
        G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ApplyAgentListRetID, function(data)
          callfun(data)
        end )
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ApplyAgentListID, sendMsg)
end

function CheckForAgent(callfun)
    
        -- 拦截代理信息
        G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.CheckAgentRetID, function(data)
          callfun(data)
        end )
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.CheckAgentID, nil)
end

-- endregion
