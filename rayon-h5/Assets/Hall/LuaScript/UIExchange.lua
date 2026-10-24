local RecTran_Bingding = this.transform:Find("_Bingding"):GetComponent("RectTransform")
--local Btn_CancelBingding = this.transform:Find("_Bingding/_CancelBingding"):GetComponent("Button")   --comment by Bob 
--local Btn_Bingding = this.transform:Find("_Bingding/_Bingding"):GetComponent("Button")               --comment by Bob 

local img_ImageTitle1 = this.transform:Find("_Bingding/Background/ImageTitle1"):GetComponent("Image")
local img_ImageTitle2 = this.transform:Find("_Bingding/Background/ImageTitle2"):GetComponent("Image")
local img_ImageTitle3 = this.transform:Find("_Bingding/Background/ImageTitle3"):GetComponent("Image")
local img_ImageTitle4 = this.transform:Find("_Bingding/Background/ImageTitle4"):GetComponent("Image")


local Btn_CloseBingding = this.transform:Find("_Bingding/ButtonClose"):GetComponent("Button")
local RecTran_Exchange = this.transform:Find("_Exchange"):GetComponent("RectTransform")
local Txt_Tips = this.transform:Find("_ExchangeFAQ/_Tips"):GetComponent("Text")
local Txt_Money = this.transform:Find("_Exchange/Moneybg/_Money"):GetComponent("Text")
local Txt_SafeMoney = this.transform:Find("_Exchange/Safebg/_Money"):GetComponent("Text")
local Txt_Account = this.transform:Find("_Exchange/infobg/_Account"):GetComponent("Text")
local Ipt_ExchangeNumber = this.transform:Find("_Exchange/_ExchangeNumber"):GetComponent("InputField")
local Btn_Exchange = this.transform:Find("_Exchange/_Exchange"):GetComponent("Button")
local RecTran_ExchangeSucess = this.transform:Find("_ExchangeSucess"):GetComponent("RectTransform")
local Txt_SucessTip = this.transform:Find("_ExchangeSucess/_SucessTip"):GetComponent("Text")
local Btn_CloseSucessWin = this.transform:Find("_ExchangeSucess/_CloseSucessWin"):GetComponent("Button")
local Btn_plus50 = this.transform:Find("_Exchange/_plus50"):GetComponent("Button")                          -- add by Bob
local Btn_plus500 = this.transform:Find("_Exchange/_plus500"):GetComponent("Button")                        -- add by Bob
local Btn_plus5000 = this.transform:Find("_Exchange/_plus5000"):GetComponent("Button")                      -- add by Bob
local Btn_plusmax = this.transform:Find("_Exchange/_plusmax"):GetComponent("Button")                        -- add by Bob
local RecTran_Text = this.transform:Find("_Exchange/_FeeD"):GetComponent("RectTransform")
local txt_Text = this.transform:Find("_Exchange/_FeeD/Text"):GetComponent("Text")                           -- add by Bob
local Btn_clear = this.transform:Find("_Exchange/_clear"):GetComponent("Button")                            -- add by Bob
local txt__Note = this.transform:Find("_Exchange/_Note"):GetComponent("Text")


--綁定銀聯卡頁面
local RecTran_BindingUP = this.transform:Find("_BindingUP"):GetComponent("RectTransform")                   -- add by Bob
local Btn_BUComfirm = this.transform:Find("_BindingUP/_BUComfirm"):GetComponent("Button")                   -- add by Bob
local RecTran_NoBindingUP = this.transform:Find("_NoBindingUP"):GetComponent("RectTransform")                 -- add by Bob
local Btn_NoBindingUP = this.transform:Find("_NoBindingUP/_Binding"):GetComponent("Button")                 -- add by Bob

local Ipt_BUAccount = this.transform:Find("_BindingUP/_BUAccount"):GetComponent("InputField")               -- add by Bob 
local Ipt_BUPlayerName  = this.transform:Find("_BindingUP/_BUPlayerName"):GetComponent("InputField")        -- add by Bob 
local Btn_BUGetCode = this.transform:Find("_BindingUP/_BUGetCode"):GetComponent("Button")                   -- add by Bob 
local Txt_BUCountDown = this.transform:Find("_BindingUP/_BUGetCode/_BUCountDown"):GetComponent("Text")      -- add by Bob 
local Ipt_BUCode = this.transform:Find("_BindingUP/_BUCode"):GetComponent("InputField")                     -- add by Bob 
--local Ipt_BUBank = this.transform:Find("_BindingUP/_BUPlayerBank/_Dropdown"):GetComponent("Dropdown")       -- add by Bob 

local img_Payment_UP = this.transform:Find("_Exchange/Payment_UP"):GetComponent("Image")                    -- add by Bob
local img_Payment_ZF = this.transform:Find("_Exchange/Payment_ZF"):GetComponent("Image")                    -- add by Bob



--綁定支付寶頁面
local RecTran_BindingZF = this.transform:Find("_BindingZF"):GetComponent("RectTransform")                   -- add by Bob
local Btn_BZComfirm = this.transform:Find("_BindingZF/_BZComfirm"):GetComponent("Button")                   -- add by Bob
local RecTran_NoBindingZF = this.transform:Find("_NoBindingZF"):GetComponent("RectTransform")               -- add by Bob
local Btn_NoBindingZF = this.transform:Find("_NoBindingZF/_Binding"):GetComponent("Button")                 -- add by Bob

