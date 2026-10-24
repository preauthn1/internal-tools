require("Module")
mmodule(...)

local G = require "Global"

RequestIpCount=0;
RequestGameRemoteVersion=0;
RequestGameServerFlag=0;
SoundStatus=true;
MusicStatus=true;

-- endregion

-- 是否使用新的Transfer機制，資料加密，回傳token
local useNewTransfer = true
function isUseNewTransfer()
	return useNewTransfer;
end

-- 游戏自动准备倒计时
local autoReadyCountNum = 3
function GetAutoReadyCountNum()
    return autoReadyCountNum
end

-- 环境定义字符串
local VERSIONTYPE_LOCAL = "local"
local VERSIONTYPE_ONLINETEST = "online"
local VERSIONTYPE_TEST = "test"
local VERSIONTYPE_FORMAL = "formal"
local ServerOption = G.GameDefine.ServerOption
--local ServerOption = CS.AppDefine.ServerOption
--ServerOption = 3
function GetServerOption()
	return ServerOption
end

function ReSetServerOption(so,pi,si)
    ServerOption=so;
    GamePlatformID=pi
    GameSiteID=si
end

-- 当前运行环境
local VERSIONTYPE = VERSIONTYPE_FORMAL

-- 非编辑器环境
if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
    VERSIONTYPE = VERSIONTYPE_FORMAL
end

-- 游戏官网
WebSiteDomain = "http://example.invalid/"
   
-- 充值地址     
ChargeDomain = "http://0.0.0.0/paycenter/index?orderid=%d"

-- serverlist地址     
ServerListDomain = { "dGYuMjg4ZG0uY29t","d3d3Ljg1enIuY24=","MjcuMTQ4LjE5Ni4xMg==","d3d3LjYzYncuY24=","d3d3LjY4eHMuY24=","d3d3Lmp1ZGVyb25nLmNu","MTgwLjEwMS43NS41NA=="}

-- User IP
UserIP = nil

-- serverlist檔案地址     
ServerListFile = nil

-- 游戏渠道编号
GameSiteID = G.GameDefine.GameSiteID

-- 游戏平台编号         
GamePlatformID = G.GameDefine.GamePlatformID      
 
-- 获取本机Ip地址
IPReuestUrl = "http://0.0.0.0:9132/client/getip"

-- 获取本机Ip地址(NEW)
IPReuestUrlNew = "http://0.0.0.0:9132/client/RegIP"

-- 资源更新地址
UpdateDomain = "http://0.0.0.0/UpdateFiles/"

-- 外网地址
DefaultServerIP = "0.0.0.0"

-- 内网和外网端口
DefaultServerPort = 10000
 
-- Android整包地址
AndroidAPPUrl = ""  

-- IOS整包地址  
IOSAPPUrl = ""  

-- IOS APP ID
IOSAPPID = ""

IOSRecharge="http://example.invalid/IAP/Pay?"

ProductIdentifiers={ CS.Utility.GetbundleIdentifier()..".100coins", CS.Utility.GetbundleIdentifier()..".500coins", CS.Utility.GetbundleIdentifier()..".1000coins" }

RechargeListDomain="http://0.0.0.0"

QRCodeCopyDomain="http://0.0.0.0"

CustomServiceDomain="http://example.invalid"

-- 登入客服
ContactService="http://0.0.0.0/client/ContactService"
-- 傳送手機資訊
RegHD = ""

-- 是否本地环境
function IsLocalCtx()
    return VERSIONTYPE == VERSIONTYPE_LOCAL
end

-- 外网
function IsOnlineTextCtx()
    return VERSIONTYPE == VERSIONTYPE_ONLINETEST
end

-- 是否测试环境
function IsTestCtx()
    return VERSIONTYPE == VERSIONTYPE_TEST
end

-- 是否正式环境
function IsFormalCtx()
    return VERSIONTYPE == VERSIONTYPE_FORMAL
end

-- 是否送審階段 Add By Rayon
Submission = false
function IsSubmissionCtx()
	--return G.Service.UIFlag.IsIOSRecharge()
	return Submission
end

