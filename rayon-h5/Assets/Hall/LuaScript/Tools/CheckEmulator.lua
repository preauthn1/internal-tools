require("Module")
mmodule(...)
local G = require "Global"

function CheckPermissions(permission, message_err)
	print("############### 檢查權限 " .. permission)
	--StartUp.setTipMessage(message, 0.8);
	
	if CS.UnityEngine.Application.platform ~= CS.UnityEngine.RuntimePlatform.Android then	
		return true
	end
	
	if CS.AndroidPermissionsSDK.s_IsCanUse == true then	
		if CS.AndroidPermissionsSDK.CheckPermissions(permission) == false then
			print("############### 沒有權限 " .. permission)
			
			G.Util.coroutine_call( function()
				local boxResult = G.Util.async_to_sync(CS.UIMessageBox.Show)(CS.MessageBoxStyle.OK, "提示", message_err, "确定");
				if boxResult == CS.MessageBoxResult.OK then
					CS.UnityEngine.Application.Quit()
				end
			end )()			
			
			return false
		else
			print("############### 已經有權限 " .. permission)
			return true
		end
	else
		print("############### 舊機制，不檢查權限")		
		return	true
	end
end

--[[
function CheckPermissions(permission, message, message_err, message_err_cb)
	print("&&&&&&&&&&&&&&&&&&&&&&&&& 檢查權限 " .. permission)
	StartUp.setTipMessage(message, 0.8);
	if CS.AndroidPermissionsSDK.CheckPermissions(permission) == false then
		print("&&&&&&&&&&&&&&&&&&&&&&&&& 沒有權限 確認權限 " .. permission)
		CS.AndroidPermissionsSDK.OnGrantButtonPress
		(
			permission, 
			function(result)
				if result == 0 then
					print("############################ 權限結果 " .. permission .. " 权限取得成功")
					StartUp.setTipMessage("权限取得成功", 0.8);
				
					idx = idx + 1
					CheckAllPermissions()
				elseif result == 1 then
					print("############################ 權限結果 " .. permission .. " 缺少权限资料")
					StartUp.setTipMessage("缺少权限资料", 0.8);
					CallbackError("缺少权限资料")
				elseif result == 2 then
					print("############################ 權限結果 " .. permission .. " 权限取得被拒")
					StartUp.setTipMessage(message_err, 0.8);
					CallbackError(message_err_cb)
				end
			end
		)
	else
		print("############################ 已經有權限 " .. permission .. " 0")
		idx = idx + 1
		CheckAllPermissions()
	end
end


idx = 1
function CheckAllPermissions()
	print("%%%%%%%%%%%%%%%%%%%%%%%%%%%% 確認權限 ", idx)
	
	-- 讀取外部檔案權限
	if idx == 1 then
		CheckPermissions("android.permission.READ_EXTERNAL_STORAGE",
						"确认读写文件权限，以启用读取更新资源功能...",
						"取得读写文件权限失败",
						"读写文件权限未启用，无法读取更新资源")
	-- 讀取手機資料權限
	elseif idx == 2 then
		CheckPermissions("android.permission.READ_PHONE_STATE",
						"确认手机权限，优化资源下载...",
						"取得读手机权限失败",
						"手机权限未启用，无法优化资源下载")
	else
		CallbackSuccess()
	end
end

function ReCheckAllPermissions()
	CheckAllPermissions()
end

CallbackSuccess = nil
CallbackError = nil
function InitCheckPermissions(SccessCallback, ErrorCallback)
	CallbackSuccess = SccessCallback
	CallbackError = ErrorCallback
	idx = 1
	
	if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then	
		print("$$$$$$$$$$$$$ 是否檢查權限 " .. tostring(CS.AndroidPermissionsSDK.s_IsCanUse))
		-- 檢查權限
		if CS.AndroidPermissionsSDK.s_IsCanUse == true then	
			print("@@@@@@@@@@@@@@@@@@@@@ 開始確認權限")
			CheckAllPermissions()
		else
			print("@@@@@@@@@@@@@@@@@@@@@ 舊機制，不檢查權限")		
			CallbackSuccess()
		end
	else
		CallbackSuccess()
	end
end
]]--

