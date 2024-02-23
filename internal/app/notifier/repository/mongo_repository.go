package repository

import (
	"context"
	"fmt"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// UserRepository interface to include FindAll method
type UserRepository interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	UpdateGrades(ctx context.Context, moodleID int, grades []models.Grade) error
	FindByMoodleID(ctx context.Context, moodleID int) (*models.User, error)
	RemoveAllGrades(ctx context.Context) error
}

type userRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserRepository returns a new instance of UserRepository
func NewUserRepository(client *mongo.Client) UserRepository {
	db := client.Database("moodbot") // Corrected database name
	return &userRepositoryImpl{
		collection: db.Collection("user"), // Corrected collection name
	}
}

// FindAll method to fetch all users
func (r *userRepositoryImpl) FindAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	cursor, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepositoryImpl) FindByMoodleID(ctx context.Context, moodleID int) (*models.User, error) {
	var user models.User
	filter := bson.M{"moodle_id": moodleID}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Error finding user by MoodleID %d: %v", moodleID, err)
		return nil, err
	}
	return &user, nil
}

//func (r *userRepositoryImpl) UpdateGrades(ctx context.Context, moodleID int, newGrades []models.Grade) error {
//	filter := bson.M{"moodle_id": moodleID}
//	update := bson.M{"$set": bson.M{"grades": newGrades}}
//	_, err := r.collection.UpdateOne(ctx, filter, update)
//	if err != nil {
//		log.Printf("Error updating grades for MoodleID %d: %v", moodleID, err)
//		return err
//	}
//	return nil
//}

func (r *userRepositoryImpl) UpdateGrades(ctx context.Context, moodleID int, newGrades []models.Grade) error {
	var user models.User
	if err := r.collection.FindOne(ctx, bson.M{"moodle_id": moodleID}).Decode(&user); err != nil {
		log.Printf("Error finding user by MoodleID %d: %v", moodleID, err)
		return err
	}

	// Map existing grades for quick lookup
	existingGradesMap := make(map[string]models.Grade)
	for _, grade := range user.Grades {
		key := fmt.Sprintf("%d-%d", grade.CourseID, grade.ID)
		existingGradesMap[key] = grade
	}

	// Prepare the final grades array including updates and new grades
	finalGrades := user.Grades // Start with existing grades
	for _, newGrade := range newGrades {
		key := fmt.Sprintf("%d-%d", newGrade.CourseID, newGrade.ID)
		if _, exists := existingGradesMap[key]; exists {
			// Replace the existing grade with the new one
			for i, grade := range finalGrades {
				if grade.CourseID == newGrade.CourseID && grade.ID == newGrade.ID {
					finalGrades[i] = newGrade
					break
				}
			}
		} else {
			// Append the new grade
			finalGrades = append(finalGrades, newGrade)
		}
	}

	// Update the user document with the final grades array
	_, err := r.collection.UpdateOne(ctx, bson.M{"moodle_id": moodleID}, bson.M{"$set": bson.M{"grades": finalGrades}})
	if err != nil {
		log.Printf("Error updating grades for MoodleID %d: %v", moodleID, err)
		return err
	}

	return nil
}

func (r *userRepositoryImpl) RemoveAllGrades(ctx context.Context) error {
	update := bson.M{"$unset": bson.M{"grades": ""}} // Remove the grades field
	_, err := r.collection.UpdateMany(ctx, bson.D{{}}, update)
	if err != nil {
		log.Printf("Error removing grades from all users: %v", err)
		return err
	}
	return nil
}
