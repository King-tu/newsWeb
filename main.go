package main

import (
	_ "newsWeb/utils"
	_ "newsWeb/models"
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

