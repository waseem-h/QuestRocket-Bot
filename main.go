package main

import (
    "fmt"
    "math/rand"
    "os"
    "os/signal"
    "runtime"
    "strings"
    "syscall"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"

    "github.com/spiri2/Quests/pkg/bot"
    "github.com/spiri2/Quests/pkg/config"
)

func init() {
    randomSeed()
    configureLogging()
}

func randomSeed() {
    rand.Seed(time.Now().UnixNano())
}

func fetchConfig() (*config.Config, error) {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    conf := &config.Config{}
    err := viper.Unmarshal(conf)
    if err != nil {
        return nil, err
    }

    return conf, err
}

func configureLogging() {
    logrus.SetReportCaller(true)
    formatter := &logrus.TextFormatter{
        TimestampFormat:        "02-01-2006 15:04:05",
        FullTimestamp:          true,
        DisableLevelTruncation: true,
        CallerPrettyfier: func(f *runtime.Frame) (string, string) {
            return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
        },
    }
    logrus.SetFormatter(formatter)
}

func formatFilePath(path string) string {
    arr := strings.Split(path, "/")
    return arr[len(arr)-1]
}

func main() {
    conf, err := fetchConfig()
    if err != nil {
        logrus.Fatal(err)
    }

    questBot, err := bot.New(conf)
    if err != nil {
        logrus.Fatal(err)
    }

    err = questBot.Init()
    if err != nil {
        logrus.Fatal(err)
    }

    logrus.Println("QuestRocket-Bot is now running")

    go questBot.Run()

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    err = questBot.Stop()
    if err != nil {
        logrus.Fatal(err)
    }
}
