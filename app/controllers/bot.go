package controllers

import "github.com/revel/revel"
import "bothub/model"
import "bothub/infrastructure"
import "fmt"
import "encoding/json"
import "github.com/mrjones/oauth"
import "strings"

type Bot struct {
	*revel.Controller
}

func (c Bot) Index() revel.Result {

	bot, ok := checkBotLogined(c)
	if !ok {
		return c.Redirect(App.Index)
	}

	var storedBot model.Bot
	var masterAlreadyRegistered bool
	vaquero, _ := infrastructure.GetVaquero()
	vaquero.Cast(bot.ScreenName, &storedBot)
	if storedBot.Master.ScreenName != "" && storedBot.Master.ProfileImageUrl != "" {
		bot.Master = storedBot.Master
		masterAlreadyRegistered = true
	}

	return c.Render(bot, masterAlreadyRegistered)
}

func (c Bot) Login() revel.Result {

	// 一時的に使えるRequestTokenを取得する
	requestToken, url, err := model.GetConsumer().GetRequestTokenAndUrl(
		fmt.Sprintf("http://%s:%s/bot/login/callback", "localhost", "19000"),
	)

	if err != nil {
		revel.ERROR.Println("リクエストトークンの取得に失敗", err)
		return c.Redirect(App.Index)
	}

	// 一時的に使えるrequestTokenが取得できたので、セッションに保存しておく
	c.Session["requestToken"] = requestToken.Token
	c.Session["requestSecret"] = requestToken.Secret

	// あとはTwitterのアプリ認証画面でログインしてoauth_verifierを取ってきてもらう
	return c.Redirect(url)
}

func (c Bot) Callback(oauth_verifier string) revel.Result {

	// 誰だテメーは、どっから来た
	if oauth_verifier == "" {
		return c.Redirect(App.Index)
	}

	// RequestTokenを回復
	requestToken := &oauth.RequestToken{
		c.Session["requestToken"],
		c.Session["requestSecret"],
	}
	// 不要になったので捨てる
	delete(c.Session, "requestToken")
	delete(c.Session, "requestSecret")

	// これと、oauth_verifierを用いてaccess_tokenを獲得する
	accessToken, err := model.GetConsumer().AuthorizeToken(requestToken, oauth_verifier)
	if err != nil {
		revel.ERROR.Println("アクセストークンの取得に失敗", err)
		return c.Redirect(App.Index)
	}

	// 成功したので、これを用いてユーザ情報を取得する
	// {{{ TODO: DRY インフラ
	resp, _ := model.GetConsumer().Get(
		//"https://api.twitter.com/1.1/statuses/mentions_timeline.json",
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		accessToken,
	)
	defer resp.Body.Close()
	var bot model.Bot
	_ = json.NewDecoder(resp.Body).Decode(&bot)
	// }}}

	bot.ProfileImageUrl = strings.Replace(bot.ProfileImageUrl, "_normal.", ".", -1)
	bot.Token = *accessToken

	c.Session["screen_name"] = bot.ScreenName
	c.Session["profile_image_url"] = bot.ProfileImageUrl

	vaquero, _ := infrastructure.GetVaquero()
	// すでにmasterが要る場合などはそれを適用する
	var storedBot model.Bot
	vaquero.Cast(bot.ScreenName, &storedBot)
	if storedBot.Master.ScreenName != "" && storedBot.Master.ProfileImageUrl != "" {
		bot.Master = storedBot.Master
	}
	vaquero.Store(bot.ScreenName, bot)

	return c.Redirect(Bot.Index)
}

func (c Bot) Confirm(master_name string) revel.Result {

	bot, ok := checkBotLogined(c)
	if !ok {
		return c.Redirect(App.Index)
	}

	bot, _ = model.FindBotByName(bot.ScreenName)

	// {{{ TODO: DRY インフラ
	resp, _ := model.GetConsumer().Get(
		"https://api.twitter.com/1.1/users/lookup.json",
		map[string]string{
			"screen_name": strings.Replace(master_name, "@", "", 1),
		},
		&bot.Token,
	)
	defer resp.Body.Close()
	accounts := []model.Master{}
	_ = json.NewDecoder(resp.Body).Decode(&accounts)
	var master model.Master
	if len(accounts) < 1 {
		return c.Redirect(Bot.Index)
	}
	master = accounts[0]
	master.ProfileImageUrl = strings.Replace(master.ProfileImageUrl, "_normal.", ".", -1)
	// }}}

	return c.Render(bot, master)
}

func (c Bot) Update(master_confirmed string) revel.Result {

	bot, ok := checkBotLogined(c)
	if !ok {
		c.Redirect(App.Index)
	}

	vaquero, _ := infrastructure.GetVaquero()
	vaquero.Cast(bot.ScreenName, bot)

	// {{{ TODO: DRY インフラ
	resp, _ := model.GetConsumer().Get(
		"https://api.twitter.com/1.1/users/lookup.json",
		map[string]string{
			"screen_name": strings.Replace(master_confirmed, "@", "", 1),
		},
		&bot.Token,
	)
	defer resp.Body.Close()
	accounts := []model.Master{}
	_ = json.NewDecoder(resp.Body).Decode(&accounts)
	var master model.Master
	if len(accounts) < 1 {
		return c.Redirect(Bot.Index)
	}
	master = accounts[0]
	master.ProfileImageUrl = strings.Replace(master.ProfileImageUrl, "_normal.", ".", -1)
	// }}}

	// Store with master info
	bot.Master = master
	vaquero.Store(bot.ScreenName, bot)
	vaquero.Set("bot."+master.ScreenName, bot.ScreenName)

	return c.Redirect(Bot.Index)
}

func checkBotLogined(c Bot) (bot *model.Bot, ok bool) {
	screenName, nameOK := c.Session["screen_name"]
	profileImageUrl, iconOK := c.Session["profile_image_url"]
	ok = (nameOK && iconOK)
	if !ok {
		return
	}
	bot = &model.Bot{
		ScreenName:      screenName,
		ProfileImageUrl: profileImageUrl,
	}
	return
}
