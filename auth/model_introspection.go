package swagger

// Introspection model
type Introspection struct {
	Token string `json:"token"`

	TokenHint string `json:"token_hint,omitempty"`
}
