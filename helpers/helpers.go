package helpers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
)

type User struct {
	Name        string `json:"name"`
	GitHub      string `json:"github"`
	X           string `json:"x"`
	Website     string `json:"website"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type UsersData struct {
	Users []User `json:"users"`
}

func LoadUsers() []User {
	fileContent, err := ioutil.ReadFile("helpers/users.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var data UsersData
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return data.Users
}

func ParseURL(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Failed to parse URL", err)
		return nil
	}
	return u
}
