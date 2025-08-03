package entity

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponseError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}