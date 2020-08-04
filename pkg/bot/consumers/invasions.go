package consumers

import (
    "encoding/json"

    "github.com/adjust/rmq"
    "github.com/bwmarrin/discordgo"

    "github.com/spiri2/Quests/pkg/bot/embeds"
    "github.com/spiri2/Quests/pkg/config"
    "github.com/spiri2/Quests/pkg/models"
    modelutils "github.com/spiri2/Quests/pkg/models/utils"
    "github.com/spiri2/Quests/pkg/utils"
)

type InvasionConsumer struct {
    session *discordgo.Session
    config  *config.Config
}

func NewInvasionConsumer(session *discordgo.Session, config *config.Config) *InvasionConsumer {
    return &InvasionConsumer{
        config:  config,
        session: session,
    }
}

func (consumer InvasionConsumer) Consume(delivery rmq.Delivery) {
    var invasion models.Invasion
    if err := json.Unmarshal([]byte(delivery.Payload()), &invasion); err != nil {
        delivery.Reject()
        return
    }

    coordinates := invasion.Lat + "," + invasion.Lng
    embed := embeds.Invasion(consumer.config.Discord.Embeds, modelutils.GetGruntFromID(invasion.Character), invasion.Name, utils.FormatToRemainingTime(invasion.InvasionEnd), coordinates)

    _, err := consumer.session.ChannelMessageSendEmbed(consumer.config.Bot.Invasion.ChannelID, embed)
    if !rejectOnError(err, delivery) {
        delivery.Ack()
    }
}
