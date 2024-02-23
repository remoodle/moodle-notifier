package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"io"
	"net/http"
	"time"
)

// CourseRepository defines the interface for course data operations
type CourseRepository interface {
	GetCourses(ctx context.Context, token string) ([]models.Course, error)
	GetUserCourses(ctx context.Context, token string, userID int) ([]models.Course, error)
}

type courseRepositoryImpl struct {
	moodleURL string
}

// NewCourseRepository returns a new instance of CourseRepository
func NewCourseRepository(moodleURL string) CourseRepository {
	return &courseRepositoryImpl{moodleURL}
}

func (r *courseRepositoryImpl) GetCourses(ctx context.Context, token string) ([]models.Course, error) {
	url := fmt.Sprintf("%s?wstoken=%s&wsfunction=core_course_get_courses_by_field&moodlewsrestformat=json", r.moodleURL, token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Courses []models.Course `json:"courses"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Filter courses based on end_date > current time
	currentTimestamp := time.Now().Unix()
	var filteredCourses []models.Course
	for _, course := range response.Courses {
		if course.EndDate > currentTimestamp {
			filteredCourses = append(filteredCourses, course)
		}
	}

	return filteredCourses, nil
}

func (r *courseRepositoryImpl) GetUserCourses(ctx context.Context, token string, userID int) ([]models.Course, error) {
	var courses []models.Course
	params := fmt.Sprintf("?wstoken=%s&wsfunction=core_enrol_get_users_courses&moodlewsrestformat=json&userid=%d", token, userID)
	url := r.moodleURL + params

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Assuming the response structure is an array of courses
	if err := json.Unmarshal(body, &courses); err != nil {
		// Handle possible "invalid token" error in the response
		var errResp struct {
			ErrorCode string `json:"errorcode"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.ErrorCode == "invalidtoken" {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	currentTimestamp := time.Now().Unix()
	var filteredCourses []models.Course
	for _, course := range courses {
		if course.EndDate > currentTimestamp {
			filteredCourses = append(filteredCourses, course)
		}
	}

	return filteredCourses, nil
}
