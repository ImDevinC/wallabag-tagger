package core

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/google/generative-ai-go/genai"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//go:embed resources/*
var templatesFS embed.FS

func renderPrompt(templatePath string, data any) string {
	// init template
	tmpl, err := template.ParseFS(templatesFS, templatePath)
	if err != nil {
		log.Fatal().Msgf("Error parsing template: %s", templatePath)
	}

	// render template
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		log.Fatal().Msgf("Error rendering now playing: %s", templatePath)
	}

	return tpl.String()
}

func GeminiGetTags(config Config, content string) (string, error) {
	var err error

	// init client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GoogleAIApiKey))
	if err != nil {
		log.Fatal().Msg("Failed to create GOOGLE AI client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"

	// submit
	p := bluemonday.StripTagsPolicy()
	contentSanitized := p.Sanitize(
		content,
	)
	prompt := renderPrompt("resources/prompt.txt", map[string]interface{}{
		"Content": contentSanitized,
	})
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))

	var output string
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Warn().Msg("Failed to generate text")
		}

		if resp.Candidates != nil {
			for _, v := range resp.Candidates {
				for _, k := range v.Content.Parts {
					output += fmt.Sprint(k.(genai.Text))
				}
			}
		}
	}

	return output, err
}
