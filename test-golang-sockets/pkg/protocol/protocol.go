package protocol

import (
	"bytes"
	"encoding/binary"
	"net"
)

// Struct to represent messages
type GuessMessage struct {
	MessageType uint8
	Number      int32
}

const (
	MessageTypeGuess    = 0
	MessageTypeResponse = 1

	AnswerTooLow  = -1
	AnswerCorrect = 0
	AnswerTooHigh = 1

	GuessMessageSize = 5
)

// So the idea is, when we want to send a message, we want to send/recv these structs
// You're going to learn to do this once, and then never again

// Want to do it this way, just so you see it
// Marshal struct into an array of bytes
// (m *GuessMessage) says "this is a function that operates on a GuessMessage called m"
func (m *GuessMessage) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Add the message type
	err := binary.Write(buf, binary.BigEndian, m.MessageType)
	if err != nil {
		return nil, err // Make errors the caller's problem
	}

	// Add the number
	err = binary.Write(buf, binary.BigEndian, m.Number)
	if err != nil {
		return nil, err
	}

	// Done, return the byte array (and no error)
	return buf.Bytes(), nil
}

func ReadGuessMessage(conn net.Conn) (*GuessMessage, error) {
	// All messages are how big?
	buffer := make([]byte, GuessMessageSize)

	_, err := conn.Read(buffer) // WARNING!
	if err != nil {
		return nil, err // TODO
	}

	msg := &GuessMessage{
		MessageType: buffer[0],
		Number:      int32(binary.BigEndian.Uint32(buffer[1:])),
	}

	return msg, nil
}