local Ipt_BZAccount = this.transform:Find("_BindingZF/_BZAccount"):GetComponent("InputField")               -- add by Bob 
local Ipt_BZPlayerName  = this.transform:Find("_BindingZF/_BZPlayerName"):GetComponent("InputField")        -- add by Bob 
local Btn_BZGetCode = this.transform:Find("_BindingZF/_BZGetCode"):GetComponent("Button")                   -- add by Bob 
local Txt_BZCountDown = this.transform:Find("_BindingZF/_BZGetCode/_BZCountDown"):GetComponent("Text")      -- add by Bob 
local Ipt_BZCode = this.transform:Find("_BindingZF/_BZCode"):GetComponent("InputField")                     -- add by Bob 


 
--選擇付款頁面-銀聯卡資訊
local Txt_BUDAccount = this.transform:Find("_Bingding/UnionPayContent/_BUDPlayerAccount"):GetComponent("Text")              -- add by Bob
local Txt_BUDPlayerName = this.transform:Find("_Bingding/UnionPayContent/_BUDPlayerName"):GetComponent("Text")              -- add by Bob
local Txt_BUDPlayerBank = this.transform:Find("_Bingding/UnionPayContent/_BUDPlayerBank"):GetComponent("Text")              -- add by Bob
local RecTran_UnionPayContent = this.transform:Find("_Bingding/UnionPayContent"):GetComponent("RectTransform")              -- add by Bob

--選擇付款頁面-支付寶資訊
local Txt_BZDAccount = this.transform:Find("_Bingding/ZifuContent/_BZDAccount"):GetComponent("Text")                    -- add by Bob
local Txt_BZDPlayerName = this.transform:Find("_Bingding/ZifuContent/_BZDPlayerName"):GetComponent("Text")              -- add by Bob
local RecTran_ZifuContent = this.transform:Find("_Bingding/ZifuContent"):GetComponent("RectTransform")                  -- add by Bob




--選擇付款頁面-選擇與修改按鈕
local Btn_ZF = this.transform:Find("_Bingding/_ZF"):GetComponent("Button")                                  -- add by Bob 
local Btn_ZF_1 = this.transform:Find("_Bingding/_ZF_1"):GetComponent("Button")                              -- add by Bob 
local Btn_UP = this.transform:Find("_Bingding/_UP"):GetComponent("Button")                                  -- add by Bob 
local Btn_UP_1 = this.transform:Find("_Bingding/_UP_1"):GetComponent("Button")                              -- add by Bob 
local Btn_ZFUpdate = this.transform:Find("_Exchange/_ZFUpdate"):GetComponent("Button")                      -- add by Bob 
local Btn_UPUpdate = this.transform:Find("_Exchange/_UPUpdate"):GetComponent("Button")                      -- add by Bob 
local Btn_Log = this.transform:Find("_Bingding/_Log"):GetComponent("Button")                                -- add by Bob 



local Img_Exchange = this.transform:Find("_Bingding/_Log"):GetComponent("Image")
local Img_ExchangeMask = this.transform:Find("_Bingding/_Log/selected"):GetComponent("Image")
local Img_ExchangeFont2 = this.transform:Find("_Bingding/_Log/ImageFont2"):GetComponent("Image")
local Img_ExchangeMask2 = this.transform:Find("_Bingding/_Log/ImageMask2"):GetComponent("Image")
local Img_ExchangeFont3 = this.transform:Find("_Bingding/_Log/ImageFont3"):GetComponent("Image")
local Img_ExchangeMask3 = this.transform:Find("_Bingding/_Log/ImageMask3"):GetComponent("Image")                       
local Img_ExchangeFont4 = this.transform:Find("_Bingding/_Log/ImageFont4"):GetComponent("Image")
local Img_ExchangeMask4 = this.transform:Find("_Bingding/_Log/ImageMask4"):GetComponent("Image")


--兌換紀錄頁面
local recTran__ExchangeLog = this.transform:Find("_ExchangeLog"):GetComponent("RectTransform")                                      -- add by Bob 
                                      
local RecTran_Item = this.transform:Find("_ExchangeLog/LogList/Viewport/Content/_Item"):GetComponent("RectTransform")               -- add by Bob
local RecTran_Item1 = this.transform:Find("_ExchangeLog/LogList/Viewport/Content/_Item1"):GetComponent("RectTransform")               -- add by Bob
--local Btn__CloseWin = this.transform:Find("_ExchangeLog/Top/_CloseWin"):GetComponent("Button")                                      -- add by Bob
local btn__commit = this.transform:Find("_ExchangeLog/SearchBar/_commit"):GetComponent("Button")                                    -- add by Bob
local dd_duration = this.transform:Find("_ExchangeLog/SearchBar/_duration"):GetComponent("Dropdown")                                -- add by Bob
local dd_type = this.transform:Find("_ExchangeLog/SearchBar/_type"):GetComponent("Dropdown")

local recTran__ExchangeFAQ = this.transform:Find("_ExchangeFAQ"):GetComponent("RectTransform")
local btn__CloseWin = this.transform:Find("_ExchangeFAQ/_CloseWin"):GetComponent("Button")
local img_FAQ1 = this.transform:Find("_ExchangeFAQ/img_title1"):GetComponent("Image")
local img_FAQ2 = this.transform:Find("_ExchangeFAQ/img_title2"):GetComponent("Image")
local img_FAQ3 = this.transform:Find("_ExchangeFAQ/img_title3"):GetComponent("Image")
local img_FAQ4 = this.transform:Find("_ExchangeFAQ/img_title4"):GetComponent("Image")

