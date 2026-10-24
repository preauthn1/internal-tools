require("Module")
mmodule(...)

local G = require "Global"

local m_IsQuit = false;

function AppQuit()
    m_IsQuit = true;
end

-- 游戏被切换到后台
local function AppPause(isPause)
    if isPause then
        print("...............游戏暂停...............")
        -- 关闭大厅连接
        G.Service.NetworkService.CloseHallSocket()
        -- 关闭游戏连接
        G.Service.NetworkService.CloseRoomSocket()
    else
        print("...............游戏恢复...............")

        -- 打开下载连接之后，游戏app再次获得焦点，则关闭
        if m_IsQuit then
            CS.UnityEngine.Application.Quit();
        end

        G.Service.NetworkService.ReLogonHall()
        if G.IsGameStart() then
            G.Service.NetworkService.ReLoginRoom()
        end
    end
end

local function AppQuit()
    print("...............游戏退出...............")
end

-- 初始化
function Init()
    CS.AppManager.Instance.OnAppPause = CS.AppManager.Instance.OnAppPause + AppPause;
    CS.AppManager.Instance.OnAppQuit = CS.AppManager.Instance.OnAppQuit + AppQuit;
end
