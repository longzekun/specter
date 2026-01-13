package core

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/specterpb"
)

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

func (s *Session) ToProtobuf() *clientpb.Session {
	return &clientpb.Session{
		ID:              s.ID,
		Name:            s.Name,
		Hostname:        s.Hostname,
		Username:        s.Username,
		UUID:            s.UUID,
		UID:             s.UID,
		GID:             s.GID,
		OS:              s.OS,
		Version:         s.Version,
		Arch:            s.Arch,
		PID:             fmt.Sprintf("%d", s.PID),
		Filename:        s.Filename,
		RemoteAddress:   s.Connection.RemoteAddress,
		TransportType:   s.Connection.Transport,
		LastMessageTime: s.Connection.GetLastMessageTime().Unix(),
		Health:          true,
	}
}

func (s *Session) Request(msgType uint32, data []byte) {
	s.Connection.Send <- &specterpb.Envelope{
		ID:   0,
		Type: msgType,
		Data: data,
	}
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
