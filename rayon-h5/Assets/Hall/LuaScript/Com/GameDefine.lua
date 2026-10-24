require("Module")
require("Com/PlatformRoute")
mmodule(...)

local G = require "Global"
 ----------APP INFO
 ProductName = "Not Assign"
 AppName = "Not Assign"
 AppleID = ""--GetTable().AppleID
 ServerOption = GetTable().ServerOption
 GamePlatformID = GetTable().PlatformID
 GameSiteID = GetTable().SiteID
 Version = 0
 Agent =1
 FlagExchage=1
 ExchangeFontType=0
 Template="0"
 Club2Agent=0
 ShowCaiJinchi=1
 Branch = GetTable().Branch

if Agent == nil then
	Agent = 1
end
if Template == nil then
	Template="0"
end

------------ JiuLiao
JiuLiaoId = ""--GetAppInfoTable().JiuLiaoId
JiuLiaoSecrety =""-- GetAppInfoTable().JiuLiaoSecrety
JiuLiaoName = ""--GetAppInfoTable().JiuLiaoName
DownloadUrl = ""--GetAppInfoTable().DownloadUrl
UIFlagShanFu = ""
		
function SetServerjson(serverjson)
   print(serverjson)
    local InfoData = nil; 
    local JiuLiaoData = nil;
    local payconfig = nil; 
    local class = {}
    local prtb = CS.MiniJSON.Json.Deserialize(serverjson)
	local iter = prtb:GetEnumerator()
    local data=nil;
    while iter:MoveNext() do
    local v = iter.Current
        if v.Key == "systemConfig" then
            data=v.Value:GetEnumerator();
          
        end
    end

	while data:MoveNext() do
		local v = data.Current


		if v.Key == "Info" then
			InfoData = v.Value;
		elseif v.Key == "JiuLiao" then
			JiuLiaoData = v.Value;
        elseif v.Key=="payConfig" then
            payconfig=v.Value;
		end
	end
    SetInfo(InfoData);
    SetJiuLiao(JiuLiaoData);
    SetPayInfo(payconfig);
end

function SetInfo(i)
    local iter=nil; 
    if i ~= nil then
		iter = i:GetEnumerator()
	    while iter:MoveNext() do
	    	local v = iter.Current
			if  v.Key == "ProductName" then
				ProductName = v.Value			
            elseif v.Key == "AppName" then
				AppName = v.Value
            elseif v.Key == "AppleID" then
				--AppleID = v.Value
				--指定充值会挂掉
				AppleID = ""
            elseif v.Key == "ServerOption" then
				--ServerOption = v.Value
            elseif v.Key == "PlatformID" then
				GamePlatformID = v.Value
            elseif v.Key == "SiteID" then
				GameSiteID = v.Value
            elseif v.Key == "Agent" then
				Agent = v.Value
            elseif v.Key == "FlagExchage" then
				FlagExchage = v.Value
            elseif v.Key == "ExchangeFontType" then
				ExchangeFontType = v.Value
            elseif v.Key == "Template" then
				Template = v.Value
				if Template == nil or Template == "" then
					Template = "0"
				end
			elseif v.Key == "Club2Agent" then
				Club2Agent = v.Value
			elseif v.Key == "ShowCaiJinchi" then
				ShowCaiJinchi = v.Value
			end
		end

        G.GameConfig.ReSetServerOption(ServerOption,GamePlatformID,GameSiteID)
    end
end

