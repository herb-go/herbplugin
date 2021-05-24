plugin=require("plugin")
function start()
    print("123")
end
function getparam(name)
    return plugin.getparam(name)
end