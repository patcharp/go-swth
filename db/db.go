package db

import (
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type Database struct {
	Config Config
	Ctx    *gorm.DB
	Job    *scheduler.Job
}

func NewDatabase(
	host string,
	port string,
	username string,
	password string,
	name string,
) Database {
	return Database{
		Config: Config{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Name:     name,
		},
	}
}

func (db *Database) Connect(prod bool) error {
	addr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		db.Config.Username,
		db.Config.Password,
		db.Config.Host,
		db.Config.Port,
		db.Config.Name,
	)
	var err error
	db.Ctx, err = gorm.Open("mysql", addr)
	if err != nil {
		return err
	}
	if err := db.startKeepAlive(); err != nil {
		return err
	}
	db.Ctx.LogMode(!prod)
	return nil
}

func (db *Database) Close() error {
	_ = db.stopKeepAlive()
	if err := db.Ctx.Close(); err != nil {
		return err
	}
	return nil
}

func (db *Database) MigrateDatabase(tables []interface{}) error {
	tx := db.Ctx.Begin()
	for _, t := range tables {
		if err := tx.AutoMigrate(t).Error; err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (db *Database) startKeepAlive() error {
	var err error
	db.Job, err = scheduler.Every(15).Seconds().Run(func() {
		if err := db.Ctx.DB().Ping(); err != nil {
			log.Errorln("Database keepalive error -:", err)
		}
	})
	return err
}

func (db *Database) stopKeepAlive() error {
	if db.Job != nil {
		db.Job.Quit <- true
	}
	return nil
}
