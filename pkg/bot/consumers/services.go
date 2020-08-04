package consumers

import (
    "github.com/adjust/rmq"
    "github.com/sirupsen/logrus"

    "github.com/spiri2/Quests/pkg/config"
    "github.com/spiri2/Quests/pkg/services"
)

type ServiceConsumer struct {
    config   *config.Config
    services services.Services
}

func NewServiceConsumer(services services.Services, config *config.Config) *ServiceConsumer {
    return &ServiceConsumer{
        config:   config,
        services: services,
    }
}

func (consumer ServiceConsumer) Consume(delivery rmq.Delivery) {
    service := consumer.services.Lookup(delivery.Payload())
    if service == nil {
        logrus.Error("Cannot lookup service with name: ", delivery.Payload())
        delivery.Reject()
        return
    }

    if err := service.Load(); err != nil {
        logrus.Error(err)
    }

    if err := service.CleanExpired(); err != nil {
        logrus.Error(err)
    }

    if err := service.Publish(); err != nil {
        logrus.Error(err)
    }
}
