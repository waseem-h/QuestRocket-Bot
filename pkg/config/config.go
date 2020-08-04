package config

import "time"

type Config struct {
    Discord  Discord  `yaml:"discord"`
    Bot      Bot      `yaml:"bot"`
    Database Database `yaml:"database"`
    Queue    Queue    `yaml:"queue"`
}

type Queue struct {
    Redis Redis `yaml:"redis"`
}

type Redis struct {
    Address  string `yaml:"address"`
    Database int    `yaml:"database"`
}

type Database struct {
    Kind     string `yaml:"kind"`
    Address  string `yaml:"address"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type Bot struct {
    DefaultRefreshInterval time.Duration `yaml:"defaultRefreshInterval"`
    CommandPrefix          string        `yaml:"commandPrefix"`
    Invasion               ServiceConfig `yaml:"invasion"`
    SpecialQuest           ServiceConfig `yaml:"specialQuest"`
    PokemonQuest           ServiceConfig `yaml:"pokemonQuest"`
    StardustQuest          ServiceConfig `yaml:"stardustQuest"`
    ItemQuest              ServiceConfig `yaml:"itemQuest"`
}

type ServiceConfig struct {
    Enabled         bool           `yaml:"enabled"`
    QueueLimit      int            `yaml:"queueLimit"`
    ChannelID       string         `yaml:"channelID"`
    RefreshInterval *time.Duration `yaml:"refreshInterval,omitempty"`
}

type Discord struct {
    Token  string `yaml:"token"`
    Embeds Embeds `yaml:"embeds"`
}

type Embeds struct {
    Footer Footer `yaml:"footer"`
}

type Footer struct {
    IconURL string `yaml:"iconURL"`
    Text    string `yaml:"text"`
}
