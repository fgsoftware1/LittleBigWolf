package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const tgToken string = "6626616536:AAFP21dn2GXQktjlwHki7z71-ax6i3AkUbo"

const (
	ConsumerKey       = "YOUR_CONSUMER_KEY"
	ConsumerSecret    = "YOUR_CONSUMER_SECRET"
	AccessToken       = "YOUR_ACCESS_TOKEN"
	AccessTokenSecret = "YOUR_ACCESS_TOKEN_SECRET"
)

func main() {
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	token := oauth1.NewToken(AccessToken, AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhook("https://192.168.1.83:8443/" + bot.Token)

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			command := update.Message.Command()
			switch command {
			case "report":
				args := update.Message.CommandArguments()
				if args != "" {
					reportUser(client, args, update.Message.Chat.ID, bot)
				} else {
					message := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide a Twitter username to report.")
					_, err := bot.Send(message)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func reportUser(client *twitter.Client, user string, chatID int64, bot *tgbotapi.BotAPI) {
	tweet, _, err := client.Statuses.Update(fmt.Sprintf("@%s This user violates Twitter policies.", user), nil)
	if err != nil {
		log.Printf("Error reporting user: %v", err)
		message := tgbotapi.NewMessage(chatID, "Error reporting user on Twitter.")
		_, err := bot.Send(message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	fmt.Printf("Reported user: @%s, Tweet ID: %d\n", user, tweet.ID)
	message := tgbotapi.NewMessage(chatID, fmt.Sprintf("User @%s reported on Twitter.", user))
	_, err = bot.Send(message)
	if err != nil {
		log.Println(err)
	}
}
