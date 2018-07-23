--
-- Created by IntelliJ IDEA.
-- User: barry
-- Date: 18/6/4
-- Time: 19:10
-- To change this template use File | Settings | File Templates.
--


-- 全局变量
main = {
    w = 9
}

function hello(name, value)
    local t = Table("h")
    print(name, value, t)
    r.set("key", "value")
    r.hmset("key", { a = 1, b = 2 })
    print(main.w)
end

function test(name, value)
    print(name, value)
end
