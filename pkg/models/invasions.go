package models

type InvasionsResponse struct {
    Invasions []Invasion `json:"invasions"`
    Meta      Metadata   `json:"meta"`
}

type Invasion struct {
    ID            string `bson:"_id" json:"_id,omitempty"`
    Name          string `bson:"name" json:"name"`
    Lat           string `bson:"lat" json:"lat"`
    Lng           string `bson:"lng" json:"lng"`
    InvasionStart string `bson:"invasion_start" json:"invasion_start"`
    InvasionEnd   string `bson:"invasion_end" json:"invasion_end"`
    Character     string `bson:"character" json:"character"`
    SiteName      string `bson:"site_name" json:"site_name"`
}
