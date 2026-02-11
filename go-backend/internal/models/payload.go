package models

type MailPayload struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Code     string `json:"code"`
}
