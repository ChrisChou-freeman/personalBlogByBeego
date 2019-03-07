package routers

import (
	"github.com/astaxie/beego"

	"MyblogByGo/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/signin", &controllers.AccountController{})
	beego.Router("/signout", &controllers.SignOutController{})
	beego.Router("/article_editor", &controllers.ArticleAddController{})
	beego.Router("/articleconten", &controllers.ArticleContentController{})
}
