package swagger

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	config "github.com/peramic/Util.Auth/go/config"
	"github.com/peramic/utils"
)

// AdminUserIDDelete deletes a user
func AdminUserIDDelete(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	if config.IsAdmin(login) {
		vars := mux.Vars(r)
		id := vars["name"]

		err := config.LoginConfig.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	} else {
		http.Error(w, "Missing privileges", http.StatusForbidden)
	}
}

// AdminUserIDGet retrieves a user
func AdminUserIDGet(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	if config.IsAdmin(login) {
		vars := mux.Vars(r)
		name := vars["name"]

		l, err := config.LoginConfig.Get(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(User{Name: l.Name, Role: l.Role, HasDefaultPassword: l.HasDefaultPassword})
		}
	} else {
		http.Error(w, "Missing privileges", http.StatusForbidden)
	}
}

// AdminUserIDPut updates a user
func AdminUserIDPut(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	if config.IsAdmin(login) {
		vars := mux.Vars(r)
		name := vars["name"]

		var u *User
		err := utils.DecodeJSONBody(w, r, &u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			l, err := config.LoginConfig.Update(name, config.Login{Name: u.Name, Password: u.Password, Role: u.Role, HasDefaultPassword: u.HasDefaultPassword})
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(User{Name: l.Name, Role: l.Role, HasDefaultPassword: l.HasDefaultPassword})
			}
		}
	} else {
		http.Error(w, "Missing privileges", http.StatusForbidden)
	}
}

// AdminUsersPost adds a user
func AdminUsersPost(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	if config.IsAdmin(login) {
		var u *User
		err := utils.DecodeJSONBody(w, r, &u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			l, err := config.LoginConfig.Add(config.Login{Name: u.Name, Password: u.Password, Role: u.Role, HasDefaultPassword: u.HasDefaultPassword})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(User{Name: l.Name, Role: l.Role, HasDefaultPassword: l.HasDefaultPassword})
			}
		}
	} else {
		http.Error(w, "Missing privileges", http.StatusForbidden)
	}
}

// AdminUsersGet retrieves all users
func AdminUsersGet(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	if config.IsAdmin(login) {
		result := make(map[string]User)
		for _, login := range config.LoginConfig.Logins {
			result[login.Name] = User{Name: login.Name, Role: login.Role, HasDefaultPassword: login.HasDefaultPassword}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	} else {
		http.Error(w, "Missing privileges", http.StatusForbidden)
	}
}
