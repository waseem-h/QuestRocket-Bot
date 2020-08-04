package filters

import "github.com/spiri2/Quests/pkg/models"

func BySiteName(object interface{}, compareTo interface{}) bool {
    return object.(models.Quest).SiteName == compareTo.(string)
}
