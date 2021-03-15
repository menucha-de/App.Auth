package config

// Login model
type Login struct {
	Name string `json:"name,omitempty"`

	Password string `json:"password,omitempty"`

	Role string `json:"role,omitempty"`

	HasDefaultPassword bool `json:"hasDefaultPassword"`
}
