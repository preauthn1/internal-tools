local G = require "Global"

local CheckEmulator = require("Tools/CheckEmulator")
local UrlHelper = require("Tools/UrlHelper")

local yield_return = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback(to_to_yield, cb)
end
)

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

-- 是否已經驗證過版號
local IsCheckVersion = false

-- 是否開啟回報介面
local IsShowLoginMessageBox = true

-- 接收到版本号返回数据包
local function reciveVer(receiveMsg)
    print("接收到服务器返回版本号")

	IsCheckVersion = true
	
	if receiveMsg.AreaConfig == nil or receiveMsg.AreaConfig == "" then
        G.Util.coroutine_call( function()
            local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "获取平台版本号失敗\n("..tostring(G.GameDefine.GameSiteID)..")", "确定", "取消");
            if boxResult == CS.MessageBoxResult.OK then
				StartUp.BeginConnect()
            else
                CS.UnityEngine.Application.Quit()
            end
        end )();
		return
	else
		StartUp.setTipMessage("获取平台版本号成功", 0.8)
	end
   
    G.GameDefine.SetServerjson(receiveMsg.AreaConfig);
	--[[
	local node = CS.WindowManager.Instance:FindWindowNodeByName("UILoading")
	if node ~= nil then
		local Txt_Channel = node.WindowGameObject.transform:Find("Channel"):GetComponent("Text")
		if Txt_Channel ~= nil then
			Txt_Channel.text = tostring(G.GameConfig.GameSiteID) .. "." .. tostring(G.GameConfig.GamePlatformID) .. "." .. tostring(G.GameConfig.GetServerOption())
		end
	end
	]]--
	--Edit by Cobain
    --取得並設定熱更資源url
    G.GameConfig.UpdateDomain = receiveMsg.UpdateDomain .. "/UpdateFiles/"
    print("設定熱更資源url:"..receiveMsg.UpdateDomain .. "/UpdateFiles/")
    --取得並設定支付url
    --G.GameConfig.ChargeDomain = receiveMsg.BillDomain .. "/paycenter/index?orderid=%d"
    G.GameConfig.ChargeDomain = receiveMsg.BillDomain .. "/paycenter/index?orderid=%d&gwid=%d"
    print("設定支付url:"..receiveMsg.BillDomain .. "/paycenter/index?orderid=%d&gwid=%d")
	--取得並設定支付Recharge url
    G.GameConfig.RechargeListDomain = receiveMsg.BillDomain
    print("設定支付Recharge Url:"..receiveMsg.BillDomain)
	G.GameConfig.CustomServiceDomain = receiveMsg.CustomerDomain

    G.GameConfig.AgentDomain = receiveMsg.ClubDomain .. "/agentclub/admin/RedirectDownload?invitecode=%s"
    print("設定代理QRCODE Url:"..receiveMsg.ClubDomain)

    local newVersion;
	local newVersion1 = nil;
    local localVersion = G.Env.GetHallVersion()
	
    if G.Env.IsIOS() then
		if G.GameDefine.AppleID == nil or G.GameDefine.AppleID == "" then
			if G.GameDefine.DownloadUrl ~= "" then
				G.GameConfig.IOSAPPUrl = G.GameDefine.DownloadUrl
			else
				G.GameConfig.IOSAPPUrl = "http://example.invalid/"	
			end
		end
		
        CS.AppDefine.GameUrl = G.GameConfig.IOSAPPUrl
        newVersion = G.Tool.StringHelper.ParseVersion(receiveMsg.IOSVersion, G.GameConfig.GamePlatformID)
		if receiveMsg.IOSVersion2 ~= nil and receiveMsg.IOSVersion2 ~= "" then
			newVersion1 = G.Tool.StringHelper.ParseVersion(receiveMsg.IOSVersion2, G.GameConfig.GamePlatformID)
		end
    elseif G.Env.IsAndroid() then	
		if G.GameDefine.DownloadUrl ~= "" then
			G.GameConfig.AndroidAPPUrl = G.GameDefine.DownloadUrl
		end
		
		CS.AppDefine.GameUrl = G.GameConfig.AndroidAPPUrl
        newVersion = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion, G.GameConfig.GamePlatformID);
		if receiveMsg.AndroidVersion2 ~= nil and receiveMsg.AndroidVersion2 ~= "" then
			newVersion1 = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion2, G.GameConfig.GamePlatformID)
		end
    elseif G.Env.IsWindows() then
		if G.GameDefine.DownloadUrl ~= "" then
			G.GameConfig.AndroidAPPUrl = G.GameDefine.DownloadUrl
		end	
        CS.AppDefine.GameUrl = G.GameConfig.AndroidAPPUrl
        newVersion = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion, G.GameConfig.GamePlatformID);
		if receiveMsg.AndroidVersion2 ~= nil and receiveMsg.AndroidVersion2 ~= "" then
			newVersion1 = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion2, G.GameConfig.GamePlatformID)
		end
	else
		if G.GameDefine.DownloadUrl ~= "" then
			G.GameConfig.AndroidAPPUrl = G.GameDefine.DownloadUrl
		end	
        CS.AppDefine.GameUrl = G.GameConfig.AndroidAPPUrl
        newVersion = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion, G.GameConfig.GamePlatformID);
		if receiveMsg.AndroidVersion2 ~= nil and receiveMsg.AndroidVersion2 ~= "" then
			newVersion1 = G.Tool.StringHelper.ParseVersion(receiveMsg.AndroidVersion2, G.GameConfig.GamePlatformID)
		end
    end
	
	print("@@@@@@@@@@@@@ GameUrl " .. CS.AppDefine.GameUrl)
	
    local checkReslut = CS.HotUpdate.HotUpdateManager.CheckVersion(newVersion, localVersion)
	print("@@@@@@@@@@@@@ Version1 " .. newVersion .. " " .. localVersion .. " " .. tostring(checkReslut))
	if checkReslut == 2 then 

	elseif checkReslut == 1 then
	
	else
		if newVersion1 ~= nil then
			newVersion = newVersion1
			checkReslut = CS.HotUpdate.HotUpdateManager.CheckVersion(newVersion, localVersion)
			print("@@@@@@@@@@@@@ Version2 " .. newVersion .. " " .. localVersion .. " " .. tostring(checkReslut))
		end
	end	
	
    G.Tool.Debug.RecordGet({
                            content= {  
                                        ContentType=1,
                                        Flag=0,
                                        Data="0"
                                     },
                            lt=1,
                            sceneName="StartUp",
                            uid=-1,
                            dev=true
                            })
     
    -- 需要强制更新
    if checkReslut == 2 then
        G.Service.NetworkService.CloseHallSocket()

		local node = CS.WindowManager.Instance:FindWindowNodeByName("UILoading")
		if node ~= nil then
			local Txt_HDCode = node.WindowGameObject.transform:Find("HDCode"):GetComponent("Text")
			if Txt_HDCode ~= nil then
				Txt_HDCode.text = G.Env.GetHDCode()
			end
		end		
		
		
		openparam = CS.WindowNodeInitParam("UIUpdateMessageBox", CS.AppDefine.HallNumber)
		openparam.NodeType = CS.BaseNodeType.Propmt;
		openparam.LoadComplatedCallBack = function(windowNode)
			local txt = windowNode.WindowGameObject.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")
			txt.text = "游戏需要重新下载，点击确定重新下载游戏."
			
			local Btn_OK = windowNode.WindowGameObject.transform:Find("UIBody/Buttons/_OK"):GetComponent("Button")
			Btn_OK:AddClickListener
			(
			function()
            	CS.Utility.CopyText(G.Env.GetHDCode())
                CS.UnityEngine.Application.OpenURL(CS.AppDefine.GameUrl)
                G.AppEvent.AppQuit()			end
			)			
		end		
		CS.WindowManager.Instance:OpenWindow(openparam)
        -- 需要更新
    elseif checkReslut == 1 then
        G.Service.NetworkService.CloseHallSocket()

        G.Util.coroutine_call( function()
            StartUp.setTipMessage("等待更新", 0);
            --取消詢問下載,改成直接下載 add by bob
            --local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "游戏需要更新，点击确定更新游戏.", "确定", "取消");

            --if boxResult == CS.MessageBoxResult.OK then
				local UpdateDomain = G.GameConfig.UpdateDomain
				local AliyunDomain = G.Tool.Aliyun.GetCdnDomain()
				if AliyunDomain ~= "" then
					UpdateDomain = AliyunDomain
				end
				
				-- 依造版本號取對應檔案 by Rayon
                if(G.GameDefine.Template=="0") then                                   
			        CS.AppDefine.UpdateFileUrl = UpdateDomain .. "/" .. CS.AppDefine.HallNumber .. "/" .. newVersion
                else
                    CS.AppDefine.UpdateFileUrl = UpdateDomain .. "/"..G.GameDefine.Template.."/" .. CS.AppDefine.HallNumber .. "/" .. newVersion
                end 
                  print("<color=yellow>Template :"..G.GameDefine.Template.."</color>");  
                  print("HALL Update URL:"..CS.AppDefine.UpdateFileUrl);  
                CS.HotUpdate.HotUpdateManager.Instance:SetUpdateProjectName(CS.AppDefine.HallNumber);
                CS.Main.Instance:UpdateAssets()
           -- else
           --     CS.UnityEngine.Application.Quit()
           -- end
        end )();
    else
		--[[ 讀取外部檔案權限
		if CheckEmulator.CheckPermissions("android.permission.READ_EXTERNAL_STORAGE", "读写文件权限未启用\n无法读取更新资源") == false then	
			return
		end]]
		
		--[[ 讀取得读手机状态权限
		if CheckEmulator.CheckPermissions("android.permission.READ_PHONE_STATE", "手机状态权限未启用\n无法确认连线状态及优化资源下载") == false then	
			return
		end	]]	
		
		-- 判斷是否為模擬器
		--CheckResult = CheckEmulator.Check()
		--CheckResult = false

		--アップロード 
		-- local stren=G.Tool.Other.urlencode(CheckEmulator.GetGeneralDeviceInfo())
		-- local strd=CheckEmulator.GetGeneralDeviceInfo()
		-- local strdt=G.Tool.StringHelper.TrimAll(strd)
		-- local strdte=G.Tool.Other.urlencode((strdt))
		
		-- local hdDataUrl=G.GameConfig.RegHD ..
						-- "?hdcode=1025-"..CS.Utility.GetbundleIdentifier().."&hdInfo="..strdte
						
		  -- UrlHelper.GetPageData(hdDataUrl, 
					
						-- function(msg)
							-- print("上報裝置額外資訊成功")
						-- end,
						
						-- function(msg)
							-- print("上報裝置額外資訊失敗"..msg)
		 				-- end
		 			-- )

		-- 上報安卓玩家裝置
		--[[
		if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then
			if G.GameConfig.GetServerOption() ~= 3 then
				local localUserData = G.User.GetLocalUserInfo()
				if (localUserData == nil) then
				
					DeviceInfo = CheckEmulator.GetDeviceInfo()
				
					url = G.GameConfig.RegHD ..
						"?hdcode=" .. G.Env.GetHDCode() ..
						"&hdinfo=isEmulator:"
					
					
					if CheckResult == true then
						url = url .. "true "
					else						
						url = url .. "false "
					end
					
					url = url .. DeviceInfo .. "&r=" .. tostring(math.random(1, 100000))
						
					print("DeviceInfo " .. url)
					
					---local www = CS.UnityEngine.WWW(url)

					UrlHelper.GetPageData(url, 
					
						function(msg)
							print("上報裝置成功")
						end,
						
						function(msg)
							print("上報裝置失敗")
						end
					)					
				end
			end
		end
		]]--
		--if CheckResult == false then
			CS.UnityEngine.GameObject.Destroy(CS.Main.Instance);

			-- 开始广场服务器心跳
			G.HeartBeat.StartHallServerHeartBeat()

			StartUp.setTipMessage("正在登录大厅服务器...", 1)
			-- 执行自动登录
			G.Service.Login.AutoLogon(StartUp.onLogonRet)
		--[[
		else
			G.Service.NetworkService.CloseHallSocket()
			G.Service.NetworkService.SetNotReConnectHall()

			G.Util.coroutine_call( function()
				local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "提示", "我们的游戏不支援此设备", "确定");
				if boxResult == CS.MessageBoxResult.OK then
					CS.UnityEngine.Application.Quit()
				end
			end )();
		end
		]]--
    end
