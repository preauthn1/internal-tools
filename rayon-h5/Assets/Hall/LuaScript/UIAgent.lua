local Btn_close = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local recTran__AgentInfo = this.transform:Find("UIBody/_AgentInfo"):GetComponent("RectTransform")
local btn__TabInfo = this.transform:Find("UIBody/_AgentTab/_TabInfo"):GetComponent("Button")
local btn__TabAgentConsult = this.transform:Find("UIBody/_AgentTab/_TabAgentConsult"):GetComponent("Button")
local recTran_Agent = this.transform:Find("UIBody/_AgentInfo/Agent"):GetComponent("RectTransform")
local btn_BtnCopyCode = this.transform:Find("UIBody/_AgentInfo/Agent/BtnCopyCode"):GetComponent("Button")
local text_NoDataTips = this.transform:Find("UIBody/_AgentConsult/_NoDataTips"):GetComponent("Text")
local recTran__AgentConsult = this.transform:Find("UIBody/_AgentConsult"):GetComponent("RectTransform")
local img_qrcode = this.transform:Find("UIBody/_AgentInfo/Agent/qrcode"):GetComponent("Image")

-- local recTran_GA = this.transform:Find("UIBody/_AgentInfo/Agent/General-Agency"):GetComponent("RectTransform")
-- --今日收益
-- local txt__TodayIncome = this.transform:Find("UIBody/_AgentInfo/Agent/TodayIncome/_TodayIncome"):GetComponent("Text")


--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---

local G = require "Global";
local groupName = "UIAgent";
local user = G.User.GetCurrentUser();
local btncol={}
local AgentContent=nil;

-- 初始化函数，实例化之后调用一次
function Start()

	btn_BtnCopyCode:AddClickListener(function()
		local url = AgentContent.AgentDetailData.QRUrl
        local genqrcode = CS.Utility.GenQRCode        
          
        CS.Utility.CopyText(url)
       
        CS.UIBubblePrompt.Show("复制成功", groupName)
    end)
 
    
    Btn_close:AddClickListener(function()
        EscapeWindow()
    end)

    btn__TabInfo:AddClickListener(function()
       ShowAgentInfo()
    end)
   
    btn__TabAgentConsult:AddClickListener(function()
       ShowAgentConsult()
    end)
end

function RefreshWindowData(windowData)

	Consult = {}
	ConsultIndex = 1
	Consult[ConsultIndex] = {}
	Consult[ConsultIndex].btn = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex):GetComponent("Button")
	Consult[ConsultIndex].content = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content"):GetComponent("Text")
	Consult[ConsultIndex].icon = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content/ImageTitle"):GetComponent("Image")
	Consult[ConsultIndex].btn.gameObject:SetActive(false)
	
	ConsultIndex = 2
	Consult[ConsultIndex] = {}
	Consult[ConsultIndex].btn = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex):GetComponent("Button")
	Consult[ConsultIndex].content = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content"):GetComponent("Text")
	Consult[ConsultIndex].icon = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content/ImageTitle"):GetComponent("Image")
	Consult[ConsultIndex].btn.gameObject:SetActive(false)
	
	ConsultIndex = 3
	Consult[ConsultIndex] = {}
	Consult[ConsultIndex].btn = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex):GetComponent("Button")
	Consult[ConsultIndex].content = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content"):GetComponent("Text")
	Consult[ConsultIndex].icon = this.transform:Find("UIBody/_AgentConsult/_BtnOpen" .. ConsultIndex .. "/_Content/ImageTitle"):GetComponent("Image")
	Consult[ConsultIndex].btn.gameObject:SetActive(false)
	
	btncol[1]=btn__TabInfo
	btncol[2]=btn__TabAgentConsult
	ShowAgentConsult();
   
end

function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(groupName)

    -- 清空当前界面的浮动提示消息
    CS.UIBubblePrompt.ClearPrompt(groupName)
end

-- 在Window下按下Esc或在Android 下按下Back
function EscapeWindow()
    CS.WindowManager.Instance:CloseWindow("UIAgent", false)
end

function RequestData(callback)
	G.Service.Agent.GetAgentData( function(content)
		AgentContent=content;
		--btn__TabInfo.gameObject:SetActive(AgentContent.AgentDetailData.InviteCode ~='')
		CS.UILoadingData.Hide()
		callback()
   end )
end




function DisplayAgent()

    if(AgentContent.RetCode==1) then
        --ipt__PromotionCode.gameObject:SetActive(false);
		
		-- 当已成为代理时
		if(AgentContent.AgentDetailData.InviteCode ~='') then
           
            recTran_Agent.gameObject:SetActive(true);  
            
           -- recTran_GA.gameObject:SetActive(AgentContent.AgentDetailData.Class == 0 );

            local url = AgentContent.AgentDetailData.QRUrl
            local genqrcode = CS.Utility.GenQRCode
          
            if(genqrcode~=nil) then
				img_qrcode.sprite=CS.Utility.GenQRCode(url);
            end
          
		else           
            recTran_Agent.gameObject:SetActive(false); 
            recTran__AgentInfo.gameObject:SetActive(false);  
		end
        
    else             
        recTran_Agent.gameObject:SetActive(false);  
        recTran__AgentInfo.gameObject:SetActive(false);       
    end

      
