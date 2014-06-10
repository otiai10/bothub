package controllers

import "github.com/revel/revel"
import "encoding/json"

type Api struct {
	*revel.Controller
}

func (c Api) Index() revel.Result {
	return c.RenderJson(map[string]string{
        "hoge": "fuga",
        "foo" : "bar",
    })
}

// curl http://localhost:9000/api/queue/add -X POST -d 'sample={"hoge":"fuga"}'
type Sample struct {
    Hoge string `json:"hoge"`
}
func (c Api) QueueAdd(sample string) revel.Result {
    var s Sample
    json.Unmarshal([]byte(sample), &s)
	return c.RenderJson(s)
}
