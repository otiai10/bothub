package model

import "github.com/mrjones/oauth"
import "github.com/otiai10/rodeo"

type Bot struct {
	Master Master
	token  oauth.AccessToken
}

func FindBotByName(name string) (bot *Bot, e error) {
	if token, e := findBotKeysByName(name); e == nil {
		bot = &Bot{
			Master{name},
			oauth.AccessToken{
				token.Token,
				token.Secret,
				make(map[string]string),
			},
		}
	}
	return
}

func findBotKeysByName(name string) (token AccessToken, e error) {
	vaquero, e := rodeo.TheVaquero(rodeo.Conf{"localhost", "6379"})
	if e != nil {
		return
	}

	// うーんこのへん生々しい
	botName := vaquero.Get(name)
	e = vaquero.Cast("token."+botName, &token)
	return
}
func (b *Bot) Tweet(text string) (e error) {
	_, e = GetConsumer().Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": text,
		},
		&b.token,
	)
	return
}
