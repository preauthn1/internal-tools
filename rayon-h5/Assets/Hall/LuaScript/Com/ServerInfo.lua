require("Module")
mmodule(...)

local G = require("Global")

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)


local yield_return_HallConnect = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback("HallConnect", to_to_yield, cb)
end
)

-- 服务器IP
local serverIP = ""
-- 服务器端口
local serverPort = 0
-- 服务器模式 0 直连 1 转发
local serverMode = 0
local co_connectHall    -- 连接大厅服务器协程
local str_StopCo = "StopCoroutine"  -- 退出协程关键字符串
local serverArea = "" -- 伺服器位置

-- 清理服务器地址
local function clearServerAddr()
    serverIP = ""
    serverPort = 0
	serverArea = ""
end

-- 发送Token给大廳服务器
function SendTokenMsg()
	local TokenData = G.User.GetFullToken()
	G.NetQ.SendHallMsg(G.Proto.BClass.Common, G.Proto.Common.ConnTokenID, { token = [REDACTED] })
end

-- 发送本地IP给服务器 0 大厅服务器　1游戏服务器
function SendLocalIPMsg(serverType)
    local localIP = G.GetLocalIP()
    if (localIP ~= nil and localIP ~= "" and serverMode == 1) then
        if (serverType == 0) then
            print("发送IP地址" .. localIP .. "给大厅服务器")
            G.NetQ.SendHallMsg(G.Proto.BClass.Common, G.Proto.Common.IPReportID, { IPAddress =G.Tool.Debug.GetPublicIP() })
        else
            print("发送IP地址" .. G.GetLocalIP() .. "给游戏服务器")
            G.NetQ.SendGameMsg(G.Proto.BClass.Common, G.Proto.Common.IPReportID, { IPAddress = G.GetLocalIP() })
        end
    end
end

-- 设置服务器信息
function SetServerInfo(ip, port, mode, area)
    serverIP = ip
    serverPort = port
    serverMode = G.Tool.Math.ParseInt(mode)
	serverArea = area
	
	print("###############SetServerInfo", ip, port, mode, area)
	
	-- IOS美國IP直接設成送審階段
	if G.Env.IsIOS() == true then
		serverArea = string.upper(serverArea)
		if serverArea == "US" then
			G.GameConfig.Submission = true
		end
	end
end

-- 获取大厅服务器信息
function GetHallServerAddr()
    if serverIP == "" or serverPort == 0 then
        return nil, nil
    end

    return serverIP, serverPort
end

local AddrIP
local AddrPort
local AddrCallback

-- 获取游戏服务器地址
function GetGameServerAddrCallback(ip, port, callback)
    -- 直连模式 - 用server給的IP
    if serverMode == 0 then
		print("直连模式")
		-- 无IP资料,直接使用大厅IP
		if ip == "" then
			callback(0, serverIP, port)		
		else
			callback(0, ip, port)
		end
		return
    end
	
	print("转发模式")
	AddrIP = ip
	AddrPort = port
	AddrCallback = callback	
	
    -- 转发模式 - 用hall IP
	local use = G.Tool.Aliyun.GetNowUse()	
	if use == 1 then
		local res, AliyunIp, AliyunPort = G.Tool.Aliyun.GetIP(port)
		if res == 0 then
			callback(res, AliyunIp, AliyunPort)
		else
			G.Tool.Aliyun.DownloadAliyunFile(ServerDownloadAliyunSuccess, ServerDownloadAliyunError)		
		end
	else
		callback(0, serverIP, port)
	end
end

function ServerDownloadAliyunError(errorCode)
	AddrCallback(1, serverIP, port)
end

function ServerDownloadAliyunSuccess()
	local res, AliyunIp, AliyunPort = G.Tool.Aliyun.GetIP(port)
	AddrCallback(res, AliyunIp, AliyunPort)
end

-- 获取游戏服务器地址
--[[
function GetGameServerAddr(ip, port)
    -- 直连模式 - 用server給的IP
    if serverMode == 0 then
        returnip, port
    end

    -- 转发模式 - 用hall IP
	return serverIP, port
end
]]--

-- 获取游戏服务器所在地區
function GetGameServerArea()
    return serverArea
