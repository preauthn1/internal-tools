require("Module")
mmodule(...)

local agentTable = nil

function Clear()
    agentTable = nil
end
 
-- 设置代理信息
function SetAgentInfo(agentList) 
    agentTable = agentList
end

-- 获取代理信息
function GetAgentInfo() 
    return agentTable
end 