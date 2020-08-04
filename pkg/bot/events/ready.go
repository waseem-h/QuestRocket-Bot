package events

import (
    "github.com/bwmarrin/discordgo"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/config"
)

func Ready(session *discordgo.Session, event *discordgo.Ready, config *config.Config) {
    err := session.UpdateStatus(0, "")
    if err != nil {
        logrus.Error(err)
    }
}
