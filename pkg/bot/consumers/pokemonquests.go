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

type PokemonQuestConsumer struct {
    session *discordgo.Session
    config  *config.Config
}

func NewPokemonQuestConsumer(session *discordgo.Session, config *config.Config) *PokemonQuestConsumer {
    return &PokemonQuestConsumer{
        config:  config,
        session: session,
    }
}

func (consumer PokemonQuestConsumer) Consume(delivery rmq.Delivery) {
    var quest models.Quest
    if err := json.Unmarshal([]byte(delivery.Payload()), &quest); err != nil {
        logrus.Error(err)
        delivery.Reject()
        return
    }

    coordinates := quest.Lat + "," + quest.Lng
    embed := embeds.PokemonQuest(consumer.config.Discord.Embeds, quest.RewardsString, quest.Name, quest.ConditionsString, utils.FormatToRemainingTime(quest.Expiration), quest.RewardsIds, coordinates)

    _, err := consumer.session.ChannelMessageSendEmbed(consumer.config.Bot.PokemonQuest.ChannelID, embed)
    if !rejectOnError(err, delivery) {
        delivery.Ack()
    }
}
