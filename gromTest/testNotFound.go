package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Topic struct {
	gorm.Model
	Title   string
	Content string
}

func main2() {
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

	var topic Topic
	//err = db.Where("id=?", 0).Find(&topic).Error
	//fmt.Println("db.Where(\"id=?\", 0).Find(&topic).Error: ", err)
	//fmt.Println("	topic的id:", topic.ID)
	//// First
	//err = db.Where("id=?", 0).First(&topic).Error
	//fmt.Println("db.Where(\"id=?\", 0).First(&topic).Error: ", err)
	//fmt.Println("	topic的id:", topic.ID)
	//// Last
	//err = db.Where("id=?", 0).Last(&topic).Error
	//fmt.Println("db.Where(\"id=?\", 0).Last(&topic).Error", err)
	//fmt.Println("	topic的id:", topic.ID)
	//
	//// 查找topic列表
	//// 这里的err值为nil
	//var topics []Topic
	//err = db.Where("id=?", 0).Find(&topics).Error
	//fmt.Println("db.Where(\"id=?\", 0).Find(&topics).Error: ", err)
	//fmt.Println("	topics的长度", len(topics))
	//// First
	//err = db.Where("id=?", 0).First(&topics).Error
	//fmt.Println("db.Where(\"id=?\", 0).First(&topics).Error: ", err)
	//fmt.Println("	topics的长度", len(topics))
	//// Last
	//err = db.Where("id=?", 0).Last(&topics).Error
	//fmt.Println("db.Where(\"id=?\", 0).Last(&topics).Error: ", err)
	//fmt.Println("	topics的长度", len(topics))

	//cnt := 0
	err = db.Model(topic).Where("id=?", 0).Update(&topic).Error
	fmt.Println("db.Where(\"id=?\", 0).Count(&cnt).Error: ", err)
}
