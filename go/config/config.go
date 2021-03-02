package config

import (
	"encoding/json"
	"os"

	loglib "github.com/peramic/logging"
)

const dirname = "./conf/auth"
const authFilename = dirname + "/auth.json"
const clientFilename = dirname + "/clients.json"
const loginFilename = dirname + "/logins.json"

var lg *loglib.Logger = loglib.GetLogger("config")

// AuthConfig contains the configuration for the auth server
var AuthConfig *AuthConfiguration

// ClientConfig contains the configured clients
var ClientConfig *ClientConfiguration

// LoginConfig contains the configured user logins
var LoginConfig *LoginConfiguration

func init() {

	authFile, err := os.Open(authFilename)
	if err == nil {
		dec := json.NewDecoder(authFile)
		err = dec.Decode(&AuthConfig)

		if err != nil {
			lg.Warning("Failed to parse config")
			AuthConfig = initAuth()
		}
	} else {
		AuthConfig = initAuth()
	}
	defer authFile.Close()

	clientFile, err := os.Open(clientFilename)
	if err == nil {
		dec := json.NewDecoder(clientFile)
		err = dec.Decode(&ClientConfig)

		if err != nil {
			lg.Warning("Failed to parse config")
			ClientConfig = initClients()
		}
	} else {
		ClientConfig = initClients()
	}
	defer clientFile.Close()

	loginFile, err := os.Open(loginFilename)
	if err == nil {
		dec := json.NewDecoder(loginFile)
		err = dec.Decode(&LoginConfig)

		if err != nil {
			lg.Warning("Failed to parse config")
			LoginConfig = initLogins()
		}
	} else {
		LoginConfig = initLogins()
	}
	defer loginFile.Close()
}

func initAuth() *AuthConfiguration {
	config := AuthConfiguration{AdministrationRole: "admin", IntrospectorClientID: "introspector", AccessTokenExpirationSeconds: (2 * 60 * 60), RefreshTokenExpirationSeconds: (24 * 60 * 60), TokenCookieKey: "TOKEN"}
	return &config
}

func initClients() *ClientConfiguration {
	clients := make(map[string]Client)
	clients["ui"] = Client{ID: "ui", Secret: "adryNYz2RNR8"}
	clients["introspector"] = Client{ID: "introspector", Secret: "y5Er1sAwY6zp"}
	config := ClientConfiguration{Clients: clients}
	return &config
}

func initLogins() *LoginConfiguration {
	logins := make(map[string]Login)
	logins["admin"] = Login{Name: "admin", Password: "admin", HasDefaultPassword: true, Role: "admin"}
	config := LoginConfiguration{Logins: logins}
	return &config
}
