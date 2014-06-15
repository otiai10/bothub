package controllers

import "github.com/revel/revel"

type Logout struct {
	*revel.Controller
}

func (c Logout) Index() revel.Result {
	delete(c.Session, "screen_name")
	delete(c.Session, "profile_image_url")
	return c.Redirect(App.Index)
}
