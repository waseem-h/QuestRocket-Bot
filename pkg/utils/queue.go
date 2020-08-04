package utils

import (
    "time"

    "github.com/adjust/rmq"

    "github.com/spiri2/Quests/pkg/constants"
)

func SetupQueue(connection rmq.Connection, name constants.ServiceType) rmq.Queue {
    queue := connection.OpenQueue(string(name))
    queue.StartConsuming(1, time.Second)
    return queue
}