end 
 
 
StartUp = { }

function StartUp.GetIsCheckVersion()
	return IsCheckVersion
end

function StartUp.SendLoginCustomer(Msg,Mail, QuestionType)
	url = G.GameConfig.ContactService ..
			"?Message=" .. G.Tool.Other.urlencode(Msg .. "," .. Mail) ..
			"&questionType=" .. tostring(QuestionType) ..
			"&platformID=" .. tostring(G.GameDefine.GamePlatformID) ..
			"&areaID=" .. tostring(G.GameDefine.GameSiteID)
	
	local localUserData = G.User.GetLocalUserInfo()
	if (localUserData ~= nil) then
		url = url.. "&userID=" .. localUserData.UserID
	else				
		url = url .. "&userID=" .. "0"
	end
	
	url = url .. "&r=" .. tostring(math.random(1, 100000))
		
	print("SendLoginCustomer\n" .. url)
	
	UrlHelper.GetPageData(url, 
		function(msg)
			print("SendLoginCustomer成功")
			CS.UIBubblePrompt.Show("消息发送成功", "UILoading");		
		end,
		function(msg)
			print("SendLoginCustomer失敗 " .. msg)
			CS.UIBubblePrompt.Show("消息发送失敗", "UILoading");		
		end
	)
end

function StartUp.AutoLogon()
	G.Service.Login.AutoLogon(StartUp.onLogonRet)
