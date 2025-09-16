package models

type AbsenceType string

const (
	Excused   AbsenceType = "excused"
	Unexcused AbsenceType = "unexcused"
	Pending   AbsenceType = "pending"
)
