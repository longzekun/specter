package transport

import (
	"fmt"
	"net"

	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	kb                   = 1024
	mb                   = kb * 1024
	gb                   = mb * 1024
	ServerMaxMessageSize = 2 * gb
)

func StartMtlsClientListener(host string, port int) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, nil, err
	}

	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(ServerMaxMessageSize), //	2GB
		grpc.MaxSendMsgSize(ServerMaxMessageSize), //	2GB
	}

	grpcServer := grpc.NewServer(options...)
	rpcpb.RegisterSpecterRPCServer(grpcServer, rpc.NewServer())
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			zap.S().Errorf("failed to serve: %v", err)
		}
	}()

	return grpcServer, listener, nil
}
