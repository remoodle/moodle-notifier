package service

import (
	"context"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"github.com/remoodle/notifier/internal/app/notifier/repository"
)

// CourseService provides an interface for course-related business logic
type CourseService interface {
	GetRelativeCourses(ctx context.Context, token string) ([]models.Course, error)
	GetUserEnrolledCourses(ctx context.Context, token string, userID int) ([]models.Course, error)
}

type courseServiceImpl struct {
	repo repository.CourseRepository
}

// NewCourseService creates a new CourseService
func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseServiceImpl{repo}
}

func (s *courseServiceImpl) GetRelativeCourses(ctx context.Context, token string) ([]models.Course, error) {
	return s.repo.GetCourses(ctx, token)
}

func (s *courseServiceImpl) GetUserEnrolledCourses(ctx context.Context, token string, userID int) ([]models.Course, error) {
	return s.repo.GetUserCourses(ctx, token, userID)
}
