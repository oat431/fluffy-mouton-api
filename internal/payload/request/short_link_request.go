package request

type ShortLinkRequest struct {
	Url        string `json:"url" validate:"required,url"`
	CustomName string `json:"custom_name,omitempty" validate:"omitempty,alphanum"`
}
