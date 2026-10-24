require("Module")
mmodule(...)

local G = require "Global"
local groupName = "VipService";

local vipDataList =
{
	-- { lv = 0, inner = "1.获得专属头相框", nextexp = 10 },
	-- { lv = 1, inner = "1.获得专属头相框", nextexp = 100 },
    -- { lv = 2, inner = "1.获得专属头相框", nextexp = 200 },
    -- { lv = 3, inner = "1.获得专属头相框", nextexp = 300 },
    -- { lv = 4, inner = "1.获得专属头相框", nextexp = 400 },
	-- { lv = 5, inner = "1.获得专属头相框", nextexp = 500 },
	-- { lv = 6, inner = "1.获得专属头相框", nextexp = 600 },
	-- { lv = 7, inner = "1.获得专属头相框", nextexp = 700 },
	-- { lv = 8, inner = "1.获得专属头相框", nextexp = 800 },
	-- { lv = 9, inner = "1.获得专属头相框", nextexp = 900 },
	-- { lv = 10, inner = "1.获得专属头相框", nextexp = 1000 },
}

local HeadData = nil

-- 檢查資料是否已經取得
function ChechVipDateFinish(callback)
	Callback = callback
	if #vipDataList == 0 then
		GetVipData()
	else
		if Callback ~= nil then
			Callback()
			Callback = nil
		end
	end
end

-- 取得資料
function GetVipDataByLevel(Level)
	for _, data in pairs(vipDataList) do
		if data.lv == Level then
			return data
		end
	end

	return nil
end

-- 取得資料
function GetVipAllData()
	return vipDataList
end

-- 取得資料數量
function GetVipDataCount()
	return #vipDataList
end

function RegRetCallBack(callback)
	G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetGradeSettingDataRetID, function(receiveMsg)
		vipDataList = {}
		
		for i, data in ipairs(receiveMsg.DataList) do
			local vipData = {}
			vipData.lv 		= data.ID
			vipData.frame 	= data.AvatarFrame
			vipData.inner 	= "1.获得专属头相框"
			vipData.inner1 	= "2.会员晋级奖励金"
			vipData.nextexp = data.Amount
			vipData.add 	= data.PromotionAward
			table.insert(vipDataList, vipData)
		end

		if callback ~= nil then
			callback()
		end
	end , groupName);
end

function GetVipData()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetGradeSettingDataID, nil)
end

function SetFrame(image, value)
	if image == nil then
		return false
	end

	--if value == 0 then
		--return false
	--end
	
	if #vipDataList == 0 then
		return false
	end
	
	image:ResetSpriteByName(CS.AppDefine.HallNumber, "f" .. tostring(value))
	
	return true
end

function SetHeadMask(image)
	if image == nil then
		return
	end

	if HeadData == nil then
		return
	end
	
	image.material = HeadData.material	
end

function SetHeadMaskAndFrame(head, frame, value)
	local result = SetFrame(frame, value)
	if result == true then
		SetHeadMask(head)	
	end
end

function SaveHeadData(data)
	HeadData = data
end