
-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成

require("Module")
mmodule(...)


function Print(...)
    arg = { ...}

    if arg == nil or #arg == 0 then
        print("table is nil.")
        return
    end

    local print_r_cache = { }
    local printStr = ""

    local t = arg
    if #arg > 1 then
        printStr = arg[1] .. ":"
        t = arg[2]
    end


    local function sub_print_r(t, indent)
        if (print_r_cache[tostring(t)]) then
            printStr = printStr .. "\n" ..(indent .. "*" .. tostring(t))
        else
            print_r_cache[tostring(t)] = true
            if (type(t) == "table") then
                for pos, val in pairs(t) do
                    if (type(val) == "table") then
                        printStr = printStr .. "\n" ..(indent .. "[" .. pos .. "] => " .. tostring(t) .. " {")
                        sub_print_r(val, indent .. string.rep(" ", string.len(pos) + 8))
                        printStr = printStr .. "\n" ..(indent .. string.rep(" ", string.len(pos) + 6) .. "}")
                    elseif (type(val) == "string") then
                        printStr = printStr .. "\n" ..(indent .. "[" .. pos .. '] => "' .. val .. '"')
                    else
                        printStr = printStr .. "\n" ..(indent .. "[" .. pos .. "] => " .. tostring(val))
                    end
                end
            else
                printStr = printStr .. "\n" ..(indent .. tostring(t))
            end
        end
    end
    if (type(t) == "table") then
        printStr = printStr .. "\n" ..(tostring(t) .. " {")
        sub_print_r(t, "  ")
        printStr = printStr .. "\n" ..("}")
    else
        sub_print_r(t, "  ")
    end
	
	print(printStr)
end

function IsSame(obj1, obj2)
    local isSame = true
    if (obj1 == nil) then obj1 = { } end
    if (obj2 == nil) then obj2 = { } end

    for k, v in pairs(obj1) do
        if (v ~= obj2[k]) then
            isSame = false
            break
        end
    end
    return isSame
end

function Find(tblist, fun)
    local o = { }
    for k, v in pairs(tblist) do
        if (fun(v)) then table.insert(o, v) end
    end
    return o
end

function First(tblist, fun)
    for k, v in pairs(tblist) do
        if (type(fun) == "function") then
            if (fun(v)) then return v end
        else
            return v
        end
    end
    return nil
end

function Exists(tblist, fun)
    for k, v in pairs(tblist) do
        if (type(fun) == "function") then
			if (fun(v)) then
				return true
			end
		else
			if (fun == v) then
				return true
			end
        end
    end
    return false
end

function Count(tblist, fun)
    local count = 0
    for k, v in pairs(tblist) do
        if (type(fun) == "function") then
			if (fun(v)) then count = count + 1 end
		else
			if (fun == v) then count = count + 1 end
		end
    end
    return count
end

function Sort(tblist, fun)
    local newtblist = { }
    for _, v in pairs(tblist) do
        table.insert(newtblist, v)
    end

    table.sort(newtblist, fun)
    return newtblist
end

function Clone(obj)
    local fun = function(obj)
        if (type(obj) ~= "table") then
            return obj
        else
            local tb = { }
            for k, v in pairs(obj) do
                tb[k] = Clone(v)
            end
            return tb
        end
    end
    return fun(obj)
end

function Pull(tblist, fun)
    for k, v in pairs(tblist) do
        if (type(fun) == "function") then
            if (fun(v)) then
                tblist[k] = nil
                return v
            end
        else
            tblist[k] = nil
            return v
        end
    end
    return nil
end
function ContainKey(tblist, fun)
    local o = { }
    for k, v in pairs(tblist) do
        if (fun(k)) then table.insert(o, k) end
    end
    return o
end

function Merge(first_table,second_table)
   for k,v in pairs(second_table) do first_table[k] = v end
   return first_table;
end

-- endregion

