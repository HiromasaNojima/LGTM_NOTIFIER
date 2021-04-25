package line

import (
	"lgtm/configs"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Notify(conf configs.Config, message string) {
	form := url.Values{}
	form.Add("message", message)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", conf.LineNotifyURL, body)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+conf.LineAccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
}
