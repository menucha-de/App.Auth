package swagger

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	loglib "github.com/peramic/logging"

	"gopkg.in/oauth2.v3"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	config "github.com/peramic/Util.Auth/go/config"
)

// AuthenticatedHandlerFunc handler func with login
type AuthenticatedHandlerFunc func(config.Login, http.ResponseWriter, *http.Request)

// AuthServer the OAuth2 server instance
var AuthServer *server.Server

var lg *loglib.Logger = loglib.GetLogger("auth")

// Init the OAuth2 server
func Init() {

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	passwordTokenConfig := &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Duration(config.AuthConfig.AccessTokenExpirationSeconds) * time.Second,
		// refresh token expiration time
		RefreshTokenExp:   time.Duration(config.AuthConfig.RefreshTokenExpirationSeconds) * time.Second,
		IsGenerateRefresh: (config.AuthConfig.RefreshTokenExpirationSeconds > 0),
	}
	manager.SetPasswordTokenCfg(passwordTokenConfig)

	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// register clients from config
	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)
	for _, client := range config.ClientConfig.Clients {
		err := clientStore.Set(client.ID, &models.Client{
			ID:     client.ID,
			Secret: client.Secret,
		})
		if err != nil {
			lg.WithError(err).Error("Failed to register client")
		}
	}

	AuthServer = server.NewDefaultServer(manager)
	AuthServer.SetAllowGetAccessRequest(false)
	AuthServer.SetAllowedGrantType(oauth2.PasswordCredentials, oauth2.Refreshing /*, oauth2.AuthorizationCode*/)
	AuthServer.SetAllowedResponseType(oauth2.Code, oauth2.Token)

	AuthServer.SetClientInfoHandler(clientCredentialHandler)
	AuthServer.SetPasswordAuthorizationHandler(loginHandler)
	AuthServer.SetClientScopeHandler(scopeHandler)
	AuthServer.SetClientAuthorizedHandler(clientAuthorizedHandler)

	AuthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		lg.WithError(err).Error("Internal error")
		return
	})
	AuthServer.SetResponseErrorHandler(func(re *errors.Response) {
		lg.WithError(re.Error).Error("Response error")
		return
	})
}

// Shutdown the OAuth2 server
func Shutdown() {

}

func clientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	// make sure, that a request with an unknown client ID returns 401 instead of 500
	_, err = AuthServer.Manager.GetClient(clientID)
	if err != nil {
		return false, errors.ErrInvalidClient // maps to 401
	}
	return true, nil
}

func loginHandler(username, password string) (userID string, err error) {
	// TODO: use system config?
	if l, ok := config.LoginConfig.Logins[username]; ok {
		if l.Password == password {
			return l.Name, nil
		}
	}
	return "", errors.ErrAccessDenied
}

func scopeHandler(clientID, scope string) (allowed bool, err error) {
	// TODO: add scope checking here (later)
	return true, nil
}

func clientCredentialHandler(r *http.Request) (clientID, clientSecret string, err error) {
	// first check basic auth, then check form values from body
	id, secret, err := server.ClientBasicHandler(r)
	if err != nil {
		// now check form
		return server.ClientFormHandler(r)
	}
	return id, secret, nil
}

func writeError(w http.ResponseWriter, err error) int {
	_, statusCode, header := AuthServer.GetErrorData(err)
	for key := range header {
		w.Header().Set(key, header.Get(key))
	}
	if statusCode > 0 {
		return statusCode
	}
	return http.StatusOK
}

func tokenError(w http.ResponseWriter, err error) error {
	// TODO: copied from oauth2.v3/server/server.go
	data, statusCode, header := AuthServer.GetErrorData(err)
	return token(w, data, header, statusCode)
}

func token(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
	// TODO: copied from oauth2.v3/server/server.go
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	for key := range header {
		w.Header().Set(key, header.Get(key))
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func handleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	// TODO: copied from oauth2.v3/server/server.go
	gt, tgr, err := AuthServer.ValidationTokenRequest(r)
	if err != nil {
		return tokenError(w, err)
	}

	ti, err := AuthServer.GetAccessToken(gt, tgr)
	if err != nil {
		return tokenError(w, err)
	}

	header := make(http.Header)
	requestCookie := r.FormValue("request_cookie")
	if strings.ToLower(requestCookie) == "true" {
		setTokenCookie(ti, &header)
	}

	return token(w, AuthServer.GetTokenData(ti), header)
}

func setTokenCookie(ti oauth2.TokenInfo, header *http.Header) {
	// adding a cookie here
	expire := time.Now().Add(ti.GetAccessExpiresIn())
	cookie := http.Cookie{
		Name:     config.AuthConfig.TokenCookieKey,
		Value:    ti.GetAccess(),
		Expires:  expire,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		// TODO: should set secure and HttpOnly
	}
	header.Set("Set-Cookie", cookie.String())
}

func getTokenFromCookie(r *http.Request) string {
	cookie, _ := r.Cookie(config.AuthConfig.TokenCookieKey)
	if cookie != nil && len(strings.TrimSpace(cookie.Value)) > 0 {
		return cookie.Value
	}
	return ""
}

// SetDefaultHeaders sets the default headers for all responses
func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}

