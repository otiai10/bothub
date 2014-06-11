package controllers

import "github.com/revel/revel"
import "encoding/json"
import "errors"
import "time"
import "bothub/model"

type Api struct {
	*revel.Controller
}

func (c Api) Index() revel.Result {
	return c.RenderJson(map[string]string{
		"hoge": "fuga",
		"foo":  "bar",
	})
}

func (c Api) QueueAdd(queue string) revel.Result {
	var p model.Payload
	if e := json.Unmarshal([]byte(queue), &p); e != nil {
		return c.RenderJson(model.Response{
			false, e.Error(), p,
		})
	}
	if e := validate(p); e != nil {
		return c.RenderJson(model.Response{
			false, e.Error(), p,
		})
	}

	model.GetObserver().Enqueue(model.NewQueueFromPayload(p))

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
