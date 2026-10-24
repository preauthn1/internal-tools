-- region *.lua
-- Date
-- 全局拦截所有游戏场景的事件，会比游戏场景中的事件先执行

local G = require "Global"
G.RequireG("ProtoConfig").Init()
function Start()   
    G.Tool.UILoading.Hide()        
end

function Destory()

end


-- endregion
