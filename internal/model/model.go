package model

type Pager struct {
    Page    int64 `json:"page" example:"1"`
    PerPage int64 `json:"per_page" example:"10"`
    Count   int64 `json:"count" example:"100"`
}

type ResourceResponse[T any] struct {
    Resource T `json:"resource"`
}
type ResourcesResponse[T any] struct {
    Resources []T   `json:"resources"`
    Pager     Pager `json:"pager"`
}

type MsgResponse struct {
    Msg string `json:"msg" example:"错误消息"` // 错误响应消息
}
