package configs

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LineAccessToken  string `json:"line_access_token"`
	LineNotifyURL    string `json:"line_notify_url"`
	QiitaUserName    string `json:"qiita_user_name"`
	QiitaAccessToken string `json:"qiita_access_token"`
	DbDataSourceName string `json:"db_data_source_name"`
	LogPath          string `json:"log_path"`
}

func ReadConfig(path string) (conf Config) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("loadConfig os.Open err:", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&conf)
	return conf
}
