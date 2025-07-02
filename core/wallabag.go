package core

import (
	"github.com/Strubbl/wallabago/v9"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func WallabagInit() {
	config, err := GetConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	wallabagConfig := wallabago.WallabagConfig{
		WallabagURL:  config.WallabagUrl,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		UserName:     config.Username,
		UserPassword: config.Password,
	}
	wallabago.SetConfig(wallabagConfig)
}

var wallabagFetchLimit = 300

func WallabagGetEntries() wallabago.Entries {
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, wallabagFetchLimit, "", 0, -1, "", "")
	if err != nil {
		log.Fatal().Msg("Cannot obtain articles from Wallabag")
	}

	return entries
}

func WallabagWriteTags(entry wallabago.Item, tags []string) {
	err := wallabago.AddEntryTags(entry.ID, tags...)
	if err != nil {
		log.Err(err).Msgf("Cannot assign tags to article: %s", entry.Title)
	}
}
