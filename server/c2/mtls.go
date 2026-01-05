package c2

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/longzekun/specter/protobuf/specterpb"
	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/constants"
	"github.com/longzekun/specter/server/core"
	handlers2 "github.com/longzekun/specter/server/handlers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	ServerMaxMessageSize = 2 * 1024 * 1024 * 1024
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
	zap.S().Infof("Accept connection from %v", conn.RemoteAddr().String())
	implantConn := core.NewImplantConnection(constants.MtlsStr, conn.RemoteAddr().String())
	defer func() {
		zap.S().Debugf("connection from %v closed", conn.RemoteAddr().String())
		conn.Close()
		implantConn.CleanUp()
	}()

	done := make(chan bool)
	//	deal recv data
	go func() {
		defer func() {
			done <- true
		}()

		handlers := handlers2.GetHandlers()
		for {
			//	recv data from implant
			envelope, err := recvEnvelopeFromImplant(conn)
			if err != nil {
				zap.S().Errorf("Error receiving envelope from implant: %v", err)
				return
			}
			implantConn.UpdateLastMessageTime()
			if envelope.ID != 0 {
				implantConn.RespMutex.Lock()
				if resp, ok := implantConn.Resp[envelope.ID]; ok {
					resp <- envelope
				}
				implantConn.RespMutex.Unlock()
			} else if handler, ok := handlers[envelope.Type]; ok {
				// first connect or beacon connect
				go func() {
					respEnvelop := handler(implantConn, envelope.Data)
					if respEnvelop != nil {
						implantConn.Send <- respEnvelop
					}
				}()
			}
		}
	}()

	//	deal send data
Loop:
	for {
		select {
		case envelope := <-implantConn.Send:
			//	send data to implant
			err := sendEnvelopeToImplant(conn, envelope)
			if err != nil {
				break Loop
			}
		case <-done:
			break Loop
		}
	}
}

func sendEnvelopeToImplant(conn net.Conn, envelope *specterpb.Envelope) error {
	data, err := proto.Marshal(envelope)
	if err != nil {
		return err
	}

	dataLengthBefore := new(bytes.Buffer)
	binary.Write(dataLengthBefore, binary.LittleEndian, uint32(len(data)))
	conn.Write(dataLengthBefore.Bytes())
	conn.Write(data)
	return nil
}

func recvEnvelopeFromImplant(conn net.Conn) (*specterpb.Envelope, error) {
	dataLengthBefore := make([]byte, 4)
	n, err := io.ReadFull(conn, dataLengthBefore)
	if err != nil || n != 4 {
		return nil, err
	}

	dataLength := int(binary.LittleEndian.Uint32(dataLengthBefore))
	zap.S().Debugf("Receive envelope from implant: %v", dataLength)
	if dataLength < 0 || dataLength > ServerMaxMessageSize {
		return nil, fmt.Errorf("invalid data length")
	}

	dataBuf := make([]byte, dataLength)
	n, err = io.ReadFull(conn, dataBuf)
	if err != nil || n != dataLength {
		return nil, err
	}

	envelope := &specterpb.Envelope{}
	err = proto.Unmarshal(dataBuf, envelope)
	if err != nil {
		return nil, err
	}

	return envelope, nil
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
