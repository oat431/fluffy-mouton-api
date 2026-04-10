package response

type AuthResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}