end

-- TF連線錯誤
local IsTfDataError = false
function GetIsTfDataError()
	return IsTfDataError
end
-- TP連限逾時
local IsTfError = 0
function GetIsTfError()
	return IsTfError
end
-- TP連限逾時
local IsTfTimeOut = 0
function GetIsTfTimeOut()
	return IsTfTimeOut
end
-- 無法連線到大廳
local IsConnectError = 0
function GetIsConnectError()
	return IsConnectError
end
-- 無法連線到大廳(直接取得)
local IsDirConnectError = 0
function GetIsDirConnectError()
	return IsDirConnectError
end
-- 抓取TF txt資料錯誤
local IsTfFileError = 0
function GetIsTfFileError()
	return IsTfFileError
end

function ClearFlag()
	print("************************* Clear server info flag")
	IsTfError = 0
	IsTfTimeOut = 0
	IsConnectError = 0
	IsDirConnectError = 0
	IsTfDataError = 0
	IsTfFileError = 0
	
	idx = 0
end

-- 取得錯誤訊息
function GetErrorCode()
	local errorValue = ""
	-- TF連線錯誤
	errorValue = errorValue .. tostring(GetIsTfError())
	
	-- TP連限逾時
	errorValue = errorValue .. tostring(GetIsTfTimeOut())
	
	-- 無法連線到大廳
	errorValue = errorValue .. tostring(GetIsConnectError())
	
	-- 無法連線到大廳(直接取得)
	errorValue = errorValue .. tostring(GetIsDirConnectError())
	
	-- TF資料錯誤
	errorValue = errorValue .. tostring(GetIsTfDataError())
	
	-- 抓取TF txt資料錯誤
	errorValue = errorValue .. tostring(GetIsTfFileError())
	
	return errorValue
end

-- 加载服务器地址
function LoadServerAddr(onSuccessCallback, onErrorCallback)	
    if (co_connectHall ~= nil) then
		coroutine.resume(co_connectHall, str_StopCo)
	end

    -- 成功回调
    local onSuccess = function()
        if onSuccessCallback ~= nil then
            onSuccessCallback()
        end
    end

    -- 失败回调
    local onError = function()
        if onErrorCallback ~= nil then
            onErrorCallback()
        end
    end

    -- 清理地址信息
    clearServerAddr()

    -- 本地环境和测试环境,使用默认的IP和端口,直连模式
    if G.GameConfig.IsLocalCtx() or G.GameConfig.IsTestCtx() or G.GameConfig.IsOnlineTextCtx()then
	
        StartUp.setTipMessage("正在连接服务器...", 0.3);

        local localaddr = G.GameConfig.DefaultServerIP .. ":" .. G.GameConfig.DefaultServerPort

        co_connectHall = coroutine.create( function()
            print("尝试连接大厅服务器." .. localaddr)
            local co_val = yield_return_HallConnect(CS.NetworkManager.Instance:SendHallConnect(G.GameConfig.DefaultServerIP, G.GameConfig.DefaultServerPort))
            if (co_val == str_StopCo) then return end
            if CS.NetworkManager.Instance.HallSocket.MyConnectState == CS.ConnectState.Connect then
                print("连接到大厅服务器成功." .. localaddr)
                SetServerInfo(G.GameConfig.DefaultServerIP, G.GameConfig.DefaultServerPort, 0, "nil")
				-- 上報token
				SendTokenMsg()
                -- 发送本地IP给大厅服务器
                SendLocalIPMsg(0)
                onSuccess()
            else
                print("无法连接到大厅服务器." .. localaddr)
                onError()
            end
        end )
		coroutine.resume(co_connectHall)
	-- 非直連模式
    else
	
        StartUp.setTipMessage("正在取得服务器资料...", 0.2);

		--Edit by Cobain
		--改輪尋
		InitSuccessCallback(onSuccessCallback, onErrorCallback)
		DoLoadServerAddr(G.GameConfig.ServerListDomain)
	end
end

