package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"MyblogByGo/models"
	"MyblogByGo/tools"
)

// ReturnData 返回数据结构体
type ReturnData struct {
	Status  bool
	Data    string
	Message string
}

// IndexController 主页面的入口
type IndexController struct {
	beego.Controller
}

// IsLogin 确认登录
func (c *IndexController) IsLogin() bool {
	islogin := c.GetSession("isLogin")
	//fmt.Println("sessionValue", islogin)
	if islogin != nil && islogin == true {
		return true
	}
	return false
}

// Get 博客主页
func (c *IndexController) Get() {
	o := orm.NewOrm()
	typeList := []*models.ArticleType{}
	articleList := []*models.Article{}
	articleListQs := o.QueryTable("Article").RelatedSel("ArticleContent", "ArticleType")
	dataCount, _ := o.QueryTable("Article").Count()

	page := c.Input().Get("page")
	pint, err := strconv.Atoi(page)
	if err != nil {
		pint = 1
	}
	pageLimit, start, stop := tools.LimitPage(pint, int(dataCount), "", "/")
	articleListQs.Limit(stop, start-1)
	c.Data["pageLimit"] = pageLimit

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("DATE_FORMAT(create_time, '%Y-%m') as mytime, count(Id) as num").From("article").GroupBy("mytime").Limit(10)
	// fmt.Println(qb.String())
	qs := qb.String()
	blogDateList := []orm.Params{}
	if _, err := o.Raw(qs).Values(&blogDateList); err == nil {
		c.Data["blogDateList"] = &blogDateList
	} else {
		fmt.Println(err)
	}

	if _, err := o.QueryTable("ArticleType").All(&typeList); err == nil {
		c.Data["typelist"] = &typeList
	} else {
		fmt.Println(err)
	}
	if _, err := articleListQs.All(&articleList); err == nil {
		c.Data["ArticleList"] = &articleList
	} else {
		fmt.Println(err)
	}
	if v := c.IsLogin(); v {
		userId := c.GetSession("userId")
		user := models.User{
			Id: userId.(int),
		}
		o.Read(&user)
		c.Data["isLogin"] = true
		c.Data["userName"] = user.Name
	}
	c.Layout = "blog/layOut.html"
	c.TplName = "blog/index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["indexTitle"] = "blog/indexTitle.html"
}

// ArticleAddController 添加文章
type ArticleAddController struct {
	beego.Controller
}

// IsLogin 确认登录结构体
func (c *ArticleAddController) IsLogin() bool {
	islogin := c.GetSession("isLogin")
	fmt.Println("sessionValue", islogin)
	if islogin != nil && islogin == true {
		return true
	}
	return false
}

// Get 后去添加文章页面
func (c *ArticleAddController) Get() {
	if v := c.IsLogin(); !v {
		url := c.URLFor("AccountController.Get")
		c.Redirect(url, 302)
	}
	posturl := c.URLFor("ArticleAddController.Post")
	c.Data["posturl"] = posturl
	c.Data["xsrf_token"] = c.XSRFToken()
	o := orm.NewOrm()
	typeList := []*models.ArticleType{}
	if _, err := o.QueryTable("ArticleType").All(&typeList); err == nil {
		c.Data["typelist"] = &typeList
	} else {
		fmt.Println(err)
	}
	c.TplName = "blog/articleEditor.html"
}

type articleForm struct {
	ArticleName    string
	ArticleType    int
	ArticleContent string
}

// Post 提交文章
func (c *ArticleAddController) Post() {
	if v := c.IsLogin(); !v {
		url := c.URLFor("AccountController.Get")
		c.Redirect(url, 302)
	}
	rd := ReturnData{
		Status:  true,
		Data:    "",
		Message: "",
	}
	af := articleForm{}
	fmt.Println(string(c.Ctx.Input.RequestBody))
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &af); err != nil {
		rd.Status = false
		rd.Message = "文章提交出错"
	} else {
		o := orm.NewOrm()
		o.Using("default")
		article := new(models.Article)
		articleContent := new(models.ArticleContent)
		articleType := new(models.ArticleType)
		articleName := af.ArticleName
		articleName = strings.Replace(articleName, " ", "", -1)
		if articleName != "" {
			article.ArticleName = articleName
		} else {
			rd.Status = false
			rd.Message = "文章标题为空"
		}
		if af.ArticleType != 0 {
			articleType.Id = af.ArticleType
			if err := o.Read(articleType); err != nil {
				rd.Status = false
				rd.Message = "文章类型错误"
			}
			article.ArticleType = articleType
		} else {
			rd.Status = false
			rd.Message = "类型为空"
		}
		articleData := af.ArticleContent
		articleData = strings.Replace(articleData, " ", "", -1)
		articleData = strings.Replace(articleData, "/n", "", -1)
		if articleData != "" {
			articleContent.Content = articleData
			article.ArticleContent = articleContent
		} else {
			rd.Status = false
			rd.Message = "文章内容为空"
		}
		if rd.Status {
			fmt.Println(articleContent)
			_, acerr := o.Insert(articleContent)
			_, aerr := o.Insert(article)
			if acerr != nil || aerr != nil {
				rd.Status = false
				rd.Message = "文章提交出错"
				o.Rollback()
			}
		}
		rdJSON, _ := json.Marshal(rd)
		c.Data["json"] = string(rdJSON)
		c.ServeJSON()
	}

}

// ArticleContentController 访问文章内容
type ArticleContentController struct {
	beego.Controller
}

// Get 访问文章页面
func (c *ArticleContentController) Get() {
	o := orm.NewOrm()
	o.Using("default")
	articleId := c.Input().Get("articleid")
	aid, err := strconv.Atoi(articleId)
	if err != nil {
		aid = 1
	}
	fmt.Println("articleid", aid)
	article := new(models.Article)
	article.Id = aid
	qerr := o.Read(article)
	o.LoadRelated(article, "ArticleContent", "ArticleType")
	typeList := []*models.ArticleType{}
	if _, err := o.QueryTable("ArticleType").All(&typeList); err == nil {
		c.Data["typelist"] = &typeList
	} else {
		fmt.Println(err)
	}
	if qerr != nil {
		type Data struct {
			ArticleName string
		}
		c.Data["Article"] = Data{ArticleName: "文章读取错误"}
	} else {
		c.Data["Article"] = article
	}
	c.Data["typelist"] = &typeList
	c.Layout = "blog/layOut.html"
	c.TplName = "blog/articleContent.html"
}
