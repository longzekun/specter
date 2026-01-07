package core

import "github.com/gofrs/uuid"

type Session struct {
	ID         string
	Name       string
	Hostname   string
	Username   string
	UUID       string
	UID        string
	GID        string
	OS         string
	Version    string
	Arch       string
	PID        int32
	Filename   string
	Connection *ImplantConnection
}

func NewSession(implantConnection *ImplantConnection) *Session {
	implantConnection.UpdateLastMessageTime()
	return &Session{
		ID:         nextSessionID(),
		Connection: implantConnection,
	}
}

func nextSessionID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
