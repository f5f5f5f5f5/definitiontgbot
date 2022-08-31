package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("place a token here")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			word := update.Message.Text
			if strings.HasPrefix(word, "/") {
				switch word {
				case "/start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send a word which definition you`d like to get, you can also use /about, /help")
					bot.Send(msg)
				case "/about":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This bot uses Cambridge dictionary to return a definition of a given word")
					bot.Send(msg)
				case "/help":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use /about to get info about this bot or send a word which definition you need")
					bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, Definition(word))
				bot.Send(msg)
			}
		}
	}
}

func Definition(incoming string) string {
	url := fmt.Sprintf("https://dictionary.cambridge.org/dictionary/english/%s", incoming)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status code error, %s , %d", url, resp.StatusCode)
	}
	page, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	text := ""
	text = page.Find(".ddef_h").Find(".db").Text()
	endtext := strings.ReplaceAll(text, ":", ";")
	if endtext == "" {
		endtext = "incorrect word, try another one"
	}
	return endtext
}
