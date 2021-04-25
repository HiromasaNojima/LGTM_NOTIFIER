package qiita

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"lgtm/configs"
)

type Item struct {
	RenderedBody  string    `json:"rendered_body"`
	Body          string    `json:"body"`
	Coediting     bool      `json:"coediting"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	Group         struct {
		CreatedAt time.Time `json:"created_at"`
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Private   bool      `json:"private"`
		UpdatedAt time.Time `json:"updated_at"`
		URLName   string    `json:"url_name"`
	} `json:"group"`
	ID             string `json:"id"`
	LikesCount     int    `json:"likes_count"`
	Private        bool   `json:"private"`
	ReactionsCount int    `json:"reactions_count"`
	Tags           []struct {
		Name     string   `json:"name"`
		Versions []string `json:"versions"`
	} `json:"tags"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	User      struct {
		Description       string `json:"description"`
		FacebookID        string `json:"facebook_id"`
		FolloweesCount    int    `json:"followees_count"`
		FollowersCount    int    `json:"followers_count"`
		GithubLoginName   string `json:"github_login_name"`
		ID                string `json:"id"`
		ItemsCount        int    `json:"items_count"`
		LinkedinID        string `json:"linkedin_id"`
		Location          string `json:"location"`
		Name              string `json:"name"`
		Organization      string `json:"organization"`
		PermanentID       int    `json:"permanent_id"`
		ProfileImageURL   string `json:"profile_image_url"`
		TeamOnly          bool   `json:"team_only"`
		TwitterScreenName string `json:"twitter_screen_name"`
		WebsiteURL        string `json:"website_url"`
	} `json:"user"`
	PageViewsCount int `json:"page_views_count"`
}

func GetAllItems(conf configs.Config, page int, items []Item) ([]Item, error) {
	temp, err := GetItems(conf, page)
	if err != nil {
		return nil, err
	}

	// 取得件数が０になるまで再帰することで全記事データの取得を行う。
	if len(temp) == 0 {
		return items, nil
	}

	for _, item := range temp {
		if item.Private {
			// 非公開の記事はスキップ
			continue
		}
		items = append(items, item)
		log.Println(fmt.Sprintf("取得した記事 : %s(%s) LGTM(%d)", item.Title, item.ID, item.LikesCount))
	}

	page++
	return GetAllItems(conf, page, items)
}

func GetItems(conf configs.Config, page int) (items []Item, err error) {
	req, err := http.NewRequest("GET", "https://qiita.com/api/v2/users/"+conf.QiitaUserName+"/items?page="+strconv.Itoa(page)+"&per_page=20", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+conf.QiitaAccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}
