local Tog_Store = this.transform:Find("UIBody/ToggleGroup/_Store"):GetComponent("Toggle")
local Tog_TakeOut = this.transform:Find("UIBody/ToggleGroup/_TakeOut"):GetComponent("Toggle")
local Tog_SetPW = this.transform:Find("UIBody/ToggleGroup/_SetPW"):GetComponent("Toggle")
local Tog_FindPw = this.transform:Find("UIBody/ToggleGroup/_FindPw"):GetComponent("Toggle")
local Btn_CloseUI = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local RecTran_NumberInput = this.transform:Find("UIBody/_NumberInput"):GetComponent("RectTransform")
local Btn_ClearAmount = this.transform:Find("UIBody/_NumberInput/_ClearAmount"):GetComponent("Button")
local Btn_Add1000 = this.transform:Find("UIBody/_NumberInput/_Add1000"):GetComponent("Button")
local Btn_Add500 = this.transform:Find("UIBody/_NumberInput/_Add500"):GetComponent("Button")

local Btn_Add100 = this.transform:Find("UIBody/_NumberInput/_Add100"):GetComponent("Button")
local Btn_Add10 = this.transform:Find("UIBody/_NumberInput/_Add10"):GetComponent("Button")
local Btn_MAX = this.transform:Find("UIBody/_NumberInput/_MAX"):GetComponent("Button")
local RecTran_MoneyInfo = this.transform:Find("UIBody/_MoneyInfo"):GetComponent("RectTransform")
local Txt_Money = this.transform:Find("UIBody/_MoneyInfo/Money/_Money"):GetComponent("Text")
local Txt_StoreAmount = this.transform:Find("UIBody/_MoneyInfo/StoreAmount/_StoreAmount"):GetComponent("Text")
local RecTran_StoreGroup = this.transform:Find("UIBody/_StoreGroup"):GetComponent("RectTransform")
local Btn_SureStore = this.transform:Find("UIBody/_StoreGroup/_SureStore"):GetComponent("Button")
local Ipt_SGStoreAmount = this.transform:Find("UIBody/_StoreGroup/_SGStoreAmount"):GetComponent("InputField")
local RecTran_TakeOutGroup = this.transform:Find("UIBody/_TakeOutGroup"):GetComponent("RectTransform")
local Btn_SureTakeOut = this.transform:Find("UIBody/_TakeOutGroup/_SureTakeOut"):GetComponent("Button")
local Ipt_TakeOutAmount = this.transform:Find("UIBody/_TakeOutGroup/_TakeOutAmount"):GetComponent("InputField")
local RecTran_ReSetPW = this.transform:Find("UIBody/_ReSetPW"):GetComponent("RectTransform")
local Ipt_ResetOldPW = this.transform:Find("UIBody/_ReSetPW/_ResetOldPW"):GetComponent("InputField")
local Ipt_ResetPW = this.transform:Find("UIBody/_ReSetPW/_ResetPW"):GetComponent("InputField")
local Ipt_RestRepeatPW = this.transform:Find("UIBody/_ReSetPW/_RestRepeatPW"):GetComponent("InputField")
local Btn_ComfirmResetPW = this.transform:Find("UIBody/_ReSetPW/_ComfirmResetPW"):GetComponent("Button")
--local Txt_SetPwTips = this.transform:Find("UIBody/_ReSetPW/_SetPwTips"):GetComponent("Text")
local RecTran_FindPW = this.transform:Find("UIBody/_FindPW"):GetComponent("RectTransform")
local RecTran_FindPWNoBinding = this.transform:Find("UIBody/_FindPWNoBinding"):GetComponent("RectTransform")
local Btn_ComfirmFindPW = this.transform:Find("UIBody/_FindPW/_ComfirmFindPW"):GetComponent("Button")
local Btn_ComfirmFindCS = this.transform:Find("UIBody/_FindPWNoBinding/_ComfirmFindCS"):GetComponent("Button")
local Ipt_RepeatFindPW = this.transform:Find("UIBody/_FindPW/_RepeatFindPW"):GetComponent("InputField")
local Ipt_FindPW = this.transform:Find("UIBody/_FindPW/_FindPW"):GetComponent("InputField")
local Btn_FindPWGetCode = this.transform:Find("UIBody/_FindPW/_FindPWGetCode"):GetComponent("Button")
local Txt_GetCodeCountDown = this.transform:Find("UIBody/_FindPW/_FindPWGetCode/_GetCodeCountDown"):GetComponent("Text")
local Ipt_FindPWCode = this.transform:Find("UIBody/_FindPW/_FindPWCode"):GetComponent("InputField")
local Txt_FindPWPhoneNumber = this.transform:Find("UIBody/_FindPW/_FindPWPhoneNumber"):GetComponent("Text")
local RecTran_BankPassWord = [REDACTED]"UIBody/_BankPassWord"):GetComponent("RectTransform")
local Btn_CloseInputBankPW = this.transform:Find("UIBody/_BankPassWord/_CloseInputBankPW"):GetComponent("Button")
local Btn_ConfirmInputBankPW = this.transform:Find("UIBody/_BankPassWord/_ConfirmInputBankPW"):GetComponent("Button")
local Ipt_BankPW = this.transform:Find("UIBody/_BankPassWord/_BankPW"):GetComponent("InputField")
local Txt_BPTips = this.transform:Find("UIBody/_BankPassWord/_BPTips"):GetComponent("Text")
local Img_StoreArrow1 = this.transform:Find("UIBody/_StoreGroup/img_arrow/Image1"):GetComponent("Image")
local Img_StoreArrow2 = this.transform:Find("UIBody/_StoreGroup/img_arrow/Image2"):GetComponent("Image")
local Img_StoreArrow3 = this.transform:Find("UIBody/_StoreGroup/img_arrow/Image3"):GetComponent("Image")
local Img_TakeOutArrow1 = this.transform:Find("UIBody/_TakeOutGroup/img_arrow/Image1"):GetComponent("Image")
local Img_TakeOutArrow2 = this.transform:Find("UIBody/_TakeOutGroup/img_arrow/Image2"):GetComponent("Image")
local Img_TakeOutArrow3 = this.transform:Find("UIBody/_TakeOutGroup/img_arrow/Image3"):GetComponent("Image")

