package model

type DataResponse[T any] struct {
	Data T `json:"data"`
}
type MsgResponse struct {
	Msg string `json:"msg"`
}
