package swagger

// Token model
type Token struct {
	TokenType string `json:"token_type,omitempty"`

	AccessToken string `json:"access_token"`

	ExpiresIn int32 `json:"expires_in,omitempty"`

	RefreshToken string `json:"refresh_token,omitempty"`

	Scope string `json:"scope,omitempty"`
}
