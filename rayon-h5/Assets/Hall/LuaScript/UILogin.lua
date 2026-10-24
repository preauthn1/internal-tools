local RecTran_Login = this.transform:Find("_Login"):GetComponent("RectTransform")
local Btn_AccountLogin = this.transform:Find("_Login/_AccountLogin"):GetComponent("Button")
local Btn_VisitorLogin = this.transform:Find("_Login/_VisitorLogin"):GetComponent("Button")
local Btn_RetrievePW = this.transform:Find("_Login/_RetrievePW"):GetComponent("Button")
local Btn_Customer = this.transform:Find("BtnCustomer"):GetComponent("Button")
local Btn_Website = this.transform:Find("BtnWebsite"):GetComponent("Button")
local Ipt_Password = [REDACTED]"_Login/_Password"):GetComponent("InputField")
local Ipt_Account = this.transform:Find("_Login/_Account"):GetComponent("InputField")
local RecTran_RetrievePW = this.transform:Find("_RetrievePW"):GetComponent("RectTransform")
local Btn_GiveupRetrievePW = this.transform:Find("_RetrievePW/_GiveupRetrievePW"):GetComponent("Button")
local Btn_ComfirmRetrievePW = this.transform:Find("_RetrievePW/_ComfirmRetrievePW"):GetComponent("Button")
local Ipt_RPPhoneNumber = this.transform:Find("_RetrievePW/_RPPhoneNumber"):GetComponent("InputField")
local Ipt_RPCode = this.transform:Find("_RetrievePW/_RPCode"):GetComponent("InputField")
local Ipt_RPNewPW = this.transform:Find("_RetrievePW/_RPNewPW"):GetComponent("InputField")
local Ipt_RPReapetPW = this.transform:Find("_RetrievePW/_RPReapetPW"):GetComponent("InputField")
local Btn_RPGetCode = this.transform:Find("_RetrievePW/_RPGetCode"):GetComponent("Button")
local Txt_RPCountDown = this.transform:Find("_RetrievePW/_RPGetCode/_RPCountDown"):GetComponent("Text")
local RecTran_Customer = this.transform:Find("_Customer"):GetComponent("RectTransform")
local Btn_CustomerClose = this.transform:Find("_Customer/_Close"):GetComponent("Button")
local Btn_CustomerSubmit = this.transform:Find("_Customer/_Submit"):GetComponent("Button")
local Ipt_CustomerContent = this.transform:Find("_Customer/_Content"):GetComponent("InputField")
local Ipt_CustomerEmail = this.transform:Find("_Customer/_Email"):GetComponent("InputField")
local Drop_QuestionType = this.transform:Find("_Customer/_QuestionType"):GetComponent("Dropdown")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---
local G = require "Global"
local gameDefine = require("Com/GameDefine")
local groupName = "UILogin";

local function onLoginRet(receiveMsg)
		
    -- 登入成功
    if receiveMsg.Code == 1 then
	        print("大厅登录成功")
        CS.UILoadingData.Hide()

        -- 关闭当前界面
        CS.WindowManager.Instance:CloseWindow("UILogin", false);

        -- 打开大厅界面
		if CS.WindowManager.Instance:FindWindowNodeByName("UIHall") == nil then
			local openparam = CS.WindowNodeInitParam("UIHall", CS.AppDefine.HallNumber)
			openparam.NodeType = 0
			CS.WindowManager.Instance:OpenWindow(openparam)
		end
	-- 邀请码错误
	elseif receiveMsg.Code == -2 then
		G.Service.NetworkService.CloseHallSocket()
		G.Service.NetworkService.SetNotReConnectHall()
		
		local openparam = CS.WindowNodeInitParam("UIBindingMessageBox", CS.AppDefine.HallNumber)
		openparam.NodeType = CS.BaseNodeType.Propmt;
		openparam.LoadComplatedCallBack = function(windowNode)
			local txt = windowNode.WindowGameObject.transform:Find("UIBody/_ContentArea/_Content"):GetComponent("Text")
			txt.text =  receiveMsg.Message	
		end
		CS.WindowManager.Instance:OpenWindow(openparam)		
	-- 登入失败
	else
        print("大厅登录失败")
        CS.UILoadingData.Hide()

        Btn_AccountLogin.enabled = true
        Btn_VisitorLogin.enabled = true

        CS.UIBubblePrompt.Show(receiveMsg.Message, "UILogin")
    end
end


