package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"io"
	"net/http"
	"strings"
)

// GradesRepository defines the interface for grade data operations
type GradesRepository interface {
	GetCourseGrades(ctx context.Context, token string, userID, courseID int) ([]models.Grade, error)
}

type gradesRepositoryImpl struct {
	moodleURL string
}

// NewGradesRepository returns a new instance of GradesRepository
func NewGradesRepository(moodleURL string) GradesRepository {
	return &gradesRepositoryImpl{moodleURL}
}

func (r *gradesRepositoryImpl) GetCourseGrades(ctx context.Context, token string, userID, courseID int) ([]models.Grade, error) {
	var grades []models.Grade
	params := fmt.Sprintf("?wstoken=%s&wsfunction=gradereport_user_get_grade_items&moodlewsrestformat=json&userid=%d&courseid=%d", token, userID, courseID)
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

	// Assuming the response structure includes a usergrades array with grade items
	var response struct {
		UserGrades []struct {
			GradeItems []struct {
				ID                  int    `json:"id"`
				ItemName            string `json:"itemname"`
				PercentageFormatted string `json:"percentageformatted"`
			} `json:"gradeitems"`
		} `json:"usergrades"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Process the response to fit the Grade model
	for _, userGrade := range response.UserGrades {
		for _, item := range userGrade.GradeItems {

			grade := models.Grade{
				ID:       item.ID,
				UserID:   userID,
				Name:     strings.Trim(item.ItemName, ""),
				CourseID: courseID,
				Grade:    item.PercentageFormatted,
			}
			if grade.Name != "" {
				grades = append(grades, grade)
			}
		}
	}

	return grades, nil
}
