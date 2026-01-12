package transport

import (
	"runtime/debug"

	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 2 * mb

func LocalListener() (*grpc.Server, *bufconn.Listener, error) {
	zap.S().Infof("starting local listener")
	listener := bufconn.Listen(bufSize)
	options := []grpc.ServerOption{
		grpc.MaxSendMsgSize(ServerMaxMessageSize),
		grpc.MaxRecvMsgSize(ServerMaxMessageSize),
	}

	options = append(options, initMiddle(true)...)
	grpcServer := grpc.NewServer(options...)
	rpcpb.RegisterSpecterRPCServer(grpcServer, rpc.NewServer())
	go func() {
		panicked := true
		defer func() {
			if panicked {
				zap.S().Errorf("panicked while starting local listener: %s", string(debug.Stack()))
			}
		}()
		if err := grpcServer.Serve(listener); err != nil {
			zap.S().Fatalf("error starting local listener: %s", err)
		} else {
			panicked = false
		}
	}()
	return grpcServer, listener, nil
}
