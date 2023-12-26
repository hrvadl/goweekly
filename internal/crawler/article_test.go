package crawler

import (
	"reflect"
	"testing"
)

func TestCreateArticleFromTokens(t *testing.T) {
	tc := []struct {
		name     string
		tokens   []string
		expected *Article
		err      bool
	}{
		{
			name:   "Should parse the correct base case",
			tokens: []string{"http://example.com", "Header", "Content", "Author"},
			expected: &Article{
				URL:         "http://example.com",
				Header:      "Header",
				Content:     "Content",
				Author:      "Author",
				IsSponsored: false,
			},
		},
		{
			name:   "Should parse the correct two-entry-content case",
			tokens: []string{"http://example.com", "Header", "Content ", "More Content", "Author"},
			expected: &Article{
				URL:         "http://example.com",
				Header:      "Header",
				Content:     "Content More Content",
				Author:      "Author",
				IsSponsored: false,
			},
		},
		{
			name:   "Should parse the correct sponsored case",
			tokens: []string{"http://example.com", "Header", "Content ", "More Content", sponsorLabel},
			expected: &Article{
				URL:         "http://example.com",
				Header:      "Header",
				Content:     "Content More Content",
				Author:      "sponsor",
				IsSponsored: true,
			},
		},
		{
			name:   "Should not parse",
			tokens: []string{"http://example.com", "Header", "Content"},
			err:    true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			res, err := newArticleFromTextTokens(tt.tokens)
			if tt.err && err == nil {
				t.Fatal("Expected to get and error")
			}

			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, res)
			}
		})
	}
}
