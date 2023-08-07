package commands

import "fmt"

type AuthReqCommand struct {
	User string
}

type AuthAttemptCommand struct {
	User     string
	Password string
}

type AuthPassCommand struct {
}

type AuthFailedCommand struct {
}

type MessageCommand struct {
	User    string
	Message string
}

type DisconnectionCommand struct {
}

type UnknownCommand struct {
}

func (u UnknownCommand) Error() string {
	return fmt.Sprintf("Unknown command")
}
