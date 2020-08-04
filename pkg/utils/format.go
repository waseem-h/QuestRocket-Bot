package utils

import (
    "fmt"
    "strconv"
    "time"

    "github.com/sirupsen/logrus"
)

func FormatToRemainingTime(nanoTime string) string {
    expireTimeInt, err := strconv.ParseInt(nanoTime, 10, 64)
    if err != nil {
        logrus.Error(err)
        return ""
    }
    diffSeconds := int(expireTimeInt - time.Now().UnixNano()/1000000000)

    hours := diffSeconds / 3600
    minutes := (diffSeconds % 3600) / 60
    seconds := diffSeconds % 60

    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
