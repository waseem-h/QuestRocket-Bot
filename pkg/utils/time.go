package utils

import (
    "strconv"
    "time"

    "github.com/nleeper/goment"
    "gopkg.in/ugjka/go-tz.v2/tz"
)

func CalculateExpirationTime(lat string, lng string) (string, error) {
    timezone, err := tz.GetZone(getPoint(lat, lng))
    if err != nil {
        return "", err
    }
    location, err := time.LoadLocation(timezone[0])
    if err != nil {
        return "", err
    }
    remoteTime, err := goment.New(time.Now().In(location))
    if err != nil {
        return "", err
    }
    return strconv.Itoa(int(remoteTime.EndOf("day").ToUnix())), err
}

func getPoint(lat string, lng string) tz.Point {
    latFloat, _ := strconv.ParseFloat(lat, 64)
    lngFloat, _ := strconv.ParseFloat(lng, 64)

    return tz.Point{
        Lat: latFloat,
        Lon: lngFloat,
    }
}
