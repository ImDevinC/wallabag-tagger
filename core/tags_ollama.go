package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/prathyushnallamothu/ollamago"
)

func OllamaGetTags(config Config, content string) (string, error) {
	ctx := context.Background()

	p := bluemonday.StripTagsPolicy()
	contentSanitized := p.Sanitize(content)

	if len(contentSanitized) == 0 {
		return "", ErrEmptyContent
	}

	prompt := renderPrompt("resources/prompt.txt", map[string]interface{}{
		"Content": contentSanitized,
	})

	client := ollamago.NewClient(
		ollamago.WithTimeout(time.Minute*2),
		ollamago.WithBaseURL(config.Ollama.URL),
	)
	streamResp, streamErr := client.GenerateStream(
		ctx,
		ollamago.GenerateRequest{
			Model:  config.Ollama.Model,
			Prompt: prompt,
			Stream: true,
		},
	)

	var fullResponse strings.Builder
	for {
		select {
		case resp, ok := <-streamResp:
			if !ok {
				if fullResponse.Len() == 0 {
					return "", errors.New("failed to get message from stream")
				}
			}
			fullResponse.WriteString(resp.Response)
			if resp.Done {
				return fullResponse.String(), nil
			}
		case err := <-streamErr:
			if err != nil {
				return "", fmt.Errorf("ollama request failed. %w", err)
			}
		case <-ctx.Done():
			return "", errors.New("request timed out")
		}
	}
}
