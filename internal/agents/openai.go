package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OpenAPIEndpoint = "https://api.openai.com/v1/chat/completions"
)

type OpenAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type OpenAPI struct {
	httpClient *http.Client
	ctx        context.Context
	apiKey     string
	model      string
}

func NewOpenAI(ctx context.Context, apiKey, model string, httpClient *http.Client) *OpenAPI {
	o := &OpenAPI{
		ctx:        ctx,
		apiKey:     apiKey,
		model:      model,
		httpClient: httpClient,
	}

	if httpClient == nil {
		o.httpClient = &http.Client{
			Timeout: time.Second * 120,
		}
	}

	return o
}

type OpenAPIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *OpenAPI) Query(messages []OpenAPIMessage) (OpenAPIResponse, error) {
	var response OpenAPIResponse

	// Default to system prompt if not provided
	if len(messages) == 0 || messages[0].Role != "system" {
		defaultSystem := OpenAPIMessage{
			Role: "system",
			Content: "You are an intelligent and supportive virtual fitness coach named GymBuddy. You help users with workouts, meals, motivation, progress tracking, and fitness education. Always respond in a friendly, practical, and encouraging tone. Use emojis to make responses fun and engaging.",
		}
		messages = append([]OpenAPIMessage{defaultSystem}, messages...)
	}

	bs, err := json.Marshal(map[string]interface{}{
		"model":    o.model,
		"messages": messages,
	})

	if err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(o.ctx, "POST", OpenAPIEndpoint, bytes.NewBuffer(bs))
	if err != nil {
		return response, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return response, fmt.Errorf("API error: %s", response.Error.Message)
	}

	

	if len(response.Choices) == 0 {
		return response, errors.New("no choices returned from API")
	}

	return response, nil
}
