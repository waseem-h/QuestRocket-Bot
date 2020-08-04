package bot

import (
    "errors"
    "sync"

    "github.com/adjust/rmq"
    "github.com/bwmarrin/discordgo"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/bot/consumers"
    "github.com/spiri2/Quests/pkg/bot/events"
    "github.com/spiri2/Quests/pkg/config"
    "github.com/spiri2/Quests/pkg/constants"
    "github.com/spiri2/Quests/pkg/database"
    "github.com/spiri2/Quests/pkg/database/mongodb"
    "github.com/spiri2/Quests/pkg/services"
    "github.com/spiri2/Quests/pkg/services/invasions"
    "github.com/spiri2/Quests/pkg/services/quests"
    "github.com/spiri2/Quests/pkg/services/raids"
    "github.com/spiri2/Quests/pkg/timer"
    "github.com/spiri2/Quests/pkg/utils"
)

type Bot struct {
    Session         *discordgo.Session
    Services        []services.Service
    Logger          *logrus.Logger
    DatabaseClient  database.Client
    QueueConnection rmq.Connection
    Queues          map[constants.ServiceType]rmq.Queue
    Consumers       map[constants.ServiceType]rmq.Consumer
    Config          *config.Config
    ServiceQueue    rmq.Queue
}

func New(config *config.Config) (*Bot, error) {

    session, err := discordgo.New("Bot " + config.Discord.Token)
    if err != nil {
        return nil, err
    }

    session.AddHandler(func(session *discordgo.Session, event *discordgo.Ready) {
        events.Ready(session, event, config)
    })
    session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageCreate) {
        events.Message(session, event, config)
    })

    var dbClient database.Client

    databaseKind := config.Database.Kind
    address := config.Database.Address
    username := config.Database.Username
    password := config.Database.Password
    name := config.Database.Name

    switch databaseKind {
    case "mongodb":
        dbClient, err = mongodb.New(address, username, password, name)
        if err != nil {
            return nil, err
        }
    default:
        return nil, errors.New("Unsupported database type specified: " + databaseKind)
    }

    connection := rmq.OpenConnection("quests", "tcp", config.Queue.Redis.Address, config.Queue.Redis.Database)

    queuesMap := make(map[constants.ServiceType]rmq.Queue)
    consumersMap := make(map[constants.ServiceType]rmq.Consumer)
    var enabledServices services.Services

    if config.Bot.Invasion.Enabled {
        invasionsQueue := utils.SetupQueue(connection, constants.Invasion)
        invasionConsumer := consumers.NewInvasionConsumer(session, config)
        invasionsQueue.AddConsumer("invasions consumer", invasionConsumer)
        queuesMap[constants.Invasion] = invasionsQueue
        consumersMap[constants.Invasion] = invasionConsumer

        enabledServices = append(enabledServices, invasions.New(dbClient, invasionsQueue, config.Bot.Invasion))
    }

    if config.Bot.SpecialQuest.Enabled {
        specialQuestsQueue := utils.SetupQueue(connection, constants.QuestSpecial)
        specialQuestConsumer := consumers.NewSpecialQuestConsumer(session, config)
        specialQuestsQueue.AddConsumer("special quests consumer", specialQuestConsumer)
        queuesMap[constants.QuestSpecial] = specialQuestsQueue
        consumersMap[constants.QuestSpecial] = specialQuestConsumer

        enabledServices = append(enabledServices, quests.New(dbClient, specialQuestsQueue, config.Bot.SpecialQuest, []int{327}, constants.QuestSpecial))
    }

    if config.Bot.PokemonQuest.Enabled {
        pokemonQuestsQueue := utils.SetupQueue(connection, constants.QuestPokemon)
        pokemonQuestConsumer := consumers.NewPokemonQuestConsumer(session, config)
        pokemonQuestsQueue.AddConsumer("pokemon quests consumer", pokemonQuestConsumer)
        queuesMap[constants.QuestPokemon] = pokemonQuestsQueue
        consumersMap[constants.QuestPokemon] = pokemonQuestConsumer

        enabledServices = append(enabledServices, quests.New(dbClient, pokemonQuestsQueue, config.Bot.PokemonQuest, utils.MakeRange(1, 649), constants.QuestPokemon))
    }

    if config.Bot.StardustQuest.Enabled {
        stardustQuestsQueue := utils.SetupQueue(connection, constants.QuestStardust)
        stardustQuestConsumer := consumers.NewStardustQuestConsumer(session, config)
        stardustQuestsQueue.AddConsumer("stardust quests consumer", stardustQuestConsumer)
        queuesMap[constants.QuestStardust] = stardustQuestsQueue
        consumersMap[constants.QuestStardust] = stardustQuestConsumer

        enabledServices = append(enabledServices, quests.New(dbClient, stardustQuestsQueue, config.Bot.StardustQuest, []int{200, 500, 1000, 1500}, constants.QuestStardust))
    }

    if config.Bot.ItemQuest.Enabled {
        itemQuestsQueue := utils.SetupQueue(connection, constants.QuestItem)
        itemQuestConsumer := consumers.NewItemQuestConsumer(session, config)
        itemQuestsQueue.AddConsumer("item quests consumer", itemQuestConsumer)
        queuesMap[constants.QuestItem] = itemQuestsQueue
        consumersMap[constants.QuestItem] = itemQuestConsumer

        enabledServices = append(enabledServices, quests.New(dbClient, itemQuestsQueue, config.Bot.ItemQuest, []int{1, 2, 3, 101, 102, 103, 104, 201, 202, 701, 702, 703, 705, 706, 1301, 1103}, constants.QuestItem))
    }

    enabledServices = append(enabledServices, raids.New(dbClient, nil, config.Bot.ItemQuest))

    serviceQueue := utils.SetupQueue(connection, "service")
    serviceConsumer := consumers.NewServiceConsumer(enabledServices, config)
    serviceQueue.AddConsumer("services consumer", serviceConsumer)


    return &Bot{
        Session:         session,
        Services:        enabledServices,
        Logger:          logrus.New(),
        DatabaseClient:  dbClient,
        QueueConnection: connection,
        Queues:          queuesMap,
        Consumers:       consumersMap,
        Config:          config,
        ServiceQueue:    serviceQueue,
    }, nil
}

