# 开发过程备忘录 
## Gob decode 
需传入指针，否则decode结果无法存入目标 
```
if err := utils.Decode(cookie, &c.Muser); err == nil && c.Muser.Id > 0
func Decode(value string, r interface{}) error {
	buff := bytes.NewBuffer([]byte(value))
	dec := gob.NewDecoder(buff)
	return dec.Decode(r)
}
```
## byte切片转string  
通过unsafe包将byte切片转换为string效率比直接使用string([]byte)要高。  
```
func UnsafeBytesToString(bytes []byte) string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return *(*string)(unsafe.Pointer(hdr))
}
```
反向转换如下：  
```
func UnsafeStringToBytes(str string) *[]byte {
	string_header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bytes_header := &reflect.SliceHeader{
		Data:string_header.Data,
		Len: string_header.Len,
		Cap: string_header.Len,
	}
	return (*[]byte)(unsafe.Pointer(bytes_header))
}
```
string与slice底层结构如下： 
```
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
  } 
type StringHeader struct {
	Data uintptr
	Len  int
  }
```
## unsafe包  
unsafe 包提供了 2 点重要的能力：
1. 任何类型的指针和 unsafe.Pointer 可以相互转换。  
2. uintptr 类型和 unsafe.Pointer 可以相互转换。  
![unsafe.Pointer.png](https://github.com/GrassInWind2019/bookms/master/readme/develop%20memo/unsafe.Pointer.png)
### Offsetof 获取成员偏移量  
```
package main
import (
"fmt"
"unsafe"
)
type Programmer struct {
name string
language string
}
func main() {
p := Programmer{"stefno", "go"}
fmt.Println(p)
name := (*string)(unsafe.Pointer(&p))
*name = "qcrao"
lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.language)))
*lang = "Golang"
fmt.Println(p)
}
```
运行代码，输出：  
{stefno go}  
{qcrao Golang}  
### unsafe.Sizeof()   
```
type Programmer struct {
name string
age int
language string
}
func main() {
p := Programmer{"stefno", 18, "go"}
fmt.Println(p)
lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Sizeof(int(0)) + unsafe.Sizeof(string(""))))
*lang = "Golang"
fmt.Println(p)
}
```
输出： 
{stefno 18 go}  
{stefno 18 Golang}  

参考： 深度解密Go语言之unsafe https://www.sohu.com/a/319106990_657921  
## 不打开新窗口访问url 
```
location.replace("http://www.csdn.net");
location="http://www.csdn.net";
window.location="http://www.csdn.net";
window.location.href="http://www.csdn.net";
```

## 隐藏input的三种方法和区别 
```
一、<input type="hidden" />
二、<input type="text" style="display:none" />
以上两种方法可以实现不留痕迹的隐藏。
三、<input type="text" style="visibility: hidden;" />
第三种方法可以实现占位隐藏（会留下空白而不显示）
```

## 页面等待一段时间后跳转 
参考https://my.oschina.net/tongjh/blog/220745?p={{page}} 
## 谷歌浏览器清缓存 
打开调试工具(mac:option + command + i, windows:ctrl + shift + i) , 按住地址栏刷新按钮，出现子菜单，选择[清空缓存并硬性重新加载]，解决 

## beego 
### 模板 
{{if eq .IsScored 0}} 
[controller.go:306]  template Execute err: template: book/bookdetail.html:131:9: executing "book/bookdetail.html" at <eq .IsScored 0>: error calling eq: invalid type for comparison
{{if compare .IsScored 0}}即可解决。 
参考http://www.144d.com/post-618.html 
### beego Can't create more than max_prepared_stmt_count statements (current value: 16382) 
使用versions-v1.11.1不会有上面问题
https://github.com/astaxie/beego/issues/3791 

## windows命令快捷键 
按Home键，快速回到命令开头，再按Crtl+e回到命令尾。 

## 创建唯一索引，删除重复数据 
```
保留id最小的重复数据一份
delete from `bookms_book` where identify in (
select * from (select identify from `bookms_book` group by identify having count(identify) > 1) temp
) and id not in (
select * from (select min(id) from `bookms_book` group by identify having count(identify) > 1) temp
);
select * from `bookms_book_record` where identify in (
select * from (select identify from `bookms_book_record` group by identify having count(identify)>1) tem
) and id not in (
select * from (select min(id) from `bookms_book_record` group by identify having count(identify)>1) tem
) order by id asc;

```

## redis+lua原子操作  
### zset基本操作  
按分数从大到小排名1~5  
zrevrangebyscore testzset 10 0 WITHSCORES LIMIT 0 5  
按从小到大排名删除第1~2个成员    
zremrangebyrank testzset 0 2  
获取成员数量  
zcard  testzset  
获取指定成员排名  
zrevrank testzset member  
添加成员  
zadd testzset 5 member  

### redis调用lua脚本  
#### EVAL命令格式  
```
EVAL script numkeys key [key ...] arg [arg ...]  
语义  
a. script即为lua脚本或lua脚本文件  
b. key一般指lua脚本操作的键，在lua脚本文件中，通过KEYS[i]获取,i从1开始  
c. arg指外部传递给lua脚本的参数，可以通过ARGV[i]获取,i从1开始  
```
#### EVAL示例  
在cmd/powershell未连接redis使用redis-cli调用lua脚本命令示例如下，其中“BookScoreRank”为key，“5 5 5 test”为参数  
```
redis-cli.exe -a 123 --eval zsetop.lua BookScoreRank , 5 5 5 test  
```
在cmd/powershell已连接redis的cli中调用lua脚本命令示例如下,其中"1"表示key的数目，其他参数同上  
```
127.0.0.1:6379> eval zsetop.lua 1 BookScoreRank 5 5 5 test  
```
#### script命令  
redis提供了以下几个script命令，用于对于脚本子系统进行控制：

script flush：清除所有的脚本缓存

script load：将脚本装入脚本缓存，不立即运行并返回其校验和

script exists：根据指定脚本校验和，检查脚本是否存在于缓存

script kill：杀死当前正在运行的脚本（防止脚本运行缓存，占用内存）

主要优势：  
减少网络开销：多个请求通过脚本一次发送，减少网络延迟

原子操作：将脚本作为一个整体执行，中间不会插入其他命令，无需使用事务

复用：客户端发送的脚本永久存在redis中，其他客户端可以复用脚本

可嵌入性：可嵌入JAVA，C#等多种编程语言，支持不同操作系统跨平台交互

通过script命令加载及执行lua脚本示例：
```
127.0.0.1:6379> script load "return 'Hello GrassInWind'"
"c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
127.0.0.1:6379> script exists "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
1) (integer) 1
127.0.0.1:6379> evalsha "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b" 0
"Hello GrassInWind"
127.0.0.1:6379> script flush
OK
127.0.0.1:6379> script exists "c66be1d9b54b3182f8d8e12f8b01a4e5c7c4af5b"
1) (integer) 0
```

#### redis log  
用于发送日志（log）的 redis.log 函数，以及相应的日志级别（level）：
redis.LOG_DEBUG
redis.LOG_VERBOSE
redis.LOG_NOTICE
redis.LOG_WARNING
默认redis server的Log级别为notice,在lua脚本中调用redis.log()打印的日志在server的日志中  
```
redis.log(redis.LOG_NOTICE, "list: ", lowestScore, lowestScoreMember)
```
用于计算 SHA1 校验和的 redis.sha1hex 函数。
用于返回错误信息的 redis.error_reply 函数和 redis.status_reply 函数。

### 开发遇到的错误  
#### Lua redis() command arguments must be strings or integers  
```
lua := redis.NewScript(1, "local len = redis.call('zcard', KEYS[1]) " +
		"if tonumber(len) >= tonumber(ARGV[1]) then " +
		"local num = len-tonumber(ARGV[1])+1 " +
		"local lowestScoreMember = redis.call('zrangebyscore',KEYS[1],0,ARGV[3],'limit',0,num) " +
		"local lowestScore = redis.call('zscore', KEYS[1],lowestScoreMember) " +
		"if tonumber(ARGV[2]) > tonumber(lowestScore) then " +
		"local index = len - tonumber(ARGV[1]) " +
		"redis.call('zremrangebyrank',KEYS[1], 0, index) end " +
		"end " +
		"redis.call('zadd', KEYS[1], ARGV[2], ARGV[4])")
```
其中"local lowestScoreMember = redis.call('zrangebyscore',KEYS[1],0,ARGV[3],'limit',0,num) " 返回的是list，故会出现如上错误  
修改如下,使用"for k,lowestScoreMember in pairs(list) do "来处理list    
```
lua := redis.NewScript(1, "local len = redis.call('zcard', KEYS[1]) " +
		"if tonumber(len) >= tonumber(ARGV[1]) then " +
		"local num = len-tonumber(ARGV[1])+1 " +
		"local list = redis.call('zrangebyscore',KEYS[1],0,ARGV[3],'limit',0,num) " +
		"for k,lowestScoreMember in pairs(list) do " +
		"local lowestScore = redis.call('zscore', KEYS[1],lowestScoreMember) " +
		"if tonumber(ARGV[2]) > tonumber(lowestScore) then " +
		"local index = len - tonumber(ARGV[1]) " +
		"redis.call('zremrangebyrank',KEYS[1], 0, index) end " +
		"end " +
		"end " +
		"redis.call('zadd', KEYS[1], ARGV[2], ARGV[4])")
```
