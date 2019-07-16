package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"newsWeb/models"
	"github.com/astaxie/beego/orm"
	"encoding/base64"
)



type UserController struct {
	beego.Controller
}


func (this *UserController) ShowRegister()  {

	this.TplName = "register.html"
}


func (this *UserController) HandleReg()  {

	//后端开发四部曲
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//2.校验数据
	if userName == "" || pwd == "" {
		fmt.Println("用户名或密码不能为空")
		this.TplName = "register.html"
		return
	}
	fmt.Println("用户名：", userName, "密码:", pwd)

	//3.处理数据 / 保存数据
	o := orm.NewOrm()

	var user models.User
	user.Name = userName
	user.Pwd = pwd

	id, err := o.Insert(&user)
	if err != nil {
		fmt.Println("用户注册失败")
		this.TplName = "register.html"
		return
	}
	fmt.Println("注册成功，用户 id = ", id)


	//4.返回数据
	//渲染
	//this.TplName = "login.html"
//	重定向
	this.Redirect("/login", 302)
}


func (this *UserController) ShowLogin()  {
	//获取cookie
	userName := this.Ctx.GetCookie("userName")

	if userName != "" {
		//解码
		dec, err := base64.StdEncoding.DecodeString(userName)
		if err != nil {
			fmt.Println("ShowLogin, base64.StdEncoding.DecodeString: ", err)
		}

		this.Data["userName"] = string(dec)
		this.Data["checked"] = "checked"
	} else {

		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}


func (this *UserController) HandleLogin()  {

	//后端开发四部曲
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//2.校验数据
	if userName == "" || pwd == "" {
		errMsg := "用户名或密码不能为空"
		fmt.Println(errMsg)

		this.Data["errMsg"] = errMsg
		this.TplName = "login.html"
		return
	}
	fmt.Println("用户名：", userName, "密码:", pwd)

	//3.处理数据
	o := orm.NewOrm()

	var user models.User
	user.Name = userName
	//user.Pwd = pwd

	err := o.Read(&user, "Name")
	if err != nil {
		errMsg := "用户名不存在"
		fmt.Println(errMsg)

		this.Data["errMsg"] = errMsg
		this.TplName = "login.html"
		return
	}

	if user.Pwd != pwd {
		errMsg := "用户名密码不一致"
		fmt.Println(errMsg)

		this.Data["errMsg"] = errMsg
		this.TplName = "login.html"
		return
	}
	fmt.Println("登陆成功！")

	//设置cookie
	remember := this.GetString("remember")
	fmt.Println("remember = ", remember)

	if remember == "on" {
		//编码
		enc := base64.StdEncoding.EncodeToString([]byte(userName))

		this.Ctx.SetCookie("userName", enc, 60 * 60)
	} else {

		this.Ctx.SetCookie("userName", userName, -1)
	}
	
	//登陆成功后 设置session
	this.SetSession("userName", userName)

	//4.返回数据
	//渲染
	//this.TplName = "index.html"
//	重定向
	this.Redirect("/article/index", 302)
}
//退出
func (this *UserController) HandleLogout()  {

//	清空session
	this.DelSession("userName")
//	跳转登陆页面
	this.Redirect("/login", 302)
}