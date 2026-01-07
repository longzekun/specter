package handlers

import (
	"strconv"

	"github.com/longzekun/specter/protobuf/specterpb"
	"github.com/longzekun/specter/server/core"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func registerSessionHandler(implantConnection *core.ImplantConnection, data []byte) *specterpb.Envelope {
	zap.S().Debugf("session register")

	register := specterpb.MsgRegister{}
	err := proto.Unmarshal(data, &register)
	if err != nil {
		return nil
	}

	session := core.NewSession(implantConnection)
	session.Name = register.Name
	session.Hostname = register.Hostname
	session.Username = register.Username
	session.UUID = register.UUID
	session.UID = register.UID
	session.GID = register.GID
	session.OS = register.OS
	session.Version = register.Version
	session.Arch = register.Arch
	pid, _ := strconv.Atoi(register.PID)
	session.PID = int32(pid)
	session.Filename = register.Filename

	core.Sessions.Add(session)
	implantConnection.CleanUp = func() {
		core.Sessions.Remove(session)
	}
	return nil
}
