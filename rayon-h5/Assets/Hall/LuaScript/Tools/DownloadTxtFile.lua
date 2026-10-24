require("Module")
mmodule(...)

local G = require "Global"

-- 重試間隔
local m_TryAccTime = 5
-- 判斷是否重間隔
local m_InvokeTime = 0.1
-- 累計時間
local m_AccTime = 0
-- 抓取完成callback
local m_onSuccessCallback
-- 抓取失敗callback
local m_onErrorCallback
-- 下載WWW
local m_www
-- 下載路徑
local m_url

function ShowLog(str)
	print(str)
end

function DownloadFile(url, onSuccessCallback, onErrorCallback)
	m_url = url
	m_onSuccessCallback = onSuccessCallback
	m_onErrorCallback = onErrorCallback
	
	BeginDownload()
end

function BeginDownload()
	m_AccTime = 0
		
	-- 支援填完整路徑
	_url = m_url .. "?abc=" .. tostring(math.random(1, 100000))
	
	ShowLog("取得aliyun檔案 " .. _url)

	m_www = CS.UnityEngine.WWW(_url)

	DoDownloadFile()
end	
	
function DoDownloadFile()
	CS.TimeManager.Instance:StopAllDelayInvoke()
				
	CS.TimeManager.Instance:DelayInvoke(m_InvokeTime, function( )
		
		m_AccTime = m_AccTime + m_InvokeTime
		
		-- 超過重新抓取間隔，重新抓取
		if m_AccTime > m_TryAccTime then
			ShowLog("!!!!!!!!!!! 取得檔案時間到")
			m_onErrorCallback()
		else
			-- 還未完成
			if m_www.isDone == false then
				DoDownloadFile()				
				
			-- 發生錯誤，重新加載
			elseif m_www.error ~= nil then
				ShowLog("!!!!!!!!!!! 取得檔案失敗:" .. m_www.error)				
				m_onErrorCallback()
			
			-- 沒有資料
			elseif m_www.text == nil then
				ShowLog("!!!!!!!!!!! 取得檔案中無資料")
				m_onErrorCallback()
			
			-- 有抓到資料
			else
				ShowLog("取得取得檔案內容\n" .. m_www.text)
				--local fileList = G.Tool.StringHelper.NewSplit(m_www.text, "\n")				
				m_onSuccessCallback(m_www.text)				
			end	
		end
	end)
end