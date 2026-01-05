package handlers

import (
	"github.com/longzekun/specter/protobuf/specterpb"
	"github.com/longzekun/specter/server/core"
)

type ServerHandler func(*core.ImplantConnection, []byte) *specterpb.Envelope

func GetHandlers() map[uint32]ServerHandler {
	return map[uint32]ServerHandler{
		specterpb.SessionRegister: registerSessionHandler,
	}
}
