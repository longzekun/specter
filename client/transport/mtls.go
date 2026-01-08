package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/longzekun/specter/client/config"
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
	clientConfig := config.GetServerConfig()
	//	use mtls connect to server,generate mtls config
	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", clientConfig.ServerHost, clientConfig.ServerPort), options...)
	if err != nil {
		return nil, nil, err
	}
	return rpcpb.NewSpecterRPCClient(connection), connection, nil
}
