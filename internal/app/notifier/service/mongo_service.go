package service

import (
	"context"
	"fmt"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"github.com/remoodle/notifier/internal/app/notifier/repository"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserGrades(ctx context.Context, moodleID int, oldGrades []models.Grade, newGrades []models.Grade) ([]models.GradeChange, error)
	RemoveAllGrades(ctx context.Context) error
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo}
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userServiceImpl) RemoveAllGrades(ctx context.Context) error {
	return s.repo.RemoveAllGrades(ctx)
}

func (s *userServiceImpl) UpdateUserGrades(ctx context.Context, moodleID int, oldGrades []models.Grade, newGrades []models.Grade) ([]models.GradeChange, error) {

	changedGrades := compareGrades(oldGrades, newGrades)

	if len(changedGrades) > 0 {

		err := s.repo.UpdateGrades(ctx, moodleID, newGrades)
		if err != nil {
			return nil, err
		}
	}

	return changedGrades, nil
}

func compareGrades(existingGrades, newGrades []models.Grade) []models.GradeChange {
	var changes []models.GradeChange

	existingMap := make(map[string]models.Grade)
	for _, grade := range existingGrades {
		key := fmt.Sprintf("%d-%d", grade.CourseID, grade.ID)
		existingMap[key] = grade
	}

	for _, newGrade := range newGrades {
		key := fmt.Sprintf("%d-%d", newGrade.CourseID, newGrade.ID)

		if existingGrade, exists := existingMap[key]; exists {

			if existingGrade.Grade != newGrade.Grade {
				changes = append(changes, models.GradeChange{
					MoodleID:      newGrade.UserID,
					CourseID:      newGrade.CourseID,
					GradeID:       newGrade.ID,
					PreviousGrade: existingGrade.Grade,
					NewGrade:      newGrade.Grade,
				})
			}
		} else {
			// This is a new grade, not found in existing grades.
			changes = append(changes, models.GradeChange{
				MoodleID:      newGrade.UserID,
				CourseID:      newGrade.CourseID,
				GradeID:       newGrade.ID,
				PreviousGrade: "", // No previous grade since it's new.
				NewGrade:      newGrade.Grade,
			})

		}
	}

	return changes
}
