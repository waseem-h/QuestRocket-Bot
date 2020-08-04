package database

import (
    "github.com/spiri2/Quests/pkg/constants"
    "github.com/spiri2/Quests/pkg/models"
)

type Client interface {
    Init() error
    Stop() error
    AddInvasions(invasions ...models.Invasion) error
    GetAllInvasions() ([]models.Invasion, error)
    DeleteInvasionsByID(ids ...string) error
    AddQuests(quests ...models.Quest) error
    GetAllQuests() ([]models.Quest, error)
    GetQuestsByType(questType constants.ServiceType) ([]models.Quest, error)
    DeleteQuestsByID(ids ...string) error
    AddExpiredQuests(quests ...models.Quest) error
    GetAllExpiredQuests() ([]models.Quest, error)
    GetExpiredQuestsByType(questType constants.ServiceType) ([]models.Quest, error)
    DeleteExpiredQuestsByID(ids ...string) error
    AddRaids(raids ...models.Raid) error
    GetAllRaids() ([]models.Raid, error)
    DeleteRaidsByID(ids ...string) error
}