-- 是否顯示充值階段 Add By Rayon
Recharge = true
function IsRechargeCtx()
	-- 有支付方式開啟才算能充值
	--[[
	if (G.Service.UIFlag.IsShowOnlinePay() == true or
		G.Service.UIFlag.IsShowAgentPay() == true or
		G.Service.UIFlag.IsShowAgentApply() == true) then
		return true
	else
		return false
	end
	]]--
	return Recharge
end

-- 本地
if VERSIONTYPE == VERSIONTYPE_LOCAL then

	-- Log開關
	CS.RuntimeDebug.s_instance:SetShow(true)
	-- FPS開關
	CS.RuntimeStats.s_instance:SetShow(true)
	
	-- serverlist地址     
	ServerListDomain = {  }
	
    -- 资源更新地址
    UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

    -- 大厅服务器IP
    DefaultServerIP = "0.0.0.0"

    -- 大厅服务器端口
    DefaultServerPort = 10000

    -- Android整包地址
    AndroidAPPUrl = ""
	
	-- IOS整包地址
	IOSAPPUrl = ""

    -- 充值地址
    ChargeDomain = ""

    -- 获取本机Ip地址
    IPReuestUrl = "http://0.0.0.0:9132/client/getip"

	-- 获取本机Ip地址(NEW)
	IPReuestUrlNew = "http://0.0.0.0:9132/client/RegIP"	
	
    QRCodeCopyDomain="http://0.0.0.0"
end

-- 外网
if VERSIONTYPE == VERSIONTYPE_ONLINETEST then

	-- Log開關
	CS.RuntimeDebug.s_instance:SetShow(true)
	-- FPS開關
	CS.RuntimeStats.s_instance:SetShow(true)

	-- serverlist地址     
	ServerListDomain = { "http://0.0.0.0:160" }
	
    -- 外网资源地址
    UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

    -- 大厅服务器IP
    DefaultServerIP = "0.0.0.0"

    -- 大厅服务器端口
    DefaultServerPort = 10000

    -- Android整包地址
    AndroidAPPUrl = "http://0.0.0.0:140/fullgame.apk"

	-- IOS整包地址
	IOSAPPUrl = ""
	
    -- 充值地址
    ChargeDomain = "http://0.0.0.0:9132/PayCenter/Index?userid=%d&paytype=%d&amount=%d&orderid=%d&platformid=%d&areaid=%d"

    -- 获取本机Ip地址
    IPReuestUrl = "http://0.0.0.0/client/getip"
	
	-- 获取本机Ip地址(NEW)
	IPReuestUrlNew = "http://0.0.0.0/client/RegIP"	
end

-- 测试
if VERSIONTYPE == VERSIONTYPE_TEST then

	-- Log開關
	CS.RuntimeDebug.s_instance:SetShow(false)
	-- FPS開關
	CS.RuntimeStats.s_instance:SetShow(false)
	
	-- serverlist地址     
	ServerListDomain = { }

    -- 资源更新地址
    UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

    -- 大厅服务器IP
    DefaultServerIP = "0.0.0.0"

    -- 大厅服务器端口
    DefaultServerPort = 10000

    -- Android整包地址
    AndroidAPPUrl = ""

	-- IOS整包地址
	IOSAPPUrl = ""
	
    -- 充值地址
    ChargeDomain = ""

    -- 获取本机Ip地址
    IPReuestUrl = "http://0.0.0.0/client/getip"
	
	-- 获取本机Ip地址(NEW)
	IPReuestUrlNew = "http://0.0.0.0/client/RegIP"	
end

-- 正式
if VERSIONTYPE == VERSIONTYPE_FORMAL then

    if ServerOption == 0 then
		--////////////////////////////////////////////////////////////////////
		--超級大贏家-內部測試
		--com.cqxx.myGametest

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(true)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(true)

		-- serverlist地址     
		ServerListDomain = { "MTkyLjE2OC4xMC4xMDA6MTYw" }
		
		-- serverlist檔案地址
		ServerListFile = "http://0.0.0.0:140"

		-- 外网资源地址
		UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = "0.0.0.0"

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://0.0.0.0:140/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = "http://0.0.0.0/paycenter/index?orderid=%d"

		-- 获取本机Ip地址
		IPReuestUrl =  "http://0.0.0.0/client/getip"
				
		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://0.0.0.0/client/RegIP"
		
        RechargeListDomain="http://0.0.0.0"

        CustomServiceDomain="http://0.0.0.0"
		
		-- 傳送手機資訊
		RegHD = "http://0.0.0.0/client/reghd"		
		
		-- 登入客服
		ContactService="http://0.0.0.0/client/ContactService"
