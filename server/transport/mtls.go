package transport

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	kb                   = 1024
	mb                   = kb * 1024
	gb                   = mb * 1024
	ServerMaxMessageSize = 2 * gb
)

func StartMtlsClientListener(host string, port int) (*grpc.Server, net.Listener, error) {
	tlsConfig := getOperatorServerTlsConfig("multiplayer")

	serverCerts := credentials.NewTLS(tlsConfig)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, nil, err
	}

	options := []grpc.ServerOption{
		grpc.Creds(serverCerts),
		grpc.MaxRecvMsgSize(ServerMaxMessageSize), //	2GB
		grpc.MaxSendMsgSize(ServerMaxMessageSize), //	2GB
	}

	options = append(options, initMiddle()...)

	grpcServer := grpc.NewServer(options...)
	rpcpb.RegisterSpecterRPCServer(grpcServer, rpc.NewServer())
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			zap.S().Errorf("failed to serve: %v", err)
		}
	}()

	return grpcServer, listener, nil
}

func getOperatorServerTlsConfig(host string) *tls.Config {
	caCert, _, err := certs.GetCertificateAuthority(certs.OperatorClientType)
	if err != nil {
		zap.S().Fatalf("failed to get operator server tls config: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	//	server tls config
	_, _, err = certs.OperatorServerGetCertificate(host)
	if err != nil {
		certs.OperatorGenerateCertificate(host)
	}

	certPEM, keyPEM, err := certs.OperatorServerGetCertificate(host)
	if err != nil {
		zap.S().Errorf("failed to get operator server tls config: %v", err)
		return nil
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		zap.S().Fatalf("failed to load operator server tls config: %v", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	return tlsConfig
}
