# bookms
bookms是一个简易的图书管理系统。  
## bookms架构示意图
![bookms架构示意图](https://raw.githubusercontent.com/GrassInWind2019/bookms/master/readme/bookms%E6%9E%B6%E6%9E%84%E7%A4%BA%E6%84%8F%E5%9B%BE.png)
## bookms功能示意图 
![bookms功能示意图](https://raw.githubusercontent.com/GrassInWind2019/bookms/master/readme/bookms%E5%8A%9F%E8%83%BD%E7%A4%BA%E6%84%8F%E5%9B%BE.png)
## bookms收藏功能sql优化 
### 优化前 
API: /usercenterfav/201 
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

## bookms使用简介 
### 创建数据库及表 
执行book.sql和user.sql即可创建。执行完毕，管理员账号为GrassInWind，密码123 
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


