package models

import "github.com/spiri2/Quests/pkg/constants"

type QuestsResponse struct {
    Quests  []Quest  `json:"quests"`
    Meta    Metadata `json:"meta"`
    Filters Filters  `json:"filters"`
}

type Filters struct {
    T3 []string `json:"t3"`
    T2 []string `json:"t2"`
    T7 []string `json:"t7"`
}

type Quest struct {
    ID               string                `bson:"_id" json:"_id,omitempty"`
    Name             string                `bson:"name" json:"name"`
    Lat              string                `bson:"lat" json:"lat"`
    Lng              string                `bson:"lng" json:"lng"`
    RewardsString    string                `bson:"rewards_string" json:"rewards_string"`
    ConditionsString string                `bson:"conditions_string" json:"conditions_string"`
    Image            string                `bson:"image" json:"image"`
    RewardsTypes     string                `bson:"rewards_types" json:"rewards_types"`
    RewardsAmounts   string                `bson:"rewards_amounts" json:"rewards_amounts"`
    RewardsIds       string                `bson:"rewards_ids" json:"rewards_ids"`
    Type             constants.ServiceType `bson:"type" json:"type,omitempty"`
    Expiration       string                `bson:"expiration" json:"expiration,omitempty"`
    SiteName         string                `bson:"site_name" json:"site_name"`
}
