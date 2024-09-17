package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
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

	GuessMessageSize = 5
)

// When we want to send a message, we want to send/recv these structs
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

// Helper for sending a guess message
func SendGuessMessage(conn net.Conn, guess int) error {
	// Build the guess message
	guessMsg := GuessMessage{
		MessageType: MessageTypeGuess,
		Number:      int32(guess),
	}

	// Turn it into a byte array
	toSend, err := guessMsg.Marshal()
	if err != nil {
		return err
	}

	// Send it
	bytesSent, err := conn.Write(toSend)
	if err != nil {
		return err
	}

	log.Printf("Sent %d bytes", bytesSent)

	return nil
}

// ALTERNATE VERSION:  Helper for building a guess message
// This version sends the message using multiple calls to conn.Write.
// This is just as reasonable, and may be required in some situations.
func SendGuessMessageV2(conn net.Conn, guess int) error {
	buf1 := new(bytes.Buffer)
	err := binary.Write(buf1, binary.BigEndian, uint8(MessageTypeGuess))
	if err != nil {
		return err
	}
	b, err := conn.Write(buf1.Bytes())
	if err != nil {
		return err
	}
	log.Printf("Sent %d bytes", b)

	buf2 := new(bytes.Buffer)
	err = binary.Write(buf2, binary.BigEndian, int32(guess))
	if err != nil {
		return err
	}
	b, err = conn.Write(buf2.Bytes())
	if err != nil {
		return err
	}
	log.Printf("Sent %d bytes", b)

	return nil
}

func ReadGuessMessage(conn net.Conn) (*GuessMessage, error) {
	// All messages are how big?
	buffer := make([]byte, GuessMessageSize)

	// Read from the socket
	// **** OLD VERSION ****
	// WARNING WARNING WARNING:  What happens if we receive fewer than 5 bytes??
	// We'll fix this later--see full version for details now.
	//_, err := conn.Read(buffer) // ignored parameter is number of bytes received
	// *******************

	// FIXED VERSION:  Read from socket (in a loop) until buffer is full (or error)
	// See lecture video for details
	b, err := io.ReadFull(conn, buffer)
	if err != nil {
		// TODO:  Graceful error handling when client disconnects normally
		return nil, err
	}
	log.Printf("%s:  Received %d bytes", conn.RemoteAddr().String(), b)

	msg := &GuessMessage{
		MessageType: buffer[0],
		Number:      int32(binary.BigEndian.Uint32(buffer[1:])),
	}

	return msg, nil
}
