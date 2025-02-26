package checkUsers

import (
	"gopher-find/cmd/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRandomUserAgent(t *testing.T) {
	agents := make(map[string]bool)
	for i := 0; i < 100; i++ {
		agent := getRandomUserAgent()
		if agent == "" {
			t.Error("Got empty user agent")
		}
		agents[agent] = true
	}

	if len(agents) < 2 {
		t.Error("Random user agent selection doesn't seem random")
	}
}

func TestURLWithUsername(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		username string
		expected string
	}{
		{
			name:     "Basic replacement",
			url:      "http://example.com/{}",
			username: "testuser",
			expected: "http://example.com/testuser",
		},
		{
			name:     "Multiple replacements",
			url:      "http://example.com/{}/profile/{}",
			username: "testuser",
			expected: "http://example.com/testuser/profile/testuser",
		},
		{
			name:     "No placeholder",
			url:      "http://example.com/profile",
			username: "testuser",
			expected: "http://example.com/profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := URLWithUsername(tt.url, tt.username)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    Response
		parameter   models.Parameter
		expectValid bool
	}{
		{
			name: "Valid response",
			response: Response{
				code: 200,
				body: "This is a valid profile page with sufficient content length. This is a valid profile page with sufficient content length. This is a valid profile page with sufficient content length",
			},
			parameter:   models.Parameter{},
			expectValid: true,
		},
		{
			name: "Invalid status code",
			response: Response{
				code: 404,
				body: "Not Found",
			},
			parameter:   models.Parameter{},
			expectValid: false,
		},
		{
			name: "Content too short",
			response: Response{
				code: 200,
				body: "short",
			},
			parameter:   models.Parameter{},
			expectValid: false,
		},
		{
			name: "Contains error indicator",
			response: Response{
				code: 200,
				body: "This page contains 404 not found message",
			},
			parameter:   models.Parameter{},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateResponse(tt.response, tt.parameter)
			if result != tt.expectValid {
				t.Errorf("Expected validity %v, got %v", tt.expectValid, result)
			}
		})
	}
}

func TestCheckForDelayedRedirect(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		expectRedirect bool
	}{
		{
			name:           "No redirect",
			responseBody:   "<html><body>Normal content</body></html>",
			expectRedirect: false,
		},
		{
			name:           "Window location redirect",
			responseBody:   "<script>window.location.href='http://example.com';</script>",
			expectRedirect: true,
		},
		{
			name:           "Meta refresh redirect",
			responseBody:   "<meta http-equiv=\"refresh\" content=\"0;url=http://example.com\">",
			expectRedirect: true,
		},
		{
			name:           "setTimeout redirect",
			responseBody:   "setTimeout(function() { window.location = 'http://example.com'; }, 1000);",
			expectRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := Response{body: tt.responseBody}
			result := checkForDelayedRedirect(response)
			if result != tt.expectRedirect {
				t.Errorf("Expected redirect detection %v, got %v", tt.expectRedirect, result)
			}
		})
	}
}

func TestDoReq(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if User-Agent header is set
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent header not set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	}))
	defer server.Close()

	// Test valid request
	response, err := doReq(server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.code != 200 {
		t.Errorf("Expected status code 200, got %d", response.code)
	}
	if response.body != "Test response" {
		t.Errorf("Expected body 'Test response', got '%s'", response.body)
	}

	// Test invalid URL
	_, err = doReq("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