func (bot *Bot) Init() error {
    err := bot.Session.Open()
    if err != nil {
        return err
    }

    err = bot.DatabaseClient.Init()
    if err != nil {
        return nil
    }

    fatalErrors := make(chan error)
    wgDone := make(chan bool)

    var wg sync.WaitGroup
    for _, service := range bot.Services {
        wg.Add(1)
        go func(s services.Service) {
            defer wg.Done()
            if err := s.PopulateCache(); err != nil {
                fatalErrors <- err
            }
        }(service)
    }

    go func() {
        wg.Wait()
        close(wgDone)
    }()

    select {
    case <-wgDone:
        break
    case err := <-fatalErrors:
        close(fatalErrors)
        return err
    }

    return nil
}

func (bot *Bot) Run() {
    timers := make([]*timer.SecondsTimer, len(bot.Services))
    for {
        for index, service := range bot.Services {
            if timers[index] != nil && timers[index].TimeRemaining() >= 0 {
                continue
            }

            bot.ServiceQueue.Publish(service.Name())

            refreshInterval := bot.Config.Bot.DefaultRefreshInterval
            if service.Config().RefreshInterval != nil {
                refreshInterval = *service.Config().RefreshInterval
            }

            timers[index] = timer.NewSecondsTimer(refreshInterval)
        }
    }
}

func (bot *Bot) Stop() error {
    err := bot.Session.Close()
    if err != nil {
        return err
    }

    return bot.DatabaseClient.Stop()
}
