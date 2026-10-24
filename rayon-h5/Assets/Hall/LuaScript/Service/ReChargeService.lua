require("Module")
mmodule(...)

local G = require "Global"
local yield_return = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback(to_to_yield, cb)
end
)

local groupName = "UIHall";

local payConfig
-- oRȦC
function GetPayConfig()
	return payConfig
end

-- dINzH
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ProxyRequestRetID, function(data)
	if data.type == "1" then
		payConfig = data
	end
end, groupName )

-- oRȸ
function GetRechargeList()
	local url = G.GameConfig.RechargeListDomain .. "/Pay/GetAmountList/"..tostring(G.GameConfig.GamePlatformID)
    -- local url = G.GameConfig.RechargeListDomain .. "/PayCenter/amountlist/"..tostring(G.GameConfig.GamePlatformID)
	sendata={};
	sendata.data=url;
	sendata.type="1";
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ProxyRequestID, sendata)
end

function GetPayAmountConfig()
    local configstr = G.User.GetCurrentUser().UserData.PayAmountConfig

    local config = { }
    for _, v in pairs(G.Tool.StringHelper.Split(configstr, "|")) do
        local payType = G.Tool.StringHelper.Split(v, ":")
        local k = payType[1]
        local list = G.Tool.StringHelper.Split(payType[2], ".")
        config[k] = list
    end
    return config
end

function GetPayAmountRange()
    

	
	
    local configstr = G.User.GetCurrentUser().UserData.PayAmountRange

    local config = { }
	if (configstr == nil or configstr == '') then
		return nil
	end
	
	for _, v in pairs(G.Tool.StringHelper.Split(configstr, "|")) do
		local payType = G.Tool.StringHelper.Split(v, ":")
		local k = payType[1]
		local list = G.Tool.StringHelper.Split(payType[2], ".")
		config[k] = list
	end
    return config
end

function GetPaymentInfoList(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList
    
    --ثeusWEd
    local rconfig = { }
     
    for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="6") then
			local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "101";
            rconfig[(pi)][6]= 6;
            rconfig[(pi)][7]= 6;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
    --
    if (CS.JiuLiaoSDK.s_IsJiuLiaoSDK ~= true) then
        rconfig = { }
    end
	
	if G.GameDefine.JiuLiaoId == "" or G.GameDefine.JiuLiaoSecrety == "" or G.GameDefine.JiuLiaoName == "" then
        rconfig = { }	
	end

    return rconfig
end

function Get888PayInfo(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList
  
  
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
            rconfig[(pi)][6]= 15;
            rconfig[(pi)][7]= 10;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
       
    --
   

    return rconfig
end

function GetCFPayInfo(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList   
    --ثeusWEd
    local rconfig = { }
        
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="16") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "103";
            rconfig[(pi)][6]= 16;
            rconfig[(pi)][7]= 11;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
    return rconfig
end
function GetBaoHePayInfo(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList   
    --ثeusWEd
    local rconfig = { }
        
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="17") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "104";
            rconfig[(pi)][6]= 17;
            rconfig[(pi)][7]= 16;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
    return rconfig
end

function GetDongGuaZFPayInfo(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList   
    --ثeusWEd
    local rconfig = { }
        
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="18") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "105";
            rconfig[(pi)][6]= 18;
            rconfig[(pi)][7]= 19;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
    return rconfig
end
function GetDongGuaWeChatPayInfo(pi)
    local config = G.User.GetCurrentUser().UserData.PaymentInfoList   
    --ثeusWEd
    local rconfig = { }
        
        for _, v in ipairs(config) do
        if(tostring(v.TypeID)=="19") then
            local arr = G.Tool.StringHelper.Split(v.AmountSort, ",");
            rconfig[(pi)]={};
            rconfig[(pi)][1]= arr[1];
            rconfig[(pi)][2]= arr[#arr];
            rconfig[(pi)][3]={};
            rconfig[(pi)][3]= arr
            rconfig[(pi)][4]= 1;
            rconfig[(pi)][5]= "106";
            rconfig[(pi)][6]= 19;
            rconfig[(pi)][7]= 19;
            rconfig[(pi)][8]= 0;
           pi=pi+1
        end
    end
    return rconfig
end
