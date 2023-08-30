package lecture

type AddLectureToUserRequest struct {
	Username    string `json:"username" validate:"required"`
	LectureCode string `json:"lectureCode" validate:"required,max=6"`
}
