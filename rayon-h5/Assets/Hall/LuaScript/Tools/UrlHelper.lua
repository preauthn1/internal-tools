require("Module")
mmodule(...)
local G = require "Global"
local yield_return = G.Util.async_to_sync
(
function(to_to_yield, cb)
    CS.TimeManager.Instance:YieldAndCallback(to_to_yield, cb)
end
)

--获取页面内容
function GetPageData(url,onSuccess, onError)
 G.Util.coroutine_call(function ()
   local www = CS.UnityEngine.WWW(url)
   yield_return(www)

   if www.error ~= nil then
        if onError ~= nil  then
            onError(www.error)
        end

        return 
   end

   if onSuccess ~= nil then
        onSuccess(www.text)
   end

end)()

end
