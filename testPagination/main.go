package main

import (
	"fmt"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          uint `gorm:"primary_key"`
	Title       string
	Body        string
	PublishedAt time.Time
}

func main() {
	db, err := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/test?parseTime=true"))
	if err != nil {
		panic(fmt.Errorf("db connection error: %s", err))
	}

	if err := db.AutoMigrate(&Post{}); err != nil {
		panic(err)
	}

	// //生成测试数据
	//for i := 0; i < 100000; i++ {
	//	p := &Post{
	//		ID: uint(i + 1),
	//		Title:       fmt.Sprint("post-", i),
	//		Body:        "bodybodybody",
	//		PublishedAt: time.Now(),
	//	}
	//	err := db.Create(p).Error
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	var posts []Post
	q := db.Model(Post{}).Where("published_at < ?", time.Now())
	p := paginator.New(adapter.NewGORMAdapter(q), 10)
	p.SetPage(2)

	if err = p.Results(&posts); err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}
}
