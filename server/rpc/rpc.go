package rpc

import (
	"context"

	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/protobuf/rpcpb"
)

type GeneralRequest interface {
}

type GeneralResponse interface {
}

type Server struct {
	rpcpb.UnimplementedSpecterRPCServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetVersion(ctx context.Context, in *commonpb.Empty) (*clientpb.Version, error) {
	version := &clientpb.Version{}

	return version, nil
}
