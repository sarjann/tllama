package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/tllama/config"
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

func PromptRequest(prompt string, context []int, firstRun bool, config config.Config) (error, RespStruct) {
	ollamaUrl, err := url.Parse(config.OllamaUrl)
	if err != nil {
		panic(err)
	}

	apiUrl := ollamaUrl.JoinPath("api")

	generateUrl := apiUrl.JoinPath("generate")

	var promptBody = ReqStruct{
		Model:   config.DefaultModel,
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

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, RespStruct{}
	}
	var respStruct RespStruct

	fmt.Println(string(body))
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return err, RespStruct{}
	}

	return nil, respStruct
}
