# personalBlogByBeego

个人博客网站用beego框架实现
-------

主页
![博客样式](https://github.com/chrischou2018/personalBlogByBeego/raw/master/static/img/blog/blogimg.png)
后台
![后台样式](https://github.com/chrischou2018/personalBlogByBeego/raw/master/static/img/blog/adminimg.png)

依赖库
```
//执行(install tools)
go get github.com/astaxie/beego
go get github.com/beego/bee
go get github.com/astaxie/beego/session/mysql
go get github.com/go-sql-driver/mysql
```

初始化数据库表(init datatables)
```
go run main.go orm syncdb
```

初始化管理员账户(init admin user)
```
//执行命令(runcommand)
go run main.go initUser
// 初始化的用户名密码可以在tools/commands.go下修改
```

初始化类型标签(init article tags)
```
go run main.go  initTags
// 初始化博客类型分类可以在tools/commands.go下修改
```

初始化session表(init session table)
```
go run main.go  initSession
```

其余数据库连接相关配置在conf/app.conf中的# self config下面自行修改

启动项目(run project)
```
go run main.go
//默认启动在127.0.0.1:8000端口
```