package protocol

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type GuessMessage struct {
	MessageType uint8
	Number      int32
}

const (
	MessageTypeGuess    = 0
	MessageTypeResponse = 1

	GuessMessageSize = 5
)

// A message is 5 bytes
// ID XX XX XX XX XX
// My number is 42 => 0x0000002a

// If you have a binary protocol, need to specify byte order
// Most network systems use big endian

// 0x00

// Method 1
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

	// Read from the socket, will block until there is SOME data
	bytesRead, err := conn.Read(buffer)

	log.Printf("Read %d bytes\n", bytesRead)

	if err != nil {
		log.Fatalln("Read error", err)
	}

	msg := GuessMessage{MessageType: buffer[0],
		Number: int32(binary.BigEndian.Uint32(buffer[1:]))}

	return msg
}
