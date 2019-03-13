package controllers

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"MyblogByGo/models"
)

// AccountController 登陆入口
type AccountController struct {
	beego.Controller
}

// Get Signin blog
func (c *AccountController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "blog/signin.html"
}

type userFrom struct {
	Username string `form:"username"`
	Pass     string `form:"pass"`
}

// Post post signin
func (c *AccountController) Post() {
	uf := userFrom{}
	if err := c.ParseForm(&uf); err != nil {
		fmt.Println("数据提交出错：", err)
	} else {
		user := models.User{
			Name: uf.Username,
		}
		o := orm.NewOrm()
		o.Using("default")
		err := o.Read(&user, "Name")
		if err == orm.ErrNoRows {
			fmt.Println("查询不到用户:", user.Name)
			c.Data["errmsg"] = "找不到用户"
			c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
			c.TplName = "blog/signin.html"
		} else if err == orm.ErrMissPK {
			fmt.Println("找不到主键")
			c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
			c.TplName = "blog/signin.html"
		} else {
			ufpass := uf.Pass
			myHs := sha512.New()
			myHs.Write([]byte(ufpass))
			myHasPas := myHs.Sum(nil)
			encodedPss := base64.StdEncoding.EncodeToString([]byte(myHasPas))
			if user.PassWord == encodedPss {
				c.SetSession("isLogin", bool(true))
				c.SetSession("userId", int(user.Id))
				// c.SetSession("userid", user.Id)
				c.Redirect("/", 302)
			} else {
				c.Data["errmsg"] = "密码错误"
				c.TplName = "blog/signin.html"
			}
		}

	}
}

// SignOutController 退出登录
type SignOutController struct {
	beego.Controller
}

// Get 退出登录
func (c *SignOutController) Get() {
	c.DestroySession()
	url := beego.URLFor("AccountController.Get")
	c.Redirect(url, 302)
}
