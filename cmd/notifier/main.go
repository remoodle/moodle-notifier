package main

import (
	"context"
	"fmt"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"github.com/remoodle/notifier/internal/app/notifier/repository"
	"github.com/remoodle/notifier/internal/app/notifier/service"
	"github.com/remoodle/notifier/internal/app/rabbitmq"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
)

const (
	keyString          = ""
	moodleURL          = ""
	mongoConnection    = ""
	rabbitMQConnection = ""
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoConnection))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := amqp.Dial(rabbitMQConnection)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ch)

	// Ensure the queue exists
	_, err = ch.QueueDeclare(
		rabbitmq.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	userRepo := repository.NewUserRepository(client)
	userService := service.NewUserService(userRepo)

	cryptoRepo := repository.NewCryptoRepository(keyString)
	cryptoService := service.NewCryptoService(cryptoRepo)

	courseRepo := repository.NewCourseRepository(moodleURL)
	courseService := service.NewCourseService(courseRepo)

	gradesRepo := repository.NewGradesRepository(moodleURL)
	gradesService := service.NewGradesService(gradesRepo)

	users, err := userService.GetAllUsers(context.TODO())
	if err != nil {
		log.Fatal("Failed to fetch users:", err)
	}
	//err = userService.RemoveAllGrades(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//return

	processUsers(ch, users, cryptoService, courseService, gradesService, userService)

}

func processUsers(ch *amqp.Channel, users []*models.User, cryptoService service.CryptoService, courseService service.CourseService, gradesService service.GradesService, userService service.UserService) {

	sort.Slice(users, func(i, j int) bool {
		return users[i].MoodleID < users[j].MoodleID
	})
	for _, user := range users {
		decryptedToken, err := cryptoService.DecryptToken(user.HashedToken)
		if err != nil {
			log.Fatal("Failed to decrypt token:", err)
		}

		courses, err := courseService.GetUserEnrolledCourses(context.Background(), decryptedToken, user.MoodleID)
		if err != nil {
			log.Fatalf("Failed to fetch user's enrolled courses: %v", err)
		}

		if len(courses) == 0 {
			fmt.Println("No courses found for the user.")
		} else {

			sort.Slice(courses, func(i, j int) bool {
				return courses[i].ID < courses[j].ID
			})

			for _, course := range courses {

				grades, err := gradesService.GetCourseGrades(context.Background(), decryptedToken, user.MoodleID, course.ID)
				if err != nil {
					log.Fatalf("Failed to fetch course grades: %v", err)
				}

				if len(grades) == 0 {
					fmt.Println("No grades found for the specified course and user.")
				} else {
					changedGrades, err := userService.UpdateUserGrades(context.Background(), user.MoodleID, user.Grades, grades)
					if err != nil {
						log.Fatalf("Failed to update grades for user %d: %v\n", user.MoodleID, err)
					}

					for _, grade := range changedGrades {
						log.Printf("Grade updated: User %d, Course %d, Grade ID: %d, Previous Grade: %s, New Grade: %s\n",
							user.MoodleID, grade.CourseID, grade.GradeID, grade.PreviousGrade, grade.NewGrade)

						rabbitmq.PublishGradeChange(ch, grade)
					}
				}

			}
		}
	}
}
