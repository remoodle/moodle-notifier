package models

// User represents the structure of our MongoDB documents
type User struct {
	ID                    string  `bson:"_id,omitempty"`
	MoodleID              int     `bson:"moodle_id"`
	TelegramID            int64   `bson:"telegram_id"`
	Username              string  `bson:"username"`
	HashedToken           string  `bson:"hashed_token"`
	FullName              string  `bson:"full_name"`
	Barcode               int     `bson:"barcode"`
	GradesNotification    bool    `bson:"grades_notification"`
	DeadlinesNotification int     `bson:"deadlines_notification"`
	IsAdmin               bool    `bson:"is_admin"`
	Grades                []Grade `bson:"grades"`
}
