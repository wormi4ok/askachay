package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Settings struct {
	Users []User `json:"users"`
}

type User struct {
	Username   string `json:"username"`
	UploadPath string `json:"upload_path"`
}

func FromFile(filePath string) (*Settings, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	c := &Settings{}
	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("can't decode config JSON: %w", err)
	}

	return c, nil
}
