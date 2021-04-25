package main

import (
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"lgtm/configs"
	"lgtm/pkg/db"
	"lgtm/pkg/line"
	"lgtm/pkg/qiita"
)

var itemMap = map[string]qiita.Item{}

func initializeMap(items []qiita.Item) {
	for _, item := range items {
		itemMap[item.ID] = item
	}
}

func errorNotifyAndTerminate(conf configs.Config, err error) {
	msg := createFatalErrorMessage(err)
	line.Notify(conf, msg)
	log.Fatalln(msg)
}

func errorNotify(conf configs.Config, err error) {
	msg := fmt.Sprintf("エラーが発生しました。 detail:%s", err.Error())
	line.Notify(conf, msg)
	log.Infoln(msg)
}

func createFatalErrorMessage(err error) string {
	return fmt.Sprintf("エラーが発生しました。 システム終了します。detail:%s", err.Error())
}

func main() {
	// config.jsonへの絶対パスを引数で入力 -> 読み取り
	flag.Parse()
	conf := configs.ReadConfig(flag.Arg(0))
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(&lumberjack.Logger{
		Filename:   conf.LogPath, // ファイル名
		MaxSize:    500,          // ローテーションするファイルサイズ(megabytes)
		MaxBackups: 3,            // 保持する古いログの最大ファイル数
		MaxAge:     365,          // 古いログを保持する日数
		LocalTime:  true,         // バックアップファイルの時刻フォーマットをサーバローカル時間指定
		Compress:   true,         // ローテーションされたファイルのgzip圧縮
	})

	items, err := qiita.GetAllItems(conf, 1, []qiita.Item{})
	if err != nil {
		log.Fatalln(createFatalErrorMessage(err))
	}

	db.Initialize(conf, items)
	initializeMap(items)

	for {
		items, err := qiita.GetAllItems(conf, 1, []qiita.Item{})
		if err != nil {
			errorNotify(conf, err)
			time.Sleep(time.Minute * 60)
			continue
		}

		for _, item := range items {
			// 記事投稿してから初回の取得、mapに該当記事のデータ存在しない場合
			beforeItem, ok := itemMap[item.ID]
			if !ok {
				itemMap[item.ID] = item
				_, err = db.InsertIntoItem(item)
				if err != nil {
					errorNotifyAndTerminate(conf, err)
				}

				continue
			}

			// LGTMの数が変化した場合
			if item.LikesCount != beforeItem.LikesCount {
				// 通知
				line.Notify(conf, fmt.Sprintf("『%s』(%s)のLGTM数が変化しました。%d -> %d", item.Title, item.ID, beforeItem.LikesCount, item.LikesCount))
				// ローカルで持ってる記事データの更新
				itemMap[item.ID] = item
				_, err = db.UpdateItem(item)
				if err != nil {
					errorNotifyAndTerminate(conf, err)
				}
			}
		}

		time.Sleep(time.Second * 60)
	}
}
