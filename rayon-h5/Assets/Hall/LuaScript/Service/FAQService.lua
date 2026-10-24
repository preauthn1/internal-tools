-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成
require("Module")
mmodule(...)

local G = require "Global"

local faqContent;

function GetFAQ(callfun)
    if (faqContent == nil) then
        -- 拦截FAQ信息
        G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetFAQRetID, function(data)
            if (data.Content ~= nil and #data.Content > 0) then
                faqContent = data.Content
                callfun(faqContent)
            end
        end )
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetFAQID, { PlatformID = G.GameConfig.GamePlatformID })
    else
        callfun(faqContent)
    end

end


-- endregion
