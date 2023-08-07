package protocol

import (
	"chatserver/commands"
	"fmt"
	"io"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer,
	}
}

func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	// naive implementation ...
	var err error
	switch v := command.(type) {
	case commands.AuthPassCommand:
		err = w.writeString(fmt.Sprintf("Authentication passed\n"))
	case commands.AuthFailedCommand:
		err = w.writeString(fmt.Sprintf("Authentication failed\n"))
	case commands.MessageCommand:
		err = w.writeString(fmt.Sprintf("%s: %s\n", v.User, v.Message))
	default:
		err = commands.UnknownCommand{}
	}
	return err
}
