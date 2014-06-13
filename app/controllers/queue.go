package controllers

import "github.com/revel/revel"
import "bothub/model"
import "bothub/infrastructure"
import "time"

type Queue struct {
	*revel.Controller
}

func (c Queue) Add(finish string, text string) revel.Result {

	// ログイン状態のチェック
	master, ok := c.checkMasterLogined()
	if !ok {
		return c.Redirect(App.Index)
	}
	// bot登録状況のチェック
	vaquero, _ := infrastructure.GetVaquero()
	if botName := vaquero.Get("bot." + master.ScreenName); botName == "" {
		return c.Redirect(App.Index)
	}

	// パラメータの妥当性チェック
	t, e := time.Parse("2006-01-02T15:04(MST)", finish+"(JST)")
	if e != nil {
		return c.Redirect(App.Index)
	}
	if t.Before(time.Now()) {
		return c.Redirect(App.Index)
	}

	queue := model.NewQueue(master.ScreenName, t, text)
	model.GetObserver().Enqueue(queue)

	return c.Redirect(App.Index)
}
func (c Queue) checkMasterLogined() (master *model.Master, ok bool) {
	screenName, nameOK := c.Session["screen_name"]
	profileImageUrl, iconOK := c.Session["profile_image_url"]
	ok = (nameOK && iconOK)
	if !ok {
		return
	}
	master = &model.Master{
		ScreenName:      screenName,
		ProfileImageUrl: profileImageUrl,
	}
	return
}
