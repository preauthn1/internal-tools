require("Module")
mmodule(...)


function HToColor(value)
    value = string.gsub(value, "#", "")
    local len = string.len(value);
    local rgb = { }
    for i = 1, 8, 2 do
        if (i > len) then
            break
        end
        table.insert(rgb, tonumber(string.sub(value, i, i + 1), 16) / 255)
    end
    if (rgb[4] == nil) then
        table.insert(rgb, 1)
    end

    return CS.UnityEngine.Color(rgb[1], rgb[2], rgb[3], rgb[4])
end

function RGBToColor(r, g, b, a)
    return CS.UnityEngine.Color(r / 255, g / 255, b / 255, a / 255)
end