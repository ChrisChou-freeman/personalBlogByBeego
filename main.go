package main

import (
	"fmt"
	"time"

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
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("DATE_FORMAT(create_time, '%Y-%m-%d') as mytime, count(Id) as num").From("article").GroupBy("mytime").Limit(10)
	// fmt.Println(qb.String())
	qs := qb.String()
	blogDateList := []orm.Params{}
	o.Raw(qs).Values(&blogDateList)
	fmt.Println(blogDateList)
	mytime, _ := time.Parse("2006/01/02", "2019/03/05")
	valueList := []orm.Params{}
	o.QueryTable("Article").Filter("CreateTime__gt", mytime).Values(&valueList, "CreateTime")
	fmt.Println(mytime)
	fmt.Println(valueList)
	//article := []orm.Params{}
	//o.QueryTable("Article").GroupBy("ArticleType").Values(&article, "ArticleType")
	for _, item := range blogDateList {
		fmt.Println(item)
	}
	// fmt.Println(article)
	tools.Mycommands(orm.RunCommand)
	beego.Run()
}
