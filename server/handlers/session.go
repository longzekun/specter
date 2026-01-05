package handlers

import (
	"github.com/longzekun/specter/protobuf/specterpb"
	"github.com/longzekun/specter/server/core"
	"go.uber.org/zap"
)

func registerSessionHandler(implantConnection *core.ImplantConnection, data []byte) *specterpb.Envelope {
	zap.S().Debugf("session register")
	return nil
}