-- ****************************************************************************
-- add by Rayon
-- 重試次數
local m_TryAccCount = 12
-- 重試間隔
local m_TryAccTime = 5
-- 判斷是否重間隔
local m_InvokeTime = 0.1
-- 累計重試次數
local m_AccCount = 0
-- 累計時間
local m_AccTime = 0
-- 抓取完成callback
local m_onSuccessCallback
-- 抓取失敗callback
local m_onErrorCallback
-- 下載WWW
local m_servelFile_www
-- add by Cobain 是否已經試過ServerListFile
local m_servelFile_www_flag = false

local idx = 0

function ShowLog(str)
	print(str)
end

function InitSuccessCallback(onSuccessCallback, onErrorCallback)
	m_onSuccessCallback = onSuccessCallback
	m_onErrorCallback = onErrorCallback
end

-- 成功回调
local onSuccessCallback = function()
	print("onSuccessCallback : ", (m_onSuccessCallback ~= nil))
	if m_onSuccessCallback ~= nil then
		m_onSuccessCallback()
	end
end

-- 失败回调
local onErrorCallback = function()
	if m_onErrorCallback ~= nil then
		m_onErrorCallback()
	end
end

function InitAndBeginLoadingWWW()
	m_AccCount = 0
	BeginLoadingWWW(0)
end

-- WWW加載和開始計時
function BeginLoadingWWW(AddCount)
	m_AccTime = 0
	m_AccCount = m_AccCount + AddCount

	-- 已達重新抓取次數上限
	if m_AccCount >= m_TryAccCount then
		IsTfFileError = IsTfFileError + 1
		onErrorCallback()
		
	-- 未達上限，重新抓取
	else
		-- 支援填完整路徑
		serverFile_url = G.GameConfig.ServerListFile
		f, p = string.find(serverFile_url, ".txt")
		if f == nil then
			serverFile_url = serverFile_url .. "/serverlistnew.txt"
		end
		serverFile_url = serverFile_url .. "?abc=" .. tostring(math.random(1, 100000))
		ShowLog("取得连接服务器 " .. serverFile_url ..  " 次數 " .. tostring(m_AccCount + 1))

		m_servelFile_www = CS.UnityEngine.WWW(serverFile_url)

		DoLoadServerFile()
	end
end

function DoLoadServerFile()

	print("DoLoadServerFile")

	CS.TimeManager.Instance:StopAllDelayInvoke()
				
	CS.TimeManager.Instance:DelayInvoke(m_InvokeTime, function( )
		
		m_AccTime = m_AccTime + m_InvokeTime
		
		-- 超過重新抓取間隔，重新抓取
		if m_AccTime > m_TryAccTime then
			ShowLog("!!!!!!!!!!! 取得連接服務器檔案時間到")
			BeginLoadingWWW(1)
		else
			-- 還未完成
			if m_servelFile_www.isDone == false then
				DoLoadServerFile()				
				
			-- 發生錯誤，重新加載
			elseif m_servelFile_www.error ~= nil then
				ShowLog("!!!!!!!!!!! 取得連接服務器檔案失敗:" .. m_servelFile_www.error)				
				BeginLoadingWWW(1)
			
			-- 沒有資料
			elseif m_servelFile_www.text == nil then
				ShowLog("!!!!!!!!!!! 連接服務器檔案中無資料")
				DoLoadServerAddr(nil)
			
			-- 有抓到資料
			else
				ShowLog("取得伺服器列表檔案內容\n" .. m_servelFile_www.text)
				local ServerListDomain = G.Tool.StringHelper.NewSplit(m_servelFile_www.text, "\n")				
				DoLoadServerAddr(ServerListDomain)				
			end	
		end
	end)
end

