package dao

import (
	"chat/model"
	"chat/util"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"xorm.io/core"
)

var messageEngine *xorm.Engine

func init() {
	connect := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", util.ConfigInfo.MySQL.User, util.ConfigInfo.MySQL.Password, util.ConfigInfo.MySQL.Host, util.ConfigInfo.MySQL.Port, util.ConfigInfo.MySQL.Database)
	var err error
	messageEngine, err = xorm.NewEngine("mysql", connect)
	if err != nil {
		log.Fatalf(err.Error())
	}
	messageEngine.ShowSQL(true)
	messageEngine.Logger().SetLevel(core.LOG_DEBUG)
	err = messageEngine.Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Info("success to connect to MySQL,connect info :", connect)
	err = messageEngine.Sync2(new(model.Message))
	if err != nil {
		log.Errorf(err.Error())
	}
}
func InsertMessage(from string, to string, content string, read int) error {
	message := model.Message{
		From:    from,
		To:      to,
		Content: content,
		IsRead:  read,
	}
	_, err := messageEngine.Table("message").Insert(&message)
	return err
}
func GetNotReadMessage(to string) (error, []model.Message) {
	list := make([]model.Message, 0)
	err := messageEngine.Table("message").Where("is_read=0").Find(&list)
	return err, list
}
func GetRecentMessage(who string) (error, []model.Message) {
	list := make([]model.Message, 0)
	err := messageEngine.Table("message").Where("from=? or to=?", who, who).Find(&list)
	return err, list
}
func ReadSomeMessage(message model.Message) error {
	_, err := messageEngine.Table("message").ID(message.Id).Update(&message)
	return err
}
