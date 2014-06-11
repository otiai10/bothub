package model

import "github.com/mrjones/oauth"

type Bot struct {
	Master Master
	token  oauth.AccessToken
}

func FindBotByName(name string) (bot *Bot, e error) {
	if keys, e := findBotKeysByName(name); e == nil {
		bot = &Bot{
			Master{name},
			oauth.AccessToken{
				keys.Token,
				keys.Secret,
				make(map[string]string),
			},
		}
	}
	return
}

type accessKeys struct {
	Token  string
	Secret string
}

func findBotKeysByName(name string) (k accessKeys, e error) {
	k = db[name]
	return
}
func (b *Bot) Tweet(text string) (e error) {
	_, e = consumer.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": text,
		},
		&b.token,
	)
	return
}
