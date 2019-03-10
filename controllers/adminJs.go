package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
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
	if tablename == "article" {
		return ArticleTableConfig
	}
	return ArticleTableConfig
}

// Get js Get 方法
func (c *AdminJsControllers) Get() {
	TableConfig := GetTableConfig("article")
	//b, _ := json.Marshal(TableConfig)
	dataList := []orm.Params{}
	qList := []string{}
	for _, item := range TableConfig {
		if item.Q != "" {
			qList = append(qList, item.Q)
		}
	}
	o := orm.NewOrm()
	o.Using("default")
	o.QueryTable("Article").Values(&dataList, qList...)
	articleTypeList := []orm.ParamsList{}
	o.QueryTable("ArticleType").ValuesList(&articleTypeList, "Id", "TypeName")
	type ReturnData struct {
		TableConfig []TableConfigArg
		DataList    []orm.Params
		GlobalDict  map[string][]orm.ParamsList
		Pager       string
	}
	rd := ReturnData{
		TableConfig: TableConfig,
		DataList:    dataList,
		GlobalDict:  map[string][]orm.ParamsList{"ArticleType": articleTypeList},
		Pager:       "",
	}
	rdb, _ := json.Marshal(rd)
	c.Data["json"] = string(rdb)
	c.ServeJSON()
}