-- ****************************************************************************
function GetTransferData(localUserData, mode)
	local data = ""
	-- 一般未加密機制
	if G.GameConfig.isUseNewTransfer() == false then
		-- 連美服
		if G.GameConfig.GetServerOption() == 3 then
			data = "/Serverlist/AppleReview?" ..
					"HDcode=" .. G.Env.GetHDCode() ..
					"&areaid=" .. tostring(G.GameConfig.GamePlatformID)
		else
			data = "/serverlist/transfer?" ..
					"gui=" .. G.Env.GetHDCode() ..
					"&channel=" .. tostring(G.GameConfig.GamePlatformID) ..
					"&simulation=" .. "0"
		end
		
		if (localUserData ~= nil) then
			data = data.. "&levelid=" .. localUserData.UserLevel ..
							"&userid=" .. localUserData.UserID ..
							"&key=" .. localUserData.LevelKey
		else				
			data = data .. "&levelid=" .. "0" ..
							"&userid=" .. "0" ..
							"&key=" .. "0"
		end
		
	-- 使用加密機制
	else
		logindata = {}
		logindata.GUI = G.Env.GetHDCode()
		logindata.Channel = tostring(G.GameConfig.GamePlatformID)
		logindata.Simulation = "0"
		logindata.time = "0"
		logindata.Action = tostring(mode)		
		
		if (localUserData ~= nil) then
			logindata.UserID = localUserData.UserID
			logindata.VIP = localUserData.UserLevel
		else
			logindata.UserID = "0"
			logindata.VIP = "0"
		end
		
		data = "/ServerNode/"
		if G.GameConfig.GetServerOption() == 3 then
			data = data .. "AppleReview"
		else
			data = data .. "Transfer"
		end
		data = data .. "?data=" .. G.Tool.Aes.Encode(logindata)
	end
	
	return data
end

function ToData(receiveData)
	local jasnData = G.Tool.Aes.Decode(receiveData)
	local data = ""
	if jasnData ~= nil then
		if jasnData["result"] == 1 then
			dataObject = jasnData["data"]
			data = dataObject["Addr"] .. ":" .. dataObject["TransferPort"] .. "|" ..
					dataObject["GiveUserIP"] .. "|" ..
					dataObject["From"] .. "|" ..
					dataObject["UserIP"]
			
			CS.PrefsUtility.SetString("Token", receiveData)
		end	
	end
	return data
end

local tempServerListDomain
function DoLoadServerAddr(ServerListDomain)
	
    StartUp.setTipMessage("正在连接服务器...", 0.2);
	
	if ServerListDomain == nil then
		ShowLog("******* ServerListDomain 資料是空的")
		onErrorCallback()
		return
	end	
	
	tempServerListDomain = ServerListDomain
	
	if G.Tool.Aliyun.IsAliyunSDK() == true then
		-- 已經初始化成功，可以直接取IP
		if G.Tool.Aliyun.IsInitSuccess() == true then
			-- 取得遊戲盾IP資料
			local res, ip, port = G.Tool.Aliyun.GetIP(G.Tool.Aliyun.GetPort())
			-- 取IP成功
			if res == 0 then
				ConnectHallServer(ip, port)
			-- 取IP失敗，直接抓檔案
			elseif res == 1 then
				G.Tool.Aliyun.DownloadAliyunFile(DownloadAliyunSuccess, DownloadAliyunError)
			end
		-- 初始化失敗，直接抓檔案
		else
			G.Tool.Aliyun.DownloadAliyunFile(DownloadAliyunSuccess, DownloadAliyunError)
		end
	else
		DoTFServer(tempServerListDomain)
	end
end

function DownloadAliyunSuccess()
	local use = G.Tool.Aliyun.GetDownloadUse()
	-- 遊戲盾
	if use == 1 then
		-- 已經初始化成功，可以直接取IP
		if G.Tool.Aliyun.IsInitSuccess() == true then
			-- 取得遊戲盾IP資料
			local res, ip, port = G.Tool.Aliyun.GetIP(G.Tool.Aliyun.GetPort())
			-- 取IP成功
			if res == 0 then
				ConnectHallServer(ip, port)
			-- 取IP失敗
			elseif res == 1 then
				m_onErrorCallback()
			end
		-- 初始化失敗
		else
			m_onErrorCallback()
		end
	-- TF Server
	elseif use == 2 then
		DoTFServer(tempServerListDomain)
	end
end

function DownloadAliyunError(errorCode)
	print("DownloadAliyunError " .. tostring(errorCode))
	-- 沒有新版本
	if errorCode == 0 then
		G.Util.coroutine_call( function()
			local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "", "游戏资源已更新，请重新启动游戏", "确定");
			if boxResult == CS.MessageBoxResult.OK then
				CS.UnityEngine.Application.Quit()
			end
		end )()
	-- key不一樣,需要重啟
	else
		if errorCode == 1 then
			IsTfError = 1
		elseif errorCode == 2 then
			IsTfTimeOut = 1
		elseif errorCode == 3 then
			IsConnectError = 1
		elseif errorCode == 4 then
			IsDirConnectError = 1
		elseif errorCode == 5 then
			IsTfDataError = 1
		end

		m_onErrorCallback()	
	end
