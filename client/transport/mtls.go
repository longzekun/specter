package transport

import (
	"context"
	"time"

	"github.com/longzekun/specter/protobuf/rpcpb"
	"google.golang.org/grpc"
)

const (
	kb                          = 1024
	mb                          = kb * 1024
	gb                          = mb * 1024
	ClientMaxReceiveMessageSize = (2 * gb) - 1

	defaultTimeout = time.Duration(10 * time.Second)
)

func MtlsConnect() (rpcpb.SpecterRPCClient, *grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	connection, err := grpc.DialContext(ctx, "127.0.0.1:7777", options...)
	if err != nil {
		return nil, nil, err
	}
	return rpcpb.NewSpecterRPCClient(connection), connection, nil
}
