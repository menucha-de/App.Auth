package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoginConfiguration type
type LoginConfiguration struct {
	Logins map[string]Login `json:"logins"`
}

// ResetPassword reset the password to default
func (c *LoginConfiguration) ResetPassword(login string) error {
	if l, ok := c.Logins[login]; ok {
		l.Password = l.Name
		l.HasDefaultPassword = true
		c.Logins[login] = l
		c.serialize()
	} else {
		return fmt.Errorf("Failed to reset password for login: %v", login)
	}
	return nil
}

// SetPassword set the password
func (c *LoginConfiguration) SetPassword(login string, password string) error {
	if l, ok := c.Logins[login]; ok {
		l.Password = password
		l.HasDefaultPassword = false
		c.Logins[login] = l
		c.serialize()
	} else {
		return fmt.Errorf("Failed to set password for login: %v", login)
	}
	return nil
}

// Get a specific user
func (c *LoginConfiguration) Get(login string) (Login, error) {
	if l, ok := c.Logins[login]; ok {
		return l, nil
	}
	return Login{}, fmt.Errorf("Failed to find user: %v", login)
}

// Delete a user
func (c *LoginConfiguration) Delete(login string) error {
	if l, ok := c.Logins[login]; ok {
		delete(c.Logins, l.Name)
		c.serialize()
	} else {
		return fmt.Errorf("Failed to delete user: %v", login)
	}
	return nil
}

// Update a user
func (c *LoginConfiguration) Update(name string, login Login) (Login, error) {
	login.Name = strings.TrimSpace(login.Name)
	if len(login.Name) == 0 || len(login.Password) == 0 {
		return Login{}, fmt.Errorf("Invalid user: %v", login)
	}
	if _, ok := c.Logins[name]; ok {
		if name != login.Name {
			delete(c.Logins, name)
		}
		c.Logins[login.Name] = login
		c.serialize()
		return login, nil
	}
	return Login{}, fmt.Errorf("Failed to update user: %v", login)
}

// Add a user
func (c *LoginConfiguration) Add(login Login) (Login, error) {
	login.Name = strings.TrimSpace(login.Name)
	if len(login.Name) == 0 || len(login.Password) == 0 {
		return Login{}, fmt.Errorf("Invalid user: %v", login)
	}
	if _, ok := c.Logins[login.Name]; !ok {
		c.Logins[login.Name] = login
		c.serialize()
		return login, nil
	}
	return Login{}, fmt.Errorf("User already exists: %v", login)
}

// IsAdmin whether the login is an administrator
func IsAdmin(login Login) bool {
	return login.Role == AuthConfig.AdministrationRole
}

func (c *LoginConfiguration) serialize() {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}
	f, err := os.Create(loginFilename)
	if err != nil {
		lg.WithError(err).Error("Failed to create or open configuration file")
	} else {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		enc.Encode(c)
	}
	defer f.Close()
}
