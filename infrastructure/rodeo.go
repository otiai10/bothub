package infrastructure

import "github.com/otiai10/rodeo"
import "github.com/revel/revel"

func GetVaquero() (*rodeo.Vaquero, error) {
	host := revel.Config.StringDefault("redis.host", "localhost")
	port := revel.Config.StringDefault("redis.port", "6379")
	return rodeo.TheVaquero(rodeo.Conf{host, port})
}
