package certs

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/longzekun/specter/server/config"
	"github.com/longzekun/specter/server/db"
	"github.com/longzekun/specter/server/db/models"
)

const (
	ECCKey = "ecc"
)

func init() {
	certsDir := config.GetServerConfig().CertsDir
	if certsDir == "" {
		panic("certsDir is empty")
	}
	err := os.MkdirAll(certsDir, 0755)
	if err != nil {
		panic(err)
	}
}

func saveCertificateToDB(caType, keyType, commonName string, certPEM, keyPEM []byte) error {
	if keyType != ECCKey {
		return errors.New("only ECC key is supported")
	}

	certModel := models.Certificate{
		CommonName:     commonName,
		CAType:         caType,
		KeyType:        keyType,
		CertificatePEM: string(certPEM),
		PrivateKeyPEM:  string(keyPEM),
	}

	dbSession := db.Session()
	result := dbSession.Create(&certModel)
	return result.Error
}

func GetCertificateFromDB(caType, keyType, commonName string) ([]byte, []byte, error) {
	if keyType != ECCKey {
		return nil, nil, errors.New("only ECC key is supported")
	}

	certModel := models.Certificate{}

	result := db.Session().Where(&models.Certificate{
		CAType:     caType,
		KeyType:    keyType,
		CommonName: commonName,
	}).First(&certModel)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	return []byte(certModel.CertificatePEM), []byte(certModel.PrivateKeyPEM), nil
}

func GenerateEccCertificate(caType string, commonName string, isCA bool, isClient bool, isOperator bool) ([]byte, []byte, error) {
	var curves []elliptic.Curve
	if isOperator {
		curves = []elliptic.Curve{
			elliptic.P256(),
		}
	} else {
		curves = []elliptic.Curve{
			elliptic.P256(),
			elliptic.P521(),
			elliptic.P384(),
		}
	}

	curve := curves[randomInt(len(curves))]

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	subject := pkix.Name{
		CommonName: commonName,
	}

	notBefore := time.Now()
	days := randomInt(365) * -1 // Within -1 year
	notBefore = notBefore.AddDate(0, 0, days)
	notAfter := notBefore.Add(randomValidFor())

	var keyUsage x509.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	var extKeyUsage []x509.ExtKeyUsage

	if isCA {
		keyUsage = x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		extKeyUsage = []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}
	} else if isClient {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	} else {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: isCA,
	}

	if !isClient {
		if ip := net.ParseIP(subject.CommonName); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, subject.CommonName)
		}
	}

	var certErr error
	var derBytes []byte
	if isCA {
		template.IsCA = isCA
		template.KeyUsage |= x509.KeyUsageCertSign
		derBytes, certErr = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	} else {
		caCert, caKey, err := GetCertificateAuthority(caType)
		if err != nil {
			panic(err)
		}
		derBytes, certErr = x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	}

	if certErr != nil {
		return nil, nil, certErr
	}

	certOut := bytes.NewBuffer([]byte{})
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := bytes.NewBuffer([]byte{})
	b, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: b})

	return certOut.Bytes(), keyOut.Bytes(), nil
}

func randomInt(max int) int {
	buf := make([]byte, 4)
	rand.Read(buf)
	i := binary.LittleEndian.Uint32(buf)
	return int(i) % max
}

func randomValidFor() time.Duration {
	validFor := 3 * (365 * 24 * time.Hour)
	return validFor
}
