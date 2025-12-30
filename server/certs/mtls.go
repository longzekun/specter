package certs

func MtlsC2ServerGenerateECCCertificate(host string) ([]byte, []byte, error) {
	certPEM, keyPEM, err := GenerateEccCertificate(MtlsServerType, host, false, false, false)
	if err != nil {
		return nil, nil, err
	}
	// save data to database
	err = saveCertifateToDB(MtlsServerType, ECCKey, host, certPEM, keyPEM)
	return certPEM, keyPEM, err
}