--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local userdata = G.User.GetCurrentUser().UserData

local groupName = "UIBank";

function AddMoneyAmount(amount)

    if Tog_TakeOut.isOn then
        local currentMoneyAmount = tonumber(Ipt_TakeOutAmount.text)

        if currentMoneyAmount == nil then
            currentMoneyAmount = 0
        end

        if amount == -1 then
            currentMoneyAmount = 0;
        elseif amount == 0 then
            currentMoneyAmount = userdata.BankAmount
        else
            currentMoneyAmount = G.Tool.UserHelper.GetSaveCoin(currentMoneyAmount + amount)
        end

        Ipt_TakeOutAmount.text = G.Tool.UserHelper.ToCashDetail(currentMoneyAmount)
    else
        local currentMoneyAmount = tonumber(Ipt_SGStoreAmount.text);

        if currentMoneyAmount == nil then
            currentMoneyAmount = 0
        end

        if amount == -1 then
            currentMoneyAmount = 0;
        elseif amount == 0 then
            currentMoneyAmount = userdata.CashAmount
        else
            currentMoneyAmount = G.Tool.UserHelper.GetSaveCoin(currentMoneyAmount + amount)
        end

        Ipt_SGStoreAmount.text = G.Tool.UserHelper.ToCashDetail(currentMoneyAmount)
    end
end