end

function ConnectHallServer(ip, port)
	G.Tool.Aliyun.SetNowUse(1)

	co_connectHall = coroutine.create( function()
		
		print("正在连接服务器...", 0.3)		
		
		local co_val = yield_return_HallConnect(CS.NetworkManager.Instance:SendHallConnect(ip, port))
		if (co_val == str_StopCo) then 
			return
		end
		
		if CS.NetworkManager.Instance.HallSocket.MyConnectState == CS.ConnectState.Connect then
			print("连接到大厅服务器成功," .. ip .. ":" .. port)
		
			--[[
			userIP = G.Tool.Debug.CheckIP()
			if userIP == "" then
				userIP = "0.0.0.0"
			end
			]]--
			
			-- userIP = tostring(math.random(50,250)) .. "." ..
			-- 	tostring(math.random(50,250)) .. "." ..
			-- 	tostring(math.random(50,250)) .. "." ..
			-- 	tostring(math.random(50,250)) 
			--G.GameConfig.UserIP = userIP
			--G.SetLocalIP(userIP)
			G.GameConfig.UserIP=G.Tool.Debug.GetPublicIP()
			G.SetLocalIP(G.Tool.Debug.GetPublicIP())
			mode = 1
			area = "CN"
			
			-- 清除登入錯誤訊息
			CS.WindowManager.Instance:CloseWindow("UILoginMessageBox", false);
			-- 设置服务器信息
			SetServerInfo(ip, port, mode, area)
			-- 上報token
			SendTokenMsg()
			-- 发送本地IP给大厅服务器
			SendLocalIPMsg(0)
			-- 连接服务器成功回调
			onSuccessCallback()
		else													
			IsConnectError = IsConnectError + 1
			
			onErrorCallback()
			print("连接到大厅服务器失败," .. ip .. ":" .. tostring(port))
		end
	end )
    coroutine.resume(co_connectHall)
end