elseif ServerOption == 1 then
		--////////////////////////////////////////////////////////////////////
		--超級大贏家-世紀富豪版
		--com.cqxx.myGame

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(true)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(true)

		-- serverlist地址
		ServerListDomain = { "NjEuMjE4LjQ0LjIzMToxNjA=" }
		
		-- serverlist檔案地址
		ServerListFile = "http://0.0.0.0:140"

		-- 外网资源地址
		UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = "0.0.0.0"

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://0.0.0.0:140/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = "http://0.0.0.0/paycenter/index?orderid=%d"

		-- 获取本机Ip地址
		IPReuestUrl =  "http://0.0.0.0/client/getip"
   
		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://0.0.0.0/client/RegIP"
   
        RechargeListDomain="http://0.0.0.0"

        CustomServiceDomain="http://0.0.0.0"
        
        QRCodeCopyDomain="http://0.0.0.0"

		-- 傳送手機資訊
		RegHD = "http://0.0.0.0/client/reghd"		

		-- 登入客服
		ContactService="http://0.0.0.0/client/ContactService"
elseif ServerOption == 2 then
		--////////////////////////////////////////////////////////////////////
		--超級大贏家-壓力測試
		--com.cqxx.myGamepress

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(true)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(true)

		-- serverlist地址     
		ServerListDomain = { "NjEuMjE4LjQ0LjIzMjoxNjA=" }
		
		-- serverlist檔案地址
		ServerListFile = "http://0.0.0.0:140"

		-- 外网资源地址
		UpdateDomain = "http://0.0.0.0:140/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = "0.0.0.0"

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://0.0.0.0:140/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = ""

		-- 获取本机Ip地址
		IPReuestUrl =  "http://0.0.0.0/client/getip"

		
		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://0.0.0.0/client/RegIP"		

		-- 傳送手機資訊
		RegHD = "http://0.0.0.0/client/reghd"
		
		-- 登入客服
		ContactService="http://0.0.0.0/client/ContactService"
    elseif ServerOption == 3 then

		--////////////////////////////////////////////////////////////////////
		--超級大贏家-美版
		--com.cqxx.myGame

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(false)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(false)

		-- serverlist地址     
		ServerListDomain = { "http://" }
		
		-- serverlist檔案地址
		ServerListFile = "http://example.invalid/serverlist_otuey.txt"

		-- 外网资源地址
		UpdateDomain = "http://example.invalid/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = ""

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://example.invalid/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = "http://example.invalid/paycenter/index?orderid=%d"

		-- 获取本机Ip地址
		IPReuestUrl =  "http://example.invalid/client/Getip"

		
		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://example.invalid/client/RegIP"

		-- 傳送手機資訊
		RegHD = "http://example.invalid/client/reghd"
		
		-- 登入客服
		ContactService="http://example.invalid/client/ContactService"
	elseif ServerOption == 4 then
		--////////////////////////////////////////////////////////////////////
		--超級大贏家-陸版
		--com.cqxx.myGame

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(false)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(false)

		-- serverlist地址     
		ServerListDomain = {
								"dGYuMmhlNy5jb20=",
								"dGYuZGlmcm9hZC5jb20=",
								"MjcuMTQ4LjE5Ni4xMg==",
								"dGYuOTRoYW9xdWFuLmNvbQ==",
								"dGYuaGVyb25nc2hlbmcuY29t",
								"dGYuMTcydnIuY29t",
								"dGYuamluYW4yMDE3LmNvbQ=="
							}
		
		-- serverlist檔案地址
		ServerListFile = "http://example.invalid"

		-- 外网资源地址
		UpdateDomain = "http://example.invalid/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = ""

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://example.invalid/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = "http://example.invalid/paycenter/index?orderid=%d"

		-- 获取本机Ip地址
		IPReuestUrl =  "http://example.invalid/client/Getip"

		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://example.invalid/client/RegIP"		
		
        RechargeListDomain="http://example.invalid"

        QRCodeCopyDomain="http://example.invalid"
        --開啟遷入客服位置
        CustomServiceDomain="http://example.invalid"

		-- 傳送手機資訊
		RegHD = "http://example.invalid/client/reghd"
		
		-- 登入客服
		ContactService="http://example.invalid/client/ContactService"
    elseif ServerOption == 5 then
		--////////////////////////////////////////////////////////////////////
		--超級大贏家-世紀富豪版
		--com.cqxx.myGame

		-- Log開關
		CS.RuntimeDebug.s_instance:SetShow(false)
		-- FPS開關
		CS.RuntimeStats.s_instance:SetShow(false)

		-- serverlist地址
		ServerListDomain = { "dGYudGVycHJvLmNvbQ==" }
		
		-- serverlist檔案地址
		ServerListFile = "http://0.0.0.0:140"

		-- 外网资源地址
		UpdateDomain = "http://example.invalid/UpdateFiles/"

		-- 大厅服务器IP
		DefaultServerIP = "0.0.0.0"

		-- 大厅服务器端口
		DefaultServerPort = 10000

		-- Android整包地址
		AndroidAPPUrl = "http://example.invalid/" .. G.GameDefine.AppName .. ".apk"

		-- IOS整包地址
		IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID
		
		-- 充值地址
		ChargeDomain = "http://0.0.0.0/paycenter/index?orderid=%d"

		-- 获取本机Ip地址
		IPReuestUrl =  "http://0.0.0.0/client/getip"
		
		-- 获取本机Ip地址(NEW)
		IPReuestUrlNew = "http://0.0.0.0/client/RegIP"
		
        RechargeListDomain="http://0.0.0.0"

        CustomServiceDomain="http://0.0.0.0"
        
        QRCodeCopyDomain="http://0.0.0.0"

		-- 傳送手機資訊
		RegHD = "http://example.invalid/client/reghd"
		
		-- 登入客服
		ContactService="http://0.0.0.0/client/ContactService"
    end
