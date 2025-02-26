package models

// Parameter represents a website's configuration for username checking
type Parameter struct {
	ErrorMsg          interface{}       `json:"errorMsg,omitempty"` // Can be string or []string
	ErrorType         string            `json:"errorType"`
	URL               string            `json:"url"`
	URLMain           string            `json:"urlMain"`
	URLProbe          string            `json:"urlProbe,omitempty"`
	UsernameClaimed   string            `json:"username_claimed"`
	UsernameUnclaimed string            `json:"username_unclaimed,omitempty"`
	Headers           map[string]string `json:"headers,omitempty"`
	RequestMethod     string            `json:"request_method,omitempty"`
	RequestPayload    interface{}       `json:"request_payload,omitempty"`
	ErrorURL          string            `json:"errorUrl,omitempty"`
	IsNSFW            bool              `json:"isNSFW,omitempty"`
	RegexCheck        string            `json:"regexCheck,omitempty"`
}

// GetErrorMessages returns all error messages as a string slice
func (p Parameter) GetErrorMessages() []string {
	var messages []string

	switch v := p.ErrorMsg.(type) {
	case string:
		messages = append(messages, v)
	case []interface{}:
		for _, msg := range v {
			if str, ok := msg.(string); ok {
				messages = append(messages, str)
			}
		}
	}

	return messages
}
