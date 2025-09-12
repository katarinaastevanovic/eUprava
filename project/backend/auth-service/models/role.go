package models

type Role string

const (
	Student Role = "STUDENT"
	Teacher Role = "TEACHER"
	Doctor  Role = "DOCTOR"
	Parent  Role = "PARENT"
)