local btn__FAQ = this.transform:Find("_Bingding/_FAQ"):GetComponent("Button")
local img_ImageFAQ1 = this.transform:Find("_Bingding/_FAQ/ImageState1"):GetComponent("Image")
local img_ImageFAQ2 = this.transform:Find("_Bingding/_FAQ/ImageState2"):GetComponent("Image")
local img_ImageFAQ3 = this.transform:Find("_Bingding/_FAQ/ImageState3"):GetComponent("Image")
local img_ImageFAQ4 = this.transform:Find("_Bingding/_FAQ/ImageState4"):GetComponent("Image")

--
local recTran_BankSelect = this.transform:Find("_BankSelect"):GetComponent("RectTransform")
local recTran_Grid = this.transform:Find("_BankSelect/AgentList/Grid"):GetComponent("RectTransform")
local btn__AgentItem = this.transform:Find("_BankSelect/_BankItem"):GetComponent("Button")
local Inp_BUBank = this.transform:Find("_BindingUP/_BUPlayerBank"):GetComponent("InputField")
local btn__BUGetBank = this.transform:Find("_BindingUP/_BUGetBank"):GetComponent("Button")
local btn__CloseBankWin = this.transform:Find("_BankSelect/_CloseBankWin"):GetComponent("Button")



 local img_immediately = this.transform:Find("_Bingding/_UP/immediately"):GetComponent("Image")
 local img_immediately_1 = this.transform:Find("_Bingding/_UP_1/immediately"):GetComponent("Image")
                   








--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
 
local userdata = G.User.GetCurrentUser().UserData
local groupName = "UIExchange";
local lcBankList={}
local exchangeitems={}
local sta ={}
local btncol={}
local prevbank=nil;
local converttype=0
local IsRejectType=0
local changeText = ""

local minexchangevalue=0 --最小兌換金額
local magnification=2 --兌換倍率

local yield_return = G.Util.async_to_sync
(
function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end
)

local showMessageBox = G.Util.async_to_sync(CS.UIMessageBox.Show)


-- region 绑定支付宝返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.BindZhifubaoRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
            CS.UILoadingData.Show("支付宝绑定成功",2);         
            
            RecTran_BindingZF.gameObject:SetActive(false);
                    
            G.Service.Login.RequestUserInfoHall()
         
            local sequence = CS.DG.Tweening.DOTween.Sequence()
            sequence:AppendInterval(2)
            sequence:AppendCallback(function()
            userdata = G.User.GetCurrentUser().UserData
            ToExchangebyZF();
            RecTran_Exchange.gameObject:SetActive(true);  
               end)
            
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);
-- add by Bob
-- region 绑定銀聯卡返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.BindBankCardRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        CS.UILoadingData.Show("银联卡绑定成功",2);    
        
        RecTran_BindingUP.gameObject:SetActive(false);
        G.Service.Login.RequestUserInfoHall()
        local sequence = CS.DG.Tweening.DOTween.Sequence()
            sequence:AppendInterval(2)
            sequence:AppendCallback(function()
            userdata = G.User.GetCurrentUser().UserData
            ToExchangebyUP();
            RecTran_Exchange.gameObject:SetActive(true)  
        end)
      
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName);

-- region 发送手机验证码
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeRetID, function(receiveMsg)

    if receiveMsg.Code == 1 then
        CS.UIBubblePrompt.Show("验证码发送成功", groupName);
        if(RecTran_BindingZF.gameObject.activeSelf) then
            Btn_BZGetCode.enabled = false;

            local sequence = CS.DG.Tweening.DOTween.Sequence();

            for i = receiveMsg.CountDown, 1, -1 do
                sequence:Append(Txt_BZCountDown:DOText(string.format("获取验证码(%s)", i), 1));
            end

            sequence:Append(Txt_BZCountDown:DOText("获取验证码", 1))

            sequence:OnComplete( function()
                Btn_BZGetCode.enabled = true;
            end );
        end
        if(RecTran_BindingUP.gameObject.activeSelf) then
            Btn_BUGetCode.enabled = false;

            local sequence = CS.DG.Tweening.DOTween.Sequence();

            for i = receiveMsg.CountDown, 1, -1 do
                sequence:Append(Txt_BUCountDown:DOText(string.format("获取验证码(%s)", i), 1));
            end

            sequence:Append(Txt_BUCountDown:DOText("获取验证码", 1))

            sequence:OnComplete( function()
                Btn_BUGetCode.enabled = true;
            end );
        end

    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end

end , groupName);
--  當交易紀錄返回结果
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetConvertLogRetID, function(receiveMsg)
    local logList =  receiveMsg.ConvertLogList
    DisplayRecord(logList);

end  , groupName)


