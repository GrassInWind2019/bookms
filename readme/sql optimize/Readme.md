# sql 优化操作流程  
## 需要修改的代码部分  
```
1.将routers/router.go中beego.Router("/usercenter/:page", &controllers.UserController{}, "*:GetUserCenter")改为
   beego.Router("/usercenter/:page", &controllers.UserController{}, "*:GetUserCenterInfo")
2.将main.go中logs.SetLevel(logs.LevelDebug)改为
  logs.SetLevel(logs.LevelWarn)
3.将conf/app.conf中runmode = dev注释或删除，排除beego的log对测试的影响
  将db_uw_database=bookms_user改为db_uw_database=bookms
  将db_ur_database=bookms_user改为db_ur_database=bookms
```
## 创建bookms所需数据库及表 
先通过optimize.sql创建数据库及表，执行完毕后，管理员账号为GrassInWind, 密码123 
然后运行bookms（Goland/liteIDE 运行main即可）测试端口号为8080 
浏览器打开http://localhost:8080/-->登录（管理员账号）-->http://localhost:8080/addcategory添加如下分类： 
（测试只需随便添加分类超过10种即可）
```
# id, pid, category_name, description, book_count, sort, status, icon
'2', '0', '文学', '文学类书籍', '0', '0', '0', ''
'3', '0', '政史', '政史类书籍', '0', '1', '0', ''
'4', '0', '技术', '技术类书籍', '0', '2', '0', ''
'5', '0', '经管', '经管类书籍', '0', '3', '0', ''
'6', '2', '小说', '小说', '505305', '4', '0', ''
'7', '2', '散文', '散文', '506799', '5', '0', ''
'8', '2', '诗词', '诗词', '504519', '6', '0', ''
'9', '4', '前端', '前端书籍', '505350', '10', '0', ''
'10', '4', '后端', '后端书籍', '501819', '11', '0', ''
'11', '4', '其他技术', '其他技术类书籍', '506208', '15', '0', ''
```
## 构造测试数据集 
执行book_user_favorite.sql生成100万测试数据（执行非常耗时，我的机器AMD Ryzen 5 3500U处理器，4核12G内存超过9小时） 
## 通过apachebench工具进行压测 
压测命令见下图： 
其中-c代表并发数，-n代表测试次数，-C表示用户cookie（可通过谷歌浏览器按下f12登录后即可查看到cookie信息） 
![压测示例.png](https://github.com/GrassInWind2019/bookms/blob/master/readme/sql%20optimize/usercenterfav%202020-02-01%20111424.png)

## 收藏图书信息sql优化（简化版） 
### 优化前 
API: /usercenter/201 
sql示例： 
```
explain select b.* from bookms_book b left join bookms_user_favorite f on f.identify=b.identify where user_id=1 limit 100 offset 20000;
+----+-------------+-------+------------+--------+---------------+----------+---------+------+--------+----------+-----------------------+
| id | select_type | table | partitions | type   | possible_keys | key      | key_len | ref  | rows   | filtered | Extra                 |
+----+-------------+-------+------------+--------+---------------+----------+---------+------+--------+----------+-----------------------+
|  1 | SIMPLE      | f     | NULL       | ALL    | NULL          | NULL     | NULL    | NULL | 202810 |    10.00 | Using where           |
|  1 | SIMPLE      | b     | NULL       | eq_ref | identify      | identify | 402     | func |      1 |   100.00 | Using index condition |
+----+-------------+-------+------------+--------+---------------+----------+---------+------+--------+----------+-----------------------+
2 rows in set, 1 warning (0.00 sec)
```
apache bench压测结果：QPS为17  
```
Document Path:          /usercenter/201
Document Length:        27886 bytes

Concurrency Level:      10
Time taken for tests:   58.626 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      28112000 bytes
HTML transferred:       27886000 bytes
Requests per second:    17.06 [#/sec] (mean)
Time per request:       586.264 [ms] (mean)
Time per request:       58.626 [ms] (mean, across all concurrent requests)
```

### 优化1 
API:/usercenter2/201 
通过将关联查询拆分为两个简单查询来提升查询效率。 
sql示例：
```
select identify from bookms_user_favorite where user_id=xxx limit 100 offset 20000
select * from bookms_book where identify in ()
```
apache bench压测结果：QPS达到了303  
```
Document Path:          /usercenter2/201
Document Length:        27980 bytes

Concurrency Level:      10
Time taken for tests:   3.298 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      28178000 bytes
HTML transferred:       27980000 bytes
Requests per second:    303.26 [#/sec] (mean)
Time per request:       32.975 [ms] (mean)
Time per request:       3.298 [ms] (mean, across all concurrent requests)
```

### 优化2 
API: /usercenter3/201 
通过子查询现将收藏的图书标识查询出来再做关联查询来提升查询效率。
sql示例： 
```
mysql> explain select * from bookms_book inner join(select identify from bookms_user_favorite where user_id=1 limit 100 offset 20000) as book_fav using(identify);
+----+-------------+----------------------+------------+--------+---------------+----------+---------+-------+--------+----------+-----------------------+
| id | select_type | table                | partitions | type   | possible_keys | key      | key_len | ref   | rows   | filtered | Extra                 |
+----+-------------+----------------------+------------+--------+---------------+----------+---------+-------+--------+----------+-----------------------+
|  1 | PRIMARY     | <derived2>           | NULL       | ALL    | NULL          | NULL     | NULL    | NULL  |  20100 |   100.00 | NULL                  |
|  1 | PRIMARY     | bookms_book          | NULL       | eq_ref | identify      | identify | 402     | func  |      1 |   100.00 | Using index condition |
|  2 | DERIVED     | bookms_user_favorite | NULL       | ref    | user_id       | user_id  | 4       | const | 101405 |   100.00 | NULL                  |
+----+-------------+----------------------+------------+--------+---------------+----------+---------+-------+--------+----------+-----------------------+
3 rows in set, 1 warning (0.00 sec)
```
apache bench压测结果：QPS达到了344  
```
Document Path:          /usercenter3/201
Document Length:        27980 bytes

Concurrency Level:      10
Time taken for tests:   2.902 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      28178000 bytes
HTML transferred:       27980000 bytes
Requests per second:    344.54 [#/sec] (mean)
Time per request:       29.024 [ms] (mean)
Time per request:       2.902 [ms] (mean, across all concurrent requests)
```