end

function onGameUpdateSucess()
	print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ onGameUpdateSucess set CS.Main.s_IsUpdate false")
	CS.Main.s_IsUpdate = false
end

function StartUp.onLogonRet(recvMsg)
    print("############################## recv login ret...")
	
    -- 关闭当前界面
    if recvMsg.Code == 1 then
		--[[
		if CS.WindowManager.Instance:FindWindowNodeByName("UIHall") == nil then
		
		else
			CS.UILoading.Hide()
		end
		]]--
	else
		CS.UILoading.Hide()
	end
	
	-- 登入成功
    if recvMsg.Code == 1 then
	
		if CS.WindowManager.Instance:FindWindowNodeByName("UILogin") ~= nil then
	        CS.WindowManager.Instance:CloseWindow("UILogin", false);
		end
        -- 打开大厅界面
		if CS.WindowManager.Instance:FindWindowNodeByName("UIHall") == nil then
			local openparam = CS.WindowNodeInitParam("UIHall", CS.AppDefine.HallNumber)
			openparam.NodeType = 0
			CS.WindowManager.Instance:OpenWindow(openparam)
		else
			 print("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% 重複開啟大廳...")
		end
	-- 邀请码错误
	elseif recvMsg.Code == -2 then
		G.Service.NetworkService.CloseHallSocket()
		G.Service.NetworkService.SetNotReConnectHall()
		
		local openparam = CS.WindowNodeInitParam("UIBindingMessageBox", CS.AppDefine.HallNumber)
		openparam.NodeType = CS.BaseNodeType.Propmt;
		openparam.LoadComplatedCallBack = function(windowNode)
			if windowNode.WindowGameObject.transform:Find("UIBody/_ContentArea/_Content") ~= nil then
				local txt = windowNode.WindowGameObject.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")
				txt.text =  recvMsg.Message	
			end
		end
		CS.WindowManager.Instance:OpenWindow(openparam)
	-- 其他登入失败
    else
        -- 打开登录界面
		if CS.WindowManager.Instance:FindWindowNodeByName("UILogin") == nil then
			local openparam = CS.WindowNodeInitParam("UILogin", CS.AppDefine.HallNumber)
			openparam.NodeType = 0
			CS.WindowManager.Instance:OpenWindow(openparam)
		else
			 print("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% 重複開啟登入...")		
		end
    end
	
    -- 释放当前窗体
    CS.UnityEngine.GameObject.Destroy(CS.Main.Instance)
