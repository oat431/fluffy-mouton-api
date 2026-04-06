package common

type ResponseDTOError struct {
	HttpCode  int
	ErrorCode string
	Message   string
}
