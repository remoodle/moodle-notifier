package models

// Grade represents the structure of a grade item for a course
type Grade struct {
	ID       int    `json:"id" bson:"grade_id"`
	Name     string `json:"name"`
	CourseID int    `json:"course_id" bson:"course_id"`
	UserID   int    `bson:"user_id"`
	Grade    string `json:"grade" bson:"grade"`
}

type GradeChange struct {
	MoodleID      int    `json:"moodle_id"`
	CourseID      int    `json:"course_id"`
	GradeID       int    `json:"grade_id"`
	PreviousGrade string `json:"previous_grade"`
	NewGrade      string `json:"new_grade"`
}
