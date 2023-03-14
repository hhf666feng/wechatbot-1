package gtp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/869413421/wechatbot/config"
)

const BASEURL = "https://api.openai.com/v1/chat/completions"

type RoleType string
type ModelType string

const ModelGpt35Turbo = "gpt-3.5-turbo"

const (
	RoleUser      RoleType = "user"
	RoleAssistant RoleType = "assistant"
	RoleSystem    RoleType = "system"
)

type Request struct {
	Model            ModelType   `json:"model"`
	Messages         []*Message  `json:"messages"`
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"top_p,omitempty"`
	N                int         `json:"n,omitempty"`
	Stream           bool        `json:"stream,omitempty"`
	Stop             interface{} `json:"stop,omitempty"`
	MaxTokens        int         `json:"max_tokens,omitempty"`
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"`
	LogitBias        interface{} `json:"logit_bias,omitempty"`
	User             string      `json:"user,omitempty"`
}
type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Choices []*Choice `json:"choices"`
	Usage   *Usage    `json:"usage"`
	Error   *Error    `json:"error,omitempty"`
}

type Message struct {
	Role    RoleType `json:"role,omitempty"`
	Content string   `json:"content"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {

	request := &Request{
		Model: ModelGpt35Turbo,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: msg,
			},
		},
		Temperature:      0.7,
		TopP:             1,
		N:                1,
		Stream:           false,
		Stop:             []string{"\r"},
		MaxTokens:        4096,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}

	requestData, err := json.Marshal(request)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &Response{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v.Message.Content
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
