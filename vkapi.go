package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type PostResponse struct {
	Response []json.RawMessage `json:"response"`
}

type GroupResponse struct {
	Groups []Group `json:"response"`
}

type Group struct {
	Gid          int    `json:"gid"`
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	IsClosed     int    `json:"is_closed"`
	Type         string `json:"type"`
	MembersCount int    `json:"members_count"`
	Description  string `json:"description"`
	Photo        string `json:"photo"`
	PhotoMedium  string `json:"photo_medium"`
	PhotoBig     string `json:"photo_big"`
}

type Post struct {
	Id          int          `json:"id"`
	FromId      int          `json:"from_id"`
	ToId        int          `json:"to_id"`
	Date        int          `json:"date"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	SignerId    int          `json:"signer_id"`
	IsPinned    int8         `json:"is_pinned"`
	Attachment  Attachment   `json:"attachment"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type  string `json:"type"`
	Photo *Photo `json:"photo"`
	Link  *Link  `json:"link"`
}

type Photo struct {
	Pid        int    `json:"pid"`
	Aid        int    `json:"aid"`
	OwnerId    int    `json:"owner_id"`
	UserId     int    `json:"user_id"`
	Src        string `json:"src"`
	SrcBig     string `json:"src_big"`
	SrcSmall   string `json:"src_small"`
	SrcXbig    string `json:"src_xbig"`
	SrcXxbig   string `json:"src_xxbig"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Text       string `json:"text"`
	Created    int    `json:"created"`
	Access_key string `json:"access_key"`
}

type Link struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageSrc    string `json:"image_src"`
	ImageBig    string `json:"image_big"`
}

func init() {
	http.DefaultClient.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	http.DefaultClient = &http.Client{
		Timeout: time.Second * 10,
	}
}

// WallGet return array of Post by domain name
func WallGet(domain string) []Post {
	posts := make([]Post, 0, 20)
	url := "http://api.vk.com/method/wall.get?domain=" + domain
	resp, err := http.Get(url)

	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var postRes PostResponse
			err := json.Unmarshal(body, &postRes)
			if err == nil {
				for i := range postRes.Response {
					var post Post
					if i > 0 {
						err := json.Unmarshal(postRes.Response[i], &post)
						if err == nil {
							posts = append(posts, post)
						}
					}
				}
			}
		}
	}
	return posts
}

// GroupsGetById return groups, where name = shortname or vk public id
func GroupsGetById(name string) (groups []Group) {
	url := "https://api.vk.com/method/groups.getById?group_id=" + name + "&fields=members_count,description"
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var groupRes GroupResponse
			err := json.Unmarshal(body, &groupRes)
			if err == nil {
				groups = groupRes.Groups
			}
		}
	}
	return
}