function DoTFServer(ServerListDomain)
	G.Tool.Aliyun.SetNowUse(2)

	co_connectHall = coroutine.create( function()
					
		
		print("正在连接服务器...", 0.3)
		local localUserData = G.User.GetLocalUserInfo()
		
		for _, url in pairs(ServerListDomain) do

			idx = idx + 1

			-- url decode
			url = "http://" .. G.Tool.Base64.decodeBase64(url)
			
			-- step1一般流程
			-- step2直接取IP
			-- 紀錄step1是否連線失敗，連線失敗時step2直接再取一次IP
			local ConnectError = false
			for step = 1, 2 do			
			
				local urlPath = url
				
				local IsDo = true
				if step == 2 then
					IsDo = ConnectError
					StartUp.setTipMessage("正在连接服务器("..tostring(idx).."-"..tostring(step)..")...", 0.2)
				else
					StartUp.setTipMessage("正在连接服务器("..tostring(idx)..")...", 0.2)				
				end
				
				if IsDo == true then
					local transferData = GetTransferData(localUserData, step)
					--print("transferData\n" .. transferData)
					urlPath = urlPath .. transferData .. "&r=" .. tostring(math.random(1, 100000))
					
					local www1 = CS.UnityEngine.Networking.UnityWebRequest.Get(urlPath)
					CS.TimeManager.Instance:YieldAndCallback(www1:SendWebRequest(), 
					function() 
						print(www1.downloadHandler.text)
						
						if www1.error == nil then
							-- 获得文本数据
							local data = www1.downloadHandler.text
							
							-- 加密資料，需解密產生同樣結構
							if G.GameConfig.isUseNewTransfer() == true then
								data = ToData(data)								
							end
							
							if data == "" then
								if step == 1 then
									IsTfDataError = IsTfDataError + 1
								end
								print("加载serverlist成功,解析数据錯誤")
							else
								print("加载serverlist成功,获取到数据 " .. data)
								
								-- 解析数据
								local arr1 = G.Tool.StringHelper.Split(data, "|")

								if #arr1 < 3 then
									print("数据结构错误," .. data .. "," .. urlPath)
								
								-- 檢查資料格式是否正確
								elseif CheckData(data) then
									local ipstr = arr1[1]
									local mode = arr1[2]
									local area = "nil"
									local userIP = "nil"
									if #arr1 >= 3 then
										area = arr1[3]
									end

									--Add by Cobain
									if #arr1 >= 4 then
										userIP = arr1[4]
									end

									local arr2 = G.Tool.StringHelper.Split(ipstr, ",")
									for _, addr in pairs(arr2) do
										if addr ~= nil and addr ~= "" then
											local arr3 = G.Tool.StringHelper.Split(addr, ":")
											if #arr3 == 2 and arr3[1] ~= "" and arr3[2] ~= "" then

												-- ip
												local ip = tostring(arr3[1])
												-- 端口
												local port = tonumber(arr3[2])
												
												-- Edit by Cobain x-forward取得UserIP
												G.GameConfig.UserIP = userIP
												G.SetLocalIP(userIP)
												
												-- 先連到kuoyu測試
												if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.WebGLPlayer then
													ip = "0.0.0.0"
													port = 5000
												end
												
												CS.TimeManager.Instance:YieldAndCallback(CS.NetworkManager.Instance:SendHallConnect(ip, port), function() 												
													-- local co_val = yield_return_HallConnect(CS.NetworkManager.Instance:SendHallConnect(ip, port))
													-- if (co_val == str_StopCo) then 
														-- return
													-- end
													
													if CS.NetworkManager.Instance.HallSocket.MyConnectState == CS.ConnectState.Connect then
														print("连接到大厅服务器成功," .. addr)
														-- 清除登入錯誤訊息
														CS.WindowManager.Instance:CloseWindow("UILoginMessageBox", false);
														print("连接到大厅服务器成功 A")
														-- 设置服务器信息
														SetServerInfo(ip, port, mode,area)
														print("连接到大厅服务器成功 B")
														-- 上報token
														SendTokenMsg()
														print("连接到大厅服务器成功 C")
														-- 发送本地IP给大厅服务器
														SendLocalIPMsg(0)
														print("连接到大厅服务器成功 D - SendLocalIPMsg")
														-- 连接服务器成功回调
														onSuccessCallback()
														print("连接到大厅服务器成功 E - 连接服务器成功回调")
														
														return
														
													else													
														if step == 1 then
															ConnectError = true
															IsConnectError = IsConnectError + 1
														else
															IsDirConnectError = IsDirConnectError + 1
														end
														print("连接到大厅服务器失败," .. addr)
													end
												end)
												
											end
										end
									end
								end
							end
						end		
					end)
					
					--[[
					local www = CS.UnityEngine.WWW(urlPath)
					print("开始加载serverlist " .. tostring(idx) .. " " .. urlPath)
					local co_val = yield_return_HallConnect(CS.TimeManager.Instance:TimeOutReturn( function() return www.isDone end, 5))
					
					if (co_val == str_StopCo) then
						return
					end
					
					if (www.isDone) then
						-- 加载url成功
						if www.error == nil then
							-- 获得文本数据
							local data = www.text
							
							-- 加密資料，需解密產生同樣結構
							if G.GameConfig.isUseNewTransfer() == true then
								data = ToData(data)
							end
							
							if data == "" then
								if step == 1 then
									IsTfDataError = IsTfDataError + 1
								end
								print("加载serverlist成功,解析数据錯誤")
							else
								print("加载serverlist成功,获取到数据 " .. data)
								
								-- 解析数据
								local arr1 = G.Tool.StringHelper.Split(data, "|")

								if #arr1 < 3 then
									print("数据结构错误," .. data .. "," .. urlPath)
								
								-- 檢查資料格式是否正確
								elseif CheckData(data) then
									local ipstr = arr1[1]
									local mode = arr1[2]
									local area = "nil"
									local userIP = "nil"
									if #arr1 >= 3 then
										area = arr1[3]
									end

									--Add by Cobain
									if #arr1 >= 4 then
										userIP = arr1[4]
									end

									local arr2 = G.Tool.StringHelper.Split(ipstr, ",")
									for _, addr in pairs(arr2) do
										if addr ~= nil and addr ~= "" then
											local arr3 = G.Tool.StringHelper.Split(addr, ":")
											if #arr3 == 2 and arr3[1] ~= "" and arr3[2] ~= "" then

												-- ip
												local ip = tostring(arr3[1])
												-- 端口
												local port = tonumber(arr3[2])
												
												-- Edit by Cobain x-forward取得UserIP
												G.GameConfig.UserIP = userIP
												G.SetLocalIP(userIP)
												
												local co_val = yield_return_HallConnect(CS.NetworkManager.Instance:SendHallConnect(ip, port))
												if (co_val == str_StopCo) then 
													return
												end
												
												if CS.NetworkManager.Instance.HallSocket.MyConnectState == CS.ConnectState.Connect then
													print("连接到大厅服务器成功," .. addr)

													-- 清除登入錯誤訊息
													CS.WindowManager.Instance:CloseWindow("UILoginMessageBox", false);
													-- 设置服务器信息
													SetServerInfo(ip, port, mode,area)
													-- 上報token
													SendTokenMsg()
													-- 发送本地IP给大厅服务器
													SendLocalIPMsg(0)
													-- 连接服务器成功回调
													onSuccessCallback()
													
													return
													
												else													
													if step == 1 then
														ConnectError = true
														IsConnectError = IsConnectError + 1
													else
														IsDirConnectError = IsDirConnectError + 1
													end
													print("连接到大厅服务器失败," .. addr)
												end
											end
										end
									end
								end
							end
						else
							if step == 1 then
								IsTfError = IsTfError + 1
							end
							print("加载serverlist失败," .. www.error)
						end
					else
						if step == 1 then
							IsTfTimeOut = IsTfTimeOut + 1
						end
						print("加载serverlist超时," .. url)
					end]]
				end
			end
		end
		--[[
		print("加载所有serverlist失败...")
		-- 连接服务器失败回调
		-- Edit by Cobain
		-- 如果ServerListDomain都失效,訪問ServerListFile
		if m_servelFile_www_flag == false then
			m_servelFile_www_flag = true
			BeginLoadingWWW(0)
		else 
			onErrorCallback()
		end
		]]
	end )
    coroutine.resume(co_connectHall)
