package server

import (
	"chatserver/commands"
	"chatserver/protocol"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
)

type ChatServer interface {
	Listen(address string) error
	Start()
	Close()
	Broadcast(interface{}) error
}

type TCPChatServer struct {
	secured   bool
	listener  net.Listener
	clients   []*protocol.UserSession
	mutex     *sync.Mutex
	tlsConfig *tls.Config
}

func NewChatServer(secure bool) ChatServer {
	return &TCPChatServer{
		secured: secure,
		clients: make([]*protocol.UserSession, 0),
		mutex:   &sync.Mutex{},
	}

}

func (s *TCPChatServer) Listen(address string) error {

	if s.secured {
		s.setUpTLSCertificate()
		l, err := tls.Listen("tcp", address, s.tlsConfig)
		if err == nil {
			s.listener = l
		}
		return err
	}

	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
	}

	return err
}

func (s *TCPChatServer) Close() {
	s.listener.Close()
}

func (s *TCPChatServer) Start() {
	for {
		// XXX: need a way to break the loop
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *TCPChatServer) accept(conn net.Conn) *protocol.UserSession {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	client := &protocol.UserSession{
		Conn:          conn,
		Writer:        protocol.NewCommandWriter(conn),
		Authenticated: false,
	}
	s.clients = append(s.clients, client)
	return client
}

func (s *TCPChatServer) remove(client *protocol.UserSession) {
	log.Printf("Removing client %s", client.Name)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	log.Printf("Closing connection from %v", client.Conn.RemoteAddr().String())
	client.Conn.Close()
}

func (s *TCPChatServer) serve(session *protocol.UserSession) {
	disconnectionRequested := false
	cmdReader := protocol.NewCommandReader(session.Conn)
	defer s.remove(session)
	for {
		if session.FailedAuthentication {
			break
		}

		if disconnectionRequested {
			break
		}

		cmd, err := cmdReader.Read(session)
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}
		if cmd != nil {
			switch v := cmd.(type) {

			case commands.AuthAttemptCommand:
				pass := v.Password
				user := session.Name
				log.Printf("Authenticating %s by password '%s'", user, pass)
				if pass != "pass" {
					log.Printf("Authentication failed for user %s", user)
					s.handleAuthFailure(session)
					break
				}

				s.handleSuccessfulAuth(session)

			case commands.AuthReqCommand:
				session.Name = v.User
				session.FailedAuthentication = false
				log.Printf("Authentication request by user %s", session.Name)

			case commands.MessageCommand:
				if !session.Authenticated {
					s.handleAuthFailure(session)
				}

				s.RegisterMessage(session, v.Message)
				go s.Broadcast(commands.MessageCommand{
					Message: v.Message,
					User:    session.Name,
				})

			case commands.DisconnectionCommand:
				log.Printf("User %s left", session.Name)
				disconnectionRequested = true

			}

		}
		if err == io.EOF {
			break
		}
	}
}

func (s *TCPChatServer) handleAuthFailure(client *protocol.UserSession) {
	client.Writer.Write(commands.AuthFailedCommand{})
	client.FailedAuthentication = true
}

func (s *TCPChatServer) handleSuccessfulAuth(client *protocol.UserSession) {
	client.Writer.Write(commands.AuthPassCommand{})
	client.Authenticated = true
}

func (s *TCPChatServer) RegisterMessage(client *protocol.UserSession, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	client.Messages = append(client.Messages, message)
}

func (s *TCPChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		client.Writer.Write(command)
	}
	return nil
}

func (s *TCPChatServer) setUpTLSCertificate() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	s.tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
}
