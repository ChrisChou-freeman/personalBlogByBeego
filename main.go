package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"

	_ "MyblogByGo/models"
	_ "MyblogByGo/routers"
	"MyblogByGo/tools"
)

func init() {
	if sessionprovider, err := beego.GetConfig("String", "sessionprovider", "--"); err == nil && sessionprovider == "mysql" {
		dbname, _ := beego.GetConfig("String", "dbname", "myblogbygo")
		sqluser, _ := beego.GetConfig("String", "sqluser", "root")
		sqlpass, _ := beego.GetConfig("String", "sqlpass", "123")
		sqlhost, _ := beego.GetConfig("String", "sqlhost", "127.0.0.1")
		sqlport, _ := beego.GetConfig("String", "sqlport", "3306")
		verification := "%s:%s@tcp(%s:%s)/%s?charset=utf8"
		verificationStr := fmt.Sprintf(verification, sqluser, sqlpass, sqlhost, sqlport, dbname)
		beego.BConfig.WebConfig.Session.SessionProviderConfig = verificationStr
	}
	beego.AddFuncMap("reContent", tools.ReContent)
}

func main() {
	tools.Mycommands(orm.RunCommand)
	beego.Run()
}
