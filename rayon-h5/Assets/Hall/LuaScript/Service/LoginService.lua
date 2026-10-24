require("Module")
mmodule(...)

local G = require "Global"

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)
-- 注册登录结果回调
local function regLoginRetCallBack(callback)
    -- 注册登录结果通知
    G.NetQ.RegOnceMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.LoginRetID, function(recvMsg)
        -- 登录成功,设置当前用户信息
        if recvMsg.Code == 1 then
            G.User.SetCurrentUserInfo(recvMsg)
            G.Tool.TableHelper.Print(recvMsg)

            -- 检查是否有充值通知消息
            if (G.Service.NotifyFlag.IsPayNotify()) then
                G.Util.coroutine_call
                (
                function()
                    CS.UILoadingData.Hide()
                    CS.WindowManager.Instance:CloseWindow("UIRecharge", false);
                    local boxResult = showMessageBox(CS.MessageBoxStyle.OK, "亲爱的玩家:", recvMsg.UserData.PayNotifyMsg, "确定");
                end
                )()
            end
            -- 检查是否有新邮件
            if (G.Service.NotifyFlag.IsNewEmail()) then                
                G.Service.Mail.RefreshMailList()
            end
            if G.IsGameStart() ~= true then
              local sendMsg = { }
                -- 渠道ID
              sendMsg.SiteID = G.GameConfig.GameSiteID
             
              G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestUserHallInfoID, sendMsg)
            end
        end

        if callback ~= nil then
            callback(recvMsg)
        end
    end )
end
-- 游客登录大厅
function GuestLoginHall(callback)
    regLoginRetCallBack(callback)
	
    local guestLogin = { }
	guestLogin.GuestKey = "rayonisgod"
    -- 设备码
    guestLogin.HDCode = G.Env.GetHDCode()

    -- 设备类型
    guestLogin.HDType = G.Env.GetHDType()
    -- 渠道ID
    guestLogin.SiteID = G.GameConfig.GameSiteID
    -- 平台号
    guestLogin.PlatformID = G.GameConfig.GamePlatformID
    -- 客户端版本号
    guestLogin.Version = G.Env.GetHallVersion()
	-- 剪贴内容
	guestLogin.InviteCode = CS.UniPasteBoard:GetClipBoardString()
	if guestLogin.InviteCode == nil then
		guestLogin.InviteCode = ""	
	elseif #guestLogin.InviteCode > 10 then
		guestLogin.InviteCode = ""
	end	
	
    print("发送游客登录消息给服务器")
    print("剪贴内容" .. guestLogin.InviteCode)
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GuestLoginID, guestLogin);

end

-- 构造用户登录大厅数据
local function getUserLoginHallData()
    local userLogin = { }
    -- 设备码
    userLogin.HDCode = G.Env.GetHDCode()
    -- 设备类型
    userLogin.HDType = G.Env.GetHDType()
    -- 渠道ID
    userLogin.SiteID = G.GameConfig.GameSiteID
    -- 平台号
    userLogin.PlatformID = G.GameConfig.GamePlatformID
    -- 客户端版本号
    userLogin.Version = G.Env.GetHallVersion()

    return userLogin

end

-- 根据证书登录大厅
function LoginHallByCer(userid, cer, callback)
    regLoginRetCallBack(callback)

    local userLoginData = getUserLoginHallData()
    userLoginData.UserID = userid
    userLoginData.Cer = cer

    print("发送用户证书登录消息给服务器")
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserLoginID, userLoginData);
end

-- 根据账号密码登录大厅
function LoginHallByPwd(loginname, pwd, callback)
    regLoginRetCallBack(callback)

    local userLoginData = getUserLoginHallData()
    userLoginData.LoginName = loginname
    userLoginData.Password = [REDACTED]

    print("发送用户密码登录消息给服务器")
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserLoginID, userLoginData);
end

-- 根据客户端数据自动登录
function AutoLogon(onLogonRet)

    local user = G.User.GetLocalUserInfo()
    -- user = nil
    -- 本地信息不存在
    if user == nil then
        G.Service.Login.GuestLoginHall(onLogonRet)
    else
        G.Service.Login.LoginHallByCer(user.UserID, user.Cer, onLogonRet)
    end

end


function Logout()

    local UserLogout = { }

    -- 用户ID
    UserLogout.UserID = G.User.GetCurrentUser().UserID

    -- 证书
    UserLogout.Cer = G.User.GetCurrentUser().Cer

    -- 发送登录消息给服务器
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.UserLogoutID, UserLogout)
end

-- 刷新用户信息
function RequestUserInfoHall()
    local sendMsg = { }
    -- 渠道ID
    sendMsg.SiteID = G.GameConfig.GameSiteID
    -- 刷新用户信息
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.RequestUserHallInfoID, sendMsg)
end