package main

import (
	"bufio"
	"fmt"
	"github.com/sarjann/tllama/config"
	"github.com/sarjann/tllama/ollama"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	context := make([]int, 0, 100)
	conf := config.LoadOrCreateConfig()
	fmt.Println("Running on server: ", conf.OllamaUrl, "With model: ", conf.Model)

	fmt.Println("\033[91m", "Started", "\033[0m")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		switch input {
		case "/clear":
			context = nil
			fmt.Println("\033[91m", "New Chat (Cleared Context)", "\033[0m")
			continue
		case "/clearconfig":
			fmt.Println("Removing config, please restart and enter new config")
			config.ClearConfig()
			return
		}

		err, respStruct := ollama.PromptRequest(input, context, true, conf)

		context = append(context, respStruct.Context...)

		if err != nil {
			fmt.Println(err)
		}
	}
}
