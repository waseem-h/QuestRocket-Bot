package raids

import (
    "github.com/spiri2/Quests/pkg/cache"
    "github.com/spiri2/Quests/pkg/convert"
    "github.com/spiri2/Quests/pkg/models"
)

type Cache struct {
    raids []models.Raid
}

var _ cache.Cache = (*Cache)(nil)

func New() *Cache {
    return &Cache{}
}

func (cache *Cache) Add(object interface{}) bool {
    raid := object.(models.Raid)
    if !cache.Contains(raid) {
        cache.raids = append(cache.raids, raid)
        return true
    }
    return false
}

func (cache *Cache) Contains(object interface{}) bool {
    for _, raid := range cache.raids {
        if raid.ID == object.(models.Raid).ID {
            return true
        }
    }

    return false
}

func (cache *Cache) GetAll() []interface{} {
    return convert.RaidSliceToInterfaceSlice(cache.raids)
}

func (cache *Cache) IndexOf(object interface{}) int {
    for index, raid := range cache.raids {
        if raid.ID == object.(models.Raid).ID {
            return index
        }
    }

    return -1
}

func (cache *Cache) Remove(object interface{}) bool {
    raid := object.(models.Raid)
    if cache.Contains(raid) {
        index := cache.IndexOf(raid)
        if index == -1 { // should never happen
            return false
        }

        cache.raids[index] = cache.raids[len(cache.raids)-1]
        cache.raids = cache.raids[:len(cache.raids)-1]

        return true
    }

    return false
}

func (cache *Cache) Filter(filterFunc cache.FilterFunc, compareTo interface{}) []interface{} {
    var filtered []interface{}

    for _, raid := range cache.raids {
        if filterFunc(raid, compareTo) {
            filtered = append(filtered, raid)
        }
    }

    return filtered
}
