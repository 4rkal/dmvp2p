package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/thedevsaddam/gojsonq"
)

func findPerson(name string) (string, string, string) {
	jq := gojsonq.New().File("./people.json")
	result := jq.From("people").Where("name", "=", name).First()

	if result == nil {
		return "", "", ""
	}

	person, ok := result.(map[string]interface{})
	if !ok {
		return "", "", ""
	}

	nameFromJSON, nameOK := person["name"].(string)
	addressFromJSON, addressOK := person["address"].(string)
	siteFromJSON, siteOK := person["website"].(string)

	if !nameOK || !addressOK || !siteOK {
		return "", "", ""
	}

	return nameFromJSON, addressFromJSON, siteFromJSON
}

func startMining(path string) (error, *exec.Cmd) {
	cmd := exec.Command(path, "-o 127.0.0.1:3333")
	if err := cmd.Start(); err != nil {
		return err, nil
	}
	return nil, cmd
}

func startP2pool(address string, settings Settings) (error, *exec.Cmd) {
	cmd := exec.Command(settings.P2poolPath,
		"--wallet", address,
		"--host", settings.Hostname,
		"--rpc-port", fmt.Sprintf("%d", settings.RPC),
		"--zmq-port", fmt.Sprintf("%d", settings.ZMQ),
	)
	if err := cmd.Start(); err != nil {
		return err, nil
	}
	return nil, cmd
}

func kill(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return err
	}
	return nil
}

func loadPeople() (People, error) {
	var people People

	response, err := http.Get("https://raw.githubusercontent.com/4rkal/dmvp2p/main/people.json")
	if err != nil {
		return people, fmt.Errorf("failed to fetch data from URL: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return people, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	byteValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return people, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(byteValue, &people)
	if err != nil {
		return people, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return people, nil
}
