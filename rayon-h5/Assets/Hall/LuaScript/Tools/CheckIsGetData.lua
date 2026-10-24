require("Module")
mmodule(...)

function Create()

	local target = {}
	local Name
	local IsGetData = false
	local CheckTime = 1
	local AccCheckTime = 0
	local SendGetFunCB = nil
	
	function Log(Text)
		--print(Text)
	end
	
	-- 初始化
	target.Init = function(sendGetFun, name, isFinish)
		Log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ Init " .. name)
		SendGetFunCB = sendGetFun
		Name = name
		IsGetData = isFinish
		
		if IsGetData == false then
			target.StartGet()
		end
	end
	
	target.StartGet = function(sendGetFun)
		Log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ sendGetFun " .. Name)
		IsGetData = false
		AccCheckTime = 0
		
		if SendGetFunCB ~= nil then
			SendGetFunCB()
		end
	end	
	
	-- 檢查是否超過時間
	target.UpdateGetData = function()
		if IsGetData == true then
			return
		end
		
		AccCheckTime = AccCheckTime + CS.UnityEngine.Time.deltaTime
		if AccCheckTime < CheckTime then
			return
		else
			AccCheckTime = 0
		end
				
		if SendGetFunCB ~= nil then
			Log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ SendGetFunCB " .. Name)
			SendGetFunCB()
		end
	end
	
	-- 設置成已取得
	target.SetGetDataFinish = function()
		Log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ SetGetDataFinish " .. Name)
		IsGetData = true
	end
	
	-- 是否已取得
	target.IsGetDataFinish = function()
		return IsGetData
	end	

	return target
end