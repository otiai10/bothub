package model

type AccessToken struct {
	Token          string
	Secret         string
	AdditionalData map[string]string
}
