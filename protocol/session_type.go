package protocol

import (
	"net"
)

type UserSession struct {
	Name                 string
	Password             string
	Messages             []string
	Writer               *CommandWriter
	Conn                 net.Conn
	Authenticated        bool
	FailedAuthentication bool
}
