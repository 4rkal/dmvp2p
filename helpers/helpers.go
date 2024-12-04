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
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var (
	p2poolProcess *exec.Cmd
	xmrigProcess  *exec.Cmd
)

type Hashrate struct {
	Total   []float64 `json:"total"`
	Highest float64   `json:"highest"`
}

type SummaryResponse struct {
	Hashrate Hashrate `json:"hashrate"`
}

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
	p2poolProcess = exec.Command(path, "--host", host, "--wallet", address, "--mini")

	if err := p2poolProcess.Start(); err != nil {
		return fmt.Errorf("failed to start p2pool: %w", err)
	}

	fmt.Println("p2pool started successfully.")
	return nil
}

func startXmrig(path string) error {
	xmrigProcess = exec.Command(path,
		"-o", "127.0.0.1:3333",
		"--http-host=127.0.0.1",
		"--http-port=9999",
		"--http-access-token=dmvp2p",
		"--api-worker-id=1",
		"--api-id=1")

	if err := xmrigProcess.Start(); err != nil {
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

func StopMining() error {
	if p2poolProcess != nil && p2poolProcess.Process != nil {
		if err := p2poolProcess.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop p2pool: %w", err)
		}
		fmt.Println("p2pool stopped successfully.")
	}

	if xmrigProcess != nil && xmrigProcess.Process != nil {
		if err := xmrigProcess.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop xmrig: %w", err)
		}
		fmt.Println("xmrig stopped successfully.")
	}

	return nil
}

func SelectFileWithDialog(label *widget.Label, settingsPath *string) {
	selectedPath, err := dialog.File().Title("Select File").Load()
	if err != nil {
		if err.Error() != "canceled" {
			fmt.Println("Error selecting file:", err)
		}
		return
	}

	*settingsPath = selectedPath
	label.SetText("Selected file: " + *settingsPath)
}

func GetXmrigStats() Hashrate {
	url := "http://127.0.0.1:9999/1/summary"
	token := "dmvp2p"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return Hashrate{}
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return Hashrate{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status: %s\n", resp.Status)
		return Hashrate{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return Hashrate{}
	}

	var summary SummaryResponse
	err = json.Unmarshal(body, &summary)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return Hashrate{}
	}

	totalHashrate := 0.0
	if len(summary.Hashrate.Total) > 0 {
		totalHashrate = summary.Hashrate.Total[0]
	}

	fmt.Printf("Total Hashrate: %.2f H/s\n", totalHashrate)
	fmt.Printf("Highest Hashrate: %.2f H/s\n", summary.Hashrate.Highest)

	return Hashrate{
		Total:   []float64{totalHashrate},
		Highest: summary.Hashrate.Highest,
	}
}
