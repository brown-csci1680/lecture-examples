package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

// A struct to represent our messages
type GuessMessage struct {
	MessageType uint8
	Number      int32
}

const (
	MessageTypeGuess    = 0
	MessageTypeResponse = 1
	MessageTypeNewGame  = 2

	GuessMessageSize = 5
)

// In order to send our message out on the wire, we need to
// turn it into a byte stream
//
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

func RecvAll(conn net.Conn, buffer []byte, n int) (int, error) {

	totalBytesRead := 0
	toRead := n

	for toRead > 0 {
		bytesRead, err := conn.Read(buffer[totalBytesRead:])
		if err == io.EOF {
			log.Println("Connection closed")
			return totalBytesRead, io.EOF
		} else if err != nil {
			log.Fatalln("Write failed:  ", err)
		}

		totalBytesRead += bytesRead
		toRead -= bytesRead
	}

	return totalBytesRead, nil
}

func ReadGuessMessage(conn net.Conn) (GuessMessage, error) {
	// Our messages are all the same size--but what would happen if they weren't?

	buffer := make([]byte, GuessMessageSize)

	//bytesRead, err := conn.Read(buffer)
	bytesRead, err := RecvAll(conn, buffer, GuessMessageSize)

	log.Printf("Read %d bytes\n", bytesRead)

	if err == io.EOF {
		log.Println("Connection closed")
		return GuessMessage{}, io.EOF
	} else if err != nil {
		log.Fatalln("Read error", err)
	}

	msg := GuessMessage{MessageType: buffer[0],
		Number: int32(binary.BigEndian.Uint32(buffer[1:]))}

	return msg, nil
}