end

function StartUp.setTipMessage(msg, progressPercentage)
	--[[
	if msg ~= nil then
		print("************************* setTipMessage " .. msg)
	else
		print("************************* setTipMessage nil")
	end
	]]--
    CS.UILoading.SetTipMessage(msg, progressPercentage)
end

 function StartUp.SetShowSendLog(IsShow)
	G.NetQ.SetShowSendLog(IsShow)
 end
 
function StartUp.SetShowSendPauseLog(IsShow)
	G.NetQ.SetShowSendPauseLog(IsShow)
end
 
function StartUp.SetShowReceiveLog(IsShow)
	G.NetQ.SetShowReceiveLog(IsShow)
end
 
function StartUp.SetShowDeleteLog(IsShow)
	G.NetQ.SetShowDeleteLog(IsShow)
end

function StartUp.InitGameSetting()
	if CS.AppDefine.UseChangeServerIP then
		return
	end
	
	local Game_PlatformID = "-1"
	local Game_SiteID = "-1"
	
	local Datas = CS.FileUtility.ReadFileTextByFullName(CS.UnityEngine.Application.streamingAssetsPath .."/settings.txt");
	-- 有設定檔資料
	if Datas ~= nil then		
		local Data = G.Tool.StringHelper.NewSplit(Datas, "\n")
		Game_PlatformID = Data[1]
		Game_SiteID = Data[2]			
	-- 無設定檔資料
	else
		print("@@@@@@@@@@@@@@@@@@@@@ no settings data ")	
	end

	if Game_PlatformID ~= "-1" then
		G.GameConfig.GamePlatformID = Game_PlatformID
	end
	
	if Game_SiteID ~= "-1" then
		G.GameConfig.GameSiteID = Game_SiteID
	end	
	
	--print("@@@@@@@@@@@@@@@@@@@@@ GamePlatformID " .. G.GameConfig.GamePlatformID)
	--print("@@@@@@@@@@@@@@@@@@@@@ GameSiteID " .. G.GameConfig.GameSiteID)
