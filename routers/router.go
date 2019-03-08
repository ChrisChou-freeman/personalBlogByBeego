package routers

import (
	"github.com/astaxie/beego"

	"MyblogByGo/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/signin", &controllers.AccountController{})
	beego.Router("/signout", &controllers.SignOutController{})
	beego.Router("/article_add", &controllers.ArticleAddController{})
	beego.Router("/article_editor", &controllers.ArticleEditorController{})
	beego.Router("/articleconten", &controllers.ArticleContentController{})
}
