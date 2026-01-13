package rpc

import (
	"context"
	"fmt"

	"github.com/longzekun/specter/client/constants"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/server/core"
)

func (s *Server) GetAllSessions(ctx context.Context, _ *commonpb.Empty) (*clientpb.Sessions, error) {
	retSessions := &clientpb.Sessions{}

	sessions := core.Sessions.GetAllSessions()
	for _, session := range sessions {
		retSessions.Sessions = append(retSessions.Sessions, session.ToProtobuf())
	}
	return retSessions, nil
}

func (s *Server) KillSession(ctx context.Context, req *clientpb.KillReq) (*commonpb.Empty, error) {
	session := core.Sessions.Get(req.SessionID)
	if session == nil {
		return &commonpb.Empty{}, fmt.Errorf("session not found")
	}
	core.Sessions.Remove(session)

	session.Request(constants.SessionQuit, []byte(""))

	return &commonpb.Empty{}, nil
}
