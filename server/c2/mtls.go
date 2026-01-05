package c2

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"

	"github.com/longzekun/specter/server/certs"
	"go.uber.org/zap"
)

func StartMutualListener(host string, port uint16) (net.Listener, error) {
	_, _, err := certs.GetCertificateFromDB(certs.MtlsServerType, certs.ECCKey, host)
	if err != nil {
		certs.MtlsC2ServerGenerateECCCertificate(host)
	}

	//	get tls config
	tlsConfig := getTlsConfig(host)

	ln, err := tls.Listen("tcp", fmt.Sprintf("%v:%v", host, port), tlsConfig)
	if err != nil {
		return nil, err
	}

	go acceptImplantConnections(ln)

	return ln, nil
}

func acceptImplantConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errType, ok := err.(*net.OpError); ok && errType.Op == "accept" {
				break // Listener was closed by the user
			}

			continue
		}

		go handleImplantConnections(conn)
	}
}

func handleImplantConnections(conn net.Conn) {

	fmt.Println("新的连接从:%v", conn.RemoteAddr().String())
	//	处理接收到的客户端的数据
	dataLengthBefore := make([]byte, 4)
	io.ReadFull(conn, dataLengthBefore)

	fmt.Println(dataLengthBefore)
	//	向客户端发送数据
}

func getTlsConfig(host string) *tls.Config {
	implantCACert, _, err := certs.GetCertificateAuthority(certs.MtlsImplantType)
	if err != nil {
		zap.S().Fatalf("get implant CA cert failed: %v", err)
	}
	implantCACertPool := x509.NewCertPool()
	implantCACertPool.AddCert(implantCACert)

	certPEM, keyPEM, err := certs.GetCertificateFromDB(certs.MtlsServerType, certs.ECCKey, host)
	if err != nil {
		zap.S().Fatalf("get certificate from DB failed: %v", err)
		return nil
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		zap.S().Fatalf("load X509KeyPair failed: %v", err)
	}

	tlsConfig := &tls.Config{
		RootCAs: implantCACertPool,
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    implantCACertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig
}
