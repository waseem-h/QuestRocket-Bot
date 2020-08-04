package invasions

import (
    "encoding/json"
    "time"

    "github.com/adjust/rmq"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/cache"
    invasionsCache "github.com/spiri2/Quests/pkg/cache/invasions"
    "github.com/spiri2/Quests/pkg/cache/invasions/filters"
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
    httpClient *httpclient.Client
    logger     *logrus.Entry
    dbClient   database.Client
    cache      cache.Cache
    queue      rmq.Queue
    config     config.ServiceConfig
}

var _ services.Service = (*Service)(nil)

func New(dbClient database.Client, queue rmq.Queue, config config.ServiceConfig) *Service {
    return &Service{
        httpClient: httpclient.New(),
        logger:     logrus.WithField("service", "invasions"),
        dbClient:   dbClient,
        cache:      invasionsCache.New(),
        queue:      queue,
        config:     config,
    }
}

func (service *Service) Name() string {
    return string(constants.Invasion)
}

func (service *Service) PopulateCache() error {
    invasions, err := service.dbClient.GetAllInvasions()
    if err != nil {
        return err
    }

    for _, invasion := range invasions {
        service.cache.Add(invasion)
    }

    service.logger.Info("Total Items in cache: ", len(invasions))

    return nil
}

func (service *Service) GetCache() cache.Cache {
    return service.cache
}

func (service *Service) Load() error {
    urlMap := constants.URLMap()
    for name, url := range urlMap {
        invasionsResponse := models.InvasionsResponse{}
        err := service.httpClient.GetJSON(url+"pokestop.php?time=0", &invasionsResponse)
        if err != nil {
            return err
        }

        filteredInvasions := make([]models.Invasion, 0)

        for _, invasion := range invasionsResponse.Invasions {
            expired, err := utils.CheckExpiry(invasion.InvasionEnd, time.Now())
            if err != nil {
                return err
            }

            if expired {
                continue
            }

            invasion.ID = hash.SHA256(invasion)
            invasion.SiteName = name          // Do not include in hash
            if !service.cache.Add(invasion) { // Already present in cache
                continue
            }
            filteredInvasions = append(filteredInvasions, invasion)
        }

        if len(filteredInvasions) > 0 { // Push documents

            service.logger.WithField("source", name).Info("Adding ", len(filteredInvasions), " new invasions to database")
            err = service.dbClient.AddInvasions(filteredInvasions...)
            if err != nil {
                return err
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
        invasions := convert.InterfaceSliceToInvasionSlice(service.cache.Filter(filters.BySiteName, name))
        if len(invasions) < queuePerSite {
            amountToQueue = len(invasions)
        }

        chosenIndexes := utils.GetNRandom(0, len(invasions), amountToQueue)
        for _, index := range chosenIndexes {
            if err := service.pushToQueue(invasions[index]); err != nil {
                return err
            }
        }
    }

    return nil
}

func (service *Service) pushToQueue(invasion models.Invasion) error {
    invasionBytes, err := json.Marshal(invasion)
    if err != nil {
        return err
    }

    service.queue.PublishBytes(invasionBytes)
    return nil
}

func (service *Service) CleanExpired() error {

    invasions := convert.InterfaceSliceToInvasionSlice(service.cache.GetAll())

    var idsToRemove []string
    for _, invasion := range invasions {
        expired, err := utils.CheckExpiry(invasion.InvasionEnd, time.Now())
        if err != nil {
            return err
        }

        if expired {
            service.cache.Remove(invasion)
            idsToRemove = append(idsToRemove, invasion.ID)
        }
    }

    if len(idsToRemove) > 0 {
        service.logger.Info("Deleting ", len(idsToRemove), " expired invasions from database")
        return service.dbClient.DeleteInvasionsByID(idsToRemove...)
    }

    return nil
}

func (service *Service) Config() config.ServiceConfig {
    return service.config
}
