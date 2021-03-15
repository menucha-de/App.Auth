package config

import (
	"encoding/json"
	"os"
)

// ClientConfiguration type
type ClientConfiguration struct {
	Clients map[string]Client `json:"clients"`
}

func (c *ClientConfiguration) serialize() {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}
	f, err := os.Create(clientFilename)
	if err != nil {
		lg.WithError(err).Error("Failed to create or open configuration file")
	} else {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		enc.Encode(c)
	}
	defer f.Close()
}
