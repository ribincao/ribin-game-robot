package robot

import (
	"ribin-game-robot/utils"

	"github.com/gorilla/websocket"
	"github.com/ribincao/ribin-game-server/logger"
	"go.uber.org/zap"
)

func DialWebsocketConn(Ip string, Port int32) (*websocket.Conn, error) {
	url := utils.GeneWebsocketURL(Ip, Port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Error("DialRoomConnError",
			zap.String("IP", Ip),
			zap.Any("Port", Port),
			zap.Error(err))
		return nil, err
	}
	return conn, nil
}
