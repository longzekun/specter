package certs

func OperatorClientGenerateCertificate(operator string) ([]byte, []byte, error) {
	cert, key, err := GenerateEccCertificate(OperatorClientType, operator, false, true, true)
	if err != nil {
		return nil, nil, err
	}
	err = saveCertificateToDB(OperatorClientType, ECCKey, operator, cert, key)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func OperatorServerGetCertificate(host string) ([]byte, []byte, error) {
	return GetCertificateFromDB(OperatorClientType, ECCKey, host)
}

func OperatorGenerateCertificate(host string) ([]byte, []byte, error) {
	cert, key, err := GenerateEccCertificate(OperatorClientType, host, false, false, true)
	if err != nil {
		return nil, nil, err
	}
	err = saveCertificateToDB(OperatorClientType, ECCKey, host, cert, key)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}
