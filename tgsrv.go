package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

var (
	bot *tgbotapi.BotAPI
)

func catch(e error) {
	if e != nil {
		log.Panic(e.Error)
	}
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

func main() {
	var err error
	tlgrmtoken, err := ioutil.ReadFile("tokentg")
	catch(err)
	tgtoken := strings.Replace(string(tlgrmtoken), "\n", "", -1)
	bot, err = tgbotapi.NewBotAPI(tgtoken)
	catch(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	c := cron.New()
	log.Println(42)
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 0h01m", func() { fmt.Println("Every hour thirty") })
	c.AddFunc("@every 0h00m05s", func() { fmt.Println("Every 5s thirty") })
	c.Start()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			user := update.Message.From
			if userNew(user) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, hello))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, someErr))
			}
		default:
			//sendRply(update.Message, "Пошел на хуй")
			msg := update.Message.Text
			pubFind(update.Message, msg)
		}
	}
	c.Stop()

}

func userNew(user *tgbotapi.User) bool {
	url := api + users + strconv.Itoa(user.ID)
	log.Println("userNew", url)
	b, err := json.Marshal(user)
	catch(err)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	catch(err)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func sendRply(message *tgbotapi.Message, txt string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	msg.ReplyToMessageID = message.MessageID
	log.Printf("1" + msg.Text)
	bot.Send(msg)
}

func pubFind(msg *tgbotapi.Message, txt string) {
	log.Println("pubFind")
	words := strings.Split(txt, " ")
	for i := range words {
		word := words[i]
		urls, err := url.Parse(word)
		catch(err)
		switch urls.Host {
		case "vk.com":
			parts := strings.Split(urls.Path, "/")
			for j := range parts {
				if parts[j] != "" {
					domain := parts[j]
					log.Println(domain)
					//bot.Send(tgbotapi.NewMessage(chatId, "Found vk domain:'"+parts[j]+"'"))
					groupDb := pubDbGet(domain)
					if groupDb.Gid == 0 {
						// public not found
						groups := GroupsGetById(domain)
						if len(groups) > 0 {
							// we have group
							groupVk := groups[0]
							// save group to DB
							if pubDbSet(groupVk) {
								// new group set
								pubSubTgAdd(groupVk, msg)
							} else {
								// group not set
								bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error create domain:'"+domain+"'"))
							}
						} else {
							// group not found
							bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error vk domain:'"+domain+"'"+" not found"))
						}

					} else {
						// public exists

						pubSubTgAdd(groupDb, msg)
					}
				}
			}
		}
	}
}

func pubDbGet(domain string) (group Group) {
	log.Println("pubDbGet")
	url := api + pubNames + domain
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err := json.Unmarshal(body, &group)
			if err == nil {
				return
			}
		}
	}
	return
}

func pubDbSet(group Group) (result bool) {
	log.Println("pubDbSet")
	domain := group.ScreenName
	b, err := json.Marshal(group)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", api+pubNames+domain, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		result = true
	}
	return
}

func pubSubTgAdd(group Group, msg *tgbotapi.Message) {

	gid := strconv.Itoa(group.Gid)
	url := api + pubSubTg + gid
	log.Println("pubSubTgAdd", url)
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("pubSubTgAdd body ", string(body))
		if err == nil {
			users := make(map[int]bool)
			//users := make([]string, 0, 100)
			json.Unmarshal(body, &users)
			delete(users, msg.From.ID)
			users[msg.From.ID] = true
			log.Println("pubSubTgAdd users ", users)
			data, err := json.Marshal(users)
			if err == nil {
				log.Println("pubSubTgAdd data ", string(data))
				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == 200 {
						bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Subscribed on vk domain:'"+group.ScreenName+"'"))
					}
				}
			}
		}
	}
	return
}

func vkWallUpd() {
	url := api + pubSubTg
	log.Println("vkWallUpd", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("pubSubTgAdd body ", string(body))
	}
}
