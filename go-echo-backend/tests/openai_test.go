package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/openai"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestOpenAI_CreateThread(t *testing.T) {
	var app = initApp("local")
	client := openai.New(app.Config)
	payload := openai.CreateThreadPayload{
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: "Hi, i need size X for 100 Jeans with grey colour",
			},
		},
	}
	_, _ = client.CreateThread(payload)
	return
}

func TestOpenAI_GetAssistantResponse(t *testing.T) {
	var app = initApp("local")
	client := openai.New(app.Config)
	payload := openai.AssistantReplyPayload{
		ThreadID: "thread_l0X3oaYyScLcxQxAyBw4jBeS",
		Limit:    1,
	}
	resp, err := client.GetAssistantResponse(payload)
	assert.NoError(t, err)
	log.Println(resp)
}
