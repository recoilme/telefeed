package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	api     = "http://badtobefat.ru/bolt"
	users   = "/usertg/"
	someErr = "Something going wrong. Try later.."
	hello   = "Hello, %username!\nJust drop me link on vk public and i will send messages from it.\nExample: https://vk.com/myakotkapub"
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

	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			user := update.Message.From
			if newUser(user) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, hello))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, someErr))
			}
		default:
			//sendRply(update.Message, "Пошел на хуй")
			msg := update.Message.Text
			addPub(update.Message.Chat.ID, msg)
		}
	}
}

func newUser(user *tgbotapi.User) bool {
	b, err := json.Marshal(user)
	catch(err)
	req, err := http.NewRequest("PUT", api+users+strconv.Itoa(user.ID), bytes.NewBuffer(b))
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

func addPub(chatId int64, txt string) {
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
					log.Printf(parts[j])
					bot.Send(tgbotapi.NewMessage(chatId, "Found vk domain:'"+parts[j]+"'"))
					getPub(parts[j])
				}
			}
		}
	}
}

func getPub(name string) {
	Bar()
	url := "http://api.vk.com/method/wall.get?domain=" + name
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		//body, err := ioutil.ReadAll(resp.Body)
		catch(err)
		//log.Println("Body", string(body))
	}
}
