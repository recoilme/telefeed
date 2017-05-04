package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

const (
	MaxInt   = int(^uint(0) >> 1)
	MinInt   = -MaxInt - 1
	api      = "http://badtobefat.ru/bolt"
	users    = "/usertg/"
	pubNames = "/pubNames/"
	pubSubTg = "/pubSubTg/"
	LastPost = "/vkpublastpost/"
	someErr  = "Something going wrong. Try later.."
	hello    = "Hello, %username!\nJust drop me link on vk public and i will send messages from it.\nExample: https://vk.com/myakotkapub"
)

var (
	bot, wrbot *tgbotapi.BotAPI
)

func init() {
	log.SetOutput(ioutil.Discard)
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

func main() {
	log.Println("main")
	var err error
	tlgrmtoken, err := ioutil.ReadFile("tokentg")
	if err != nil {
		log.Fatal(err)
	}
	writetoken, err := ioutil.ReadFile("vkwriter")
	if err != nil {
		log.Fatal(err)
	}
	tgtoken := strings.Replace(string(tlgrmtoken), "\n", "", -1)
	wrtoken := strings.Replace(string(writetoken), "\n", "", -1)
	bot, err = tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		log.Fatal(err)
	}
	wrbot, err = tgbotapi.NewBotAPI(wrtoken)
	wrbot.Debug = true
	if err != nil {
		log.Fatal(err)
	}
	_ = bot
	_ = wrbot
	c := cron.New()
	c.AddFunc("@every 0h10m00s", parseVk)
	c.Start()
	c.Stop()

	log.Println("end")
}

func parseVk() {
	domains := vkdomains()
	for i := range domains {
		domain := domains[i]
		log.Println(domain.ScreenName)
		users := domUsers(domains[i])
		saveposts(domain, users)
	}
}

func domUsers(domain Group) (users map[int]bool) {
	mask := api + pubSubTg + "%d"
	url := fmt.Sprintf(mask, domain.Gid)
	log.Println(url)
	b := httpGet(url)
	if b != nil {
		json.Unmarshal(b, &users)
	}
	return users
}

func lastPostIdGet(domain Group) int {
	postId := MinInt
	mask := api + LastPost + "%d"
	url := fmt.Sprintf(mask, domain.Gid)
	b := httpGet(url)
	if b != nil {
		json.Unmarshal(b, &postId)
	}
	return postId
}

func lastPostIdSet(domain Group, lastPostId int) int {
	postId := MinInt
	mask := api + LastPost + "%d"
	url := fmt.Sprintf(mask, domain.Gid)
	b, err := json.Marshal(lastPostId)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		postId = lastPostId
	} else {
		log.Println(err)
	}
	return postId
}

func saveposts(domain Group, users map[int]bool) {
	log.Println(domain)
	var lastPost = lastPostIdGet(domain)
	log.Println("last", lastPost)
	posts := WallGet(domain.ScreenName)
	last := len(posts) - 1
	for i := range posts {
		post := posts[last-i]
		if post.Id <= lastPost {
			continue
		}
		lastPost = lastPostIdSet(domain, post.Id)
		url := fmt.Sprintf("http://badtobefat.ru/bolt/%d/%s", post.OwnerID*(-1), fmt.Sprintf("%010d", post.Id))
		b, _ := json.Marshal(post)
		httpPut(url, b)
		log.Println(post.Id)
		pubpost(domain, post, users)
		break
	}
}

func httpGet(url string) []byte {
	log.Println("httpGet", url)
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return body
		}
	}
	return nil
}

func httpPut(url string, b []byte) {
	log.Println("httpPut", url)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
}

func getpost() (post Post) {
	postid := "126993367/0000001170"

	url := api + "/" + postid
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err := json.Unmarshal(body, &post)
			if err == nil {
				return
			}
		}
	}
	return
}

func vkdomains() (domains []Group) {
	var domainNames []string
	url := api + "/pubNames/Last?cnt=1000000&order=desc&vals=false"
	log.Println("vkdomains", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &domainNames)
		if err == nil {
			for i := range domainNames {
				domainName := domainNames[i]
				//log.Println("domainName", domainName)
				b := httpGet(api + pubNames + domainName)
				if b != nil {
					var domain Group
					err := json.Unmarshal(b, &domain)
					if err == nil {
						domains = append(domains, domain)
					}
				}
			}
		} else {
			log.Println(err)
		}

	}
	return
}

