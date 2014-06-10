package controllers

import "github.com/revel/revel"
import "encoding/json"
import "errors"

type Api struct {
	*revel.Controller
}

func (c Api) Index() revel.Result {
	return c.RenderJson(map[string]string{
		"hoge": "fuga",
		"foo":  "bar",
	})
}

// curl http://localhost:9000/api/queue/add -X POST -d 'sample={"hoge":"fuga"}'
type Sample struct {
	Hoge string `json:"hoge"`
}
type Response struct {
	OK       bool   `json:"ok"`
	Message  string `json:"message"`
	Accepted Sample `json:"accepted"`
}

func (c Api) QueueAdd(sample string) revel.Result {
	var s Sample
	if e := json.Unmarshal([]byte(sample), &s); e != nil {
		return c.RenderJson(Response{
			false, e.Error(), s,
		})
	}
	if e := validate(s); e != nil {
		return c.RenderJson(Response{
			false, e.Error(), s,
		})
	}
	return c.RenderJson(Response{
		true, "OK", s,
	})
}
func validate(s Sample) (e error) {
	if s.Hoge == "" {
		return errors.New("Missing parameter `hoge`")
	}
	return
}
