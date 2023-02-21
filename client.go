package gptclient

import (
	"encoding/json"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	baseURL     = "https://api.openai.com/v1/completions"
	maxTokens   = 2000
	temperature = 0.7
	engine      = "text-davinci-003"
)

type chatGPTResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []choice               `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type choice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	LogProbs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type chatGPTRequest struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

func Dialogue(msg, apiKey string) (string, error) {
	request := chatGPTRequest{
		Model:            engine,
		Prompt:           msg,
		MaxTokens:        maxTokens,
		Temperature:      temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(baseURL)
	req.SetBody(requestBody)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &fasthttp.Client{}

	if err := client.DoTimeout(req, resp, 60*time.Second); err != nil {
		return "", err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return "", err
	}

	response := &chatGPTResponse{}
	if err := json.Unmarshal(resp.Body(), response); err != nil {
		return "", err
	}

	var reply string

	if response.Choices != nil && len(response.Choices) > 0 {
		reply = response.Choices[0].Text
	}

	log.Printf("response text: %s", reply)

	return reply, nil
}


