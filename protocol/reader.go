package protocol

import (
	"bufio"
	"chatserver/commands"
	"io"
	"log"
	"strconv"
	"strings"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read(session *UserSession) (interface{}, error) {
	bytes := make([]byte, 256)
	n, err := r.reader.Read(bytes)
	if err != nil {
		return nil, err
	}
	var message string
	if string(bytes[n-2:]) == "\r\n" {
		message = string(bytes[:n-3])
	} else {
		message = string(bytes[:n-2])
	}

	commandType, err := strconv.Atoi(message[:1])
	if err != nil {
		return nil, err
	}

	messageParts := strings.Split(message, " ")
	switch commandType {

	case commands.MESSAGE:
		message := messageParts[1:]

		return commands.MessageCommand{
			User:    session.Name,
			Message: strings.Join(message, " "),
		}, nil

	case commands.AUTH_REQUEST:
		user := messageParts[1]

		return commands.AuthReqCommand{
			User: user,
		}, nil

	case commands.AUTH_ATTEMPT:
		user := messageParts[1]
		pass := messageParts[2]

		return commands.AuthAttemptCommand{
			User:     user,
			Password: pass,
		}, nil

	case commands.DISCONNECTION:
		return commands.DisconnectionCommand{}, nil

	default:
		log.Printf("Unknown command: %v", commandType)
	}
	return nil, commands.UnknownCommand{}
}
