require("Module")
mmodule(...)

local G = require "Global"
local hallPName = CS.AppDefine.HallNumber

-- 音效开关状态 获取/设置 传值设置，不传获取
function SoundStatus(val)
    if (val == nil) then
        return CS.SoundManager.Instance.IsSoundOn
    end

    if (val == true) then
        CS.SoundManager.Instance.IsSoundOn = true
    else
        CS.SoundManager.Instance.IsSoundOn = false
		CS.SoundManager.Instance:CloseSound()
    end
end

-- 音乐开状状态 获取/设置 传值设置，不传获取
function MusicStatus(val)
    if (val == nil) then
        return CS.SoundManager.Instance.IsMusicOn
    end

    if (val == true) then
        CS.SoundManager.Instance.IsMusicOn = true
    else
        CS.SoundManager.Instance.IsMusicOn = false
    end
end

-- 播放游戏音乐
function PlayGameBGM(audioClipName)
    CS.SoundManager.Instance:PlayBGM(G.GetGamePName(), audioClipName)
end
-- 播放游戏音效
function PlayGameSound(audioClipName, isLoop)
    isLoop = G.CaseWhen(isLoop == nil, false, isLoop)
    CS.SoundManager.Instance:PlaySound(G.GetGamePName(), audioClipName, isLoop)
    --print("播放游戏音效", audioClipName)
end
-- 播放游戏背景音乐
function PlayHallBGM(audioClipName)
    CS.SoundManager.Instance:PlayBGM(hallPName, audioClipName)
end

-- 播放大厅音效
function PlayHallSound(audioClipName, isLoop)
    isLoop = G.CaseWhen(isLoop == nil, false, isLoop)
    CS.SoundManager.Instance:PlaySound(hallPName, audioClipName, isLoop)
    --print("播放大厅音效", audioClipName)
end

-- 关闭循环音效
function CloseLoopSound(audioClipName)
    CS.SoundManager.Instance:CloseLoopSound(audioClipName)
end

-- 关闭音效
function CloseSound()
    CS.SoundManager.Instance:CloseSound()
end
-- 关闭音乐
function CloseBGMSound()
    CS.SoundManager.Instance:CloseBGMSound()
end