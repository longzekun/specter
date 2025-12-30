package c2

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"

	"encoding/pem"

	"github.com/longzekun/specter/server/certs"
)

func StartMutualListener(host string, port uint16) (net.Listener, error) {
	certPem, keyPem, err := certs.GenerateEccCertificate("mtls-server", host, true, false, true)
	if err != nil {
		return nil, err
	}

	certBlock, _ := pem.Decode(certPem)
	if certBlock == nil {
		return nil, fmt.Errorf("pem decode failed")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	mtlsCACertPool := x509.NewCertPool()
	mtlsCACertPool.AddCert(cert)

	mtlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: mtlsCACertPool,
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    mtlsCACertPool,
		Certificates: []tls.Certificate{mtlsCert},
		MinVersion:   tls.VersionTLS12,
	}

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