-- 初始化函数，实例化之后调用一次
function Start()
	local ExchangeFontType = tonumber(G.GameDefine.ExchangeFontType)
	
    img_ImageTitle1.gameObject:SetActive(ExchangeFontType == 0)
    img_ImageTitle2.gameObject:SetActive(ExchangeFontType == 1)
    img_ImageTitle3.gameObject:SetActive(ExchangeFontType == 2)
    img_ImageTitle4.gameObject:SetActive(ExchangeFontType == 3)
	
    img_ImageFAQ1.gameObject:SetActive(ExchangeFontType == 0)
    img_ImageFAQ2.gameObject:SetActive(ExchangeFontType == 1)
    img_ImageFAQ3.gameObject:SetActive(ExchangeFontType == 2)
    img_ImageFAQ4.gameObject:SetActive(ExchangeFontType == 3)
	
    img_FAQ1.gameObject:SetActive(ExchangeFontType == 0)
    img_FAQ2.gameObject:SetActive(ExchangeFontType == 1)
    img_FAQ3.gameObject:SetActive(ExchangeFontType == 2)
    img_FAQ4.gameObject:SetActive(ExchangeFontType == 3)
	
	if ExchangeFontType == 0 then
		changeText = "兑换"
	elseif ExchangeFontType == 1 then
		changeText = "福利"
		Img_Exchange.sprite = Img_ExchangeFont2.sprite
		Img_ExchangeMask.sprite = Img_ExchangeMask2.sprite	
	elseif ExchangeFontType == 2 then
		changeText = "清算"
		Img_Exchange.sprite = Img_ExchangeFont3.sprite
		Img_ExchangeMask.sprite = Img_ExchangeMask3.sprite
	elseif ExchangeFontType == 3 then
		changeText = "提现"
		Img_Exchange.sprite = Img_ExchangeFont4.sprite
		Img_ExchangeMask.sprite = Img_ExchangeMask4.sprite
	end
	
	local TipsText = string.gsub(userdata.ConvertRateTipMsg, "兑换", changeText)
	Txt_Tips.text = TipsText	
	
    --當輸入兌換金額後,手續費變動
    Ipt_ExchangeNumber.onValueChanged:AddListener
    (
    function(r)
        if r==nil or tostring(r)=="" then
            txt_Text.text=0
            return 
        end
       if(tonumber(r)==0) then
       txt_Text.text=0
       elseif (tonumber(r)<100 and tonumber(r)>0)  then
            txt_Text.text=2
       
       else
            txt_Text.text=math.modf(tonumber(r)*0.02)
       end
    end
    )
    
    btn__CloseWin:AddClickListener
    (
    function()
        recTran__ExchangeFAQ.gameObject:SetActive(false);
    end
    )
    btn__FAQ:AddClickListener
    (
    function()
       ToFAQ();
    end
    )
   
    
    Btn_NoBindingUP:AddClickListener
    (
    function()
       BindingUP()
    end
    )
    Btn_NoBindingZF:AddClickListener
    (
    function()
       BindingZF()
    end
    )

    Btn_clear:AddClickListener
    (
    function()
         Ipt_ExchangeNumber.text=0
    end
    )
    Btn_plus50:AddClickListener(function()Plus(50)end)

    Btn_plus500:AddClickListener(function()Plus(500)end)
   
    Btn_plus5000:AddClickListener(function()Plus(5000)end)
   
    Btn_plusmax:AddClickListener
    (
    function()       
          Ipt_ExchangeNumber.text=math.modf((((userdata.CashAmount/100)-10)/1.02)/50)*50    
    end
    )
    btn__CloseBankWin:AddClickListener
    (
    function()       
         recTran_BankSelect.gameObject:SetActive(false) 
    end
    )
    btn__BUGetBank:AddClickListener
    (
    function()       
         recTran_BankSelect.gameObject:SetActive(true) 
    end
    )

    Btn_Exchange:AddClickListener
    (
    function()
       
        local exAmount = G.Tool.UserHelper.GetSaveCoin(G.Tool.StringHelper.IsEmpty(Ipt_ExchangeNumber.text) and 0 or Ipt_ExchangeNumber.text)
		
		-- 检查最小兑换金额
		if exAmount < minexchangevalue then
			CS.UIBubblePrompt.Show(changeText .. "金额最小为" .. minexchangevalue, groupName)
			return
		end
		
        if exAmount <= 0 then
            CS.UIBubblePrompt.Show(changeText .. "金额不能为 0 或负数", groupName)
            return
        end
       
        if exAmount >   userdata.CashAmount then
            CS.UIBubblePrompt.Show(changeText .. "金额不能超过携带金额", groupName)
            return
        end

       


        if G.Tool.StringHelper.IsEmpty(Ipt_ExchangeNumber.text) then
            CS.UIBubblePrompt.Show("请输入" .. changeText .. "数量", groupName)
            return
        end

        if math.fmod( tonumber(Ipt_ExchangeNumber.text), 50 )~=0  then
            CS.UIBubblePrompt.Show("必須是50的倍數才可兌換。", groupName)
            return
        end

        G.Util.coroutine_call
        (
        function()
            
			local msgtext
			if magnification > 0 then
				msgtext = changeText .. Ipt_ExchangeNumber.text .. "金币，系统将收取"..txt_Text.text.."手续费用，确定要" .. changeText .."？"
			else
				msgtext = "确定要" .. changeText .. Ipt_ExchangeNumber.text .. "金币？"			
			end
						
            local boxResult = showMessageBox(CS.MessageBoxStyle.OKCancelClose, "确认", msgtext, "确定", "取消");
            if boxResult == CS.MessageBoxResult.OK then
                local sendMsg = { }

                -- 渠道ID
                sendMsg.Amount = G.Tool.UserHelper.GetSaveCoin(Ipt_ExchangeNumber.text)
                sendMsg.ConvertType=converttype

                -- 兑换
                print("发送兑换请求")
                G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.ConvertMoneyID, sendMsg)
            end
        end
        )()
    end
    )

    Ipt_ExchangeNumber.onEndEdit:AddListener
    (
    function()
        if G.Tool.StringHelper.IsEmpty(Ipt_ExchangeNumber.text) then
            Ipt_ExchangeNumber.text = "0"
        else
            Ipt_ExchangeNumber.text = tostring(G.Tool.Math.ParseInt(Ipt_ExchangeNumber.text))
        end
    end
    )
       Btn_CloseSucessWin:AddClickListener(Close)

       Btn_CloseBingding:AddClickListener(Close)      
     
      
    
