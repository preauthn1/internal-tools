require("Module")
mmodule(...)

local G = require "Global"
local Json = require("Tools/Json")
local String = require("string");
local Lockbox = require("Tools/lockbox/lockbox")
local Array = require("Tools/lockbox/util/array")
local Stream = require("Tools/lockbox/util/stream")

Lockbox.ALLOW_INSECURE = true;

--[[
local PKCS7Padding = require("Tools/lockbox/padding/pkcs7")
local ZeroPadding = require("Tools/lockbox/padding/zero")
]]--
local Padding = require("Tools/lockbox/padding/pkcs7")


local AES128Cipher = require("Tools/lockbox/cipher/aes128")

local key = "V1M5WE7X8WE2TWEK"
local iv = "DNW123EXGMTR83K7"

--[[
local ECBMode = require("Tools/lockbox/cipher/mode/ecb")
local CBCMode = require("Tools/lockbox/cipher/mode/cbc")
]]--
local mode = require("Tools/lockbox/cipher/mode/cbc")

function ShowLog(Text)
	--print(Text)
end

function str2hex(str)
	local hex = ""
	for i=0, string.len(str)-1, 1 do
		local k = i+1
		local b = string.byte(str, k, k)
		hex = hex..string.format("%02x", b)
	end
	return hex
end

function hex2str(hex)
	local str, n = hex:gsub("(%x%x)[ ]?", function (word)
		return string.char(tonumber(word, 16))
	end)
	return str
end

function check_hex_lua(str)
    if "string" ~= type(str) then
        return false
    end

    for i = 1, #str do
        local ord = str:byte(i)
        if not (
            (48 <= ord and ord <= 57) or
            (65 <= ord and ord <= 70) or
            (97 <= ord and ord <= 102)
        ) then
            return false
        end
    end
    return true
end

function Encode(data)
	local aes = mode.Cipher()
	aes.setKey(Stream.toArray(Stream.fromString(key)))
	aes.setBlockCipher(AES128Cipher)
	aes.setPadding(Padding)
	aes.init() 
	aes.update(Stream.fromString(iv))
	aes.update(Stream.fromString(Json.encode(data)))
	aes.finish()
	ShowLog("Encode Json\n" .. Json.encode(data))	
	hexData = aes.asHex()
	ShowLog("Encode\n" .. hexData)	
	return hexData
end

function Decode(data)
	ShowLog("Decode data\n" .. data)
	if check_hex_lua(data) == false then
		print("######## Decode error !!回傳資料資料格式不是Hex")
		return nil
	end
	local aes = mode.Decipher()
	aes.setKey(Stream.toArray(Stream.fromString(key)))
	aes.setBlockCipher(AES128Cipher)
	aes.setPadding(Padding)
	aes.init() 
	aes.update(Stream.fromString(iv))
	aes.update(Stream.fromHex(data))
	aes.finish()
	k = aes.asHex()
	ShowLog("Decode hex data\n" .. k)
	strData = hex2str(k)
	ShowLog("Decode Json\n" .. strData)
	local decodeJsonData = Json.decode(strData .. "\n")
	return decodeJsonData
end

function DecodeNoJson(data)
	ShowLog("Decode data\n" .. data)
	if check_hex_lua(data) == false then
		print("######## Decode error !!回傳資料資料格式不是Hex")
		return nil
	end
	local aes = mode.Decipher()
	aes.setKey(Stream.toArray(Stream.fromString(key)))
	aes.setBlockCipher(AES128Cipher)
	aes.setPadding(Padding)
	aes.init() 
	aes.update(Stream.fromString(iv))
	aes.update(Stream.fromHex(data))
	aes.finish()
	k = aes.asHex()
	ShowLog("Decode hex data\n" .. k)
	strData = hex2str(k)
	ShowLog("Decode\n" .. strData)
	local decodeJsonData = strData .. "\n"
	return decodeJsonData
end