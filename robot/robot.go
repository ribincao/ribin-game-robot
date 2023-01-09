package robot

import (
	"sync"

	"github.com/ribincao/ribin-game-server/logger"
	"go.uber.org/zap"
)

var robotMng sync.Map

func AddRobot(robotId string, robot *Robot) string {
	_, ok := robotMng.Load(robotId)
	if ok {
		logger.Error("RobotId Repeat", zap.String("RobotId", robotId))
		return "[ERROR]RobotId Repeat"
	}
	robotMng.Store(robotId, robot)
	return "[INFO]Robot Add Success."
}

type Robot struct {
	Id   string
	Name string
}

func (r *Robot) EnterRoom(roomId string) {

}
