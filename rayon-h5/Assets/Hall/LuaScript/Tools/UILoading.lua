require("Module")
mmodule(...)

local G = require("Global")


function ShowTxt(tips, time, isshade, callfun)
    local uiData = CS.LoadingDataUIData()
    uiData.TipsContent = tips
    uiData.MaxLastTime = time 
    uiData.IsFrozenUI = G.CaseWhen(isshade == nil, false, isshade)
    if (callfun ~= nil) then
        uiData.MaxTimeCloseWindowCallBack = callfun
    end

    CS.UILoadingData.Show(uiData)
end
 
function Hide()    
    CS.UILoadingData.Hide()
end