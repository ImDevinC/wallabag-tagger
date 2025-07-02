package core

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Strubbl/wallabago/v9"
	"github.com/rs/zerolog/log"
)

var ErrEmptyContent = errors.New("empty content")

type Tags struct {
	Tag []string `json:"tags"`
}

func isSkipEntry(entry wallabago.Item) bool {
	isSkip := false
	if len(entry.Tags) >= 1 { // if already has tags
		for _, tag := range entry.Tags {
			if strings.HasPrefix(tag.Label, "llm") {
				isSkip = true
				continue
			}
		}
	}

	return isSkip
}

func LLMTags() {
	config, err := GetConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	entries := WallabagGetEntries()
	log.Debug().Msgf("found %d items", len(entries.Embedded.Items))
	for _, entry := range entries.Embedded.Items {
		// skip if already tagged via LLM
		isSkip := isSkipEntry(entry)
		if isSkip {
			continue
		}
		log.Info().Msgf("Processing article: %s", entry.Title)

		var tagsStr string
		if config.GoogleAIApiKey != "" {
			// get tags from llm
			tagsStr, err = GeminiGetTags(config, entry.Content)
		} else if config.Ollama.URL != "" && config.Ollama.Model != "" {
			tagsStr, err = OllamaGetTags(config, entry.Content)
		} else {
			log.Error().Msg("no llm config")
			return
		}

		if err != nil {
			log.Error().Err(err).Msgf("failed to get content for article: %s", entry.Title)
			if errors.Is(err, ErrEmptyContent) {
				WallabagWriteTags(entry, []string{"llm-no-content"})
			}
			continue
		}

		// convert json-string to Tags struct
		var tags Tags
		err := json.Unmarshal([]byte(tagsStr), &tags)
		if err != nil {
			log.Error().Msgf("Cannot unmarshal tags: %s", tagsStr)
		}

		// add tags prefix so it doesn't conflict with manually-assigned tags
		var tagsWithPrefix []string
		for _, tag := range tags.Tag {
			tagsWithPrefix = append(tagsWithPrefix, "llm-"+tag)
		}

		// update entry tags
		WallabagWriteTags(entry, tagsWithPrefix)
	}
}
