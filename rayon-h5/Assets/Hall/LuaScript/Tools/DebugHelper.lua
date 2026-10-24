require("Module")
mmodule(...)
local G = require "Global"

local PublicIP="";

local yield_return = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback(to_to_yield, cb)
end
)

--获取页面内容
function GetPageData(url,onSuccess, onError)
 G.Util.coroutine_call(function ()
   local www = CS.UnityEngine.WWW(url)
   yield_return(www)

   if www.error ~= nil then
        if onError ~= nil  then
            onError(www.error)
        end

        return 
   end

   if onSuccess ~= nil then
        onSuccess(www.text,onError)
   end

end)()

end

function PostPageData(url,content,onSuccess, onError)
   local f=CS.UnityEngine.WWWForm();
  
   for key,value in pairs(content) do   
    f:AddField(key, value);
   end


   G.Util.coroutine_call(function ()
   local www = CS.UnityEngine.WWW(url, f)
   yield_return(www)

   if www.error ~= nil then
        if onError ~= nil  then
            onError(www.error)
        end

        return 
   end

   if onSuccess ~= nil then
        onSuccess(www.text)
   end

end)()
  

end

function GetPublicIP()

 return PublicIP;
end

function RecordGet(c)
   

    if(PublicIP=="")then
		if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android or CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.IPhonePlayer then
			CheckIP()			
		end
    end
   

    local url= "http://example.invalid/api/target/record/"
   -- local url= "http://localhost:1132/api/target/record"
    
    local localIP2=G.GetLocalIP();
    -- if (type(CS.UnityEngine.Network.player.ipAddress) == "table") then
        -- print("<color=yellow>IP~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~</color>")
        -- G.Tool.TableHelper.Print(CS.UnityEngine.Network.player.ipAddress)
    -- else
       -- localIP2=CS.UnityEngine.Network.player.ipAddress
    -- end
    local deviceinfo="0";
    local ServerOP=G.GameConfig.GetServerOption()
    local DeviceId=CS.Utility.GetUUID
    --當取不到uuid時版本為舊版本,則回傳預設值 update by bob
    -- if(uuid==nil) then
    --   DeviceId= CS.UnityEngine.SystemInfo.deviceUniqueIdentifier;
    -- end
    if(c.dev~=nil and c.dev)then
        deviceinfo=(G.Tool.Other.urlencode(GetDeviceInfo()))
    end
    
    local data="?AreaId="..G.GameDefine.GameSiteID.."&"..
               "BundleId="..G.Tool.Other.urlencode(CS.Utility.GetbundleIdentifier()).."&"..
               "deviceinfo=".. deviceinfo.."&"..
               "IP="..G.Tool.Other.urlencode( PublicIP..",".. G.GetLocalIP() .. "," ..localIP2 ).."&"..
               "LogContent=".. G.Tool.Other.urlencode(G.Tool.Json.encode(c.content)).."&"..
               "LogType=".. tostring(c.lt).."&"..
               "PlatformId=".. G.GameDefine.GamePlatformID .."&"..
               "SceneName=".. c.sceneName .."&"..
               "UserID="..c.uid.."&"..
               "DeviceId="..DeviceId.."&"..
               "ServerOP="..ServerOP;
               data=G.Tool.StringHelper.TrimAll(data)
   -- data=G.Tool.Other.urlencode(data)
   print(url..data)
    GetPageData((url..data),
        function(msg)
			print("Record成功"..msg)
		end,
		function(msg)
			print("Record失敗"..msg)
		end)
    
end

function SendTelegram(t,c)

   local chatid="CLIENTNOTIFY"
   if(c=="c")then
    chatid="CST";
   elseif(c=="r")then
    chatid="RAYON";
   elseif(c=="b")then
    chatid="BOB";
   end


    local url= "http://example.invalid/api/target/telegramSend?".."ChatId="..chatid.."&text="..G.Tool.StringHelper.TrimAll(t);
    -- local url= "http://localhost:1132/api/target/record"
    
    print(url)
    GetPageData((url),
    function(msg)
      print("Send成功"..msg)
    end,
    function(msg)
      print("Send失敗"..msg)
    end)
    
