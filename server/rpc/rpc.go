package rpc

import (
	"github.com/longzekun/specter/protobuf/rpcpb"
)

type Server struct {
	rpcpb.UnimplementedSpecterRPCServer
}

func NewServer() *Server {
	return &Server{}
}