function SetJiuLiao(j)
    local iter=nil; 
    if j ~= nil then
		iter = j:GetEnumerator()
	    while iter:MoveNext() do
	    	local v = iter.Current
			if  v.Key == "JiuLiaoId" then
				JiuLiaoId = v.Value
            elseif v.Key == "JiuLiaoSecrety" then
				JiuLiaoSecrety = v.Value
            elseif v.Key == "JiuLiaoName" then
				JiuLiaoName = v.Value
            elseif v.Key == "DownloadUrl" then
				DownloadUrl = v.Value
				CS.UnityEngine.PlayerPrefs.SetString("WebsiteUrl", DownloadUrl)
				CS.UnityEngine.PlayerPrefs.Save()
			end
		end
		
		if JiuLiaoId == nil then
			JiuLiaoId = ""
		end		
		if JiuLiaoSecrety == nil then
			JiuLiaoSecrety = ""
		end
		if JiuLiaoName == nil then
			JiuLiaoName = ""
		end
		if DownloadUrl == nil then
			DownloadUrl = ""
		end
 	else
		JiuLiaoId = ""
		JiuLiaoSecrety =""
		JiuLiaoName = ""
		DownloadUrl = ""
    end
end

function SetPayInfo(j)
    G.User.payinfo={}
	G.User.payinfo.UIFlagPay = false
	G.User.payinfo.UIFlagAgent = false
	G.User.payinfo.UIFlagApplyAgent = false
	G.User.payinfo.UIFlagJiuLiao = false
	G.User.payinfo.UIFlagShanFu = false
	G.User.payinfo.UIClub2AgentMode1 = false
	G.User.payinfo.UIClub2AgentMode2 = false
	G.User.payinfo.UIClub2AgentMode3 = false
	G.User.payinfo.UIClub2AgentMode4 = false
	G.User.payinfo.UIClub2AgentMode5 = false
	G.User.payinfo.UIClub2AgentMode6 = false
	
    G.User.andor={}
	
    local iter=nil; 
    if j ~= nil then
		iter = j:GetEnumerator()
	    while iter:MoveNext() do
	    	local v = iter.Current
			if  v.Key == "UIFlagPay" then
				G.User.payinfo.UIFlagPay = v.Value
            elseif v.Key == "UIFlagAgent" then
				G.User.payinfo.UIFlagAgent = v.Value
            elseif v.Key == "UIFlagApplyAgent" then
				G.User.payinfo.UIFlagApplyAgent = v.Value
            elseif v.Key == "UIFlagJiuLiao" then
				G.User.payinfo.UIFlagJiuLiao = v.Value
			elseif v.Key == "UIFlagShanFu" then
				G.User.payinfo.UIFlagShanFu = v.Value
			elseif v.Key == "UIClub2AgentMode1" then
				G.User.payinfo.UIClub2AgentMode1 = v.Value
			elseif v.Key == "UIClub2AgentMode2" then
				G.User.payinfo.UIClub2AgentMode2 = v.Value
			elseif v.Key == "UIClub2AgentMode3" then
				G.User.payinfo.UIClub2AgentMode3 = v.Value
			elseif v.Key == "UIClub2AgentMode4" then
				G.User.payinfo.UIClub2AgentMode4 = v.Value
			elseif v.Key == "UIClub2AgentMode5" then
				G.User.payinfo.UIClub2AgentMode5 = v.Value
			elseif v.Key == "UIClub2AgentMode6" then
				G.User.payinfo.UIClub2AgentMode6 = v.Value
			elseif v.Key == "ConvertAndOr" then
				G.User.andor.ConvertAndOr = (v.Value == false)
			elseif v.Key == "AllGameAndOr" then
				G.User.andor.AllGameAndOr = (v.Valuer == false)
			elseif v.Key == "VIPAgentAndOr" then
				G.User.andor.VIPAgentAndOr = (v.Value == false)
			elseif v.Key == "CSAndOr" then
				G.User.andor.CSAndOr = (v.Value == false)
			end
		end
 	else
 			
    end
	--G.User.payinfo.UIClub2AgentMode6 = true	
end




if CS.AppDefine.UseChangeServerIP then
	-- 游戏渠道编号 
	GameSiteID = CS.AppDefine.GameSiteID
	-- 游戏平台编号         
	GamePlatformID = CS.AppDefine.GamePlatformID

    ServerOption=CS.AppDefine.ServerOption
end

-- 游戏渠道编号
--GameSiteID = 2
-- 游戏平台编号         
--GamePlatformID = 2

 
