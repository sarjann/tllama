package main

import (
	"bufio"
	"fmt"
	"main/tllama/ollama"
	"main/tllama/config"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	context := make([]int, 0, 100)
	config := config.LoadOrCreateConfig()
	fmt.Println(config)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		if input == "/clear" {
			context = nil
			fmt.Println("\033[91m", "New Chat (Cleared Context)", "\033[0m")
			continue
		}

		err, respStruct := ollama.PromptRequest(input, context, true, config)
		fmt.Println(respStruct)
		fmt.Println(respStruct.Response)
		for !respStruct.Done {
			err, respStruct = ollama.PromptRequest("", context, false, config)
			fmt.Println(respStruct.Response)
			fmt.Println(respStruct)
			// fmt.Println("\033[92m", respStruct.Response, "\033[0m")
		}
		fmt.Println("done")

		context = append(context, respStruct.Context...)
		fmt.Println("new context", context)
		if err != nil {
			fmt.Println(err)
		}
	}
}
