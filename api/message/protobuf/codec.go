package protobuf

import (
	"errors"

	"github.com/eosswedenorg/thalos/api/message"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=$DST_DIR $SRC_DIR/addressbook.proto

func Encode(v any) ([]byte, error) {
	switch p := v.(type) {
	case message.HeartBeat:
		return EncodeHeartbeat(p)
	case message.ActionTrace:
		return EncodeActionTrace(p)
	}

	return nil, errors.New("invalid type")
}

func EncodeHeartbeat(hb message.HeartBeat) ([]byte, error) {
	return proto.Marshal(&Heartbeat{
		BlockNum:                 hb.BlockNum,
		HeadBlocknum:             hb.HeadBlockNum,
		LastIrreversibleBlocknum: hb.LastIrreversibleBlockNum,
	})
}

func EncodeActionTrace(act message.ActionTrace) ([]byte, error) {
	return proto.Marshal(&ActionTrace{
		TxId:     act.TxID,
		Name:     act.Name,
		Contract: act.Contract,
		Receiver: act.Receiver,
		Data:     act.Data,
		HexData:  act.HexData,
	})
}

func Decode(data []byte, v any) error {
	switch p := v.(type) {
	case *message.HeartBeat:
		return DecodeHeartbeat(data, p)
	case *message.ActionTrace:
		return DecodeActionTrace(data, p)
	}
	return errors.New("invalid type")
}

func DecodeHeartbeat(data []byte, hb *message.HeartBeat) error {
	msg := &Heartbeat{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}

	hb.BlockNum = msg.BlockNum
	hb.HeadBlockNum = msg.HeadBlocknum
	hb.LastIrreversibleBlockNum = msg.LastIrreversibleBlocknum

	return nil
}

func DecodeActionTrace(b []byte, act *message.ActionTrace) error {
	msg := &ActionTrace{}
	if err := proto.Unmarshal(b, msg); err != nil {
		return err
	}

	*act = message.ActionTrace{
		TxID:     msg.TxId,
		Name:     msg.Name,
		Contract: msg.Contract,
		Receiver: msg.Receiver,
		Data:     msg.Data,
		HexData:  msg.HexData,
	}

	return nil
}

func init() {
	message.RegisterCodec("protobuf", message.Codec{
		Encoder: Encode,
		Decoder: Decode,
	})
}