--       Btn__CloseWin:AddClickListener(function()
--       recTran__ExchangeLog.gameObject:SetActive(false);
--       end
--       )
      
      
      
        
    -- Add by Bob
    --選擇修改綁定銀聯卡
    Btn_ZFUpdate:AddClickListener
    (
   function()
          BindingZF()
   end
   )
   -- Add by Bob
   --選擇使用支付寶兌換
    Btn_ZF:AddClickListener
    (
    function()      
       
       DisplayZF()
        
    end
    )
    Btn_ZF_1:AddClickListener
    (
    function()      
       
       DisplayZF()
        
    end
    )	
    -- Add by Bob
    --選擇修改綁定銀聯卡
    Btn_UPUpdate:AddClickListener
    (
   function()
          BindingUP()
   end
   )
    -- Add by Bob
    --選擇使用銀聯卡兌換
    Btn_UP:AddClickListener
    (
   function()

        DisplayUP()
    end
    )
    Btn_UP_1:AddClickListener
    (
   function()

        DisplayUP()
    end
    )
	
    Btn_Log:AddClickListener
    (
   function()
       DisplayLog()
    end
    )


     -- Add by Bob
     -- 送出綁定支付寶
    Btn_BZComfirm:AddClickListener
    (
    function()

            -- 检查用户名
            if G.Tool.StringHelper.IsEmpty(Ipt_BZPlayerName.text) then
                CS.UIBubblePrompt.Show("姓名不能为空。", groupName)
                return
            end
           
             -- 检查用户名長度

            -- if G.Tool.StringHelper.length(Ipt_BZPlayerName.text)<2 or G.Tool.StringHelper.length(Ipt_BZPlayerName.text)>6 then

            --     CS.UIBubblePrompt.Show("姓名字数只能能介于2~6之间。", groupName)
            --     return
            -- end

            -- 支付宝账号不能为空
            if G.Tool.StringHelper.IsEmpty(Ipt_BZAccount.text) then
                CS.UIBubblePrompt.Show("账号不能为空。", groupName)
                return
            end 
              -- 驗證碼不能为空
              if G.Tool.StringHelper.IsEmpty(Ipt_BZCode.text) then
                CS.UIBubblePrompt.Show("验证码不能为空。", groupName)
                return
            end      
            
            local sendMsg = { }

            -- 账户
            sendMsg.Zhifubao = Ipt_BZAccount.text

            -- 姓名
            sendMsg.RealName = Ipt_BZPlayerName.text

            -- 手機驗證碼
            sendMsg.VCode = Ipt_BZCode.text

            -- codetype 4:第一次綁定 5:修改
           -- userdata.IsBindZhifubao~=true and userdata.IsBindBankCard ~= true)
             sendMsg.CodeType = userdata.IsBindZhifubao and 5 or 4
            -- 发送登录游戏服务消息
            G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.BindZhifubaoID, sendMsg)        
        
    end
    )
    -- Add by Bob
    Btn_BUComfirm:AddClickListener
    (
    function()
       
            -- 检查用户名
            if G.Tool.StringHelper.IsEmpty(Ipt_BUPlayerName.text) then
                CS.UIBubblePrompt.Show("姓名不能为空。", groupName)
                return
            end
             -- 检查用户名長度

           -- if G.Tool.StringHelper.length(Ipt_BUPlayerName.text)<2 or G.Tool.StringHelper.length(Ipt_BUPlayerName.text)>6 then

           --      CS.UIBubblePrompt.Show("姓名字数只能能介于2~6之间。", groupName)
           --      return
           --  end

            -- 銀聯卡账号不能为空
            if G.Tool.StringHelper.IsEmpty(Ipt_BUAccount.text) then
                CS.UIBubblePrompt.Show("账号不能为空。", groupName)
                return
            end

           
             -- 銀聯卡账号不能为空
            if (lcBankList[Inp_BUBank.text] == nil)    then
                CS.UIBubblePrompt.Show("銀行不能为空。", groupName)
                return
            end   
            
             -- 驗證碼不能为空
              if G.Tool.StringHelper.IsEmpty(Ipt_BUCode.text) then
                CS.UIBubblePrompt.Show("验证码不能为空。", groupName)
                return
            end       

            local sendMsg = { }

            -- 账户
            sendMsg.BankCardNum = Ipt_BUAccount.text

            -- 姓名
            sendMsg.RealName = Ipt_BUPlayerName.text

            --銀行代碼(下拉式選單)
            sendMsg.BankCode =lcBankList[Inp_BUBank.text]

             --手機驗證碼
            sendMsg.VCode = Ipt_BUCode.text

             -- codetype 4:第一次綁定 5:修改
            sendMsg.CodeType = userdata.IsBindBankCard and 5 or 4

            -- 发送登录游戏服务消息
            G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.BindBankCardID, sendMsg)
        
    end
    )

    -- 绑定获取验证码--銀聯卡
    Btn_BUGetCode:AddClickListener
    (
    function()

        -- 检查用户名
            if G.Tool.StringHelper.IsEmpty(Ipt_BUPlayerName.text) then
                CS.UIBubblePrompt.Show("姓名不能为空。", groupName)
                return
            end

            -- 銀聯卡账号不能为空
            if G.Tool.StringHelper.IsEmpty(Ipt_BUAccount.text) then
                CS.UIBubblePrompt.Show("账号不能为空。", groupName)
                return
            end

           
             -- 銀聯卡账号不能为空
            if lcBankList[Inp_BUBank.text] == nil   then
                CS.UIBubblePrompt.Show("銀行不能为空。", groupName)
                return
            end   
            
          
        local sendMsg = { }

        -- 手机号
        sendMsg.Tel = ""
        
        -- 类型
        sendMsg.CodeType = userdata.IsBindBankCard and 5 or 4

        -- 绑定手机验证码
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeID, sendMsg)

    end
    )
      -- 绑定获取验证码--支付寶
    Btn_BZGetCode:AddClickListener
    (
    function()

    -- 检查用户名
            if G.Tool.StringHelper.IsEmpty(Ipt_BZPlayerName.text) then
                CS.UIBubblePrompt.Show("姓名不能为空。", groupName)
                return
            end

            -- 支付宝账号不能为空
            if G.Tool.StringHelper.IsEmpty(Ipt_BZAccount.text) then
                CS.UIBubblePrompt.Show("账号不能为空。", groupName)
                return
            end 
                  
    
        local sendMsg = { }
        -- 手机号
        sendMsg.Tel = ""
        -- 类型
        sendMsg.CodeType = userdata.IsBindZhifubao and 5 or 4
        -- 绑定手机验证码
        G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.SendPhoneVCodeID, sendMsg)

    end
    )
    -- 送出查詢兌換紀錄
    btn__commit:AddClickListener
    (
    function()
     
       
        if exchangeitems == nil or #exchangeitems == 0 then
        --Txt_NoMailSign.gameObject:SetActive(true);
        else
        for i = 1, #exchangeitems, 1 do
           CS.UnityEngine.GameObject.Destroy(exchangeitems[i].gameObject)
        end
        exchangeitems={}
    end
            local typ ={}
            --送出條件 全部與審核不通過皆用99送出查全部
            typ[0] = 99
            typ[1] = 0 
            typ[2] = 1 
            typ[3] = 2 
            typ[4] = 99
            typ[5] = 3
         
            if(dd_type.value==3) then
                IsRejectType=1;
            else 
                IsRejectType=0;
            end
            
         
         G.Service.ExchangeLog.RefreshLogList(typ[dd_type.value],(dd_duration.value+1))


    end)
  
   
