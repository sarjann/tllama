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

		if strings.HasPrefix(input, "/changemodel ") {
			fmt.Println("Changing active model")
			model := strings.TrimPrefix(input, "/changemodel ")
			model = strings.TrimSpace(model)

			if model == "" {
				fmt.Println("No model provided")
			}

			conf.Model = model
			fmt.Println("Changed model to: ", model)
			context = nil
			conf.Save()
			continue
		}

		if strings.HasPrefix(input, "/pullmodel") {
			fmt.Println("Pulling new model")
			model := strings.TrimPrefix(input, "/pullmodel")
			model = strings.TrimSpace(model)

			if model == "" {
				fmt.Println("No model provided")
			}
			err, _ := ollama.PullModel(model, conf)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("Changed model to: ", model)

			conf.Model = model
			context = nil
			conf.Save()
			continue
		}

		if strings.HasPrefix(input, "/deletemodel") {
			fmt.Println("Deleting model")
			model := strings.TrimPrefix(input, "/deletemodel")
			model = strings.TrimSpace(model)

			if model == "" {
				fmt.Println("No model provided")
			}
			err := ollama.DeleteModel(model, conf)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("Deleted model: ", model)
			continue
		}

		err, respStruct := ollama.PromptRequest(input, context, true, conf)

		context = append(context, respStruct.Context...)

		if err != nil {
			fmt.Println(err)
		}
	}
}
