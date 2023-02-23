package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Msg struct {
	Prompt string
}

type Chat struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func main() {
	e := gin.Default()
	e.Use(Cors())
	e.POST("/chat", func(c *gin.Context) {
		var p Msg
		data, _ := c.GetRawData()
		_ = json.Unmarshal(data, &p)
		url := "https://api.openai.com/v1/completions"
		s := `{"model": "text-davinci-003", "prompt": "` + p.Prompt + `", "temperature": 0, "max_tokens": 2048}`
		body := strings.NewReader(s)
		request, err := http.NewRequest("POST", url, body)
		if err != nil {
			fmt.Println(err.Error())
		}

		
		// 填入你的keyid
		request.Header.Add("Authorization", "Bearer XXXXXXXXXXXXXXXXX")


		request.Header.Add("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println(err.Error())
		}
		resbody, _ := ioutil.ReadAll(response.Body)
		var chat Chat
		_ = json.Unmarshal(resbody, &chat)
		log.Println(p.Prompt, ": ", chat.Choices[0].Text)
		c.JSON(200, chat)
	})
	_ = e.Run(":8084")
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}
