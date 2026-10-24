require("Module")
mmodule(...)

--[[
    url解码
    @access public
    @param string 要处理的字符串
    @return string
]]
function urldecode(string)
    local decode = string.gsub(string, '%%(%x%x)', function(hex)
        return string.char(tonumber(hex, 16))
    end)
    return decode
end

--[[
    url编码
    @access public
    @param string 要处理的字符串
    @return string
]]
function urlencode(string)
    local encode = string.gsub(string, '([^%w%.%- ])', function(chr)
        return string.format('%%%02X', string.byte(chr))
    end)
    return encode
end
