package utils

import "github.com/astaxie/beego"


func init()  {
	beego.AddFuncMap("prePage", ShowPrePage)
	beego.AddFuncMap("nextPage", ShowNextPage)
}

func ShowPrePage(pageIndex int) int {

	if pageIndex <= 1 {
		return 1
	}
	return pageIndex - 1
}

func ShowNextPage(pageIndex, pageCount int) int {

	if pageIndex >= pageCount {
		return pageCount
	}
	return pageIndex + 1
}
