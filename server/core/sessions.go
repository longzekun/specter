package core

import "sync"

var (
	Sessions = sessions{
		sessions: &sync.Map{},
	}
)

type sessions struct {
	sessions *sync.Map
}

func (s *sessions) Add(session *Session) {
	s.sessions.Store(session.ID, session)

	//	push session to client
}

func (s *sessions) Remove(session *Session) {
}
