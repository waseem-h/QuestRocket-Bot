package consumers

import (
    "encoding/json"

    "github.com/adjust/rmq"
    "github.com/bwmarrin/discordgo"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/bot/embeds"
    "github.com/spiri2/Quests/pkg/config"
    "github.com/spiri2/Quests/pkg/models"
    "github.com/spiri2/Quests/pkg/utils"
)

type StardustQuestConsumer struct {
    session *discordgo.Session
    config  *config.Config
}

func NewStardustQuestConsumer(session *discordgo.Session, config *config.Config) *StardustQuestConsumer {
    return &StardustQuestConsumer{
        config:  config,
        session: session,
    }
}

func (consumer StardustQuestConsumer) Consume(delivery rmq.Delivery) {
    var quest models.Quest
    if err := json.Unmarshal([]byte(delivery.Payload()), &quest); err != nil {
        logrus.Error(err)
        delivery.Reject()
        return
    }

    coordinates := quest.Lat + "," + quest.Lng
    embed := embeds.StardustQuest(consumer.config.Discord.Embeds, quest.RewardsString, quest.Name, quest.ConditionsString, utils.FormatToRemainingTime(quest.Expiration), coordinates)

    _, err := consumer.session.ChannelMessageSendEmbed(consumer.config.Bot.StardustQuest.ChannelID, embed)
    if !rejectOnError(err, delivery) {
        delivery.Ack()
    }
}