end

function CheckData(data)
	local arr1 = G.Tool.StringHelper.Split(data, "|")
	if #arr1 < 3 then
		return false
	else
		-- 檢查IP格式(包含port)
		local ipstr = arr1[1]
		if GetIPType(ipstr, true) ~= 1 then
			print("ip error")
			return false
		end

		-- 檢查模式，必須是數字
		local mode = arr1[2]
		local n = tonumber(mode);
		if n ==  nil then
			print("mode error")
			return false
		end

		-- 檢查自己的IP格式(不包含port)
		if #arr1 >= 4 then
			local userIP = arr1[4] 
			if GetIPType(userIP, false) ~= 1 then
				print("self ip error")
				return false
			end
		end
		--print("########## 資料檢驗通過")
		return true
	end 
end

--[[
    [0] = "Error",
    [1] = "IPv4",
    [2] = "string",
]]--
function GetIPType(ip, hasPort)
	-- must pass in a string value
	if ip == nil or type(ip) ~= "string" then
		return 0
	end

	-- check for format 0.0.0.0:10000 for ipv4
	local chunks
	if hasPort == false then
		chunks = {ip:match("(%d+)%.(%d+)%.(%d+)%.(%d+)")}
	else
		chunks = {ip:match("(%d+)%.(%d+)%.(%d+)%.(%d+)%:(%d+)")}
	end

	if (#chunks == 4 or #chunks == 5) then
	-- IP
	for i = 1, 4 do
		if (tonumber(chunks[i]) < 0 or tonumber(chunks[i]) > 255) then
			return 0
		end
	end

	-- port
	if #chunks == 5 then
		if (tonumber(chunks[5]) < 0 or tonumber(chunks[5]) > 65535) then
			return 0
		end  
	end

		return 1
	else
		return 0
	end

	-- if we get here, assume we've been given a random string
	return 2
end

