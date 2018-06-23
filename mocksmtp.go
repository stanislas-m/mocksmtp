package Server

import (
	"errors"
	"sync"

	"github.com/gobuffalo/buffalo/mail"
)

// ErrNoMessage is returned when no message was caught.
var ErrNoMessage = errors.New("no message sent")

// Server is an in-memory implementation for buffalo `Sender`
// interface. It's intended to catch sent messages for test purposes.
type Server struct {
	messages []mail.Message
	mutex    sync.RWMutex
}

// Send implements buffalo `Sender` interface, to send mails.
func (s *Server) Send(m mail.Message) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messages = append(s.messages, m)
	return nil
}

// LastMessage gets the last sent message, if it exists.
// It returns `NoMessage` error if there is not last message.
func (s *Server) LastMessage() (mail.Message, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	l := len(s.Messages())
	if l == 0 {
		return mail.Message{}, ErrNoMessage
	}

	return s.Messages()[l-1], nil
}

// Messages gets the list of sent messages.
func (s *Server) Messages() []mail.Message {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.messages
}

// Count gets the amount of sent messages.
func (s *Server) Count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.messages)
}

// Clear destroys all sent messages.
func (s *Server) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = []mail.Message{}
}

// New constructs a new Server.
func New() *Server {
	return &Server{
		messages: []mail.Message{},
		mutex:    sync.RWMutex{},
	}
}
