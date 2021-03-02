package swagger

// Auth model
type Auth struct {
	Username string `json:"username"`

	Password string `json:"password"`

	RefreshToken string `json:"refresh_token"`

	GrantType string `json:"grant_type"`
}
