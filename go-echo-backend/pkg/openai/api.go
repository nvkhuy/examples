package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"io"
	"net/http"
)

type Client struct {
	cfg    *config.Configuration
	logger *logger.Logger
}

func New(cfg *config.Configuration) *Client {
	var client = &Client{
		cfg:    cfg,
		logger: logger.New("openai"),
	}

	return client
}

type CreateThreadPayload struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type CreateThreadResp struct {
	ID        string      `json:"id,omitempty"`
	Object    string      `json:"object,omitempty"`
	CreatedAt int64       `json:"created_at,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

func (c *Client) CreateThread(payload CreateThreadPayload) (resp CreateThreadResp, err error) {
	url := "https://api.openai.com/v1/threads"
	method := "POST"
	payloadBytes, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.OpenAIKey))
	req.Header.Add("OpenAI-Beta", "assistants=v1")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

type CreateMessagePayload struct {
	ThreadID string `json:"thread_id"`
	Role     string `json:"role"`
	Content  string `json:"content"`
}

type MessageContent struct {
	Type string `json:"type"`
	Text struct {
		Value       string   `json:"value"`
		Annotations []string `json:"annotations"`
	} `json:"text"`
}

type CreateMessageResp struct {
	ID          string                 `json:"id"`
	Object      string                 `json:"object"`
	CreatedAt   int64                  `json:"created_at"`
	ThreadID    string                 `json:"thread_id"`
	Role        string                 `json:"role"`
	Content     []MessageContent       `json:"content"`
	FileIDs     []string               `json:"file_ids"`
	AssistantID interface{}            `json:"assistant_id"`
	RunID       interface{}            `json:"run_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func (c *Client) CreateMessage(payload CreateMessagePayload) (resp CreateMessageResp, err error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", payload.ThreadID)
	method := "POST"

	payloadBytes, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.OpenAIKey))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

type RunAssistantPayload struct {
	ThreadID    string `json:"thread_id,omitempty"`
	AssistantID string `json:"assistant_id,omitempty"`
}

type RunAssistantResp struct {
	ID           string                 `json:"id"`
	Object       string                 `json:"object"`
	CreatedAt    int64                  `json:"created_at"`
	AssistantID  string                 `json:"assistant_id"`
	ThreadID     string                 `json:"thread_id"`
	Status       string                 `json:"status"`
	StartedAt    int64                  `json:"started_at"`
	ExpiresAt    int64                  `json:"expires_at"`
	CancelledAt  int64                  `json:"cancelled_at"`
	FailedAt     int64                  `json:"failed_at"`
	CompletedAt  int64                  `json:"completed_at"`
	LastError    int64                  `json:"last_error"`
	Model        string                 `json:"model"`
	Instructions string                 `json:"instructions"`
	Tools        []interface{}          `json:"tools"`
	FileIDs      []interface{}          `json:"file_ids"`
	Metadata     map[string]interface{} `json:"metadata"`
}

func (c *Client) RunAssistant(payload RunAssistantPayload) (resp RunAssistantResp, err error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", payload.ThreadID)
	method := "POST"

	if payload.AssistantID == "" {
		payload.AssistantID = c.cfg.OpenAIAssistant
	}

	payloadBytes, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.OpenAIKey))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

type RunStatusPayload struct {
	ThreadID string `json:"thread_id"`
	RunID    string `json:"run_id"`
}

func (c *Client) RunStatus(payload RunStatusPayload) (resp RunAssistantResp, err error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs/%s", payload.ThreadID, payload.RunID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.OpenAIKey))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

type AssistantReplyPayload struct {
	ThreadID string `json:"thread_id"`
	Limit    int    `json:"limit"`
}

type TextContent struct {
	Type string `json:"type"`
	Text struct {
		Value       string   `json:"value"`
		Annotations []string `json:"annotations"`
	} `json:"text"`
}

type ThreadMessage struct {
	ID          string                 `json:"id"`
	Object      string                 `json:"object"`
	CreatedAt   int64                  `json:"created_at"`
	ThreadID    string                 `json:"thread_id"`
	Role        string                 `json:"role"`
	Content     []TextContent          `json:"content"`
	FileIDs     []interface{}          `json:"file_ids"`
	AssistantID string                 `json:"assistant_id"`
	RunID       string                 `json:"run_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type AssistantReplyResp struct {
	Object  string          `json:"object"`
	Data    []ThreadMessage `json:"data"`
	FirstID string          `json:"first_id"`
	LastID  string          `json:"last_id"`
	HasMore bool            `json:"has_more"`
}

func (c *Client) GetAssistantResponse(payload AssistantReplyPayload) (resp AssistantReplyResp, err error) {
	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", payload.ThreadID)
	if payload.Limit != 0 {
		url = url + fmt.Sprintf("?limit=%d", payload.Limit)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v1")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.OpenAIKey))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}
