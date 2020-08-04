package embeds

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/bwmarrin/discordgo"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/config"
)

func Invasion(embeds config.Embeds, gruntName string, stopName string, endTime string, coordinates string) *discordgo.MessageEmbed {
    description := "**Stop Name:** " + stopName + "\n**Expires in:** " + endTime
    iconURL := "https://github.com/Gitanjan18/StopInvasion/blob/master/type_" + gruntName + ".png?raw=true"
    name := "Team Rocket Invasion"

    return messageEmbed(embeds, strings.ToUpper(gruntName), description, iconURL, name, coordinates)
}

func PokemonQuest(embeds config.Embeds, reward string, stop string, condition string, endTime string, task string, coordinates string) *discordgo.MessageEmbed {
    description := "**Stop Name:** " + stop + "\n**Condition:** " + condition + "\n**Expires In:** " + endTime
    iconURL := "https://github.com/Gitanjan18/sprites/blob/master/" + task + ".png?raw=true"
    name := "Pokemon Quest"

    return messageEmbed(embeds, reward, description, iconURL, name, coordinates)
}

func StardustQuest(embeds config.Embeds, reward string, stop string, condition string, endTime string, coordinates string) *discordgo.MessageEmbed {
    description := "**Stop Name:** " + stop + "\n**Condition:** " + condition + "\n**Expires In:** " + endTime
    iconURL := "https://github.com/Gitanjan18/quest/blob/master/stardust_painted.png?raw=true"
    name := "Stardust Quest"

    return messageEmbed(embeds, reward, description, iconURL, name, coordinates)
}

func ItemQuest(embeds config.Embeds, reward string, stop string, condition string, endTime string, itemID string, coordinates string) *discordgo.MessageEmbed {
    itemIDInt, err := strconv.Atoi(itemID)
    if err != nil {
        logrus.Error(err)
    }

    description := "**Stop Name:** " + stop + "\n**Condition:** " + condition + "\n**Expires In:** " + endTime
    iconURL := "https://github.com/Gitanjan18/quest/blob/master/Item_" + fmt.Sprintf("%04d", itemIDInt) + ".png?raw=true"
    name := "Item Quest"

    return messageEmbed(embeds, reward, description, iconURL, name, coordinates)
}

func messageEmbed(embeds config.Embeds, title string, description string, iconURL string, name string, coordinates string) *discordgo.MessageEmbed {
    return &discordgo.MessageEmbed{
        Description: description,
        URL:         "https://discord.gg/ACaXz4",
        Color:       12390624,
        Footer: &discordgo.MessageEmbedFooter{
            IconURL: embeds.Footer.IconURL,
            Text:    embeds.Footer.Text,
        },
        Thumbnail: &discordgo.MessageEmbedThumbnail{
            URL: iconURL,
        },
        Author: &discordgo.MessageEmbedAuthor{
            Name:    title,
            URL:     "https://discord.gg/ACaXz4",
            IconURL: iconURL,
        },
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:  coordinates,
                Value: "[[Google Maps]](https://www.google.com/maps?" + coordinates + ") [[iOS Maps]](http://maps.apple.com/?address=" + coordinates + "&t=m)",
            },
        },
    }
}
