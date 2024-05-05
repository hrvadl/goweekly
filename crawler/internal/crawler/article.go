package crawler

import (
	"fmt"
	"strings"
)

const (
	minTokens    = 4
	sponsorLabel = "sponsor"
)

const (
	urlTokenIdx = iota
	headerTokenIdx
	startContentIdx
)

type Article struct {
	URL         string
	Header      string
	Content     string
	Author      string
	IsSponsored bool
}

/*
Extracts Article's fields from tokens

	If HTML was parsed with getArticleFromStream()
	token slice order should be like this:
	n := len(tokens)
	tokens[0] - URL
	tokens[1] - Article header
	tokens[2]:token[n-1] - Content
	tokens[n-1] - Author or Sponsored label *
*/
func newArticleFromTextTokens(tokens []string) (*Article, error) {
	length := len(tokens)
	if length < minTokens {
		return nil, fmt.Errorf("there should be at least 4 tokens, but got %d", len(tokens))
	}

	return &Article{
		URL:         tokens[urlTokenIdx],
		Header:      tokens[headerTokenIdx],
		Content:     strings.Join(tokens[startContentIdx:length-1], ""),
		Author:      tokens[length-1],
		IsSponsored: tokens[length-1] == sponsorLabel,
	}, nil
}
