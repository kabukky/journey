package.path = getCurrentDir() .. "/testmod/?.lua;" .. package.path

TM = require "testmod"

function register()
	return {"test"}
end

function test()
	return TM.foo()
end