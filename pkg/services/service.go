package services

import (
    "github.com/spiri2/Quests/pkg/cache"
    "github.com/spiri2/Quests/pkg/config"
)

type Services []Service

type Service interface {
    Load() error
    CleanExpired() error
    Name() string
    Publish() error
    Config() config.ServiceConfig
    cache.Caching
}

func (s Services) Lookup(name string) Service {
    for _, service := range s {
        if service.Name() == name {
            return service
        }
    }
    return nil
}
