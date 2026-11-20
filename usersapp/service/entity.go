package service

const (
	GoogleOAuthProvider = "google"
	GithubOAuthProvider = "github"
)

type User struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsActive    bool   `json:"is_active"`
}

type OAuthAccount struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Provider string `json:"provider"`
}
