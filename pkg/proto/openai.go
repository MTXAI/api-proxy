package proto

// ChatResponse completed response on site: https://platform.openai.com/docs/api-reference/chat
type ChatResponse struct {
	Model string             `json:"model"`
	Usage *ChatResponseUsage `json:"usage"`
	Error *ChatResponseError `json:"error"`
}

type ChatResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatResponseError struct {
	Type string `json:"type"`
	Code string `json:"code"`
}
