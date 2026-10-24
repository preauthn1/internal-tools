require("Module")
mmodule(...)

local G = require "Global"
local groupName = "BindingTestSerice";

local BingdingPhone = false     -- 公告队列
local BindingZFDone = false
local BindingUPDone = false
local TestMode=false

function IsBindingTestMode()
 return TestMode;
end

function IsBingdingPhone()
 return BingdingPhone;
end
function IsBindingZFDone()
 return BindingZFDone;
end
function IsBindingUPDone()
 return BindingUPDone;
end



          

