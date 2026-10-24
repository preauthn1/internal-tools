require("Module")
mmodule(...)
local G = require "Global"
-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成

function Init(this)
    local btn_close = this.transform:Find("shade"):GetComponent("Button")
    local btn_ok = this.transform:Find("btn_ok"):GetComponent("Button")
    local btn_cancel = this.transform:Find("btn_cancel"):GetComponent("Button")

    btn_ok:AddClickListener( function()
        CS.WindowManager.Instance:CloseWindow("UIExitGame", false)
        G.Service.Game.ExitGame()
    end )

    local closeWin = function()
        CS.WindowManager.Instance:CloseWindow("UIExitGame", false)
    end
    btn_cancel:AddClickListener(closeWin)
    btn_close:AddClickListener(closeWin)

end
-- endregion
