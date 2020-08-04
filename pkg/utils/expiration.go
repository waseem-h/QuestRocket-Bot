package utils

import (
    "strconv"
    "time"
)

func CheckExpiry(expireTime string, from time.Time) (bool, error) {
    expireTimeInt, err := strconv.ParseInt(expireTime, 10, 64)
    if err != nil {
        return false, err
    }
    diffSeconds := int(expireTimeInt - from.UnixNano()/1000000000)
    if diffSeconds <= 0 {
        return true, nil
    }

    return false, nil
}
