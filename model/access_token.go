package model

import "github.com/revel/revel"
import "github.com/otiai10/rodeo"
import "github.com/mrjones/oauth"

type AccessToken struct {
	Token          string
	Secret         string
	AdditionalData map[string]string
}

func StoreAccessTokenByName(name string, token *oauth.AccessToken) (e error) {

	// {{{ TODO: DRY  インフラにもっていく？
	host := revel.Config.StringDefault("redis.host", "localhost")
	port := revel.Config.StringDefault("redis.port", "6379")
	vaquero, _ := rodeo.TheVaquero(rodeo.Conf{host, port})
	// }}}

	return vaquero.Store("token."+name, token)
}

func FindAccessTokenFromName(name string) (token oauth.AccessToken, e error) {

	// {{{ TODO: DRY  インフラにもっていく？
	host := revel.Config.StringDefault("redis.host", "localhost")
	port := revel.Config.StringDefault("redis.port", "6379")
	vaquero, _ := rodeo.TheVaquero(rodeo.Conf{host, port})
	// }}}

	e = vaquero.Cast("token."+name, &token)
	return
}