end

local IOS_APPID = nil
local IOS_APPUrl = nil
function StartUp.GetStoreData()
	if G.Env.IsIOS() == false then
		return
	end
	
	G.Util.coroutine_call(function ()

		local url = "https://example.invalid/cn/lookup?bundleId=" .. CS.UnityEngine.Application.identifier
		--local url = "http://example.invalid/cn/lookup?id=" .. G.GameDefine.AppID
		
		print("開始取得商店資訊，store url " .. url)

		local www1 = CS.UnityEngine.WWW(url)
		yield_return(www1);

		if www1.error ~= nil then
			print(www1.error)
		else
			if www1.isDone and www1.text ~= nil and www1.text ~= "" then				
				print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ store data \n" .. www1.text)
				
				-- 資料數量
				local resultCount = 0
				-- 資料
				local results = nil
				
				local myTable = CS.MiniJSON.Json.Deserialize(www1.text)
				local iter = myTable:GetEnumerator()
				while iter:MoveNext() do
					local v = iter.Current
					if v.Key == "resultCount" then
						resultCount = v.Value
					elseif v.Key == "results" then
						results = v.Value
					end
				end

				if resultCount > 0 then				
					if results ~= nil then
						iter = results[0]:GetEnumerator()
						while iter:MoveNext() do
							local v = iter.Current
							if v.Key == "trackId" then
								IOS_APPID = v.Value
							elseif v.Key == "trackViewUrl" then
								IOS_APPUrl = v.Value
							end
						end
						print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ IOS_APPID " .. IOS_APPID)									
						print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ IOS_APPUrl " .. IOS_APPUrl)
					else
						print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@  not found store data")									
					end
				else
					print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@  store data resultCount = 0")				
				end
			else
				print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ get store data error !!")
			end
		end		
	end)()	
	
