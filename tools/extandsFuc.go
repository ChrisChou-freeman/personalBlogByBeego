package tools

import (
	"fmt"
	"math"
	"strings"

	"github.com/astaxie/beego"
)

// LimitPage 分页组建
func LimitPage(CurreentPage, AllCount int, FilterArgs, url string) (string, int, int) {
	pageCount, _ := beego.GetConfig("Int", "pagecount", 10)
	showPage := 10
	if AllCount < 1 {
		AllCount = 1
	}
	allPage := AllCount / pageCount.(int)
	mod := math.Mod(float64(AllCount), float64(pageCount.(int)))
	if mod > 0 {
		allPage++
	}
	htmlList := []string{}
	pageHalf := (showPage - 1) / 2
	start := 0
	stop := 0
	var previous string
	var next string
	if allPage < showPage {
		start = 1
		stop = allPage
	} else {
		if CurreentPage < pageHalf+1 {
			start = 1
			stop = showPage
		} else {
			if CurreentPage >= allPage-pageHalf {
				start = allPage - showPage
				stop = allPage
			} else {
				start = CurreentPage - pageHalf
				stop = CurreentPage + pageHalf
			}
		}
	}
	if CurreentPage <= 1 {
		previous = "<li class='page-item disabled'><a href='#' class='page-link'>上一页</a></li>"
	} else {
		as := "<li class='page-item'><a href='%v?page=%v%v' class='page-link'  style='cursor:pointer;text-decoration:none;'>上一页<span aria-hidden='true'>&laquo;</span></a></li>"
		previous = fmt.Sprintf(as, url, CurreentPage-1, FilterArgs)
	}
	htmlList = append(htmlList, previous)
	for i := start; i <= stop; i++ {
		temp := ""
		if CurreentPage == i {
			temp = "<li class='page-item active'><a href='%v?page=%v%v' class='page-link' >%v</a></li>"
			temp = fmt.Sprintf(temp, url, i, FilterArgs, i)
		} else {
			temp = "<li class='page-item'><a href='%v?page=%v%v' class='page-link' >%v</a></li>"
			temp = fmt.Sprintf(temp, url, i, FilterArgs, i)
		}
		htmlList = append(htmlList, temp)
	}
	if CurreentPage >= allPage {
		next = "<li class='page-item disabled'><a href='#' class='page-link'>下一页</a></li>"
	} else {
		as := "<li class='page-item'><a href='%v?page=%v%v' class='page-link' >下一页</a></li>"
		next = fmt.Sprintf(as, url, CurreentPage+1, FilterArgs)
	}
	htmlList = append(htmlList, next)
	data := strings.Join(htmlList, "")

	dataStart := 0
	dataStop := 0
	dataStart = (CurreentPage - 1) * pageCount.(int)
	dataStop = pageCount.(int)
	return data, dataStart, dataStop
}
