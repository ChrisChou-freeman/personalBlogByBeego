package tools

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/astaxie/beego/orm"

	"MyblogByGo/models"
)

var commandMaps = map[string]func(){
	"initTags": initTags,
	"initUser": initUser,
}

// initTags 初始化标签函数
func initTags() {
	tags := []models.ArticleType{
		{TypeName: "Python"},
		{TypeName: "Go"},
		{TypeName: "Mysql"},
		{TypeName: "Django"},
		{TypeName: "Beego"},
		{TypeName: "杂谈"},
	}
	o := orm.NewOrm()
	o.Using("default")
	for _, item := range tags {
		if created, _, err := o.ReadOrCreate(&item, "TypeName"); err == nil {
			if created {
				fmt.Println("create new ArticleType:", item.TypeName)
			} else {
				fmt.Println("该元素以及存在:", item.TypeName)
			}
		}
	}
	os.Exit(0)
}

// initUser 初始化管理员
func initUser() {
	user := new(models.User)
	user.Name = "cris"
	userpass := "admin123"
	myHs := sha512.New()
	myHs.Write([]byte(userpass))
	myHasPas := myHs.Sum(nil)
	encodedPss := base64.StdEncoding.EncodeToString([]byte(myHasPas))
	user.PassWord = encodedPss
	o := orm.NewOrm()
	o.Using("default")
	if id, err := o.Insert(user); err == nil {
		fmt.Println("用户已经创建：", user.Name, id)
	} else {
		fmt.Println("error:", err)
	}
	os.Exit(0)
}

// Mycommands 自定义命令
func Mycommands(com func()) {
	if len(os.Args) < 2 {
		return
	} else if elm, ok := commandMaps[os.Args[1]]; ok {
		elm()
	} else {
		com()
	}
}
