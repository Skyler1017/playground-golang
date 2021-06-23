package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:password@tcp(localhost:3306)/test")
	if err != nil {
		panic("连接数据库失败")
	}
	// db.LogMode(true)
	defer db.Close()

	// 自动迁移模式，会自动在数据库中生成topics表
	// 以及对应的列(id, title, content，created_at,
	// updated_at, deleted_at）
	db.AutoMigrate(&Topic{})

	// 创建记录
	db.Create(&Topic{Title: "1", Content: "1"})

	t := &Topic{
		Title:   "new title",
		Content: "new content",
	}
	db.Model(&Topic{}).Update(&t)
	fmt.Printf("%+v", t)
}