-- 初始化函数，实例化之后调用一次
function Start()

	if CS.Utility.GetbundleIdentifier() == "com.stage65.test" then
    	Btn_Customer.gameObject:SetActive(false);
    	Btn_Website.gameObject:SetActive(false);
    end  


    -- 游客登录
    Btn_VisitorLogin:AddClickListener
    (
    function()

        Btn_VisitorLogin.enabled = false

        local uiData = CS.LoadingDataUIData()

        uiData.TipsContent = "登录中，请稍后。。。";
        uiData.MaxLastTime = 5; 
        uiData.MaxTimeCloseWindowCallBack =
        (
        function()
            CS.UIBubblePrompt.Show("登录失败", "Login")
            Btn_VisitorLogin.enabled = true
        end
        )

        CS.UILoadingData.Show(uiData)

        G.Service.Login.GuestLoginHall(onLoginRet)
    end
    )

    -- 账户登录
    Btn_AccountLogin:AddClickListener( function()
        local loginname = Ipt_Account.text
        local pwd = [REDACTED]

        -- 检查用户名
        if G.Tool.StringHelper.IsEmpty(loginname) then
            CS.UIBubblePrompt.Show("用户名不能为空", "Login")
            return
        end

        -- 检查密码
        if G.Tool.StringHelper.IsEmpty(pwd) then
            CS.UIBubblePrompt.Show("密码不能为空", "Login")
            return
        end

        Btn_AccountLogin.enabled = false

        local uiData = CS.LoadingDataUIData()
        uiData.TipsContent = "登录中，请稍后。。。";
        uiData.MaxLastTime = 5; 
        uiData.MaxTimeCloseWindowCallBack =
        (
        function()
            CS.UIBubblePrompt.Show("登录失败", "Login")
            Btn_AccountLogin.enabled = true
        end
        )

        CS.UILoadingData.Show(uiData)

        G.Service.Login.LoginHallByPwd(loginname, pwd, onLoginRet)
    end
    )

    -- 找回密码
    Btn_RetrievePW:AddClickListener( function()
        RecTran_Login.gameObject:SetActive(false);
        RecTran_RetrievePW.gameObject:SetActive(true);
    end
    )

    -- 放弃找回
    Btn_GiveupRetrievePW:AddClickListener( function()
        RecTran_Login.gameObject:SetActive(true);
        RecTran_RetrievePW.gameObject:SetActive(false);
    end
    )
	
	local _WebsiteUrl = CS.UnityEngine.PlayerPrefs.GetString("WebsiteUrl", "")
	if _WebsiteUrl == "" then
		Btn_Website.gameObject:SetActive(false)
	end
		
	-- 前往官网
	Btn_Website:AddClickListener(function()
       --CS.UnityEngine.Application.OpenURL("http://example.invalid/")
       CS.UnityEngine.Application.OpenURL(_WebsiteUrl)
	end
	)
	
	-- 聯絡客服
	local opencs=G.Platform.IsShowStartUpCS(G.GameConfig.GamePlatformID)
	Btn_Customer.gameObject:SetActive(opencs)
	Btn_Customer:AddClickListener(function()
        RecTran_RetrievePW.gameObject:SetActive(false);
        RecTran_Login.gameObject:SetActive(false);
		RecTran_Customer.gameObject:SetActive(true);
	end
	)
	
	Btn_CustomerClose:AddClickListener(function()
        RecTran_Login.gameObject:SetActive(true);
		RecTran_Customer.gameObject:SetActive(false);
	end
	)

    -- 发送验证码
    Btn_RPGetCode:AddClickListener( function()
        -- 判断手机号码
        if G.Tool.StringHelper.IsEmpty(Ipt_RPPhoneNumber.text) then
            CS.UIBubblePrompt.Show("请输入手机号码", groupName)
            return
        end

        local sendMsg = { }

        -- 手机号
        sendMsg.Tel = Ipt_RPPhoneNumber.text

        -- 类型
        sendMsg.CodeType = 2

        -- 绑定手机验证码
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeID, sendMsg)
    end
    )

    Btn_ComfirmRetrievePW:AddClickListener( function()
        if G.Tool.StringHelper.IsEmpty(Ipt_RPPhoneNumber.text) then
            CS.UIBubblePrompt.Show("请输入手机号码", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_RPCode.text) then
            CS.UIBubblePrompt.Show("请输入验证码", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_RPNewPW.text) then
            CS.UIBubblePrompt.Show("请输入密码", groupName)
            return
        end

        if Ipt_RPNewPW.text ~= Ipt_RPReapetPW.text then
            CS.UIBubblePrompt.Show("输入密码和确认密码不一致，请验证输入", groupName)
            return
        end

        local sendMsg = { }

        -- 账户
        sendMsg.Tel = Ipt_RPPhoneNumber.text

        -- 验证码
        sendMsg.VCode = Ipt_RPCode.text

        -- 密码
        sendMsg.Password = [REDACTED]

        -- 发送绑定用户消息
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.FindSetPwdByPhoneID, sendMsg)
    end
    )
	
	Btn_CustomerSubmit:AddClickListener( function()
		if G.Tool.StringHelper.IsEmpty(Ipt_CustomerContent.text) then
            CS.UIBubblePrompt.Show("请输入消息内容", groupName)
            return
        end
		if G.Tool.StringHelper.IsEmpty(Ipt_CustomerEmail.text) then
            CS.UIBubblePrompt.Show("请输入邮箱号码", groupName)
            return
        end
		
		local sendMsg = { }
        sendMsg.Msg = Ipt_CustomerContent.text .. "," .. Ipt_CustomerEmail.text
        sendMsg.QuestionType = Drop_QuestionType.value +1
        sendMsg.PlatformID = gameDefine.GamePlatformID 
        sendMsg.AreaID = gameDefine.GameSiteID
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ContactServiceID, sendMsg)
    end
    )
end

-- region 发送手机验证码
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeRetID, function(receiveMsg)

    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("验证码发送成功", groupName);

        Btn_RPGetCode.enabled = false;

        local sequence = CS.DG.Tweening.DOTween.Sequence();

        for i = receiveMsg.CountDown, 1, -1 do
            sequence:Append(Txt_RPCountDown:DOText(string.format("获取验证码(%s)", i), 1));
        end

        sequence:Append(Txt_RPCountDown:DOText("获取验证码", 1))

        sequence:OnComplete( function()
            Btn_RPGetCode.enabled = true;
        end );
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end

end , groupName);

-- region 修改密码返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.FindSetPwdByPhoneRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("找回密码成功", groupName);

        RecTran_RetrievePW.gameObject:SetActive(false);
        RecTran_Login.gameObject:SetActive(true);

        Ipt_Account.text = Ipt_RPPhoneNumber.text;
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- region 聯絡客服結果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ContactServiceRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("消息发送成功", "UIHall");
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);


function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UISendMsgToGM", false);
end


-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end


function EscapeWindow()
    if RecTran_RetrievePW.gameObject.activeSelf then
        RecTran_RetrievePW.gameObject:SetActive(false);
    else
        CS.UnityEngine.Application.Quit()
    end
end