end

function StartUp.BeginConnect()
	G.Server.ClearFlag()
	G.Server.LoadServerAddr(
	(
	function()
		StartUp.setTipMessage("连接服务器成功...", 0.4)
		-- 注册消息接收函数
		G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.HallVerInfoID, reciveVer)
		 local sendMsg = { }

        -- 渠道號
		print("######### GameSiteID", G.GameDefine.GameSiteID)
        sendMsg.AreaId = G.GameDefine.GameSiteID
		-- 发送获取版本号消息
		G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestHallVerID,sendMsg)

		StartUp.setTipMessage("正在检查更新...", 0.6)
		
		if G.GameDefine.Branch ~= G.Platform.GetBranch() then
			StartUp.setTipMessage("大厅和平台资料错误...", 0.6)
		end		
	end ),
	( 
	function()
		StartUp.setTipMessage("无法连接到服务器,请检查您的网络是否正常.", 0.4)
		
		if IsShowLoginMessageBox == true then
			local openparam = CS.WindowManager.Instance:FindWindowNodeByName("UILoginMessageBox")
			if openparam == nil then
				openparam = CS.WindowNodeInitParam("UILoginMessageBox", CS.AppDefine.HallNumber)
				openparam.NodeType = CS.BaseNodeType.Propmt;
				CS.WindowManager.Instance:OpenWindow(openparam)
			end
		else
			-- 顯示錯誤訊息
			local text = "无法连接到服务器，点击确定重新连线."
			text = text .. "\n" .. "(" .. G.Server.GetErrorCode() .. ")"
			
			G.Util.coroutine_call( function()
				local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", text, "确定", "取消");
				if boxResult == CS.MessageBoxResult.OK then
					StartUp.BeginConnect()
				else
					CS.UnityEngine.Application.Quit()
				end
			end )()
		end
	end ))
end

