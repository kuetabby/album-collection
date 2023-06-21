package shared

type BaseResponse struct {
	Status string `json:"status"`
	Message *string `json:"message"`
	Error bool `json:"error"`
}

type SingleResponse[T any] struct {
	BaseResponse
	Content *T `json:"content"`
}