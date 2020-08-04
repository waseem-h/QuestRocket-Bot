package invasions

import (
    "github.com/spiri2/Quests/pkg/cache"
    "github.com/spiri2/Quests/pkg/convert"
    "github.com/spiri2/Quests/pkg/models"
)

type Cache struct {
    invasions []models.Invasion
}

var _ cache.Cache = (*Cache)(nil)

func New() *Cache {
    return &Cache{}
}

func (cache *Cache) Add(object interface{}) bool {
    invasion := object.(models.Invasion)
    if !cache.Contains(invasion) {
        cache.invasions = append(cache.invasions, invasion)
        return true
    }
    return false
}

func (cache *Cache) Contains(object interface{}) bool {
    for _, invasion := range cache.invasions {
        if invasion.ID == object.(models.Invasion).ID {
            return true
        }
    }
    return false
}

func (cache *Cache) GetAll() []interface{} {
    return convert.InvasionSliceToInterfaceSlice(cache.invasions)
}

func (cache *Cache) IndexOf(object interface{}) int {
    for index, invasion := range cache.invasions {
        if invasion.ID == object.(models.Invasion).ID {
            return index
        }
    }

    return -1
}

func (cache *Cache) Remove(object interface{}) bool {
    invasion := object.(models.Invasion)
    if cache.Contains(invasion) {
        index := cache.IndexOf(invasion)
        if index == -1 { // should never happen
            return false
        }

        cache.invasions[index] = cache.invasions[len(cache.invasions)-1]
        cache.invasions = cache.invasions[:len(cache.invasions)-1]

        return true
    }

    return false
}

func (cache *Cache) Filter(filterFunc cache.FilterFunc, compareTo interface{}) []interface{} {
    var filtered []interface{}

    for _, invasion := range cache.invasions {
        if filterFunc(invasion, compareTo) {
            filtered = append(filtered, invasion)
        }
    }

    return filtered
}
