package controllers

import "github.com/revel/revel"
import "bothub/model"
import "fmt"
import "encoding/json"
import "github.com/mrjones/oauth"
import "strings"

type Master struct {
	*revel.Controller
}

func (c Master) Login() revel.Result {

	// 一時的に使えるRequestTokenを取得する
	requestToken, url, err := model.GetConsumer().GetRequestTokenAndUrl(
		fmt.Sprintf("http://%s:%s/master/login/callback", "localhost", "19000"),
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

func (c Master) Callback(oauth_verifier string) revel.Result {

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
	var master model.Master
	_ = json.NewDecoder(resp.Body).Decode(&master)
	// }}}

	master.ProfileImageUrl = strings.Replace(master.ProfileImageUrl, "_normal.", ".", -1)
	// TODO: master.Token = *accessToken

	c.Session["screen_name"] = master.ScreenName
	c.Session["profile_image_url"] = master.ProfileImageUrl

	return c.Redirect(App.Index)
}
