package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/longzekun/specter/client/config"
	"github.com/longzekun/specter/protobuf/rpcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	kb                          = 1024
	mb                          = kb * 1024
	gb                          = mb * 1024
	ClientMaxReceiveMessageSize = (2 * gb) - 1

	defaultTimeout = time.Duration(10 * time.Second)
)

type TokenAuth struct {
	token string
}

func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.token,
	}, nil
}

func (TokenAuth) RequireTransportSecurity() bool {
	return true
}

func MtlsConnect() (rpcpb.SpecterRPCClient, *grpc.ClientConn, error) {
	clientAuthConfig := config.GetClientAuthConfig()
	if clientAuthConfig == nil {
		return nil, nil, fmt.Errorf("clientAuthConfig is nil")
	}

	if clientAuthConfig.Token == "" ||
		clientAuthConfig.Lhost == "" ||
		clientAuthConfig.PrivateKey == "" ||
		clientAuthConfig.PublicKey == "" {
		return nil, nil, fmt.Errorf("clientAuthConfig is nil")
	}
	//	use mtls connect to server,generate mtls config
	tlsConfig, err := getTlsConfig(clientAuthConfig.CACertificate, clientAuthConfig.PublicKey, clientAuthConfig.PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	transportCreds := credentials.NewTLS(tlsConfig)

	//token
	callCreds := credentials.PerRPCCredentials(TokenAuth{token: clientAuthConfig.Token})
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithPerRPCCredentials(callCreds),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", clientAuthConfig.Lhost, clientAuthConfig.Lport), options...)
	if err != nil {
		return nil, nil, err
	}
	return rpcpb.NewSpecterRPCClient(connection), connection, nil
}

func getTlsConfig(caCertificate string, certificate string, privateKey string) (*tls.Config, error) {
	certPem, err := tls.X509KeyPair([]byte(certificate), []byte(privateKey))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCertificate))
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certPem},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return RootOnlyVerifyCertificate(caCertificate, rawCerts)
		},
	}
	return tlsConfig, nil
}

func RootOnlyVerifyCertificate(caCertificate string, rawCerts [][]byte) error {
	return nil
}