func pubpost(domain Group, p Post, users map[int]bool) {
	log.Println("pubpost", p.Id)
	var vkcnt int64 = -1001067277325 //myakotka
	//var fwd int64 = 366035536        //telefeed

	var t = strings.Replace(p.Text, "&lt;br&gt;", "\n", -1)
	if t != "" {
		t = t + "\n"
	}
	link := fmt.Sprintf("vk.com/wall%d_%d", domain.Gid*(-1), p.Id)
	tag := strings.Replace(domain.ScreenName, ".", "", -1)
	txt := fmt.Sprintf("%s#%s 🔗 %s", t, tag, link)
	log.Println("txt:", txt)
	if len(p.Attachments) == 0 || len(txt) > 250 {
		res, err := wrbot.Send(tgbotapi.NewMessage(vkcnt, txt))
		if err == nil {
			for user := range users {
				log.Println(user)
				bot.Send(tgbotapi.NewForward(int64(user), vkcnt, res.MessageID))
			}
		}
	}
	for i := range p.Attachments {
		att := p.Attachments[i]
		log.Println(att.Type)
		switch att.Type {
		case "photo":
			var photo = att.Photo.Photo807
			if photo == "" {
				photo = att.Photo.Photo604
			}
			log.Println(photo)
			b := httpGet(photo)
			if b != nil {
				bb := tgbotapi.FileBytes{Name: photo, Bytes: b}
				msg := tgbotapi.NewPhotoUpload(vkcnt, bb)
				msg.Caption = txt
				res, err := wrbot.Send(msg)
				if err == nil {
					for user := range users {

						bot.Send(tgbotapi.NewForward(int64(user), vkcnt, res.MessageID))
					}
				}
			}
		case "video":
			//fmt.Printf("%+v\n", att.Video)
			if att.Video.Duration > 600 {
				continue
			}
			urlv := fmt.Sprintf("https://vk.com/video%d_%d", att.Video.OwnerID, att.Video.ID)
			b := httpGet(urlv)
			if b != nil {
				cnt := string(b)
				var pos360 = strings.Index(cnt, ".360.mp4")
				if pos360 < 0 {
					pos360 = strings.Index(cnt, ".240.mp4")
				}
				if pos360 < 0 || pos360 < 200 {
					break
				}
				poshttp := strings.Index(cnt[pos360-200:], "https") + pos360 - 200 //cnt.find("https:",pos360-200)
				if poshttp > 0 {
					s := strings.Replace(cnt[poshttp:pos360+8], "\\/", "/", -1)
					if s != "" {
						//post video
						vidb := httpGet(s)
						bb := tgbotapi.FileBytes{Name: s, Bytes: vidb}
						msg := tgbotapi.NewVideoUpload(vkcnt, bb)
						msg.Caption = txt
						res, err := wrbot.Send(msg)
						if err == nil {
							for user := range users {
								bot.Send(tgbotapi.NewForward(int64(user), vkcnt, res.MessageID))
							}
						}
					}
				}
			}
		case "doc":
			//fmt.Printf("%+v\n", att.Doc)
			b := httpGet(att.Doc.URL)
			if b != nil {
				bb := tgbotapi.FileBytes{Name: "tmp." + att.Doc.Ext, Bytes: b}
				msg := tgbotapi.NewDocumentUpload(vkcnt, bb)
				msg.Caption = txt
				res, err := wrbot.Send(msg)
				if err == nil {
					for user := range users {
						bot.Send(tgbotapi.NewForward(int64(user), vkcnt, res.MessageID))
					}
				}
			}
		case "link":
			//fmt.Printf("%+v\n", att.Link)
			var desc = ""
			if len(txt) <= 250 {
				desc = att.Link.URL + "\n" + txt
			} else {
				desc = att.Link.URL
			}
			msg := tgbotapi.NewMessage(vkcnt, desc)
			res, err := wrbot.Send(msg)
			if err == nil {
				for user := range users {

					bot.Send(tgbotapi.NewForward(int64(user), vkcnt, res.MessageID))
				}
			}

		}
	}

}
