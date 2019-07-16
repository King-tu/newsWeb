package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
    "fmt"
)

func init() {

    beego.InsertFilter("/article/*", beego.BeforeExec, filtersFunc)

    //beego.Router("/", &controllers.MainController{})
    //beego.Router("/", &controllers.ArticleController{}, "get:ShowIndex")

    //注册
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleReg")
    //登陆
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
    //退出
    beego.Router("/logout", &controllers.UserController{}, "get:HandleLogout")

    //展示首页
    beego.Router("/article/index", &controllers.ArticleController{}, "get:ShowIndex")

    beego.Router("/article/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")

    beego.Router("/article/content", &controllers.ArticleController{}, "get:ShowContent")
    beego.Router("/article/editArticle", &controllers.ArticleController{}, "get:ShowEditArticle;post:HandleEditArticle")

    beego.Router("/article/delete", &controllers.ArticleController{}, "get:HandleDelete")

    beego.Router("/article/addType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")

    beego.Router("/article/delType", &controllers.ArticleController{}, "get:HandleDelType")

}

func filtersFunc(ctx *context.Context) {

    //获取用户登陆状态( session )
    //userName := this.GetSession("userName")   //controler中获取session的方法
	//context 上下文
    userName := ctx.Input.Session("userName")


    if userName == nil {

        ctx.Redirect(302, "/login")

        fmt.Println("路由过滤器已过滤。。。")
        return
    }

}