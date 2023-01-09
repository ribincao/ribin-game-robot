package main

import (
	"ribin-game-robot/ui"

	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
)

func initLogger() {
	path := "./conf.yaml"
	config.ParseConf(path, config.GlobalConfig)
	logger.InitLogger(config.GlobalConfig.LogConfig)
}

func main() {
	initLogger()
	ui.RunMenu()
}
