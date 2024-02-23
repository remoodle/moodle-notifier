package rabbitmq

import (
	"encoding/json"
	"github.com/remoodle/notifier/internal/app/notifier/models"
	"github.com/streadway/amqp"
	"log"
)

const QueueName = "gradeChanges"

func PublishGradeChange(ch *amqp.Channel, gradeChange models.GradeChange) {
	body, err := json.Marshal(gradeChange)
	if err != nil {
		log.Fatalf("Failed to marshal grade change: %v", err)
	}

	err = ch.Publish(
		"",        // exchange
		QueueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
