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

	"github.com/robfig/cron"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	api      = "http://badtobefat.ru/bolt"
	users    = "/usertg/"
	pubNames = "/pubNames/"
	pubSubTg = "/pubSubTg/"
	someErr  = "Something going wrong. Try later.."
	hello    = "Hello, %username!\nJust drop me link on vk public and i will send messages from it.\nExample: https://vk.com/myakotkapub"
)

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
	var bot, wrbot *tgbotapi.BotAPI
	bot, err = tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		log.Fatal(err)
	}
	wrbot, err = tgbotapi.NewBotAPI(wrtoken)
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New()
	//c.AddFunc("@every 0h12m05s", func() { fmt.Println("Every 5s thirty") })
	c.Start()
	c.Stop()
	domains := vkdomains()
	for i := range domains {
		log.Println(domains[i])
		//saveposts(domains[i])
	}
	p := getpost()
	txt := strings.Replace(p.Text, "&lt;br&gt;", "\n", -1)
	log.Println(txt)
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
				bb := tgbotapi.FileBytes{Name: "image.jpg", Bytes: b}
				msg := tgbotapi.NewPhotoUpload(-1001067277325, bb)
				res, err := wrbot.Send(msg)
				if err == nil {
					fmt.Printf("%+v", res.MessageID)
					fmt.Printf("%+v", res.Photo)
					msg := tgbotapi.NewForward(1263310, -1001067277325, res.MessageID)
					r, err := bot.Send(msg)
					log.Println("fwd", r, err)
				}
			}
			fmt.Printf("%+v - -", bot.GetChat)
			//msg := tgbotapi.NewForward(1263310, -1001119114536, 1)
			//_, err := bot.Send(msg)

			//resp, err := bot.UploadFile("sendPhoto", params, "some", nil)
			//res, err := bot.UploadFile("1", nil, "", nil)
			//msg := tgbotapi.Upload//.ChatUploadPhoto() .NewPhotoUpload(-1001119114536, photo)
			//msg.Caption = "Test"
			//res, err := bot.Send(msg)

			//log.Println(resp, err)
			//fmt.Printf("%+v", photo)
		}
	}
	log.Println("end")
}

func saveposts(domain string) {
	log.Println(domain)
	posts := WallGet(domain)
	for i := range posts {
		post := posts[i]
		if post.MarkedAsAds == 0 {
			url := fmt.Sprintf("http://badtobefat.ru/bolt/%d/%s", post.OwnerID*(-1), fmt.Sprintf("%010d", post.Id))
			log.Println("url", url)
			b, err := json.Marshal(post)
			if err == nil {
				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err == nil {
					defer resp.Body.Close()

				} else {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		}
		log.Println(post.Id)
		if i == 2 {
			//break
		}
	}
}

func httpGet(url string) []byte {
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

func vkdomains() (domains []string) {
	url := api + "/pubNames/Last?cnt=1000000&order=desc&vals=false"
	log.Println("vkdomains", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &domains)
		if err != nil {
			log.Println(err)
		}
	}
	return
}
