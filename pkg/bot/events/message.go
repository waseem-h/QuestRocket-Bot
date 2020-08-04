package events

import (
    "github.com/bwmarrin/discordgo"

    "github.com/spiri2/Quests/pkg/config"
)

func Message(session *discordgo.Session, event *discordgo.MessageCreate, config *config.Config) {
    // Ignore all bots
    if event.Author.Bot {
        return
    }

    if event.Content != "" && string(event.Content[0]) != config.Bot.CommandPrefix {
        return
    }
}
