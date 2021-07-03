package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"

	"money/core/config"
)

var db *gorm.DB

type FetchTask struct {
	gorm.Model
	Project string `gorm:"type:varchar(32);index:idx_name,unique"`
	Url     string `gorm:"type:varchar(255);index:idx_name,unique"`
	Parser  string `gorm:"type:varchar(255)"`
	Page    string `gorm:"type:LONGTEXT"`
}

func AddFetchTask(task *FetchTask) error {
	return db.Model(&FetchTask{}).FirstOrCreate(task, task).Error
}

func GetFetchTask(task *FetchTask) error {
	err := db.Model(&FetchTask{}).Where("deleted_at is null").Take(task).Error
	return err
}

func UpdateFetchTask(task *FetchTask) error {
	return db.Model(task).Updates(task).Error
}

func SetUp(conf *config.Config) {

	var err error

	dbs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DbName,
	)
	fmt.Println(dbs)

	db, err = gorm.Open("mysql", dbs)
	if err != nil {
		panic(err)
	}
	// 设置全局表名禁用复数
	db.SingularTable(true)

	if err := db.Model(&FetchTask{}).CreateTable(&FetchTask{}); err != nil {
		logrus.Error(err)
	}
}
