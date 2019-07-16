package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"unique"`
	Pwd string

	//设置多对多的反向关系
	Article []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(50);unique"`
	Content string `orm:"size(500)"`
	Time time.Time `orm:"type(datatime);auto_now_add"`	// auto_now
										// auto_now 每次 model 保存时都会对时间自动更新 //auto_now_add 第一次保存时才设置时间
	ReadCount int `orm:"default(0)"`
	Image string `orm:"null"`
	//设置一对多关系 外健约束	 `orm:"外健；主表删除时，该字段的值不变"`
	ArticleType *ArticleType `orm:"rel(fk);on_delete(do_nothing)"`

	//设置多对多关系(m2m: many2many)
	User []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int `orm:"pk;auto"`

	TypeName string `orm:"size(50);unique"`
	//设置一对多的反向关系
	Article []*Article `orm:"reverse(many)"`
}

func init() {
//	注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8&loc=Local")

//	注册表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

//	跑起来(创建表)
//				参1：表别名；参2：强制更新表属性，会清空表数据；参3：是否显示命令执行过程中的详细信息
//	orm.RunSyncdb("default", true, true)
	orm.RunSyncdb("default", false, true)

}