end

--////////////////////////////////////////////////////////////////////
--超級大贏家-伺服器
--com.cqxx.myGameserver
if CS.AppDefine.UseChangeServerIP then
    -- Log開關
    CS.RuntimeDebug.s_instance:SetShow(true)
    -- FPS開關
    CS.RuntimeStats.s_instance:SetShow(true)

    -- 大厅服务器IP
    DefaultServerIP = CS.AppDefine.DefaultServerIP

    -- 大厅服务器端口
    DefaultServerPort = CS.AppDefine.DefaultServerPort

    -- serverlist檔案地址
    ServerListFile = "http://" .. DefaultServerIP .. ":" .. DefaultServerPort

    -- serverlist地址     
    ServerListDomain = { CS.AppDefine.ServerListDomain  }

    -- 充值地址
    ChargeDomain = "http://" .. CS.AppDefine.ChargeDomainIP .. "/paycenter/index?orderid=%d"

    -- 获取本机Ip地址
    IPReuestUrl =  "http://" .. CS.AppDefine.IPReuestIP .. "/client/getip"

	-- 获取本机Ip地址(NEW)
	IPReuestUrlNew = "http://" .. CS.AppDefine.IPReuestIP .. "/client/RegIP"	
	
    -- 外网资源地址
    UpdateDomain = "http://" .. CS.AppDefine.UpdateDomainIP .. ":"..CS.AppDefine.UpdateDomainPort.."/UpdateFiles/"

    -- Android整包地址
    AndroidAPPUrl = "http://" .. CS.AppDefine.UpdateDomainIP .. ":"..CS.AppDefine.UpdateDomainPort.."/" .. G.GameDefine.AppName .. ".apk"

    -- IOS整包地址
    IOSAPPUrl = "https://example.invalid/app/id" .. G.GameDefine.AppleID


    if CS.AppDefine.LoginProcessType == 1 then
        ServerListFile = nil
        VERSIONTYPE = VERSIONTYPE_FORMAL
		local data = tostring(DefaultServerIP) .. ":" .. tostring(DefaultServerPort)
		data = require("Tools/Base64").encodeBase64(data)
		ServerListDomain = { data }		
    elseif CS.AppDefine.LoginProcessType == 2 then
        ServerListFile = nil
        VERSIONTYPE = VERSIONTYPE_TEST
    end
end
