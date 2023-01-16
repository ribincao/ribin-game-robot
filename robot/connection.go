package robot

import (
	"fmt"
	"ribin-game-robot/utils"
	"time"

	"go.uber.org/atomic"

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

func DialWrapConn(playerId string, roomId string) *WrapConnection {
	wrapConn := NewWrapConnection(playerId, roomId)
	wrapConn.roomConn, _ = DialRoomConn("localhost", 8080)
	return wrapConn
}

type WrapConnection struct {
	playerId   string
	roomId     string
	roomConn   *websocket.Conn
	isClose    atomic.Bool
	SeqCounter *atomic.Int32
}

func NewWrapConnection(playId string, roomId string) *WrapConnection {
	return &WrapConnection{
		playerId:   playId,
		roomId:     roomId,
		SeqCounter: atomic.NewInt32(0),
	}
}

func (wc *WrapConnection) SendMessage(req *base.Client2ServerReq) error {
	client2serverReq, _ := MarshalAndEncode(req)
	return wc.roomConn.WriteMessage(websocket.BinaryMessage, client2serverReq)
}

func (wc *WrapConnection) GetSeq() string {
	return fmt.Sprintf("%s-%d", wc.playerId, wc.SeqCounter.Add(1))
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
			seq := wc.GetSeq()
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