end



-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end  

function CreateTimeToString(CreateTime)
	local spdate=''
	for token in string.gmatch(CreateTime, "[^%s]+") do
		-- 避免下一行沒資料卻換行 By Rayon
		if token ~= nil then
			if spdate ~= '' then
				spdate = spdate .."\n"
			end	
			spdate = spdate..token
		end
	end
	
	return spdate
end

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
   
	-- 986 1.兌換倍率為0 2.顯示最小兑换金额提示 3.可以自行输入金额 4.最低兑换金额50 5.关掉金额按钮
	if G.Platform.GetBranch() == 1 then	
		magnification = 0
		minexchangevalue = 50
		RecTran_Text.gameObject:SetActive(false)
		txt__Note.gameObject:SetActive(true)
		Ipt_ExchangeNumber.interactable = true
		Btn_plus50.gameObject:SetActive(false)
		Btn_plus500.gameObject:SetActive(false)
		Btn_plus5000.gameObject:SetActive(false)
		Btn_plusmax.gameObject:SetActive(false)
	else
		magnification = 2
		minexchangevalue = 0
		RecTran_Text.gameObject:SetActive(true)
		txt__Note.gameObject:SetActive(false)
		Ipt_ExchangeNumber.interactable = false
		Btn_plus50.gameObject:SetActive(true)
		Btn_plus500.gameObject:SetActive(true)
		Btn_plus5000.gameObject:SetActive(true)
		Btn_plusmax.gameObject:SetActive(true)
	end   
   
    -- Add by Bob
    G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.GetBankListID, nil)
    -- 兌換紀錄處理
     local logList = G.Service.ExchangeLog.GetLogList();
     DisplayRecord(logList);
    -- 兌換紀錄處理  End

    RecTran_Bingding.gameObject:SetActive(true);

    --btncol[1]=Btn_ZF
    --btncol[2]=Btn_UP
    btncol[3]=Btn_Log

	-- 是否优先显示银联
	if G.Platform.IsChangeExchangeOrder() == false then
		btncol[1]=Btn_ZF
		btncol[2]=Btn_UP
		Btn_ZF_1.gameObject:SetActive(false)
		Btn_UP_1.gameObject:SetActive(false)
		
		if(userdata.ConvertAlipayShow==0 and userdata.ConvertBankShow==1) then
			--  DisplayZF()
			DisplayUP()
		else
			DisplayZF() 
		end
	else
		btncol[1]=Btn_ZF_1
		btncol[2]=Btn_UP_1
		Btn_ZF.gameObject:SetActive(false)
		Btn_UP.gameObject:SetActive(false)

		if(userdata.ConvertBankShow==0 and userdata.ConvertAlipayShow==1) then
			DisplayZF()
		else
			DisplayUP() 
		end
	end
	
	-- 銀聯秒到標籤
	if G.Platform.IsSExchangeMaskUP() == true then
		img_immediately.gameObject:SetActive(true)
		img_immediately_1.gameObject:SetActive(true)
	else
		img_immediately.gameObject:SetActive(false)
		img_immediately_1.gameObject:SetActive(false)	
	end
	
    btncol[1].transform:Find("Disable").gameObject:SetActive(userdata.ConvertAlipayShow==0); 
    btncol[1].transform:GetComponent("Button").interactable=(userdata.ConvertAlipayShow==1);    
    btncol[2].transform:Find("Disable").gameObject:SetActive(userdata.ConvertBankShow==0); 
    btncol[2].transform:GetComponent("Button").interactable=(userdata.ConvertBankShow==1); 

      -- 绑定或修改支付保
    if windowData ~= nil and windowData == 1 then
         if(userdata.ConvertAlipayShow==0 and userdata.ConvertBankShow==1) then            
            if(userdata.IsBindBankCard) then
                DisplayUP() 
            else
                BindingUP()
            end    
        else            
            BindingZF()
        end 
      
    end
    -- 绑定或修改銀聯
    if windowData ~= nil and windowData == 2 then
         if(userdata.ConvertAlipayShow==1 and userdata.ConvertBankShow==0) then            
            if(userdata.IsBindZhifubao) then
                DisplayZF() 
            else
                BindingZF()
            end 
        else          
            BindingUP() 
        end 
      
    end
    -- 绑定或修改支付保
     if windowData ~= nil and windowData == 3 then
        if(userdata.ConvertAlipayShow==0 and userdata.ConvertBankShow==1) then            
            if(userdata.IsBindBankCard) then
                DisplayUP() 
            else
                BindingUP()
            end 
        else          
            DisplayZF()
        end 
    end
    -- 绑定或修改銀連
    if windowData ~= nil and windowData == 4 then
        if(userdata.ConvertAlipayShow==1 and userdata.ConvertBankShow==0) then            
            if(userdata.IsBindZhifubao) then
                DisplayZF() 
            else
                BindingZF()
            end          
        else          
            DisplayUP() 
        end 
    end



