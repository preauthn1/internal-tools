local G = require "Global"

-- 接收消息table
local windowName = "UIActivity"



local yield_return = G.Util.async_to_sync(function(to_to_yield, cb)
    this:YieldAndCallback(to_to_yield, cb)
end)

local togPrefab = this.transform:Find("UIBody/Toggle/Root/ToggleGroup/Activity").gameObject
local buttonCloss = this.transform:Find("UIBody/ButtonClose"):GetComponent("Button")
local txtLoading = this.transform:Find("UIBody/Loading"):GetComponent("Text")
local imgDetail = this.transform:Find("UIBody/Detail"):GetComponent("Image")

local imgDetailTmp = {}
local selectSeqID = nil;

local isOpen = false;

function Init(this)
end

function Start()
	buttonCloss:AddClickListener(function ()
		CS.WindowManager.Instance:CloseWindow(windowName,false)
	end)
	togPrefab:SetActive(false)
	-- region 
	G.NetQ.RegMsgHook(G.Proto.BClass.Hall, G.Proto.Hall.CampaignGetListRetID,  ReceiveCampaignGetListRetID, windowName);
	SendCampaignGetListID()
end


-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)
	selectSeqID = tonumber(windowData)
end
	
local CompareTicket = function (a,b)
  return a.SeqID<b.SeqID
 end

function SendCampaignGetListID()
	local param = {}
	param.UserID = G.User.GetCurrentUserID()
	G.NetQ.SendMsg(G.Proto.BClass.Hall, G.Proto.Hall.CampaignGetListID, param)
end
	
function ReceiveCampaignGetListRetID(receiveMsg)
	
	local activity = receiveMsg.CampaignData or {}
	
	if #activity == 0 then 
		txtLoading.text = "目前无活动"
	else
		for i=1,#activity do 
			table.sort(activity,CompareTicket)
			
			

			local togClone = CS.UnityEngine.Object.Instantiate(togPrefab)
			togClone.transform:SetParent(togPrefab.transform.parent)
			togClone.transform.localPosition = CS.UnityEngine.Vector3.zero
			togClone.transform.localScale = CS.UnityEngine.Vector3.one
			togClone.gameObject.name = "Activity"..i
			togClone:SetActive(true)
			local tog = togClone:GetComponent("Toggle")
			local txtName = togClone.transform:Find("Name"):GetComponent("Text")
			local outlineName = txtName.gameObject:GetComponent("Outline")
			
			txtName.text = activity[i].Title
			tog.onValueChanged:AddListener(
				function(isOn)
					if isOn == true then 
						txtName.color = CS.UnityEngine.Color(255/255,243/255,202/255,255/255)
						outlineName.effectColor = CS.UnityEngine.Color(190/255,68/255,29/255,205/255)
						if imgDetailTmp[i] ~= nil then
							imgDetail.sprite = imgDetailTmp[i]
							imgDetail.color = CS.UnityEngine.Color.white
							
						else
							SetDetailImage(activity[i].Picture,i) 
						end
					else
						txtName.color = CS.UnityEngine.Color(165/255,91/255,41/255,255/255)
						outlineName.effectColor = CS.UnityEngine.Color.clear
					end
				end
			)
			if isOpen == false then
			        --沒選擇就顯示第一個活動或顯示選擇活動
					if (selectSeqID == nil and i == 1) or activity[i].SeqID == selectSeqID then
					tog.isOn = true
					isOpen = true
				end
			end
		end
		if isOpen == false then
			this.transform:Find("UIBody/Toggle/Root/ToggleGroup/Activity1").gameObject:GetComponent("Toggle").isOn = true
		end	
	end
end


--下載圖片並顯示於Detail
function SetDetailImage(url,saveIndex)
	G.Util.coroutine_call(function()
		imgDetail.color = CS.UnityEngine.Color.clear
		txtLoading.text = "读取中.."
        local www = CS.UnityEngine.WWW(url)
        yield_return(www);
        if www.error ~= nil then
			txtLoading.text = ""
            print(www.error)
        else
            if www.isDone and www.texture ~= nil and www.texture.width ~= 8 then
                imgDetail.sprite = CS.UnityEngine.Sprite.Create(www.texture, CS.UnityEngine.Rect(0, 0, www.texture.width, www.texture.height), CS.UnityEngine.Vector2(0, 0));
				imgDetail.color = CS.UnityEngine.Color.white
				imgDetailTmp[saveIndex] = imgDetail.sprite --存入到暫存
				txtLoading.text = ""
            end
        end
    end)()
