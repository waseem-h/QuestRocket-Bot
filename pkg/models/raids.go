package models

type RaidsResponse struct {
    Raids    []Raid    `json:"raids"`
    Weathers []Weather `json:"weathers"`
    Meta     Metadata  `json:"meta"`
}

type Weather struct {
    CellID  string `json:"cell_id"`
    Weather string `json:"weather"`
}

type Raid struct {
    ID             string `bson:"_id" json:"_id,omitempty"`
    GymName        string `json:"gym_name"`
    CellID         string `json:"cell_id"`
    ExRaidEligible string `json:"ex_raid_eligible"`
    Sponsor        string `json:"sponsor"`
    Lat            string `json:"lat"`
    Lng            string `json:"lng"`
    RaidSpawn      string `json:"raid_spawn"`
    RaidStart      string `json:"raid_start"`
    RaidEnd        string `json:"raid_end"`
    PokemonID      string `json:"pokemon_id"`
    Level          string `json:"level"`
    Cp             string `json:"cp"`
    Team           string `json:"team"`
    Move1          string `json:"move1"`
    Move2          string `json:"move2"`
    IsExclusive    string `json:"is_exclusive"`
    Form           string `json:"form"`
    Gender         string `json:"gender"`
    SiteName       string `bson:"site_name" json:"site_name"`
}