--获取裝置内容
function Check()
	print("############################# CheckEmulator")
	
	if CS.UnityEngine.Application.platform == CS.UnityEngine.RuntimePlatform.Android then	
	
		Data = CS.TerproSDK_New.Check()
		Data = string.gsub(Data, "-", ".")
		print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		print(Data)
		
		-- mac address "08:00:27:0e:38:b3" for bluestacks4
		checkData = { "mayOnEmulator:true", "08:00:27:0e:38:b3" }
		for _, name in pairs(checkData) do
			q=string.match(Data, name) 
			if q ~= nil then
				print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name)
				return true
			end		
		end
		
		Data1 = CS.TerproSDK_New.Check1()
		print("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		print(Data1)
		
		-- for old bluestacks(ro.kernel.qemu:0 and CPU)
		-- mac address 12:34:56:78:9a:bc for 夜神模擬器
		q=string.match(Data1, "ro.kernel.qemu:0")
		if q == nil then
			print("not find ro.kernel.qemu:0")
			
			checkData = { "Intel", "intel", "AMD", "amd", "Amd" }
		
			for _, name in pairs(checkData) do
				_, q=string.find(Data, name) 
				if q ~= nil then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name)
					return true
				end		
			end	
		end
		
		checkData = { "12:34:56:78:9a:bc" }
		for _, name in pairs(checkData) do
			q=string.match(Data1, name) 
			if q ~= nil then
				print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name)
				return true
			end		
		end
		
		local width = CS.UnityEngine.Screen.width
		local height = CS.UnityEngine.Screen.height		
		print("widthh " .. width .. " height " .. height)
		
		checkData = { "SM%-G955N", "SM%-N950F" }
		for _, name in pairs(checkData) do
			q=string.match(Data1, name) 
			if q ~= nil then
				print("find " .. name)
			
				local isFind = false
				if width == 1920 and height == 1080 then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name .. " 1920*1080")
					return true
				elseif width == 1440 and height == 900 then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name .. " 1440*900")
					return true
				elseif width == 1280 and height == 720 then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name .. " 1280*720")
					return true
				elseif width == 800 and height == 600 then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name .. " 800*600")
					return true
				elseif width < 800 then
					print("find Emulator !!!!!!!!!!!!!!!!! Check!! " .. name .. " <800")
					return true
				end
			else
				print("not find " .. name)			
			end		
		end			
	end
	
	return false
end

function GetDeviceInfo()
	Data = CS.TerproSDK_New.Check()
	Data = string.gsub(Data,"\n","")

	Data = string.gsub(Data,"最终方法获取IMEI","")
	Data = string.gsub(Data,"最终方法获取MAC地址","")
	Data = string.gsub(Data,"最终方法获取AndroidID","")
	Data = string.gsub(Data,"是否模拟器","")
	Data = string.gsub(Data,"真实 反Hook Proxy代理获取Deviceid level1","")
	Data = string.gsub(Data,"真实 反Hook Proxy代理获取Deviceid level2","")
	Data = string.gsub(Data,"真实 反Hook Proxy代理获取Deviceid","")
	Data = string.gsub(Data,"真实 ITelephonyUtil反Hook 获取DeviceId level1","")
	Data = string.gsub(Data,"真实 ITelephonyUtil反Hook 获取DeviceId","")
	Data = string.gsub(Data,"自定义ServiceManager获取getDeviceId level2","")
	Data = string.gsub(Data,"getCpuInfo"," ")
	
	--[[
	Data = string.gsub(Data,"getMacAddress","1")
	Data = string.gsub(Data,"getAndroidId","2")
	Data = string.gsub(Data,"mayOnEmulator","3")
	Data = string.gsub(Data,"getDeviceIdLevel0","4")
	Data = string.gsub(Data,"getDeviceIdLevel1","5")
	Data = string.gsub(Data,"getDeviceIdLevel2","6")
	Data = string.gsub(Data,"getDeviceIdLeveHl0","7")
	Data = string.gsub(Data,"getDeviceIdLevelH1","8")
	Data = string.gsub(Data,"getDeviceIdLeveHl2","9")
	Data = string.gsub(Data,"getCpuInfo","10")
	Data = string.gsub(Data,"getDeviceId","0")
	Data = string.gsub(Data,"ro.product.cpu.abi","11")
	]]--
	--print("Data\n" .. Data)

	Data1 = CS.TerproSDK_New.Check1()
	Data1 = string.gsub(Data1,"\n","")

	Data1 = string.gsub(Data1,"系统API反射获取序列号","")
	Data1 = string.gsub(Data1,"直接通过 Build Serial","")
	Data1 = string.gsub(Data1,"通过ADB Build Serial","")
	Data1 = string.gsub(Data1,"直接native获取  Serial","")
	Data1 = string.gsub(Data1,"通过系统API获取MAC地址","")
	Data1 = string.gsub(Data1,"Iwifmanager 获取mac level 0","")
	Data1 = string.gsub(Data1,"通过NetworkInterface获取MAC地址","")
	Data1 = string.gsub(Data1,"系统API获取手机型号 （作假）","")
	Data1 = string.gsub(Data1,"通过系统API获取ANDROID_ID","")
	Data1 = string.gsub(Data1,"(XPOSED可以HOOK)","")
	Data1 = string.gsub(Data1,"反射获取系统 ANDROID_IDISettingUtils","")
	Data1 = string.gsub(Data1,"反射获取系统 ANDROID_ID ISettingUtils level2","")
	Data1 = string.gsub(Data1,"native ro.product.manufacturer","")
	Data1 = string.gsub(Data1,"native ro.product.model","")
	Data1 = string.gsub(Data1,"native ro.product.device","")
	Data1 = string.gsub(Data1,"native ro.kernel.qemu","")

	--[[	
	Data1 = string.gsub(Data1,"getSerialNumber","20")
	Data1 = string.gsub(Data1,"getJavaSerialNumber","21")
	Data1 = string.gsub(Data1,"Build.SERIAL","22")
	Data1 = string.gsub(Data1,"getSerialno","23")
	Data1 = string.gsub(Data1,"ro.serialno","24")
	Data1 = string.gsub(Data1,"getMacAddressOS","25")
	Data1 = string.gsub(Data1,"getMacAddressWifi","26")
	Data1 = string.gsub(Data1,"getMacAddressByWlan0","27")
	Data1 = string.gsub(Data1,"getPhoneManufacturer","28")
	Data1 = string.gsub(Data1,"getAndroidId","29")
	Data1 = string.gsub(Data1,"getAndroidProperty","30")
	Data1 = string.gsub(Data1,"ro.product.manufacturer","31")
	Data1 = string.gsub(Data1,"ro.product.model","32")
	Data1 = string.gsub(Data1,"ro.product.device","33")
	Data1 = string.gsub(Data1,"ro.kernel.qemu","34")
	Data1 = string.gsub(Data1,"ro.product.name","35")
	]]--
	--print("Data1\n" .. Data1)
	
	DeviceInfo = Data .. " " .. Data1
	
	DeviceInfo = string.gsub(DeviceInfo,"get","")
	DeviceInfo = string.gsub(DeviceInfo,"  "," ")
		
	return DeviceInfo
