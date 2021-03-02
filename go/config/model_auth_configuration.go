package config

import (
	"encoding/json"
	"os"
)

// AuthConfiguration type
type AuthConfiguration struct {
	IntrospectorClientID          string `json:"introspectorClientId"`
	AdministrationRole            string `json:"administrationRole"`
	TokenCookieKey                string `json:"tokenCookieKey"`
	AccessTokenExpirationSeconds  int64  `json:"accessTokenExpirationSeconds"`
	RefreshTokenExpirationSeconds int64  `json:"refreshTokenExpirationSeconds"`
}

func (c *AuthConfiguration) serialize() {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}
	f, err := os.Create(authFilename)
	if err != nil {
		lg.WithError(err).Error("Failed to create or open configuration file")
	} else {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		enc.Encode(c)
	}
	defer f.Close()
}
