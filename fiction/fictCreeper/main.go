package main

import (
	"flag"
	"fmt"
	"github.com/mm-ooto/base/common/config"
	"github.com/mm-ooto/base/common/orm"
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/constC"
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/fiction"
	"github.com/mm-ooto/summarizesStudy/fiction/fictCreeper/model"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

var (
	platformSourceType = flag.Int("platformSourceType", 1, "小说来源，默认1：新笔趣阁")
	bookUrl            = flag.String("bookUrl", "http://www.paoshuzw.com/10/10489/", "小说地址，格式：http://www.paoshuzw.com/10/10489/")
)

func init() {
	flag.Parse()
	cmdMysqlConfig := &orm.CmdMysqlConfig{
		MysqlURL: fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True&loc=Local", config.AppConfig.String("user"), config.AppConfig.String("password"), config.AppConfig.String("database"), config.AppConfig.String("charset")),
	}
	cmdMysqlConfig.Debug, _ = config.AppConfig.Bool("debug")
	cmdMysqlConfig.MysqlIdle, _ = config.AppConfig.Int("maxIdle")
	cmdMysqlConfig.MysqlMaxOpen, _ = config.AppConfig.Int("maxOpen")
	cmdMysqlConfig.MysqlConnMaxLifetime, _ = config.AppConfig.Int("maxLifetime")
	//init mysql
	orm.DBv2Init(cmdMysqlConfig)
	model.AutoMigrate()
}

func main() {
	fmt.Println(*platformSourceType,*bookUrl)
	fun := fiction.FictionCreeper(constC.Paoshuzw, "http://www.paoshuzw.com/10/10489/")
	fun()
	//crontab
	cornTab := cron.New()
	cornTab.AddFunc("0 0 0 * * ?", fun) //每天0点跑一次
	cornTab.Start()
	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal
}
