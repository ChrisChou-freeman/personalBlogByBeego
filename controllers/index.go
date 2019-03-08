package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"MyblogByGo/models"
	"MyblogByGo/tools"
)

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

	argmap := map[string]string{
		"searchdate":  "",
		"articletype": "",
	}
	for key := range argmap {
		arg := c.Input().Get(key)
		if arg != "" {
			argmap[key] = arg
		}
	}

	if argmap["searchdate"] != "" {
		mytime, _ := time.Parse("2006-01", argmap["searchdate"])
		v := strings.Split(argmap["searchdate"], "-")
		y := v[0]
		m := v[1]
		mint, _ := strconv.Atoi(m)
		nexmint := 0
		nexm := ""
		if mint == 12 {
			nexmint = 1
		} else {
			nexmint = mint + 1
		}
		if nexmint < 10 {
			nexm = fmt.Sprintf("0%v", nexmint)
		}
		nexMdateStr := fmt.Sprintf("%v-%v", y, nexm)
		nexMdateTime, _ := time.Parse("2006-01", nexMdateStr)
		articleListQs = articleListQs.Filter("CreateTime__gte", mytime).Filter("CreateTime__lte", nexMdateTime)
	} else if argmap["articletype"] != "" {
		articleTypeID := argmap["articletype"]
		typeID, _ := strconv.Atoi(articleTypeID)
		articleListQs = articleListQs.Filter("ArticleType", typeID)
	}

	urlarg := ""
	for key, item := range argmap {
		if item != "" {
			temp := fmt.Sprintf("&%v=%v", key, item)
			urlarg += temp
		}
	}

	page := c.Input().Get("page")
	pint, err := strconv.Atoi(page)
	if err != nil {
		pint = 1
	}
	pageLimit, start, stop := tools.LimitPage(pint, int(dataCount), urlarg, "/")
	articleListQs = articleListQs.Limit(stop, start)
	c.Data["pageLimit"] = pageLimit

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("DATE_FORMAT(create_time, '%Y-%m') as mytime, count(Id) as num").From("article").GroupBy("mytime").Limit(10)
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