-- 主入口函数。从这里开始lua逻辑
function StartUp.Main()
    print("StartUp Main....")

    print("game run ctx init....")

	CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated = CS.HotUpdate.HotUpdateManager.OnHotUpdateComplated + onGameUpdateSucess

	
    G.Init()
	
	if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
		G.Tool.Debug.CheckIP();
	end
	
    G.Tool.Debug.Record(G.GameDefine.GameSiteID,G.Tool.Debug.GetContent(1,0,("0,0")), -1,"StartUp",-1)

	-- 讀取StreamingAssets中的設定檔
	--StartUp.InitGameSetting()
	
    -- 资源热更新地址
    CS.AppDefine.UpdateFileUrl = G.GameConfig.UpdateDomain	
	-- 避免熱更完再次初始化兩個
	UICustomer = CS.WindowManager.Instance:FindWindowNodeByName("UICustomer")
	if UICustomer ~= nil then
		CS.UnityEngine.GameObject.Destroy(UICustomer.WindowGameObject)
	end
	
    -- 显示加载界面
    CS.UILoading.Show(function(WindowNode)	    	
		
		-- 关闭背景
		if G.Platform.isChangeLoadingBG() == true then
			if WindowNode.WindowGameObject.transform:Find("Bg_Texture_43/Image") ~= nil then
				local img_BG = WindowNode.WindowGameObject.transform:Find("Bg_Texture_43/Image"):GetComponent("Image")
				img_BG.gameObject:SetActive(false)
			end
		end
		
		-- 切换LOGO
		if G.Platform.isChangeLoadingLogo() == true then
			if WindowNode.WindowGameObject.transform:Find("LOGO_62") ~= nil then
				local img_LOGO_62 = WindowNode.WindowGameObject.transform:Find("LOGO_62"):GetComponent("Image")
				img_LOGO_62.gameObject:SetActive(false)
			end
			
			if WindowNode.WindowGameObject.transform:Find("LOGO_62_1") ~= nil then
				local img_LOGO_62_1 = WindowNode.WindowGameObject.transform:Find("LOGO_62_1"):GetComponent("Image")
				img_LOGO_62_1.gameObject:SetActive(true)
			end			
		end

	
		if WindowNode.WindowGameObject.transform:Find("Channel") ~= nil then
			local Txt_Channel = WindowNode.WindowGameObject.transform:Find("Channel"):GetComponent("Text")
			--Txt_Channel.text = ""
			Txt_Channel.text = tostring(G.GameConfig.GameSiteID) .. "." .. tostring(G.GameConfig.GamePlatformID) .. "." .. tostring(G.GameConfig.GetServerOption()) .. "." .. G.Platform.GetBranch()
		end
		
		if WindowNode.WindowGameObject.transform:Find("_Common") ~= nil then
			Btn_Common = WindowNode.WindowGameObject.transform:Find("_Common"):GetComponent("Button")
		end
		
		if WindowNode.WindowGameObject.transform:Find("_Website") ~= nil then
			Btn_Website = WindowNode.WindowGameObject.transform:Find("_Website"):GetComponent("Button")
		end
		
		if WindowNode.WindowGameObject.transform:Find("TextCompany") ~= nil then
			TextCompany = WindowNode.WindowGameObject.transform:Find("TextCompany"):GetComponent("Text")			
			local copyright = G.Platform.GetCopyright()
			TextCompany.text = copyright
		end
		
		local _WebsiteUrl = CS.UnityEngine.PlayerPrefs.GetString("WebsiteUrl", "")
		if _WebsiteUrl == "" then
			Btn_Website.gameObject:SetActive(false)
		end
		
		--強制關閉
		if CS.Utility.GetbundleIdentifier() == "com.stage65.test" then
			Btn_Common.gameObject:SetActive(false)
			Btn_Website.gameObject:SetActive(false)
		end

		--判斷是否開啟客服
        local opencs=G.Platform.IsShowStartUpCS(G.GameConfig.GamePlatformID)
        Btn_Common.gameObject:SetActive(opencs)

		if IsShowLoginMessageBox == true then
			Btn_Common:AddClickListener
			(
			function()
				local openparam = CS.WindowNodeInitParam("UICustomer", CS.AppDefine.HallNumber)
				openparam.NodeType = CS.BaseNodeType.Propmt;
				CS.WindowManager.Instance:OpenWindow(openparam)
			end
			)
			Btn_Website:AddClickListener
			(
			function()
				--[[
				local URL = "http://example.invalid/"
			    if G.Env.IsIOS() then
					if G.GameDefine.DownloadUrl ~= "" then
						URL = G.GameDefine.DownloadUrl
					end
				end
				]]--
				CS.UnityEngine.Application.OpenURL(_WebsiteUrl)
			end
			)
		else
			Btn_Common.gameObject:SetActive(false)
		end
	end)
	
	-- 不能移動到UILoading，否則5月底前舊版本會有問題，因為不會再次載入UILoading介面
	G.Util.coroutine_call
	(
	function()		
		G.SetGameInitFinish()
		
		if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
			StartUp.IPGet()
		else
			StartUp.BeginConnect()
			--StartUp.IPGet()
		end
	end
	)()	
end

function StartUp.IPGet()
	print("取自身IP中...")
	if(G.Tool.Debug.GetPublicIP()=="")then
		G.Tool.Debug.CheckIP(	function()
			print("重新尝试取得IP...")
			G.GameConfig.RequestIpCount=G.GameConfig.RequestIpCount+1;
			StartUp.IPGet();
		end);
	else
		G.GameConfig.RequestIpCount=0;		
		StartUp.BeginConnect();
	end

end



function StartUp.Shutdown()
    G.Clear()
end