local Btn_close = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local ipt__PromotionCode = this.transform:Find("UIBody/_AgentInfo/_PromotionCode"):GetComponent("InputField")
local btn__Binding = this.transform:Find("UIBody/_AgentInfo/_PromotionCode/_Binding"):GetComponent("Button")
local recTran__AgentInfo = this.transform:Find("UIBody/_AgentInfo"):GetComponent("RectTransform")
local recTran__ApplyAgent = this.transform:Find("UIBody/_ApplyAgent"):GetComponent("RectTransform")

local txt__AgentID = this.transform:Find("UIBody/_AgentInfo/Normal/_MyAgent/_AgentID"):GetComponent("Text")
local btn__ApplyAgent = this.transform:Find("UIBody/_ApplyAgent/_BtnApplyAgent"):GetComponent("Button")
local img_AgentCode = this.transform:Find("UIBody/_AgentInfo/Agent/AgentCode"):GetComponent("Image")
local txt__AgentCode = this.transform:Find("UIBody/_AgentInfo/Agent/AgentCode/_AgentCode"):GetComponent("Text")
local txt__LV = this.transform:Find("UIBody/_AgentInfo/Agent/_LV/_Value"):GetComponent("Text")
local txt__CurrLv = this.transform:Find("UIBody/_AgentInfo/Agent/_CurrLv/_Value"):GetComponent("Text")
local txt__NextLv = this.transform:Find("UIBody/_AgentInfo/Agent/_NextLv/_Value"):GetComponent("Text")
local img__CurrLv = this.transform:Find("UIBody/_AgentInfo/Agent/_CurrLv/icon"):GetComponent("Image")
local img__NextLv = this.transform:Find("UIBody/_AgentInfo/Agent/_NextLv/icon"):GetComponent("Image")
 
local txt__ApplyStatus = this.transform:Find("UIBody/_ApplyAgent/_BtnApplyAgent/_ApplyStatus"):GetComponent("Text")

local btn__TabInfo = this.transform:Find("UIBody/_AgentTab/_TabInfo"):GetComponent("Button")
local btn__TabAgentApply = this.transform:Find("UIBody/_AgentTab/_TabAgentApply"):GetComponent("Button")

local recTran_Normal = this.transform:Find("UIBody/_AgentInfo/Normal"):GetComponent("RectTransform")
local recTran_Agent = this.transform:Find("UIBody/_AgentInfo/Agent"):GetComponent("RectTransform")
local sld_Slider = this.transform:Find("UIBody/_AgentInfo/Agent/Slider"):GetComponent("Slider")
local btn_BtnCopyCode = this.transform:Find("UIBody/_AgentInfo/Agent/BtnCopyCode"):GetComponent("Button")


local img_qrcode = this.transform:Find("UIBody/_AgentInfo/Agent/qrcode"):GetComponent("Image")














--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local groupName = "UIAgent";
local user = G.User.GetCurrentUser();
local btncol={}
local AgentContent=nil;

-- 初始化函数，实例化之后调用一次
function Start()

btn_BtnCopyCode:AddClickListener
    (
    function()
        CS.Utility.CopyText(AgentContent.AgentDetailData.InviteCode)
        CS.UIBubblePrompt.Show("复制成功", groupName)

    end)
 btn__Binding:AddClickListener
    (
    function()
        if G.Tool.StringHelper.IsEmpty(ipt__PromotionCode.text) then
            CS.UIBubblePrompt.Show("请输入邀请码", groupName)
            return
        end
         if string.len(ipt__PromotionCode.text)<5 then
            CS.UIBubblePrompt.Show("邀请码长度不得低于5码", groupName)
            return
        end
        local JoinAgentList  = { }
    
        JoinAgentList.UserID = user.UserID
        JoinAgentList.InviteCode = ipt__PromotionCode.text
        JoinAgentList.RetCode = 0
        JoinAgentList.RetMsg="000"

        G.Service.Agent.ApplyForDownline( (function(data)
                 if(data.RetCode==1) then
                     RequestData();
                  else
                     CS.UIBubblePrompt.Show(data.RetMsg, groupName)
                  end
        end),JoinAgentList )
       
    end
    )
     btn__ApplyAgent:AddClickListener
    (
    function()

    ApplyForAgent();

    end)
    Btn_close:AddClickListener( function()
        EscapeWindow()
    end )

     btn__TabInfo:AddClickListener( function()
       ShowAgentInfo()
    end )

    btn__TabAgentApply:AddClickListener( function()
       ShowApplyAgent()
    end )
    
   
