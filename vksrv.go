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

	"github.com/robfig/cron"
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
	c := cron.New()
	//c.AddFunc("@every 0h12m05s", func() { fmt.Println("Every 5s thirty") })
	c.Start()
	c.Stop()
	domains := vkdomains()
	for i := range domains {
		//log.Println(domains[i])
		saveposts(domains[i])
	}
	time.Sleep(1000 * time.Millisecond)
	log.Println("end")
}

func saveposts(domain string) {
	log.Println(domain)
	posts := WallGet(domain)
	for i := range posts {
		post := posts[i]
		if post.MarkedAsAds == 0 {
			url := fmt.Sprintf("http://badtobefat.ru/bolt/%d/%d", post.OwnerID, post.Id)
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
