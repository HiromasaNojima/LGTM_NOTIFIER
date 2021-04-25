package db

import (
	"database/sql"
	"lgtm/configs"
	"lgtm/pkg/qiita"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize(conf configs.Config, items []qiita.Item) {
	var err error
	db, err = sql.Open("mysql", conf.DbDataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Second)
}

func InsertIntoItem(item qiita.Item) (sql.Result, error) {
	return db.Exec("INSERT INTO item(id, title, likes_count) values(?, ?, ?)", item.ID, item.Title, item.LikesCount)
}

func UpdateItem(item qiita.Item) (sql.Result, error) {
	return db.Exec("UPDATE item SET likes_count = ? WHERE id = ?", item.LikesCount, item.ID)
}
