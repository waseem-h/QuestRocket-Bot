package convert

import "github.com/spiri2/Quests/pkg/models"

func InvasionSliceToInterfaceSlice(invasions []models.Invasion) []interface{} {
    slice := make([]interface{}, len(invasions))
    for index, value := range invasions {
        slice[index] = value
    }

    return slice
}

func QuestSliceToInterfaceSlice(quests []models.Quest) []interface{} {
    slice := make([]interface{}, len(quests))
    for index, value := range quests {
        slice[index] = value
    }

    return slice
}

func RaidSliceToInterfaceSlice(raids []models.Raid) []interface{} {
    slice := make([]interface{}, len(raids))
    for index, value := range raids {
        slice[index] = value
    }

    return slice
}

func InterfaceSliceToInvasionSlice(slice []interface{}) []models.Invasion {
    invasionSlice := make([]models.Invasion, len(slice))
    for index, invasion := range slice {
        invasionSlice[index] = invasion.(models.Invasion)
    }

    return invasionSlice
}

func InterfaceSliceToQuestSlice(slice []interface{}) []models.Quest {
    questSlice := make([]models.Quest, len(slice))
    for index, quest := range slice {
        questSlice[index] = quest.(models.Quest)
    }

    return questSlice
}

func InterfaceSliceToRaidSlice(slice []interface{}) []models.Raid {
    raidSlice := make([]models.Raid, len(slice))
    for index, raid := range slice {
        raidSlice[index] = raid.(models.Raid)
    }

    return raidSlice
}
