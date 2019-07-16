package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"fmt"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
	"strings"
)

type ArticleController struct {
	beego.Controller
}

//展示 首页
func (this *ArticleController) ShowIndex()  {

	//获取用户登陆状态( session )
/*	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login",302)
		return
	}*/

	//fmt.Println("ShowIndex userName = ",userName)

	o := orm.NewOrm()
	//定义切片，保存从数据库查询到的记录
	var articles []models.Article


	typeName := this.GetString("select")
	//fmt.Println("typeName = ", typeName)
	//查询表的所有记录
	//qs := o.QueryTable("article")
	qs := o.QueryTable("Article")

	//每页展示的文章数
	pageSize := 2
	//总的文章数
	// RelatedSel("ArticleType") 关联查询外健
	var count int64
	var err error
	if typeName == "" {

		count, err = qs.RelatedSel("ArticleType").Count()

		if err != nil {
			fmt.Println("ShowIndex.Count1 err:", err)
		}
	} else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		if err != nil {
			fmt.Println("ShowIndex.Count2 err:", err)
		}
	}
	//总页数
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	//页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		//fmt.Println("ShowIndex.GetInt err:", err)	//0 strconv.Atoi: parsing "": invalid syntax
		pageIndex = 1
	}

	if typeName == "" {
		//orm默认惰性查询（不会主动关联查询外键信息）
		//_, err = qs.Limit(pageSize, (pageIndex - 1) * pageSize).All(&articles)

		//手动关联查询外键		//"ArticleType"
		_, err = qs.Limit(pageSize, (pageIndex - 1) * pageSize).RelatedSel().All(&articles)
		if err != nil {
			fmt.Println("ShowIndex.Limit2 err:", err)
		}
	} else {
		//orm默认惰性查询（不会主动关联查询外键信息）
		//_, err = qs.Limit(pageSize, (pageIndex - 1) * pageSize).Filter("ArticleType__TypeName", typeName).All(&articles)

		//手动关联查询外键		//"ArticleType"
		_, err = qs.Limit(pageSize, (pageIndex - 1) * pageSize).RelatedSel().Filter("ArticleType__TypeName", typeName).All(&articles)
		if err != nil {
			fmt.Println("ShowIndex.Limit2 err:", err)
		}
	}


	var articleTypes []models.ArticleType

	//连接Redis
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("连接Redis失败, err : ", err)
		return
	}
	defer conn.Close()

	//查询redis中是否存在相应的key-value
	//res, err := redis.Bytes(conn.Do("get", "articleTypes"))
	res, err := redis.String(conn.Do("get", "articleTypes"))

	if len(res) != 0 && err != nil {
		fmt.Println("redis.Bytes err : ", err)
		return
	}
	//len(res)==0 表示redis中不存在 相应的key-value
	if len(res) == 0 {
		//从mysql查询文章类型列表
		_, err = o.QueryTable("ArticleType").All(&articleTypes)

		if err != nil {
			fmt.Println("ShowIndex.QueryTable(ArticleType) err:", err)
		}

		//序列化
		//容器，存放序列化后的结果
		var buffer bytes.Buffer
		//编码器
		enc := gob.NewEncoder(&buffer)
		//编码
		enc.Encode(&articleTypes)

		//序列化后的结果转为字节切片，保存到redis中
		//conn.Do("set", "articleTypes", buffer.Bytes())

		//序列化后的结果转为字符串，保存到redis中
		conn.Do("set", "articleTypes", buffer.String())
		//测试代码
		fmt.Println("从mysql中获取数据")
	} else {

		//解码
		//获取解码器
		dec := gob.NewDecoder(strings.NewReader(res))
		//解码
		dec.Decode(&articleTypes)
		//测试代码
		fmt.Println("从redis中获取数据", articleTypes)
	}

	this.Data["count"] = count
	this.Data["pageCount"] = int(pageCount)
	this.Data["pageIndex"] = pageIndex

	this.Data["articles"] = articles

	this.Data["articleTypes"] = articleTypes

	this.Data["typeName"] = typeName

	//视图布局
	this.Data["htmlTitle"] = "后台管理页面"
	this.Layout = "layout.html"
	this.TplName = "index.html"

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts.html"

	/*
	//测试代码
		//fmt.Println(typeName)
		//for _, v := range articles {
		//	fmt.Printf("文章名称%s, 文章类型%s\n", v.Title, v.ArticleType.TypeName)
		}
		for _, v := range articles {
			fmt.Printf("文章名称%s, 文章类型%v\n", v.Title, v.ArticleType)
		}*/
}