end


function Record(aid,content,pid,sceneName,uid,lt)
	--[[
    if(PublicIP=="")then
        CheckIP()
    end
    lt=lt or 1

    local url= "http://0.0.0.0:6500/api/target/record/"
    local localIP2="";
    if (type(CS.UnityEngine.Network.player.ipAddress) == "table") then
        print("<color=yellow>IP~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~</color>")
        G.Tool.Print(CS.UnityEngine.Network.player.ipAddress)
    else
       localIP2=CS.UnityEngine.Network.player.ipAddress
    end
    
    local data={
        AreaId=aid,
        BundleId=CS.Utility.GetbundleIdentifier(),
        deviceinfo=GetDeviceInfo(),
        IP=( PublicIP..",".. G.GetLocalIP() .. "," ..localIP2 ),
        LogContent=G.Tool.Json.encode(content),
        LogType=lt,
        PlatformId=pid,
        SceneName=sceneName,
        UserID=uid
    }
   
    PostPageData(url,data,
        function(msg)
			print("Record成功"..msg)
		end,
		function(msg)
			print("Record失敗"..msg)
		end)
    ]]
end

function GetDeviceInfo()
local deviceinfo=
{
    batteryLevel=tostring(CS.UnityEngine.SystemInfo.batteryLevel),   
    batteryStatus= tostring(CS.UnityEngine.SystemInfo.batteryStatus),
    copyTextureSupport=tostring(CS.UnityEngine.SystemInfo.copyTextureSupport),
    deviceName=CS.UnityEngine.SystemInfo.deviceName,
    deviceType=tostring(CS.UnityEngine.SystemInfo.deviceType),
    deviceModel=tostring(CS.UnityEngine.SystemInfo.deviceModel),
    deviceUniqueIdentifier=tostring(CS.UnityEngine.SystemInfo.deviceUniqueIdentifier),
    graphicsDeviceID=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceID),
    graphicsDeviceName=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceName),
    graphicsDeviceType=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceType),
    graphicsDeviceVendor=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVendor),
    graphicsDeviceVendorID=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVendorID),
    graphicsDeviceVersion=tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVersion),
    graphicsMemorySize=tostring(CS.UnityEngine.SystemInfo.graphicsMemorySize),
    graphicsMultiThreaded=tostring(CS.UnityEngine.SystemInfo.graphicsMultiThreaded),
    graphicsShaderLevel=tostring(CS.UnityEngine.SystemInfo.graphicsShaderLevel),
    processorCount=tostring(CS.UnityEngine.SystemInfo.processorCount),
    processorFrequency=tostring(CS.UnityEngine.SystemInfo.processorFrequency),
    processorType=tostring(CS.UnityEngine.SystemInfo.processorType),
    supportedRenderTargetCount=tostring(CS.UnityEngine.SystemInfo.supportedRenderTargetCount),
    supports2DArrayTextures=tostring(CS.UnityEngine.SystemInfo.supports2DArrayTextures),
    supports32bitsIndexBuffer=tostring(CS.UnityEngine.SystemInfo.supports32bitsIndexBuffer),
    supports3DRenderTextures=tostring(CS.UnityEngine.SystemInfo.supports3DRenderTextures),
    supports3DTextures=tostring(CS.UnityEngine.SystemInfo.supports3DTextures),
    supportsAccelerometer=tostring(CS.UnityEngine.SystemInfo.supportsAccelerometer),
    supportsAudio=tostring(CS.UnityEngine.SystemInfo.supportsAudio),
    systemMemorySize=tostring(CS.UnityEngine.SystemInfo.systemMemorySize),
    unsupportedIdentifier=tostring(CS.UnityEngine.SystemInfo.unsupportedIdentifier),
    usesReversedZBuffer=tostring(CS.UnityEngine.SystemInfo.graphicsShaderLevel),
}
 -- print(G.Tool.StringHelper.string.utf8len(G.Tool.Json.encode(deviceinfo)));
 return G.Tool.Json.encode(deviceinfo);

