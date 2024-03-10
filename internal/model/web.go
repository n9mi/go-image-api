package model

type ErrorResponse struct {
	Code     int      `json:"code"`
	Messages []string `json:"messages"`
}