-- 初始化函数，实例化之后调用一次
function Start()

    Btn_CloseUI:AddClickListener
    (
    function()
        EscapeWindow()
    end
    )

    Tog_Store.onValueChanged:AddListener
    (
    function(isOn)
        if isOn then
            RecTran_NumberInput.gameObject:SetActive(isOn)
            RecTran_MoneyInfo.gameObject:SetActive(isOn)
            Ipt_SGStoreAmount.text='';
        end

        RecTran_StoreGroup.gameObject:SetActive(isOn)
    end
    )

    Tog_TakeOut.onValueChanged:AddListener
    (
    function(isOn)
        if isOn then
            RecTran_NumberInput.gameObject:SetActive(isOn)
            RecTran_MoneyInfo.gameObject:SetActive(isOn)
            Ipt_TakeOutAmount.text='';
        end

        RecTran_TakeOutGroup.gameObject:SetActive(isOn)
    end
    )

    Tog_SetPW.onValueChanged:AddListener
    (
    function(isOn)
        if isOn then
            RecTran_NumberInput.gameObject:SetActive(false)
            RecTran_MoneyInfo.gameObject:SetActive(false)
        end

        RecTran_ReSetPW.gameObject:SetActive(isOn)
    end
    )

    Tog_FindPw.onValueChanged:AddListener
    (
    function(isOn)
        if isOn then
            RecTran_NumberInput.gameObject:SetActive(false)
            RecTran_MoneyInfo.gameObject:SetActive(false)
        end
        if((userdata.IsBindGuest == true)) then
            RecTran_FindPW.gameObject:SetActive(isOn)
        else
           RecTran_FindPWNoBinding.gameObject:SetActive(isOn)
        end
    end
    )


    Btn_SureStore:AddClickListener
    (
    function()
        
        local storeAmount = G.Tool.UserHelper.GetSaveCoin(G.Tool.StringHelper.IsEmpty(Ipt_SGStoreAmount.text)and 0 or Ipt_SGStoreAmount.text )

        if storeAmount <= 0 then
            CS.UIBubblePrompt.Show("存款金额不能为 0 或负数", gameGroup)
        elseif storeAmount > userdata.CashAmount then
            CS.UIBubblePrompt.Show("存款金额不能超过携带金额", gameGroup)
        else
            local sendMsg = { }

            -- 金额
            sendMsg.Amount = storeAmount;

            -- 类型 0 存款 1 取款
            sendMsg.OP = 0

            print("------------存款-------------")
            G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.DepositMoneyID, sendMsg)
        end

    end
    )

    Btn_SureTakeOut:AddClickListener
    (
    function()
    
        local takeOutAmount = G.Tool.UserHelper.GetSaveCoin(G.Tool.StringHelper.IsEmpty(Ipt_TakeOutAmount.text) and 0 or Ipt_TakeOutAmount.text)

        if takeOutAmount <= 0 then
            CS.UIBubblePrompt.Show("取款金额不能为 0 或负数", gameGroup)
        elseif takeOutAmount > userdata.BankAmount then
            CS.UIBubblePrompt.Show("取款金额不能超过存款金额", gameGroup)
        else
            RecTran_BankPassWord.gameObject:SetActive(true);
            Ipt_BankPW.text = "";
        end
    end
    )

    Btn_CloseInputBankPW:AddClickListener
    (
    function()
        RecTran_BankPassWord.gameObject:SetActive(false)
    end
    )

    Btn_ConfirmInputBankPW:AddClickListener
    (
    function()
        RecTran_BankPassWord.gameObject:SetActive(false)

        local sendMsg = { }

        -- 金额
        sendMsg.Amount = G.Tool.UserHelper.GetSaveCoin(Ipt_TakeOutAmount.text)

        -- 密码
        sendMsg.MoneyPassword = [REDACTED]

        -- 类型   0 存款 1 取款
        sendMsg.OP = 1

        print("------------取款-------------")
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.DepositMoneyID, sendMsg)
    end
    )

    Btn_Add10:AddClickListener
    (
    function()
        AddMoneyAmount(10)
    end
    )

    Btn_Add100:AddClickListener
    (
    function()
        AddMoneyAmount(100)
    end
    )
    Btn_Add500:AddClickListener
    (
    function()
        AddMoneyAmount(500)
    end
    )

    Btn_Add1000:AddClickListener
    (
    function()
        AddMoneyAmount(1000)
    end
    )

    Btn_MAX:AddClickListener
    (
    function()
        AddMoneyAmount(0)
    end
    )

    Btn_ClearAmount:AddClickListener
    (
    function()
        AddMoneyAmount(-1)
    end
    )

    Ipt_SGStoreAmount.onEndEdit:AddListener
    (
    function()
        if G.Tool.StringHelper.IsEmpty(Ipt_SGStoreAmount.text) then
            Ipt_SGStoreAmount.text = G.Tool.UserHelper.ToCashDetail(0)
        else
            Ipt_SGStoreAmount.text = G.Tool.UserHelper.ToCashDetail(Ipt_SGStoreAmount.text, false)
        end
    end
    )

    Ipt_TakeOutAmount.onEndEdit:AddListener
    (
    function()
        if G.Tool.StringHelper.IsEmpty(Ipt_TakeOutAmount.text) then
            Ipt_TakeOutAmount.text = G.Tool.UserHelper.ToCashDetail(0)
        else
            Ipt_TakeOutAmount.text = G.Tool.UserHelper.ToCashDetail(Ipt_TakeOutAmount.text, false)
        end
    end
    )
	
	local opencs=G.Platform.IsShowStartUpCS(G.GameConfig.GamePlatformID)
	Btn_ComfirmFindCS.gameObject:SetActive(opencs)	
    Btn_ComfirmFindCS:AddClickListener
     (
    function()
    
     local openparam = CS.WindowNodeInitParam("UISendMsgToGM", CS.AppDefine.HallNumber)
        CS.WindowManager.Instance:OpenWindow(openparam)
          CS.WindowManager.Instance:CloseWindow("UIBank", false);
    end)
    Btn_ComfirmFindPW:AddClickListener
    (
    function()

        if G.Tool.StringHelper.IsEmpty(Ipt_FindPWCode.text) then
            CS.UIBubblePrompt.Show("请输入验证码", groupName)
            return
        end

        if G.Tool.StringHelper.IsEmpty(Ipt_FindPW.text) then
            CS.UIBubblePrompt.Show("请输入密码", groupName)
            return
        end

        if Ipt_FindPW.text ~= Ipt_RepeatFindPW.text then
            CS.UIBubblePrompt.Show("输入密码和确认密码不一致，请验证输入", groupName)
            return
        end

        local sendMsg = { }

        -- 验证码
        sendMsg.VCode = Ipt_FindPWCode.text

        -- 密码
        sendMsg.Password = [REDACTED]

        print("发送消息 FindSetBankPwdByPhoneID")

        -- 发送绑定用户消息
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.FindSetBankPwdByPhoneID, sendMsg)
    end
    )

    Btn_FindPWGetCode:AddClickListener
    (
    function()
        local sendMsg = { }

        -- 手机号
        sendMsg.Tel = userdata.LoginName;

        -- 类型
        sendMsg.CodeType = 3

        --  找回密码验证码
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeID, sendMsg)
    end
    )

    Btn_ComfirmResetPW:AddClickListener
    (
    function()

        if G.Tool.StringHelper.IsEmpty(Ipt_ResetOldPW.text) then
            CS.UIBubblePrompt.Show("请输入旧密码", groupName)
            return
        end
        if G.Tool.StringHelper.IsEmpty(Ipt_ResetPW.text) then
            CS.UIBubblePrompt.Show("请输入新密码", groupName)
            return
        end
        if Ipt_ResetPW.text ~= Ipt_RestRepeatPW.text then
            CS.UIBubblePrompt.Show("新密码和确认密码不一致，请验证输入", groupName)
            return
        end

        local sendMsg = { }

        -- 旧密码
        sendMsg.OldPwd = [REDACTED]

        -- 新密码
        sendMsg.NewPassword = [REDACTED]

        -- 发送绑定用户消息
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ModifyBankPasswordID, sendMsg)
    end
    )
