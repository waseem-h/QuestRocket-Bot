package quests

import (
    "github.com/spiri2/Quests/pkg/cache"
    "github.com/spiri2/Quests/pkg/convert"
    "github.com/spiri2/Quests/pkg/models"
)

type Cache struct {
    quests []models.Quest
}

var _ cache.Cache = (*Cache)(nil)

func New() *Cache {
    return &Cache{}
}

func (cache *Cache) Add(object interface{}) bool {
    quest := object.(models.Quest)
    if !cache.Contains(quest) {
        cache.quests = append(cache.quests, quest)
        return true
    }
    return false
}

func (cache *Cache) Contains(object interface{}) bool {
    for _, quest := range cache.quests {
        if quest.ID == object.(models.Quest).ID {
            return true
        }
    }
    return false
}

func (cache *Cache) GetAll() []interface{} {
    return convert.QuestSliceToInterfaceSlice(cache.quests)
}

func (cache *Cache) IndexOf(object interface{}) int {
    for index, quest := range cache.quests {
        if quest.ID == object.(models.Quest).ID {
            return index
        }
    }

    return -1
}

func (cache *Cache) Remove(object interface{}) bool {
    quest := object.(models.Quest)
    if cache.Contains(quest) {
        index := cache.IndexOf(quest)
        if index == -1 { // should never happen
            return false
        }

        cache.quests[index] = cache.quests[len(cache.quests)-1]
        cache.quests = cache.quests[:len(cache.quests)-1]

        return true
    }

    return false
}

func (cache *Cache) Filter(filterFunc cache.FilterFunc, compareTo interface{}) []interface{} {
    var filtered []interface{}

    for _, quest := range cache.quests {
        if filterFunc(quest, compareTo) {
            filtered = append(filtered, quest)
        }
    }

    return filtered
}
