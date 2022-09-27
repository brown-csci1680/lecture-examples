package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"time"
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

func SendGuess(num int, conn net.Conn) {
	buf1 := new(bytes.Buffer)
	err := binary.Write(buf1, binary.BigEndian, uint8(MessageTypeGuess))
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}
	_, err = conn.Write(buf1.Bytes())
	if err != nil {
		log.Fatalln("Write", err)
	}

	buf2 := new(bytes.Buffer)
	err = binary.Write(buf2, binary.BigEndian, int32(num))
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}
	_, err = conn.Write(buf2.Bytes())
	if err != nil {
		log.Fatalln("Write", err)
	}
}

func RecvAll(conn net.Conn, buffer []byte, n int, timeout bool) (int, error) {
	if timeout {
		// Let's say we want to optionally have this read timeout if nothing was received
		// within 5 seconds
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	}
	bytesRead, err := io.ReadFull(conn, buffer)
	if timeout {
		// Remove the timeout deadline so that other reads
		// on this socket aren't affected
		conn.SetReadDeadline(time.Time{})
	}

	// Check for various error conditions
	// We might want to handle these in some way here; however, it might
	// be better to pass them up to the caller to handle them in a more graceful way
	// (ie, like stopping other parts of the application)
	if err == io.EOF {
		log.Println("Connection closed")
		return bytesRead, io.EOF
	} else if os.IsTimeout(err) {
		// Could handle timeouts differently here, if desired
		return bytesRead, err
	} else if err != nil {
		return bytesRead, err
	} else { // No error
		return bytesRead, nil
	}
}

func ReadGuessMessage(conn net.Conn, timeout bool) (GuessMessage, error) {
	// Our messages are all the same size--but what would happen if they weren't?

	buffer := make([]byte, GuessMessageSize)

	//bytesRead, err := conn.Read(buffer)
	bytesRead, err := RecvAll(conn, buffer, GuessMessageSize, timeout)

	log.Printf("Read %d bytes\n", bytesRead)

	if err == io.EOF {
		log.Println("Connection closed")
		return GuessMessage{}, io.EOF
	} else if err != nil {
		return GuessMessage{}, err
	}

	msg := GuessMessage{MessageType: buffer[0],
		Number: int32(binary.BigEndian.Uint32(buffer[1:]))}

	return msg, nil
}
