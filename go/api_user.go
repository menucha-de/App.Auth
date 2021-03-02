package swagger

import (
	"encoding/json"
	"net/http"

	config "github.com/peramic/Util.Auth/go/config"
	"github.com/peramic/utils"
)

// UserInfoGet get you current user information
func UserInfoGet(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	data := User{Name: login.Name, Role: login.Role, HasDefaultPassword: login.HasDefaultPassword}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// UserPasswordDelete reset your password
func UserPasswordDelete(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	err := config.LoginConfig.ResetPassword(login.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// UserPasswordPut set a new password
func UserPasswordPut(login config.Login, w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	var p *Password
	err := utils.DecodeJSONBody(w, r, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		err = config.LoginConfig.SetPassword(login.Name, p.NewPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
