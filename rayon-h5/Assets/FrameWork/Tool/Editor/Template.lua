local G = require "Global"

-- 接收消息table
local pojectName = G.GetGamePName()

function Init(this)
end

function Start()
end

function Update()
end

function OnTriggerEnter2D(collision)
end

function OnTriggerStay2D(collision)
end

function OnTriggerExit2D(collision)
end

function OnCollisionEnter2D(collision)
end

function OnCollisionStay2D(collision)
end

function OnCollisionExit2D(collision)
end

function OnDestroy()

    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(pojectName)
	
end