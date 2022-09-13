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

	GuessMessageSize = 5
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

func ReadGuessMessage(conn net.Conn) GuessMessage {
	// Our messages are all the same size--but what would happen if they weren't?

	buffer := make([]byte, GuessMessageSize)

	bytesRead, err := conn.Read(buffer)

	log.Printf("Read %d bytes\n", bytesRead)

	if err != nil {
		log.Fatalln("Read error", err)
	}

	msg := GuessMessage{MessageType: buffer[0],
		Number: int32(binary.BigEndian.Uint32(buffer[1:]))}

	return msg
}
