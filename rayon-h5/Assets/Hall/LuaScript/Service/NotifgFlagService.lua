require("Module")
mmodule(...)

local G = require "Global"

-- 是否有充值到账通知
function IsPayNotify()
    return G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.NotifyFlag, 1, 1) == "1"
end


-- 是否有新邮件
function IsNewEmail()    
    return G.Tool.StringHelper.SubString(G.User.GetCurrentUser().UserData.NotifyFlag, 2, 1) == "1"
end

