package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const api = "http://badtobefat.ru/bolt"
const users = "/usertg/"

var bot *tgbotapi.BotAPI

func panic(e error) {
	if e != nil {
		log.Panic(e)
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
	panic(err)
	bot, err = tgbotapi.NewBotAPI(string(tlgrmtoken))
	if err != nil {
		log.Panic(err)
	}

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
				log.Printf("user created")
				sendRply(update.Message, "Привет, "+user.String()+"!\nПросто кинь мне ссылкуна паблик во вконтосе, который ты хочешь читать.")
			} else {
				sendRply(update.Message, "Something going wrong. Try later..")
			}
		default:
			sendRply(update.Message, "Пошел на хуй")
		}
	}
}

func newUser(user *tgbotapi.User) bool {
	b, err := json.Marshal(user)
	panic(err)
	req, err := http.NewRequest("PUT", api+users+strconv.Itoa(user.ID), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	panic(err)
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
