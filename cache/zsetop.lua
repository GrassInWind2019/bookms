-- redis zset operations
-- argv[capacity maxScore newMemberScore member]
-- 执行示例 eval zsetop.lua mtest , 3 5 5 test1
-- 获取键和参数
local key,cap,maxSetScore,newMemberScore,member = KEYS[1],ARGV[1],ARGV[2],ARGV[3],ARGV[4]
redis.log(redis.LOG_NOTICE, "key=", key,",cap=", cap,",maxSetScore=", maxSetScore,",newMemberScore=", newMemberScore,",member=", member)
local len = redis.call('zcard', key);
-- len need not nil, otherwise will occur "attempt to compare nil with number"
if len then
    if tonumber(len) >= tonumber(cap)
    then
        local num = tonumber(len)-tonumber(cap)+1
        local list = redis.call('zrangebyscore',key,0,maxSetScore,'limit',0,num)
        redis.log(redis.LOG_NOTICE,"key=",key,"maxSetScore=",maxSetScore, "num=",num)
        for k,lowestScoreMember in pairs(list) do
            local lowestScore = redis.call('zscore', key,lowestScoreMember)
            redis.log(redis.LOG_NOTICE, "list: ", lowestScore, lowestScoreMember)
            if tonumber(newMemberScore) > tonumber(lowestScore)
            then
                local rank = redis.call('zrevrank',key,member)
                -- rank is nil indicate new member is not exist in set, need remove the lowest score member
                if not rank then
                    local index = tonumber(len) - tonumber(cap);
                    -- if list has more than 2 items, zremrangebyrank will remove all items and the second round lowestScore will be nil and error
                    -- "user_script:15: attempt to compare nil with number" occured
                    redis.call('zremrangebyrank',key, 0, index)
                    -- redis.call('zrem',KEYS[1], lowestScoreMember)
                end
                redis.call('zadd', key, newMemberScore, member);
                break
            end
        end
    else
        -- set is not full yet, add new member directly
        redis.call('zadd', key, newMemberScore, member);
    end
end