require("Module")
mmodule(...)

local G = require "Global"
-- ǷIOSֵ
function IsIOSRecharge()
    return false
end

-- Ƿʾ֧
function IsShowOnlinePay()
    return G.User.payinfo.UIFlagPay
end

-- Ƿʾ֧
function IsShowAgentPay()
	if G.User.payinfo.UIFlagAgent == false then
		return false
	end
	
	local user = G.User.GetCurrentUser()
	local userdata = user.UserData	
	local isAnd = G.Platform.getRechargeOperationAnd()
	local result = G.User.CheckOpenData(isAnd, userdata.ShowVIPAgent, userdata.PlaytimeVIPAgent, userdata.ShowVIPAgentBindPhone)
	return result
end

-- Ƿʾ
function IsShowAgentApply()
    return G.User.payinfo.UIFlagApplyAgent
end

-- Ƿʾ
function IsShowJiuLiao()
	if CS.JiuLiaoSDK.s_IsJiuLiaoSDK == true then
		return G.User.payinfo.UIFlagJiuLiao
	else
		return false
	end
end

function IsShow888()
	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
  
    --ثeusWEd
    local rconfig = { }
     
    
        
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="15") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "102";
           pi=pi+1
        end
    end
		return #rconfig>0
	
end

function IsShowCF()
	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
    --ثeusWEd
    local rconfig = { }
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="16") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            
           pi=pi+1
        end
    end
		return #rconfig>0
	
end

function IsShowBaoHe()
	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
    --ثeusWEd
    local rconfig = { }
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="17") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            
           pi=pi+1
        end
    end
		return #rconfig>0
	
end

function IsShowDongGuaZF()
	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
    --ثeusWEd
    local rconfig = { }
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="18") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            
           pi=pi+1
        end
    end
		return #rconfig>0
	
end

function IsShowDongGuaWechat()
	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
    --ثeusWEd
    local rconfig = { }
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="19") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            
           pi=pi+1
        end
    end
		return #rconfig>0
	
end



-- Ƿʾ
function IsShowShanFu()
	-- SDKǷ֧Ԯ
	if CS.HenrySDK.s_IsHenrySDK == false then
		return false
	end
	
	-- 
	if G.User.payinfo.UIFlagShanFu == false then
		return false
	end

	local config = G.User.GetCurrentUser().UserData.PaymentInfoList
	local pi=1 
    --ثeusWEd
    local rconfig = { }
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="0") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];            
           pi=pi+1
        end
    end


	if(#rconfig==0)then
		return false;
	end
	
	local user = G.User.GetCurrentUser()
	local userdata = user.UserData
	local isAnd = G.Platform.getRechargeOperationAnd()
	local result = G.User.CheckOpenData(isAnd, userdata.ShowVIPAgent, userdata.PlaytimeVIPAgent, userdata.ShowVIPAgentBindPhone)
	return result
end

-- Ƿʾģʽ1
function IsShowClub2AgentMode1()
    return G.User.payinfo.UIClub2AgentMode1
end

-- Ƿʾģʽ1
function IsShowClub2AgentMode1()
    return G.User.payinfo.UIClub2AgentMode1
end

-- Ƿʾģʽ3
function IsShowClub2AgentMode2()
    return G.User.payinfo.UIClub2AgentMode2
end

-- Ƿʾģʽ3
function IsShowClub2AgentMode3()
    return G.User.payinfo.UIClub2AgentMode3
end

-- Ƿʾģʽ4
function IsShowClub2AgentMode4()
    return G.User.payinfo.UIClub2AgentMode4
end

-- Ƿʾģʽ5
function IsShowClub2AgentMode5()
    return G.User.payinfo.UIClub2AgentMode5
end