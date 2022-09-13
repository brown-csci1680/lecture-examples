package protocol

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const (
	MessageTypeGuess    = 0
	MessageTypeResponse = 1
	MessageTypeNewGame  = 2
)

type GuessMessage struct {
	MessageType uint8
	Number      int32
}

func (m *GuessMessage) Marshal() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, m.MessageType)
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}

	err = binary.Write(buf, binary.BigEndian, m.Number)
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}

	return buf.Bytes()
}

func ReadGuessMessage(conn net.TCPConn) *GuessMessage {
	return nil
}
