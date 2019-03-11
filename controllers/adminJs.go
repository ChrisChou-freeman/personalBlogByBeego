package controllers

import (
	"MyblogByGo/models"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// AdminJsControllers 后台js事件处理
type AdminJsControllers struct {
	beego.Controller
}

// TableConfigArg 表格配置参数
type TableConfigArg struct {
	Q       string
	Title   string
	Display bool
	Text    map[string]string
	Attrs   map[string]string
}

// GetTableConfig Admin 数据表格配置
func GetTableConfig(tablename string) []TableConfigArg {
	ArticleTableConfig := []TableConfigArg{
		{
			Q:       "",
			Title:   "选项",
			Display: true,
			Text:    map[string]string{"content": "<input type='checkbox' />", "kwargs": "{}"},
			Attrs:   map[string]string{},
		},
		{
			Q:       "Id",
			Title:   "ID",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@Id"}`},
			Attrs:   map[string]string{},
		},
		{
			Q:       "ArticleName",
			Title:   "文章标题",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@ArticleName"}`},
			Attrs:   map[string]string{"name": "ArticleName", "origin": "@ArticleName", "edit-enable": "true", "edit-type": "input"},
		},
		{
			Q:       "ArticleType__Id",
			Title:   "文章类型",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@@ArticleType"}`},
			Attrs:   map[string]string{"name": "ArticleType", "origin": "@ArticleType__Id", "edit-enable": "true", "edit-type": "select", "global-name": "ArticleType"},
		},
		{
			Q:       "",
			Title:   "操作",
			Display: true,
			Text:    map[string]string{"content": "<a href='/assetdetail-{m}.html'>{n}</a>", "kwargs": `{"n": "文章","m": "@id"}`},
			Attrs:   map[string]string{},
		},
	}
	UserTableConfig := []TableConfigArg{
		{
			Q:       "",
			Title:   "选项",
			Display: true,
			Text:    map[string]string{"content": "<input type='checkbox' />", "kwargs": "{}"},
			Attrs:   map[string]string{},
		},
		{
			Q:       "Id",
			Title:   "ID",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@Id"}`},
			Attrs:   map[string]string{},
		},
		{
			Q:       "Name",
			Title:   "用户名",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@Name"}`},
			Attrs:   map[string]string{"name": "Name", "origin": "@Name", "edit-enable": "true", "edit-type": "input"},
		},
		{
			Q:       "PassWord",
			Title:   "密码",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@PassWord"}`},
			Attrs:   map[string]string{"name": "PassWord", "origin": "@PassWord", "edit-enable": "true", "edit-type": "input"},
		},
		{
			Q:       "About",
			Title:   "关于",
			Display: true,
			Text:    map[string]string{"content": "{n}", "kwargs": `{"n": "@About"}`},
			Attrs:   map[string]string{"name": "About", "origin": "@About", "edit-enable": "true", "edit-type": "input"},
		},
	}
	var tableConfig []TableConfigArg
	switch {
	case tablename == "article":
		tableConfig = ArticleTableConfig
	case tablename == "user":
		tableConfig = UserTableConfig
	}
	return tableConfig
}

// Get js Get 方法
func (c *AdminJsControllers) Get() {
	type ReturnData struct {
		TableConfig []TableConfigArg
		DataList    []orm.Params
		GlobalDict  map[string][]orm.ParamsList
		Pager       string
	}
	rd := ReturnData{}
	o := orm.NewOrm()
	o.Using("default")
	dataList := []orm.Params{}
	tableName := ""
	dataType := c.Input().Get("datatype")
	if dataType == "article" {
		tableName = "Article"
		articleTypeList := []orm.ParamsList{}
		o.QueryTable("ArticleType").ValuesList(&articleTypeList, "Id", "TypeName")
		rd.GlobalDict = map[string][]orm.ParamsList{"ArticleType": articleTypeList}
	} else if dataType == "user" {
		tableName = "User"
	} else {
		rd.GlobalDict = map[string][]orm.ParamsList{}
	}
	TableConfig := GetTableConfig(dataType)
	rd.TableConfig = TableConfig
	if tableName != "" {
		qList := []string{}
		for _, item := range TableConfig {
			if item.Q != "" {
				qList = append(qList, item.Q)
			}
		}
		_, err := o.QueryTable(tableName).Values(&dataList, qList...)
		if err == nil {
			rd.DataList = dataList
		} else {
			rd.DataList = []orm.Params{}
		}
	} else {
		rd.DataList = []orm.Params{}
	}
	rdb, _ := json.Marshal(rd)
	c.Data["json"] = string(rdb)
	c.ServeJSON()
}

type adminArticleForm struct {
	Id          int
	ArticleName string
	ArticleType int
}

type AdminUserForm struct {
	Id       int
	Name     string
	PassWord string
	About    string
}

func (af adminArticleForm) AdminArticleEditorFuc(rd *ReturnData) {
	o := orm.NewOrm()
	o.Using("default")
	article := new(models.Article)
	articleType := new(models.ArticleType)
	articleName := af.ArticleName
	articleName = strings.Replace(articleName, " ", "", -1)
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
	} else {
		rd.Status = false
		rd.Message = "类型为空"
	}
	if rd.Status {
		article.Id = af.Id
		rerr := o.Read(article, "Id")
		if rerr != nil {
			fmt.Println(article.Id)
			rd.Status = false
			rd.Message = "找不到文章"
		} else {
			article.ArticleName = af.ArticleName
			article.ArticleType = articleType
			_, aerr := o.Update(article)
			if aerr != nil {
				rd.Status = false
				rd.Message = "文章提交出错"
				o.Rollback()
			}
		}
	}
}

func (uf AdminUserForm) AdminUserEditorForm(rd *ReturnData) {

}

func (c *AdminJsControllers) Post() {
	datatype := c.Input().Get("datatype")
	o := orm.NewOrm()
	o.Using("default")
	switch {
	case datatype == "article":

	}
}
