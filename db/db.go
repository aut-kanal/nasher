package db

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/nasher/configuration"
)

var (
	db   *gorm.DB
	once sync.Once
)

func getSqliteDB(path string) (*gorm.DB, error) {
	return gorm.Open("sqlite3", path)
}

func getPostgresDB(host, port, name, user, password string) (*gorm.DB, error) {
	return gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, password,
	))
}

// GetInstance returns a gorm DB singleton
func GetInstance() *gorm.DB {
	once.Do(func() {
		err := initDB()
		if err != nil {
			logrus.Errorf("can't init database")
		}
	})
	return db
}

// TODO: call this func when config changes
func initDB() error {
	var err error
	dialect := configuration.GetInstance().GetString("db.dialect")
	switch dialect {
	case "postgres":
		host := configuration.GetInstance().GetString("db.host")
		port := configuration.GetInstance().GetString("db.port")
		dbname := configuration.GetInstance().GetString("db.name")
		user := configuration.GetInstance().GetString("db.user")
		password := configuration.GetInstance().GetString("db.password")
		logrus.Debug("initializing Postgres connection")
		db, err = getPostgresDB(host, port, dbname, user, password)
	case "sqlite3":
		filePath := configuration.GetInstance().GetString("db.path")
		logrus.Debug("initializing Sqlite3 connection")
		db, err = getSqliteDB(filePath)
	}

	return err
}

// Close singleton DB instance
func Close() error {
	return db.Close()
}
