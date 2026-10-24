-- 模块定义
function mmodule(mname)
    local M = { }
    _G[mname] = M
    package.loaded[mname] = M
    setmetatable(M, { __index = _G })
    setfenv(2, M)
end