package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"

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
	url := "https://raw.githubusercontent.com/4rkal/dmvp2p/refs/heads/main/helpers/users.json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	fileContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
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

func startP2Pool(host, address, path string) error {
	cmd := exec.Command(path, "--host", host, "--wallet --mini", address)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start p2pool: %w", err)
	}

	fmt.Println("p2pool started successfully.")
	return nil
}

func startXmrig(path string) error {
	cmd := exec.Command(path, "-o", "127.0.0.1:3333")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start xmrig: %w", err)
	}

	fmt.Println("xmrig started successfully.")
	return nil
}

func StartMining(host, address, xmrigpath, p2poolpath string) error {
	err := startP2Pool(host, address, p2poolpath)
	if err != nil {
		return err
	}

	err2 := startXmrig(xmrigpath)
	if err2 != nil {
		return err2
	}
	return nil
}
