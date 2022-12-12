## 介绍
worm是一款方便易用的Go语言ORM库，worm具有使用简单，运行性能高，功能强大的特点。具体特征如下：
* 通过Struct的Tag与数据库字段进行映射，让您免于经常拼写SQL的麻烦。
* 支持Struct映射、原生SQL以及SQL builder三种模式来操作数据库，并且Struct映射、原生SQL以及SQL builder可混合使用。
* Struct映射、SQL builder支持链式API，可使用Where, And, Or, ID, In, Limit, GroupBy, OrderBy, Having等函数构造查询条件。
* 可通过Join、LeftJoin、RightJoin来进行数据库表之间的关联查询。
* 支持事务支持，可在会话中开启事务，在事务中可以混用Struct映射、原生SQL以及SQL builder来操作数据库。
* 支持预编译模式访问数据库，会话开启预编译模式后，任何SQL都会使用缓存的Statement，可以提升数据库访问效率。
* 支持使用数据库的读写分离模式来访问一组数据库。
* 支持在Insert，Update，Delete，Get, Find操作中使用钩子方法。
* 支持SQL日志的输出，并且可通过日志钩子自定义SQL日志的输出，日志钩子支持Context。
* 可根据数据库自动生成库表对应的Model结构体。

目前worm支持的数据库有：mysql、postgres、sqlite、sqlserver。

## 安装
go get github.com/haming123/wego/worm

