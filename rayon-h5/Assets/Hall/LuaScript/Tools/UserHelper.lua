-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成
require("Module")
mmodule(...)

local G = require "Global"
local coinShowRate = 100
local function clearSuffix(str)
    if (string.find(str, ".") > 0) then
        str = string.gsub(str, "0+$", "")
        str = string.gsub(str, "%.$", "")
    end
    return str
end
-- 显示完整金币数
function ToCashDetail(num, isConvert)
    num = tonumber(num)
    if (isConvert == nil or isConvert == true) then
        num = num / coinShowRate
    end
    if (coinShowRate > 1) then
        num = G.Tool.Math.Decimal(num, 2)
        return string.format("%.2f", num)
    else
        return string.format("%s", G.Tool.Math.ParseInt(num * coinShowRate))
    end
end

-- 显示简短金币数
function ToCash(num, isConvert)

    num = tonumber(num)
    if (isConvert ~= false) then
        num = num / coinShowRate
    end
    local pointNum = 0
    if (coinShowRate > 1) then
        pointNum = 2
        num = G.Tool.Math.Decimal(num, 2)
    else
        num = G.Tool.Math.ParseInt(num)
    end
    if (num < 100000) then
        return clearSuffix(string.format("%." .. pointNum .. "f", num))
    elseif (num < 100000000) then
        num = G.Tool.Math.Decimal(num / 10000, 2)
        num = clearSuffix(string.format("%2f", num))
        return string.format("%s万", num);
    else
        num = G.Tool.Math.Decimal(num / 100000000, 2)
        num = clearSuffix(string.format("%2f", num))
        return string.format("%s亿", num);
    end
end
-- 转换为单位的货币显示
function ToCashUnit(num)
    return ToCash(num) .. "金币"
end

-- 获取保存的金币数
function GetSaveCoin(num)
    num = tonumber(num)
    return G.Tool.Math.ParseInt(num * coinShowRate)
end

--[[
Example usage:							Output:
print(FormatNum(9))						9
print(FormatNum(999))					999
print(FormatNum(1000))					1,000
print(FormatNum('1333444.10'))			1,333,444.10
print(FormatNum('US$1333400'))			US$1,333,400
print(FormatNum('-$22333444.56'))		-$22,333,444.56
print(FormatNum('($22333444.56)'))		($22,333,444.56)
print(FormatNum('NEG $22333444.563'))	NEG $22,333,444.563
]]
function FormatNum(n) -- credit http://richard.warburton.it
	local left,num,right = string.match(n,'^([^%d]*%d)(%d*)(.-)$')
	return left..(num:reverse():gsub('(%d%d%d)','%1,'):reverse())..right
end


function sec2time(t)
    local tt=t;
    local str=""

    if(tt>86400)then
       local d= math.floor(tt/86400);  
       tt= tt-(86400*d);
       str=str..d.."天";
    end
    if(tt>3600)then
       local h= math.floor(tt/3600);  
       tt= tt-(3600*h);
       str=str..h.."时";
    end
    if(tt>60)then
       local m= math.floor(tt/60);  
       tt= tt-(60*m);
       str=str..m.."分";
    end
    str=str..tt.."秒";
   
    return str;
end


-- endregion
