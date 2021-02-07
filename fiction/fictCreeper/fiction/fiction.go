package fiction

import (
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/constC"
	"log"
)

type Fiction interface {
	FictionHotBook(bookUrl string) error
}

func NewFiction(source constC.PlatformSourceType) (f Fiction, err error) {
	if !constC.CheckPlatformSourceInSlice(source){
		log.Println("暂不支持该网站！")
	}
	switch source {
	case constC.Paoshuzw:
		f = new(PaoShuzw)
	default:
	}
	return
}

//platformSourceType:平台来源;bookUrl:小说地址
func FictionCreeper(platformSourceType constC.PlatformSourceType,bookUrl string)func()  {
	fic, err := NewFiction(platformSourceType)
	if err != nil {
		log.Printf("newFiction error:%s\n", err.Error())
	}
	return func() {
		fic.FictionHotBook(bookUrl)
	}
}
