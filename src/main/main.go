package main

import (
    "kdbp",
    "github.com/spf13/viper"
)

//setup configuration
viper.SetConfigName("config")
viper.AddConfigPath("./config/")
viper.ReadInConfig()

func main() {
    kdbp.RunBenchmark()
}
