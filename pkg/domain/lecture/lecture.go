package lecture

import (
	"context"
	"fmt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func GetLecturesFromUser(ctx context.Context, db *database.DB, id types.UserID) ([]gendb.GetUserLecturesRow, error) {
	lectures, err := db.GetUserLectures(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(lectures) == 0 {
		return nil, fmt.Errorf("lecture not found in the given user")
	}

	return lectures, nil
}

func AddLectureToUser(ctx context.Context, db *database.DB, username, lectureCode string) (gendb.AddLectureToUserRow, error) {
	user, err := db.GetUserCredentialsWithUsername(ctx, username)
	if err != nil {
		return gendb.AddLectureToUserRow{}, err
	}

	lecture, err := db.GetLectureByCode(ctx, lectureCode)
	if err != nil {
		return gendb.AddLectureToUserRow{}, err
	}

	row, err := db.AddLectureToUser(ctx, gendb.AddLectureToUserParams{
		UserID:    user.ID,
		LectureID: lecture.ID,
	})
	if err != nil {
		return gendb.AddLectureToUserRow{}, err
	}

	return row, nil
}

func RemoveLectureFromUser(ctx context.Context, db *database.DB, uid types.UserID, lid types.LectureID) error {
	row, err := db.RemoveLectureFromUser(ctx, gendb.RemoveLectureFromUserParams{
		UserID:    uid,
		LectureID: lid,
	})
	if err != nil {
		return err
	}
	if row <= 0 {
		return fmt.Errorf("given lecture not found in the user")
	}

	return nil
}
