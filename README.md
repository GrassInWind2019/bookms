# bookms
bookms是一个简易的图书管理系统。  
## bookms架构示意图
![bookms架构示意图](https://raw.githubusercontent.com/GrassInWind2019/bookms/master/readme/bookms%E6%9E%B6%E6%9E%84%E7%A4%BA%E6%84%8F%E5%9B%BE.png)
## bookms功能示意图 
![bookms功能示意图](https://raw.githubusercontent.com/GrassInWind2019/bookms/master/readme/bookms%E5%8A%9F%E8%83%BD%E7%A4%BA%E6%84%8F%E5%9B%BE.png)

## bookms图书评分排行  
使用redis的zset保存排行数据，使用lua脚本实现评分排行更新的原子操作。  
lua脚本如下：此脚本用于在用户对书评分后更新书的评分数据。KEYS[1]表示zset的key，ARGV[1]表示期望的zset最大存储成员数量，ARGV[2]表示评分上限，默认评分下限是0，ARGV[3]表示待添加的评分，ARGV[4]表示待添加的成员名称。  
相关redis命令：  
ZCARD key   获取有序集合的成员数  
ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT]  通过分数返回有序集合指定区间内的成员(从小到大的顺序)  
ZREMRANGEBYRANK key start stop 移除有序集合中给定的排名区间的所有成员  
ZADD key score1 member1 [score2 member2]  向有序集合添加一个或多个成员，或者更新已存在成员的分数  
```
-- redis zset operations
-- argv[capacity maxScore newMemberScore member]
-- 执行示例 redis-cli.exe --eval zsetop.lua mtest , 3 5 5 test1
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
```
在cache/cache.go的init函数中读取Lua脚本并通过redisgo包的NewScript函数加载这个脚本，在使用时通过返回的指针调用lua.Do()即可。
```
func init() {
    ...
	file, err := os.Open(zsetopFileName)
	if err != nil {
		panic("open"+zsetopFileName+" "+err.Error())
	}
	bytes,err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}
	zsetopScript = utils.UnsafeBytesToString(bytes)
	logs.Debug(zsetopScript)
	lua =redis.NewScript(1,zsetopScript)
}
func ZaddWithCap(key,member string, score float32, maxScore, cap int) (reply interface{}, err error) {
	c := pool.Get()
	//Do optimistically evaluates the script using the EVALSHA command. If script not exist, will use eval command.
	reply, err= lua.Do(c,key,cap,maxScore,score,member)
	return
}
```

## bookms使用简介 
### 下载code  
1. git clone git@github.com:GrassInWind2019/bookms.git  
2. 下载依赖库  
  go mod download  
3. 下载mysql5.7以后版本并安装  
  下载地址： https://dev.mysql.com/downloads/mysql/  
4. 下载redis并安装  
    windows版本下载地址：https://github.com/MicrosoftArchive/redis/releases 
    选择3.2.100版本  
5. 下载nginx并安装(运行本项目不是必须的)    
    下载地址： http://nginx.org/en/download.html  
### 创建数据库及表 
通过一些工具如Navicat/Mysql Workbench等执行book.sql和user.sql即可创建。执行完毕，管理员账号为GrassInWind，密码123 
### 配置redis 
默认redis连接密码为123，可通过conf/app.conf修改 
```
#session
sessionon = true
ProviderConfig=localhost:6379,1000,123
#redis
server=127.0.0.1:6379
password=123
max_active_connection=1000
max_idle_connection=50
db=0
```
### 添加图书分类 
在登录后访问http://localhost:8080/addcategory 即可添加分类 
### 添加图书 
管理员在个人中心页面可通过添加图书按钮添加 
### nginx配置 
nginx配置文件见nginx_conf目录，本人所用nginx版本为nginx-1.16.1。 
通过nginx -t检查配置文件有没有问题。  
### nginx使用 
1. 为了使用nginx,需要以不同的端口运行多个bookms。使用beego的bee工具通过在项目目录下执行bee pack即可完成打包bookms运行所需要的文件。  
2. 将bookms.tar.gz拷贝到两个目录并解压，修改conf/app.conf中“httpport = 8080”为不同的端口如8081,8082，然后直接运行bookms.exe。  
3. 运行nginx（nginx配置端口8080），然后访问http://localhost:8080 即可访问bookms。  
```
dongs@LAPTOP-V7V47H0L MINGW64 /d/GoCode/study/bookms (master)
$ bee pack
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0
2020/02/28 11:35:07 INFO     ▶ 0001 Packaging application on 'D:\GoCode\study\bookms'...
2020/02/28 11:35:07 INFO     ▶ 0002 Building application...
2020/02/28 11:35:07 INFO     ▶ 0003 Using: GOOS=windows GOARCH=amd64
2020/02/28 11:35:12 SUCCESS  ▶ 0004 Build Successful!
2020/02/28 11:35:12 INFO     ▶ 0005 Writing to output: D:\GoCode\study\bookms\bookms.tar.gz
2020/02/28 11:35:12 INFO     ▶ 0006 Excluding relpath prefix: .
2020/02/28 11:35:12 INFO     ▶ 0007 Excluding relpath suffix: .go:.DS_Store:.tmp
2020/02/28 11:35:16 SUCCESS  ▶ 0008 Application packed!
```
## bookms收藏功能sql优化 
### 优化前 
API: /usercenterfav/201  
sql示例： 
```
select id,user_id,identify from bookms_book_record where user_id=1 limit 100 offset 20000  
mysql> explain select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,bc.category_id,c.category_name,
(case when r.lend_status=0 then '可借' 
when r.lend_status=5 then '已下架' 
when r.lend_status=1 and r.user_id=1 then '正在借阅' 
when r.lend_status=1 and r.user_id<>1 then '不可借' 
end) as lend_status from bookms_book_record r 
left join bookms_book b using(identify) 
left join bookms_book_category bc using(identify) 
inner join bookms_user_favorite f using(identify) 
left join bookms_category c on bc.category_id=c.id where f.user_id=1 limit 100 offset 20000;
+----+-------------+-------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+--------+----------+-------------+
| id | select_type | table | partitions | type   | possible_keys                     | key                               | key_len | ref                   | rows   | filtered | Extra       |
+----+-------------+-------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+--------+----------+-------------+
|  1 | SIMPLE      | f     | NULL       | ALL    | idx_bookms_user_favorite_identify | NULL                              | NULL    | NULL                  | 202810 |    10.00 | Using where |
|  1 | SIMPLE      | r     | NULL       | ref    | idx_bookms_book_record_identify   | idx_bookms_book_record_identify   | 302     | bookms.f.identify     |      1 |   100.00 | NULL        |
|  1 | SIMPLE      | b     | NULL       | eq_ref | identify                          | identify                          | 402     | func                  |      1 |   100.00 | Using where |
|  1 | SIMPLE      | bc    | NULL       | ref    | idx_bookms_book_category_identify | idx_bookms_book_category_identify | 302     | bookms.f.identify     |      1 |   100.00 | NULL        |
|  1 | SIMPLE      | c     | NULL       | eq_ref | PRIMARY                           | PRIMARY                           | 4       | bookms.bc.category_id |      1 |   100.00 | NULL        |
+----+-------------+-------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+--------+----------+-------------+
5 rows in set, 1 warning (0.00 sec)
```
apachebench压测结果：QPS约为1  
```
Document Path:          /usercenterfav/201
Document Length:        21876 bytes

Concurrency Level:      10
Time taken for tests:   91.497 seconds
Complete requests:      100
Failed requests:        0
Total transferred:      2207400 bytes
HTML transferred:       2187600 bytes
Requests per second:    1.09 [#/sec] (mean)
Time per request:       9149.661 [ms] (mean)
Time per request:       914.966 [ms] (mean, across all concurrent requests)
```
### 优化1 
API: /usercenterfav3/201  
先通过子查询(select fav.id,fav.user_id,fav.identify from bookms_user_favorite fav where fav.user_id=1 limit 100 offset 20000)查出收藏的图书的标识信息再关联查询。  
sql示例：  
```
mysql> explain select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,bc.category_id,c.category_name,
(case when r.lend_status=0 then '可借' 
when r.lend_status=5 then '已下架' 
when r.lend_status=1 and r.user_id=1 then '正在借阅' 
when r.lend_status=1 and r.user_id<>1 then '不可借' 
end) as lend_status from bookms_book_record r 
left join bookms_book b using(identify) 
left join bookms_book_category bc using(identify) 
inner join (select fav.id,fav.user_id,fav.identify from bookms_user_favorite fav where fav.user_id=1 limit 100 offset 20000) f using(user_id) 
left join bookms_category c on bc.category_id=c.id where f.user_id=1 limit 100 offset 20000;
+----+-------------+------------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+---------+----------+----------------------------------------------------+
| id | select_type | table      | partitions | type   | possible_keys                     | key                               | key_len | ref                   | rows    | filtered | Extra                                              |
+----+-------------+------------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+---------+----------+----------------------------------------------------+
|  1 | PRIMARY     | <derived2> | NULL       | ref    | <auto_key0>                       | <auto_key0>                       | 4       | const                 |      10 |   100.00 | NULL                                               |
|  1 | PRIMARY     | r          | NULL       | ALL    | NULL                              | NULL                              | NULL    | NULL                  | 1004634 |    10.00 | Using where; Using join buffer (Block Nested Loop) |
|  1 | PRIMARY     | b          | NULL       | eq_ref | identify                          | identify                          | 402     | func                  |       1 |   100.00 | Using where                                        |
|  1 | PRIMARY     | bc         | NULL       | ref    | idx_bookms_book_category_identify | idx_bookms_book_category_identify | 302     | bookms.r.identify     |       1 |   100.00 | NULL                                               |
|  1 | PRIMARY     | c          | NULL       | eq_ref | PRIMARY                           | PRIMARY                           | 4       | bookms.bc.category_id |       1 |   100.00 | NULL                                               |
|  2 | DERIVED     | fav        | NULL       | ALL    | NULL                              | NULL                              | NULL    | NULL                  |  202810 |    10.00 | Using where                                        |
+----+-------------+------------+------------+--------+-----------------------------------+-----------------------------------+---------+-----------------------+---------+----------+----------------------------------------------------+
6 rows in set, 1 warning (0.00 sec)
```
apachebench压测结果：QPS相比优化前提升了近40倍  
```
Document Path:          /usercenterfav3/201
Document Length:        21733 bytes

Concurrency Level:      10
Time taken for tests:   24.206 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      21931000 bytes
HTML transferred:       21733000 bytes
Requests per second:    41.31 [#/sec] (mean)
Time per request:       242.064 [ms] (mean)
Time per request:       24.206 [ms] (mean, across all concurrent requests)
```
### 优化2 
API: /usercenterfav2/201  
通过将多表关联查询拆分成多个简单查询来提升查询效率。  
sql示例：  
```
select (case 
when r.lend_status=0 then '可借' 
when r.lend_status=5 then '已下架' 
when r.lend_status=1 and r.user_id=1 then '正在借阅' 
when r.lend_status=1 and r.user_id<>1 then '不可借' 
end) as lend_status,identify from bookms_book_record r where r.identify in(...)

select book_name,cover,author,identify from bookms_book where identify in(...)

select category_id,identify from bookms_book_category where identify in(...)

select category_name from bookms_category where id in(...)

```
apachebench压测结果：QPS提升了近170倍  
```
Document Path:          /usercenterfav2/201
Document Length:        21168 bytes

Concurrency Level:      10
Time taken for tests:   5.599 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      21366000 bytes
HTML transferred:       21168000 bytes
Requests per second:    178.62 [#/sec] (mean)
Time per request:       55.986 [ms] (mean)
Time per request:       5.599 [ms] (mean, across all concurrent requests)
```
### bookms收藏功能优化操作流程 
详见https://github.com/GrassInWind2019/bookms/tree/master/readme/sql%20optimize 

## 开发备忘录 
详见 https://github.com/GrassInWind2019/bookms/blob/master/readme/develop%20memo/Readme.md 
