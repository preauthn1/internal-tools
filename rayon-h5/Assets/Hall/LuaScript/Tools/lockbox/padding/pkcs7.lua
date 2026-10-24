local Stream = require("Tools/lockbox/util/stream");

local PKCS7Padding = function(blockSize,byteCount)

	-- 修改 https://example.invalid/p/c49b9768a6e9
	local paddingCount = blockSize - (byteCount % blockSize);
	--local paddingCount = blockSize - ((byteCount -1) % blockSize) + 1;
	local bytesLeft = paddingCount;

	local stream = function()
		if bytesLeft > 0 then
			bytesLeft = bytesLeft - 1;
			return paddingCount;
		else
			return nil;
		end
	end

	return stream;
end

return PKCS7Padding;
