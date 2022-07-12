package models

type Parameter struct {
	ErrorMsg          string `json:"errorMsg,omitempty"`
	FalsePositive 	  bool `json:"falsePositive,omitempty"`	
	RegexCheck        string `json:"regexCheck,omitempty"`
	ErrorType         string `json:"errorType"`
	URL               string `json:"url"`
	URLMain           string `json:"urlMain"`
	URLProbe          string `json:"urlProbe,omitempty"`
	UsernameClaimed   string `json:"username_claimed"`
	UsernameUnclaimed string `json:"username_unclaimed"`
}
