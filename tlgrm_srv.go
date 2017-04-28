package main

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("364483768:AAFhyU95D609MLVMQNNFzd3ZOAxIgyIHMN0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://badtobefat.ru/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	//go http.ListenAndServeTLS("https://badtobefat.ru", "/etc/letsencrypt/live/badtobefat.ru/cert.pem", "/etc/letsencrypt/live/badtobefat.ru/privkey.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
