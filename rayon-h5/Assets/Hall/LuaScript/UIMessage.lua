local G = require "Global"

-- 接收消息table
local pojectName = G.GetGamePName()

local showTime = 2
local accShowTime = 0
local movespeed = 30

function Init(this)
end

function Start()
	accShowTime = 0
end

function Update()
	accShowTime = accShowTime + CS.UnityEngine.Time.deltaTime
	
	if accShowTime < 2 then
		local pos = this.gameObject.transform.localPosition
		pos.y = pos.y + CS.UnityEngine.Time.deltaTime * movespeed
		this.gameObject.transform.localPosition = pos
	end
	
	if accShowTime > showTime then
		CS.UnityEngine.GameObject.Destroy(this.gameObject);	
	end
end