end

function DisplayAgentConsult()
	recTran__AgentConsult.gameObject:SetActive(true);   
	ConsultIndex = 1
	Consult[ConsultIndex].btn.gameObject:SetActive(not G.Tool.StringHelper.IsEmpty(AgentContent.AgentDetailData.OfficeUser1));
	Consult[ConsultIndex].content.text = AgentContent.AgentDetailData.OfficeUser1;      
	local index = AgentContent.AgentDetailData.User1Type	
	Consult[ConsultIndex].btn.image:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentCopy" .. index)
	Consult[ConsultIndex].icon:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentIcon" .. index)
	Consult[ConsultIndex].btn:AddClickListener (function()
		CS.Utility.CopyText(Consult[1].content.text)
		CS.UIBubblePrompt.Show("复制成功", groupName)
		if(tostring(AgentContent.AgentDetailData.User1Type) == "1")then
			CS.UnityEngine.Application.OpenURL("weixin://")
		elseif(tostring(AgentContent.AgentDetailData.User1Type) == "2")then
			CS.UnityEngine.Application.OpenURL("mqqwpa://im/chat?chat_type=wpa&uin=" .. Consult[1].content.text)
		end
	end)
	
	ConsultIndex = 2
	Consult[ConsultIndex].btn.gameObject:SetActive(not G.Tool.StringHelper.IsEmpty(AgentContent.AgentDetailData.OfficeUser2));
	Consult[ConsultIndex].content.text = AgentContent.AgentDetailData.OfficeUser2;
	local index = AgentContent.AgentDetailData.User2Type
	Consult[ConsultIndex].btn.image:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentCopy" .. index)
	Consult[ConsultIndex].icon:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentIcon" .. index)
	Consult[ConsultIndex].btn:AddClickListener (function()
		CS.Utility.CopyText(Consult[2].content.text)
		CS.UIBubblePrompt.Show("复制成功", groupName)
		if(tostring(AgentContent.AgentDetailData.User2Type) == "1")then
			CS.UnityEngine.Application.OpenURL("weixin://")
		elseif(tostring(AgentContent.AgentDetailData.User2Type) == "2")then
			CS.UnityEngine.Application.OpenURL("mqqwpa://im/chat?chat_type=wpa&uin=" .. Consult[2].content.text)
		end
	end)
	
	ConsultIndex = 3
	Consult[ConsultIndex].btn.gameObject:SetActive(not G.Tool.StringHelper.IsEmpty(AgentContent.AgentDetailData.OfficeUser3));
	Consult[ConsultIndex].content.text = AgentContent.AgentDetailData.OfficeUser3;
	local index = AgentContent.AgentDetailData.User3Type
	Consult[ConsultIndex].btn.image:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentCopy" .. index)
	Consult[ConsultIndex].icon:ResetSpriteByName(CS.AppDefine.HallNumber, "AgentIcon" .. index)
	Consult[ConsultIndex].btn:AddClickListener (function()
		CS.Utility.CopyText(Consult[3].content.text)
		CS.UIBubblePrompt.Show("复制成功", groupName)
		if(tostring(AgentContent.AgentDetailData.User3Type) == "1")then
			CS.UnityEngine.Application.OpenURL("weixin://")
		elseif(tostring(AgentContent.AgentDetailData.User3Type) == "2")then
			CS.UnityEngine.Application.OpenURL("mqqwpa://im/chat?chat_type=wpa&uin=" .. Consult[3].content.text)
		end
	end)
	
	if (Consult[1].btn.gameObject.activeSelf == false and
		Consult[2].btn.gameObject.activeSelf == false and
		Consult[3].btn.gameObject.activeSelf == false) then
		text_NoDataTips.gameObject:SetActive(true)
	else
		text_NoDataTips.gameObject:SetActive(false)
	end
end

--展示代理资讯
function ShowAgentInfo()
	CS.UILoadingData.Show("资料读取中......");
	RequestData(function()

		DisplayAgent();

	end);
	recTran__AgentInfo.gameObject:SetActive(true); 
	recTran__AgentConsult.gameObject:SetActive(false);
  
	SelectedBtn(2)

end


function ShowAgentConsult()
	CS.UILoadingData.Show("资料读取中......");
	RequestData(function()

		DisplayAgentConsult();
		if (AgentContent.AgentDetailData.InviteCode ~='') then
			SelectedBtn(1);
		end

	end);
	recTran__AgentInfo.gameObject:SetActive(false); 
	recTran__AgentConsult.gameObject:SetActive(true);
end

function SelectedBtn(index)
btncol[1].gameObject:SetActive(index==1);
btncol[2].gameObject:SetActive(index==2);

end