end


function GetGeneralDeviceInfo()
	local msg="";
  
        msg=msg.."设备与系统信息:";
  
   
  
         --设备的模型
  
          msg=msg..("设备模型"..tostring(CS.UnityEngine.SystemInfo.deviceModel))..",";
  
          --设备的名称
  
          msg=msg..("设备名称"..tostring(CS.UnityEngine.SystemInfo.deviceName))..",";
 
          --设备的类型
  
          msg=msg..("设备类型（PC电脑，掌上型）"..tostring(CS.UnityEngine.SystemInfo.deviceType))..",";
  
          --系统内存大小
  
          msg=msg..("系统内存大小MB"..tostring(CS.UnityEngine.SystemInfo.systemMemorySize))..",";
  
          --操作系统
  
          msg=msg..("操作系统"..tostring(CS.UnityEngine.SystemInfo.operatingSystem))..",";
  
          --设备的唯一标识符
 
          msg=msg..("设备唯一标识符"..tostring(CS.UnityEngine.SystemInfo.deviceUniqueIdentifier))..",";
 
         --显卡设备标识ID
  
          msg=msg..("显卡ID"..tostring(CS.UnityEngine.SystemInfo.graphicsDeviceID))..",";
  
          --显卡名称
  
          msg=msg..("显卡名称".. tostring(CS.UnityEngine.SystemInfo.graphicsDeviceName))..",";
  
          --显卡类型
  
          msg=msg..("显卡类型"..tostring(CS.UnityEngine.SystemInfo.graphicsDeviceType))..",";
  
          --显卡供应商
  
          msg=msg..("显卡供应商".. tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVendor))..",";
  
          --显卡供应唯一ID
  
          msg=msg..("显卡供应唯一ID".. tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVendorID))..",";
  
          --显卡版本号
  
          msg=msg..("显卡版本号"..tostring(CS.UnityEngine.SystemInfo.graphicsDeviceVersion))..",";
  
          --显卡内存大小
  
          msg=msg..("显存大小MB"..tostring(CS.UnityEngine.SystemInfo.graphicsMemorySize))..",";
  
          --显卡是否支持多线程渲染
  
          msg=msg..("显卡是否支持多线程渲染"..tostring(CS.UnityEngine.SystemInfo.graphicsMultiThreaded))..",";
  
          --支持的渲染目标数量
  
          msg=msg..("支持的渲染目标数量".. tostring(CS.UnityEngine.SystemInfo.supportedRenderTargetCount));

	return msg
end