end

-- region 兑换回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.ConvertMoneyRetID, function(receiveMsg)
    if receiveMsg.Code == 1 then
        G.Service.ExchangeLog.RefreshLogList()
        Txt_SucessTip.text = receiveMsg.Message;
        RecTran_Exchange.gameObject:SetActive(false);
        RecTran_ExchangeSucess.gameObject:SetActive(true);
    else
        CS.UIBubblePrompt.Show(receiveMsg.Message, groupName);
    end
end , groupName); 

-- 銀行代碼回调
G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.GetBankListRetID, function(receiveMsg)
   
      local BankList =receiveMsg.BankList 
      
        if BankList == nil or #BankList == 0 then
            print("No Bank")
        else
         
            for i = 1, #BankList, 1 do
                --local op1 = CS.UnityEngine.UI.Dropdown.OptionData(BankList[i].BankName);
                --Ipt_BUBank.options:Add(op1);               
                lcBankList[BankList[i].BankName]=BankList[i].BankCode;
                local clone = CS.UnityEngine.GameObject.Instantiate(btn__AgentItem)
	            if(clone == nil)then return nil end
	            clone.transform:SetParent(recTran_Grid.transform)
	            clone.transform.localPosition = CS.UnityEngine.Vector3.zero
	            clone.transform.localScale = CS.UnityEngine.Vector3.one
                clone.gameObject.name=BankList[i].BankCode;
                clone.transform:Find("Name"):GetComponent("Text").text=BankList[i].BankName;
                clone.gameObject:SetActive(true);
                clone:GetComponent("Button"):AddClickListener( function()
						Inp_BUBank.text=BankList[i].BankName;
						recTran_BankSelect.gameObject:SetActive(false);
                        clone.transform:Find("_Select").gameObject:SetActive(true);
                        if(prevbank~=nil)then prevbank:SetActive(false); end
                        prevbank=clone.transform:Find("_Select").gameObject;

                        
					
				end )

                
            end


            -- local TempInt = Ipt_BUBank.value;
--             Ipt_BUBank.value = Ipt_BUBank.value+1;
--             Ipt_BUBank.value = TempInt;
        end
    
end , groupName); 


function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIExchange", false);
end

G.Service.ExchangeLog.RegLogRetCallBack( function()
end
)

function Close()
        CS.WindowManager.Instance:CloseWindow("UIExchange", false);
end

function Plus(Amt)
       if G.Tool.StringHelper.IsEmpty(Ipt_ExchangeNumber.text) then
         Ipt_ExchangeNumber.text=Amt
        else
          Ipt_ExchangeNumber.text=tonumber(Ipt_ExchangeNumber.text)+Amt
        end
end

function BindingUP()
 RecTran_BindingZF.gameObject:SetActive(false);
 RecTran_BindingUP.gameObject:SetActive(true);
 RecTran_NoBindingZF.gameObject:SetActive(false);
 RecTran_NoBindingUP.gameObject:SetActive(false); 
 recTran__ExchangeLog.gameObject:SetActive(false);
  RecTran_Exchange.gameObject:SetActive(false); 

SelectedBtn(2)
end

function BindingZF()
 RecTran_BindingUP.gameObject:SetActive(false);
 RecTran_BindingZF.gameObject:SetActive(true);
 RecTran_NoBindingZF.gameObject:SetActive(false);
 RecTran_NoBindingUP.gameObject:SetActive(false); 
 recTran__ExchangeLog.gameObject:SetActive(false);
 RecTran_Exchange.gameObject:SetActive(false); 
 SelectedBtn(1)
end

function DisplayZF()
 RecTran_BindingZF.gameObject:SetActive(false);
 RecTran_BindingUP.gameObject:SetActive(false);
 RecTran_NoBindingZF.gameObject:SetActive(false);
 RecTran_NoBindingUP.gameObject:SetActive(false);
 recTran__ExchangeLog.gameObject:SetActive(false);
 RecTran_Exchange.gameObject:SetActive(false); 