end

function GetContent(contenttype,flag,data)
    
    --ContentType 指定使用哪种搜集方式
    --  1   ：   使用者大厅点选行为（仅记录X,Y）再依照FLAG来区分点击事件

    local p =
    {
        ContentType=contenttype,
        Flag=flag,
        Data=data
    }
return p
end

function CheckIP(callback) 
  local toaip="https://example.invalid"
    -- 支援游戏盾
  if G.Tool.Aliyun.IsAliyunSDK() == true then
    -- 已經初始化成功，可以直接取IP
    if G.Tool.Aliyun.IsInitSuccess() == true then
       -- 取得遊戲盾IP資料
       local res, ip, port = G.Tool.Aliyun.getProxyTcpByDomain("web", 80)

        -- 取IP成功
        if res == 0 then
          print("<color=yellow>res=0</color>")
             toaip="http://"..ip..":"..port
        else
           print("<color=yellow>res~=0</color>")

        end 
    else
      print("<color=yellow>IsInitSuccess~=true</color>")
    end 
  else
      print("<color=yellow>not support Aliyun</color>")

  end
  print("<color=yellow>ToaIO</color>",toaip)
   GetPageData(toaip,
     function(msg,goerror) 
                print("TOA Module")
               
            if(string.find(msg, 'IP:', 1, true)~=nil)then

            else
                if(goerror~=nil)then
                  goerror(); 
                  return
                end                       
            end     
                       
            local startindex=string.find(msg, 'IP:', 1, true)+3;
            local endindex=G.Tool.StringHelper.string.utf8len(msg)+1;           
            local Len=endindex-startindex; 
			      local ip=G.Tool.StringHelper.SubString(msg,startindex,Len);  
            
            --判断IP是否正常
            if(PublicIP=="") then
              print("GetData:"..ip)
              print("iptype:"..G.Tool.StringHelper.GetIPType(ip))
                if((G.Tool.StringHelper.GetIPType(ip)==1) or (G.Tool.StringHelper.GetIPType(ip)==2))then
                   PublicIP=ip
                    print("Public IP Set Already"..ip)
                    if(callback~=nil)then                        
                        callback(); 
                    end
                else
                  --失败往下去继续走
                  if(goerror~=nil)then
                    goerror(); 
                  end  
                end
            else
                if(callback~=nil)then
                    print("Goback"..ip)
                    callback(); 
                end
            end

                  
          
		end,
    function(msg)
       GetPageData("http://example.invalid/ic.asp", function(msg,goerror)             
                        print("Public IP from ip138")

            if(string.find(msg, '[', 1, true)~=nil)then

            else
                if(goerror~=nil)then
                  goerror(); 
                  return
                end                       
            end     
                      
            local startindex=string.find(msg, '[', 1, true)+1; 
            local endindex=string.find(msg, ']', 1, true);           
            local Len=endindex-startindex;           
            local ip=G.Tool.StringHelper.SubString(msg,startindex,Len);  
            --判断IP是否正常
            if(PublicIP=="") then
              print("GetData:"..ip)
              print("iptype:"..G.Tool.StringHelper.GetIPType(ip))
                if((G.Tool.StringHelper.GetIPType(ip)==1) or (G.Tool.StringHelper.GetIPType(ip)==2))then
                   PublicIP=ip
                    print("Public IP Set Already"..ip)
                    if(callback~=nil)then
                        
                        callback(); 
                    end
                else
                  --失败往下去继续走
                  if(goerror~=nil)then
                    goerror(); 
                  end  
                end
            else
                if(callback~=nil)then
                    print("Goback"..ip)
                    callback(); 
                end
            end

                         
    end,
    function(msg)
                     print("Get Public Ip 失敗")
                     if(G.GameConfig.RequestIpCount>5)then

                     PublicIP="0.0.0.0"

                     end

                     if(callback~=nil)then
                        callback(); 
                     end

    
    end);--2 end
		end);--1 end

  
end

function Print(content,color)
    print("<color="..color..">"..tostring(content).."</color>")
end

