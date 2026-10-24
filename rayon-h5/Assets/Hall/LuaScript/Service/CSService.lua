require("Module")
mmodule(...)

local G = require "Global"

local CSList = {}

local groupName = "CSService"

function ClearCSList()
    CSList = nil
end

function HaveUnRead()
     
    return CSList.Count>0;
end

function GetUnReadCSCount()
    if (CSList ~= nil) then
		return CSList.Count;
	else
		return 0
	end
end

function InitLoadCS()
    if (CSList == nil) then
        RefreshCSList();
    end
end

function RefreshCSList()
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.CheckCSNotifyID, nil)
end
 
function RegCSRetCallBack(callback)
    G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.CheckCSNotifyRetID, function(receiveMsg)
        CSList = receiveMsg

        if callback ~= nil then
            callback()
        end
    end , groupName)
end
-- ƳʼϢ
function RemoveRegCSMsg()
    G.NetQ.RemoveMsgHookByGroup(groupName)
end