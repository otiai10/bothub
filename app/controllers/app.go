package controllers

import "github.com/revel/revel"
import "bothub/infrastructure"
import "bothub/model"
import "time"
import "fmt"
import "math/rand"

type App struct {
	*revel.Controller
}

func (c App) Index(message string) revel.Result {

	master, loginedAsMaster := c.checkMasterLogined()

	if !loginedAsMaster {
		return c.Render()
	}

	// bot登録状況のチェック
	vaquero, _ := infrastructure.GetVaquero()
	var bot model.Bot
	var botName string
	if botName = vaquero.Get("bot." + master.ScreenName); botName == "" {
		loginedAsMaster = false
		return c.Render(loginedAsMaster, master)
	}
	_ = vaquero.Cast(botName, &bot)

	rand.Seed(time.Now().Unix())
	message = c.Message("bot.default" + fmt.Sprintf("%02d", rand.Intn(5)))

	return c.Render(loginedAsMaster, master, bot, message)
}
func (c App) checkMasterLogined() (bot *model.Master, ok bool) {
	screenName, nameOK := c.Session["screen_name"]
	profileImageUrl, iconOK := c.Session["profile_image_url"]
	ok = (nameOK && iconOK)
	if !ok {
		return
	}
	bot = &model.Master{
		ScreenName:      screenName,
		ProfileImageUrl: profileImageUrl,
	}
	return
}
