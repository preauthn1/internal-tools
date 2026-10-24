require("Module")
mmodule(...)

local G = require "Global"

local logList

local groupName = "ExchangelogService"

function ClearLoadLogList()
    logList = nil
end

function GetLogList()
    return logList
end

function InitLoadLogList()
    if (logList == nil) then
        RefreshLogList()
    end
end

function RefreshLogList(status,range)
local sendMsg = { }
            -- //A {0:"??",1:"?֤",2:"?ֳq?",3:"w",4:"??ݱ`",5:"ݱ`??",6:"¦W?",7:"????",8:"???",9:"æ~?",10:"?? ",11:"h??",99:""}
            -- //1(@P),2(bӤ),3(@Ӥ),4(TӤ)
                sendMsg.TnyStatus = (status==nil and 99 or status)
                sendMsg.LogRange =  (range==nil and 4 or range)
             
                G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetConvertLogID, sendMsg)
   
end
 
function RegLogRetCallBack(callback)
    G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetConvertLogRetID, function(receiveMsg)
        logList = receiveMsg.ConvertLogList
        print( G.Service.ExchangeLog.GetLogList())
        if callback ~= nil then
            callback()
        end
    end , groupName)
end
-- ƳʼϢ
function RemoveRegLogMsg()
    G.NetQ.RemoveMsgHookByGroup(groupName)
end

