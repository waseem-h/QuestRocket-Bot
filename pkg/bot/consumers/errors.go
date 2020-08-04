package consumers

import (
    "github.com/adjust/rmq"
    "github.com/sirupsen/logrus"
)

func rejectOnError(err error, delivery rmq.Delivery) bool {
    if err != nil {
        logrus.Error(err)
        delivery.Reject()
        return true
    }
    return false
}
