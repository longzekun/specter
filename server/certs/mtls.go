package certs

func MtlsC2ServerGenerateECCCertificate(host string) ([]byte, []byte, error) {
	certPEM, keyPEM, err := GenerateEccCertificate(MtlsServerType, host, false, false, false)
	if err != nil {
		return nil, nil, err
	}
	// save data to database
	err = saveCertificateToDB(MtlsServerType, ECCKey, host, certPEM, keyPEM)
	return certPEM, keyPEM, err
}

func MtlsC2ImplantGenerateECCCertificate(name string) ([]byte, []byte, error) {
	certPEM, keyPEM, err := GenerateEccCertificate(MtlsImplantType, name, false, true, false)
	if err != nil {
		return nil, nil, err
	}
	//	save implant
	err = saveCertificateToDB(MtlsImplantType, ECCKey, name, certPEM, keyPEM)
	return certPEM, keyPEM, err
}
