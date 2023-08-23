package types

type (
	UserID    uint64
	RoleID    uint64
	LectureID uint64
	Role      string
)

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
)
