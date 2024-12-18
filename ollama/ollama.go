package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sarjann/tllama/config"
	"net/http"
	"net/url"
)

type RespPromptStruct struct {
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

type ReqPromptStruct struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context"`
}

func PromptRequest(prompt string, context []int, firstRun bool, conf config.Config) (error, RespPromptStruct) {
	ollamaUrl, err := url.Parse(conf.OllamaUrl)
	if err != nil {
		panic(err)
	}

	apiUrl := ollamaUrl.JoinPath("api")

	generateUrl := apiUrl.JoinPath("generate")

	var promptBody = ReqPromptStruct{
		Model:   conf.Model,
		Prompt:  prompt,
		Stream:  true,
		Context: context,
	}

	promptJSONBody, err := json.Marshal(promptBody)
	if err != nil {
		return err, RespPromptStruct{}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", generateUrl.String(), bytes.NewReader(promptJSONBody))
	if err != nil {
		return err, RespPromptStruct{}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err, RespPromptStruct{}
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var respStruct RespPromptStruct
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

	return nil, RespPromptStruct{}
}

type ReqPullModelStruct struct {
	Name   string `json:"name"`
	Stream bool   `json:"stream"`
}

type RespPullModelStruct struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int    `json:"total,omitempty"`
	Completed int    `json:"completed,omitempty"`
}

func PullModel(modelName string, conf config.Config) (error, RespPullModelStruct) {
	ollamaUrl, err := url.Parse(conf.OllamaUrl)
	if err != nil {
		panic(err)
	}

	apiUrl := ollamaUrl.JoinPath("api")

	pullUrl := apiUrl.JoinPath("pull")

	var pullBody = ReqPullModelStruct{
		Name:   modelName,
		Stream: true,
	}

	pullJSONBody, err := json.Marshal(pullBody)
	if err != nil {
		return err, RespPullModelStruct{}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", pullUrl.String(), bytes.NewReader(pullJSONBody))
	if err != nil {
		return err, RespPullModelStruct{}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err, RespPullModelStruct{}
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var respStruct RespPullModelStruct
		err := decoder.Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(respStruct.Status)
		if respStruct.Completed != 0 {
			fmt.Print("\033[92m", respStruct.Completed, "\033[0m\r")
			if respStruct.Completed == respStruct.Total {
				fmt.Println("Finished with total: ", respStruct.Total)
				return nil, respStruct
			}
			fmt.Print("\r")
		}
		// fmt.Println(respStruct.Status)
	}

	defer resp.Body.Close()

	return nil, RespPullModelStruct{}
}

type ReqDeleteModelStruct struct {
	Name string `json:"name"`
}

func DeleteModel(modelName string, conf config.Config) error {
	ollamaUrl, err := url.Parse(conf.OllamaUrl)
	if err != nil {
		panic(err)
	}

	apiUrl := ollamaUrl.JoinPath("api")

	deleteUrl := apiUrl.JoinPath("delete")

	client := &http.Client{}

	var deleteBody = ReqDeleteModelStruct{
		Name: modelName,
	}
	deleteJSONBody , err := json.Marshal(deleteBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", deleteUrl.String(), bytes.NewReader(deleteJSONBody))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Failed to delete model")
	}
	fmt.Println(resp.Status)
	return nil
}