end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
	
	sequenceStoreArrow:Kill()
	sequenceTakeOut:Kill()
end  

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
    Txt_Money.text = G.Tool.UserHelper.ToCashDetail(userdata.CashAmount)
    Txt_StoreAmount.text = G.Tool.UserHelper.ToCashDetail(userdata.BankAmount)
    --Tog_FindPw.gameObject:SetActive(userdata.IsBindGuest == true)
    --Txt_SetPwTips.text = userdata.BankPwdTipMsg
    Txt_BPTips.text = userdata.BankPwdTipMsg

    Txt_FindPWPhoneNumber.text = userdata.LoginName

    --Ipt_SGStoreAmount.text = G.Tool.UserHelper.ToCashDetail(0)
    --Ipt_TakeOutAmount.text = G.Tool.UserHelper.ToCashDetail(0)
	AnimRunArrow()
end

-- region 存取款结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.DepositMoneyRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then

        if receiveMsg.OP == 0 then
            CS.UIBubblePrompt.Show("存款完成。", groupName);
        else
            CS.UIBubblePrompt.Show("取款成功。", groupName);
        end

        userdata.CashAmount = receiveMsg.CurrentMoney
        userdata.BankAmount = receiveMsg.CurrentBank

        RefreshWindowData(nil);
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);


