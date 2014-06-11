package model

import "time"

import "github.com/revel/revel"

var instance_observer *Observer

type Observer struct {
	Enq   chan Queue
	Timed chan Queue
}

func GetObserver() *Observer {
	if instance_observer == nil {
		instance_observer = &Observer{
			make(chan Queue),
			make(chan Queue),
		}
	}
	return instance_observer
}
func (obs *Observer) Enqueue(q Queue) (e error) {
	obs.Enq <- q
	return
}
func (obs *Observer) WaitQueue() {
	for {
		q := <-obs.Enq
		obs.generateRoutine(q)
	}
}

// なんかObserverってmodelでいいのかな
func (obs *Observer) ListenQueue() {
	for {
		timed := <-obs.Timed
		bot, _ := FindBotByName("bot." + timed.Master.Name)
		e := bot.Tweet("@" + timed.Master.Name + " " + timed.Text)
		if e != nil {
			revel.WARN.Println("Tweet Failed!!!!!!!!", e)
		}
	}
}
func (obs *Observer) generateRoutine(q Queue) {
	go func() {
		dur := q.Time.Sub(time.Now())
		<-time.After(dur)
		obs.Timed <- q
	}()
}