//展示 添加文章页面
func (this *ArticleController) ShowAddArticle()  {

	o := orm.NewOrm()
	var articleTypes []models.ArticleType

	o.QueryTable("ArticleType").All(&articleTypes)
	num, err := o.QueryTable("ArticleType").All(&articleTypes)
	if num != 0 && err != nil {
		fmt.Println("ShowAddType,QueryTable失败：", err)
		this.Redirect("/article/index", 302)
		return
	}

	this.Data["articleTypes"] = articleTypes

	this.Data["htmlTitle"] = "添加文章内容"
	this.Layout = "layout.html"
	this.TplName = "add.html"

}

//处理 文章添加事务
func (this *ArticleController) HandleAddArticle()  {

	//获取前端数据
	typeName := this.GetString("select")
	//fmt.Println("typeName = ", typeName)

	articleName := this.GetString("articleName")
	content := this.GetString("content")
	file, head, err := this.GetFile("uploadname")

	//校验数据
	if articleName == "" || content == "" || err != nil {
		this.Data["errmsg"] = "获取数据错误"
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	//对上传的文集及进行校验
	if  head.Size > 1024*1024*10 {
		this.Data["errmsg"] = "图片太大,请重新上传"
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	//判断文件类型
	ext := path.Ext(head.Filename)
	//fmt.Println(ext)
	if ext != ".jpg" && ext != ".png" {
		this.Data["errmsg"] = "文件格式不对,请重新上传"
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	//防止重名
	fileName := time.Now().Format("20060102150405")

	//err = this.SaveToFile("uploadname", "static/img/" + fileName + ext) // static/img/ 也可以
	err = this.SaveToFile("uploadname", "./static/img/" + fileName + ext)


	if err != nil {
		fmt.Println("文件保存失败:", err)
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var article models.Article
	//插入带类型的文章
	//定义指针要用new，否则报错： invalid memory address or nil pointer dereference
	articleType := new(models.ArticleType)
	articleType.TypeName = typeName
	//测试
	fmt.Println("articleType = ", *articleType)

	err = o.Read(articleType, "TypeName")
	if err != nil {
		fmt.Println("HandleAddArticle, articleType获取失败", err)
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	article.Title = articleName
	article.Content = content
	//此处不能加 .
	article.Image = "/static/img/" + fileName + ext
	article.ArticleType = articleType

	id, err := o.Insert(&article)
	if err != nil {
		fmt.Println("插入数据失败")
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}
	fmt.Println("插入成功, id = ", id)

	//返回
	//this.TplName = "index.html"
	this.Redirect("/article/index", 302)
}

//展示 详情页
func (this *ArticleController) ShowContent() {

	//获取数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		fmt.Println("id获取失败：", err)
		//this.Data["errmsg"] = err
		//this.TplName = "index.html"

		this.Redirect("/article/index", 302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	err = o.Read(&article)
	if err != nil {
		fmt.Println("文章详情不存在：", err)
		//this.TplName = "index.html"
		this.Redirect("/article/index", 302)
		return
	}

	userName := this.GetSession("userName")

//多对多插入
	//获取多对多操作对象
	m2m := o.QueryM2M(&article, "User")
	//获取插入对象，即文章
		/*上文已获取*/

	//获取要被插入的内容
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	m2m.Add(user)

/*	this.Data["title"] = article.Title
	this.Data["content"] = article.Content
	this.Data["readConunt"] = article.ReadCount
	this.Data["createTime"] = article.Time
*/
//查询多对多
	//没有去重，同一个人多次查看文章会重复显示
	//o.LoadRelated(&article, "User")
	//											   表名__字段名__文章结构体的id属性
	o.QueryTable("User").Filter("Article__Article__Id", article.Id).Distinct().All(&article.User)

	//m2ms := o.QueryM2M("User", "TypeName")



	this.Data["article"] = article

	//fmt.Println("article.Image = ", article.Image)

	this.Data["htmlTitle"] = "文章详情"
	this.Layout = "layout.html"
	this.TplName = "content.html"

	//更新阅读次数
	article.ReadCount += 1
	_, err = o.Update(&article)
	if err != nil {
		fmt.Println("更新失败：", err)
		return
	}
}

//展示 文章编辑页面
func (this *ArticleController) ShowEditArticle(){
	//获取前端数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		fmt.Println("ShowEditArticle，id获取失败:", err)

		this.Redirect("/article/index", 302)
		//this.TplName = "index.html"

		return
	}
	//查询数据库
	o := orm.NewOrm()
	article := models.Article{Id:id}

	err = o.Read(&article)
	if err != nil {
		fmt.Println("ShowEditArticle，Read:", err)
		this.Redirect("/article/index", 302)
		return
	}
	//返回数据给前端
	this.Data["article"] = article
	//指定前端视图
	this.Data["htmlTitle"] = "更新文章内容"
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

//处理 编辑文章事务
func (this *ArticleController) HandleEditArticle() {
//	1. 获取前端数据
	//获取要修改的文章的id
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("HandleEditArticle，id获取失败:", err)
		this.Redirect("/article/index", 302)
		return
	}
	//测试
	//fmt.Println("id = ", id)

	//根据id 查询数据库，文章是否存在
	o := orm.NewOrm()
	article := models.Article{Id:id}

	err = o.Read(&article)
	if err != nil {
		fmt.Println("文章不存在:", err)
		this.Redirect("/article/index", 302)
		return
	}

	//获取前端修改后的数据
	title := this.GetString("articleName")
	content := this.GetString("content")
	file, head, err := this.GetFile("uploadname")

//	2. 校验数据
	if title == "" || content == "" || err != nil {
		this.Data["article"] = article

		this.Data["errmsg"] = "获取数据错误"
		this.Layout = "layout.html"
		this.TplName = "update.html"
		return
	}
	defer file.Close()

//3. 处理数据
	article.Title = title
	article.Content = content

//对上传的文件 进行校验
//	判断文件大小
	if  head.Size > 1024*1024*10 {

		this.Data["article"] = article
		this.Data["errmsg"] = "图片太大,请重新上传"

		this.Layout = "layout.html"
		this.TplName = "update.html"
		return
	}

	//判断文件类型
	ext := path.Ext( head.Filename )
	//fmt.Println(ext)
	if ext != ".jpg" && ext != ".png" {

		this.Data["article"] = article
		this.Data["errmsg"] = "图片格式不对,请重新上传"

		this.Layout = "layout.html"
		this.TplName = "update.html"
		return
	}

	//防止重名
	imageName := time.Now().Format("20060102150405")

	//err = this.SaveToFile("uploadname", "./static/img/" + fileName + ext)
	err = this.SaveToFile("uploadname", "static/img/" + imageName + ext)

	if err != nil {

		this.Data["article"] = article
		this.Data["errmsg"] = "图片保存失败，请重试"

		this.Layout = "layout.html"
		this.TplName = "update.html"
		fmt.Println("图片保存失败:", err)
		return
	}

	//	3. 处理数据
	//article.Title = title
	//article.Content = content
	article.Image = "static/img/" + imageName + ext

	_, err = o.Update(&article)

	if err != nil {

		this.Data["article"] = article
		this.Data["errmsg"] = "文章更新失败，请重试"

		this.Layout = "layout.html"
		this.TplName = "update.html"
		fmt.Println("文章更新失败：", err)
		return
	}

	//	4. 返回数据给前端
	//this.TplName = "index.html"
	this.Redirect("/article/index", 302)

}

func (this *ArticleController) HandleDelete()  {
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("HandleDelete,获取id失败：", err)
		this.Redirect("/article/index", 302)
		return
	}

	o := orm.NewOrm()
	article := models.Article{Id:id}

	_, err = o.Delete(&article)
	if err != nil {
		fmt.Println("HandleDelete,删除失败：", err)
		this.Redirect("/article/index", 302)
		return
	}

	fmt.Println("删除成功")
	this.Redirect("/article/index", 302)
}
//展示添加类型页面
func (this *ArticleController) ShowAddType() {

	o := orm.NewOrm()
	var articleTypes []*models.ArticleType
	num, err := o.QueryTable("ArticleType").All(&articleTypes)
	if num != 0 && err != nil {
		fmt.Println("ShowAddType,QueryTable失败：", err)
		this.Redirect("/article/index", 302)
		return
	}

	this.Data["articleTypes"] = articleTypes
	//视图布局
	this.Data["htmlTitle"] = "编辑文章类型"
	this.Layout = "layout.html"
	this.TplName = "addType.html"

//	LayoutSections 传递js
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts.html"

}

// 处理添加文章类型 事务
func (this *ArticleController)HandleAddType() {

	typeName := this.GetString("typeName")

	if typeName == "" {
		fmt.Println("文章类型名获取失败")
		this.Redirect("/article/index", 302)
		return
	}

	o := orm.NewOrm()
	articleType := models.ArticleType{TypeName:typeName}

	id, err := o.Insert(&articleType)
	if err != nil {
		fmt.Println("HandleAddType,Insert失败：", err)
		this.Redirect("/article/index", 302)
		return
	}
	fmt.Println("文章类型添加成功，id=", id)

	this.Redirect("/article/addType", 302)
}

//删除文章类型
func (this *ArticleController) HandleDelType()  {

	id, err := this.GetInt("id")

	if err != nil {
		fmt.Println("HandleDelType.GetInt err：", err)
		this.Redirect("/article/index", 302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	articleType := models.ArticleType{Id:id}

	_, err = o.Delete(&articleType)

	if err != nil {
		fmt.Println("HandleDelType.Delete err：", err)
		this.Redirect("/article/index", 302)
		return
	}

	this.Redirect("/article/addType", 302)
}