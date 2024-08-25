package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	OllamaUrl string `json:"ollama_url"`
	Model     string `json:"model"`
}

var usr, _ = user.Current()
var dir = usr.HomeDir
var ConfigPath = dir + "/.config/ttlama.conf"

func validateConfig(config Config) {
	resp, err := http.Get(config.OllamaUrl)
	if err != nil {
		fmt.Sprintln("Failed to reach ollama url %s", config.OllamaUrl)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Sprintln("None valid response from ollama url, expect 200 %s", config.OllamaUrl)
		panic(err)
	}

	if config.Model == "" {
		panic("Default model cannot be empty")
	}
}

func LoadOrCreateConfig() Config {
	var config Config
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		fmt.Sprintln("File not found at %s, creating new config", ConfigPath)
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Please enter the ollama url path (eg. http://localhost:11434)")
		ollamaUrl, _ := reader.ReadString('\n')
		ollamaUrl = strings.Replace(ollamaUrl, "\n", "", -1)

		fmt.Println("Please enter the default model (e.g. llama3.1), currently this should be preinstalled")
		defaultModel, _ := reader.ReadString('\n')
		defaultModel = strings.Replace(defaultModel, "\n", "", -1)

		config = Config{ollamaUrl, defaultModel}
		validateConfig(config)

		fmt.Println("Writing to config")
		configJSONBody, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}

		file, err := os.Create(ConfigPath)
		fmt.Println(ConfigPath)
		if err != nil {
			panic(err)
		}
		_, err = file.Write(configJSONBody)
		if err != nil {
			panic(err)
		}
	} else {
		err = json.Unmarshal(data, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

func ClearConfig() {
	err := os.Remove(ConfigPath)
	if err != nil {
		fmt.Println("No config to remove")
	}
}
