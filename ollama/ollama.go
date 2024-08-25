package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sarjann/tllama/config"
	"net/http"
	"net/url"
)

type RespStruct struct {
	Model                string `json:"model"`
	Created_at           string `json:"created_at"`
	Response             string `json:"response"`
	Done                 bool   `json:"done"`
	Done_reason          string `json:"done_reason,omitempty"`
	Context              []int  `json:"context,omitempty"`
	Total_duration       int    `json:"total_duration,omitempty"`
	Load_duration        int    `json:"load_duration,omitempty"`
	Prompt_eval_count    int    `json:"prompt_eval_count,omitempty"`
	Prompt_eval_duration int    `json:"prompt_eval_duration,omitempty"`
	Eval_count           int    `json:"eval_count,omitempty"`
	Eval_duration        int    `json:"eval_duration,omitempty"`
}

type ReqStruct struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context"`
}

func PromptRequest(prompt string, context []int, firstRun bool, conf config.Config) (error, RespStruct) {
	ollamaUrl, err := url.Parse(conf.OllamaUrl)
	if err != nil {
		panic(err)
	}

	apiUrl := ollamaUrl.JoinPath("api")

	generateUrl := apiUrl.JoinPath("generate")

	var promptBody = ReqStruct{
		Model:   conf.Model,
		Prompt:  prompt,
		Stream:  true,
		Context: context,
	}

	promptJSONBody, err := json.Marshal(promptBody)
	if err != nil {
		return err, RespStruct{}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", generateUrl.String(), bytes.NewReader(promptJSONBody))
	if err != nil {
		return err, RespStruct{}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err, RespStruct{}
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var respStruct RespStruct
		err := decoder.Decode(&respStruct)
		if err != nil {
			break
		}
		fmt.Print("\033[92m", respStruct.Response, "\033[0m")
		if respStruct.Done {
			println()
			return nil, respStruct
		}
	}

	defer resp.Body.Close()

	return nil, RespStruct{}
}
