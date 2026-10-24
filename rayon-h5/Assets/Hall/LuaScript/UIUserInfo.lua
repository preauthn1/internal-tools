local Btn_Back = this.transform:Find("UIBody/_Back"):GetComponent("Button")

local img_ZFIsBinding = this.transform:Find("UIBody/_Information/_ZFIsBinding"):GetComponent("RectTransform")
local img_UPIsBinding = this.transform:Find("UIBody/_Information/_UPIsBinding"):GetComponent("RectTransform")
local img_Phone = this.transform:Find("UIBody/_Information/_Phone"):GetComponent("RectTransform")
local Btn_ZFIsBinding = this.transform:Find("UIBody/_Information/_ZFIsBinding/_ZF"):GetComponent("Button")
local Btn_UPIsBinding = this.transform:Find("UIBody/_Information/_UPIsBinding/_UP"):GetComponent("Button")
local Btn_ZFUpdate = this.transform:Find("UIBody/_Information/_ZFIsBinding/_ZFUpdate"):GetComponent("Button")
local Btn_UPUpdate = this.transform:Find("UIBody/_Information/_UPIsBinding/_UPUpdate"):GetComponent("Button")

local btn__CopyId = this.transform:Find("UIBody/_Information/ID/_CopyId"):GetComponent("Button")
local btn__SwitchAccount = this.transform:Find("UIBody/_Information/_ButtonGroup/_SwitchAccount"):GetComponent("Button")

local btn_VipInfo = this.transform:Find("UIBody/_Information/_VipInfo"):GetComponent("Button")
local img_VipIcon = this.transform:Find("UIBody/_Information/_VipInfo/_Icon"):GetComponent("Image")



local RecTran_BingdingPhone = this.transform:Find("UIBody/_BingdingPhone"):GetComponent("RectTransform")
local Btn_BPComfirm = this.transform:Find("UIBody/_BingdingPhone/_BPComfirm"):GetComponent("Button")
local Ipt_BPReapetPW = this.transform:Find("UIBody/_BingdingPhone/_BPReapetPW"):GetComponent("InputField")
local Ipt_BPNewPW = this.transform:Find("UIBody/_BingdingPhone/_BPNewPW"):GetComponent("InputField")
local Btn_BPGetCode = this.transform:Find("UIBody/_BingdingPhone/_BPGetCode"):GetComponent("Button")

local Txt_BPCountDown = this.transform:Find("UIBody/_BingdingPhone/_BPGetCode/_BPCountDown"):GetComponent("Text")
local Ipt_BPCode = this.transform:Find("UIBody/_BingdingPhone/_BPCode"):GetComponent("InputField")
local Ipt_BPPhoneNumber = this.transform:Find("UIBody/_BingdingPhone/_BPPhoneNumber"):GetComponent("InputField")
local RecTran_ChangePassword = [REDACTED]"UIBody/_ChangePassword"):GetComponent("RectTransform")
local Txt_CPPhoneNumber = this.transform:Find("UIBody/_ChangePassword/_CPPhoneNumber"):GetComponent("Text")
local Btn_CPConfirm = this.transform:Find("UIBody/_ChangePassword/_CPConfirm"):GetComponent("Button")
local Ipt_CPRepeatPW = this.transform:Find("UIBody/_ChangePassword/_CPRepeatPW"):GetComponent("InputField")
local Ipt_CPNewPW = this.transform:Find("UIBody/_ChangePassword/_CPNewPW"):GetComponent("InputField")
local Ipt_CPOldPW = this.transform:Find("UIBody/_ChangePassword/_CPOldPW"):GetComponent("InputField")
local RecTran_Information = this.transform:Find("UIBody/_Information"):GetComponent("RectTransform")
local Txt_ID = this.transform:Find("UIBody/_Information/ID/_ID"):GetComponent("Text")
local Img_HeadIcon = this.transform:Find("UIBody/_Information/_ChangeHead/_HeadIcon"):GetComponent("Image")
local Img_ImageFrame = this.transform:Find("UIBody/_Information/_ChangeHead/_HeadFrame"):GetComponent("Image")

