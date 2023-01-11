package robot

import (
	"sync"

	"github.com/ribincao/ribin-game-server/logger"
	"go.uber.org/zap"
)

var robotMng sync.Map

func AddRobot(robotId string) string {
	_, ok := robotMng.Load(robotId)
	if ok {
		logger.Error("RobotId Repeat", zap.String("RobotId", robotId))
		return "[ERROR]RobotId Repeat"
	}
	robot := NewRobot(robotId)
	robotMng.Store(robotId, robot)
	return "[INFO]Robot Add Success."
}

type Robot struct {
	Id       string
	wrapconn *WrapConnection
}

func NewRobot(robotId string) *Robot {
	return &Robot{
		Id: robotId,
	}
}

func (r *Robot) EnterRoom(roomId string) error {
	wrapConn, err := DialWrapConn(r.Id, roomId)
	r.wrapconn = wrapConn
	return err
}
