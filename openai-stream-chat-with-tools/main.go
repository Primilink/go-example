package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
)

// ========================================
// API TYPES
// ========================================

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

type StreamChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role      string     `json:"role,omitempty"`
			Content   string     `json:"content,omitempty"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id,omitempty"`
				Type     string `json:"type,omitempty"`
				Function struct {
					Name      string `json:"name,omitempty"`
					Arguments string `json:"arguments,omitempty"`
				} `json:"function"`
			} `json:"tool_calls,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// ========================================
// MATH TOOLS DEFINITION
// ========================================

var mathTools = []Tool{
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "add",
			Description: "Add two numbers together",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{"type": "number", "description": "First number"},
					"b": map[string]interface{}{"type": "number", "description": "Second number"},
				},
				"required": []string{"a", "b"},
			},
		},
	},
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "subtract",
			Description: "Subtract second number from first number",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{"type": "number", "description": "Number to subtract from"},
					"b": map[string]interface{}{"type": "number", "description": "Number to subtract"},
				},
				"required": []string{"a", "b"},
			},
		},
	},
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "multiply",
			Description: "Multiply two numbers together",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{"type": "number", "description": "First number"},
					"b": map[string]interface{}{"type": "number", "description": "Second number"},
				},
				"required": []string{"a", "b"},
			},
		},
	},
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "divide",
			Description: "Divide first number by second number",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{"type": "number", "description": "Dividend (number to divide)"},
					"b": map[string]interface{}{"type": "number", "description": "Divisor (number to divide by)"},
				},
				"required": []string{"a", "b"},
			},
		},
	},
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "power",
			Description: "Raise a number to a power",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"base":     map[string]interface{}{"type": "number", "description": "Base number"},
					"exponent": map[string]interface{}{"type": "number", "description": "Exponent"},
				},
				"required": []string{"base", "exponent"},
			},
		},
	},
	{
		Type: "function",
		Function: ToolFunction{
			Name:        "sqrt",
			Description: "Calculate the square root of a number",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"n": map[string]interface{}{"type": "number", "description": "Number to find square root of"},
				},
				"required": []string{"n"},
			},
		},
	},
}

// ========================================
// TOOL EXECUTION
// ========================================

func executeTool(name string, argsJSON string) (string, error) {
	var args map[string]float64
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	var result float64
	var resultStr string

	switch name {
	case "add":
		result = args["a"] + args["b"]
		resultStr = fmt.Sprintf("%.6g + %.6g = %.6g", args["a"], args["b"], result)

	case "subtract":
		result = args["a"] - args["b"]
		resultStr = fmt.Sprintf("%.6g - %.6g = %.6g", args["a"], args["b"], result)

	case "multiply":
		result = args["a"] * args["b"]
		resultStr = fmt.Sprintf("%.6g × %.6g = %.6g", args["a"], args["b"], result)

	case "divide":
		if args["b"] == 0 {
			return "Error: division by zero", nil
		}
		result = args["a"] / args["b"]
		resultStr = fmt.Sprintf("%.6g ÷ %.6g = %.6g", args["a"], args["b"], result)

	case "power":
		result = math.Pow(args["base"], args["exponent"])
		resultStr = fmt.Sprintf("%.6g ^ %.6g = %.6g", args["base"], args["exponent"], result)

	case "sqrt":
		if args["n"] < 0 {
			return "Error: cannot calculate square root of negative number", nil
		}
		result = math.Sqrt(args["n"])
		resultStr = fmt.Sprintf("√%.6g = %.6g", args["n"], result)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	return resultStr, nil
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

// StreamResponse holds accumulated tool calls during streaming
type StreamResponse struct {
	Content   string
	ToolCalls []ToolCall
}

func (c *OpenAIClient) ChatStreamWithTools(messages []Message, tools []Tool, onContent func(string)) (*StreamResponse, error) {
	reqBody := ChatRequest{
		Model:    "gpt-4o-mini",
		Messages: messages,
		Tools:    tools,
		Stream:   true,
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

	// Accumulate response
	result := &StreamResponse{
		ToolCalls: []ToolCall{},
	}

	// Track tool calls being built
	toolCallsMap := make(map[int]*ToolCall)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta

		// Handle content
		if delta.Content != "" {
			result.Content += delta.Content
			onContent(delta.Content)
		}

		// Handle tool calls
		for _, tc := range delta.ToolCalls {
			if _, exists := toolCallsMap[tc.Index]; !exists {
				toolCallsMap[tc.Index] = &ToolCall{
					ID:   tc.ID,
					Type: tc.Type,
				}
				toolCallsMap[tc.Index].Function.Name = tc.Function.Name
			}
			toolCallsMap[tc.Index].Function.Arguments += tc.Function.Arguments
		}
	}

	// Convert map to slice
	for i := 0; i < len(toolCallsMap); i++ {
		if tc, exists := toolCallsMap[i]; exists {
			result.ToolCalls = append(result.ToolCalls, *tc)
		}
	}

	return result, nil
}

// ========================================
// MAIN
// ========================================

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	client := NewOpenAIClient(apiKey)
	messages := []Message{
		{
			Role: "system",
			Content: `You are a helpful math assistant. You have access to math tools for calculations.
Always use the tools for any mathematical operations - don't calculate in your head.
Available tools: add, subtract, multiply, divide, power, sqrt.
After getting tool results, explain the answer clearly.`,
		},
	}

	fmt.Println("=== OpenAI Chat with Math Tools ===")
	fmt.Println("Try: 'What is 25 * 4 + 10?' or 'Calculate the square root of 144'")
	fmt.Println("Type 'quit' to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

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

		messages = append(messages, Message{Role: "user", Content: input})

		// Loop to handle multiple tool calls
		for {
			fmt.Print("Assistant: ")

			resp, err := client.ChatStreamWithTools(messages, mathTools, func(content string) {
				fmt.Print(content)
			})

			if err != nil {
				fmt.Printf("\nError: %v\n\n", err)
				messages = messages[:len(messages)-1]
				break
			}

			// If no tool calls, we're done
			if len(resp.ToolCalls) == 0 {
				fmt.Println()
				if resp.Content != "" {
					messages = append(messages, Message{Role: "assistant", Content: resp.Content})
				}
				break
			}

			// Handle tool calls
			fmt.Println() // newline after any partial content

			// Add assistant message with tool calls
			assistantMsg := Message{
				Role:      "assistant",
				Content:   resp.Content,
				ToolCalls: resp.ToolCalls,
			}
			messages = append(messages, assistantMsg)

			// Execute each tool and add results
			for _, tc := range resp.ToolCalls {
				fmt.Printf("  [Calling %s(%s)]\n", tc.Function.Name, tc.Function.Arguments)

				result, err := executeTool(tc.Function.Name, tc.Function.Arguments)
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
				}
				fmt.Printf("  [Result: %s]\n", result)

				// Add tool result message
				messages = append(messages, Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})
			}

			// Continue loop to get final response after tool execution
		}
		fmt.Println()
	}
}
