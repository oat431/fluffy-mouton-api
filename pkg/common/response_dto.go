package common

type ResponseDTO[T any] struct {
	Data   *T                `json:"data"`
	Status ResponseDTOStatus `json:"status"`
	Error  *ResponseDTOError `json:"error"`
}
