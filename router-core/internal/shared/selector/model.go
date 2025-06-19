package selector

type SelectorRequest struct {
	UserPrompt string `json:"user_prompt"`
}

type SelectorResponse struct {
	ToolID  int    `json:"tool_id"`
	Message string `json:"message"`
}