## 文档
请点击：[详细文档](http://39.108.252.54:8080/docs/worm/worm)

## 快速开始
### 创建实体类
```
//建表语句
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  `passwd` varchar(32) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);
```
数据库表user对应的实体类的定义如下：
```Go
type User struct {
	Id          int64   	`db:"id;autoincr"`
	Name        string  	`db:"name"`
	Age         int64   	`db:"age"`
	Passwd      string  	`db:"passwd"`
	Created     time.Time	`db:"created;n_update"`
}
func (ent *User) TableName() string {
	return "user"
}
```
worm使用名称为"db"的Tag映射数据库字段，"db"后面是字段的名称，autoincr用于说明该字段是自增ID，n_update用于说明该字段不可用于update语句中。

### 创建DbEngine
本文中的例子使用的都是mysql数据库。若要创建一个mysql数据库的DbEngine，您可调用worm.NewMysql()函数或者调用worm.InitMysql()函数。
```Go
package main
import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/haming123/wego/worm"
)
func main() {
	var err error
	cnnstr := "user:pwd@tcp(127.0.0.1:3306)/db?charset=utf8&parseTime=True"
	dbcnn, err := sql.Open("mysql", cnnstr)
	if err != nil {
		log.Println(err)
		return
	}
	err = dbcnn.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	err = worm.InitMysql(dbcnn)
	if err != nil {
		log.Println(err)
		return
	}
}
```

### 插入记录
```Go
user := User{Name:"name1", Age: 21, Created: time.Now()}
id, err := worm.Model(&user).Insert()
//insert into user set name=?,age=?,passwd=?,created=?
```

### 更新数据
```Go
//通过ID更新数据
user := User{Name:"name2", Age: 22}
affected, err := worm.Model(&user).ID(1).Update()
//update user set name=?,age=?,passwd=? where id=?
```
为了防止误操作，更新数据时必须指定查询条件，若没有指定查询条件worm不会执行该操作。

### 删除数据
```Go
//通过ID删除数据
affected, err := worm.Model(&User{}).ID(1).Delete()
//delete from user where id=?
```
为了防止误操作，删除数据时必须指定查询条件，若没有指定查询条件worm不会执行该操作。

### 查询单条记录
```Go
//通过ID查询数据
user := User{}
_, err := worm.Model(&user).ID(1).Get()
//select id,name,age,passwd,created from user where id=? limit 1
```

### 查询多条记录
```Go
//查询全部数据
users := []User{}
err := worm.Model(&User{}).Find(&users)
//select id,name,age,passwd,created from user

//使用limit查询数据
users = []User{}
err := worm.Model(&User{}).Limit(5).Offset(2).Find(&users)
//select id,name,age,passwd,created from user limit 2, 5

//使用order查询数据
users = []User{}
err = worm.Model(&User{}).OrderBy("name asc").Find(&users)
//select id,name,age,passwd,created from users order by name asc
```

### 查询指定的字段
worm允许通过Select方法选择特定的字段, 或者使用Omit方法排除特定字段
```Go
//查询指定的字段
user := User{}
_, err := worm.Model(&user).Select("id", "name", "age").ID(1).Get()
// select id,name,age from users where id=1 limit 1

user = User{}
_, err = worm.Model(&user).Omit("passwd").ID(1).Get()
//select id,name,age,created from users where id=1 limit 1
```

### 条件查询
worm支持链式API，可使用Where, And, Or, ID, In等函数构造查询条件。
```Go
users := []User{}
err := worm.Model(&User{}).Where("age>?", 0).Find(&users)
//select id,name,age,passwd,created from user where age>?

//and
users = []User{}
err := worm.Model(&User{}).Where("age>?", 0).And("id<?", 10).Find(&users)
//select id,name,age,passwd,created from user where age>? and id<?

//like
users = []User{}
err := worm.Model(&User{}).Where("name like ?", "%name%").Find(&users)
//select id,name,age,passwd,created from user where name like ?

//in
users = []User{}
err := worm.Model(&User{}).Where("age>?", 0).AndIn("id", 5,6).Find(&users)
//select id,name,age,passwd,created from user where age>? and id in (?,?)

```
worm占位符统一使用?，worm会根据数据库类型,自动替换占位符，例如postgresql数据库把?替换成$1,$2...

### 获取记录条数
```Go
num, err := Model(&User{}).Where("id>?", 0).Count()
//select count(1) from user where id>?

num, err := Model(&User{}).Where("id>?", 0).Count("name")
//select count(name) from user where id>?
```

### 检测记录是否存在
```Go
has, err := worm.Model(&User{}).Where("id=?", 1).Exist()
//select id,name,age,passwd,created from user where id=? limit 1
```

### Update或Insert数据
```Go
var user = User{Name:"name2", Age: 22}
affected,insertId,err := worm.Model(&user).UpdateOrInsert(1)
//根据id判断是Update还是Insert， 若id>0，则调用Update，否则调用Insert
//insert into user set name=?,age=?,passwd=?,created=?
//update user set name=?,age=?,passwd=? where id=?
```

### 批量新增
```Go
users := []User{User{DB_name:"batch1", Age: 33}, User{DB_name:"batch2", Age: 33} }
res, err := Model(&User{}).BatchInsert(&users)
//insert into user set name=?,age=?,created=?
//TX: time=141.074ms affected=2; prepare
```

### 使用原生SQL
```Go
//插入记录
res, err :=  worm.SQL("insert into user set name=?, age=?", "name1", 1).Exec()
//更新记录
res, err =  worm.SQL("update user set name=? where id=?", "name2", id).Exec()
//删除记录
res, err = worm.SQL("delete from user where id=?", id).Exec()

//查询单条记录
var name string; var age int64
has, err = worm.SQL("select name,age from user where id=?", 6).Get(&name, &age)

//查询单条记录到model对象
var ent User
has, err := worm.SQL("select * from user where id=?", 6).GetModel(&ent)

//查询多条记录到model数组
var users []User
err := worm.SQL( "select * from user where id>?", 5).FindModel(&users)
```

### 使用Sql Builder
```Go
//插入记录
id, err := worm.Table("user").Value("name", "name1").Value("age", 21).Insert()
//更新记录
affected, err := worm.Table("user").Value("name", "name2").Where("id=?", id).Update()
//删除记录
affected, err = worm.Table("user").Where("id=?", id).Delete()

//查询单条记录
var name string; var age int64
has, err = worm.Table("user").Select("name", "age").Where("id=?", 6).Get(&name, &age)

//查询单条记录到model对象
var ent User
has, err := worm.Table("user").Select("*").Where("id=?", 6).GetModel(&ent)

//查询多条记录
var users []User
err :=  worm.Table("user").Select("*").Where("id>?", 0).FindModel(&users)
```

### 关联查询
定义一个book表以及相应的实体类：
```
CREATE TABLE `book` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(16) NOT NULL DEFAULT '',
  `price` decimal(11,2) NOT NULL DEFAULT 0.0,
  PRIMARY KEY (`id`)
);
```
```Go
type Book struct {
	Id          int64   	`db:"id;autoincr"`
	Name        string  	`db:"name"`
	Author  	int64       `db:"author"`
	Price       float32     `db:"price"`
}
func (ent *Book) TableName() string {
	return "book"
}
```
在book表中，通过author字段与user表的id字段相关联。若要查询一个用户购买的书，在worm中可以通过Join来查询：
```Go
type UserBook struct {
    User
    Book
}
var datas []UserBook
md := worm.Model(&User{}).Select("id","name","age").TableAlias("u")
md.Join(&Book{}, "b", "b.author=u.id", "name")
err := md.Where("u.id>?", 0).Find(&datas)
//select u.id as u_id,u.name as u_name,u.age as u_age,b.name as b_name from user u join book b on b.author=u.id where u.id>0
```
除了Join，您还可以使用LeftJoin以及RightJoin进行左连接、右连接查询。

### 事务处理
当使用事务处理时，需要创建 Session对象，并开启数据库事务。在进行事务处理时，在事务中可以混用Model模式、原生SQL模式以及SQL builder模式来操作数据库：
```Go
tx := worm.NewSession()
tx.TxBegin()
user := User{Name:"name1", Age: 21, Created: time.Now()}
id, err := tx.Model(&user).Insert()
if err != nil{
    tx.TxRollback()
    return
}
_, err = tx.Table("user").Value("age", 20).Value("name", "zhangsan").Where("id=?", id).Update()
if err != nil{
    tx.TxRollback()
    return
}
_, err = tx.SQL("delete from user where id=?", id).Exec()
if err != nil{
    tx.TxRollback()
    return
}
tx.TxCommit()
```

### 使用SQL语句预编译
worm支持SQL语句的预编译，使用SQL语句预编译可以提升数据库访问的效率。在worm中可以通过三种方式开启SQL语句预编译：全局开启、会话中开启、语句中开启。
```Go
//全局开启，所有操作都会创建并缓存预编译
worm.UsePrepare(true)

//会话中开启, 会话中的sql语句会创建并缓存预编译
dbs := worm.NewSession()
dbs.UsePrepare(true)

//语句中使用UsePrepare
user := User{}
_, err := Model(&user).UsePrepare(true).Where("id=?", 1).Get()
```
