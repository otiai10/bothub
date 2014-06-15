package controllers

import "github.com/revel/revel"
import "errors"
import "time"
import "bothub/model"

type Api struct {
	*revel.Controller
}

func (c Api) Index() revel.Result {
	return c.Render()
}

func (c Api) QueueAdd(user string, finish int64, text string) revel.Result {

	p := model.Payload{
		Master: user,
		Finish: finish,
		Text:   text,
	}
	if e := validate(p); e != nil {
		revel.ERROR.Println(e.Error())
		return c.RenderJson(model.Response{
			false, e.Error(), p,
		})
	}
	// TODO: このマスターに対するbotのチェック

	model.GetObserver().Enqueue(model.NewQueueFromPayload(p))

	revel.INFO.Printf("%+v\n", model.NewQueueFromPayload(p))

	return c.RenderJson(model.Response{
		true, "OK", p,
	})
}
func validate(p model.Payload) (e error) {
	if p.Master == "" {
		return errors.New("Missing parameter `user`")
	}
	if p.Finish < time.Now().Unix() {
		return errors.New("Invalid finish time in the past")
	}
	if p.Text == "" {
		return errors.New("Missing parameter `text`")
	}
	return
}
