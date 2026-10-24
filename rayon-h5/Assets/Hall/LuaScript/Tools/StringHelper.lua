require("Module")
mmodule(...)

local G = require "Global";

-- 删除所有字符串中的空格
function TrimAll(str)
    return string.gsub(str, "%s+", "")
end

-- 删除字符串首尾空格
function Trim(str)
    return string.gsub(str, "^%s*(.-)%s*$", "%1")
end

-- 删除字符串首空格
function TrimLeft(str)
    return string.gsub(str, "^%s*(.-)", "%1")
end

-- 删除字符串尾空格
function TrimRight(str)
    return string.gsub(str, "(.-)%s*$", "%1")
end

-- 分隔字符串返回一个Table
function Split(str, splitStr)
    if (str == nil or str == "") then
        return nil
    end

    local strAry = { }
    local afterStr = splitStr
    if (splitStr == ".") then
        splitStr = "%."
    end
    for v in string.gmatch(str .. afterStr, "(.-)" .. splitStr) do
        table.insert(strAry, v)
    end
    return strAry
end

-- 分隔並清除空白 By Rayon
function NewSplit(str, splitStr)
    if (str == nil or str == "") then
        return nil
    end

	local strAry={}
	for v in string.gmatch(str, "([^"..splitStr.."]+)") do
		str1 = string.gsub(v, "%s", "")
		table.insert(strAry, str1)
	end
	return strAry
end

-- 检查字符串是否为空
function IsEmpty(str)
    if str == nil or str == "" then
        return true
    end

    return TrimAll(str) == ""
end

function Sub(str, startIndex, endIndex)
    return string.sub(str, startIndex, endIndex)
end

function SubString(str, startIndex, len)
    return string.sub(str, startIndex, startIndex + len - 1)
end

-- 解析当前平台游戏的版本号 {平台号}_{版本号} 1_1.0.1,2_1.0.1
function ParseVersion(versionStr, platformID)
    local platformVersions = Split(versionStr, ",");

    platformID = tostring(platformID)

	-- 修正判斷版號錯誤
    --for i = 1, #platformVersions - 1, 1 do
    for i = 1, #platformVersions, 1 do
        local platform_Version = Split(platformVersions[i], "_");
        if platform_Version[1] == platformID then
            return platform_Version[2];
        end
    end

	print("!!!!!!!!!!!!!!!!!!!!!! 解析当前平台游戏的版本号错误 !!")
	
    return "99.0.0"
end

function chSize(char)  
    if not char then  
        return 0  
    elseif char > 240 then  
        return 4  
    elseif char > 225 then  
        return 3  
    elseif char > 192 then  
        return 2  
    else  
        return 1  
    end  
end  

function IsUpgrade_v1(v,HD)
   
    local arrv=Split(v,".");
		
    if (HD==1) then
		if G.GameDefine.GamePlatformID == 80 then
			return true
		else
			return  (tonumber(arrv[2]) ~= 0);
		end
    elseif (HD==2) then
        return  (tonumber(arrv[2]) ~= 0);
    else
        return true;
    end
end
  
  

  
-- 字數計算
function length(str)  
    local length = 0;  
    local currentIndex = 1;  
    while currentIndex <= #str do  
        local char = string.byte(str, currentIndex)  
        currentIndex = currentIndex + chSize(char)  
        length = length + 1  
    end  
    return length;  
end  
  

-- 字數計算
function string.utf8len(input)
    local len  = string.len(input)
    local left = len
    local cnt  = 0
    local arr  = {0, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc}
    while left ~= 0 do
        local tmp = string.byte(input, -left)
        local i   = #arr
        while arr[i] do
            if tmp >= arr[i] then
                left = left - i
                break
            end
            i = i - 1
        end
        cnt = cnt + 1
    end
    return cnt
end

function GetIPType(ip)
    -- must pass in a string value
    if ip == nil or type(ip) ~= "string" then
        return 0
    end

    -- check for format 0.0.0.0 for ipv4
    local chunks = {ip:match("(%d+)%.(%d+)%.(%d+)%.(%d+)")}
    if (#chunks == 4) then
        for _,v in pairs(chunks) do
            if (tonumber(v) < 0 or tonumber(v) > 255) then
                return 0
            end
        end
        return 1
    else
        return 0
    end

    -- check for ipv6 format, should be 8 'chunks' of numbers/letters
    local _, chunks = ip:gsub("[%a%d]+%:?", "")
    if chunks == 8 then
        return 2
    end

    -- if we get here, assume we've been given a random string
    return 3
end

