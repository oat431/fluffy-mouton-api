package response

type ShortLinkDTO struct {
	ShortLink    string `json:"short_link"`
	OriginalLink string `json:"original_link"`
	LinkType     string `json:"link_type"`
}
