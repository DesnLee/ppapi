package model

type ResourceResponse[T any] struct {
	Resource T `json:"resource"`
}
type ResourcesResponse[T any] struct {
	Resources []T `json:"resources"`
	Pager     struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Count   int `json:"count"`
	} `json:"pager"`
}

type MsgResponse struct {
	Msg string `json:"msg" example:"错误消息"` // 错误响应消息
}
