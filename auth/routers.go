package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Route for HTTP server
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes for HTTP server
type Routes []Route

// NewRouter creates a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Server")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"AdminUserIdDelete",
		strings.ToUpper("Delete"),
		"/users/{name}",
		ValidateToken(AdminUserIDDelete),
	},

	Route{
		"AdminUserIdGet",
		strings.ToUpper("Get"),
		"/users/{name}",
		ValidateToken(AdminUserIDGet),
	},

	Route{
		"AdminUserIdPut",
		strings.ToUpper("Put"),
		"/users/{name}",
		ValidateToken(AdminUserIDPut),
	},

	Route{
		"AdminUsersPost",
		strings.ToUpper("Post"),
		"/users",
		ValidateToken(AdminUsersPost),
	},

	Route{
		"AdminUsersGet",
		strings.ToUpper("Get"),
		"/users",
		ValidateToken(AdminUsersGet),
	},

	Route{
		"Oauth2AuthorizePost",
		strings.ToUpper("Post"),
		"/oauth2/authorize",
		Oauth2AuthorizePost,
	},

	Route{
		"Oauth2IntrospectPost",
		strings.ToUpper("Post"),
		"/oauth2/introspect",
		Oauth2IntrospectPost,
	},

	Route{
		"Oauth2ValidatePost",
		strings.ToUpper("Post"),
		"/oauth2/validate",
		Oauth2ValidatePost,
	},

	Route{
		"Oauth2InvalidateTokenPost",
		strings.ToUpper("Post"),
		"/oauth2/invalidate_token",
		Oauth2InvalidateTokenPost,
	},

	Route{
		"Oauth2TokenPost",
		strings.ToUpper("Post"),
		"/oauth2/token",
		Oauth2TokenPost,
	},

	Route{
		"UserInfoGet",
		strings.ToUpper("Get"),
		"/user/info",
		ValidateToken(UserInfoGet),
	},

	Route{
		"UserPasswordDelete",
		strings.ToUpper("Delete"),
		"/user/password",
		ValidateToken(UserPasswordDelete),
	},

	Route{
		"UserPasswordPut",
		strings.ToUpper("Put"),
		"/user/password",
		ValidateToken(UserPasswordPut),
	},
}
