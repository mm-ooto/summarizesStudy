package model

import (
	"github.com/mm-ooto/base/common/orm"
)

//热门书籍
type HotBook struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Name       string `gorm:"not null;type:varchar(10);unique" json:"name"`                        //书名
	Intro      string `gorm:"not null;type:varchar(200)" json:"intro"`                      //简介
	Author     string `gorm:"not null;type:varchar(10)" json:"author"`                      //作者
	CoverImage string `gorm:"not null;type:varchar(100)" json:"cover_image"`                //封面
	BookUrl    string `gorm:"not null;type:varchar(100)" json:"book_url"`                   //书籍地址
	Status     uint8  `gorm:"not null;default:0;index" json:"status"`                       //状态
	Category   string `gorm:"not null;type:varchar(10);default:'';index" json:"category"` //分类
	LastUpdateTime int64 `gorm:"not null;type:int(10);default:0" json:"last_update_time"`//最后更新时间
}

//书籍章节
type BookChapter struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	BookId   int    `gorm:"not null;type:int(10)" json:"book_id"`      //书籍id
	Charpter string `gorm:"not null;type:varchar(20)" json:"charpter"` //章节
	CharpterUrl string `gorm:"not null;type:varchar(50);unique" json:"charpter_url"` //章节地址
	Title    string `gorm:"not null;type:varchar(20)" json:"title"`    //标题
}

//章节内容
type ChapterContent struct {
	ID            uint   `gorm:"primary_key" json:"id"`
	BookChapterId int    `gorm:"not null"  json:"book_chapter_id"` //章节id
	Content       string `gorm:"not null" json:"content"`          //章节内容
}

func AutoMigrate() {
	orm.Odbv2.AutoMigrate(&HotBook{},&BookChapter{},&ChapterContent{})
}
