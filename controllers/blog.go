package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"MyblogByGo/models"
)

// ReturnData 返回数据结构体
type ReturnData struct {
	Status  bool
	Data    string
	Message string
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

// Get 添加文章页面
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

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("DATE_FORMAT(create_time, '%Y-%m') as mytime, count(Id) as num").From("article").GroupBy("mytime").Limit(10)
	qs := qb.String()
	blogDateList := []orm.Params{}
	if _, err := o.Raw(qs).Values(&blogDateList); err == nil {
		c.Data["blogDateList"] = &blogDateList
	} else {
		fmt.Println(err)
	}

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
		c.Data["Article"] = &Data{ArticleName: "文章读取错误"}
	} else {
		c.Data["Article"] = article
	}
	c.Data["typelist"] = &typeList
	c.Layout = "blog/layOut.html"
	c.TplName = "blog/articleContent.html"
}

// ArticleEditorController 编辑文章
type ArticleEditorController struct {
	beego.Controller
}

// IsLogin 确认登录结构体
func (c *ArticleEditorController) IsLogin() bool {
	islogin := c.GetSession("isLogin")
	fmt.Println("sessionValue", islogin)
	if islogin != nil && islogin == true {
		return true
	}
	return false
}

// Get 获取文章编辑页面
func (c *ArticleEditorController) Get() {
	c.Data["xsrf_token"] = c.XSRFToken()
	articleId := c.Input().Get("articleid")
	aId, _ := strconv.Atoi(articleId)
	article := new(models.Article)
	article.Id = aId
	o := orm.NewOrm()
	o.Using("default")
	qerr := o.Read(article, "Id")
	if qerr != nil {
		type Data struct {
			ArticleName string
		}
		c.Data["Article"] = &Data{ArticleName: "文章读取错误"}
	} else {
		o.LoadRelated(article, "ArticleContent", "ArticleType")
	}

	typeList := []*models.ArticleType{}
	if _, err := o.QueryTable("ArticleType").All(&typeList); err == nil {
		c.Data["typelist"] = &typeList
	} else {
		fmt.Println(err)
	}
	c.Data["posturl"] = c.URLFor("ArticleEditorController.Post")
	c.Data["Article"] = article
	c.TplName = "blog/articleEditor.html"
}

// articleForm 文章更新form
type articleEditorForm struct {
	Id             int
	ArticleName    string
	ArticleType    int
	ArticleContent string
}

// Post 提交编辑文章内容
func (c *ArticleEditorController) Post() {
	if v := c.IsLogin(); !v {
		url := c.URLFor("AccountController.Get")
		c.Redirect(url, 302)
	}
	rd := ReturnData{
		Status:  true,
		Data:    "",
		Message: "",
	}
	af := articleEditorForm{}
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
		article.Id = af.Id
		if articleName == "" {
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
		if articleData == "" {
			rd.Status = false
			rd.Message = "文章内容为空"
		}
		if rd.Status {
			rerr := o.Read(article, "Id")
			if rerr != nil {
				fmt.Println(article.Id)
				rd.Status = false
				rd.Message = "找不到文章"
			} else {
				article.ArticleName = af.ArticleName
				article.ArticleType = articleType
				_, aerr := o.Update(article)
				o.QueryTable("ArticleContent").Filter("Article__id", article.Id).One(articleContent)
				articleContent.Content = af.ArticleContent
				_, cerr := o.Update(articleContent)
				if aerr != nil || cerr != nil {
					rd.Status = false
					rd.Message = "文章提交出错"
					o.Rollback()
				}
			}
		}
		rdJSON, _ := json.Marshal(rd)
		c.Data["json"] = string(rdJSON)
		c.ServeJSON()
	}
}
