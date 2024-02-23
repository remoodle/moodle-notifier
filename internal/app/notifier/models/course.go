package models

// Course represents the structure of a course fetched from Moodle
type Course struct {
	ID                int    `json:"id"`
	Name              string `json:"displayname"`
	EnrolledUserCount int    `json:"enrolled_user_count"`
	Category          int    `json:"category"`
	Completed         bool   `json:"completed"`
	StartDate         int64  `json:"startdate"`
	EndDate           int64  `json:"enddate"`
}
