
### for linux mysql5.x.x

### install
```
yum install mysql-devel.x86_64
make
```
* 留意mysql库位置不同时需修改 cgo CFLAGS: -I/usr/include/mysql
* 编译完成生成mysqlHttp.so并拷贝到select @@plugin_dir; 结果路径目录下

### usage
```
select @@plugin_dir;

drop function httpPost;
create function jsonObject returns string soname 'mysqlHttp.so';
create function httpPost returns string soname 'mysqlHttp.so';

select jsonObject("id", Host) into @temp from user limit 1;
select @temp;
select httpPost('https://www.caoxianjie.cn/todo.php', @temp);
```


### 触发器
```
/* INSERT插入操作的触发器 */  
DELIMITER |  
DROP TRIGGER IF EXISTS mytable_insert;  
CREATE TRIGGER mytable_insert  
AFTER INSERT ON mytable  
FOR EACH ROW BEGIN  
    SELECT json_object("action", "insert", "id", id) into @temp FROM mytable WHERE id = NEW.id LIMIT 1);  
    SELECT httpPost('https://www.caoxianjie.cn/todo.php', @temp);  
END |  
DELIMITER ;  
  
/* UPDATE更新操作的触发器 */  
DELIMITER |  
DROP TRIGGER IF EXISTS mytable_update;  
CREATE TRIGGER mytable_update  
AFTER UPDATE ON mytable  
FOR EACH ROW BEGIN  
    SELECT json_object("action", "update", "id", id) into @temp FROM mytable WHERE id = OLD.id LIMIT 1);  
    SELECT httpPost('https://www.caoxianjie.cn/todo.php', @temp);  
END |  
DELIMITER ;  
  
/* DELETE删除操作的触发器 */  
DELIMITER |  
DROP TRIGGER IF EXISTS mytable_delete;  
CREATE TRIGGER mytable_delete  
AFTER DELETE ON mytable  
FOR EACH ROW BEGIN  
    SELECT json_object("action", "delete", "id", id) into @temp FROM mytable WHERE id = OLD.id LIMIT 1);  
    SELECT httpPost('https://www.caoxianjie.cn/todo.php', @temp); 
END |  
DELIMITER ;  
```

### thanks
* https://dev.mysql.com/doc/refman/8.0/en/adding-udf.html
* https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices

### support
* qq group 233415606
