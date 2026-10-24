require("Module")
mmodule(...)

local G = require "Global"

function Init(this)
    local btn_shade = this.transform:Find("shade"):GetComponent("Button")
    local btn_close = this.transform:Find("close"):GetComponent("Button")
    local btn_music = this.transform:Find("sound/music"):GetComponent("Button")
    local btn_sound = this.transform:Find("sound/musicef"):GetComponent("Button")

    local sound_state = G.Tool.Sound.SoundStatus()
    local music_state = G.Tool.Sound.MusicStatus()
    local music_on = btn_music.transform:Find("on").gameObject
    local music_off = btn_music.transform:Find("off").gameObject

    local sound_on = btn_sound.transform:Find("on").gameObject
    local sound_off = btn_sound.transform:Find("off").gameObject

	-- 顯示版號
	versionData = this.transform:Find("version");
	if versionData ~= nil then
		local txt_version = versionData:GetComponent("Text")
		if txt_version ~= nul then
			local localVersion = CS.HotUpdate.HotUpdateManager.Instance:GetProjectVersion(G.GetGamePName());
			txt_version.text = "V." .. tostring(localVersion)
		end
	end
	
    local function event_music()
        if (music_state == true) then
            music_on:SetActive(true)
            music_off:SetActive(false)
        else
            music_on:SetActive(false)
            music_off:SetActive(true)
        end
    end
    btn_music:AddClickListener( function()
        music_state = not music_state
        G.Tool.Sound.MusicStatus(music_state)
        G.GameConfig.MusicStatus=music_state;
        event_music()
    end )


    local function event_sound()
        if (sound_state == true) then
            sound_on:SetActive(true)
            sound_off:SetActive(false)
        else
            sound_on:SetActive(false)
            sound_off:SetActive(true)
        end
    end

    btn_sound:AddClickListener( function()
        sound_state = not sound_state
        G.Tool.Sound.SoundStatus(sound_state)
        G.GameConfig.SoundStatus=sound_state;
        event_sound()
    end )

    event_music()
    event_sound()

    local closeFun = function()
        CS.WindowManager.Instance:CloseWindow("UISetGame", false)
    end
    if (btn_shade ~= nil) then
        btn_shade:AddClickListener(closeFun)
    end

    btn_close:AddClickListener(closeFun)
end