local Btn_ChangeHead = this.transform:Find("UIBody/_Information/_ChangeHead"):GetComponent("Button")
local Txt_Money = this.transform:Find("UIBody/_Information/OwnMoney/_Money"):GetComponent("Text")
local Txt_StoreMoney = this.transform:Find("UIBody/_Information/StoreMoney/_StoreMoney"):GetComponent("Text")
local Txt_GiftMoney = this.transform:Find("UIBody/_Information/GiftMoney/_GiftMoney"):GetComponent("Text")
local Btn_GoRecharge = this.transform:Find("UIBody/_Information/OwnMoney/_Recharge"):GetComponent("Button")
local Btn_GoBank = this.transform:Find("UIBody/_Information/StoreMoney/_Bank"):GetComponent("Button")
local Btn_GiftGet = this.transform:Find("UIBody/_Information/GiftMoney/_Get"):GetComponent("Button")
local Rect_GiftGet = this.transform:Find("UIBody/_Information/GiftMoney"):GetComponent("RectTransform")

local Btn_BindingAccount = this.transform:Find("UIBody/_Information/_ButtonGroup/_BindingAccount"):GetComponent("Button")
local Btn_ChangePassword = [REDACTED]"UIBody/_Information/_Phone/_ChangePassword"):GetComponent("Button")
local Txt_Phone = this.transform:Find("UIBody/_Information/_Phone/txt_phone"):GetComponent("Text")
local Txt_Sex = this.transform:Find("UIBody/_Information/_Sex"):GetComponent("Text")
local Btn_Close = this.transform:Find("UIBody/_Information/_Close"):GetComponent("Button")
local Btn_CP_Close = this.transform:Find("UIBody/_ChangePassword/_Close"):GetComponent("Button")
local Btn_BP_Close = this.transform:Find("UIBody/_BingdingPhone/_Close"):GetComponent("Button")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

-- 初始化函数，实例化之后调用一次

local G = require "Global";

local groupName = "UIUserInfo";
local user = G.User.GetCurrentUser()
local userData = user.UserData
local selfopen=true;
local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)

-- region 刷新用户个人信息
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.UserHallInfoID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        G.User.SetCurrentUserInfo(receiveMsg)
        user = G.User.GetCurrentUser()
        userData = user.UserData

        -- 刷新成功
        --comment by Bob
        --Txt_BZDAccount.text = userData.Zhifubao
        --Txt_BZDPlayerName.text = userData.RealName

       
        Btn_BindingAccount.gameObject:SetActive(userData.IsBindGuest ~= true)
       
        Btn_ChangePassword.gameObject:SetActive(userData.IsBindGuest)
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);



-- region 修改密码返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ModifyPasswordRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("密码修改成功", groupName);

        RecTran_ChangePassword.gameObject:SetActive(false);
        RecTran_Information.gameObject:SetActive(true);
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- region 发送手机验证码
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeRetID, function(receiveMsg)

    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("验证码发送成功", groupName);

        Btn_BPGetCode.enabled = false;

        local sequence = CS.DG.Tweening.DOTween.Sequence();

        for i = receiveMsg.CountDown, 1, -1 do
            sequence:Append(Txt_BPCountDown:DOText(string.format("获取验证码(%s)", i), 1));
        end

        sequence:Append(Txt_BPCountDown:DOText("获取验证码", 1))

        sequence:OnComplete( function()
            Btn_BPGetCode.enabled = true;
        end );
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end

end , groupName);


-- region 绑定手机
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.BindGuestAccountRetID, function(receiveMsg)

    if receiveMsg.Code == 1 then
        CS.UILoadingData.Show("绑定成功",2);
        --EscapeWindow()
        --CS.UIBubblePrompt.Show("绑定成功", groupName);
        RecTran_BingdingPhone.gameObject:SetActive(false)
        RecTran_Information.gameObject:SetActive(true);
        G.Service.Login.RequestUserInfoHall();
        local sequence = CS.DG.Tweening.DOTween.Sequence()
            sequence:AppendInterval(2)
            sequence:AppendCallback(function()
            userdata = G.User.GetCurrentUser().UserData
            RefreshWindowData()
               end) 
        
        
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end

end , groupName);

