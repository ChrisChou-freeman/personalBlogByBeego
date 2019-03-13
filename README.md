# personalBlogByBeego

个人博客网站用beego框架实现
-------

依赖库
```
//执行
go get github.com/astaxie/beego
go get github.com/astaxie/beego/session/mysql
go get github.com/go-sql-driver/mysql
```

初始化管理员账户
```
//执行命令
go run main.go initUser
// 初始化的用户名密码可以在tools/commands.go下修改
```

初始化分类标签
```
go run main.go  initTags
```

初始化session表
```
go run main.go  initSession
```

其余数据库连接相关配置在conf/app.conf中的# self config下面自行修改