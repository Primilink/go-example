package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ========================================
// TYPES - matching OpenAI API structure
// ========================================

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// For streaming responses
type StreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// ========================================
// CLIENT
// ========================================

type OpenAIClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1",
		Client:  &http.Client{},
	}
}

// ========================================
// NON-STREAMING CHAT
// ========================================

func (c *OpenAIClient) Chat(messages []Message, model string) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chatResp, nil
}

// ========================================
// STREAMING CHAT
// ========================================

func (c *OpenAIClient) ChatStream(messages []Message, model string, onChunk func(content string)) error {
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read SSE stream
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// SSE format: "data: {json}"
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // Skip malformed chunks
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			onChunk(chunk.Choices[0].Delta.Content)
		}
	}

	return nil
}

// ========================================
// MAIN - Interactive Chat Demo
// ========================================

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable not set")
		fmt.Println("Set it with: export OPENAI_API_KEY=your-key-here")
		os.Exit(1)
	}

	client := NewOpenAIClient(apiKey)
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant. Keep responses concise."},
	}

	fmt.Println("=== OpenAI Chat in Go ===")
	fmt.Println("Type 'quit' to exit, 'stream' to toggle streaming mode")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	streaming := true

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "stream" {
			streaming = !streaming
			fmt.Printf("Streaming mode: %v\n\n", streaming)
			continue
		}

		// Add user message to history
		messages = append(messages, Message{Role: "user", Content: input})

		fmt.Print("Assistant: ")

		if streaming {
			// Streaming response
			var fullResponse strings.Builder
			err := client.ChatStream(messages, "gpt-4o-mini", func(content string) {
				fmt.Print(content)
				fullResponse.WriteString(content)
			})
			fmt.Println()

			if err != nil {
				fmt.Printf("Error: %v\n\n", err)
				// Remove failed message from history
				messages = messages[:len(messages)-1]
				continue
			}

			// Add assistant response to history
			messages = append(messages, Message{Role: "assistant", Content: fullResponse.String()})
		} else {
			// Non-streaming response
			resp, err := client.Chat(messages, "gpt-4o-mini")
			if err != nil {
				fmt.Printf("Error: %v\n\n", err)
				// Remove failed message from history
				messages = messages[:len(messages)-1]
				continue
			}

			if len(resp.Choices) > 0 {
				content := resp.Choices[0].Message.Content
				fmt.Println(content)
				messages = append(messages, Message{Role: "assistant", Content: content})

				// Show token usage
				fmt.Printf("(tokens: %d prompt, %d completion, %d total)\n",
					resp.Usage.PromptTokens,
					resp.Usage.CompletionTokens,
					resp.Usage.TotalTokens)
			}
		}
		fmt.Println()
	}
}
