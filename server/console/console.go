package console

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	client_auth_config "github.com/longzekun/specter/client/config"
	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/db"
	"github.com/longzekun/specter/server/db/models"
	"gorm.io/gorm"
)

var namePattern = regexp.MustCompile("^[a-zA-Z0-9_-]*$")

func NewOperatorClientConfig(operatorName string, lhost string, lport uint16) ([]byte, error) {
	if !namePattern.MatchString(operatorName) {
		return nil, fmt.Errorf("invalid operator name: %s", operatorName)
	}
	if operatorName == "" {
		return nil, fmt.Errorf("invalid operator name: %s", operatorName)
	}

	//
	var record models.Operator
	err := db.Session().Where(&models.Operator{Name: operatorName}).First(&record).Error
	if err == nil {
		return nil, fmt.Errorf("operator already exists: %s", operatorName)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if lhost == "" {
		return nil, fmt.Errorf("invalid host name: %s", lhost)
	}

	if lport == 0 {
		return nil, fmt.Errorf("invalid port number: %d", lport)
	}

	rawToken := models.GenerateOperatorToken()
	digest := sha256.Sum256([]byte(rawToken))
	dbOperator := &models.Operator{
		Name:  operatorName,
		Token: hex.EncodeToString(digest[:]),
	}

	err = db.Session().Save(dbOperator).Error
	if err != nil {
		return nil, err
	}

	//	cert
	publicKey, privateKey, err := certs.OperatorClientGenerateCertificate(operatorName)
	if err != nil {
		return nil, err
	}
	caCertPem, _, _ := certs.GetCertificateAuthorityPem(certs.OperatorClientType)

	clientAuthConfig := &client_auth_config.ClientAuthConfig{
		Operator:      operatorName,
		Token:         rawToken,
		Lhost:         lhost,
		Lport:         int(lport),
		CACertificate: string(caCertPem),
		PrivateKey:    string(privateKey),
		PublicKey:     string(publicKey),
	}

	clientConfig, err := json.MarshalIndent(clientAuthConfig, "", "    ")
	if err != nil {
		return nil, err
	}

	return clientConfig, nil
}
