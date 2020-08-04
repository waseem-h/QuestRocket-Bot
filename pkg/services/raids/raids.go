package raids

import (
    "time"

    "github.com/adjust/rmq"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/cache"
    raidsCache "github.com/spiri2/Quests/pkg/cache/raids"
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
        logger: logrus.WithFields(logrus.Fields{
            "service": "raids",
        }),
        dbClient: dbClient,
        cache:    raidsCache.New(),
        queue:    queue,
        config:   config,
    }
}

func (service *Service) Name() string {
    return string(constants.Raid)
}

func (service *Service) GetCache() cache.Cache {
    return service.cache
}

func (service *Service) PopulateCache() error {
    raids, err := service.dbClient.GetAllRaids()
    if err != nil {
        return err
    }

    for _, raid := range raids {
        service.cache.Add(raid)
    }

    service.logger.Info("Total Items in cache: ", len(raids))

    return nil
}

func (service *Service) Load() error {
    urlMap := constants.URLMap()
    for name, url := range urlMap {
        raidsResponse := models.RaidsResponse{}
        err := service.httpClient.GetJSON(url+ "raids.php?time=0", &raidsResponse)
        if err != nil {
            return err
        }

        filteredRaids := make([]models.Raid, 0)

        for _, raid := range raidsResponse.Raids {
            expired, err := utils.CheckExpiry(raid.RaidEnd, time.Now())
            if err != nil {
                return err
            }

            if expired {
                continue
            }

            raid.ID = hash.SHA256(raid)
            raid.SiteName = name
            if !service.cache.Add(raid) {
                continue
            }
            filteredRaids = append(filteredRaids, raid)
        }

        if len(filteredRaids) > 0 { // Push documents

            service.logger.WithField("source", name).Info("Adding ", len(filteredRaids), " new raids to database")
            err = service.dbClient.AddRaids(filteredRaids...)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (service *Service) Publish() error {

    return nil
}


func (service *Service) CleanExpired() error {
    raids := convert.InterfaceSliceToRaidSlice(service.cache.GetAll())

    var idsToRemove []string
    for _, raid := range raids {
        expired, err := utils.CheckExpiry(raid.RaidEnd, time.Now()) // Check for expired raids by end time
        if err != nil {
            return err
        }

        if expired {
            service.cache.Remove(raid)
            idsToRemove = append(idsToRemove, raid.ID)
            continue
        }

        expired, err = utils.CheckExpiry(raid.RaidStart, time.Now()) // Check for expired raids which are started but not updated
        if err != nil {
            return err
        }

        if expired && raid.PokemonID == "0" {
            service.cache.Remove(raid)
            idsToRemove = append(idsToRemove, raid.ID)
        }
    }

    if len(idsToRemove) > 0 {
        service.logger.Info("Deleting ", len(idsToRemove), " expired raids from database")
        return service.dbClient.DeleteRaidsByID(idsToRemove...)
    }

    return nil
}

func (service *Service) Config() config.ServiceConfig {
    return service.config
}
