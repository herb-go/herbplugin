function start()
    print("printed")
end
function getparam(name)
    return system.getparam(name)
end

m=require("module")
m.Print()
m=dofile("../testscripts/module.lua")
m.Print()
f=loadfile("module.lua")
m=f()
m.Print()
