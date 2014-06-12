package model

import "github.com/mrjones/oauth"
import "github.com/revel/revel"

// TODO: じっさいこのへんもinfrastructureじゃねえかと
var consumer oauth.Consumer
var consumer_initialized = false

func GetConsumer() *oauth.Consumer {
	if !consumer_initialized {
		key := revel.Config.StringDefault("twitter.consumer_key", "")
		secret := revel.Config.StringDefault("twitter.consumer_secret", "")
		consumer = *oauth.NewConsumer(
			key,
			secret,
			oauth.ServiceProvider{
				AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
				RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
				AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
			},
		)
	}
	return &consumer
}
