package quests

import (
    "encoding/json"
    "strconv"
    "strings"
    "time"

    "github.com/adjust/rmq"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/cache"
    questsCache "github.com/spiri2/Quests/pkg/cache/quests"
    "github.com/spiri2/Quests/pkg/cache/quests/filters"
    "github.com/spiri2/Quests/pkg/config"
    "github.com/spiri2/Quests/pkg/constants"
    "github.com/spiri2/Quests/pkg/convert"
    "github.com/spiri2/Quests/pkg/database"
    "github.com/spiri2/Quests/pkg/hash"
    "github.com/spiri2/Quests/pkg/httpclient"
    "github.com/spiri2/Quests/pkg/models"
    "github.com/spiri2/Quests/pkg/services"
    "github.com/spiri2/Quests/pkg/utils"
)

type Service struct {
    httpClient   *httpclient.Client
    logger       *logrus.Entry
    dbClient     database.Client
    cache        cache.Cache
    expiredCache cache.Cache
    questIDs     []int
    questType    constants.ServiceType
    queue        rmq.Queue
    config       config.ServiceConfig
}

var _ services.Service = (*Service)(nil)

func New(dbClient database.Client, queue rmq.Queue, config config.ServiceConfig, questIDs []int, questType constants.ServiceType) *Service {
    return &Service{
        httpClient: httpclient.New(),
        logger: logrus.WithFields(logrus.Fields{
            "service":   "quests",
            "questType": questType,
        }),
        dbClient:     dbClient,
        questIDs:     questIDs,
        cache:        questsCache.New(),
        expiredCache: questsCache.New(),
        questType:    questType,
        queue:        queue,
        config:       config,
    }
}

func (service *Service) Name() string {
    return string(service.questType)
}

func (service *Service) GetCache() cache.Cache {
    return service.cache
}

func (service *Service) PopulateCache() error {

    quests, err := service.dbClient.GetQuestsByType(service.questType)
    if err != nil {
        return err
    }

    for _, quest := range quests {
        service.cache.Add(quest)
    }

    service.logger.Info("Total Items in cache: ", len(quests))

    expiredQuests, err := service.dbClient.GetExpiredQuestsByType(service.questType)
    if err != nil {
        return err
    }

    for _, expiredQuest := range expiredQuests {
        service.expiredCache.Add(expiredQuest)
    }

    service.logger.Info("Total Items in expired cache: ", len(expiredQuests))

    return nil
}

func (service *Service) createQueryString(ids []int) string {
    var objects []string
    var prefix string
    postfix := "&time=0"

    switch service.questType {
    case constants.QuestItem:
        prefix = "quests[]=7,0,290"
    }

    for _, id := range ids {
        switch service.questType {
        case constants.QuestPokemon, constants.QuestSpecial:
            objects = append(objects, "quests[]=7,0,"+strconv.Itoa(id))
        case constants.QuestStardust:
            objects = append(objects, "quests[]=3,"+strconv.Itoa(id)+",0")
        case constants.QuestItem:
            objects = append(objects, "quests[]=2,0,"+strconv.Itoa(id))
        }
    }

    return prefix + strings.Join(objects, "&") + postfix
}

func (service *Service) Load() error {
    urlMap := constants.URLMap()
    for name, url := range urlMap {

        for _, chunk := range utils.ChunkInt(service.questIDs, 150) {
            questsResponse := models.QuestsResponse{}
            requestURL := url + "quests.php?" + service.createQueryString(chunk)
            err := service.httpClient.GetJSON(requestURL, &questsResponse)
            if err != nil {
                return err
            }

            var filteredQuests []models.Quest

            for _, quest := range questsResponse.Quests {
                estimatedExpiration, err := utils.CalculateExpirationTime(quest.Lat, quest.Lng)
                if err != nil {
                    return err
                }

                expired, err := utils.CheckExpiry(estimatedExpiration, time.Now())
                if err != nil {
                    return err
                }

                if expired {
                    continue
                }

                quest.Type = service.questType
                quest.ID = hash.SHA256(quest)
                // Add expiration after computing the Hash
                quest.Expiration = estimatedExpiration
                quest.SiteName = name

                if service.expiredCache.Contains(quest) || !service.cache.Add(quest) { // Already present in cache or expired cache
                    continue
                }

                filteredQuests = append(filteredQuests, quest)
            }

            if len(filteredQuests) > 0 { // Push documents
                service.logger.WithField("source", name).Info("Adding ", len(filteredQuests), " new quests to database")
                err = service.dbClient.AddQuests(filteredQuests...)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func (service *Service) Publish() error {
    urlMap := constants.URLMap()
    queuePerSite := service.config.QueueLimit / len(urlMap)
    for name := range urlMap {
        amountToQueue := queuePerSite
        quests := convert.InterfaceSliceToQuestSlice(service.cache.Filter(filters.BySiteName, name))
        if len(quests) < queuePerSite {
            amountToQueue = len(quests)
        }

        chosenIndexes := utils.GetNRandom(0, len(quests), amountToQueue)
        for _, index := range chosenIndexes {
            if err := service.pushToQueue(quests[index]); err != nil {
                return err
            }
        }
    }

    return nil
}

func (service *Service) pushToQueue(quest models.Quest) error {
    invasionBytes, err := json.Marshal(quest)
    if err != nil {
        return err
    }

    service.queue.PublishBytes(invasionBytes)
    return nil
}

func (service *Service) CleanExpired() error {
    quests := convert.InterfaceSliceToQuestSlice(service.cache.GetAll())
    expiredQuests := convert.InterfaceSliceToQuestSlice(service.expiredCache.GetAll())

    var questsToRemove []string
    var expiredQuestsToAdd []models.Quest
    for _, quest := range quests {
        expired, err := utils.CheckExpiry(quest.Expiration, time.Now())
        if err != nil {
            return err
        }

        if expired {
            service.cache.Remove(quest)
            service.expiredCache.Add(quest)

            questsToRemove = append(questsToRemove, quest.ID)
            expiredQuestsToAdd = append(expiredQuestsToAdd, quest)
        }
    }

    var expiredQuestsToRemove []string
    for _, expiredQuest := range expiredQuests {
        expired, err := utils.CheckExpiry(expiredQuest.Expiration, time.Now().Add(-48*time.Hour))
        if err != nil {
            return err
        }

        if expired {
            service.expiredCache.Remove(expiredQuest)
            expiredQuestsToRemove = append(expiredQuestsToRemove, expiredQuest.ID)
        }
    }

    if len(questsToRemove) > 0 {
        service.logger.Info("Deleting ", len(questsToRemove), " quests from database")
        if err := service.dbClient.DeleteQuestsByID(questsToRemove...); err != nil {
            return err
        }

        if err := service.dbClient.AddExpiredQuests(expiredQuestsToAdd...); err != nil {
            return err
        }
    }

    if len(expiredQuestsToRemove) > 0 {
        service.logger.Info("Deleting ", len(expiredQuestsToRemove), " expired quests from database")
        return service.dbClient.DeleteExpiredQuestsByID(expiredQuestsToRemove...)
    }

    return nil
}

func (service *Service) Config() config.ServiceConfig {
    return service.config
}
