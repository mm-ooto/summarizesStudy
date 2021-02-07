package fiction

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mm-ooto/base/common/functions"
	"github.com/mm-ooto/base/common/orm"
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/model"
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/utils"
	"log"
	"strings"
	"time"
)

type PaoShuzw struct {
}

func (p *PaoShuzw) FictionHotBook(bookUrl string) (err error) {
	doc, err := utils.GetDocument(bookUrl)
	if err != nil {
		return err
	}
	getHotBook(doc, bookUrl)
	return
}

//小说基础信息入库
func getHotBook(doc *goquery.Document, bookUrl string) {
	log.Println("-------------进行中-------------")
	hotBook := &model.HotBook{}
	mainInfo := doc.Find("#maininfo")
	hotBook.Name = mainInfo.Find("#info h1").Text()
	db := orm.DBv2Begin()
	defer db.DBv2Rollback()
	db.Model(&model.HotBook{}).Where("book_url=?", bookUrl).Last(hotBook)

	mainInfo.Find("#info p").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			hotBook.Author = strings.Split(selection.Text(), "：")[1]
		}
		if i == 2 {
			lastUpdateTime := strings.Split(selection.Text(), "：")[1]
			lut, _ := time.Parse("2006-01-02 15:04:05", lastUpdateTime)
			hotBook.LastUpdateTime = lut.Unix()
			return
		}
	})
	hotBook.BookUrl = bookUrl
	mainInfo.Find("#intro p").Each(func(i int, selection *goquery.Selection) {
		if i == 1 {
			hotBook.Intro = selection.Text()
			return
		}
	})
	hotBook.CoverImage, _ = doc.Find("#fmimg img").Attr("src")
	doc.Find(".con_top a").Each(func(i int, selection *goquery.Selection) {
		if i == 2 {
			hotBook.Category = selection.Text()
			return
		}
	})
	gdb := db.Model(&model.HotBook{})
	if hotBook.ID > 0 {
		gdb = gdb.Where("id=?", hotBook.ID)
	}
	gdb.Save(hotBook).First(hotBook) //create or update
	db.DBv2Commit()
	chn := make(chan uint)
	go getBookChapter(chn, hotBook, doc)
	<-chn
	log.Println("-------------已结束-------------")
	return
}

//小说章节信息入库
func getBookChapter(chn chan uint, hotBook *model.HotBook, doc *goquery.Document) {
	db := orm.DBv2Begin()
	defer func() {
		db.DBv2Rollback()
		chn <- hotBook.ID
		close(chn)
		log.Printf("getBookChapter 程序退出！bookName：【%s】\n", hotBook.Name)
	}()
	bookChapters := make([]model.BookChapter, 0)
	bookChapter := model.BookChapter{}
	var charpter []string
	db.Model(&model.BookChapter{}).Select("charpter").Scan(&charpter)
	doc.Find("#list dl dd").Each(func(i int, selection *goquery.Selection) {
		aSelection := selection.Find("a")
		atext := strings.Split(aSelection.Text(), " ")
		//if len(atext) == 2 && i < 60 {//用来限制爬取的章节
		if len(atext) == 2 && !functions.InArray(atext[0], charpter) {
			charpterUrl, _ := aSelection.Attr("href")
			bookChapter = model.BookChapter{
				BookId:      int(hotBook.ID),
				Charpter:    atext[0],
				CharpterUrl: charpterUrl,
				Title:       atext[1],
			}
			db.Model(&model.BookChapter{}).Save(&bookChapter).Select("id").First(&bookChapter)
			bookChapters = append(bookChapters, bookChapter)
			bookChapter = model.BookChapter{}
		}
		//}
	})
	db.DBv2Commit()
	if len(bookChapters) > 0 {
		log.Println(len(bookChapters))
		chapterContentChan := make(chan string, len(bookChapters))
		existChan := make(chan bool, 1)
		go getChapterContent(chapterContentChan, bookChapters)
		go readchapterContent(chapterContentChan, existChan)
		for {
			_, ok := <-existChan
			if !ok {
				break
			}
		}
	}
	return
}

func readchapterContent(chapterContentChan chan string, existChan chan bool) {
	for {
		v, ok := <-chapterContentChan
		if !ok {
			break //chan close
		}
		log.Printf("成功接收数据：%s\n", v)
	}
	existChan <- true
	close(existChan)
}

//小说章节内容入库
func getChapterContent(chapterContentChan chan string, bookChapters []model.BookChapter) {
	db := orm.DBv2Begin()
	defer func() {
		db.DBv2Rollback()
		close(chapterContentChan)
	}()
	chapterContents := make([]map[string]interface{}, 0)
	for _, v := range bookChapters {
		chapterContentChan <- fmt.Sprintf("%s %s", v.Charpter, v.Title)
		url := fmt.Sprintf("http://www.paoshuzw.com%s", v.CharpterUrl)
		doc, err := utils.GetDocument(url)
		if err != nil {
			fmt.Println("GetDocument error:", err.Error())
			return
		}
		content := doc.Find("#content").Text()
		chapterContents = append(chapterContents, map[string]interface{}{
			"BookChapterId": v.ID,
			"Content":       content,
		})
	}
	gdb := db.Model(&model.ChapterContent{}).Create(chapterContents)
	if err := gdb.Error; err != nil {
		log.Printf("err:%s", err.Error())
		return
	}
	db.DBv2Commit()
	return
}
