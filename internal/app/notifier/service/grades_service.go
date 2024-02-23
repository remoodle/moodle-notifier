package service

import (
	"context"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"github.com/remoodle/notifier/internal/app/notifier/repository"
)

// GradesService provides an interface for grade-related business logic
type GradesService interface {
	GetCourseGrades(ctx context.Context, token string, userID, courseID int) ([]models.Grade, error)
}

type gradesServiceImpl struct {
	repo repository.GradesRepository
}

// NewGradesService creates a new GradesService
func NewGradesService(repo repository.GradesRepository) GradesService {
	return &gradesServiceImpl{repo}
}

func (s *gradesServiceImpl) GetCourseGrades(ctx context.Context, token string, userID, courseID int) ([]models.Grade, error) {
	return s.repo.GetCourseGrades(ctx, token, userID, courseID)
}
