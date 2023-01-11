package robot

import (
	"ribin-game-robot/utils"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ribincao/ribin-game-server/codec"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-protocol/base"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func MarshalAndEncode(input protoreflect.ProtoMessage) ([]byte, error) {
	inputJSON, err := proto.Marshal(input)
	if err != nil {
		return nil, err
	}
	return codec.DefaultCodec.Encode(inputJSON, codec.RPC)
}

func DecodeAndUnmarshal(input []byte) (*base.Server2ClientRsp, error) {
	rsp := new(base.Server2ClientRsp)
	decode, _ := codec.DefaultCodec.Decode(input)
	err := proto.Unmarshal(decode.Data, rsp)
	return rsp, err
}

func DialRoomConn(Ip string, Port int32) (*websocket.Conn, error) {
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

func DialWrapConn(playerId string, roomId string) (*WrapConnection, error) {
	wrapConn := NewWrapConnection(playerId, roomId)
	wrapConn.roomConn, _ = DialRoomConn("localhost", 8080)
	err := wrapConn.EnterRoom()
	return wrapConn, err
}

type WrapConnection struct {
	playerId string
	roomId   string
	roomConn *websocket.Conn
	isClose  *atomic.Bool
}

func NewWrapConnection(playId string, roomId string) *WrapConnection {
	return &WrapConnection{
		playerId: playId,
		roomId:   roomId,
	}
}

func (wc *WrapConnection) EnterRoom() error {
	enterRoomReq := &base.Client2ServerReq{
		Cmd: base.Client2ServerReqCmd_E_CMD_ROOM_ENTER,
		Seq: "TEST",
		Body: &base.ReqBody{
			PlayerId:     wc.playerId,
			RoomId:       wc.roomId,
			EnterRoomReq: &base.EnterRoomReq{},
		},
	}
	req, _ := MarshalAndEncode(enterRoomReq)
	// go wc.RoomHeartBeat()
	// go wc.ReadMessage()
	return wc.roomConn.WriteMessage(websocket.BinaryMessage, req)
}

func (wc *WrapConnection) ReadMessage() {
	ticker := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-ticker.C:
			_, p, err := wc.roomConn.ReadMessage()
			if err != nil {
				logger.Error("ReadMessageError", zap.Error(err))
				return
			}
			rsp, err := DecodeAndUnmarshal(p)
			logger.Debug("Rsp", zap.Any("Rsp", rsp), zap.Error(err))
		default:
			if wc.isClose.Load() {
				wc.roomConn.Close()
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (wc *WrapConnection) RoomHeartBeat() {
	ticker := time.NewTicker(time.Millisecond * 2000)
	for {
		select {
		case <-ticker.C:
			seq := "Test"
			heartbeatReq := &base.Client2ServerReq{
				Cmd: base.Client2ServerReqCmd_E_CMD_HEART_BEAT,
				Seq: seq,
				Body: &base.ReqBody{
					RoomId:   wc.roomId,
					PlayerId: wc.playerId,
				},
			}
			req, _ := MarshalAndEncode(heartbeatReq)
			wc.roomConn.WriteMessage(websocket.BinaryMessage, req)
		default:
			if wc.isClose.Load() {
				wc.roomConn.Close()
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}
