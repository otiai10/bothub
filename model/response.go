package model

type Response struct {
	OK       bool    `json:"ok"`
	Message  string  `json:"message"`
	Accepted Payload `json:"accepted"`
}