if userdata.IsBindZhifubao then
        RecTran_Exchange.gameObject:SetActive(true);
        RecTran_NoBindingZF.gameObject:SetActive(false);   
        ToExchangebyZF();
    else       
        BindingZF() 
    end
SelectedBtn(1)
end

function DisplayUP()
RecTran_BindingZF.gameObject:SetActive(false);
 RecTran_BindingUP.gameObject:SetActive(false);
 RecTran_NoBindingZF.gameObject:SetActive(false);
 RecTran_NoBindingUP.gameObject:SetActive(false); 
 recTran__ExchangeLog.gameObject:SetActive(false);
 RecTran_Exchange.gameObject:SetActive(false);

    if userdata.IsBindBankCard then
        RecTran_Exchange.gameObject:SetActive(true);
        RecTran_NoBindingUP.gameObject:SetActive(false);   
        ToExchangebyUP();
    else       
        BindingUP()   
    end
    SelectedBtn(2)
end

function DisplayLog()

 RecTran_BindingUP.gameObject:SetActive(false);
 RecTran_BindingZF.gameObject:SetActive(false);
 RecTran_NoBindingZF.gameObject:SetActive(false);
 RecTran_NoBindingUP.gameObject:SetActive(false); 
 recTran__ExchangeLog.gameObject:SetActive(true);
 RecTran_Exchange.gameObject:SetActive(false); 
 SelectedBtn(3)
end

function SelectedBtn(index)
btncol[1].transform:Find("selected").gameObject:SetActive(index==1);
btncol[2].transform:Find("selected").gameObject:SetActive(index==2);
btncol[3].transform:Find("selected").gameObject:SetActive(index==3);

end

function ToExchangebyZF()        
                converttype=1
                img_Payment_ZF.gameObject:SetActive(true);
                img_Payment_UP.gameObject:SetActive(false);
                Btn_ZFUpdate.gameObject:SetActive(true);
                Btn_UPUpdate.gameObject:SetActive(false);
                Txt_Money.text = G.Tool.UserHelper.ToCashDetail(userdata.CashAmount);
                Txt_SafeMoney.text=G.Tool.UserHelper.ToCashDetail(userdata.BankAmount);
                Txt_Account.text = userdata.Zhifubao;
                Ipt_ExchangeNumber.text = "0";
                
end

function ToExchangebyUP()     
                converttype=2
                
                img_Payment_ZF.gameObject:SetActive(false);
                img_Payment_UP.gameObject:SetActive(true);
                Btn_ZFUpdate.gameObject:SetActive(false);
                Btn_UPUpdate.gameObject:SetActive(true);
                Txt_Money.text = G.Tool.UserHelper.ToCashDetail(userdata.CashAmount);
                Txt_SafeMoney.text=G.Tool.UserHelper.ToCashDetail(userdata.BankAmount);
                Txt_Account.text = userdata.BankName .."   ".. userdata.BankCardNum;
                Ipt_ExchangeNumber.text = "0";
   
	-- 986银行卡不显示修改按钮
	if G.Platform.GetBranch() == 1 then
		Btn_UPUpdate.gameObject:SetActive(false)
	end	   
end

function ToFAQ()        
                --Txt_Tips.text = userdata.ConvertRateTipMsg
				
                recTran__ExchangeFAQ.gameObject:SetActive(true);
end

function DisplayRecord(logList)
        --不屬於審核不通過之清單
          sta[0] = "提出" .. changeText
          sta[1] = "审核中"
          sta[2] = "审核通过"
          sta[3] = "已打款"
          



  if logList == nil or #logList == 0 then
        --Txt_NoMailSign.gameObject:SetActive(true);
  else
        
       local indexi=1
            for i = 1, #logList, 1 do

             if  sta[tonumber(logList[i].Status)]~=nil and IsRejectType==1 then  
        
             else

                local logItem 
                if math.fmod( i, 2 )~=0 then
                    logItem=CS.UnityEngine.GameObject.Instantiate(RecTran_Item.gameObject);
                else
                    logItem=CS.UnityEngine.GameObject.Instantiate(RecTran_Item1.gameObject);
                end
             
                exchangeitems[indexi]=logItem;
                indexi=indexi+1;
                --logItem.transform.parent = RecTran_Item.gameObject.transform.parent;
                logItem.transform:SetParent(RecTran_Item.gameObject.transform.parent.transform)
                logItem.transform.localScale = CS.UnityEngine.Vector3(1, 1, 1)
                logItem.transform.localPosition = CS.UnityEngine.Vector3.zero
                logItem:SetActive(true)
           
                logItem.name = tostring(i)

			    -- 避免下一行沒資料卻換行 By Rayon
			    local spdate = CreateTimeToString(logList[i].CreateTime)			
			    --[[
                local spdate=''
                for token in string.gmatch(logList[i].CreateTime, "[^%s]+") do
				    spdate=spdate..token.."\n"
                end
			    ]]--
            
                logItem.transform:Find("Amt"):GetComponent("Text").text = logList[i].Amount/100
                logItem.transform:Find("fee"):GetComponent("Text").text = logList[i].TaxAmount/100
                logItem.transform:Find("type"):GetComponent("Text").text =(sta[tonumber(logList[i].Status)]~=nil) and (sta[tonumber(logList[i].Status)]) or "<color=red>审核不通过，请洽客服人员</color>"
                logItem.transform:Find("payment"):GetComponent("Text").text = logList[i].Payment
                logItem.transform:Find("SendTime"):GetComponent("Text").text =spdate
             
              end
            end
        
    end

end

  

