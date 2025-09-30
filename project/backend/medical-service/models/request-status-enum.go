package models

type TypeOfRequest string

const (
	REQUESTED TypeOfRequest = "REQUESTED"
	APPROVED  TypeOfRequest = "APPROVED"
	REJECTED  TypeOfRequest = "REJECTED"
	FINISHED  TypeOfRequest = "FINISHED"
)
