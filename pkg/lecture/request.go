package lecture

type AddLectureToUserRequest struct {
	Username    string `json:"username"`
	LectureCode string `json:"lectureCode"`
}
