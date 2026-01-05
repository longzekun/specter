package core

import (
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/longzekun/specter/protobuf/specterpb"
)

type ImplantConnection struct {
	ID                   string
	Send                 chan *specterpb.Envelope
	RespMutex            *sync.Mutex
	Resp                 map[int64]chan *specterpb.Envelope
	Transport            string
	RemoteAddress        string
	LastMessageTime      time.Time
	LastMessageTimeMutex *sync.Mutex
	CleanUp              func()
}

func NewImplantConnection(transport string, remoteAddress string) *ImplantConnection {
	return &ImplantConnection{
		ID:                   generateImplantConnectionID(),
		Send:                 make(chan *specterpb.Envelope),
		RespMutex:            &sync.Mutex{},
		Resp:                 make(map[int64]chan *specterpb.Envelope),
		Transport:            transport,
		RemoteAddress:        remoteAddress,
		LastMessageTimeMutex: &sync.Mutex{},
		CleanUp:              func() {},
	}
}

func (conn *ImplantConnection) UpdateLastMessageTime() {
	conn.LastMessageTimeMutex.Lock()
	defer conn.LastMessageTimeMutex.Unlock()
	conn.LastMessageTime = time.Now()
}
func (conn *ImplantConnection) GetLastMessageTime() time.Time {
	conn.LastMessageTimeMutex.Lock()
	defer conn.LastMessageTimeMutex.Unlock()
	return conn.LastMessageTime
}

func generateImplantConnectionID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
