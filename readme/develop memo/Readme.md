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
