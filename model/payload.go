package model

type Payload struct {
	Master string `json:"user"`
	Finish int64  `json:"finish"`
	Text   string `json:"text"`
}
