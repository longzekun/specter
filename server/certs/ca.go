package certs

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/longzekun/specter/server/config"
	"go.uber.org/zap"
)

const (
	MtlsServerType  = "mtls-server"
	MtlsImplantType = "mtls-implant"
)

func getCertDir() string {
	certsDir := config.GetServerConfig().CertsDir
	if certsDir == "" {
		panic("certsDir is empty")
	}
	return certsDir
}

// SetupCAs
func SetupCAs() {
	//	mtls-server CA
	GenerateCertificateAuthority(MtlsServerType, "")
	//	mtls-implant CA
	GenerateCertificateAuthority(MtlsImplantType, "")
}

func GenerateCertificateAuthority(CAType string, commonName string) (*x509.Certificate, *ecdsa.PrivateKey) {
	certsDir := getCertDir()
	certFilePath := filepath.Join(certsDir, fmt.Sprintf("%s-ca-cert.pem", CAType))
	if _, err := os.Stat(certFilePath); os.IsNotExist(err) {
		zap.S().Infof("Generated CA certificate at %s", certFilePath)
		certBytes, keyBytes, err := GenerateEccCertificate(CAType, commonName, true, false, false)
		if err != nil {
			panic(err)
		}
		//	save cert and key to file
		SaveCertAndKeyToFile(CAType, certBytes, keyBytes)
	}
	cert, key, err := GetCertificateAuthority(CAType)
	if err != nil {
		panic(err)
	}
	return cert, key
}

func GetCertificateAuthority(CAType string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	certPem, keyPem, err := GetCertificateAuthorityPem(CAType)
	if err != nil {
		return nil, nil, err
	}

	certBlock, _ := pem.Decode(certPem)
	if certBlock == nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(keyPem)
	if keyBlock == nil {
		return nil, nil, err
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	key, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("key is not ECDSA")
	}

	return cert, key, nil
}

func GetCertificateAuthorityPem(CAType string) ([]byte, []byte, error) {
	certsDir := getCertDir()
	certFilePath := filepath.Join(certsDir, fmt.Sprintf("%s-ca-cert.pem", CAType))
	keyFilePath := filepath.Join(certsDir, fmt.Sprintf("%s-ca-key.pem", CAType))

	certPem, err := os.ReadFile(certFilePath)
	if err != nil {
		return nil, nil, err
	}

	keyPem, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, nil, err
	}

	return certPem, keyPem, nil
}

func SaveCertAndKeyToFile(CAType string, certBytes []byte, keyBytes []byte) {
	certDir := getCertDir()
	certFilePath := filepath.Join(certDir, fmt.Sprintf("%s-ca-cert.pem", CAType))
	keyFilePath := filepath.Join(certDir, fmt.Sprintf("%s-ca-key.pem", CAType))

	err := os.WriteFile(certFilePath, certBytes, 0600)
	if err != nil {
		zap.S().Warnf("Failed to save CA certificate at %s", certFilePath)
	}

	err = os.WriteFile(keyFilePath, keyBytes, 0600)
	if err != nil {
		zap.S().Warnf("Failed to save CA key at %s", keyFilePath)
	}
}