// Oauth2TokenPost obtains a new token
func Oauth2TokenPost(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	handleTokenRequest(w, r)
}

// Oauth2AuthorizePost obtains a new token
func Oauth2AuthorizePost(w http.ResponseWriter, r *http.Request) {
	// TODO: test and finish (several handlers missing)
	// actual login needs to be implemented as a separate call
	SetDefaultHeaders(w)
	err := AuthServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// ValidateToken internally validates a token
func ValidateToken(f AuthenticatedHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token oauth2.TokenInfo
		t := getTokenFromCookie(r)
		if t != "" {
			token, _ = AuthServer.Manager.LoadAccessToken(t)
		}
		if token == nil {
			token, _ = AuthServer.ValidationBearerToken(r)
		}
		if token == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user := config.LoginConfig.Logins[token.GetUserID()]
		f(user, w, r)
	})
}

// Oauth2IntrospectPost token introspection following standards (always returns HTTP 200)
func Oauth2IntrospectPost(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	status := http.StatusOK
	data := IntrospectionResult{}

	clientID, clientSecret, err := AuthServer.ClientInfoHandler(r)
	if err != nil {
		status = writeError(w, err)
	} else {
		cli, err := AuthServer.Manager.GetClient(clientID)
		if err != nil {
			status = writeError(w, errors.ErrInvalidClient)
		} else if clientID != config.AuthConfig.IntrospectorClientID {
			status = writeError(w, errors.ErrInvalidClient)
		} else if clientSecret != cli.GetSecret() {
			status = writeError(w, errors.ErrInvalidClient)
		} else {
			token := r.FormValue("token")
			tokenHint := r.FormValue("token_hint")
			var loadToken func(string) (oauth2.TokenInfo, error) = nil
			if tokenHint == "access_token" {
				loadToken = AuthServer.Manager.LoadAccessToken
			} else if tokenHint == "refresh_token" {
				loadToken = AuthServer.Manager.LoadRefreshToken
			} else {
				status = writeError(w, errors.ErrInvalidRequest)
			}
			if loadToken != nil {
				_, err := loadToken(token)
				data = IntrospectionResult{Active: err == nil}
			}
		}
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Oauth2ValidatePost token introspection or validation which returns 200 or 401
func Oauth2ValidatePost(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	status := http.StatusNoContent

	clientID, clientSecret, err := AuthServer.ClientInfoHandler(r)
	if err != nil {
		status = writeError(w, err)
	} else {
		cli, err := AuthServer.Manager.GetClient(clientID)
		if err != nil {
			status = writeError(w, errors.ErrInvalidClient)
		} else if clientID != config.AuthConfig.IntrospectorClientID {
			status = writeError(w, errors.ErrInvalidClient)
		} else if clientSecret != cli.GetSecret() {
			status = writeError(w, errors.ErrInvalidClient)
		} else {

			setCookie := false
			var tokenHint string
			token := getTokenFromCookie(r)
			if token != "" {
				setCookie = true
				tokenHint = "access_token"
			} else {
				token = r.FormValue("token")
				tokenHint = r.FormValue("token_hint")
			}

			var loadToken func(string) (oauth2.TokenInfo, error) = nil
			if tokenHint == "access_token" {
				loadToken = AuthServer.Manager.LoadAccessToken
			} else if tokenHint == "refresh_token" {
				loadToken = AuthServer.Manager.LoadRefreshToken
			} else {
				status = writeError(w, errors.ErrInvalidRequest)
			}
			if loadToken != nil {
				ti, err := loadToken(token)
				if err != nil {
					status = http.StatusUnauthorized
				}
				if ti != nil && setCookie {
					header := make(http.Header)
					setTokenCookie(ti, &header)
					for key := range header {
						w.Header().Set(key, header.Get(key))
					}
				}
			}
		}
	}

	w.WriteHeader(status)
}

// Oauth2InvalidateTokenPost invalidates a token
func Oauth2InvalidateTokenPost(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	status := http.StatusOK
	data := Token{}

	clientID, clientSecret, err := AuthServer.ClientInfoHandler(r)
	if err != nil {
		status = writeError(w, err)
	} else {
		cli, err := AuthServer.Manager.GetClient(clientID)
		if err != nil {
			status = writeError(w, errors.ErrInvalidClient)
		} else if clientSecret != cli.GetSecret() {
			status = writeError(w, errors.ErrInvalidClient)
		} else {
			token, err := AuthServer.ValidationBearerToken(r)
			if err != nil {
				status = writeError(w, errors.ErrInvalidRequest)
			} else {
				err = AuthServer.Manager.RemoveAccessToken(token.GetAccess())
				if err != nil {
					status = writeError(w, err)
				} else {
					data = Token{AccessToken: token.GetAccess()}
				}
			}
		}
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
