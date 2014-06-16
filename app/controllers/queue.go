package controllers

import "github.com/revel/revel"
import "bothub/model"
import "bothub/infrastructure"
import "time"
import "math/rand"
import "fmt"

type Queue struct {
	*revel.Controller
}

func (c Queue) GetAcceptedMessage(queue model.Queue) string {
	rand.Seed(time.Now().Unix())
	return c.Message("bot.accepted"+fmt.Sprintf("%02d", rand.Intn(5)), queue.Time.Format("02日15時04分"))
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
		if t, e = time.Parse("2006-01-02T15:04:05(MST)", finish+"(JST)"); e != nil {
			c.Flash.Out["invalid"] = c.Message("bot.invalid.invalid_time")
			return c.Redirect(App.Index)
		}
	}
	if t.Before(time.Now()) {
		c.Flash.Out["invalid"] = c.Message("bot.invalid.time_past")
		return c.Redirect(App.Index)
	}
	if text == "" {
		c.Flash.Out["invalid"] = c.Message("bot.invalid.text_nil")
		return c.Redirect(App.Index)
	}

	queue := model.NewQueue(master.ScreenName, t, text)
	model.GetObserver().Enqueue(queue)

	c.Flash.Out["accepted"] = c.GetAcceptedMessage(queue)

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
