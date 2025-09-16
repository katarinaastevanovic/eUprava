package models

type TypeOfCertificate string

const (
	CERT_REGULAR   TypeOfCertificate = "REGULAR"
	CERT_PE        TypeOfCertificate = "PE"
	CERT_EXCURSION TypeOfCertificate = "EXCURSION"
	CERT_SICKNESS  TypeOfCertificate = "SICKNESS"
)