end

function RefreshWindowData(windowData)

  btncol[1]=btn__TabInfo
  btncol[2]=btn__TabAgentApply
  ShowAgentInfo();
   
end


-- 在Window下按下Esc或在Android 下按下Back
function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIAgent", false)
end

function RequestData()
     G.Service.Agent.GetAgentData( function(content)
      AgentContent=content;
      DisplayAgent();
   end )

end

function ApplyForAgent()
 G.Service.Agent.ApplyForAgent( function(content)
    if(content.RetCode==-1) then
        CS.UIBubblePrompt.Show(content.RetMsg, groupName)
    else
        CS.UIBubblePrompt.Show("您的申请已成功送出", groupName) 
        btn__ApplyAgent.interactable = false;
   end
   end )
end

function CheckForAgent()
 G.Service.Agent.CheckForAgent( function(content)
 
   if(content.RetCode<=-1) then
       txt__ApplyStatus.text=content.RetMsg;
     
     --   CS.UIBubblePrompt.Show(content.RetMsg, groupName)
    if(content.RetCode==-3) then
        btn__ApplyAgent.interactable = false;
        btn__ApplyAgent.transform:Find("Disabled").gameObject:SetActive(false);
    end
   end

    
   end )
end


function DisplayAgent()

    if(AgentContent.RetCode==1) then
         ipt__PromotionCode.gameObject:SetActive(false);
        
       
       if(AgentContent.AgentDetailData.ParentUserID==0) then
            txt__AgentID.text="总代理";
       else
            txt__AgentID.text=AgentContent.AgentDetailData.ParentUserID;
       end
        --判斷是否可以申請代理
       if(AgentContent.AgentDetailData.CanBeAgent==1 and AgentContent.AgentDetailData.Class==99) then
        btn__TabAgentApply.gameObject:SetActive(true);
        CheckForAgent()
       end

       if(AgentContent.AgentDetailData.InviteCode ~='') then
          recTran_Normal.gameObject:SetActive(false);
          recTran_Agent.gameObject:SetActive(true);  
          
          txt__AgentCode.text=AgentContent.AgentDetailData.InviteCode
          img_qrcode.sprite=CS.Utility.GenQRCode("http://0.0.0.0:130/agentclub/admin/RedirectDownload?userID="..user.UserID);
          txt__LV.text=tostring(AgentContent.AgentDetailData.TotalFlows/100).."/"..tostring(AgentContent.AgentDetailData.UpgradeLimit/100)
          local t1 =( tonumber(AgentContent.AgentDetailData.TotalFlows)/AgentContent.AgentDetailData.UpgradeLimit);
        
         
          sld_Slider.value=t1
          txt__CurrLv.text=AgentContent.AgentDetailData.GradeName
          img__CurrLv:ResetSpriteByName(CS.AppDefine.HallNumber, "agentlv"..AgentContent.AgentDetailData.Level, true)
			
          txt__NextLv.text=AgentContent.AgentDetailData.NextGradeName
          img__NextLv:ResetSpriteByName(CS.AppDefine.HallNumber, "agentlv"..((tonumber(AgentContent.AgentDetailData.Level)+1)>3 and "0" or tostring(tonumber(AgentContent.AgentDetailData.Level)+1)) , true)
		
       else
          recTran_Normal.gameObject:SetActive(true);
          recTran_Agent.gameObject:SetActive(false);  
       end
        
      else
        
        recTran_Normal.gameObject:SetActive(false);
        recTran_Agent.gameObject:SetActive(false);       
        ipt__PromotionCode.gameObject:SetActive(true);
      end
end

function ShowAgentInfo()
  RequestData();
  recTran__AgentInfo.gameObject:SetActive(true); 
  recTran__ApplyAgent.gameObject:SetActive(false);
  
  SelectedBtn(1)

end
function ShowApplyAgent()
 recTran__AgentInfo.gameObject:SetActive(false); 
 recTran__ApplyAgent.gameObject:SetActive(true);

 SelectedBtn(2)
end

function SelectedBtn(index)
btncol[1].transform:Find("selected").gameObject:SetActive(index==1);
btncol[2].transform:Find("selected").gameObject:SetActive(index==2);

end