-- region 发送手机验证码
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeRetID, function(receiveMsg)

    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("验证码发送成功。", groupName);

        Btn_FindPWGetCode.enabled = false;

        local sequence = CS.DG.Tweening.DOTween.Sequence();

        for i = receiveMsg.CountDown, 1, -1 do
            sequence:Append(Txt_GetCodeCountDown:DOText(string.format("获取验证码(%s)", i), 1));
        end

        sequence:Append(Txt_GetCodeCountDown:DOText("获取验证码", 1))

        sequence:OnComplete( function()
            Btn_FindPWGetCode.enabled = true;
        end );
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end

end , groupName);
 
 
-- region 修改密码返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ModifyBankPasswordRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("修改密码成功", groupName);
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- region 找回银行密码
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.FindSetBankPwdByPhoneRetID, function(receiveMsg)
  --  if receiveMsg.Code == 1 then
  --      CS.UIBubblePrompt.Show("找回密码成功", groupName);
  --  else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
  --  end
end , groupName);



function EscapeWindow()
	
    if RecTran_BankPassWord.gameObject.activeSelf then
        RecTran_BankPassWord.gameObject:SetActive(false)
    else
        CS.WindowManager.Instance:CloseWindow("UIBank", false);
    end
end

function AnimRunArrow()
	Img_StoreArrow1.gameObject:SetActive(false)
	Img_StoreArrow2.gameObject:SetActive(false)
	Img_StoreArrow3.gameObject:SetActive(false)
	local timer = 0.2
	if sequenceStoreArrow ~= nil then
		sequenceStoreArrow:Kill()
	end
	sequenceStoreArrow = CS.DG.Tweening.DOTween.Sequence()
	sequenceStoreArrow:SetLoops(-1,CS.DG.Tweening.LoopType.Restart)
	sequenceStoreArrow:AppendCallback(function()
		Img_StoreArrow1.gameObject:SetActive(true)
		Img_StoreArrow2.gameObject:SetActive(false)
		Img_StoreArrow3.gameObject:SetActive(false)
	end)
	sequenceStoreArrow:AppendInterval(timer)
	sequenceStoreArrow:AppendCallback(function()
		Img_StoreArrow1.gameObject:SetActive(false)
		Img_StoreArrow2.gameObject:SetActive(true)
		Img_StoreArrow3.gameObject:SetActive(false)
	end)
	sequenceStoreArrow:AppendInterval(timer)
	sequenceStoreArrow:AppendCallback(function()
		Img_StoreArrow1.gameObject:SetActive(false)
		Img_StoreArrow2.gameObject:SetActive(false)
		Img_StoreArrow3.gameObject:SetActive(true)
	end)
	sequenceStoreArrow:AppendInterval(timer)
	-----------------------------------------------------------
	Img_TakeOutArrow1.gameObject:SetActive(false)
	Img_TakeOutArrow2.gameObject:SetActive(false)
	Img_TakeOutArrow3.gameObject:SetActive(false)
	local timer = 0.2
	if sequenceTakeOut ~= nil then
		sequenceTakeOut:Kill()
	end
	sequenceTakeOut = CS.DG.Tweening.DOTween.Sequence()
	sequenceTakeOut:SetLoops(-1,CS.DG.Tweening.LoopType.Restart)
	sequenceTakeOut:AppendCallback(function()
		Img_TakeOutArrow1.gameObject:SetActive(true)
		Img_TakeOutArrow2.gameObject:SetActive(false)
		Img_TakeOutArrow3.gameObject:SetActive(false)
	end)
	sequenceTakeOut:AppendInterval(timer)
	sequenceTakeOut:AppendCallback(function()
		Img_TakeOutArrow1.gameObject:SetActive(false)
		Img_TakeOutArrow2.gameObject:SetActive(true)
		Img_TakeOutArrow3.gameObject:SetActive(false)
	end)
	sequenceTakeOut:AppendInterval(timer)
	sequenceTakeOut:AppendCallback(function()
		Img_TakeOutArrow1.gameObject:SetActive(false)
		Img_TakeOutArrow2.gameObject:SetActive(false)
		Img_TakeOutArrow3.gameObject:SetActive(true)
	end)
	sequenceTakeOut:AppendInterval(timer)
end