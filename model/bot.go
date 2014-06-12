package model

import "bothub/infrastructure"
import "github.com/mrjones/oauth"
import "fmt"

type Bot struct {
	Master          Master
	ScreenName      string `json:"screen_name"`
	ProfileImageUrl string `json:"profile_image_url"`
	Token           oauth.AccessToken
}

func FindBotByName(botName string) (bot *Bot, e error) {
	if vaquero, e := infrastructure.GetVaquero(); e == nil {
		e = vaquero.Cast(botName, &bot)
	}
	if bot == nil {
		e = fmt.Errorf("Bot name `%s` not found", botName)
	}
	return
}

func FindBotByMasterName(masterName string) (bot *Bot, e error) {
	if vaquero, e := infrastructure.GetVaquero(); e == nil {
		if botName := vaquero.Get("bot." + masterName); botName != "" {
			return FindBotByName(botName)
		}
		e = fmt.Errorf("Bot related to master name `%s` not found", masterName)
	}
	return
}

func (b *Bot) Tweet(text string) (e error) {
	_, e = GetConsumer().Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": text,
		},
		&b.Token,
	)
	return
}
