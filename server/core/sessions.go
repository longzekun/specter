package core

import (
	"sync"
)

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
	s.sessions.Delete(session.ID)

	//	push session quit message to all client
}

func (s *sessions) Get(sessionID string) *Session {
	v, ok := s.sessions.Load(sessionID)
	if !ok {
		return nil
	}
	return v.(*Session)
}

func (s *sessions) GetAllSessions() []*Session {
	all := []*Session{}
	s.sessions.Range(func(_, v any) bool {
		all = append(all, v.(*Session))
		return true
	})
	return all
}