function Start()

	-- 送審階段不開啟 Add By Rayon
	if G.GameConfig.IsSubmissionCtx() == true then
		local _Btn = this.transform:Find("UIBody/_Information/_ButtonGroup")
		if _Btn ~= nil then
			_Btn.gameObject:SetActive(false)
		end
		
		local _StoreMoney = this.transform:Find("UIBody/_Information/Number/StoreMoney/")
		if _StoreMoney ~= nil then
			_StoreMoney.gameObject:SetActive(false)
		end		
	end

    Btn_Back:AddClickListener
    (
    function()
        EscapeWindow();
    end
    )
      btn__CopyId:AddClickListener
    (
    function()
        CS.Utility.CopyText(user.UserID)
        CS.UIBubblePrompt.Show("复制成功", groupName);
    end
    )
      btn__SwitchAccount:AddClickListener
    (
    function()
        
        EscapeWindow()

        G.Tool.Sound.CloseBGMSound()
        CS.WindowManager.Instance:CloseWindow("UIHall", false);
         G.ClearLoginData()
        local openparam = CS.WindowNodeInitParam("UILogin", CS.AppDefine.HallNumber)
        openparam.NodeType = 0
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )

	
	btn_VipInfo:AddClickListener
    (
    function()
        
        EscapeWindow()

        local openparam = CS.WindowNodeInitParam("UIVipInfo", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )
	-- 頭相框是0表示是舊機制,關掉VIP按鈕
	if userData.AvatarFrame == 0 then
		btn_VipInfo.gameObject:SetActive(false)
	end
	
    -- 绑定手机获取验证码
    Btn_BPGetCode:AddClickListener
    (
    function()

        -- 判断手机号码
        if G.Tool.StringHelper.IsEmpty(Ipt_BPPhoneNumber.text) then
            CS.UIBubblePrompt.Show("请输入手机号码", groupName)
            return
        end

        local sendMsg = { }

        -- 手机号
        sendMsg.Tel = Ipt_BPPhoneNumber.text

        -- 类型
        sendMsg.CodeType = 1

        -- 绑定手机验证码
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeID, sendMsg)

    end
    )
    Btn_GoRecharge:AddClickListener
    (
    function()
        
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
        local openparam = CS.WindowNodeInitParam("UIRecharge", CS.AppDefine.HallNumber)
         CS.WindowManager.Instance:OpenWindow(openparam)

    end)
     Btn_GoBank:AddClickListener
    (
    function()
        
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
        local openparam = CS.WindowNodeInitParam("UIBank", CS.AppDefine.HallNumber)
                CS.WindowManager.Instance:OpenWindow(openparam)

    end)
	
	Btn_GiftGet:AddClickListener
    (
    function()
        
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
        local openparam = CS.WindowNodeInitParam("UIGiftCode", CS.AppDefine.HallNumber)
                CS.WindowManager.Instance:OpenWindow(openparam)

    end)	

	if G.User.payinfo.UIClub2AgentMode6 == false then
		Rect_GiftGet.gameObject:SetActive(false)
	end
	
    -- 绑定支付寶按钮
    Btn_ZFUpdate:AddClickListener
    (
    function()
       if userData.IsBindGuest  ~= true then
       G.Util.coroutine_call(function()
             local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "绑定支付宝需要先绑定用户信息，确认绑定吗？", "绑定", "取消");
             if boxResult == CS.MessageBoxResult.OK then
                RecTran_BingdingPhone.gameObject:SetActive(true)
                RecTran_Information.gameObject:SetActive(false)
             end
        end)() 

       else
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);

        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = 1;
        CS.WindowManager.Instance:OpenWindow(openparam)

      end
    end)
      -- 绑定銀聯按钮
    Btn_UPUpdate:AddClickListener
    (
    function()
       if userData.IsBindGuest ~= true then
       G.Util.coroutine_call(function()
                   local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "绑定银联需要先绑定用户信息，确认绑定吗？", "绑定", "取消");
                     if boxResult == CS.MessageBoxResult.OK then
                        RecTran_BingdingPhone.gameObject:SetActive(true)
                        RecTran_Information.gameObject:SetActive(false)
                     end
                end)()
       else

         CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);

        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = 2;
        CS.WindowManager.Instance:OpenWindow(openparam)
      end
    end)

    -- 绑定支付寶按钮
    Btn_ZFIsBinding:AddClickListener
    (
    function()
       if userData.IsBindGuest  ~= true then
       G.Util.coroutine_call(function()
             local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "绑定支付宝需要先绑定用户信息，确认绑定吗？", "绑定", "取消");
             if boxResult == CS.MessageBoxResult.OK then
                RecTran_BingdingPhone.gameObject:SetActive(true)
                RecTran_Information.gameObject:SetActive(false)
             end
        end)() 

       else
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);

        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = 3;
        CS.WindowManager.Instance:OpenWindow(openparam)

      end
    end)
      -- 绑定銀聯按钮
    Btn_UPIsBinding:AddClickListener
    (
    function()
       if userData.IsBindGuest ~= true then
       G.Util.coroutine_call(function()
                   local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancel, "", "绑定银联需要先绑定用户信息，确认绑定吗？", "绑定", "取消");
                     if boxResult == CS.MessageBoxResult.OK then
                        RecTran_BingdingPhone.gameObject:SetActive(true)
                        RecTran_Information.gameObject:SetActive(false)
                     end
                end)()
       else

         CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);

        local openparam = CS.WindowNodeInitParam("UIExchange", CS.AppDefine.HallNumber)
        openparam.WindowData = 4;
        CS.WindowManager.Instance:OpenWindow(openparam)
      end
    end)


    -- 绑定手机按钮
    Btn_BPComfirm:AddClickListener
    (
    function()

        if G.Tool.StringHelper.IsEmpty(Ipt_BPPhoneNumber.text) then
            CS.UIBubblePrompt.Show("请输入手机号码", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_BPCode.text) then
            CS.UIBubblePrompt.Show("请输入验证码", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_BPNewPW.text) then
            CS.UIBubblePrompt.Show("请输入密码", groupName)
            return
        end

        if Ipt_BPNewPW.text ~= Ipt_BPReapetPW.text then
            CS.UIBubblePrompt.Show("输入密码和确认密码不一致，请验证输入", groupName)
            return
        end

        local sendMsg = { }

        -- 账户
        sendMsg.Tel = Ipt_BPPhoneNumber.text

        -- 验证码
        sendMsg.VCode = Ipt_BPCode.text

        -- 密码
        sendMsg.Password = [REDACTED]

        -- 发送绑定用户消息
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.BindGuestAccountID, sendMsg)
    end
    )

    -- 修改密码按钮
    Btn_CPConfirm:AddClickListener
    (
    function()

        if G.Tool.StringHelper.IsEmpty(Ipt_CPOldPW.text) then
            CS.UIBubblePrompt.Show("旧密码不能为空", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_CPNewPW.text) then
            CS.UIBubblePrompt.Show("请输入新密码", groupName)
            return
        end

        if Ipt_CPNewPW.text ~= Ipt_CPRepeatPW.text then
            CS.UIBubblePrompt.Show("新密码和确认密码不一致，请验证输入", groupName)
            return
        end

        local sendMsg = { }

        -- 手机号
        sendMsg.OldPwd = [REDACTED]

        -- 类型
        sendMsg.NewPassword = [REDACTED]

        -- 发送修改密码消息
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ModifyPasswordID, sendMsg)
    end
    )

    Btn_ChangeHead:AddClickListener
    (
    function()
        --local openparam = CS.WindowNodeInitParam("UIChooseHeadIcon", CS.AppDefine.HallNumber)
        local openparam = CS.WindowNodeInitParam("UIChooseHead", CS.AppDefine.HallNumber)

        openparam.WindowCloseCallBack =
        (
        function()
            Img_HeadIcon:ResetSpriteByName(CS.AppDefine.HallNumber, tostring(userData.HeadID))
			G.Service.Vip.SetFrame(Img_ImageFrame, userData.AvatarFrame)
            Txt_Sex.text = G.CaseWhen(userData.Sex == 1, "女", "男")
        end
        )

        CS.WindowManager.Instance:OpenWindow(openparam)
    end
    )

    Btn_BindingAccount:AddClickListener
    (
    function()
        RecTran_BingdingPhone.gameObject:SetActive(true)
        RecTran_Information.gameObject:SetActive(false)
    end
    )

    

    Btn_ChangePassword:[REDACTED]
    (
    function()
		Txt_CPPhoneNumber.text = userData.LoginName;
        Ipt_CPOldPW.text = "";
        Ipt_CPNewPW.text = "";
        Ipt_CPRepeatPW.text = "";

        RecTran_ChangePassword.gameObject:SetActive(true)
        RecTran_Information.gameObject:SetActive(false)
    end
    )
     Btn_BP_Close:AddClickListener
    (
    function()
        EscapeWindow();
        --CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
    end
    )
     Btn_CP_Close:AddClickListener
    (
    function()
        EscapeWindow();
        --CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
    end
    )

    Btn_Close:AddClickListener
    (
    function()
        EscapeWindow();
        --CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
    end
    )
