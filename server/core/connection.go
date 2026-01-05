package core

import (
	"sync"
	"time"
)

type ImplantConnection struct {
	ID                   string
	Send                 chan []byte
	RespMutex            *sync.Mutex
	Resp                 map[int64]chan []byte
	Transport            string
	RemoteAddress        string
	LastMessageTime      time.Time
	LastMessageTimeMutex *sync.Mutex
	CleanUp              func()
}
