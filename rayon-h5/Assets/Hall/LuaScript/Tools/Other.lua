require("Module")
mmodule(...)

--[[
    url解码
    @access public
    @param string 要处理的字符串
    @return string
]]
function urldecode(data)
	
    local decode = string.gsub(data, '%%(%x%x)', function(hex)
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
function urlencode(data)
	
    local encode = string.gsub(data, '([^%w%.%- ])', function(chr)
        return string.format('%%%02X', string.byte(chr))
    end)

    return encode
end