end

function OnDestroy()
    -- 回收当前页面注册的消息
    G.NetQ.RemoveMsgHookByGroup(windowName)
end


	
	
-- local receiveMsg = {}
	-- receiveMsg.CampaignData = {}
	-- receiveMsg.CampaignData[1] = {}
	-- receiveMsg.CampaignData[1].ID = 123
	-- receiveMsg.CampaignData[1].Picture = "https://example.invalid/transparent250_/ha-ha-simpsons-png.png"
	-- receiveMsg.CampaignData[1].SeqID = 5
	-- receiveMsg.CampaignData[1].Title = "A"
	-- receiveMsg.CampaignData[1].CampaignTag = 0
	-- receiveMsg.CampaignData[1].CampaignBehavior = 0
	-- receiveMsg.CampaignData[1].Link = "https://example.invalid/wiki/MVC"
	-- receiveMsg.CampaignData[1].StartTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")}))
	-- receiveMsg.CampaignData[1].EndTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")+1}))
	
	-- receiveMsg.CampaignData[2] = {}
	-- receiveMsg.CampaignData[2].ID = 456
	-- receiveMsg.CampaignData[2].Picture = "https://example.invalid/transparent250_/png-imagenes.png"
	-- receiveMsg.CampaignData[2].SeqID = 4
	-- receiveMsg.CampaignData[2].Title = "B"
	-- receiveMsg.CampaignData[2].CampaignTag = 0
	-- receiveMsg.CampaignData[2].CampaignBehavior = 0
	-- receiveMsg.CampaignData[2].Link = "https://example.invalid/wiki/MVC"
	-- receiveMsg.CampaignData[2].StartTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")}))
	-- receiveMsg.CampaignData[2].EndTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")+1}))
	
	-- receiveMsg.CampaignData[3] = {}
	-- receiveMsg.CampaignData[3].ID = 789
	-- receiveMsg.CampaignData[3].Picture = "https://example.invalid/news2/2018/64/62/q5qYoaGVk6eYqA.jpg"
	-- receiveMsg.CampaignData[3].SeqID = 3
	-- receiveMsg.CampaignData[3].Title = "C"
	-- receiveMsg.CampaignData[3].CampaignTag = 0
	-- receiveMsg.CampaignData[3].CampaignBehavior = 0
	-- receiveMsg.CampaignData[3].Link = "https://example.invalid/wiki/MVC"
	-- receiveMsg.CampaignData[3].StartTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")}))
	-- receiveMsg.CampaignData[3].EndTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")+1}))
	
	-- receiveMsg.CampaignData[4] = {}
	-- receiveMsg.CampaignData[4].ID = "1"
	-- receiveMsg.CampaignData[4].Picture = "https://example.invalid/FGrHQeN.jpg"
	-- receiveMsg.CampaignData[4].SeqID = 2
	-- receiveMsg.CampaignData[4].Title = "D"
	-- receiveMsg.CampaignData[4].CampaignTag = 0
	-- receiveMsg.CampaignData[4].CampaignBehavior = 0
	-- receiveMsg.CampaignData[4].Link = "https://example.invalid/wiki/MVC"
	-- receiveMsg.CampaignData[4].StartTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")}))
	-- receiveMsg.CampaignData[4].EndTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")+1}))
	
	-- receiveMsg.CampaignData[5] = {}
	-- receiveMsg.CampaignData[5].ID = 5265456
	-- receiveMsg.CampaignData[5].Picture = "https://example.invalid/images?q=tbn:ANd9GcS_H6KtZpi5F9NqkDuOwOYb3LZI86CTIv-YQTGT8Wn6HVq6l8nzXA"
	-- receiveMsg.CampaignData[5].SeqID = 1
	-- receiveMsg.CampaignData[5].Title = "E"
	-- receiveMsg.CampaignData[5].CampaignTag = 0
	-- receiveMsg.CampaignData[5].CampaignBehavior = 0
	-- receiveMsg.CampaignData[5].Link = "https://example.invalid/wiki/MVC"
	-- receiveMsg.CampaignData[5].StartTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")}))
	-- receiveMsg.CampaignData[5].EndTime = os.date("%Y-%m-%d",os.time({year=os.date("%Y"),month=os.date("%m"),day=os.date("%d")+1}))
	
	-- ReceiveCampaignGetListRetID(receiveMsg)
	
	
	
	
	
	
	