end
 
-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)

    img_VipIcon:ResetSpriteByName(CS.AppDefine.HallNumber, "vip" .. tostring(userData.UserLevel))
	G.Service.Vip.SetFrame(Img_ImageFrame, userData.AvatarFrame)

    Img_HeadIcon:ResetSpriteByName(CS.AppDefine.HallNumber, tostring(userData.HeadID))
    Txt_ID.text = tostring(user.UserID);

    Txt_Money.text = G.Tool.UserHelper.ToCashDetail(userData.CashAmount);
    Txt_StoreMoney.text = G.Tool.UserHelper.ToCashDetail(userData.BankAmount);
    Txt_GiftMoney.text = G.Tool.UserHelper.ToCashDetail(userData.CouponAmount);

    Btn_BindingAccount.gameObject:SetActive(userData.IsBindGuest ~= true)
    btn__SwitchAccount.gameObject:SetActive(userData.IsBindGuest == true)
    Btn_ChangePassword.gameObject:SetActive(userData.IsBindGuest)
   
    Txt_Sex.text = G.CaseWhen(userData.Sex == 1, "女", "男")
    if(userData.IsBindZhifubao) then
         img_ZFIsBinding.transform:Find("img_bindingnotyet").gameObject:SetActive(false)
         Btn_ZFUpdate.gameObject:SetActive(true)
		 Btn_ZFIsBinding.gameObject:SetActive(false)
         img_ZFIsBinding.transform:Find("txt_binding"):GetComponent("Text").text=userData.Zhifubao;
         img_ZFIsBinding.transform:Find("txt_binding").gameObject:SetActive(true)
    end
    if(userData.IsBindBankCard) then
         img_UPIsBinding.transform:Find("img_bindingnotyet").gameObject:SetActive(false)
         Btn_UPUpdate.gameObject:SetActive(true)
		 Btn_UPIsBinding.gameObject:SetActive(false)
         img_UPIsBinding.transform:Find("txt_binding"):GetComponent("Text").text=userData.BankCardNum;
         img_UPIsBinding.transform:Find("txt_binding").gameObject:SetActive(true)
    end
     if(userData.IsBindGuest) then
        img_Phone.transform:Find("img_bindingnotyet").gameObject:SetActive(false);
        Txt_Phone.gameObject:SetActive(true)
        Txt_Phone.text=userData.LoginName;
    end

    
    -- 绑定账号
    if windowData ~= nil and windowData == 1 then
        selfopen=false
        RecTran_BingdingPhone.gameObject:SetActive(true)
        RecTran_Information.gameObject:SetActive(false)
    end

	local isAnd = G.Platform.getExchangeOperationAnd()
	local result = G.User.CheckOpenData(isAnd, userData.ShowConvertLimit, userData.PlaytimeConvert, userData.ShowConvertBindPhone)
	
	-- 支付绑绑定
	Btn_ZFIsBinding.gameObject:SetActive(result)
	if userData.IsBindZhifubao == true then
		Btn_ZFIsBinding.gameObject:SetActive(false)
	end
	
	-- 银联绑定
	Btn_UPIsBinding.gameObject:SetActive(result)
	if(userData.IsBindBankCard) then
		Btn_UPIsBinding.gameObject:SetActive(false)
	end
	
	-- 没绑定手机就不显示支付版和银联
	if G.Platform.getRoleBindZfiAndUp() == true then
		-- 绑定支付寶
		img_ZFIsBinding.gameObject:SetActive(result)
		
		-- 绑定銀聯
		img_UPIsBinding.gameObject:SetActive(result)
	end	
	
	-- 986银行卡不显示修改按钮
	if G.Platform.GetBranch() == 1 then
		Btn_UPUpdate.gameObject:SetActive(false)
	end	
end


function EscapeWindow()
  
    if RecTran_BingdingPhone.gameObject.activeSelf then
        if(selfopen) then
            RecTran_BingdingPhone.gameObject:SetActive(false)
            RecTran_Information.gameObject:SetActive(true)
        else
            CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
        end

    elseif RecTran_ChangePassword.gameObject.activeSelf then
        RecTran_ChangePassword.gameObject:SetActive(false)
        RecTran_Information.gameObject:SetActive(true)
    else
        CS.WindowManager.Instance:CloseWindow("UIUserInfo", false);
    end
end