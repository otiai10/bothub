package controllers

import "github.com/revel/revel"
import "bothub/infrastructure"
import "bothub/model"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	master, loginedAsMaster := c.checkMasterLogined()

	// bot登録状況のチェック
	vaquero, _ := infrastructure.GetVaquero()
	if botName := vaquero.Get("bot." + master.ScreenName); botName == "" {
        loginedAsMaster = false
	}

	return c.Render(loginedAsMaster, master)
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
